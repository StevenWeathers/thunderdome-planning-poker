<script>
    import Sockette from 'sockette'
    import { onDestroy, onMount } from 'svelte'
    import {
        dndzone,
        SHADOW_ITEM_MARKER_PROPERTY_NAME,
    } from 'svelte-dnd-action'

    import AddGoal from '../components/storyboard/AddGoal.svelte'
    import PageLayout from '../components/PageLayout.svelte'
    import UserCard from '../components/storyboard/UserCard.svelte'
    import InviteUser from '../components/storyboard/InviteUser.svelte'
    import ColumnForm from '../components/storyboard/ColumnForm.svelte'
    import StoryForm from '../components/storyboard/StoryForm.svelte'
    import ColorLegendForm from '../components/storyboard/ColorLegendForm.svelte'
    import PersonasForm from '../components/storyboard/PersonasForm.svelte'
    import UsersIcon from '../components/icons/Users.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import EditIcon from '../components/icons/PencilIcon.svelte'
    import DownCarrotIcon from '../components/icons/ChevronDown.svelte'
    import CommentIcon from '../components/icons/CommentIcon.svelte'
    import DeleteStoryboard from '../components/storyboard/DeleteStoryboard.svelte'
    import EditStoryboard from '../components/storyboard/EditStoryboard.svelte'
    import UpCarrotIcon from '../components/icons/ChevronUp.svelte'
    import { AppConfig, appRoutes, PathPrefix } from '../config'
    import { warrior as user } from '../stores.js'
    import { _ } from '../i18n.js'
    import BecomeFacilitator from '../components/user/BecomeFacilitator.svelte'

    export let storyboardId
    export let notifications
    export let router
    export let eventTag

    const { AllowRegistration, AllowGuests } = AppConfig
    const loginOrRegister =
        AllowRegistration || AllowGuests ? appRoutes.register : appRoutes.login

    const hostname = window.location.origin
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'

    let JoinPassRequired = false
    let socketError = false
    let socketReconnecting = false
    let storyboard = {
        goals: [],
        users: [],
        colorLegend: [],
        personas: [],
        facilitators: [],
        facilitatorCode: '',
        joinCode: '',
    }
    let showUsers = false
    let showColorLegend = false
    let showColorLegendForm = false
    let showPersonas = false
    let showPersonasForm = null
    let editColumn = null
    let activeStory = null
    let showDeleteStoryboard = false
    let showEditStoryboard = false
    let joinPasscode = ''
    let collapseGoals = []

    const onSocketMessage = function (evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'join_code_required':
                JoinPassRequired = true
                break
            case 'join_code_incorrect':
                notifications.danger($_('incorrectPassCode'))
                break
            case 'init':
                JoinPassRequired = false
                storyboard = JSON.parse(parsedEvent.value)
                eventTag('join', 'storyboard', '')
                break
            case 'user_joined':
                storyboard.users = JSON.parse(parsedEvent.value)
                const joinedUser = storyboard.users.find(
                    w => w.id === parsedEvent.userId,
                )
                notifications.success(`${joinedUser.name} joined.`)
                break
            case 'user_retreated':
                const leftUser = storyboard.users.find(
                    w => w.id === parsedEvent.userId,
                )
                storyboard.users = JSON.parse(parsedEvent.value)

                notifications.danger(`${leftUser.name} left.`)
                break
            case 'storyboard_updated':
                storyboard = JSON.parse(parsedEvent.value)
                break
            case 'goal_added':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'goal_revised':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'goal_deleted':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'column_added':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'column_updated':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'story_added':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'story_updated':
                storyboard.goals = JSON.parse(parsedEvent.value)
                if (activeStory) {
                    let activeStoryFound = false
                    for (let goal of storyboard.goals) {
                        for (let column of goal.columns) {
                            for (let story of column.stories) {
                                if (story.id === activeStory.id) {
                                    activeStory = story
                                    break
                                }
                            }
                            if (activeStoryFound) {
                                break
                            }
                        }
                        if (activeStoryFound) {
                            break
                        }
                    }
                }
                break
            case 'story_moved':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'story_deleted':
                storyboard.goals = JSON.parse(parsedEvent.value)
                break
            case 'personas_updated':
                storyboard.personas = JSON.parse(parsedEvent.value)
                break
            case 'storyboard_edited':
                const revisedStoryboard = JSON.parse(parsedEvent.value)
                storyboard.name = revisedStoryboard.storyboardName
                storyboard.joinCode = revisedStoryboard.joinCode
                break
            case 'storyboard_conceded':
                // storyboard over, goodbye.
                notifications.warning($_('storyboardDeleted'))
                router.route(appRoutes.storyboards)
                break
            default:
                break
        }
    }

    const ws = new Sockette(
        `${socketExtension}://${window.location.host}${PathPrefix}/api/storyboard/${storyboardId}`,
        {
            timeout: 2e3,
            maxAttempts: 15,
            onmessage: onSocketMessage,
            onerror: () => {
                socketError = true
                eventTag('socket_error', 'storyboard', '')
            },
            onclose: e => {
                if (e.code === 4004) {
                    eventTag('not_found', 'storyboard', '', () => {
                        router.route(appRoutes.storyboards)
                    })
                } else if (e.code === 4001) {
                    eventTag('socket_unauthorized', 'storyboard', '', () => {
                        user.delete()
                        router.route(`${appRoutes.login}/${storyboardId}`)
                    })
                } else if (e.code === 4003) {
                    eventTag('socket_duplicate', 'storyboard', '', () => {
                        notifications.danger($_('duplicateStoryboardSession'))
                        router.route(`${appRoutes.storyboards}`)
                    })
                } else if (e.code === 4002) {
                    eventTag(
                        'storyboard_user_abandoned',
                        'storyboard',
                        '',
                        () => {
                            router.route(appRoutes.storyboards)
                        },
                    )
                } else {
                    socketReconnecting = true
                    eventTag('socket_close', 'storyboard', '')
                }
            },
            onopen: () => {
                socketError = false
                socketReconnecting = false
                eventTag('socket_open', 'storyboard', '')
            },
            onmaximum: () => {
                socketReconnecting = false
                eventTag(
                    'socket_error',
                    'storyboard',
                    'Socket Reconnect Max Reached',
                )
            },
        },
    )

    onDestroy(() => {
        eventTag('leave', 'storyboard', '', () => {
            ws.close()
        })
    })

    const sendSocketEvent = (type, value) => {
        ws.send(
            JSON.stringify({
                type,
                value,
            }),
        )
    }

    // event handlers
    function handleDndConsider(e) {
        const goalIndex = e.target.dataset.goalindex
        const columnIndex = e.target.dataset.columnindex

        storyboard.goals[goalIndex].columns[columnIndex].stories =
            e.detail.items
        storyboard.goals = storyboard.goals
    }

    function handleDndFinalize(e) {
        const goalIndex = e.target.dataset.goalindex
        const columnIndex = e.target.dataset.columnindex
        const storyId = e.detail.info.id

        storyboard.goals[goalIndex].columns[columnIndex].stories =
            e.detail.items
        storyboard.goals = storyboard.goals

        const matchedStory = storyboard.goals[goalIndex].columns[
            columnIndex
        ].stories.find(i => i.id === storyId)

        if (matchedStory) {
            const goalId = storyboard.goals[goalIndex].id
            const columnId = storyboard.goals[goalIndex].columns[columnIndex].id

            // determine what story to place story before in target column
            const matchedStoryIndex =
                storyboard.goals[goalIndex].columns[
                    columnIndex
                ].stories.indexOf(matchedStory)
            const sibling =
                storyboard.goals[goalIndex].columns[columnIndex].stories[
                    matchedStoryIndex + 1
                ]
            const placeBefore = sibling ? sibling.id : ''

            sendSocketEvent(
                'move_story',
                JSON.stringify({
                    storyId,
                    goalId,
                    columnId,
                    placeBefore,
                }),
            )
            eventTag('story_move', 'storyboard', '')
        }
    }

    function authStoryboard(e) {
        e.preventDefault()

        sendSocketEvent('auth_storyboard', joinPasscode)
        eventTag('auth_storyboard', 'storyboard', '')
    }

    const addStory = (goalId, columnId) => () => {
        sendSocketEvent(
            'add_story',
            JSON.stringify({
                goalId,
                columnId,
            }),
        )
        eventTag('story_add', 'storyboard', '')
    }

    const addStoryColumn = goalId => () => {
        sendSocketEvent(
            'add_column',
            JSON.stringify({
                goalId,
            }),
        )
        eventTag('column_add', 'storyboard', '')
    }

    const deleteColumn = columnId => () => {
        sendSocketEvent('delete_column', columnId)
        eventTag('column_delete', 'storyboard', '')
        toggleColumnEdit()()
    }

    const handleAddFacilitator = userId => () => {
        sendSocketEvent(
            'facilitator_add',
            JSON.stringify({
                userId,
            }),
        )
    }

    const handleRemoveFacilitator = userId => () => {
        sendSocketEvent(
            'facilitator_remove',
            JSON.stringify({
                userId,
            }),
        )
    }

    function concedeStoryboard() {
        eventTag('concede_storyboard', 'storyboard', '', () => {
            sendSocketEvent('concede_storyboard', '')
        })
    }

    function abandonStoryboard() {
        eventTag('abandon_storyboard', 'storyboard', '', () => {
            sendSocketEvent('abandon_storyboard', '')
        })
    }

    function toggleUsersPanel() {
        showColorLegend = false
        showPersonas = false
        showUsers = !showUsers
        eventTag('show_users', 'storyboard', `show: ${showUsers}`)
    }

    function toggleColorLegend() {
        showUsers = false
        showPersonas = false
        showColorLegend = !showColorLegend
        eventTag('show_colorlegend', 'storyboard', `show: ${showColorLegend}`)
    }

    function togglePersonas() {
        showUsers = false
        showColorLegend = false
        showPersonas = !showPersonas
        eventTag('show_personas', 'storyboard', `show: ${showPersonas}`)
    }

    function toggleColumnEdit(column) {
        return () => {
            editColumn = editColumn != null ? null : column
        }
    }

    function toggleEditLegend() {
        showColorLegend = false
        showColorLegendForm = !showColorLegendForm
        eventTag(
            'show_edit_legend',
            'storyboard',
            `show: ${showColorLegendForm}`,
        )
    }

    const toggleEditPersona = persona => () => {
        showPersonas = false
        showPersonasForm = showPersonasForm != null ? null : persona
        eventTag(
            'show_edit_personas',
            'storyboard',
            `show: ${showPersonasForm}`,
        )
    }

    const toggleDeleteStoryboard = () => {
        showDeleteStoryboard = !showDeleteStoryboard
    }

    let showAddGoal = false
    let reviseGoalId = ''
    let reviseGoalName = ''

    const toggleAddGoal = goalId => () => {
        if (goalId) {
            const goalName = storyboard.goals.find(p => p.id === goalId).name
            reviseGoalId = goalId
            reviseGoalName = goalName
        } else {
            reviseGoalId = ''
            reviseGoalName = ''
        }
        showAddGoal = !showAddGoal
        eventTag('show_goal_add', 'storyboard', `show: ${showAddGoal}`)
    }

    const handleGoalAdd = goalName => {
        sendSocketEvent('add_goal', goalName)
        eventTag('goal_add', 'storyboard', '')
    }

    const handleGoalRevision = updatedGoal => {
        sendSocketEvent('revise_goal', JSON.stringify(updatedGoal))
        eventTag('goal_edit_name', 'storyboard', '')
    }

    const handleGoalDeletion = goalId => () => {
        sendSocketEvent('delete_goal', goalId)
        eventTag('goal_delete', 'storyboard', '')
    }

    const handleColumnRevision = column => {
        sendSocketEvent('revise_column', JSON.stringify(column))
        eventTag('column_revise', 'storyboard', '')
    }

    const handleLegendRevision = legend => {
        sendSocketEvent('revise_color_legend', JSON.stringify(legend))
        eventTag('color_legend_revise', 'storyboard', '')
    }

    const handlePersonaAdd = persona => {
        sendSocketEvent('add_persona', JSON.stringify(persona))
        eventTag('persona_add', 'storyboard', '')
    }

    const handlePersonaRevision = persona => {
        sendSocketEvent('revise_persona', JSON.stringify(persona))
        eventTag('persona_revise', 'storyboard', '')
    }

    const handleDeletePersona = personaId => () => {
        sendSocketEvent('delete_persona', personaId)
        eventTag('persona_delete', 'storyboard', '')
    }

    function handleStoryboardEdit(revisedStoryboard) {
        sendSocketEvent('edit_storyboard', JSON.stringify(revisedStoryboard))
        eventTag('edit_storyboard', 'storyboard', '')
        toggleEditStoryboard()
    }

    function toggleEditStoryboard() {
        showEditStoryboard = !showEditStoryboard
    }

    const toggleStoryForm = story => () => {
        activeStory = activeStory != null ? null : story
    }

    function toggleGoalCollapse(goalId) {
        return () => {
            const goalIndex = collapseGoals.indexOf(goalId)
            if (goalIndex > -1) {
                delete collapseGoals[goalIndex]
            } else {
                collapseGoals.push(goalId)
            }
            collapseGoals = collapseGoals
        }
    }

    let showBecomeFacilitator = false

    function becomeFacilitator(facilitatorCode) {
        sendSocketEvent('facilitator_self', facilitatorCode)
        eventTag('become_facilitator', 'retro', '')
        toggleBecomeFacilitator()
    }

    function toggleBecomeFacilitator() {
        showBecomeFacilitator = !showBecomeFacilitator
        eventTag('toggle_become_facilitator', 'retro', '')
    }

    $: isFacilitator =
        storyboard.facilitators && storyboard.facilitators.includes($user.id)

    onMount(() => {
        if (!$user.id) {
            router.route(`${loginOrRegister}/storyboard/${storyboardId}`)
        }
    })
</script>

<style>
    .story-gray {
        @apply border-gray-400;
    }

    .story-gray:hover {
        @apply border-gray-800;
    }

    .story-red {
        @apply border-red-400;
    }

    .story-red:hover {
        @apply border-red-800;
    }

    .story-orange {
        @apply border-orange-400;
    }

    .story-orange:hover {
        @apply border-orange-800;
    }

    .story-yellow {
        @apply border-yellow-400;
    }

    .story-yellow:hover {
        @apply border-yellow-800;
    }

    .story-green {
        @apply border-green-400;
    }

    .story-green:hover {
        @apply border-green-800;
    }

    .story-teal {
        @apply border-teal-400;
    }

    .story-teal:hover {
        @apply border-teal-800;
    }

    .story-blue {
        @apply border-blue-400;
    }

    .story-blue:hover {
        @apply border-blue-800;
    }

    .story-indigo {
        @apply border-indigo-400;
    }

    .story-indigo:hover {
        @apply border-indigo-800;
    }

    .story-purple {
        @apply border-purple-400;
    }

    .story-purple:hover {
        @apply border-purple-800;
    }

    .story-pink {
        @apply border-pink-400;
    }

    .story-pink:hover {
        @apply border-pink-800;
    }

    .colorcard-gray {
        @apply bg-gray-400;
    }

    .colorcard-red {
        @apply bg-red-400;
    }

    .colorcard-orange {
        @apply bg-orange-400;
    }

    .colorcard-yellow {
        @apply bg-yellow-400;
    }

    .colorcard-green {
        @apply bg-green-400;
    }

    .colorcard-teal {
        @apply bg-teal-400;
    }

    .colorcard-blue {
        @apply bg-blue-400;
    }

    .colorcard-indigo {
        @apply bg-indigo-400;
    }

    .colorcard-purple {
        @apply bg-purple-400;
    }

    .colorcard-pink {
        @apply bg-pink-400;
    }
</style>

<svelte:head>
    <title>{$_('storyboard')} {storyboard.name} | {$_('appName')}</title>
</svelte:head>

{#if storyboard.name && !socketReconnecting && !socketError}
    <div class="w-full">
        <div
            class="px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b border-t border-gray-400 dark:border-gray-700 flex
        flex-wrap"
        >
            <div class="w-1/3">
                <h1 class="text-3xl font-bold leading-tight dark:text-gray-200">
                    {storyboard.name}
                </h1>
            </div>
            <div class="w-2/3 text-right">
                <div>
                    {#if isFacilitator}
                        <HollowButton
                            color="green"
                            onClick="{toggleAddGoal()}"
                            additionalClasses="mr-2"
                            testid="goal-add"
                        >
                            {$_('storyboardAddGoal')}
                        </HollowButton>
                        <HollowButton
                            color="blue"
                            onClick="{toggleEditStoryboard}"
                            testid="storyboard-edit"
                        >
                            {$_('editStoryboard')}
                        </HollowButton>
                        <HollowButton
                            color="red"
                            onClick="{toggleDeleteStoryboard}"
                            additionalClasses="mr-2"
                            testid="storyboard-delete"
                        >
                            {$_('deleteStoryboard')}
                        </HollowButton>
                    {:else}
                        <HollowButton
                            color="blue"
                            onClick="{toggleBecomeFacilitator}"
                            testid="become-facilitator"
                        >
                            {$_('becomeFacilitator')}
                        </HollowButton>
                        <HollowButton
                            color="red"
                            onClick="{abandonStoryboard}"
                            testid="storyboard-leave"
                        >
                            {$_('leaveStoryboard')}
                        </HollowButton>
                    {/if}
                    <div class="inline-block relative">
                        <HollowButton
                            color="indigo"
                            additionalClasses="transition ease-in-out duration-150"
                            onClick="{togglePersonas}"
                            testid="personas-toggle"
                        >
                            {$_('personas')}
                            <DownCarrotIcon additionalClasses="ml-1" />
                        </HollowButton>
                        {#if showPersonas}
                            <div
                                class="origin-top-right absolute right-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
                            >
                                <div
                                    class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
                                >
                                    <div class="p-2">
                                        {#each storyboard.personas as persona}
                                            <div class="mb-1 w-full">
                                                <div>
                                                    <span class="font-bold">
                                                        {persona.name}
                                                    </span>
                                                    {#if isFacilitator}
                                                        &nbsp;|&nbsp;
                                                        <button
                                                            on:click="{toggleEditPersona(
                                                                persona,
                                                            )}"
                                                            class="text-orange-500
                                                        hover:text-orange-800"
                                                            data-testid="persona-edit"
                                                        >
                                                            {$_('edit')}
                                                        </button>
                                                        &nbsp;|&nbsp;
                                                        <button
                                                            on:click="{handleDeletePersona(
                                                                persona.id,
                                                            )}"
                                                            class="text-red-500
                                                        hover:text-red-800"
                                                            data-testid="persona-delete"
                                                        >
                                                            {$_('delete')}
                                                        </button>
                                                    {/if}
                                                </div>
                                                <span class="text-sm">
                                                    {persona.role}
                                                </span>
                                            </div>
                                        {/each}
                                    </div>

                                    {#if isFacilitator}
                                        <div class="p-2 text-right">
                                            <HollowButton
                                                color="green"
                                                onClick="{toggleEditPersona({
                                                    id: '',
                                                    name: '',
                                                    role: '',
                                                    description: '',
                                                })}"
                                                testid="persona-add"
                                            >
                                                {$_('addPersona')}
                                            </HollowButton>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>
                    <div class="inline-block relative">
                        <HollowButton
                            color="teal"
                            additionalClasses="transition ease-in-out duration-150"
                            onClick="{toggleColorLegend}"
                            testid="colorlegend-toggle"
                        >
                            {$_('colorLegend')}
                            <DownCarrotIcon additionalClasses="ml-1" />
                        </HollowButton>
                        {#if showColorLegend}
                            <div
                                class="origin-top-right absolute right-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
                            >
                                <div
                                    class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
                                >
                                    <div class="p-2">
                                        {#each storyboard.color_legend as color}
                                            <div class="mb-1 flex w-full">
                                                <span
                                                    class="p-4 mr-2 inline-block
                                                colorcard-{color.color}"></span>
                                                <span
                                                    class="inline-block align-middle
                                                {color.legend === ''
                                                        ? 'text-gray-300 dark:text-gray-500'
                                                        : 'text-gray-600 dark:text-gray-200'}"
                                                >
                                                    {color.legend ||
                                                        $_(
                                                            'colorLegendNotSpecified',
                                                        )}
                                                </span>
                                            </div>
                                        {/each}
                                    </div>

                                    {#if isFacilitator}
                                        <div class="p-2 text-right">
                                            <HollowButton
                                                color="orange"
                                                onClick="{toggleEditLegend}"
                                                testid="colorlegend-edit"
                                            >
                                                {$_('editColorLegend')}
                                            </HollowButton>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>
                    <div class="inline-block relative">
                        <HollowButton
                            color="orange"
                            additionalClasses="transition ease-in-out duration-150"
                            onClick="{toggleUsersPanel}"
                            testid="users-toggle"
                        >
                            <UsersIcon
                                additionalClasses="mr-1"
                                height="18"
                                width="18"
                            />
                            {$_('users')}
                            <DownCarrotIcon additionalClasses="ml-1" />
                        </HollowButton>
                        {#if showUsers}
                            <div
                                class="origin-top-right absolute right-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
                            >
                                <div
                                    class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
                                >
                                    {#each storyboard.users as usr, index (usr.id)}
                                        {#if usr.active}
                                            <UserCard
                                                user="{usr}"
                                                sendSocketEvent="{sendSocketEvent}"
                                                showBorder="{index !==
                                                    storyboard.users.length -
                                                        1}"
                                                facilitators="{storyboard.facilitators}"
                                                handleAddFacilitator="{handleAddFacilitator}"
                                                handleRemoveFacilitator="{handleRemoveFacilitator}"
                                            />
                                        {/if}
                                    {/each}

                                    <div class="p-2">
                                        <InviteUser
                                            hostname="{hostname}"
                                            storyboardId="{storyboard.id}"
                                        />
                                    </div>
                                </div>
                            </div>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
        {#each storyboard.goals as goal, goalIndex (goal.id)}
            <div data-goalid="{goal.id}">
                <div
                    class="flex px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b-2 border-gray-400 dark:border-gray-700 {goalIndex >
                    0
                        ? 'border-t-2'
                        : ''}"
                >
                    <div class="w-3/4 relative">
                        <div class="font-bold dark:text-gray-200 text-xl">
                            <h2 class="inline-block align-middle pt-1">
                                <button
                                    on:click="{toggleGoalCollapse(goal.id)}"
                                >
                                    {#if collapseGoals.includes(goal.id)}
                                        <DownCarrotIcon
                                            additionalClasses="mr-1"
                                        />
                                    {:else}
                                        <UpCarrotIcon
                                            additionalClasses="mr-1"
                                        />
                                    {/if}
                                </button>{goal.name}
                            </h2>
                        </div>
                    </div>
                    <div class="w-1/4 text-right">
                        {#if isFacilitator}
                            <HollowButton
                                color="green"
                                onClick="{addStoryColumn(goal.id)}"
                                btnSize="small"
                                testid="column-add"
                            >
                                {$_('storyboardAddColumn')}
                            </HollowButton>
                            <HollowButton
                                color="orange"
                                onClick="{toggleAddGoal(goal.id)}"
                                btnSize="small"
                                additionalClasses="ml-2"
                                testid="goal-edit"
                            >
                                {$_('edit')}
                            </HollowButton>
                            <HollowButton
                                color="red"
                                onClick="{handleGoalDeletion(goal.id)}"
                                btnSize="small"
                                additionalClasses="ml-2"
                                testid="goal-delete"
                            >
                                {$_('delete')}
                            </HollowButton>
                        {/if}
                    </div>
                </div>
                {#if !collapseGoals.includes(goal.id)}
                    <section class="flex px-2" style="overflow-x: scroll">
                        {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
                            <div class="flex-none my-4 mx-2 w-40">
                                <div class="flex-none">
                                    <div class="w-full mb-2">
                                        <div class="flex">
                                            <span
                                                class="font-bold flex-grow truncate dark:text-gray-300"
                                                title="{goalColumn.name}"
                                                data-testid="column-name"
                                            >
                                                {goalColumn.name}
                                            </span>
                                            <button
                                                on:click="{toggleColumnEdit(
                                                    goalColumn,
                                                )}"
                                                class="flex-none font-bold text-xl
                                        border-dashed border-2 border-gray-400 dark:border-gray-600
                                        hover:border-green-500 text-gray-600 dark:text-gray-400
                                        hover:text-green-500 py-1 px-2"
                                                title="{$_(
                                                    'storyboardEditColumn',
                                                )}"
                                                data-testid="column-edit"
                                            >
                                                <EditIcon />
                                            </button>
                                        </div>
                                    </div>
                                    <div class="w-full">
                                        <div class="flex">
                                            <button
                                                on:click="{addStory(
                                                    goal.id,
                                                    goalColumn.id,
                                                )}"
                                                class="flex-grow font-bold text-xl py-1
                                        px-2 border-dashed border-2
                                        border-gray-400 dark:border-gray-600 hover:border-green-500
                                        text-gray-600 dark:text-gray-400 hover:text-green-500"
                                                title="{$_(
                                                    'storyboardAddStoryToColumn',
                                                )}"
                                                data-testid="story-add"
                                            >
                                                +
                                            </button>
                                        </div>
                                    </div>
                                </div>
                                <div
                                    class="w-full relative"
                                    style="min-height: 160px;"
                                    data-goalid="{goal.id}"
                                    data-columnid="{goalColumn.id}"
                                    data-goalIndex="{goalIndex}"
                                    data-columnindex="{columnIndex}"
                                    use:dndzone="{{
                                        items: goalColumn.stories,
                                        type: 'story',
                                        dropTargetStyle: '',
                                        dropTargetClasses: [
                                            'outline',
                                            'outline-2',
                                            'outline-indigo-500',
                                            'dark:outline-yellow-400',
                                        ],
                                    }}"
                                    on:consider="{handleDndConsider}"
                                    on:finalize="{handleDndFinalize}"
                                >
                                    {#each goalColumn.stories as story (story.id)}
                                        <div
                                            class="relative max-w-xs shadow bg-white dark:bg-gray-700 dark:text-white border-l-4
                                    story-{story.color} border my-4
                                    cursor-pointer"
                                            style="list-style: none;"
                                            data-goalid="{goal.id}"
                                            data-columnid="{goalColumn.id}"
                                            data-storyid="{story.id}"
                                            on:click="{toggleStoryForm(story)}"
                                        >
                                            <div>
                                                <div>
                                                    <div
                                                        class="h-20 p-1 text-sm
                                                overflow-hidden {story.closed
                                                            ? 'line-through'
                                                            : ''}"
                                                        title="{story.name}"
                                                        data-testid="story-name"
                                                    >
                                                        {story.name}
                                                    </div>
                                                    <div class="h-8">
                                                        <div
                                                            class="flex content-center
                                                    p-1 text-sm"
                                                        >
                                                            <div
                                                                class="w-1/2
                                                        text-gray-600 dark:text-gray-300"
                                                            >
                                                                {#if story.comments.length > 0}
                                                                    <span
                                                                        class="inline-block
                                                                align-middle"
                                                                    >
                                                                        {story
                                                                            .comments
                                                                            .length}
                                                                        <CommentIcon
                                                                        />
                                                                    </span>
                                                                {/if}
                                                            </div>
                                                            <div
                                                                class="w-1/2 text-right"
                                                            >
                                                                {#if story.points > 0}
                                                                    <span
                                                                        class="px-2
                                                                bg-gray-300 dark:bg-gray-500
                                                                inline-block
                                                                align-middle"
                                                                    >
                                                                        {story.points}
                                                                    </span>
                                                                {/if}
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            {#if story[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
                                                <div
                                                    class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible opacity-50 max-w-xs shadow bg-white dark:bg-gray-700 dark:text-white border-l-4
                                    story-{story.color} border
                                    cursor-pointer"
                                                    style="list-style: none;"
                                                    data-goalid="{goal.id}"
                                                    data-columnid="{goalColumn.id}"
                                                    data-storyid="{story.id}"
                                                    on:click="{toggleStoryForm(
                                                        story,
                                                    )}"
                                                >
                                                    <div>
                                                        <div>
                                                            <div
                                                                class="h-20 p-1 text-sm
                                                overflow-hidden {story.closed
                                                                    ? 'line-through'
                                                                    : ''}"
                                                                title="{story.name}"
                                                            >
                                                                {story.name}
                                                            </div>
                                                            <div class="h-8">
                                                                <div
                                                                    class="flex content-center
                                                    p-1 text-sm"
                                                                >
                                                                    <div
                                                                        class="w-1/2
                                                        text-gray-600"
                                                                    >
                                                                        {#if story.comments.length > 0}
                                                                            <span
                                                                                class="inline-block
                                                                align-middle"
                                                                            >
                                                                                {story
                                                                                    .comments
                                                                                    .length}
                                                                                <CommentIcon
                                                                                />
                                                                            </span>
                                                                        {/if}
                                                                    </div>
                                                                    <div
                                                                        class="w-1/2 text-right"
                                                                    >
                                                                        {#if story.points > 0}
                                                                            <span
                                                                                class="px-2
                                                                bg-gray-300
                                                                inline-block
                                                                align-middle"
                                                                            >
                                                                                {story.points}
                                                                            </span>
                                                                        {/if}
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            {/if}
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        {/each}
                    </section>
                {/if}
            </div>
        {/each}
    </div>
{:else}
    <PageLayout>
        <div class="flex items-center">
            <div class="flex-1 text-center">
                {#if JoinPassRequired}
                    <div class="flex justify-center">
                        <div class="w-full md:w-1/2 lg:w-1/3">
                            <form
                                on:submit="{authStoryboard}"
                                class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                                name="authStoryboard"
                            >
                                <div class="mb-4">
                                    <label
                                        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                                        for="storyboardJoinCode"
                                    >
                                        {$_('passCodeRequired')}
                                    </label>
                                    <input
                                        bind:value="{joinPasscode}"
                                        placeholder="{$_('enterPasscode')}"
                                        class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                                        id="storyboardJoinCode"
                                        name="storyboardJoinCode"
                                        type="password"
                                        required
                                    />
                                </div>

                                <div class="text-right">
                                    <SolidButton type="submit"
                                        >{$_('joinStoryboard')}
                                    </SolidButton>
                                </div>
                            </form>
                        </div>
                    </div>
                {:else if socketReconnecting}
                    <h1
                        class="text-5xl text-orange-500 leading-tight font-bold"
                    >
                        {$_('reloadingStoryboard')}
                    </h1>
                {:else if socketError}
                    <h1 class="text-5xl text-red-500 leading-tight font-bold">
                        {$_('joinStoryboardError')}
                    </h1>
                {:else}
                    <h1 class="text-5xl text-green-500 leading-tight font-bold">
                        {$_('loadingStoryboard')}
                    </h1>
                {/if}
            </div>
        </div>
    </PageLayout>
{/if}

{#if showAddGoal}
    <AddGoal
        handleGoalAdd="{handleGoalAdd}"
        toggleAddGoal="{toggleAddGoal()}"
        handleGoalRevision="{handleGoalRevision}"
        goalId="{reviseGoalId}"
        goalName="{reviseGoalName}"
    />
{/if}

{#if editColumn}
    <ColumnForm
        handleColumnRevision="{handleColumnRevision}"
        toggleColumnEdit="{toggleColumnEdit()}"
        column="{editColumn}"
        deleteColumn="{deleteColumn}"
    />
{/if}

{#if activeStory}
    <StoryForm
        toggleStoryForm="{toggleStoryForm()}"
        story="{activeStory}"
        sendSocketEvent="{sendSocketEvent}"
        eventTag="{eventTag}"
        notifications="{notifications}"
        colorLegend="{storyboard.color_legend}"
        users="{storyboard.users}"
    />
{/if}

{#if showColorLegendForm}
    <ColorLegendForm
        handleLegendRevision="{handleLegendRevision}"
        toggleEditLegend="{toggleEditLegend}"
        colorLegend="{storyboard.color_legend}"
    />
{/if}

{#if showPersonasForm}
    <PersonasForm
        toggleEditPersona="{toggleEditPersona()}"
        persona="{showPersonasForm}"
        handlePersonaAdd="{handlePersonaAdd}"
        handlePersonaRevision="{handlePersonaRevision}"
    />
{/if}

{#if showEditStoryboard}
    <EditStoryboard
        storyboardName="{storyboard.name}"
        handleStoryboardEdit="{handleStoryboardEdit}"
        toggleEditStoryboard="{toggleEditStoryboard}"
        joinCode="{storyboard.joinCode}"
        facilitatorCode="{storyboard.facilitatorCode}"
    />
{/if}

{#if showDeleteStoryboard}
    <DeleteStoryboard
        toggleDelete="{toggleDeleteStoryboard}"
        handleDelete="{concedeStoryboard}"
    />
{/if}

{#if showBecomeFacilitator}
    <BecomeFacilitator
        handleBecomeFacilitator="{becomeFacilitator}"
        toggleBecomeFacilitator="{toggleBecomeFacilitator}"
    />
{/if}

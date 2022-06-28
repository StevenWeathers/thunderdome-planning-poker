<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import Checkin from '../components/checkin/Checkin.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import PencilIcon from '../components/icons/PencilIcon.svelte'
    import TrashIcon from '../components/icons/TrashIcon.svelte'
    import BlockedPing from '../components/checkin/BlockedPing.svelte'
    import Gauge from '../components/Gauge.svelte'
    import { _ } from '../i18n.js'
    import { warrior as user } from '../stores.js'
    import { appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'
    import {
        formatDayForInput,
        getTimezoneName,
        subtractDays,
    } from '../dateUtils.js'
    import TableRow from '../components/table/TableRow.svelte'
    import RowCol from '../components/table/RowCol.svelte'
    import HeadCol from '../components/table/HeadCol.svelte'
    import Table from '../components/table/Table.svelte'
    import CommentForm from '../components/checkin/CommentForm.svelte'
    import UserIcon from '../components/icons/UserIcon.svelte'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId
    export let departmentId
    export let teamId

    let showCheckin = false
    let checkins = []
    let now = new Date()
    let maxNegativeDate
    let selectedDate
    let selectedCheckin
    let stats = {
        participants: 0,
        pPerc: 0,
        goals: 0,
        gPerc: 0,
        blocked: 0,
        bPerc: 0,
    }

    let team = {
        id: teamId,
        name: '',
    }
    let organization = {
        id: organizationId,
        name: '',
    }
    let department = {
        id: departmentId,
        name: '',
    }
    let users = []
    let userCount = 1

    let organizationRole = ''
    let departmentRole = ''
    let teamRole = ''

    let showCommentForm = false
    let selectedCheckinId = null

    const apiPrefix = '/api'
    $: orgPrefix = departmentId
        ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
        : `${apiPrefix}/organizations/${organizationId}`
    $: teamPrefix = organizationId
        ? `${orgPrefix}/teams/${teamId}`
        : `${apiPrefix}/teams/${teamId}`

    function getTeam() {
        xfetch(teamPrefix)
            .then(res => res.json())
            .then(function (result) {
                team = result.data.team
                teamRole = result.data.teamRole

                if (departmentId) {
                    department = result.data.department
                    departmentRole = result.data.departmentRole
                }
                if (organizationId) {
                    organization = result.data.organization
                    organizationRole = result.data.organizationRole
                }
            })
            .catch(function () {
                notifications.danger($_('teamGetError'))
                eventTag('team_checkin_team', 'engagement', 'failure')
            })
    }

    function getCheckins() {
        xfetch(
            `${teamPrefix}/checkins?date=${selectedDate}&tz=${getTimezoneName()}`,
        )
            .then(res => res.json())
            .then(function (result) {
                checkins = result.data
            })
            .catch(function () {
                notifications.danger($_('getCheckinsError'))
                eventTag('team_checkin_checkins', 'engagement', 'failure')
            })
    }

    let userMap = {}

    function getUsers() {
        xfetch(`${teamPrefix}/users?limit=1000&offset=0`)
            .then(res => res.json())
            .then(function (result) {
                users = result.data
                userCount = result.meta.count
                userMap = users.reduce((prev, cur) => {
                    prev[cur.id] = cur.name
                    return prev
                }, {})
            })
            .catch(function () {
                notifications.danger($_('teamGetUsersError'))
                eventTag('team_checkin_users', 'engagement', 'failure')
            })
    }

    function toggleCheckin(checkin) {
        showCheckin = !showCheckin
        if (checkin) {
            selectedCheckin = checkin
        } else {
            selectedCheckin = null
        }
    }

    function handleCheckin(checkin) {
        const body = {
            ...checkin,
        }

        xfetch(`${teamPrefix}/checkins`, { body })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCheckin()
                eventTag('team_checkin_create', 'engagement', 'success')
            })
            .catch(function (error) {
                if (Array.isArray(error)) {
                    error[1].json().then(function (result) {
                        if (result.error === 'REQUIRES_TEAM_USER') {
                            notifications.danger(
                                $_('teamUserRequiredToCheckin'),
                            )
                        } else {
                            notifications.danger($_('checkinError'))
                        }
                    })
                } else {
                    notifications.danger($_('checkinError'))
                }
                eventTag('team_checkin_create', 'engagement', 'failure')
            })
    }

    function handleCheckinEdit(checkinId, checkin) {
        xfetch(`${teamPrefix}/checkins/${checkinId}`, {
            body: checkin,
            method: 'PUT',
        })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCheckin()
                eventTag('team_checkin_edit', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('updateCheckinError'))
                eventTag('team_checkin_edit', 'engagement', 'failure')
            })
    }

    function handleCheckinDelete(checkinId) {
        xfetch(`${teamPrefix}/checkins/${checkinId}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                eventTag('team_checkin_delete', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('deleteCheckinError'))
                eventTag('team_checkin_delete', 'engagement', 'failure')
            })
    }

    const toggleCommentForm = checkinId => () => {
        showCommentForm = !showCommentForm
        selectedCheckinId = checkinId
    }

    function handleCheckinComment(comment) {
        const body = {
            ...comment,
        }

        xfetch(`${teamPrefix}/checkins/${selectedCheckinId}/comments`, { body })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCommentForm(null)()
                eventTag('team_checkin_comment', 'engagement', 'success')
            })
            .catch(function (error) {
                if (Array.isArray(error)) {
                    error[1].json().then(function (result) {
                        if (result.error === 'REQUIRES_TEAM_USER') {
                            notifications.danger(
                                $_('teamUserRequiredToComment'),
                            )
                        } else {
                            notifications.danger($_('checkinCommentError'))
                        }
                    })
                } else {
                    notifications.danger($_('checkinCommentError'))
                }
                eventTag('team_checkin_comment', 'engagement', 'failure')
            })
    }

    const handleCommentDelete = (checkinId, commentId) => () => {
        xfetch(`${teamPrefix}/checkins/${checkinId}/comments/${commentId}`, {
            method: 'DELETE',
        })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                eventTag('team_checkin_comment_delete', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('checkinCommentDeleteError'))
                eventTag('team_checkin_comment_delete', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$user.id || !validateUserIsRegistered($user)) {
            router.route(appRoutes.login)
            return
        }

        selectedDate = formatDayForInput(now)
        maxNegativeDate = formatDayForInput(subtractDays(now, 60))

        getTeam()
        getCheckins()
        getUsers()
    })

    function calculateCheckinStats() {
        const ucs = []
        stats.blocked = 0
        stats.goals = 0

        checkins.map(c => {
            // @todo - remove once multiple same day checkins are prevented
            if (!ucs.includes(c.user.id)) {
                ucs.push(c.user.id)

                if (c.blockers !== '') {
                    ++stats.blocked
                }
                if (c.goalsMet) {
                    ++stats.goals
                }
            }
        })
        stats.participants = ucs.length

        stats.pPerc = Math.round((100 * stats.participants) / (userCount || 1))
        stats.gPerc = Math.round(
            (100 * stats.goals) / (stats.participants || 1),
        )
        stats.bPerc = Math.round(
            (100 * stats.blocked) / (stats.participants || 1),
        )

        stats = stats

        return stats
    }

    $: isAdmin =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole === 'ADMIN'
    $: isTeamMember =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole !== ''

    $: checkStats = checkins && userCount && calculateCheckinStats()
    $: alreadyCheckedIn =
        checkins && checkins.find(c => c.user.id === $user.id) !== undefined
</script>

<style global>
    ::-webkit-calendar-picker-indicator {
        margin-left: 0px;
        background-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="16" height="15" viewBox="0 0 24 24"><path fill="%2322c55e" d="M20 3h-1V1h-2v2H7V1H5v2H4c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 18H4V8h16v13z"/></svg>');
    }
</style>

<svelte:head>
    <title>{$_('team')} {team.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex sm:flex-wrap">
        <div class="md:grow">
            <h1
                class="text-3xl font-semibold font-rajdhani leading-none uppercase dark:text-white"
            >
                {$_('checkIn')}
                <ChevronRight class="w-8 h-8" />
                <input
                    type="date"
                    id="checkindate"
                    bind:value="{selectedDate}"
                    min="{maxNegativeDate}"
                    max="{formatDayForInput(now)}"
                    on:change="{getCheckins}"
                    class="bg-transparent"
                />
            </h1>

            {#if organizationId}
                <div
                    class="text-xl font-semibold font-rajdhani dark:text-white"
                >
                    <span class="uppercase">{$_('organization')}</span>
                    <ChevronRight />
                    <a
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        href="{appRoutes.organization}/{organization.id}"
                        >{organization.name}</a
                    >
                    {#if departmentId}
                        &nbsp;
                        <ChevronRight />
                        <span class="uppercase">{$_('department')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                            >{department.name}</a
                        >
                        <ChevronRight />
                        <span class="uppercase">{$_('team')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                        >
                            {team.name}
                        </a>
                    {:else}
                        <ChevronRight />
                        <span class="uppercase">{$_('team')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                            href="{appRoutes.organization}/{organization.id}/team/{team.id}"
                        >
                            {team.name}
                        </a>
                    {/if}
                </div>
            {:else}
                <div
                    class="text-2xl font-semibold font-rajdhani dark:text-white"
                >
                    <span class="uppercase">{$_('team')}</span>
                    <ChevronRight />
                    <a
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        href="{appRoutes.team}/{team.id}"
                    >
                        {team.name}
                    </a>
                </div>
            {/if}
        </div>
        <div class="md:pl-2 md:shrink text-right">
            <SolidButton
                additionalClasses="font-rajdhani uppercase text-2xl"
                onClick="{toggleCheckin}"
                disabled="{selectedDate !== formatDayForInput(now) ||
                    alreadyCheckedIn}"
                >{$_('checkIn')}
            </SolidButton>
        </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 my-4">
        <div class="px-2 md:px-4">
            <Gauge
                text="{$_('participation')}"
                percentage="{stats.pPerc}"
                stat="{stats.pPerc}"
                count="{stats.participants} / {userCount}"
            />
        </div>
        <div class="px-2 md:px-4">
            <Gauge
                text="{$_('goalsMet')}"
                percentage="{stats.gPerc}"
                color="green"
                stat="{stats.gPerc}"
                count="{stats.goals} / {stats.participants}"
            />
        </div>
        <div class="px-2 md:px-4">
            <Gauge
                text="{$_('blocked')}"
                percentage="{stats.bPerc}"
                color="red"
                stat="{stats.bPerc}"
                count="{stats.blocked} / {stats.participants}"
            />
        </div>
    </div>

    <div class="mt-8">
        <Table>
            <tr slot="header">
                <HeadCol>
                    {$_('name')}
                </HeadCol>
                <HeadCol>{$_('yesterday')}</HeadCol>
                <HeadCol>{$_('today')}</HeadCol>
                <HeadCol>{$_('blockers')}</HeadCol>
                <HeadCol>{$_('discuss')}</HeadCol>
                <HeadCol>{$_('comments')}</HeadCol>
            </tr>
            <tbody slot="body" let:class="{className}" class="{className}">
                {#each checkins as checkin, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            <div class="flex items-center">
                                <div class="">
                                    <div
                                        class="relative cursor-pointer w-14 h-14"
                                    >
                                        <div class="relative w-full h-full">
                                            <div
                                                class="w-full h-full bg-gray-200 rounded-full shadow"
                                            >
                                                <UserAvatar
                                                    warriorId="{checkin.user
                                                        .id}"
                                                    avatar="{checkin.user
                                                        .avatar}"
                                                    gravatarHash="{checkin.user
                                                        .gravatarHash}"
                                                    options="{{
                                                        class: 'w-full h-full rounded-full',
                                                    }}"
                                                />
                                            </div>
                                            {#if checkin.goalsMet}
                                                <div
                                                    class="absolute bottom-0 w-1/4 h-1/4 rounded-full shadow-md"
                                                >
                                                    <svg
                                                        xmlns="http://www.w3.org/2000/svg"
                                                        viewBox="0 0 1024 1024"
                                                    >
                                                        <circle
                                                            class="text-white fill-current"
                                                            cx="512"
                                                            cy="512"
                                                            r="512"></circle>
                                                        <circle
                                                            fill="green"
                                                            class="fill-current text-green-500"
                                                            cx="512"
                                                            cy="512"
                                                            r="384"></circle>
                                                        <path
                                                            class="text-white fill-current"
                                                            d="M456.4 576.1L334.8 454.4l-81.1 81.1 121.6 121.7 81.1 81.1 81.1-81.1 243.3-243.3-81.1-81.1z"
                                                        ></path>
                                                    </svg>
                                                </div>
                                            {/if}
                                            {#if checkin.blockers !== ''}
                                                <BlockedPing />
                                            {/if}
                                            <div
                                                class="hidden absolute top-0 right-0 w-1/4 h-1/4 bg-white rounded-full shadow-md"
                                            >
                                                <!-- emoji -->
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                <div class="ml-4">
                                    <div class="text-sm font-medium">
                                        {checkin.user.name}
                                    </div>
                                    {#if checkin.user.id === $user.id || isAdmin}
                                        <div>
                                            <button
                                                on:click="{() => {
                                                    toggleCheckin(checkin)
                                                }}"
                                                class="text-blue-500"
                                                title="{$_('edit')}"
                                            >
                                                <span class="sr-only"
                                                    >{$_('edit')}</span
                                                >
                                                <PencilIcon />
                                            </button>
                                            <button
                                                on:click="{() => {
                                                    handleCheckinDelete(
                                                        checkin.id,
                                                    )
                                                }}"
                                                class="text-red-500"
                                                title="Delete"
                                            >
                                                <span class="sr-only"
                                                    >{$_('delete')}</span
                                                >
                                                <TrashIcon />
                                            </button>
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </RowCol>
                        <RowCol>
                            <div class="unreset whitespace-pre-wrap">
                                {@html checkin.yesterday}
                            </div>
                        </RowCol>
                        <RowCol>
                            <div class="unreset whitespace-pre-wrap">
                                {@html checkin.today}
                            </div>
                        </RowCol>
                        <RowCol>
                            <div class="unreset whitespace-pre-wrap">
                                {@html checkin.blockers}
                            </div>
                        </RowCol>
                        <RowCol>
                            <div class="unreset whitespace-pre-wrap">
                                {@html checkin.discuss}
                            </div>
                        </RowCol>
                        <RowCol>
                            <div class="whitespace-normal">
                                {#each checkin.comments as comment}
                                    <p
                                        class="pb-2 mb-2 border-b border-gray-300 dark:border-gray-600"
                                    >
                                        <span>{comment.comment}</span><br />
                                        <span
                                            ><UserIcon class="h-4 w-4" />
                                            {userMap[comment.user_id] ||
                                                '...'}</span
                                        >
                                        {#if comment.user_id === $user.id || comment.user_id === isAdmin}
                                            <button
                                                class="ml-2 text-red-500"
                                                on:click="{handleCommentDelete(
                                                    checkin.id,
                                                    comment.id,
                                                )}"
                                            >
                                                {$_('delete')}
                                            </button>
                                        {/if}
                                    </p>
                                {/each}
                            </div>
                            <div class="text-right pt-2">
                                <SolidButton
                                    onClick="{toggleCommentForm(checkin.id)}"
                                    >{$_('addComment')}
                                </SolidButton>
                            </div>
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>
    </div>

    {#if showCheckin}
        {#if selectedCheckin}
            <Checkin
                teamId="{team.id}"
                userId="{$user.id}"
                checkinId="{selectedCheckin.id}"
                yesterday="{selectedCheckin.yesterday}"
                today="{selectedCheckin.today}"
                blockers="{selectedCheckin.blockers}"
                discuss="{selectedCheckin.discuss}"
                goalsMet="{selectedCheckin.goalsMet}"
                toggleCheckin="{toggleCheckin}"
                handleCheckin="{handleCheckin}"
                handleCheckinEdit="{handleCheckinEdit}"
            />
        {:else}
            <Checkin
                teamId="{team.id}"
                userId="{$user.id}"
                toggleCheckin="{toggleCheckin}"
                handleCheckin="{handleCheckin}"
                handleCheckinEdit="{handleCheckinEdit}"
            />
        {/if}
    {/if}

    {#if showCommentForm}
        <CommentForm
            checkinId="{selectedCheckinId}"
            toggleForm="{toggleCommentForm(null)}"
            handleComment="{handleCheckinComment}"
        />
    {/if}
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import AddUser from '../components/user/AddUser.svelte'
    import DeleteConfirmation from '../components/DeleteConfirmation.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import CreateBattle from '../components/battle/CreateBattle.svelte'
    import CreateRetro from '../components/retro/CreateRetro.svelte'
    import CreateStoryboard from '../components/storyboard/CreateStoryboard.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CountryFlag from '../components/user/CountryFlag.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'
    import Table from '../components/table/Table.svelte'
    import HeadCol from '../components/table/HeadCol.svelte'
    import TableRow from '../components/table/TableRow.svelte'
    import RowCol from '../components/table/RowCol.svelte'
    import Modal from '../components/Modal.svelte'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId
    export let departmentId
    export let teamId

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig

    const battlesPageLimit = 1000
    const retrosPageLimit = 1000
    const storyboardsPageLimit = 1000
    const usersPageLimit = 1000

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
    let battles = []
    let retros = []
    let storyboards = []
    let showCreateBattle = false
    let showCreateRetro = false
    let showCreateStoryboard = false
    let showAddUser = false
    let showRemoveUser = false
    let showRemoveBattle = false
    let showRemoveRetro = false
    let showRemoveStoryboard = false
    let removeBattleId = null
    let removeRetroId = null
    let removeStoryboardId = null
    let removeUserId = null
    let usersPage = 1
    let battlesPage = 1
    let retrosPage = 1
    let storyboardsPage = 1

    let organizationRole = ''
    let departmentRole = ''
    let teamRole = ''

    const apiPrefix = '/api'
    $: orgPrefix = departmentId
        ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
        : `${apiPrefix}/organizations/${organizationId}`
    $: teamPrefix = organizationId
        ? `${orgPrefix}/teams/${teamId}`
        : `${apiPrefix}/teams/${teamId}`

    $: currentPageUrl = teamPrefix
        .replace('/api', '')
        .replace('organizations', 'organization')
        .replace('departments', 'department')
        .replace('teams', 'team')

    function toggleAddUser() {
        showAddUser = !showAddUser
    }

    function toggleCreateBattle() {
        showCreateBattle = !showCreateBattle
    }

    function toggleCreateRetro() {
        showCreateRetro = !showCreateRetro
    }

    function toggleCreateStoryboard() {
        showCreateStoryboard = !showCreateStoryboard
    }

    const toggleRemoveUser = userId => () => {
        showRemoveUser = !showRemoveUser
        removeUserId = userId
    }

    const toggleRemoveBattle = battleId => () => {
        showRemoveBattle = !showRemoveBattle
        removeBattleId = battleId
    }

    const toggleRemoveRetro = retroId => () => {
        showRemoveRetro = !showRemoveRetro
        removeRetroId = retroId
    }

    const toggleRemoveStoryboard = storyboardId => () => {
        showRemoveStoryboard = !showRemoveStoryboard
        removeStoryboardId = storyboardId
    }

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

                getBattles()
                getRetros()
                getStoryboards()
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('teamGetError'))
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(
            `${teamPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                users = result.data
            })
            .catch(function () {
                notifications.danger($_('teamGetUsersError'))
            })
    }

    function getBattles() {
        if (FeaturePoker) {
            const battlesOffset = (battlesPage - 1) * battlesPageLimit
            xfetch(
                `${teamPrefix}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
            )
                .then(res => res.json())
                .then(function (result) {
                    battles = result.data
                })
                .catch(function () {
                    notifications.danger($_('teamGetBattlesError'))
                })
        }
    }

    function getRetros() {
        if (FeatureRetro) {
            const retrosOffset = (retrosPage - 1) * retrosPageLimit
            xfetch(
                `${teamPrefix}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`,
            )
                .then(res => res.json())
                .then(function (result) {
                    retros = result.data
                })
                .catch(function () {
                    notifications.danger($_('teamGetRetrosError'))
                })
        }
    }

    function getStoryboards() {
        if (FeatureStoryboard) {
            const storyboardsOffset =
                (storyboardsPage - 1) * storyboardsPageLimit
            xfetch(
                `${teamPrefix}/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}`,
            )
                .then(res => res.json())
                .then(function (result) {
                    storyboards = result.data
                })
                .catch(function () {
                    notifications.danger($_('teamGetStoryboardsError'))
                })
        }
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role,
        }

        xfetch(`${teamPrefix}/users`, { body })
            .then(function () {
                eventTag('team_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success($_('userAddSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userAddError'))
                eventTag('team_add_user', 'engagement', 'failure')
            })
    }

    function handleUserRemove() {
        xfetch(`${teamPrefix}/users/${removeUserId}`, { method: 'DELETE' })
            .then(function () {
                eventTag('team_remove_user', 'engagement', 'success')
                toggleRemoveUser(null)()
                notifications.success($_('userRemoveSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userRemoveError'))
                eventTag('team_remove_user', 'engagement', 'failure')
            })
    }

    function handleBattleRemove() {
        xfetch(`${teamPrefix}/battles/${removeBattleId}`, { method: 'DELETE' })
            .then(function () {
                eventTag('team_remove_battle', 'engagement', 'success')
                toggleRemoveBattle(null)()
                notifications.success($_('battleRemoveSuccess'))
                getBattles()
            })
            .catch(function () {
                notifications.danger($_('battleRemoveError'))
                eventTag('team_remove_battle', 'engagement', 'failure')
            })
    }

    function handleRetroRemove() {
        xfetch(`${teamPrefix}/retros/${removeRetroId}`, { method: 'DELETE' })
            .then(function () {
                eventTag('team_remove_retro', 'engagement', 'success')
                toggleRemoveRetro(null)()
                notifications.success($_('retroRemoveSuccess'))
                getRetros()
            })
            .catch(function () {
                notifications.danger($_('retroRemoveError'))
                eventTag('team_remove_retro', 'engagement', 'failure')
            })
    }

    function handleStoryboardRemove() {
        xfetch(`${teamPrefix}/storyboards/${removeStoryboardId}`, {
            method: 'DELETE',
        })
            .then(function () {
                eventTag('team_remove_storyboard', 'engagement', 'success')
                toggleRemoveStoryboard(null)()
                notifications.success($_('storyboardRemoveSuccess'))
                getStoryboards()
            })
            .catch(function () {
                notifications.danger($_('storyboardRemoveError'))
                eventTag('team_remove_storyboard', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || !validateUserIsRegistered($warrior)) {
            router.route(appRoutes.login)
            return
        }

        getTeam()
    })

    $: isAdmin =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole === 'ADMIN'
    $: isTeamMember =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole !== ''
</script>

<svelte:head>
    <title>{$_('team')} {team.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex mb-6 lg:mb-8">
        <div class="flex-1">
            <h1 class="text-3xl font-semibold font-rajdhani dark:text-white">
                <span class="uppercase">{$_('team')}</span>
                <ChevronRight class="w-8 h-8" />
                {team.name}
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
                    >
                        {organization.name}
                    </a>
                    {#if departmentId}
                        &nbsp;
                        <ChevronRight />
                        <span class="uppercase">{$_('department')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                        >
                            {department.name}
                        </a>
                    {/if}
                </div>
            {/if}
        </div>
        <div class="flex-1 text-right">
            <SolidButton
                additionalClasses="font-rajdhani uppercase text-2xl"
                href="{`${currentPageUrl}/checkin`}"
                >Checkins
            </SolidButton>
        </div>
    </div>

    {#if FeaturePoker}
        <div class="w-full mb-6 lg:mb-8">
            <div class="flex w-full">
                <div class="flex-1">
                    <h2
                        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                    >
                        {$_('battles')}
                    </h2>
                </div>
                <div class="flex-1 text-right">
                    {#if isTeamMember}
                        <SolidButton onClick="{toggleCreateBattle}"
                            >Create Battle
                        </SolidButton>
                    {/if}
                </div>
            </div>

            <div class="flex flex-wrap">
                {#each battles as battle}
                    <div
                        class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                    >
                        <div class="flex flex-wrap items-center p-4">
                            <div
                                class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                            >
                                <span data-testid="battle-name"
                                    >{battle.name}</span
                                >
                            </div>
                            <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                                {#if isAdmin}
                                    <HollowButton
                                        onClick="{toggleRemoveBattle(
                                            battle.id,
                                        )}"
                                        color="red"
                                    >
                                        {$_('remove')}
                                    </HollowButton>
                                {/if}
                                <HollowButton
                                    href="{appRoutes.battle}/{battle.id}"
                                >
                                    {$_('battleJoin')}
                                </HollowButton>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>

        {#if showCreateBattle}
            <Modal closeModal="{toggleCreateBattle}">
                <CreateBattle
                    apiPrefix="{teamPrefix}"
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </Modal>
        {/if}
    {/if}

    {#if FeatureRetro}
        <div class="w-full mb-6 lg:mb-8">
            <div class="flex w-full">
                <div class="flex-1">
                    <h2
                        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                    >
                        Retros
                    </h2>
                </div>
                <div class="flex-1 text-right">
                    {#if isTeamMember}
                        <SolidButton onClick="{toggleCreateRetro}"
                            >Create Retro
                        </SolidButton>
                    {/if}
                </div>
            </div>

            <div class="flex flex-wrap">
                {#each retros as retro}
                    <div
                        class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                    >
                        <div class="flex flex-wrap items-center p-4">
                            <div
                                class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                            >
                                <span data-testid="retro-name"
                                    >{retro.name}</span
                                >
                            </div>
                            <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                                {#if isAdmin}
                                    <HollowButton
                                        onClick="{toggleRemoveRetro(retro.id)}"
                                        color="red"
                                    >
                                        {$_('remove')}
                                    </HollowButton>
                                {/if}
                                <HollowButton
                                    href="{appRoutes.retro}/{retro.id}"
                                >
                                    Join Retro
                                </HollowButton>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>

        {#if showCreateRetro}
            <Modal closeModal="{toggleCreateRetro}">
                <CreateRetro
                    apiPrefix="{teamPrefix}"
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </Modal>
        {/if}
    {/if}

    {#if FeatureStoryboard}
        <div class="w-full mb-6 lg:mb-8">
            <div class="flex w-full">
                <div class="flex-1">
                    <h2
                        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                    >
                        Storyboards
                    </h2>
                </div>
                <div class="flex-1 text-right">
                    {#if isTeamMember}
                        <SolidButton onClick="{toggleCreateStoryboard}"
                            >Create Storyboard
                        </SolidButton>
                    {/if}
                </div>
            </div>

            <div class="flex flex-wrap">
                {#each storyboards as storyboard}
                    <div
                        class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                    >
                        <div class="flex flex-wrap items-center p-4">
                            <div
                                class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                            >
                                <span data-testid="storyboard-name"
                                    >{storyboard.name}</span
                                >
                            </div>
                            <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                                {#if isAdmin}
                                    <HollowButton
                                        onClick="{toggleRemoveStoryboard(
                                            storyboard.id,
                                        )}"
                                        color="red"
                                    >
                                        {$_('remove')}
                                    </HollowButton>
                                {/if}
                                <HollowButton
                                    href="{appRoutes.storyboard}/{storyboard.id}"
                                >
                                    Join Storyboard
                                </HollowButton>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>

        {#if showCreateStoryboard}
            <Modal closeModal="{toggleCreateStoryboard}">
                <CreateStoryboard
                    apiPrefix="{teamPrefix}"
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </Modal>
        {/if}
    {/if}

    <div class="w-full">
        <div class="flex w-full">
            <div class="w-4/5">
                <h2
                    class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                >
                    {$_('users')}
                </h2>
            </div>
            <div class="w-1/5">
                <div class="text-right">
                    {#if isAdmin}
                        <SolidButton onClick="{toggleAddUser}">
                            {$_('userAdd')}
                        </SolidButton>
                    {/if}
                </div>
            </div>
        </div>

        <Table>
            <tr slot="header">
                <HeadCol>
                    {$_('name')}
                </HeadCol>
                <HeadCol>
                    {$_('email')}
                </HeadCol>
                <HeadCol>
                    {$_('role')}
                </HeadCol>
                <HeadCol type="action">
                    <span class="sr-only">Actions</span>
                </HeadCol>
            </tr>
            <tbody slot="body" let:class="{className}" class="{className}">
                {#each users as user, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            <div class="flex items-center">
                                <div class="flex-shrink-0 h-10 w-10">
                                    <UserAvatar
                                        warriorId="{user.id}"
                                        avatar="{user.avatar}"
                                        gravatarHash="{user.gravatarHash}"
                                        width="48"
                                        class="h-10 w-10 rounded-full"
                                    />
                                </div>
                                <div class="ml-4">
                                    <div
                                        class="font-medium text-gray-900 dark:text-gray-200"
                                    >
                                        {user.name}
                                        {#if user.country}
                                            &nbsp;
                                            <CountryFlag
                                                country="{user.country}"
                                                additionalClass="inline-block"
                                                width="32"
                                                height="24"
                                            />
                                        {/if}
                                    </div>
                                </div>
                            </div>
                        </RowCol>
                        <RowCol>
                            {user.email}
                        </RowCol>
                        <RowCol>
                            <div
                                className="text-sm text-gray-500 dark:text-gray-300"
                            >
                                {user.role}
                            </div>
                        </RowCol>
                        <RowCol type="action">
                            {#if isAdmin}
                                <HollowButton
                                    onClick="{toggleRemoveUser(user.id)}"
                                    color="red"
                                >
                                    {$_('remove')}
                                </HollowButton>
                            {/if}
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>
    </div>

    {#if showAddUser}
        <AddUser toggleAdd="{toggleAddUser}" handleAdd="{handleUserAdd}" />
    {/if}

    {#if showRemoveUser}
        <DeleteConfirmation
            toggleDelete="{toggleRemoveUser(null)}"
            handleDelete="{handleUserRemove}"
            permanent="{false}"
            confirmText="{$_('removeUserConfirmText')}"
            confirmBtnText="{$_('removeUser')}"
        />
    {/if}

    {#if showRemoveBattle}
        <DeleteConfirmation
            toggleDelete="{toggleRemoveBattle(null)}"
            handleDelete="{handleBattleRemove}"
            permanent="{false}"
            confirmText="{$_('removeBattleConfirmText')}"
            confirmBtnText="{$_('removeBattle')}"
        />
    {/if}

    {#if showRemoveRetro}
        <DeleteConfirmation
            toggleDelete="{toggleRemoveRetro(null)}"
            handleDelete="{handleRetroRemove}"
            permanent="{false}"
            confirmText="Are you sure you want to remove this retro from the team?"
            confirmBtnText="Remove Retro"
        />
    {/if}

    {#if showRemoveStoryboard}
        <DeleteConfirmation
            toggleDelete="{toggleRemoveStoryboard(null)}"
            handleDelete="{handleStoryboardRemove}"
            permanent="{false}"
            confirmText="Are you sure you want to remove this storyboard from the team?"
            confirmBtnText="Remove Storyboard"
        />
    {/if}
</PageLayout>

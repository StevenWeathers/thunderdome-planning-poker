<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import AddUser from '../components/AddUser.svelte'
    import RemoveUser from '../components/RemoveUser.svelte'
    import RemoveBattle from '../components/RemoveBattle.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import CreateBattle from '../components/CreateBattle.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId
    export let departmentId
    export let teamId

    const battlesPageLimit = 1000
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
    let showAddUser = false
    let showRemoveUser = false
    let showRemoveBattle = false
    let removeBattleId = null
    let removeUserId = null
    let usersPage = 1
    let battlesPage = 1

    let organizationRole = ''
    let departmentRole = ''
    let teamRole = ''

    const apiPrefix = '/api'
    $: orgPrefix = departmentId
        ? `${apiPrefix}/organization/${organizationId}/department/${departmentId}`
        : `${apiPrefix}/organization/${organizationId}`
    $: teamPrefix = organizationId
        ? `${orgPrefix}/team/${teamId}`
        : `${apiPrefix}/team/${teamId}`

    function toggleAddUser() {
        showAddUser = !showAddUser
    }

    const toggleRemoveUser = userId => () => {
        showRemoveUser = !showRemoveUser
        removeUserId = userId
    }

    const toggleRemoveBattle = battleId => () => {
        showRemoveBattle = !showRemoveBattle
        removeBattleId = battleId
    }

    function getTeam() {
        xfetch(teamPrefix)
            .then(res => res.json())
            .then(function(result) {
                team = result.team
                teamRole = result.teamRole

                if (departmentId) {
                    department = result.department
                    departmentRole = result.departmentRole
                }
                if (organizationId) {
                    organization = result.organization
                    organizationRole = result.organizationRole
                }

                getBattles()
                getUsers()
            })
            .catch(function(error) {
                notifications.danger($_('teamGetError'))
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(`${teamPrefix}/users/${usersPageLimit}/${usersOffset}`)
            .then(res => res.json())
            .then(function(result) {
                users = result
            })
            .catch(function(error) {
                notifications.danger($_('teamGetUsersError'))
            })
    }

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit
        xfetch(`${teamPrefix}/battles/${battlesPageLimit}/${battlesOffset}`)
            .then(res => res.json())
            .then(function(result) {
                battles = result
            })
            .catch(function(error) {
                notifications.danger($_('teamGetBattlesError'))
            })
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role,
        }

        xfetch(`${teamPrefix}/users`, { body })
            .then(function() {
                eventTag('team_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success($_('userAddSuccess'))
                getUsers()
            })
            .catch(function() {
                notifications.danger($_('userAddError'))
                eventTag('team_add_user', 'engagement', 'failure')
            })
    }

    function handleUserRemove() {
        const body = {
            id: removeUserId,
        }

        xfetch(`${teamPrefix}/user`, { body, method: 'DELETE' })
            .then(function() {
                eventTag('team_remove_user', 'engagement', 'success')
                toggleRemoveUser(null)()
                notifications.success($_('userRemoveSuccess'))
                getUsers()
            })
            .catch(function() {
                notifications.danger($_('userRemoveError'))
                eventTag('team_remove_user', 'engagement', 'failure')
            })
    }

    function handleBattleRemove() {
        const body = {
            id: removeBattleId,
        }

        xfetch(`${teamPrefix}/battle`, { body, method: 'DELETE' })
            .then(function() {
                eventTag('team_remove_battle', 'engagement', 'success')
                toggleRemoveBattle(null)()
                notifications.success($_('battleRemoveSuccess'))
                getBattles()
            })
            .catch(function() {
                notifications.danger($_('battleRemoveError'))
                eventTag('team_remove_battle', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || $warrior.rank === 'PRIVATE') {
            router.route(appRoutes.login)
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
    <h1 class="text-3xl font-bold">{$_('team')}: {team.name}</h1>
    {#if organizationId}
        <div class="font-bold">
            {$_('organization')}
            <ChevronRight class="inline-block" />
            <a
                class="text-blue-500 hover:text-blue-800"
                href="{appRoutes.organization}/{organization.id}">
                {organization.name}
            </a>
            {#if departmentId}
                &nbsp;
                <ChevronRight class="inline-block" />
                {$_('department')}
                <ChevronRight class="inline-block" />
                <a
                    class="text-blue-500 hover:text-blue-800"
                    href="{appRoutes.organization}/{organization.id}/department/{department.id}">
                    {department.name}
                </a>
            {/if}
        </div>
    {/if}

    <div class="w-full mt-4">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded flex">
            <div class="w-full md:w-1/2 lg:w-3/5 md:pr-4">
                <div class="flex w-full">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('battles')}
                    </h2>
                </div>

                <table class="table-fixed w-full">
                    <thead>
                        <tr>
                            <th class="w-2/6 px-4 py-2">{$_('name')}</th>
                            <th class="w-1/6 px-4 py-2"></th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each battles as battle}
                            <tr>
                                <td class="border px-4 py-2">{battle.name}</td>
                                <td class="border px-4 py-2 text-right">
                                    {#if isAdmin}
                                        <HollowButton
                                            onClick="{toggleRemoveBattle(battle.id)}"
                                            color="red">
                                            {$_('remove')}
                                        </HollowButton>
                                    {/if}
                                    <HollowButton
                                        href="{appRoutes.battle}/{battle.id}">
                                        {$_('battleJoin')}
                                    </HollowButton>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>

            <div class="w-full md:w-1/2 lg:w-2/5 md:pl-2 xl:pl-4">
                {#if isTeamMember}
                    <h2 class="mb-4 text-2xl font-bold leading-tight">
                        {$_('pages.myBattles.createBattle.title')}
                    </h2>
                    <CreateBattle
                        apiPrefix="{teamPrefix}"
                        {notifications}
                        {router}
                        {eventTag}
                        {xfetch} />
                {/if}
            </div>
        </div>
    </div>

    <div class="w-full mt-4">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">{$_('users')}</h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleAddUser}">
                                {$_('userAdd')}
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">{$_('name')}</th>
                        <th class="w-2/6 px-4 py-2">{$_('email')}</th>
                        <th class="w-1/6 px-4 py-2">{$_('role')}</th>
                        <th class="w-1/6 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each users as usr}
                        <tr>
                            <td class="border px-4 py-2">{usr.name}</td>
                            <td class="border px-4 py-2">{usr.email}</td>
                            <td class="border px-4 py-2">{usr.role}</td>
                            <td class="border px-4 py-2 text-right">
                                {#if isAdmin}
                                    <HollowButton
                                        onClick="{toggleRemoveUser(usr.id)}"
                                        color="red">
                                        {$_('remove')}
                                    </HollowButton>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>

    {#if showAddUser}
        <AddUser toggleAdd="{toggleAddUser}" handleAdd="{handleUserAdd}" />
    {/if}

    {#if showRemoveUser}
        <RemoveUser
            toggleRemove="{toggleRemoveUser(null)}"
            handleRemove="{handleUserRemove}" />
    {/if}

    {#if showRemoveBattle}
        <RemoveBattle
            toggleRemove="{toggleRemoveBattle(null)}"
            handleRemove="{handleBattleRemove}" />
    {/if}
</PageLayout>

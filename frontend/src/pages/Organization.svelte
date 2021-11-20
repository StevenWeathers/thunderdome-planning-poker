<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import CreateDepartment from '../components/CreateDepartment.svelte'
    import CreateTeam from '../components/CreateTeam.svelte'
    import AddUser from '../components/AddUser.svelte'
    import RemoveUser from '../components/RemoveUser.svelte'
    import DeleteTeam from '../components/DeleteTeam.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'
    import { validateUserIsRegistered } from '../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId

    const departmentsPageLimit = 1000
    const teamsPageLimit = 1000
    const usersPageLimit = 1000

    let organization = {
        id: organizationId,
        name: '',
        createdDate: '',
        updateDate: '',
    }
    let role = 'MEMBER'
    let departments = []
    let teams = []
    let users = []
    let showCreateDepartment = false
    let showCreateTeam = false
    let showAddUser = false
    let showRemoveUser = false
    let removeUserId = null
    let showDeleteTeam = false
    let deleteTeamId = null
    let teamsPage = 1
    let usersPage = 1
    let departmentsPage = 1

    function toggleCreateDepartment() {
        showCreateDepartment = !showCreateDepartment
    }

    function toggleCreateTeam() {
        showCreateTeam = !showCreateTeam
    }

    function toggleAddUser() {
        showAddUser = !showAddUser
    }

    const toggleRemoveUser = userId => () => {
        showRemoveUser = !showRemoveUser
        removeUserId = userId
    }

    const toggleDeleteTeam = teamId => () => {
        showDeleteTeam = !showDeleteTeam
        deleteTeamId = teamId
    }

    function getOrganization() {
        xfetch(`/api/organizations/${organizationId}`)
            .then(res => res.json())
            .then(function (result) {
                organization = result.data.organization
                role = result.data.role

                getDepartments()
                getTeams()
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('organizationGetError'))
            })
    }

    function getDepartments() {
        const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit
        xfetch(
            `/api/organizations/${organizationId}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                departments = result.data
            })
            .catch(function () {
                notifications.danger($_('organizationGetDepartmentsError'))
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/organizations/${organizationId}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                teams = result.data
            })
            .catch(function () {
                notifications.danger($_('organizationGetTeamsError'))
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(
            `/api/organizations/${organizationId}/users?limit=${usersPageLimit}&offset=${usersOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                users = result.data
            })
            .catch(function () {
                notifications.danger($_('organizationGetUsersError'))
            })
    }

    function createDepartmentHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/organizations/${organizationId}/departments`, { body })
            .then(res => res.json())
            .then(function (result) {
                eventTag('create_department', 'engagement', 'success', () => {
                    router.route(
                        `${appRoutes.organization}/${organizationId}/department/${result.data.id}`,
                    )
                })
            })
            .catch(function () {
                notifications.danger($_('departmentCreateError'))
                eventTag('create_department', 'engagement', 'failure')
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/organizations/${organizationId}/teams`, { body })
            .then(res => res.json())
            .then(function () {
                eventTag('create_organization_team', 'engagement', 'success')
                toggleCreateTeam()
                notifications.success($_('teamCreateSuccess'))
                getTeams()
            })
            .catch(function () {
                notifications.danger($_('teamCreateError'))
                eventTag('create_organization_team', 'engagement', 'failure')
            })
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role,
        }

        xfetch(`/api/organizations/${organizationId}/users`, { body })
            .then(function () {
                eventTag('organization_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success($_('userAddSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userAddError'))
                eventTag('organization_add_user', 'engagement', 'failure')
            })
    }

    function handleUserRemove() {
        xfetch(`/api/organizations/${organizationId}/users/${removeUserId}`, {
            method: 'DELETE',
        })
            .then(function () {
                eventTag('organization_remove_user', 'engagement', 'success')
                toggleRemoveUser(null)()
                notifications.success($_('userRemoveSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userRemoveError'))
                eventTag('organization_remove_user', 'engagement', 'failure')
            })
    }

    function handleDeleteTeam() {
        xfetch(`/api/organizations/${organizationId}/teams/${deleteTeamId}`, {
            method: 'DELETE',
        })
            .then(function () {
                eventTag('organization_delete_team', 'engagement', 'success')
                toggleDeleteTeam(null)()
                notifications.success($_('teamDeleteSuccess'))
                getTeams()
            })
            .catch(function () {
                notifications.danger($_('teamDeleteError'))
                eventTag('organization_delete_team', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || !validateUserIsRegistered($warrior)) {
            router.route(appRoutes.login)
        }

        getOrganization()
    })

    $: isAdmin = role === 'ADMIN'
</script>

<svelte:head>
    <title>{$_('organization')} {organization.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <h1 class="mb-4 text-3xl font-bold">
        {$_('organization')}: {organization.name}
    </h1>

    <div class="w-full mb-4">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('departments')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleCreateDepartment}">
                                {$_('departmentCreate')}
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">{$_('name')}</th>
                    </tr>
                </thead>
                <tbody>
                    {#each departments as department}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="{appRoutes.organization}/{organizationId}/department/{department.id}"
                                    class="text-blue-500 hover:text-blue-800"
                                >
                                    {department.name}
                                </a>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>

    <div class="w-full mb-4">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('teams')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleCreateTeam}">
                                {$_('teamCreate')}
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-4/6 px-4 py-2">{$_('name')}</th>
                        <th class="w-2/6 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each teams as team}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="{appRoutes.organization}/{organizationId}/team/{team.id}"
                                    class="text-blue-500 hover:text-blue-800"
                                >
                                    {team.name}
                                </a>
                            </td>
                            <td class="border px-4 py-2 text-right">
                                {#if isAdmin}
                                    <HollowButton
                                        onClick="{toggleDeleteTeam(team.id)}"
                                        color="red"
                                    >
                                        {$_('delete')}
                                    </HollowButton>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('users')}
                    </h2>
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
                                        color="red"
                                    >
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

    {#if showCreateDepartment}
        <CreateDepartment
            toggleCreate="{toggleCreateDepartment}"
            handleCreate="{createDepartmentHandler}"
        />
    {/if}

    {#if showCreateTeam}
        <CreateTeam
            toggleCreate="{toggleCreateTeam}"
            handleCreate="{createTeamHandler}"
        />
    {/if}

    {#if showAddUser}
        <AddUser toggleAdd="{toggleAddUser}" handleAdd="{handleUserAdd}" />
    {/if}

    {#if showRemoveUser}
        <RemoveUser
            toggleRemove="{toggleRemoveUser(null)}"
            handleRemove="{handleUserRemove}"
        />
    {/if}

    {#if showDeleteTeam}
        <DeleteTeam
            toggleDelete="{toggleDeleteTeam(null)}"
            handleDelete="{handleDeleteTeam}"
        />
    {/if}
</PageLayout>

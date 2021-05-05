<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import CreateDepartment from '../components/CreateDepartment.svelte'
    import CreateTeam from '../components/CreateTeam.svelte'
    import AddUser from '../components/AddUser.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

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

    function getOrganization() {
        xfetch(`/api/organization/${organizationId}`)
            .then(res => res.json())
            .then(function(result) {
                organization = result.organization
                role = result.role

                getDepartments()
                getTeams()
                getUsers()
            })
            .catch(function(error) {
                notifications.danger('Error getting organization')
            })
    }

    function getDepartments() {
        const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit
        xfetch(
            `/api/organization/${organizationId}/departments/${departmentsPageLimit}/${departmentsOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                departments = result
            })
            .catch(function(error) {
                notifications.danger('Error getting organization departments')
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/organization/${organizationId}/teams/${teamsPageLimit}/${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                teams = result
            })
            .catch(function(error) {
                notifications.danger('Error getting organization teams')
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(
            `/api/organization/${organizationId}/users/${usersPageLimit}/${usersOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                users = result
            })
            .catch(function(error) {
                notifications.danger('Error getting organization users')
            })
    }

    function createDepartmentHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/organization/${organizationId}/departments`, { body })
            .then(res => res.json())
            .then(function(department) {
                eventTag('create_department', 'engagement', 'success', () => {
                    router.route(
                        `${appRoutes.organization}/${organizationId}/department/${department.id}`,
                    )
                })
            })
            .catch(function(error) {
                notifications.danger(
                    'Error attempting to create organization department',
                )
                eventTag('create_department', 'engagement', 'failure')
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/organization/${organizationId}/teams`, { body })
            .then(res => res.json())
            .then(function(organization) {
                eventTag('create_organization_team', 'engagement', 'success')
                toggleCreateTeam()
                notifications.success('Team created successfully.')
                getTeams()
            })
            .catch(function(error) {
                notifications.danger('Error attempting to create organization team')
                eventTag('create_organization_team', 'engagement', 'failure')
            })
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role
        }

        xfetch(`/api/organization/${organizationId}/users`, { body })
            .then(function() {
                eventTag('organization_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success('User added successfully.')
                getUsers()
            })
            .catch(function() {
                notifications.danger('Error attempting to add user to organization')
                eventTag('organization_add_user', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || $warrior.rank === 'PRIVATE') {
            router.route(appRoutes.login)
        }

        getOrganization()
    })

    $: isAdmin = role === 'ADMIN'
</script>

<PageLayout>
    <h1 class="mb-4 text-3xl font-bold">Organization: {organization.name}</h1>

    <div class="w-full mb-4">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        Departments
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleCreateDepartment}">
                                Create Department
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">Name</th>
                    </tr>
                </thead>
                <tbody>
                    {#each departments as department}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="/organization/{organizationId}/department/{department.id}"
                                    class="text-blue-500 hover:text-blue-800">
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
                        Teams
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleCreateTeam}">
                                Create Team
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">Name</th>
                    </tr>
                </thead>
                <tbody>
                    {#each teams as team}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="/organization/{organizationId}/team/{team.id}"
                                    class="text-blue-500 hover:text-blue-800">
                                    {team.name}
                                </a>
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
                        Users
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        {#if isAdmin}
                            <HollowButton onClick="{toggleAddUser}">
                                Add User
                            </HollowButton>
                        {/if}
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">Name</th>
                        <th class="w-2/6 px-4 py-2">Email</th>
                        <th class="w-1/6 px-4 py-2">Role</th>
                        <th class="w-1/6 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each users as usr}
                        <tr>
                            <td class="border px-4 py-2">{usr.name}</td>
                            <td class="border px-4 py-2">{usr.email}</td>
                            <td class="border px-4 py-2">{usr.role}</td>
                            <td class="border px-4 py-2 text-right"></td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>

    {#if showCreateDepartment}
        <CreateDepartment
            toggleCreate="{toggleCreateDepartment}"
            handleCreate="{createDepartmentHandler}" />
    {/if}

    {#if showCreateTeam}
        <CreateTeam toggleCreate="{toggleCreateTeam}" handleCreate={createTeamHandler} />
    {/if}

    {#if showAddUser}
        <AddUser toggleAdd="{toggleAddUser}" handleAdd={handleUserAdd} />
    {/if}
</PageLayout>

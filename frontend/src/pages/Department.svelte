<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import CreateTeam from '../components/CreateTeam.svelte'
    import AddUser from '../components/AddUser.svelte'
    import RemoveUser from '../components/RemoveUser.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId
    export let departmentId

    const teamsPageLimit = 1000
    const usersPageLimit = 1000

    let organization = {
        id: organizationId,
        name: ''
    }
    let department = {
        id: departmentId,
        name: ''
    }
    let departmentRole = ''
    let organizationRole = ''
    let teams = []
    let users = []
    let showCreateTeam = false
    let showAddUser = false
    let showRemoveUser = false
    let removeUserId = null
    let teamsPage = 1
    let usersPage = 1

    function toggleCreateTeam() {
        showCreateTeam = !showCreateTeam
    }

    function toggleAddUser() {
        showAddUser = !showAddUser
    }

    const toggleRemoveUser = (userId) => () => {
        showRemoveUser = !showRemoveUser
        removeUserId = userId
    }

    function getDepartment() {
        xfetch(`/api/organization/${organizationId}/department/${departmentId}`)
            .then(res => res.json())
            .then(function(result) {
                department = result.department
                organization = result.organization
                organizationRole = result.organizationRole
                departmentRole = result.departmentRole

                getTeams()
                getUsers()
            })
            .catch(function(error) {
                notifications.danger('Error getting department')
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/organization/${organizationId}/department/${departmentId}/teams/${teamsPageLimit}/${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                teams = result
            })
            .catch(function(error) {
                notifications.danger('Error getting department teams')
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(
            `/api/organization/${organizationId}/department/${departmentId}/users/${usersPageLimit}/${usersOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                users = result
            })
            .catch(function(error) {
                notifications.danger('Error getting department users')
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/organization/${organizationId}/department/${departmentId}/teams`, { body })
            .then(res => res.json())
            .then(function(organization) {
                eventTag('create_department_team', 'engagement', 'success')
                toggleCreateTeam()
                notifications.success('Team created successfully.')
                getTeams()
            })
            .catch(function(error) {
                notifications.danger('Error attempting to create department team')
                eventTag('create_department_team', 'engagement', 'failure')
            })
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role
        }

        xfetch(`/api/organization/${organizationId}/department/${departmentId}/users`, { body })
            .then(function() {
                eventTag('department_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success('User added successfully.')
                getUsers()
            })
            .catch(function() {
                notifications.danger('Error attempting to add user to department')
                eventTag('department_add_user', 'engagement', 'failure')
            })
    }

    function handleUserRemove() {
        const body = {
            id: removeUserId
        }

        xfetch(`/api/organization/${organizationId}/department/${departmentId}/user`, { body, method: 'DELETE' })
            .then(function() {
                eventTag('department_remove_user', 'engagement', 'success')
                toggleRemoveUser(null)()
                notifications.success('User removed successfully.')
                getUsers()
            })
            .catch(function() {
                notifications.danger('Error attempting to remove user from department')
                eventTag('department_remove_user', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || $warrior.rank === 'PRIVATE') {
            router.route(appRoutes.login)
        }

        getDepartment()
    })

    $: isAdmin = organizationRole === 'ADMIN' || departmentRole === 'ADMIN'
</script>

<PageLayout>
    <h1 class="text-3xl font-bold">Department: {department.name}</h1>
    <div class="font-bold mb-4">
        Organization <ChevronRight class="inline-block" /> <a class="text-blue-500 hover:text-blue-800" href="{appRoutes.organization}/{organization.id}">{organization.name}</a>
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
                                    href="{appRoutes.organization}/{organizationId}/department/{departmentId}/team/{team.id}"
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
                            <td class="border px-4 py-2 text-right">
                                {#if isAdmin}
                                    <HollowButton onClick="{toggleRemoveUser(usr.id)}" color="red">
                                        Remove
                                    </HollowButton>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>

    {#if showCreateTeam}
        <CreateTeam toggleCreate="{toggleCreateTeam}" handleCreate={createTeamHandler} />
    {/if}

    {#if showAddUser}
        <AddUser toggleAdd="{toggleAddUser}" handleAdd={handleUserAdd} />
    {/if}

    {#if showRemoveUser}
        <RemoveUser toggleRemove={toggleRemoveUser(null)} handleRemove={handleUserRemove} />
    {/if}
</PageLayout>

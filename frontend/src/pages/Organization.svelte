<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CreateDepartment from '../components/user/CreateDepartment.svelte'
    import CreateTeam from '../components/user/CreateTeam.svelte'
    import AddUser from '../components/user/AddUser.svelte'
    import DeleteConfirmation from '../components/DeleteConfirmation.svelte'
    import CountryFlag from '../components/user/CountryFlag.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { appRoutes } from '../config.js'
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
            return
        }

        getOrganization()
    })

    $: isAdmin = role === 'ADMIN'
</script>

<svelte:head>
    <title>{$_('organization')} {organization.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
        <span class="uppercase">{$_('organization')}</span>
        <ChevronRight class="w-8 h-8" />
        {organization.name}
    </h1>

    <div class="w-full mb-6 lg:mb-8">
        <div class="flex w-full">
            <div class="w-4/5">
                <h2
                    class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                >
                    {$_('departments')}
                </h2>
            </div>
            <div class="w-1/5">
                <div class="text-right">
                    {#if isAdmin}
                        <SolidButton onClick="{toggleCreateDepartment}">
                            {$_('departmentCreate')}
                        </SolidButton>
                    {/if}
                </div>
            </div>
        </div>

        <div class="w-full">
            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
                        >
                            <table
                                class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                            >
                                <thead class="bg-gray-50 dark:bg-gray-800">
                                    <tr>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('name')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('dateCreated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('dateUpdated')}
                                        </th>
                                        <th></th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                                >
                                    {#each departments as department, i}
                                        <tr
                                            class:bg-slate-100="{i % 2 !== 0}"
                                            class:dark:bg-gray-800="{i % 2 !==
                                                0}"
                                        >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                <a
                                                    href="{appRoutes.organization}/{organizationId}/department/{department.id}"
                                                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                                >
                                                    {department.name}
                                                </a>
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    department.createdDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    department.updatedDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                            >
                                                {#if isAdmin}
                                                    <HollowButton
                                                        onClick="{toggleDeleteTeam(
                                                            department.id,
                                                        )}"
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
                </div>
            </div>
        </div>
    </div>

    <div class="w-full mb-6 lg:mb-8">
        <div class="flex w-full">
            <div class="w-4/5">
                <h2
                    class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                >
                    {$_('teams')}
                </h2>
            </div>
            <div class="w-1/5">
                <div class="text-right">
                    {#if isAdmin}
                        <SolidButton onClick="{toggleCreateTeam}">
                            {$_('teamCreate')}
                        </SolidButton>
                    {/if}
                </div>
            </div>
        </div>

        <div class="w-full">
            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
                        >
                            <table
                                class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                            >
                                <thead class="bg-gray-50 dark:bg-gray-800">
                                    <tr>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('name')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('dateCreated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                        >
                                            {$_('dateUpdated')}
                                        </th>
                                        <th></th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                                >
                                    {#each teams as team, i}
                                        <tr
                                            class:bg-slate-100="{i % 2 !== 0}"
                                            class:dark:bg-gray-800="{i % 2 !==
                                                0}"
                                        >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                <a
                                                    href="{appRoutes.organization}/{organizationId}/team/{team.id}"
                                                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                                >
                                                    {team.name}
                                                </a>
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    team.createdDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    team.updatedDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                            >
                                                {#if isAdmin}
                                                    <HollowButton
                                                        onClick="{toggleDeleteTeam(
                                                            team.id,
                                                        )}"
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
                </div>
            </div>
        </div>
    </div>

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

        <div class="flex flex-col">
            <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div
                    class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                >
                    <div
                        class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
                    >
                        <table
                            class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                        >
                            <thead class="bg-gray-50 dark:bg-gray-800">
                                <tr>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                    >
                                        {$_('name')}
                                    </th>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                    >
                                        {$_('email')}
                                    </th>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                    >
                                        {$_('role')}
                                    </th>
                                    <th scope="col" class="relative px-6 py-3">
                                        <span class="sr-only">Actions</span>
                                    </th>
                                </tr>
                            </thead>
                            <tbody
                                class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                            >
                                {#each users as user, i}
                                    <tr
                                        class:bg-slate-100="{i % 2 !== 0}"
                                        class:dark:bg-gray-800="{i % 2 !== 0}"
                                    >
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="flex items-center">
                                                <div
                                                    class="flex-shrink-0 h-10 w-10"
                                                >
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
                                                        class="font-medium text-gray-900"
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
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {user.email}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <div class="text-sm text-gray-500">
                                                {user.role}
                                            </div>
                                        </td>
                                        <td
                                            class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                        >
                                            {#if isAdmin}
                                                <HollowButton
                                                    onClick="{toggleRemoveUser(
                                                        user.id,
                                                    )}"
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
            </div>
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
        <DeleteConfirmation
            toggleDelete="{toggleRemoveUser(null)}"
            handleDelete="{handleUserRemove}"
            permanent="{false}"
            confirmText="{$_('removeUserConfirmText')}"
            confirmBtnText="{$_('removeUser')}"
        />
    {/if}

    {#if showDeleteTeam}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteTeam(null)}"
            handleDelete="{handleDeleteTeam}"
            confirmText="{$_('deleteTeamConfirmText')}"
            confirmBtnText="{$_('deleteTeam')}"
        />
    {/if}
</PageLayout>

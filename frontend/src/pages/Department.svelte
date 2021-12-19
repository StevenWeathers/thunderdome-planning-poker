<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CreateTeam from '../components/user/CreateTeam.svelte'
    import AddUser from '../components/user/AddUser.svelte'
    import DeleteConfirmation from '../components/DeleteConfirmation.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import CountryFlag from '../components/user/CountryFlag.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'

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
        name: '',
    }
    let department = {
        id: departmentId,
        name: '',
    }
    let departmentRole = ''
    let organizationRole = ''
    let teams = []
    let users = []
    let showCreateTeam = false
    let showAddUser = false
    let showRemoveUser = false
    let removeUserId = null
    let showDeleteTeam = false
    let deleteTeamId = null
    let teamsPage = 1
    let usersPage = 1

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

    function getDepartment() {
        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}`,
        )
            .then(res => res.json())
            .then(function (result) {
                department = result.data.department
                organization = result.data.organization
                organizationRole = result.data.organizationRole
                departmentRole = result.data.departmentRole

                getTeams()
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('departmentGetError'))
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                teams = result.data
            })
            .catch(function () {
                notifications.danger($_('departmentTeamsGetError'))
            })
    }

    function getUsers() {
        const usersOffset = (usersPage - 1) * usersPageLimit
        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/users?limit=${usersPageLimit}&offset=${usersOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                users = result.data
            })
            .catch(function () {
                notifications.danger($_('departmentUsersGetError'))
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/teams`,
            { body },
        )
            .then(res => res.json())
            .then(function () {
                eventTag('create_department_team', 'engagement', 'success')
                toggleCreateTeam()
                notifications.success($_('teamCreateSuccess'))
                getTeams()
            })
            .catch(function () {
                notifications.danger($_('teamCreateError'))
                eventTag('create_department_team', 'engagement', 'failure')
            })
    }

    function handleUserAdd(email, role) {
        const body = {
            email,
            role,
        }

        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/users`,
            { body },
        )
            .then(function () {
                eventTag('department_add_user', 'engagement', 'success')
                toggleAddUser()
                notifications.success($_('userAddSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userAddError'))
                eventTag('department_add_user', 'engagement', 'failure')
            })
    }

    function handleUserRemove() {
        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/users/${removeUserId}`,
            { method: 'DELETE' },
        )
            .then(function () {
                eventTag('department_remove_user', 'engagement', 'success')
                toggleRemoveUser(null)()
                notifications.success($_('userRemoveSuccess'))
                getUsers()
            })
            .catch(function () {
                notifications.danger($_('userRemoveError'))
                eventTag('department_remove_user', 'engagement', 'failure')
            })
    }

    function handleDeleteTeam() {
        xfetch(
            `/api/organizations/${organizationId}/departments/${departmentId}/teams/${deleteTeamId}`,
            { method: 'DELETE' },
        )
            .then(function () {
                eventTag('department_delete_team', 'engagement', 'success')
                toggleDeleteTeam(null)()
                notifications.success($_('teamDeleteSuccess'))
                getTeams()
            })
            .catch(function () {
                notifications.danger($_('teamDeleteError'))
                eventTag('department_delete_team', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || !validateUserIsRegistered($warrior)) {
            router.route(appRoutes.login)
            return
        }

        getDepartment()
    })

    $: isAdmin = organizationRole === 'ADMIN' || departmentRole === 'ADMIN'
</script>

<svelte:head>
    <title>{$_('department')} {department.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="mb-6 lg:mb-8 dark:text-white">
        <h1 class="text-3xl font-semibold font-rajdhani">
            <span class="uppercase">{$_('department')}</span>
            <ChevronRight class="w-8 h-8" />
            {department.name}
        </h1>
        <div class="text-xl font-semibold font-rajdhani">
            <span class="uppercase">{$_('organization')}</span>
            <ChevronRight />
            <a
                class="text-blue-500 hover:text-blue-800"
                href="{appRoutes.organization}/{organization.id}"
            >
                {organization.name}
            </a>
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
                                                    href="{appRoutes.organization}/{organizationId}/department/{departmentId}/team/{team.id}"
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
                                                class="px-6 py-4 whitespace-nowrap"
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
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                    >
                                        {$_('name')}
                                    </th>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                    >
                                        {$_('email')}
                                    </th>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
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

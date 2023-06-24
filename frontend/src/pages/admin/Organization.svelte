<script lang="ts">
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import DeleteConfirmation from '../../components/DeleteConfirmation.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import ChevronRight from '../../components/icons/ChevronRight.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.ts'
    import RowCol from '../../components/table/RowCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import Table from '../../components/table/Table.svelte'
    import { validateUserIsAdmin } from '../../validationUtils'

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
    let showDeleteDepartment = false
    let showDeleteOrganization = false
    let deleteTeamId = null
    let deleteDeptId = null
    let teamsPage = 1
    let usersPage = 1
    let departmentsPage = 1

    const toggleDeleteTeam = teamId => () => {
        showDeleteTeam = !showDeleteTeam
        deleteTeamId = teamId
    }

    const toggleDeleteDepartment = deptId => () => {
        showDeleteDepartment = !showDeleteDepartment
        deleteDeptId = deptId
    }

    const toggleDeleteOrganization = () => {
        showDeleteOrganization = !showDeleteOrganization
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

    function handleDeleteDepartment() {
        xfetch(
            `/api/organizations/${organizationId}/departments/${deleteDeptId}`,
            {
                method: 'DELETE',
            },
        )
            .then(function () {
                eventTag(
                    'organization_delete_department',
                    'engagement',
                    'success',
                )
                toggleDeleteDepartment(null)()
                notifications.success($_('departmentDeleteSuccess'))
                getDepartments()
            })
            .catch(function () {
                notifications.danger($_('departmentDeleteError'))
                eventTag(
                    'organization_delete_department',
                    'engagement',
                    'failure',
                )
            })
    }

    function handleDeleteOrganization() {
        xfetch(`/api/organizations/${organizationId}`, {
            method: 'DELETE',
        })
            .then(function () {
                eventTag('organization_delete', 'engagement', 'success')
                toggleDeleteTeam()
                notifications.success($_('organizationDeleteSuccess'))
                router.route(appRoutes.adminOrganizations)
            })
            .catch(function () {
                notifications.danger($_('organizationDeleteError'))
                eventTag('organization_delete', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
            return
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
            return
        }

        getOrganization()
    })
</script>

<svelte:head>
    <title>{$_('organization')} {organization.name} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="organizations">
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
        </div>

        <div class="w-full">
            <Table>
                <tr slot="header">
                    <HeadCol>
                        {$_('name')}
                    </HeadCol>
                    <HeadCol>
                        {$_('dateCreated')}
                    </HeadCol>
                    <HeadCol>
                        {$_('dateUpdated')}
                    </HeadCol>
                    <HeadCol type="action">
                        <span class="sr-only">{$_('actions')}</span>
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    {#each departments as department, i}
                        <TableRow itemIndex="{i}">
                            <RowCol>
                                <a
                                    href="{appRoutes.adminOrganizations}/{organizationId}/department/{department.id}"
                                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                >
                                    {department.name}
                                </a>
                            </RowCol>
                            <RowCol>
                                {new Date(
                                    department.createdDate,
                                ).toLocaleString()}
                            </RowCol>
                            <RowCol>
                                {new Date(
                                    department.updatedDate,
                                ).toLocaleString()}
                            </RowCol>
                            <RowCol type="action">
                                <HollowButton
                                    onClick="{toggleDeleteDepartment(
                                        department.id,
                                    )}"
                                    color="red"
                                >
                                    {$_('delete')}
                                </HollowButton>
                            </RowCol>
                        </TableRow>
                    {/each}
                </tbody>
            </Table>
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
        </div>

        <div class="w-full">
            <Table>
                <tr slot="header">
                    <HeadCol>
                        {$_('name')}
                    </HeadCol>
                    <HeadCol>
                        {$_('dateCreated')}
                    </HeadCol>
                    <HeadCol>
                        {$_('dateUpdated')}
                    </HeadCol>
                    <HeadCol type="action">
                        <span class="sr-only">{$_('actions')}</span>
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    {#each teams as team, i}
                        <TableRow itemIndex="{i}">
                            <RowCol>
                                <a
                                    href="{appRoutes.adminOrganizations}/{organizationId}/team/{team.id}"
                                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                >
                                    {team.name}
                                </a>
                            </RowCol>
                            <RowCol>
                                {new Date(team.createdDate).toLocaleString()}
                            </RowCol>
                            <RowCol>
                                {new Date(team.updatedDate).toLocaleString()}
                            </RowCol>
                            <RowCol type="action">
                                <HollowButton
                                    onClick="{toggleDeleteTeam(team.id)}"
                                    color="red"
                                >
                                    {$_('delete')}
                                </HollowButton>
                            </RowCol>
                        </TableRow>
                    {/each}
                </tbody>
            </Table>
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
                                <div class="ms-4">
                                    <div
                                        class="font-medium text-gray-900 dark:text-gray-200"
                                    >
                                        <a
                                            data-testid="user-name"
                                            href="{appRoutes.adminUsers}/{user.id}"
                                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                            >{user.name}</a
                                        >
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
                            <span data-testid="user-email">{user.email}</span>
                        </RowCol>
                        <RowCol>
                            <div
                                class="text-sm text-gray-500 dark:text-gray-300"
                            >
                                {user.role}
                            </div>
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>
    </div>

    <div class="w-full text-center mt-8">
        <HollowButton onClick="{toggleDeleteOrganization}" color="red">
            {$_('deleteOrganization')}
        </HollowButton>
    </div>

    {#if showDeleteTeam}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteTeam(null)}"
            handleDelete="{handleDeleteTeam}"
            confirmText="{$_('deleteTeamConfirmText')}"
            confirmBtnText="{$_('deleteTeam')}"
        />
    {/if}

    {#if showDeleteDepartment}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteDepartment(null)}"
            handleDelete="{handleDeleteDepartment}"
            confirmText="{$_('deleteDepartmentConfirmText')}"
            confirmBtnText="{$_('deleteDepartment')}"
        />
    {/if}

    {#if showDeleteOrganization}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteOrganization}"
            handleDelete="{handleDeleteOrganization}"
            confirmText="{$_('deleteOrganizationConfirmText')}"
            confirmBtnText="{$_('deleteOrganization')}"
        />
    {/if}
</AdminPageLayout>

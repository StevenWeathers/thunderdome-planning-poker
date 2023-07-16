<script lang="ts">
    import { onMount } from 'svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import { warrior } from '../../stores'
    import LL from '../../i18n/i18n-svelte'
    import { AppConfig, appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import RowCol from '../../components/table/RowCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import CheckboxIcon from '../../components/icons/CheckboxIcon.svelte'
    import CommentIcon from '../../components/icons/CommentIcon.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import ChevronRight from '../../components/icons/ChevronRight.svelte'
    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Table from '../../components/table/Table.svelte'
    import DeleteConfirmation from '../../components/DeleteConfirmation.svelte'

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig

    export let xfetch
    export let router
    export let notifications
    export let organizationId
    export let departmentId
    export let teamId

    const battlesPageLimit = 1000
    const retrosPageLimit = 1000
    const retroActionsPageLimit = 5
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
    let retroActions = []
    let storyboards = []
    let usersPage = 1
    let battlesPage = 1
    let retrosPage = 1
    let retroActionsPage = 1
    let storyboardsPage = 1
    let totalRetroActions = 0
    let completedActionItems = false

    let showDeleteTeam = false

    const apiPrefix = '/api'
    $: orgPrefix = departmentId
        ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
        : `${apiPrefix}/organizations/${organizationId}`
    $: teamPrefix = organizationId
        ? `${orgPrefix}/teams/${teamId}`
        : `${apiPrefix}/teams/${teamId}`

    let showRetroActionComments = false
    let selectedRetroAction = null
    const toggleRetroActionComments = id => () => {
        showRetroActionComments = !showRetroActionComments
        selectedRetroAction = id
    }

    function getTeam() {
        xfetch(teamPrefix)
            .then(res => res.json())
            .then(function (result) {
                team = result.data.team

                if (departmentId) {
                    department = result.data.department
                }
                if (organizationId) {
                    organization = result.data.organization
                }

                getBattles()
                getRetros()
                getRetrosActions()
                getStoryboards()
                getUsers()
            })
            .catch(function () {
                notifications.danger($LL.teamGetError())
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
                notifications.danger($LL.teamGetUsersError())
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
                    notifications.danger(
                        $LL.teamGetBattlesError({
                            friendly: AppConfig.FriendlyUIVerbs,
                        }),
                    )
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
                    notifications.danger($LL.teamGetRetrosError())
                })
        }
    }

    function getRetrosActions() {
        if (FeatureRetro) {
            const offset = (retroActionsPage - 1) * retroActionsPageLimit
            xfetch(
                `${teamPrefix}/retro-actions?limit=${retroActionsPageLimit}&offset=${offset}&completed=${completedActionItems}`,
            )
                .then(res => res.json())
                .then(function (result) {
                    retroActions = result.data
                    totalRetroActions = result.meta.count
                })
                .catch(function () {
                    notifications.danger($LL.teamGetRetroActionsError())
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
                    notifications.danger($LL.teamGetStoryboardsError())
                })
        }
    }

    function deleteTeam() {
        xfetch(`${teamPrefix}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(function () {
                if (departmentId) {
                    router.route(
                        `${appRoutes.adminOrganizations}/${organizationId}/department/${departmentId}`,
                    )
                } else if (organizationId) {
                    router.route(
                        `${appRoutes.adminOrganizations}/${organizationId}`,
                    )
                } else {
                    router.route(appRoutes.adminTeams)
                }
            })
            .catch(function () {
                notifications.danger($LL.teamDeleteError())
            })
    }

    const changeRetroActionPage = evt => {
        retroActionsPage = evt.detail
        getRetrosActions()
    }

    const changeRetroActionCompletedToggle = () => {
        retroActionsPage = 1
        getRetrosActions()
    }

    function toggleDeleteTeam() {
        showDeleteTeam = !showDeleteTeam
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

        getTeam()
    })
</script>

<svelte:head>
    <title>{$LL.team()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="teams">
    <div class="px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
        >
            {team.name}
        </h1>

        {#if organizationId}
            <div class="text-xl font-semibold font-rajdhani dark:text-white">
                <span class="uppercase">{$LL.organization()}</span>
                <ChevronRight />
                <a
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    href="{appRoutes.adminOrganizations}/{organization.id}"
                >
                    {organization.name}
                </a>
                {#if departmentId}
                    &nbsp;
                    <ChevronRight />
                    <span class="uppercase">{$LL.department()}</span>
                    <ChevronRight />
                    <a
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        href="{appRoutes.adminOrganizations}/{organization.id}/department/{department.id}"
                    >
                        {department.name}
                    </a>
                {/if}
            </div>
        {/if}
    </div>

    <div class="w-full">
        <div class="mb-4">
            <Table>
                <tr slot="header">
                    <HeadCol>
                        {$LL.dateCreated()}
                    </HeadCol>
                    <HeadCol>
                        {$LL.dateUpdated()}
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    <TableRow itemIndex="{0}">
                        <RowCol>
                            {new Date(team.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(team.updatedDate).toLocaleString()}
                        </RowCol>
                    </TableRow>
                </tbody>
            </Table>
        </div>
        <div>
            {#if FeaturePoker}
                <div class="w-full mb-6 lg:mb-8">
                    <div class="flex w-full">
                        <div class="flex-1">
                            <h2
                                class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                            >
                                {$LL.battles({
                                    friendly: AppConfig.FriendlyUIVerbs,
                                })}
                            </h2>
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
                                    <div
                                        class="w-full md:w-1/2 md:mb-0 md:text-right"
                                    >
                                        <HollowButton
                                            href="{appRoutes.game}/{battle.id}"
                                        >
                                            {$LL.battleJoin({
                                                friendly:
                                                    AppConfig.FriendlyUIVerbs,
                                            })}
                                        </HollowButton>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}

            {#if FeatureRetro}
                <div class="w-full mb-6 lg:mb-8">
                    <div class="flex w-full">
                        <div class="flex-1">
                            <h2
                                class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                            >
                                {$LL.retros()}
                            </h2>
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
                                    <div
                                        class="w-full md:w-1/2 md:mb-0 md:text-right"
                                    >
                                        <HollowButton
                                            href="{appRoutes.retro}/{retro.id}"
                                        >
                                            {$LL.joinRetro()}
                                        </HollowButton>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>

                    {#if retros.length}
                        <div class="w-full pt-4 px-4">
                            <div class="w-full">
                                <h3
                                    class="text-xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                                >
                                    {$LL.retroActionItems()}
                                </h3>

                                <div class="text-right mb-4">
                                    <div
                                        class="relative inline-block w-10 me-2 align-middle select-none transition duration-200 ease-in"
                                    >
                                        <input
                                            type="checkbox"
                                            name="completedActionItems"
                                            id="completedActionItems"
                                            bind:checked="{completedActionItems}"
                                            on:change="{changeRetroActionCompletedToggle}"
                                            class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                                        />
                                        <label
                                            for="completedActionItems"
                                            class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                                        >
                                        </label>
                                    </div>
                                    <label
                                        for="completedActionItems"
                                        class="dark:text-gray-300"
                                        >{$LL.showCompletedActionItems()}</label
                                    >
                                </div>
                            </div>

                            <Table>
                                <tr slot="header">
                                    <HeadCol>{$LL.actionItem()}</HeadCol>
                                    <HeadCol>{$LL.completed()}</HeadCol>
                                    <HeadCol>{$LL.comments()}</HeadCol>
                                </tr>
                                <tbody
                                    slot="body"
                                    let:class="{className}"
                                    class="{className}"
                                >
                                    {#each retroActions as item, i}
                                        <TableRow itemIndex="{i}">
                                            <RowCol>
                                                <div
                                                    class="whitespace-pre-wrap"
                                                >
                                                    {item.content}
                                                </div>
                                            </RowCol>
                                            <RowCol>
                                                <input
                                                    type="checkbox"
                                                    id="{i}Completed"
                                                    checked="{item.completed}"
                                                    class="opacity-0 absolute h-6 w-6"
                                                    disabled
                                                />
                                                <div
                                                    class="bg-white dark:bg-gray-800 border-2 rounded-md
                                            border-gray-400 dark:border-gray-300 w-6 h-6 flex flex-shrink-0
                                            justify-center items-center me-2
                                            focus-within:border-blue-500 dark:focus-within:border-sky-500"
                                                >
                                                    <CheckboxIcon />
                                                </div>
                                                <label
                                                    for="{i}Completed"
                                                    class="select-none"></label>
                                            </RowCol>
                                            <RowCol>
                                                <CommentIcon
                                                    width="22"
                                                    height="22"
                                                />
                                                <button
                                                    class="text-lg text-blue-400 dark:text-sky-400"
                                                    on:click="{toggleRetroActionComments(
                                                        item.id,
                                                    )}"
                                                >
                                                    &nbsp;{item.comments.length}
                                                </button>
                                            </RowCol>
                                        </TableRow>
                                    {/each}
                                </tbody>
                            </Table>

                            {#if totalRetroActions > retroActionsPageLimit}
                                <div class="pt-6 flex justify-center">
                                    <Pagination
                                        bind:current="{retroActionsPage}"
                                        num_items="{totalRetroActions}"
                                        per_page="{retroActionsPageLimit}"
                                        on:navigate="{changeRetroActionPage}"
                                    />
                                </div>
                            {/if}
                        </div>
                    {/if}
                </div>
            {/if}

            {#if FeatureStoryboard}
                <div class="w-full mb-6 lg:mb-8">
                    <div class="flex w-full">
                        <div class="flex-1">
                            <h2
                                class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                            >
                                {$LL.storyboards()}
                            </h2>
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
                                    <div
                                        class="w-full md:w-1/2 md:mb-0 md:text-right"
                                    >
                                        <HollowButton
                                            href="{appRoutes.storyboard}/{storyboard.id}"
                                        >
                                            {$LL.joinStoryboard()}
                                        </HollowButton>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}

            <div class="w-full">
                <div class="flex w-full">
                    <div class="w-4/5">
                        <h2
                            class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                        >
                            {$LL.users()}
                        </h2>
                    </div>
                </div>

                <Table>
                    <tr slot="header">
                        <HeadCol>
                            {$LL.name()}
                        </HeadCol>
                        <HeadCol>
                            {$LL.email()}
                        </HeadCol>
                        <HeadCol>
                            {$LL.role()}
                        </HeadCol>
                    </tr>
                    <tbody
                        slot="body"
                        let:class="{className}"
                        class="{className}"
                    >
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
                                    <span data-testid="user-email"
                                        >{user.email}</span
                                    >
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

            {#if !organizationId}
                <div class="text-center mt-4">
                    <HollowButton
                        color="red"
                        onClick="{toggleDeleteTeam}"
                        testid="team-delete"
                    >
                        {$LL.deleteTeam()}
                    </HollowButton>
                </div>
            {/if}

            {#if showDeleteTeam}
                <DeleteConfirmation
                    toggleDelete="{toggleDeleteTeam}"
                    handleDelete="{deleteTeam}"
                    confirmText="{$LL.deleteTeamConfirmText()}"
                    confirmBtnText="{$LL.deleteTeam()}"
                />
            {/if}
        </div>
    </div>
</AdminPageLayout>

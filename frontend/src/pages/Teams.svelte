<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CreateOrganization from '../components/user/CreateOrganization.svelte'
    import CreateTeam from '../components/user/CreateTeam.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'
    import RowCol from '../components/table/RowCol.svelte'
    import TableRow from '../components/table/TableRow.svelte'
    import HeadCol from '../components/table/HeadCol.svelte'
    import Table from '../components/table/Table.svelte'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const organizationsPageLimit = 1000
    const teamsPageLimit = 1000
    const { OrganizationsEnabled } = AppConfig

    let organizations = []
    let teams = []
    let showCreateOrganization = false
    let showCreateTeam = false
    let organizationsPage = 1
    let teamsPage = 1

    function toggleCreateOrganization() {
        showCreateOrganization = !showCreateOrganization
    }

    function toggleCreateTeam() {
        showCreateTeam = !showCreateTeam
    }

    function getOrganizations() {
        const orgsOffset = (organizationsPage - 1) * organizationsPageLimit
        xfetch(
            `/api/users/${$warrior.id}/organizations?limit=${organizationsPageLimit}&offset=${orgsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                organizations = result.data
            })
            .catch(function () {
                notifications.danger($_('getOrganizationsError'))
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/users/${$warrior.id}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                teams = result.data
            })
            .catch(function () {
                notifications.danger($_('getTeamsError'))
            })
    }

    function createOrganizationHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/users/${$warrior.id}/organizations`, { body })
            .then(res => res.json())
            .then(function (result) {
                eventTag('create_organization', 'engagement', 'success', () => {
                    router.route(`${appRoutes.organization}/${result.data.id}`)
                })
            })
            .catch(function () {
                notifications.danger($_('createOrgError'))
                eventTag('create_organization', 'engagement', 'failure')
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/users/${$warrior.id}/teams`, { body })
            .then(res => res.json())
            .then(function (result) {
                eventTag('create_team', 'engagement', 'success', () => {
                    router.route(`${appRoutes.team}/${result.data.id}`)
                })
            })
            .catch(function () {
                notifications.danger($_('teamCreateError'))
                eventTag('create_team', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || !validateUserIsRegistered($warrior)) {
            router.route(appRoutes.login)
            return
        }

        if (OrganizationsEnabled) {
            getOrganizations()
        }
        getTeams()
    })
</script>

<svelte:head>
    <title>{$_('organizations')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    {#if OrganizationsEnabled}
        <div class="w-full mb-6 lg:mb-8">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2
                        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                    >
                        {$_('organizations')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <SolidButton onClick="{toggleCreateOrganization}">
                            {$_('organizationCreate')}
                        </SolidButton>
                    </div>
                </div>
            </div>

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
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    {#each organizations as organization, i}
                        <TableRow itemIndex="{i}">
                            <RowCol>
                                <a
                                    href="{appRoutes.organization}/{organization.id}"
                                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                >
                                    {organization.name}
                                </a>
                            </RowCol>
                            <RowCol>
                                {new Date(
                                    organization.createdDate,
                                ).toLocaleString()}
                            </RowCol>
                            <RowCol>
                                {new Date(
                                    organization.updatedDate,
                                ).toLocaleString()}
                            </RowCol>
                        </TableRow>
                    {/each}
                </tbody>
            </Table>
        </div>
    {/if}

    <div class="w-full">
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
                    <SolidButton onClick="{toggleCreateTeam}">
                        {$_('teamCreate')}
                    </SolidButton>
                </div>
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
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    {#each teams as team, i}
                        <TableRow itemIndex="{i}">
                            <RowCol>
                                <a
                                    href="{appRoutes.team}/{team.id}"
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
                        </TableRow>
                    {/each}
                </tbody>
            </Table>
        </div>
    </div>

    {#if showCreateOrganization}
        <CreateOrganization
            toggleCreate="{toggleCreateOrganization}"
            handleCreate="{createOrganizationHandler}"
        />
    {/if}

    {#if showCreateTeam}
        <CreateTeam
            toggleCreate="{toggleCreateTeam}"
            handleCreate="{createTeamHandler}"
        />
    {/if}
</PageLayout>

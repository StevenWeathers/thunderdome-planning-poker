<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CreateOrganization from '../components/user/CreateOrganization.svelte'
    import CreateTeam from '../components/user/CreateTeam.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const organizationsPageLimit = 1000
    const teamsPageLimit = 1000

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

        getOrganizations()
        getTeams()
    })
</script>

<svelte:head>
    <title>{$_('organizations')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
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
                                </tr>
                            </thead>
                            <tbody
                                class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                            >
                                {#each organizations as org, i}
                                    <tr
                                        class:bg-slate-100="{i % 2 !== 0}"
                                        class:dark:bg-gray-800="{i % 2 !== 0}"
                                    >
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <a
                                                href="{appRoutes.organization}/{org.id}"
                                                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                            >
                                                {org.name}
                                            </a>
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                org.createdDate,
                                            ).toLocaleString()}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                org.updatedDate,
                                            ).toLocaleString()}
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
                                                    href="/team/{team.id}"
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

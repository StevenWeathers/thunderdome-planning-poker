<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import CreateOrganization from '../components/CreateOrganization.svelte'
    import CreateTeam from '../components/CreateTeam.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

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
            .then(function(result) {
                organizations = result.data
            })
            .catch(function() {
                notifications.danger('Error getting organizations')
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(
            `/api/users/${$warrior.id}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
        )
            .then(res => res.json())
            .then(function(result) {
                teams = result.data
            })
            .catch(function() {
                notifications.danger('Error getting teams')
            })
    }

    function createOrganizationHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/users/${$warrior.id}/organizations`, { body })
            .then(res => res.json())
            .then(function(result) {
                eventTag('create_organization', 'engagement', 'success', () => {
                    router.route(`${appRoutes.organization}/${result.data.id}`)
                })
            })
            .catch(function() {
                notifications.danger('Error attempting to create organization')
                eventTag('create_organization', 'engagement', 'failure')
            })
    }

    function createTeamHandler(name) {
        const body = {
            name,
        }

        xfetch(`/api/users/${$warrior.id}/teams`, { body })
            .then(res => res.json())
            .then(function(result) {
                eventTag('create_team', 'engagement', 'success', () => {
                    router.route(`${appRoutes.team}/${result.data.id}`)
                })
            })
            .catch(function() {
                notifications.danger('Error attempting to create team')
                eventTag('create_team', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || $warrior.rank === 'PRIVATE') {
            router.route(appRoutes.login)
        }

        getOrganizations()
        getTeams()
    })
</script>

<svelte:head>
    <title>{$_('organizations')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('organizations')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <HollowButton onClick="{toggleCreateOrganization}">
                            {$_('organizationCreate')}
                        </HollowButton>
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
                    {#each organizations as org}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="{appRoutes.organization}/{org.id}"
                                    class="text-blue-500 hover:text-blue-800">
                                    {org.name}
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
                        {$_('teams')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <HollowButton onClick="{toggleCreateTeam}">
                            {$_('teamCreate')}
                        </HollowButton>
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
                    {#each teams as team}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="/team/{team.id}"
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

    {#if showCreateOrganization}
        <CreateOrganization
            toggleCreate="{toggleCreateOrganization}"
            handleCreate="{createOrganizationHandler}" />
    {/if}

    {#if showCreateTeam}
        <CreateTeam
            toggleCreate="{toggleCreateTeam}"
            handleCreate="{createTeamHandler}" />
    {/if}
</PageLayout>

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
        xfetch(`/api/organizations/${organizationsPageLimit}/${orgsOffset}`)
            .then(res => res.json())
            .then(function(result) {
                organizations = result
            })
            .catch(function(error) {
                notifications.danger('Error getting organizations')
            })
    }

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(`/api/teams/${teamsPageLimit}/${teamsOffset}`)
            .then(res => res.json())
            .then(function(result) {
                teams = result
            })
            .catch(function(error) {
                notifications.danger('Error getting teams')
            })
    }

    function createOrganizationHandler(name) {
        const body = {
            name,
        }

        xfetch('/api/organizations', { body })
            .then(res => res.json())
            .then(function(organization) {
                eventTag('create_organization', 'engagement', 'success', () => {
                    router.route(`${appRoutes.organization}/${organization.id}`)
                })
            })
            .catch(function(error) {
                notifications.danger('Error attempting to create organization')
                eventTag('create_organization', 'engagement', 'failure')
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

<PageLayout>
    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold text-center mb-4">
                        Organizations
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <HollowButton onClick="{toggleCreateOrganization}">
                            Create Organization
                        </HollowButton>
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
                    {#each organizations as org}
                        <tr>
                            <td class="border px-4 py-2">
                                <a
                                    href="/organization/{org.id}"
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
                    <h2 class="text-2xl md:text-3xl font-bold text-center mb-4">
                        Teams
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <HollowButton onClick="{toggleCreateTeam}">
                            Create Team
                        </HollowButton>
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
        <CreateTeam toggleCreate="{toggleCreateTeam}" />
    {/if}
</PageLayout>

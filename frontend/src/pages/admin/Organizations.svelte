<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    // export let eventTag

    const organizationsPageLimit = 100

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        battleCount: 0,
        planCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
    }
    let organizations = []
    let organizationsPage = 1

    function getAppStats() {
        xfetch('/api/admin/stats')
            .then(res => res.json())
            .then(function (result) {
                appStats = result.data
            })
            .catch(function () {
                notifications.danger($_('applicationStatsError'))
            })
    }

    function getOrganizations() {
        const organizationsOffset =
            (organizationsPage - 1) * organizationsPageLimit
        xfetch(
            `/api/admin/organizations?limit=${organizationsPageLimit}&offset=${organizationsOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                organizations = result.data
            })
            .catch(function () {
                notifications.danger($_('getOrganizationsError'))
            })
    }

    const changePage = evt => {
        organizationsPage = evt.detail
        getOrganizations()
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

        getAppStats()
        getOrganizations()
    })
</script>

<svelte:head>
    <title>
        {$_('organizations')}
        {$_('pages.admin.title')} | {$_('appName')}
    </title>
</svelte:head>

<AdminPageLayout activePage="organizations">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase">
            {$_('organizations')}
        </h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"
                        >
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead class="bg-gray-50">
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
                                            {$_('dateCreated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('dateUpdated')}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    {#each organizations as org, i}
                                        <tr class:bg-slate-100="{i % 2 !== 0}">
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {org.name}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    org.createdDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
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

            {#if appStats.organizationCount > organizationsPageLimit}
                <div class="pt-6 flex justify-center">
                    <Pagination
                        bind:current="{organizationsPage}"
                        num_items="{appStats.organizationCount}"
                        per_page="{organizationsPageLimit}"
                        on:navigate="{changePage}"
                    />
                </div>
            {/if}
        </div>
    </div>
</AdminPageLayout>

<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    // export let eventTag

    const apikeysPageLimit = 100

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        battleCount: 0,
        planCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
        apikeyCount: 0,
    }
    let apikeys = []
    let apikeysPage = 1

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

    function getApiKeys() {
        const apikeysOffset = (apikeysPage - 1) * apikeysPageLimit
        xfetch(
            `/api/admin/apikeys?limit=${apikeysPageLimit}&offset=${apikeysOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                apikeys = result.data
            })
            .catch(function () {
                notifications.danger($_('getApikeysError'))
            })
    }

    const changePage = evt => {
        apikeysPage = evt.detail
        getApiKeys()
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
        getApiKeys()
    })
</script>

<svelte:head>
    <title>{$_('apiKeys')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="apikeys">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('apiKeys')}
        </h1>
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
                                        {$_('prefix')}
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
                                        {$_('active')}
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
                                {#each apikeys as apikey, i}
                                    <tr
                                        class:bg-slate-100="{i % 2 !== 0}"
                                        class:dark:bg-gray-800="{i % 2 !== 0}"
                                    >
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {apikey.name}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {apikey.prefix}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {apikey.userId}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {#if apikey.active}
                                                <span class="text-green-600"
                                                    ><CheckIcon /></span
                                                >
                                            {/if}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                apikey.createdDate,
                                            ).toLocaleString()}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                apikey.updatedDate,
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

        {#if appStats.apikeyCount > apikeysPageLimit}
            <div class="pt-6 flex justify-center">
                <Pagination
                    bind:current="{apikeysPage}"
                    num_items="{appStats.apikeyCount}"
                    per_page="{apikeysPageLimit}"
                    on:navigate="{changePage}"
                />
            </div>
        {/if}
    </div>
</AdminPageLayout>

<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

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
            .then(function(result) {
                appStats = result
            })
            .catch(function(error) {
                notifications.danger('Error getting application stats')
            })
    }

    function getApiKeys() {
        const apikeysOffset = (apikeysPage - 1) * apikeysPageLimit
        xfetch(`/api/admin/apikeys/${apikeysPageLimit}/${apikeysOffset}`)
            .then(res => res.json())
            .then(function(result) {
                apikeys = result
            })
            .catch(function(error) {
                notifications.danger('Error getting apikeys')
            })
    }

    const changePage = evt => {
        apikeysPage = evt.detail
        getApiKeys()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
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
        <h1 class="text-3xl md:text-4xl font-bold">{$_('apiKeys')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-3/12 px-4 py-2">{$_('name')}</th>
                        <th class="w-1/12 px-4 py-2">{$_('prefix')}</th>
                        <th class="w-2/12 px-4 py-2">{$_('email')}</th>
                        <th class="w-1/12 px-4 py-2">{$_('active')}</th>
                        <th class="w-2/12 px-4 py-2">{$_('dateCreated')}</th>
                        <th class="w-2/12 px-4 py-2">{$_('dateUpdated')}</th>
                        <th class="w-1/12 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each apikeys as apikey}
                        <tr>
                            <td class="border px-4 py-2">{apikey.name}</td>
                            <td class="border px-4 py-2">{apikey.prefix}</td>
                            <td class="border px-4 py-2">{apikey.warriorId}</td>
                            <td class="border px-4 py-2">{apikey.active}</td>
                            <td class="border px-4 py-2">
                                {new Date(apikey.createdDate).toLocaleString()}
                            </td>
                            <td class="border px-4 py-2">
                                {new Date(apikey.updatedDate).toLocaleString()}
                            </td>
                            <td class="border px-4 py-2"></td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            {#if appStats.apikeyCount > apikeysPageLimit}
                <div class="pt-6 flex justify-center">
                    <Pagination
                        bind:current="{apikeysPage}"
                        num_items="{appStats.apikeyCount}"
                        per_page="{apikeysPageLimit}"
                        on:navigate="{changePage}" />
                </div>
            {/if}
        </div>
    </div>
</AdminPageLayout>

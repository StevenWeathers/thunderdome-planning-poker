<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import RowCol from '../../components/table/RowCol.svelte'

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
        <Table>
            <tr slot="header">
                <HeadCol>
                    {$_('name')}
                </HeadCol>
                <HeadCol>
                    {$_('prefix')}
                </HeadCol>
                <HeadCol>
                    {$_('email')}
                </HeadCol>
                <HeadCol>
                    {$_('active')}
                </HeadCol>
                <HeadCol>
                    {$_('dateCreated')}
                </HeadCol>
                <HeadCol>
                    {$_('dateUpdated')}
                </HeadCol>
            </tr>
            <tbody slot="body" let:class="{className}" class="{className}">
                {#each apikeys as apikey, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            {apikey.name}
                        </RowCol>
                        <RowCol>
                            {apikey.prefix}
                        </RowCol>
                        <RowCol>
                            {apikey.userId}
                        </RowCol>
                        <RowCol>
                            {#if apikey.active}
                                <span class="text-green-600"><CheckIcon /></span
                                >
                            {/if}
                        </RowCol>
                        <RowCol>
                            {new Date(apikey.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(apikey.updatedDate).toLocaleString()}
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>

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

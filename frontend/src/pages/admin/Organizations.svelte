<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
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
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('organizations')}
        </h1>
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
                {#each organizations as org, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            {org.name}
                        </RowCol>
                        <RowCol>
                            {new Date(org.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(org.updatedDate).toLocaleString()}
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>

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
</AdminPageLayout>

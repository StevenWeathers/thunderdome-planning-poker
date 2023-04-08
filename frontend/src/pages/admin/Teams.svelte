<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import RowCol from '../../components/table/RowCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import Table from '../../components/table/Table.svelte'

    export let xfetch
    export let router
    export let notifications

    const teamsPageLimit = 100

    let teamCount = 0
    let teams = []
    let teamsPage = 1

    function getTeams() {
        const teamsOffset = (teamsPage - 1) * teamsPageLimit
        xfetch(`/api/admin/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`)
            .then(res => res.json())
            .then(function (result) {
                teams = result.data
                teamCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getTeamsError'))
            })
    }

    const changePage = evt => {
        teamsPage = evt.detail
        getTeams()
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

        getTeams()
    })
</script>

<svelte:head>
    <title>{$_('teams')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="teams">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('teams')}
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
                {#each teams as team, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            <a
                                href="{appRoutes.adminTeams}/{team.id}"
                                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                >{team.name}</a
                            >
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

        {#if teamCount > teamsPageLimit}
            <div class="pt-6 flex justify-center">
                <Pagination
                    bind:current="{teamsPage}"
                    num_items="{teamCount}"
                    per_page="{teamsPageLimit}"
                    on:navigate="{changePage}"
                />
            </div>
        {/if}
    </div>
</AdminPageLayout>

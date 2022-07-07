<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
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

    const battlesPageLimit = 100
    let battleCount = 0
    let battles = []
    let battlesPage = 1
    let activeBattles = false

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit
        xfetch(
            `/api/battles?limit=${battlesPageLimit}&offset=${battlesOffset}&active=${activeBattles}`,
        )
            .then(res => res.json())
            .then(function (result) {
                battles = result.data
                battleCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getBattlesError'))
            })
    }

    const changePage = evt => {
        battlesPage = evt.detail
        getBattles()
    }

    const changeActiveBattlesToggle = () => {
        battlesPage = 1
        getBattles()
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

        getBattles()
    })
</script>

<style>
    .toggle-checkbox:checked {
        @apply right-0;
        @apply border-green-400;
        border-color: #68d391;
    }

    .toggle-checkbox:checked + .toggle-label {
        @apply bg-green-400;
        background-color: #68d391;
    }
</style>

<svelte:head>
    <title>{$_('battles')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="battles">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('battles')}
        </h1>
    </div>

    <div class="w-full">
        <div class="text-right mb-4">
            <div
                class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in"
            >
                <input
                    type="checkbox"
                    name="activeBattles"
                    id="activeBattles"
                    bind:checked="{activeBattles}"
                    on:change="{changeActiveBattlesToggle}"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                    for="activeBattles"
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
            </div>
            <label for="activeBattles" class="dark:text-gray-300"
                >{$_('showActiveBattles')}</label
            >
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
                <HeadCol type="action">
                    <span class="sr-only">{$_('actions')}</span>
                </HeadCol>
            </tr>
            <tbody slot="body" let:class="{className}" class="{className}">
                {#each battles as battle, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            <a
                                href="{appRoutes.admin}/battles/{battle.id}"
                                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                >{battle.name}</a
                            >
                        </RowCol>
                        <RowCol>
                            {new Date(battle.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(battle.updatedDate).toLocaleString()}
                        </RowCol>
                        <RowCol type="action">
                            <HollowButton href="{appRoutes.battle}/{battle.id}">
                                {$_('battleJoin')}
                            </HollowButton>
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>

        {#if battleCount > battlesPageLimit}
            <div class="pt-6 flex justify-center">
                <Pagination
                    bind:current="{battlesPage}"
                    num_items="{battleCount}"
                    per_page="{battlesPageLimit}"
                    on:navigate="{changePage}"
                />
            </div>
        {/if}
    </div>
</AdminPageLayout>

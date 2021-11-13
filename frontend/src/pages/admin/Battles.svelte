<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const battlesPageLimit = 100
    let battleCount = 0
    let battles = []
    let battlesPage = 1

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit
        xfetch(`/api/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`)
            .then(res => res.json())
            .then(function (result) {
                battles = result.data
                battleCount = result.meta.count
            })
            .catch(function () {
                notifications.danger('Error getting battles')
            })
    }

    const changePage = evt => {
        battlesPage = evt.detail
        getBattles()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
        }

        getBattles()
    })
</script>

<svelte:head>
    <title>{$_('battles')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="battles">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('battles')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-3/12 px-4 py-2">{$_('name')}</th>
                        <th class="w-3/12 px-4 py-2">{$_('dateCreated')}</th>
                        <th class="w-3/12 px-4 py-2">{$_('dateUpdated')}</th>
                        <th class="w-3/12 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each battles as battle}
                        <tr>
                            <td class="border px-4 py-2">{battle.name}</td>
                            <td class="border px-4 py-2"
                                >{new Date(
                                    battle.createdDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border px-4 py-2"
                                >{new Date(
                                    battle.updatedDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border px-4 py-2 text-right">
                                <HollowButton
                                    color="blue"
                                    href="{appRoutes.admin}/battles/{battle.id}"
                                >
                                    Details
                                </HollowButton>
                                <HollowButton
                                    href="{appRoutes.battle}/{battle.id}"
                                >
                                    Join
                                </HollowButton>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

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
    </div>
</AdminPageLayout>

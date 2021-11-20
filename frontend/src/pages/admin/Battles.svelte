<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

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
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
        }

        getBattles()
    })
</script>

<style>
    :global(.toggle-checkbox:checked) {
        @apply right-0;
        @apply border-green-400;
        border-color: #68d391;
    }

    :global(.toggle-checkbox:checked + .toggle-label) {
        @apply bg-green-400;
        background-color: #68d391;
    }
</style>

<svelte:head>
    <title>{$_('battles')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="battles">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('battles')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="text-right">
                <div
                    class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in"
                >
                    <input
                        type="checkbox"
                        name="activeBattles"
                        id="activeBattles"
                        bind:checked="{activeBattles}"
                        on:change="{getBattles}"
                        class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                    />
                    <label
                        for="activeBattles"
                        class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                    >
                    </label>
                </div>
                <label for="activeBattles">Show Active Battles</label>
            </div>
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('name')}</th>
                        <th class="flex-1 p-2">{$_('dateCreated')}</th>
                        <th class="flex-1 p-2">{$_('dateUpdated')}</th>
                        <th class="flex-1 p-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each battles as battle}
                        <tr>
                            <td class="border p-2">
                                <a
                                    href="{appRoutes.admin}/battles/{battle.id}"
                                    class="no-underline text-blue-500 hover:text-blue-800"
                                    >{battle.name}</a
                                >
                            </td>
                            <td class="border p-2"
                                >{new Date(
                                    battle.createdDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border p-2"
                                >{new Date(
                                    battle.updatedDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border p-2 text-right">
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

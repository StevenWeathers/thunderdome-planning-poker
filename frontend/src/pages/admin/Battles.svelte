<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

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
                    on:change="{getBattles}"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                    for="activeBattles"
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
            </div>
            <label for="activeBattles">{$_('showActiveBattles')}</label>
        </div>

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
                                        {$_('dateCreated')}
                                    </th>
                                    <th
                                        scope="col"
                                        class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                    >
                                        {$_('dateUpdated')}
                                    </th>
                                    <th scope="col" class="relative px-6 py-3">
                                        <span class="sr-only">Actions</span>
                                    </th>
                                </tr>
                            </thead>
                            <tbody
                                class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                            >
                                {#each battles as battle, i}
                                    <tr
                                        class:bg-slate-100="{i % 2 !== 0}"
                                        class:dark:bg-gray-800="{i % 2 !== 0}"
                                    >
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            <a
                                                href="{appRoutes.admin}/battles/{battle.id}"
                                                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                                >{battle.name}</a
                                            >
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                battle.createdDate,
                                            ).toLocaleString()}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {new Date(
                                                battle.updatedDate,
                                            ).toLocaleString()}
                                        </td>
                                        <td
                                            class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                        >
                                            <HollowButton
                                                href="{appRoutes.battle}/{battle.id}"
                                            >
                                                {$_('battleJoin')}
                                            </HollowButton>
                                        </td>
                                    </tr>
                                {/each}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

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

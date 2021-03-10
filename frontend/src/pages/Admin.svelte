<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import CreateWarrior from '../components/CreateWarrior.svelte'
    import Pagination from '../components/Pagination.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const warriorsPageLimit = 100

    let appStats = {
        unregisteredWarriorCount: 0,
        registeredWarriorCount: 0,
        battleCount: 0,
        planCount: 0,
    }
    let warriors = []
    let showCreateWarrior = false
    let warriorsPage = 1

    function toggleCreateWarrior() {
        showCreateWarrior = !showCreateWarrior
    }

    function createWarrior(
        warriorName,
        warriorEmail,
        warriorPassword1,
        warriorPassword2,
    ) {
        const body = {
            warriorName,
            warriorEmail,
            warriorPassword1,
            warriorPassword2,
        }

        xfetch('/api/admin/warrior', { body })
            .then(function() {
                eventTag('admin_create_warrior', 'engagement', 'success')

                getWarriors()
                toggleCreateWarrior()
            })
            .catch(function(error) {
                notifications.danger('Error encountered creating warrior')
                eventTag('admin_create_warrior', 'engagement', 'failure')
            })
    }

    xfetch('/api/admin/stats')
        .then(res => res.json())
        .then(function(result) {
            appStats = result
        })
        .catch(function(error) {
            notifications.danger('Error getting application stats')
        })

    function getWarriors() {
        const warriorsOffset = (warriorsPage - 1) * warriorsPageLimit
        xfetch(`/api/admin/warriors/${warriorsPageLimit}/${warriorsOffset}`)
            .then(res => res.json())
            .then(function(result) {
                warriors = result
            })
            .catch(function(error) {
                notifications.danger('Error getting warriors')
            })
    }

    function promoteWarrior(warriorId) {
        return function() {
            const body = {
                warriorId,
            }

            xfetch('/api/admin/promote', { body })
                .then(function() {
                    eventTag('admin_promote_warrior', 'engagement', 'success')

                    getWarriors()
                })
                .catch(function(error) {
                    notifications.danger('Error encountered promoting warrior')
                    eventTag('admin_promote_warrior', 'engagement', 'failure')
                })
        }
    }

    function demoteWarrior(warriorId) {
        return function() {
            const body = {
                warriorId,
            }

            xfetch('/api/admin/demote', { body })
                .then(function() {
                    eventTag('admin_demote_warrior', 'engagement', 'success')

                    getWarriors()
                })
                .catch(function(error) {
                    notifications.danger('Error encountered demoting warrior')
                    eventTag('admin_demote_warrior', 'engagement', 'failure')
                })
        }
    }

    const changePage = evt => {
        warriorsPage = evt.detail
        getWarriors()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.register)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
        }

        getWarriors()
    })
</script>

<PageLayout>
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">
            {$_('pages.admin.title')}
        </h1>
    </div>
    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="w-1/4">
                    <div class="mb-2 font-bold">
                        {$_('pages.admin.counts.unregistered')}
                    </div>
                    {appStats.unregisteredWarriorCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">
                        {$_('pages.admin.counts.registered')}
                    </div>
                    {appStats.registeredWarriorCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">
                        {$_('pages.admin.counts.battles')}
                    </div>
                    {appStats.battleCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">
                        {$_('pages.admin.counts.plans')}
                    </div>
                    {appStats.planCount}
                </div>
            </div>
        </div>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-4/5">
                    <h2 class="text-2xl md:text-3xl font-bold text-center mb-4">
                        {$_('pages.admin.registeredWarriors.title')}
                    </h2>
                </div>
                <div class="w-1/5">
                    <div class="text-right">
                        <HollowButton onClick="{toggleCreateWarrior}">
                            {$_('actions.warrior.create')}
                        </HollowButton>
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-2/6 px-4 py-2">
                            {$_('pages.admin.registeredWarriors.name')}
                        </th>
                        <th class="w-2/6 px-4 py-2">
                            {$_('pages.admin.registeredWarriors.email')}
                        </th>
                        <th class="w-1/6 px-4 py-2">
                            {$_('pages.admin.registeredWarriors.verified')}
                        </th>
                        <th class="w-1/6 px-4 py-2">
                            {$_('pages.admin.registeredWarriors.rank')}
                        </th>
                        <th class="w-1/6 px-4 py-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each warriors as warrior}
                        <tr>
                            <td class="border px-4 py-2">{warrior.name}</td>
                            <td class="border px-4 py-2">{warrior.email}</td>
                            <td class="border px-4 py-2">{warrior.verified}</td>
                            <td class="border px-4 py-2">{warrior.rank}</td>
                            <td class="border px-4 py-2">
                                {#if warrior.rank !== 'GENERAL'}
                                    <HollowButton
                                        onClick="{promoteWarrior(warrior.id)}"
                                        color="blue">
                                        {$_('actions.warrior.promote')}
                                    </HollowButton>
                                {:else}
                                    <HollowButton
                                        onClick="{demoteWarrior(warrior.id)}"
                                        color="blue">
                                        {$_('actions.warrior.demote')}
                                    </HollowButton>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            {#if appStats.registeredWarriorCount > warriorsPageLimit}
            <div class="pt-6 flex justify-center">
                <Pagination bind:current={warriorsPage} num_items={appStats.registeredWarriorCount} per_page={warriorsPageLimit} on:navigate={changePage} />
            </div>
            {/if}
        </div>
    </div>

    {#if showCreateWarrior}
        <CreateWarrior
            toggleCreate="{toggleCreateWarrior}"
            handleCreate="{createWarrior}"
            notifications />
    {/if}
</PageLayout>

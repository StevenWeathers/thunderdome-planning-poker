<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import CreateWarrior from '../../components/CreateWarrior.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CountryFlag from '../../components/CountryFlag.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const warriorsPageLimit = 100

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        battleCount: 0,
        planCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
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
            router.route(appRoutes.login)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
        }

        getAppStats()
        getWarriors()
    })
</script>

<AdminPageLayout activePage="users">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('users')}</h1>
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
                            {$_('warriorCreate')}
                        </HollowButton>
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="w-3/12 p-2">
                            {$_('pages.admin.registeredWarriors.name')}
                        </th>
                        <th class="w-3/12 p-2">
                            {$_('pages.admin.registeredWarriors.email')}
                        </th>
                        <th class="w-3/12 p-2">
                            {$_('pages.admin.registeredWarriors.company')}
                        </th>
                        <th class="w-2/12 p-2">
                            {$_('pages.admin.registeredWarriors.rank')}
                        </th>
                        <th class="w-1/12 p-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each warriors as warrior}
                        <tr>
                            <td class="border p-2">
                                {warrior.name}
                                {#if warrior.country}
                                    &nbsp;
                                    <CountryFlag
                                        country="{warrior.country}"
                                        size="{16}"
                                        additionalClass="inline-block" />
                                {/if}
                            </td>
                            <td class="border p-2">
                                {warrior.email}
                                {#if warrior.verified}
                                    &nbsp;
                                    <span
                                        class="text-green-600"
                                        title="{$_('pages.admin.registeredWarriors.verified')}">
                                        <CheckIcon />
                                    </span>
                                {/if}
                            </td>
                            <td class="border p-2">
                                <div>{warrior.company}</div>
                                {#if warrior.jobTitle}
                                    <div class="text-gray-700 text-sm">
                                        {$_('pages.admin.registeredWarriors.jobTitle')}:
                                        {warrior.jobTitle}
                                    </div>
                                {/if}
                            </td>
                            <td class="border p-2">{warrior.rank}</td>
                            <td class="border p-2">
                                {#if warrior.rank !== 'GENERAL'}
                                    <HollowButton
                                        onClick="{promoteWarrior(warrior.id)}"
                                        color="blue">
                                        {$_('promote')}
                                    </HollowButton>
                                {:else}
                                    <HollowButton
                                        onClick="{demoteWarrior(warrior.id)}"
                                        color="blue">
                                        {$_('demote')}
                                    </HollowButton>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            {#if appStats.registeredUserCount > warriorsPageLimit}
                <div class="pt-6 flex justify-center">
                    <Pagination
                        bind:current="{warriorsPage}"
                        num_items="{appStats.registeredUserCount}"
                        per_page="{warriorsPageLimit}"
                        on:navigate="{changePage}" />
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
</AdminPageLayout>

<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const { CleanupGuestsDaysOld, CleanupBattlesDaysOld, APIEnabled } = appConfig

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        battleCount: 0,
        planCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
        apikeyCount: 0
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

    function cleanBattles() {
        xfetch('/api/admin/clean-battles', { method: 'DELETE' })
            .then(function() {
                eventTag('admin_clean_battles', 'engagement', 'success')

                getAppStats()
            })
            .catch(function(error) {
                notifications.danger('Error encountered cleaning battles')
                eventTag('admin_clean_battles', 'engagement', 'failure')
            })
    }

    function cleanGuests() {
        xfetch('/api/admin/clean-guests', { method: 'DELETE' })
            .then(function() {
                eventTag('admin_clean_guests', 'engagement', 'success')

                getAppStats()
            })
            .catch(function(error) {
                notifications.danger('Error encountered cleaning guests')
                eventTag('admin_clean_guests', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
        }

        getAppStats()
    })
</script>

<AdminPageLayout activePage="admin">
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
                    {appStats.unregisteredUserCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">
                        {$_('pages.admin.counts.registered')}
                    </div>
                    {appStats.registeredUserCount}
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
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Organizations</div>
                    {appStats.organizationCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Departments</div>
                    {appStats.departmentCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Teams</div>
                    {appStats.teamCount}
                </div>
                <div class="w-1/4">
                    {#if APIEnabled}
                        <div class="mb-2 font-bold">API Keys</div>
                        {appStats.apikeyCount}
                    {/if}
                </div>
            </div>
        </div>
    </div>

    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="text-center p-2 md:p-4 bg-white shadow-lg rounded text-xl">
                <div class="text-2xl md:text-3xl font-bold text-center mb-4">
                    {$_('pages.admin.maintenance.title')}
                </div>
                <HollowButton onClick="{cleanGuests}" color="red">
                    {$_('pages.admin.maintenance.cleanGuests', {
                        values: { daysOld: CleanupGuestsDaysOld },
                    })}
                </HollowButton>

                <HollowButton onClick="{cleanBattles}" color="red">
                    {$_('pages.admin.maintenance.cleanBattles', {
                        values: { daysOld: CleanupBattlesDaysOld },
                    })}
                </HollowButton>
            </div>
        </div>
    </div>
</AdminPageLayout>

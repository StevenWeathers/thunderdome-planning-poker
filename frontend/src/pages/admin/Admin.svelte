<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { AppConfig, appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const { CleanupGuestsDaysOld, CleanupBattlesDaysOld, ExternalAPIEnabled } =
        AppConfig

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        battleCount: 0,
        planCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
        apikeyCount: 0,
        activeBattleCount: 0,
        activeBattleUserCount: 0,
    }

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

    function cleanBattles() {
        xfetch('/api/maintenance/clean-battles', { method: 'DELETE' })
            .then(function () {
                eventTag('admin_clean_battles', 'engagement', 'success')

                getAppStats()
            })
            .catch(function () {
                notifications.danger($_('oldBattleCleanError'))
                eventTag('admin_clean_battles', 'engagement', 'failure')
            })
    }

    function cleanGuests() {
        xfetch('/api/maintenance/clean-guests', { method: 'DELETE' })
            .then(function () {
                eventTag('admin_clean_guests', 'engagement', 'success')

                getAppStats()
            })
            .catch(function () {
                notifications.danger($_('oldGuestsCleanError'))
                eventTag('admin_clean_guests', 'engagement', 'failure')
            })
    }

    function lowercaseEmails() {
        xfetch('/api/maintenance/lowercase-emails', { method: 'PATCH' })
            .then(function () {
                eventTag('admin_lowercase_emails', 'engagement', 'success')
                notifications.success($_('lowercaseEmailsSuccess'))

                getAppStats()
            })
            .catch(function () {
                notifications.danger($_('lowercaseEmailsError'))
                eventTag('admin_lowercase_emails', 'engagement', 'failure')
            })
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
    })
</script>

<svelte:head>
    <title>{$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="admin">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase">
            {$_('pages.admin.title')}
        </h1>
    </div>
    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl"
            >
                <div class="w-1/2">
                    <div class="mb-2 font-bold">{$_('battlesActive')}</div>
                    {appStats.activeBattleCount}
                </div>
                <div class="w-1/2">
                    <div class="mb-2 font-bold">{$_('battlesActiveUsers')}</div>
                    {appStats.activeBattleUserCount}
                </div>
            </div>
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl"
            >
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
                md:pb-4 bg-white shadow-lg rounded text-xl"
            >
                <div class="{ExternalAPIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">{$_('organizations')}</div>
                    {appStats.organizationCount}
                </div>
                <div class="{ExternalAPIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">{$_('departments')}</div>
                    {appStats.departmentCount}
                </div>
                <div class="{ExternalAPIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">{$_('teams')}</div>
                    {appStats.teamCount}
                </div>
                {#if ExternalAPIEnabled}
                    <div class="w-1/4">
                        <div class="mb-2 font-bold">{$_('apiKeys')}</div>
                        {appStats.apikeyCount}
                    </div>
                {/if}
            </div>
        </div>
    </div>

    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="text-center p-2 md:p-4 bg-white shadow-lg rounded text-xl"
            >
                <div
                    class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase text-center mb-4"
                >
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

                <HollowButton onClick="{lowercaseEmails}" color="red">
                    {$_('maintenanceLowercaseEmails')}
                </HollowButton>
            </div>
        </div>
    </div>
</AdminPageLayout>

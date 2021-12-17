<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import UserIcon from '../../components/icons/UserIcon.svelte'
    import UserRankGuest from '../../components/icons/UserRankGuest.svelte'
    import UserRankRegistered from '../../components/icons/UserRankRegistered.svelte'
    import KeyIcon from '../../components/icons/Key.svelte'
    import LightingBolt from '../../components/icons/LightningBolt.svelte'
    import UserGroupIcon from '../../components/icons/UserGroup.svelte'
    import UsersIcon from '../../components/icons/Users.svelte'
    import OfficeBuildingIcon from '../../components/icons/OfficeBuilding.svelte'
    import DocumentTextIcon from '../../components/icons/DocumentText.svelte'
    import ShieldExclamationIcon from '../../components/icons/ShieldExclamation.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { AppConfig, appRoutes } from '../../config.js'
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
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-8">
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-blue-400 text-white">
                        <UserRankGuest width="28" height="28" />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('pages.admin.counts.unregistered')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.unregisteredUserCount}
                    </h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-indigo-500 text-white">
                        <UserRankRegistered width="28" height="28" />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('pages.admin.counts.registered')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.registeredUserCount}
                    </h3>
                </div>
            </div>
        </div>
        {#if ExternalAPIEnabled}
            <div class="bg-white border rounded shadow-lg p-2">
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-cyan-500 text-white">
                            <KeyIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5 class="font-bold uppercase text-gray-500">
                            {$_('apiKeys')}
                        </h5>
                        <h3 class="font-bold text-3xl">
                            {appStats.apikeyCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-orange-500 text-white">
                        <ShieldExclamationIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('pages.admin.counts.battles')}
                    </h5>
                    <h3 class="font-bold text-3xl">{appStats.battleCount}</h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-teal-500 text-white">
                        <DocumentTextIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('pages.admin.counts.plans')}
                    </h5>
                    <h3 class="font-bold text-3xl">{appStats.planCount}</h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-red-500 text-white">
                        <LightingBolt />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('battlesActive')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.activeBattleCount}
                    </h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-green-500 text-white">
                        <UserIcon width="28" height="28" />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('battlesActiveUsers')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.activeBattleUserCount}
                    </h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-sky-500 text-white">
                        <OfficeBuildingIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('organizations')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.organizationCount}
                    </h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-rose-500 text-white">
                        <UserGroupIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('departments')}
                    </h5>
                    <h3 class="font-bold text-3xl">
                        {appStats.departmentCount}
                    </h3>
                </div>
            </div>
        </div>
        <div class="bg-white border rounded shadow-lg p-2">
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-purple-500 text-white">
                        <UsersIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5 class="font-bold uppercase text-gray-500">
                        {$_('teams')}
                    </h5>
                    <h3 class="font-bold text-3xl">{appStats.teamCount}</h3>
                </div>
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

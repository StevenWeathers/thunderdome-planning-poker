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
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import CheckCircle from '../../components/icons/CheckCircle.svelte'
    import SmileCircle from '../../components/icons/SmileCircle.svelte'
    import QuestionCircle from '../../components/icons/QuestionCircle.svelte'
    import FrownCircle from '../../components/icons/FrownCircle.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { AppConfig, appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const {
        CleanupGuestsDaysOld,
        CleanupBattlesDaysOld,
        CleanupRetrosDaysOld,
        CleanupStoryboardsDaysOld,
        ExternalAPIEnabled,
        FeaturePoker,
        FeatureRetro,
        FeatureStoryboard,
        OrganizationsEnabled,
    } = AppConfig

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
        teamCheckinsCount: 0,
        retroCount: 0,
        activeRetroCount: 0,
        activeRetroUserCount: 0,
        retroItemCount: 0,
        retroActionCount: 0,
        storyboardCount: 0,
        activeStoryboardCount: 0,
        activeStoryboardUserCount: 0,
        storyboardGoalCount: 0,
        storyboardColumnCount: 0,
        storyboardStoryCount: 0,
        storyboardPersonaCount: 0,
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

    function cleanRetros() {
        xfetch('/api/maintenance/clean-retros', { method: 'DELETE' })
            .then(function () {
                eventTag('admin_clean_retros', 'engagement', 'success')

                getAppStats()
            })
            .catch(function () {
                notifications.danger($_('oldRetrosCleanError'))
                eventTag('admin_clean_retros', 'engagement', 'failure')
            })
    }

    function cleanStoryboards() {
        xfetch('/api/maintenance/clean-storyboards', { method: 'DELETE' })
            .then(function () {
                eventTag('admin_clean_storyboards', 'engagement', 'success')

                getAppStats()
            })
            .catch(function () {
                notifications.danger($_('oldStoryboardsCleanError'))
                eventTag('admin_clean_storyboards', 'engagement', 'failure')
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
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('pages.admin.title')}
        </h1>
    </div>
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-8">
        <div
            class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
        >
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-blue-400 text-white">
                        <UserRankGuest width="28" height="28" />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5
                        class="font-bold uppercase text-gray-500 dark:text-gray-400"
                    >
                        {$_('pages.admin.counts.unregistered')}
                    </h5>
                    <h3 class="font-bold text-3xl dark:text-white">
                        {appStats.unregisteredUserCount}
                    </h3>
                </div>
            </div>
        </div>
        <div
            class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
        >
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-indigo-500 text-white">
                        <UserRankRegistered width="28" height="28" />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5
                        class="font-bold uppercase text-gray-500 dark:text-gray-400"
                    >
                        {$_('pages.admin.counts.registered')}
                    </h5>
                    <h3 class="font-bold text-3xl dark:text-white">
                        {appStats.registeredUserCount}
                    </h3>
                </div>
            </div>
        </div>
        {#if ExternalAPIEnabled}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-cyan-500 text-white">
                            <KeyIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('apiKeys')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.apikeyCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
        {#if FeaturePoker}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-red-500 text-white">
                            <ShieldExclamationIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('pages.admin.counts.battles')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.battleCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-teal-500 text-white">
                            <DocumentTextIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('pages.admin.counts.plans')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.planCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-yellow-500 text-white">
                            <LightingBolt />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('battlesActive')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeBattleCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-green-500 text-white">
                            <UserIcon width="28" height="28" />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('battlesActiveUsers')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeBattleUserCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
        {#if OrganizationsEnabled}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-orange-500 text-white">
                            <OfficeBuildingIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('organizations')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.organizationCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-rose-500 text-white">
                            <UserGroupIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('departments')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.departmentCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
        <div
            class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
        >
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-purple-500 text-white">
                        <UsersIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5
                        class="font-bold uppercase text-gray-500 dark:text-gray-400"
                    >
                        {$_('teams')}
                    </h5>
                    <h3 class="font-bold text-3xl dark:text-white">
                        {appStats.teamCount}
                    </h3>
                </div>
            </div>
        </div>
        <div
            class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
        >
            <div class="flex flex-row items-center">
                <div class="flex-shrink pr-4">
                    <div class="rounded p-3 bg-lime-500 text-white">
                        <CheckIcon />
                    </div>
                </div>
                <div class="flex-1 text-right md:text-center">
                    <h5
                        class="font-bold uppercase text-gray-500 dark:text-gray-400"
                    >
                        {$_('teamCheckins')}
                    </h5>
                    <h3 class="font-bold text-3xl dark:text-white">
                        {appStats.teamCheckinsCount}
                    </h3>
                </div>
            </div>
        </div>
        {#if FeatureRetro}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-fuchsia-500 text-white">
                            <FrownCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('retros')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.retroCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-amber-500 text-white">
                            <QuestionCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('activeRetros')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeRetroCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-pink-500 text-white">
                            <UserIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('activeRetroUsers')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeRetroUserCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-emerald-500 text-white">
                            <SmileCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('retroItems')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.retroItemCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-violet-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('retroActionItems')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.retroActionCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
        {#if FeatureStoryboard}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-fuchsia-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('storyboards')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.storyboardCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-amber-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('activeStoryboards')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeStoryboardCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-pink-500 text-white">
                            <UserIcon />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('activeStoryboardUsers')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.activeStoryboardUserCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-emerald-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('storyboardGoals')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.storyboardGoalCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-violet-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('storyboardColumns')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.storyboardColumnCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-violet-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('storyboardStories')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.storyboardStoryCount}
                        </h3>
                    </div>
                </div>
            </div>
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-shrink pr-4">
                        <div class="rounded p-3 bg-violet-500 text-white">
                            <CheckCircle />
                        </div>
                    </div>
                    <div class="flex-1 text-right md:text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400"
                        >
                            {$_('storyboardPersonas')}
                        </h5>
                        <h3 class="font-bold text-3xl dark:text-white">
                            {appStats.storyboardPersonaCount}
                        </h3>
                    </div>
                </div>
            </div>
        {/if}
    </div>

    <div class="w-full">
        <div
            class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase text-center mb-4 dark:text-white"
        >
            {$_('pages.admin.maintenance.title')}
        </div>

        <div class="grid grid-cols-3 gap-4">
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-1 text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400 mb-2"
                        >
                            {$_('pages.admin.maintenance.cleanGuests', {
                                values: { daysOld: CleanupGuestsDaysOld },
                            })}
                        </h5>
                        <HollowButton onClick="{cleanGuests}" color="red">
                            {$_('execute')}
                        </HollowButton>
                    </div>
                </div>
            </div>

            {#if FeaturePoker}
                <div
                    class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
                >
                    <div class="flex flex-row items-center">
                        <div class="flex-1 text-center">
                            <h5
                                class="font-bold uppercase text-gray-500 dark:text-gray-400 mb-2"
                            >
                                {$_('pages.admin.maintenance.cleanBattles', {
                                    values: { daysOld: CleanupBattlesDaysOld },
                                })}
                            </h5>
                            <HollowButton onClick="{cleanBattles}" color="red">
                                {$_('execute')}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/if}
            {#if FeatureRetro}
                <div
                    class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
                >
                    <div class="flex flex-row items-center">
                        <div class="flex-1 text-center">
                            <h5
                                class="font-bold uppercase text-gray-500 dark:text-gray-400 mb-2"
                            >
                                {$_('adminCleanOldRetros', {
                                    values: { daysOld: CleanupRetrosDaysOld },
                                })}
                            </h5>
                            <HollowButton onClick="{cleanRetros}" color="red">
                                {$_('execute')}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/if}
            {#if FeatureStoryboard}
                <div
                    class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
                >
                    <div class="flex flex-row items-center">
                        <div class="flex-1 text-center">
                            <h5
                                class="font-bold uppercase text-gray-500 dark:text-gray-400 mb-2"
                            >
                                {$_('adminCleanOldStoryboards', {
                                    values: {
                                        daysOld: CleanupStoryboardsDaysOld,
                                    },
                                })}
                            </h5>
                            <HollowButton
                                onClick="{cleanStoryboards}"
                                color="red"
                            >
                                {$_('execute')}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/if}
            <div
                class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
            >
                <div class="flex flex-row items-center">
                    <div class="flex-1 text-center">
                        <h5
                            class="font-bold uppercase text-gray-500 dark:text-gray-400 mb-2"
                        >
                            {$_('maintenanceLowercaseEmails')}
                        </h5>
                        <HollowButton onClick="{lowercaseEmails}" color="red">
                            {$_('execute')}
                        </HollowButton>
                    </div>
                </div>
            </div>
        </div>
    </div>
</AdminPageLayout>

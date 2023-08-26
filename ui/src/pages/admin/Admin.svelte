<script lang="ts">
  import { onMount } from 'svelte';

  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import UserIcon from '../../components/icons/UserIcon.svelte';
  import UserRankRegisteredIcon from '../../components/icons/UserRankRegisteredIcon.svelte';
  import LightingBoltIcon from '../../components/icons/LightningBoltIcon.svelte';
  import UsersIcon from '../../components/icons/UsersIcon.svelte';
  import OfficeBuildingIcon from '../../components/icons/OfficeBuildingIcon.svelte';
  import ShieldExclamationIcon from '../../components/icons/ShieldExclamationIcon.svelte';
  import CheckCircleIcon from '../../components/icons/CheckCircleIcon.svelte';
  import FrownCircleIcon from '../../components/icons/FrownCircleIcon.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import SmileCircleIcon from '../../components/icons/SmileCircleIcon.svelte';
  import QuestionCircleIcon from '../../components/icons/QuestionCircleIcon.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import UserGroupIcon from '../../components/icons/UserGroupIcon.svelte';
  import DocumentTextIcon from '../../components/icons/DocumentTextIcon.svelte';
  import KeyIcon from '../../components/icons/KeyIcon.svelte';
  import UserRankGuestIcon from '../../components/icons/UserRankGuestIcon.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;

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
  } = AppConfig;

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
  };

  function getAppStats() {
    xfetch('/api/admin/stats')
      .then(res => res.json())
      .then(function (result) {
        appStats = result.data;
      })
      .catch(function () {
        notifications.danger($LL.applicationStatsError());
      });
  }

  function cleanBattles() {
    xfetch('/api/maintenance/clean-battles', { method: 'DELETE' })
      .then(function () {
        eventTag('admin_clean_battles', 'engagement', 'success');

        getAppStats();
      })
      .catch(function () {
        notifications.danger(
          $LL.oldBattleCleanError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
        eventTag('admin_clean_battles', 'engagement', 'failure');
      });
  }

  function cleanRetros() {
    xfetch('/api/maintenance/clean-retros', { method: 'DELETE' })
      .then(function () {
        eventTag('admin_clean_retros', 'engagement', 'success');

        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldRetrosCleanError());
        eventTag('admin_clean_retros', 'engagement', 'failure');
      });
  }

  function cleanStoryboards() {
    xfetch('/api/maintenance/clean-storyboards', { method: 'DELETE' })
      .then(function () {
        eventTag('admin_clean_storyboards', 'engagement', 'success');

        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldStoryboardsCleanError());
        eventTag('admin_clean_storyboards', 'engagement', 'failure');
      });
  }

  function cleanGuests() {
    xfetch('/api/maintenance/clean-guests', { method: 'DELETE' })
      .then(function () {
        eventTag('admin_clean_guests', 'engagement', 'success');

        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldGuestsCleanError());
        eventTag('admin_clean_guests', 'engagement', 'failure');
      });
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getAppStats();
  });
</script>

<svelte:head>
  <title>{$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="admin">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.admin()}
    </h1>
  </div>
  <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-8">
    <div
      class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded shadow-lg p-2"
    >
      <div class="flex flex-row items-center">
        <div class="flex-shrink pe-4">
          <div class="rounded p-3 bg-blue-400 text-white">
            <UserRankGuestIcon width="28" height="28" />
          </div>
        </div>
        <div class="flex-1 text-right md:text-center">
          <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
            {$LL.guestUsers()}
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
        <div class="flex-shrink pe-4">
          <div class="rounded p-3 bg-indigo-500 text-white">
            <UserRankRegisteredIcon width="28" height="28" />
          </div>
        </div>
        <div class="flex-1 text-right md:text-center">
          <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
            {$LL.registeredUsers()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-cyan-500 text-white">
              <KeyIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.apiKeys()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-red-500 text-white">
              <ShieldExclamationIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.battles({
                friendly: AppConfig.FriendlyUIVerbs,
              })}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-teal-500 text-white">
              <DocumentTextIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.plans({
                friendly: AppConfig.FriendlyUIVerbs,
              })}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-yellow-500 text-white">
              <LightingBoltIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.battlesActive({
                friendly: AppConfig.FriendlyUIVerbs,
              })}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-green-500 text-white">
              <UserIcon width="28" height="28" />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.battlesActiveUsers({
                friendly: AppConfig.FriendlyUIVerbs,
              })}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-orange-500 text-white">
              <OfficeBuildingIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.organizations()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-rose-500 text-white">
              <UserGroupIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.departments()}
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
        <div class="flex-shrink pe-4">
          <div class="rounded p-3 bg-purple-500 text-white">
            <UsersIcon />
          </div>
        </div>
        <div class="flex-1 text-right md:text-center">
          <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
            {$LL.teams()}
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
        <div class="flex-shrink pe-4">
          <div class="rounded p-3 bg-lime-500 text-white">
            <CheckIcon />
          </div>
        </div>
        <div class="flex-1 text-right md:text-center">
          <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
            {$LL.teamCheckins()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-fuchsia-500 text-white">
              <FrownCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.retros()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-amber-500 text-white">
              <QuestionCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.activeRetros()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-pink-500 text-white">
              <UserIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.activeRetroUsers()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-emerald-500 text-white">
              <SmileCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.retroItems()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-violet-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.retroActionItems()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-fuchsia-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.storyboards()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-amber-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.activeStoryboards()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-pink-500 text-white">
              <UserIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.activeStoryboardUsers()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-emerald-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.storyboardGoals()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-violet-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.storyboardColumns()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-violet-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.storyboardStories()}
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
          <div class="flex-shrink pe-4">
            <div class="rounded p-3 bg-violet-500 text-white">
              <CheckCircleIcon />
            </div>
          </div>
          <div class="flex-1 text-right md:text-center">
            <h5 class="font-bold uppercase text-gray-500 dark:text-gray-400">
              {$LL.storyboardPersonas()}
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
      {$LL.maintenance()}
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
              {$LL.cleanGuests({
                daysOld: CleanupGuestsDaysOld,
              })}
            </h5>
            <HollowButton onClick="{cleanGuests}" color="red">
              {$LL.execute()}
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
                {$LL.cleanBattles[AppConfig.FriendlyUIVerbs]({
                  daysOld: CleanupBattlesDaysOld,
                })}
              </h5>
              <HollowButton onClick="{cleanBattles}" color="red">
                {$LL.execute()}
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
                {$LL.adminCleanOldRetros({
                  daysOld: CleanupRetrosDaysOld,
                })}
              </h5>
              <HollowButton onClick="{cleanRetros}" color="red">
                {$LL.execute()}
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
                {$LL.adminCleanOldStoryboards({
                  daysOld: CleanupStoryboardsDaysOld,
                })}
              </h5>
              <HollowButton onClick="{cleanStoryboards}" color="red">
                {$LL.execute()}
              </HollowButton>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</AdminPageLayout>

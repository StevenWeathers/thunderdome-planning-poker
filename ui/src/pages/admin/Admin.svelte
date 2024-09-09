<script lang="ts">
  import { onMount } from 'svelte';

  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import {
    BarChart2,
    Building,
    CircleUser,
    Columns2,
    FileText,
    Gamepad,
    Ghost,
    Key,
    LayoutDashboard,
    LibraryBig,
    MessageCircleQuestion,
    Network,
    RefreshCcw,
    Smile,
    SquareCheckBig,
    User,
    UserRound,
    Users,
    Zap,
  } from 'lucide-svelte';

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
    estimationScaleCount: 0,
    publicEstimationScaleCount: 0,
    organizationEstimationScaleCount: 0,
    teamEstimationScaleCount: 0,
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
        notifications.danger($LL.oldBattleCleanError());
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

  $: statGroups = [
    {
      title: $LL.users(),
      bgColor: 'bg-indigo-500',
      stats: [
        {
          name: 'guestUsers',
          count: appStats.unregisteredUserCount,
          icon: Ghost,
          active: true,
        },
        {
          name: 'registeredUsers',
          count: appStats.registeredUserCount,
          icon: CircleUser,
          active: true,
        },
        {
          name: 'apiKeys',
          count: appStats.apikeyCount,
          icon: Key,
          active: ExternalAPIEnabled,
        },
      ],
    },
    {
      title: $LL.organizations(),
      bgColor: 'bg-orange-500 dark:bg-orange-400',
      stats: [
        {
          name: 'organizations',
          count: appStats.organizationCount,
          icon: Building,
          active: OrganizationsEnabled,
        },
        {
          name: 'departments',
          count: appStats.departmentCount,
          icon: Network,
          active: OrganizationsEnabled,
        },
        {
          name: 'teams',
          count: appStats.teamCount,
          icon: Users,
          active: true,
        },
      ],
    },
    {
      title: $LL.battles(),
      bgColor: 'bg-red-500',
      stats: [
        {
          name: 'battles',
          count: appStats.battleCount,
          icon: Gamepad,
          active: FeaturePoker,
        },
        {
          name: 'plans',
          count: appStats.planCount,
          icon: FileText,
          active: FeaturePoker,
        },
        {
          name: 'battlesActive',
          count: appStats.activeBattleCount,
          icon: Zap,
          active: FeaturePoker,
        },
        {
          name: 'battlesActiveUsers',
          count: appStats.activeBattleUserCount,
          icon: User,
          active: FeaturePoker,
        },
      ],
    },
    {
      title: $LL.retros(),
      bgColor: 'bg-green-500 dark:bg-lime-400',
      stats: [
        {
          name: 'retros',
          count: appStats.retroCount,
          icon: RefreshCcw,
          active: FeatureRetro,
        },
        {
          name: 'retroItems',
          count: appStats.retroItemCount,
          icon: Smile,
          active: FeatureRetro,
        },
        {
          name: 'retroActionItems',
          count: appStats.retroActionCount,
          icon: SquareCheckBig,
          active: FeatureRetro,
        },
        {
          name: 'activeRetros',
          count: appStats.activeRetroCount,
          icon: MessageCircleQuestion,
          active: FeatureRetro,
        },
        {
          name: 'activeRetroUsers',
          count: appStats.activeRetroUserCount,
          icon: User,
          active: FeatureRetro,
        },
      ],
    },
    {
      title: $LL.storyboards(),
      bgColor: 'bg-emerald-500 dark:bg-emerald-400',
      stats: [
        {
          name: 'storyboards',
          count: appStats.storyboardCount,
          icon: LayoutDashboard,
          active: FeatureStoryboard,
        },
        {
          name: 'storyboardGoals',
          count: appStats.storyboardGoalCount,
          icon: SquareCheckBig,
          active: FeatureStoryboard,
        },
        {
          name: 'storyboardColumns',
          count: appStats.storyboardColumnCount,
          icon: Columns2,
          active: FeatureStoryboard,
        },
        {
          name: 'storyboardStories',
          count: appStats.storyboardStoryCount,
          icon: LibraryBig,
          active: FeatureStoryboard,
        },
        {
          name: 'storyboardPersonas',
          count: appStats.storyboardPersonaCount,
          icon: UserRound,
          active: FeatureStoryboard,
        },
        {
          name: 'activeStoryboards',
          count: appStats.activeStoryboardCount,
          icon: Zap,
          active: FeatureStoryboard,
        },
        {
          name: 'activeStoryboardUsers',
          count: appStats.activeStoryboardUserCount,
          icon: User,
          active: FeatureStoryboard,
        },
      ],
    },
    {
      title: $LL.estimationScales(),
      bgColor: 'bg-yellow-500',
      stats: [
        {
          name: 'estimationScales',
          count: appStats.estimationScaleCount,
          icon: BarChart2,
          active: FeaturePoker,
        },
        {
          name: 'publicEstimationScales',
          count: appStats.publicEstimationScaleCount,
          icon: BarChart2,
          active: FeaturePoker,
        },
        {
          name: 'organizationEstimationScales',
          count: appStats.organizationEstimationScaleCount,
          icon: BarChart2,
          active: FeaturePoker,
        },
        {
          name: 'teamEstimationScales',
          count: appStats.teamEstimationScaleCount,
          icon: BarChart2,
          active: FeaturePoker,
        },
      ],
    },
  ];
</script>

<svelte:head>
  <title>{$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="admin">
  <div class="grid grid-cols-2 gap-4 mb-4">
    {#each statGroups as group}
      {#if group.stats.some(stat => stat.active)}
        <div
          class="bg-white dark:bg-gray-800 border dark:border-gray-700 rounded-lg shadow-lg p-4"
        >
          <h2
            class="text-2xl font-semibold mb-2 dark:text-white font-semibold font-rajdhani uppercase"
          >
            {group.title}
          </h2>
          <div class="grid grid-cols-2 gap-4">
            {#if group.stats.filter(stat => stat.active).length > 0}
              {#each group.stats.filter(stat => stat.active) as stat}
                <div
                  class="bg-gray-100 dark:bg-gray-700 rounded-lg p-2 transition-all duration-300 hover:shadow-md hover:scale-105"
                >
                  <div class="flex gap-2 content-center">
                    <div class="flex-none items-center content-center">
                      <div
                        class="w-12 h-12 justify-center text-center content-center rounded-full {group.bgColor} text-white mr-2"
                      >
                        <svelte:component
                          this="{stat.icon}"
                          width="26"
                          height="26"
                          class="mx-auto"
                        />
                      </div>
                    </div>
                    <div class="flex-grow">
                      <h3 class="text-lg font-medium dark:text-gray-200 mb-2">
                        {$LL[stat.name]()}
                      </h3>
                      <p class="text-3xl font-bold dark:text-white">
                        {stat.count}
                      </p>
                    </div>
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        </div>
      {/if}
    {/each}
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
                {$LL.cleanBattles({
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

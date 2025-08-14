<script lang="ts">
  import { onMount } from 'svelte';

  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import {
    ChartNoAxesColumn,
    Building,
    CircleUser,
    Columns2,
    CreditCard,
    FileText,
    Ghost,
    Key,
    LayoutDashboard,
    LibraryBig,
    MapPinCheckInside,
    MessageCircleQuestion,
    Network,
    RefreshCcw,
    Smile,
    SquareCheckBig,
    SquareDashedKanban,
    User,
    UserRound,
    Users,
    Vote,
    Zap,
    Package,
    Activity,
  } from 'lucide-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const {
    CleanupGuestsDaysOld,
    CleanupBattlesDaysOld,
    CleanupRetrosDaysOld,
    CleanupStoryboardsDaysOld,
    ExternalAPIEnabled,
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    FeatureProject,
    OrganizationsEnabled,
  } = AppConfig;

  let appStats = $state({
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
    userSubscriptionActiveCount: 0,
    teamSubscriptionActiveCount: 0,
    orgSubscriptionActiveCount: 0,
    publicRetroTemplateCount: 0,
    retroTemplateCount: 0,
    organizationRetroTemplateCount: 0,
    teamRetroTemplateCount: 0,
    projectCount: 0,
  });

  let isLoading = $state(false);

  function getAppStats() {
    isLoading = true;
    xfetch('/api/admin/stats')
      .then(res => res.json())
      .then(function (result) {
        appStats = result.data;
      })
      .catch(function () {
        notifications.danger($LL.applicationStatsError());
      })
      .finally(() => {
        isLoading = false;
      });
  }

  function cleanBattles() {
    xfetch('/api/maintenance/clean-battles', { method: 'DELETE' })
      .then(function () {
        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldBattleCleanError());
      });
  }

  function cleanRetros() {
    xfetch('/api/maintenance/clean-retros', { method: 'DELETE' })
      .then(function () {
        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldRetrosCleanError());
      });
  }

  function cleanStoryboards() {
    xfetch('/api/maintenance/clean-storyboards', { method: 'DELETE' })
      .then(function () {
        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldStoryboardsCleanError());
      });
  }

  function cleanGuests() {
    xfetch('/api/maintenance/clean-guests', { method: 'DELETE' })
      .then(function () {
        getAppStats();
      })
      .catch(function () {
        notifications.danger($LL.oldGuestsCleanError());
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

  // Calculate totals for overview cards
  let overviewStats = $derived([
    {
      title: 'Total Users',
      value: appStats.registeredUserCount + appStats.unregisteredUserCount,
      icon: Users,
      color: 'from-blue-500 to-blue-600',
      active: true
    },
    {
      title: 'Active Sessions',
      value: appStats.activeBattleUserCount + appStats.activeRetroUserCount + appStats.activeStoryboardUserCount,
      icon: Activity,
      color: 'from-green-500 to-green-600',
      active: true
    },
    {
      title: 'Total Subscriptions',
      value: appStats.userSubscriptionActiveCount + appStats.teamSubscriptionActiveCount + appStats.orgSubscriptionActiveCount,
      icon: CreditCard,
      color: 'from-indigo-500 to-indigo-600',
      active: AppConfig.SubscriptionsEnabled
    }
  ]);

  let statGroups = $derived([
    {
      title: $LL.users(),
      active: true,
      gradient: 'from-indigo-500 to-indigo-600',
      iconBg: 'bg-gradient-to-r from-indigo-500 to-indigo-600',
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
        {
          name: 'activeSubscriptions',
          count: appStats.userSubscriptionActiveCount,
          icon: CreditCard,
          active: AppConfig.SubscriptionsEnabled,
        },
      ],
    },
    {
      title: $LL.organizations(),
      active: OrganizationsEnabled,
      gradient: 'from-orange-500 to-orange-600',
      iconBg: 'bg-gradient-to-r from-orange-500 to-orange-600',
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
          name: 'activeSubscriptions',
          count: appStats.orgSubscriptionActiveCount,
          icon: CreditCard,
          active: AppConfig.SubscriptionsEnabled,
        },
      ],
    },
    {
      title: $LL.teams(),
      active: true,
      gradient: 'from-blue-500 to-sky-500',
      iconBg: 'bg-gradient-to-r from-blue-500 to-sky-500',
      stats: [
        {
          name: 'teams',
          count: appStats.teamCount,
          icon: Users,
          active: true,
        },
        {
          name: 'teamCheckins',
          count: appStats.teamCheckinsCount,
          icon: MapPinCheckInside,
          active: true,
        },
        {
          name: 'activeSubscriptions',
          count: appStats.teamSubscriptionActiveCount,
          icon: CreditCard,
          active: AppConfig.SubscriptionsEnabled,
        },
      ],
    },
    {
      title: $LL.battles(),
      active: FeaturePoker,
      gradient: 'from-red-500 to-red-600',
      iconBg: 'bg-gradient-to-r from-red-500 to-red-600',
      stats: [
        {
          name: 'battles',
          count: appStats.battleCount,
          icon: Vote,
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
      active: FeatureRetro,
      gradient: 'from-green-500 to-lime-500',
      iconBg: 'bg-gradient-to-r from-green-500 to-lime-500',
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
      active: FeatureStoryboard,
      gradient: 'from-emerald-500 to-emerald-600',
      iconBg: 'bg-gradient-to-r from-emerald-500 to-emerald-600',
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
      title: $LL.projects(),
      active: FeatureProject,
      gradient: 'from-sky-500 to-sky-600',
      iconBg: 'bg-gradient-to-r from-sky-500 to-sky-600',
      stats: [
        {
          name: 'projects',
          count: appStats.projectCount,
          icon: Package,
          active: FeatureProject,
        },
      ],
    },
    {
      title: $LL.estimationScales(),
      active: FeaturePoker,
      gradient: 'from-yellow-500 to-yellow-600',
      iconBg: 'bg-gradient-to-r from-yellow-500 to-yellow-600',
      stats: [
        {
          name: 'estimationScales',
          count: appStats.estimationScaleCount,
          icon: ChartNoAxesColumn,
          active: FeaturePoker,
        },
        {
          name: 'publicEstimationScales',
          count: appStats.publicEstimationScaleCount,
          icon: ChartNoAxesColumn,
          active: FeaturePoker,
        },
        {
          name: 'organizationEstimationScales',
          count: appStats.organizationEstimationScaleCount,
          icon: ChartNoAxesColumn,
          active: FeaturePoker,
        },
        {
          name: 'teamEstimationScales',
          count: appStats.teamEstimationScaleCount,
          icon: ChartNoAxesColumn,
          active: FeaturePoker,
        },
      ],
    },
    {
      title: $LL.retroTemplates(),
      active: FeatureRetro,
      gradient: 'from-pink-500 to-pink-600',
      iconBg: 'bg-gradient-to-r from-pink-500 to-pink-600',
      stats: [
        {
          name: 'retroTemplates',
          count: appStats.retroTemplateCount,
          icon: SquareDashedKanban,
          active: FeatureRetro,
        },
        {
          name: 'publicRetroTemplates',
          count: appStats.publicRetroTemplateCount,
          icon: SquareDashedKanban,
          active: FeatureRetro,
        },
        {
          name: 'organizationRetroTemplates',
          count: appStats.organizationRetroTemplateCount,
          icon: SquareDashedKanban,
          active: FeatureRetro,
        },
        {
          name: 'teamRetroTemplates',
          count: appStats.teamRetroTemplateCount,
          icon: SquareDashedKanban,
          active: FeatureRetro,
        },
      ],
    },
  ]);
</script>

<svelte:head>
  <title>{$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="admin">
  <!-- Loading State -->
  {#if isLoading}
    <div class="flex items-center justify-center py-8">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
    </div>
  {:else}
    <!-- Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      {#each overviewStats.filter(stat => stat.active) as stat}
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-100 dark:border-gray-700 overflow-hidden group hover:shadow-lg transition-all duration-300">
          <div class="p-4">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <p class="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">{stat.title}</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{stat.value.toLocaleString()}</p>
              </div>
              <div class="w-12 h-12 bg-gradient-to-r {stat.color} rounded-lg flex items-center justify-center shadow-md group-hover:scale-110 transition-transform duration-300">
                <stat.icon class="w-6 h-6 text-white" />
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <!-- Detailed Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-6 gap-4 mb-6">
      {#each statGroups.filter(g => g.active) as group}
        {#if group.stats.some(stat => stat.active)}
          <div class="bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 rounded-lg shadow-md overflow-hidden">
            <!-- Group Header -->
            <div class="bg-gradient-to-r {group.gradient} px-4 py-2">
              <h2 class="text-sm font-bold text-white font-rajdhani uppercase tracking-wide">
                {group.title}
              </h2>
            </div>
            
            <!-- Stats Content -->
            <div class="p-3">
              <div class="space-y-2">
                {#each group.stats.filter(stat => stat.active) as stat}
                  <div class="flex items-center p-2 bg-gray-50 dark:bg-gray-700/50 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-all duration-200 group cursor-pointer">
                    <div class="flex-shrink-0 mr-3">
                      <div class="w-8 h-8 {group.iconBg} rounded flex items-center justify-center shadow-sm group-hover:shadow-md transition-shadow duration-200">
                        <stat.icon class="w-4 h-4 text-white" />
                      </div>
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="text-xs font-medium text-gray-900 dark:text-gray-100 truncate">
                        {$LL[stat.name]()}
                      </p>
                      <p class="text-lg font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors duration-200">
                        {stat.count.toLocaleString()}
                      </p>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        {/if}
      {/each}
    </div>

    <!-- Maintenance Section -->
    <div class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-700 rounded-lg p-6 border border-gray-200 dark:border-gray-600">
      <div class="text-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white font-rajdhani uppercase tracking-wide mb-1">
          {$LL.maintenance()}
        </h2>
        <p class="text-sm text-gray-600 dark:text-gray-400">Cleanup operations to maintain system performance</p>
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
        <!-- Clean Guests -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-all duration-300 group">
          <div class="text-center">
            <div class="w-12 h-12 bg-gradient-to-r from-red-500 to-red-600 rounded-lg flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
              <Ghost class="w-6 h-6 text-white" />
            </div>
            <h3 class="font-bold text-gray-700 dark:text-gray-300 mb-2 text-xs uppercase tracking-wide">
              {$LL.cleanGuests({ daysOld: CleanupGuestsDaysOld })}
            </h3>
            <HollowButton onClick={cleanGuests} color="red" class="w-full text-xs py-1">
              {$LL.execute()}
            </HollowButton>
          </div>
        </div>

        <!-- Clean Battles -->
        {#if FeaturePoker}
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-all duration-300 group">
            <div class="text-center">
              <div class="w-12 h-12 bg-gradient-to-r from-red-500 to-red-600 rounded-lg flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
                <Vote class="w-6 h-6 text-white" />
              </div>
              <h3 class="font-bold text-gray-700 dark:text-gray-300 mb-2 text-xs uppercase tracking-wide">
                {$LL.cleanBattles({ daysOld: CleanupBattlesDaysOld })}
              </h3>
              <HollowButton onClick={cleanBattles} color="red" class="w-full text-xs py-1">
                {$LL.execute()}
              </HollowButton>
            </div>
          </div>
        {/if}

        <!-- Clean Retros -->
        {#if FeatureRetro}
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-all duration-300 group">
            <div class="text-center">
              <div class="w-12 h-12 bg-gradient-to-r from-red-500 to-red-600 rounded-lg flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
                <RefreshCcw class="w-6 h-6 text-white" />
              </div>
              <h3 class="font-bold text-gray-700 dark:text-gray-300 mb-2 text-xs uppercase tracking-wide">
                {$LL.adminCleanOldRetros({ daysOld: CleanupRetrosDaysOld })}
              </h3>
              <HollowButton onClick={cleanRetros} color="red" class="w-full text-xs py-1">
                {$LL.execute()}
              </HollowButton>
            </div>
          </div>
        {/if}

        <!-- Clean Storyboards -->
        {#if FeatureStoryboard}
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 hover:shadow-md transition-all duration-300 group">
            <div class="text-center">
              <div class="w-12 h-12 bg-gradient-to-r from-red-500 to-red-600 rounded-lg flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform duration-300">
                <LayoutDashboard class="w-6 h-6 text-white" />
              </div>
              <h3 class="font-bold text-gray-700 dark:text-gray-300 mb-2 text-xs uppercase tracking-wide">
                {$LL.adminCleanOldStoryboards({ daysOld: CleanupStoryboardsDaysOld })}
              </h3>
              <HollowButton onClick={cleanStoryboards} color="red" class="w-full text-xs py-1">
                {$LL.execute()}
              </HollowButton>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</AdminPageLayout>
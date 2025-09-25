<script lang="ts">
  import { AppConfig, appRoutes } from '../config';
  import type { ApiClient } from '../types/apiclient';
  import LL from '../i18n/i18n-svelte';
  import type { Team } from '../types/team';
  import PokerGame from './poker/PokerGame.svelte';
  import type { Retro } from '../types/retro';
  import type { Storyboard } from '../types/storyboard';
  import { user } from '../stores';
  import { onMount } from 'svelte';
  import BoxList from '../components/BoxList.svelte';
  import PageLayout from '../components/PageLayout.svelte';
  import {
    ArrowRight,
    LayoutDashboard,
    RefreshCcw,
    Users,
    Vote,
    ChevronDown,
    Filter,
    Gamepad2,
    RotateCcw,
    CheckCircle,
  } from 'lucide-svelte';
  import { validateUserIsRegistered } from '../validationUtils';

  interface Props {
    xfetch: ApiClient;
    notifications: {
      success: (msg: string) => void;
      danger: (msg: string) => void;
      info: (msg: string) => void;
    };
    router: any;
  }

  let { xfetch, notifications, router }: Props = $props();
  let teams: Team[] = $state([]);
  let games: PokerGame[] = $state([]);
  let gameCount: number = $state(0);
  let retros: Retro[] = $state([]);
  let retroCount: number = $state(0);
  let storyboards: Storyboard[] = $state([]);
  let storyboardCount: number = $state(0);
  let selectedTeam: Team | null = $state(null);
  let showTeamDropdown: boolean = $state(false);

  function loadDashboardData() {
    Promise.any([loadTeams(), loadGames(), loadRetros(), loadStoryboards()]);
  }

  async function loadTeams() {
    if (!validateUserIsRegistered($user)) {
      return Promise.resolve();
    }
    try {
      xfetch(`/api/users/${$user.id}/teams`)
        .then(res => res.json())
        .then(function (result) {
          teams = result.data;
        })
        .catch(function () {
          notifications.danger($LL.getTeamsError());
        });
    } catch (error) {
      console.error('Error loading teams:', error);
    }
  }

  async function loadGames() {
    if (!AppConfig.FeaturePoker) {
      return;
    }
    try {
      const prefix = selectedTeam ? `/api/teams/${selectedTeam.id}/battles` : `/api/users/${$user.id}/battles`;
      xfetch(`${prefix}?limit=4&offset=0`)
        .then(res => res.json())
        .then(function (result) {
          games = result.data;
          gameCount = result.meta.count;
        })
        .catch(function () {
          notifications.danger($LL.myBattlesError());
        });
    } catch (error) {
      console.error('Error loading poker games:', error);
    }
  }

  async function loadRetros() {
    if (!AppConfig.FeatureRetro) {
      return;
    }
    try {
      const prefix = selectedTeam ? `/api/teams/${selectedTeam.id}/retros` : `/api/users/${$user.id}/retros`;
      xfetch(`${prefix}?limit=4&offset=0`)
        .then(res => res.json())
        .then(function (result) {
          retros = result.data;
          retroCount = result.meta.count;
        })
        .catch(function () {
          notifications.danger($LL.getRetrosErrorMessage());
        });
    } catch (error) {
      console.error('Error loading retros:', error);
    }
  }

  async function loadStoryboards() {
    if (!AppConfig.FeatureStoryboard) {
      return;
    }
    try {
      const prefix = selectedTeam ? `/api/teams/${selectedTeam.id}/storyboards` : `/api/users/${$user.id}/storyboards`;
      xfetch(`${prefix}?limit=4&offset=0`)
        .then(res => res.json())
        .then(function (result) {
          storyboards = result.data;
          storyboardCount = result.meta.count;
        })
        .catch(function (error) {
          notifications.danger($LL.getStoryboardsErrorMessage());
        });
    } catch (error) {
      console.error('Error loading storyboards:', error);
    }
  }

  function selectTeam(team: Team | null) {
    selectedTeam = team;
    showTeamDropdown = false;
    // Reload data with new team filter
    loadGames();
    loadRetros();
    loadStoryboards();
  }

  function toggleTeamDropdown() {
    showTeamDropdown = !showTeamDropdown;
  }

  const {} = AppConfig;

  onMount(() => {
    if (!$user) {
      router.route('/login');
      return;
    }
    loadDashboardData();
  });
</script>

<svelte:head>
  <title>{$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <!-- Background with subtle gradient -->
  <div>
    <!-- Hero Section with Team Filter -->
    <div class="mb-8">
      <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-6">
        <!-- Welcome Text -->
        <div class="flex-1">
          <h1 class="text-4xl font-bold tracking-tight text-slate-900 dark:text-white sm:text-5xl lg:text-6xl">
            {$LL.welcomeBack()}
            <span
              class="bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400 bg-clip-text text-transparent"
            >
              {$user.name}
            </span>
          </h1>
          <p class="mt-4 text-xl text-slate-600 dark:text-slate-300">
            {$LL.dashboardSubtitle()}
          </p>
        </div>

        {#if validateUserIsRegistered($user)}
          <!-- Team Filter Dropdown -->
          <div class="flex-shrink-0 lg:mt-2">
            <div class="flex items-center space-x-4">
              <div class="flex items-center space-x-2">
                <Filter class="h-5 w-5 text-slate-500 dark:text-slate-400" />
                <span class="text-sm font-medium text-slate-700 dark:text-slate-300">{$LL.filterByTeam()}</span>
              </div>

              <div class="relative">
                <button
                  onclick={toggleTeamDropdown}
                  class="inline-flex items-center justify-between rounded-xl bg-white dark:bg-slate-800 px-4 py-2.5 text-sm font-medium text-slate-700 dark:text-slate-300 shadow-sm ring-1 ring-slate-300 dark:ring-slate-600 hover:bg-slate-50 dark:hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 min-w-[200px]"
                >
                  <div class="flex items-center space-x-2">
                    <Users class="h-4 w-4" />
                    <span>
                      {selectedTeam ? selectedTeam.name : $LL.allTeams()}
                    </span>
                  </div>
                  <ChevronDown
                    class="h-4 w-4 ml-2 transition-transform duration-200 {showTeamDropdown ? 'rotate-180' : ''}"
                  />
                </button>

                {#if showTeamDropdown}
                  <div
                    class="absolute right-0 mt-2 w-64 origin-top-right rounded-xl bg-white dark:bg-slate-800 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
                  >
                    <div class="py-2">
                      <button
                        onclick={() => selectTeam(null)}
                        class="flex w-full items-center px-4 py-2.5 text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 {!selectedTeam
                          ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300'
                          : ''}"
                      >
                        <div class="flex items-center space-x-3">
                          <div
                            class="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-100 dark:bg-slate-700"
                          >
                            <Users class="h-4 w-4 text-slate-600 dark:text-slate-400" />
                          </div>
                          <div class="text-left">
                            <div class="font-medium">{$LL.allTeams()}</div>
                            <div class="text-xs text-slate-500 dark:text-slate-400">
                              {$LL.showContentForAllTeams()}
                            </div>
                          </div>
                        </div>
                      </button>

                      {#each teams as team}
                        <button
                          onclick={() => selectTeam(team)}
                          class="flex w-full items-center px-4 py-2.5 text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 {selectedTeam?.id ===
                          team.id
                            ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300'
                            : ''}"
                        >
                          <div class="flex items-center space-x-3">
                            <div
                              class="flex h-8 w-8 items-center justify-center rounded-lg bg-emerald-100 dark:bg-emerald-900/30"
                            >
                              <Users class="h-4 w-4 text-emerald-600 dark:text-emerald-400" />
                            </div>
                            <div class="text-left">
                              <div class="font-medium">{team.name}</div>
                            </div>
                          </div>
                        </button>
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>

              {#if selectedTeam}
                <button
                  onclick={() => selectTeam(null)}
                  class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 font-medium"
                >
                  {$LL.clearFilter()}
                </button>
              {/if}
            </div>
            {#if selectedTeam}
              {@const linkPrefix = selectedTeam.department_id
                ? `/organization/${selectedTeam.organization_id}/department/${selectedTeam.department_id}/team/${selectedTeam.id}`
                : selectedTeam.organization_id
                  ? `/organization/${selectedTeam.organization_id}/team/${selectedTeam.id}`
                  : `/team/${selectedTeam.id}`}
              <div class="flex flex-col sm:flex-row gap-3 mt-4 justify-end">
                <a
                  href={linkPrefix}
                  class="inline-flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-white bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 dark:from-blue-500 dark:to-blue-600 dark:hover:from-blue-600 dark:hover:to-blue-700 rounded-lg shadow-sm hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-slate-800 transition-all duration-200"
                >
                  <Users class="h-4 w-4" />
                  {$LL.teamPage()}
                </a>
                <a
                  href="{linkPrefix}/checkin"
                  class="inline-flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-lime-700 dark:text-lime-300 bg-lime-50 hover:bg-lime-100 dark:bg-lime-900/20 dark:hover:bg-lime-900/30 border border-lime-200 dark:border-lime-700 hover:border-lime-300 dark:hover:border-lime-600 rounded-lg shadow-sm hover:shadow-md focus:outline-none focus:ring-2 focus:ring-lime-500 focus:ring-offset-2 dark:focus:ring-offset-slate-800 transition-all duration-200"
                >
                  <CheckCircle class="h-4 w-4" />
                  {$LL.teamCheckins()}
                </a>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>

    <!-- Main Content - Responsive Grid Layout with Equal Heights -->
    <div class="py-12">
      {#if AppConfig.RequireTeams && !selectedTeam}
        <div
          class="mb-8 rounded-2xl bg-yellow-50 dark:bg-yellow-900/20 p-6 ring-1 ring-yellow-200 dark:ring-yellow-700 flex items-center space-x-4"
        >
          <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-yellow-100 dark:bg-yellow-800">
            <Users class="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
          </div>
          <div>
            <p class="text-sm text-yellow-700 dark:text-yellow-300">
              {$LL.dashboardTeamsRequired()}
            </p>
          </div>
        </div>
      {:else}
        <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-8 grid-rows-1">
          <!-- Poker Games Section -->
          {#if AppConfig.FeaturePoker}
            <section class="flex flex-col h-full">
              <div class="flex items-center justify-between mb-6 flex-shrink-0">
                <div class="flex items-center space-x-3">
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-red-500 to-red-600 shadow-lg shadow-red-500/25"
                  >
                    <Vote class="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">
                      {$LL.myBattles()}
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400">
                      {$LL.planningPokerSessions()}
                    </p>
                  </div>
                </div>
              </div>

              {#if gameCount > 0}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex flex-col"
                >
                  <div class="flex-1">
                    <BoxList
                      items={games}
                      itemType="battle"
                      pageRoute={appRoutes.game}
                      joinBtnText={$LL.battleJoin()}
                      showOwner={false}
                      showOwnerName={true}
                      ownerNameField="teamName"
                      showFacilitatorIcon={true}
                      facilitatorsKey="leaders"
                      showCompletedStories={true}
                    />
                  </div>

                  {#if gameCount > 4}
                    <div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700 flex-shrink-0">
                      <a
                        href={appRoutes.games}
                        class="group w-full flex items-center justify-between p-4 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 hover:border-gray-400 dark:hover:border-gray-500 shadow-sm hover:shadow-md transition-all"
                      >
                        <div class="flex items-center space-x-3">
                          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100 dark:bg-gray-700">
                            <Vote class="h-4 w-4 text-gray-600 dark:text-gray-400" />
                          </div>
                          <div>
                            <div class="font-semibold text-gray-900 dark:text-white">
                              {$LL.viewAllPokerSessions()}
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">
                              {$LL.totalPokerSessions({ total: gameCount })}
                            </div>
                          </div>
                        </div>
                        <ArrowRight
                          class="h-5 w-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:translate-x-1 transition-all"
                        />
                      </a>
                    </div>
                  {/if}
                </div>
              {:else}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex items-center justify-center"
                >
                  <div class="text-center">
                    <div
                      class="mx-auto h-24 w-24 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-6"
                    >
                      <Gamepad2 class="h-12 w-12 text-slate-400" />
                    </div>
                    <p class="text-xl text-slate-600 dark:text-slate-400 mb-2">
                      {selectedTeam
                        ? $LL.noSessionsFoundForTeam({
                            teamName: selectedTeam.name,
                          })
                        : $LL.noGamesFound()}
                    </p>
                    <p class="text-sm text-slate-500 dark:text-slate-500">
                      {selectedTeam ? $LL.trySelectingDifferentTeamForPoker() : $LL.startFirstPlanningPokerSession()}
                    </p>
                  </div>
                </div>
              {/if}
            </section>
          {/if}

          <!-- Retros Section -->
          {#if AppConfig.FeatureRetro}
            <section class="flex flex-col h-full">
              <div class="flex items-center justify-between mb-6 flex-shrink-0">
                <div class="flex items-center space-x-3">
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 shadow-lg shadow-purple-500/25"
                  >
                    <RefreshCcw class="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">
                      {$LL.myRetros()}
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400">
                      {$LL.sprintRetrospectives()}
                    </p>
                  </div>
                </div>
              </div>

              {#if retroCount > 0}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex flex-col"
                >
                  <div class="flex-1">
                    <BoxList
                      items={retros}
                      itemType="retro"
                      pageRoute={appRoutes.retro}
                      ownerField="ownerId"
                      showOwnerName={true}
                      ownerNameField="teamName"
                      joinBtnText={$LL.joinRetro()}
                    />
                  </div>

                  {#if retroCount > 4}
                    <div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700 flex-shrink-0">
                      <a
                        href={appRoutes.retros}
                        class="group w-full flex items-center justify-between p-4 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 hover:border-gray-400 dark:hover:border-gray-500 shadow-sm hover:shadow-md transition-all"
                      >
                        <div class="flex items-center space-x-3">
                          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100 dark:bg-gray-700">
                            <RefreshCcw class="h-4 w-4 text-gray-600 dark:text-gray-400" />
                          </div>
                          <div>
                            <div class="font-semibold text-gray-900 dark:text-white">
                              {$LL.viewAllRetros()}
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">
                              {$LL.totalRetros({ total: retroCount })}
                            </div>
                          </div>
                        </div>
                        <ArrowRight
                          class="h-5 w-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:translate-x-1 transition-all"
                        />
                      </a>
                    </div>
                  {/if}
                </div>
              {:else}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex items-center justify-center"
                >
                  <div class="text-center">
                    <div
                      class="mx-auto h-24 w-24 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-6"
                    >
                      <RotateCcw class="h-12 w-12 text-slate-400" />
                    </div>
                    <p class="text-xl text-slate-600 dark:text-slate-400 mb-2">
                      {selectedTeam
                        ? $LL.noRetrosFoundForTeam({
                            teamName: selectedTeam.name,
                          })
                        : $LL.noRetrosFound()}
                    </p>
                    <p class="text-sm text-slate-500 dark:text-slate-500">
                      {selectedTeam ? $LL.trySelectingDifferentTeamForRetros() : $LL.startFirstSprintRetrospective()}
                    </p>
                  </div>
                </div>
              {/if}
            </section>
          {/if}

          <!-- Storyboards Section -->
          {#if AppConfig.FeatureStoryboard}
            <section class="flex flex-col h-full">
              <div class="flex items-center justify-between mb-6 flex-shrink-0">
                <div class="flex items-center space-x-3">
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 shadow-lg shadow-blue-500/25"
                  >
                    <LayoutDashboard class="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">
                      {$LL.myStoryboards()}
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400">
                      {$LL.agileStoryMapping()}
                    </p>
                  </div>
                </div>
              </div>

              {#if storyboardCount > 0}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex flex-col"
                >
                  <div class="flex-1">
                    <BoxList
                      items={storyboards}
                      itemType="storyboard"
                      showOwnerName={true}
                      ownerNameField="teamName"
                      pageRoute={appRoutes.storyboard}
                      joinBtnText={$LL.joinStoryboard()}
                    />
                  </div>

                  {#if storyboardCount > 4}
                    <div class="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700 flex-shrink-0">
                      <a
                        href={appRoutes.storyboards}
                        class="group w-full flex items-center justify-between p-4 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 hover:border-gray-400 dark:hover:border-gray-500 shadow-sm hover:shadow-md transition-all"
                      >
                        <div class="flex items-center space-x-3">
                          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-100 dark:bg-gray-700">
                            <LayoutDashboard class="h-4 w-4 text-gray-600 dark:text-gray-400" />
                          </div>
                          <div>
                            <div class="font-semibold text-gray-900 dark:text-white">
                              {$LL.viewAllStoryboards()}
                            </div>
                            <div class="text-sm text-gray-500 dark:text-gray-400">
                              {$LL.totalStoryboards({
                                total: storyboardCount,
                              })}
                            </div>
                          </div>
                        </div>
                        <ArrowRight
                          class="h-5 w-5 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:translate-x-1 transition-all"
                        />
                      </a>
                    </div>
                  {/if}
                </div>
              {:else}
                <div
                  class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 flex-1 flex items-center justify-center"
                >
                  <div class="text-center">
                    <div
                      class="mx-auto h-24 w-24 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-6"
                    >
                      <LayoutDashboard class="h-12 w-12 text-slate-400" />
                    </div>
                    <p class="text-xl text-slate-600 dark:text-slate-400 mb-2">
                      {selectedTeam
                        ? $LL.noStoryboardsFoundForTeam({
                            teamName: selectedTeam.name,
                          })
                        : $LL.noStoryboardsFound()}
                    </p>
                    <p class="text-sm text-slate-500 dark:text-slate-500">
                      {selectedTeam ? $LL.trySelectingDifferentTeamForStoryboards() : $LL.startFirstStoryboard()}
                    </p>
                  </div>
                </div>
              {/if}
            </section>
          {/if}
        </div>
      {/if}
    </div>
  </div>
</PageLayout>

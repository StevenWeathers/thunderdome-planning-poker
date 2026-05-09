<script lang="ts">
  import { AppConfig, appRoutes } from '../config';
  import type { ApiClient } from '../types/apiclient';
  import type { NotificationService } from '../types/notifications';
  import type { SessionUser } from '../types/user';
  import LL from '../i18n/i18n-svelte';
  import type { Team } from '../types/team';
  import PokerGame from './poker/PokerGame.svelte';
  import type { PokerGame as PokerGameType } from '../types/poker';
  import type { Retro, RetroAction } from '../types/retro';
  import type { Storyboard } from '../types/storyboard';
  import { user } from '../stores';
  import { onMount } from 'svelte';
  import ActionComments from '../components/retro/ActionComments.svelte';
  import ActionItemAssignees from '../components/retro/ActionItemAssignees.svelte';
  import ActionItemEdit from '../components/retro/ActionItemEdit.svelte';
  import ActionsMenu from '../components/global/ActionsMenu.svelte';
  import BoxList from '../components/BoxList.svelte';
  import HollowButton from '../components/global/HollowButton.svelte';
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
    Mail,
    MessageSquareMore,
    Pencil,
    UsersRound,
  } from '@lucide/svelte';
  import { validateUserIsRegistered } from '../validationUtils';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
  }

  let { xfetch, notifications, router }: Props = $props();
  let teams: Team[] = $state([]);
  let games: PokerGameType[] = $state([]);
  let gameCount: number = $state(0);
  let retros: Retro[] = $state([]);
  let retroCount: number = $state(0);
  let retroActions: RetroAction[] = $state([]);
  let retroActionCount: number = $state(0);
  let showRetroActionComments: boolean = $state(false);
  let selectedRetroActionId: string = $state('');
  let showRetroActionEdit: boolean = $state(false);
  let showRetroActionAssignees: boolean = $state(false);
  let selectedRetroAction: RetroAction | null = $state(null);
  let storyboards: Storyboard[] = $state([]);
  let storyboardCount: number = $state(0);
  type InviteKind = 'organization' | 'department' | 'team';
  type BaseInvite = {
    invite_id: string;
    email: string;
    role: string;
    created_date: string;
    expire_date: string;
  };
  type OrganizationInvite = BaseInvite & { organization_id: string; organization_name: string };
  type DepartmentInvite = BaseInvite & { department_id: string; department_name: string };
  type TeamInvite = BaseInvite & { team_id: string; team_name: string };
  type PendingInvitesResponse = {
    organizationInvites: OrganizationInvite[];
    departmentInvites: DepartmentInvite[];
    teamInvites: TeamInvite[];
  };
  type PendingInviteItem = {
    inviteId: string;
    inviteType: InviteKind;
    name: string;
    role: string;
    expireDate: string;
    href: string;
  };

  let selectedTeam: Team | null = $state(null);
  let showTeamDropdown: boolean = $state(false);
  let isVerifiedUser: boolean = $state(false);
  let rejectingInviteIds: string[] = $state([]);
  let pendingInvites: PendingInvitesResponse = $state({
    organizationInvites: [],
    departmentInvites: [],
    teamInvites: [],
  });

  let pendingInviteItems: PendingInviteItem[] = $derived([
    ...pendingInvites.organizationInvites.map(invite => ({
      inviteId: invite.invite_id,
      inviteType: 'organization' as const,
      name: invite.organization_name,
      role: invite.role,
      expireDate: invite.expire_date,
      href: `${appRoutes.invite}/organization/${invite.invite_id}`,
    })),
    ...pendingInvites.departmentInvites.map(invite => ({
      inviteId: invite.invite_id,
      inviteType: 'department' as const,
      name: invite.department_name,
      role: invite.role,
      expireDate: invite.expire_date,
      href: `${appRoutes.invite}/department/${invite.invite_id}`,
    })),
    ...pendingInvites.teamInvites.map(invite => ({
      inviteId: invite.invite_id,
      inviteType: 'team' as const,
      name: invite.team_name,
      role: invite.role,
      expireDate: invite.expire_date,
      href: `${appRoutes.invite}/team/${invite.invite_id}`,
    })),
  ]);

  function loadDashboardData() {
    Promise.any([loadTeams(), loadGames(), loadRetros(), loadRetroActions(), loadStoryboards()]);
    loadPendingInvites();
  }

  $effect(() => {
    if ($user?.id && validateUserIsRegistered($user) && typeof $user.verified === 'boolean') {
      isVerifiedUser = $user.verified;
    }
  });

  async function resolveVerifiedUser() {
    if (!validateUserIsRegistered($user) || !$user.id) {
      isVerifiedUser = false;
      return false;
    }

    try {
      const result = await xfetch(`/api/users/${$user.id}`).then(res => res.json());
      isVerifiedUser = !!result.data?.verified;

      user.update({
        ...$user,
        verified: isVerifiedUser,
      } as SessionUser);

      return isVerifiedUser;
    } catch (error) {
      console.error('Error resolving verified status:', error);
      isVerifiedUser = false;
      return false;
    }
  }

  function getInviteTypeLabel(inviteType: InviteKind) {
    switch (inviteType) {
      case 'organization':
        return $LL.organization();
      case 'department':
        return $LL.department();
      case 'team':
        return $LL.team();
    }
  }

  function getInviteTypeBadgeClasses(inviteType: InviteKind) {
    switch (inviteType) {
      case 'organization':
        return 'bg-amber-100 text-amber-900 ring-1 ring-amber-300 dark:bg-amber-500/15 dark:text-amber-200 dark:ring-amber-400/30';
      case 'department':
        return 'bg-sky-100 text-sky-900 ring-1 ring-sky-300 dark:bg-sky-500/15 dark:text-sky-200 dark:ring-sky-400/30';
      case 'team':
        return 'bg-emerald-100 text-emerald-900 ring-1 ring-emerald-300 dark:bg-emerald-500/15 dark:text-emerald-200 dark:ring-emerald-400/30';
    }
  }

  function formatInviteDate(date: string) {
    return new Date(date).toLocaleString([], {
      year: 'numeric',
      month: 'numeric',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    });
  }

  function getInviteApiPath(invite: PendingInviteItem) {
    return `/api/users/${$user.id}/invite/${invite.inviteType}/${invite.inviteId}`;
  }

  function removePendingInvite(invite: PendingInviteItem) {
    pendingInvites = {
      organizationInvites:
        invite.inviteType === 'organization'
          ? pendingInvites.organizationInvites.filter(item => item.invite_id !== invite.inviteId)
          : pendingInvites.organizationInvites,
      departmentInvites:
        invite.inviteType === 'department'
          ? pendingInvites.departmentInvites.filter(item => item.invite_id !== invite.inviteId)
          : pendingInvites.departmentInvites,
      teamInvites:
        invite.inviteType === 'team'
          ? pendingInvites.teamInvites.filter(item => item.invite_id !== invite.inviteId)
          : pendingInvites.teamInvites,
    };
  }

  async function rejectInvite(invite: PendingInviteItem) {
    if (rejectingInviteIds.includes(invite.inviteId)) {
      return;
    }

    rejectingInviteIds = [...rejectingInviteIds, invite.inviteId];

    try {
      await xfetch(getInviteApiPath(invite), { method: 'DELETE' });
      removePendingInvite(invite);
      notifications.success($LL.inviteRejected());
    } catch (error) {
      console.error('Error rejecting invite:', error);
      notifications.danger($LL.inviteRejectError());
    } finally {
      rejectingInviteIds = rejectingInviteIds.filter(id => id !== invite.inviteId);
    }
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

  async function loadRetroActions() {
    if (!AppConfig.FeatureRetro || !validateUserIsRegistered($user)) {
      retroActions = [];
      retroActionCount = 0;
      return;
    }

    try {
      const params = new URLSearchParams({
        limit: '100',
        offset: '0',
        completed: 'false',
      });

      if (selectedTeam) {
        params.set('teamId', selectedTeam.id);
      }

      xfetch(`/api/users/${$user.id}/retro-actions?${params.toString()}`)
        .then(res => res.json())
        .then(function (result) {
          retroActions = result.data;
          retroActionCount = result.meta.count;
          const currentId = selectedRetroAction?.id;
          selectedRetroAction = currentId ? (retroActions.find(action => action.id === currentId) ?? null) : null;
        })
        .catch(function () {
          notifications.danger('Failed to fetch retro action items.');
        });
    } catch (error) {
      console.error('Error loading retro action items:', error);
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

  async function loadPendingInvites() {
    if (!validateUserIsRegistered($user) || !isVerifiedUser) {
      pendingInvites = {
        organizationInvites: [],
        departmentInvites: [],
        teamInvites: [],
      };
      return Promise.resolve();
    }

    try {
      xfetch(`/api/users/${$user.id}/invites`)
        .then(res => res.json())
        .then(function (result) {
          pendingInvites = {
            organizationInvites: result.data.organizationInvites ?? [],
            departmentInvites: result.data.departmentInvites ?? [],
            teamInvites: result.data.teamInvites ?? [],
          };
        })
        .catch(function () {
          notifications.danger($LL.getUserInvitesError());
        });
    } catch (error) {
      console.error('Error loading pending invites:', error);
    }
  }

  function selectTeam(team: Team | null) {
    selectedTeam = team;
    showTeamDropdown = false;
    // Reload data with new team filter
    loadGames();
    loadRetros();
    loadRetroActions();
    loadStoryboards();
  }

  function toggleTeamDropdown() {
    showTeamDropdown = !showTeamDropdown;
  }

  function toggleRetroActionComments(actionId: string | null) {
    showRetroActionComments = !showRetroActionComments;
    selectedRetroActionId = actionId ?? '';
  }

  function toggleRetroActionEdit(action: RetroAction | null) {
    showRetroActionEdit = !showRetroActionEdit;
    selectedRetroAction = action;
  }

  function toggleRetroActionAssignees(action: RetroAction | null) {
    showRetroActionAssignees = !showRetroActionAssignees;
    selectedRetroAction = action;
  }

  function handleRetroActionEdit(action: RetroAction) {
    xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
      method: 'PUT',
      body: {
        content: action.content,
        completed: action.completed,
      },
    })
      .then(function () {
        loadRetroActions();
        toggleRetroActionEdit(null);
        notifications.success($LL.updateActionItemSuccess());
      })
      .catch(function () {
        notifications.danger($LL.updateActionItemError());
      });
  }

  function handleRetroActionDelete(action: RetroAction) {
    return () => {
      xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
        method: 'DELETE',
      })
        .then(function () {
          loadRetroActions();
          toggleRetroActionEdit(null);
          notifications.success($LL.deleteActionItemSuccess());
        })
        .catch(function () {
          notifications.danger($LL.deleteActionItemError());
        });
    };
  }

  function handleRetroActionAssigneeAdd(retroId: string, actionId: string, userId: string) {
    xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
      method: 'POST',
      body: {
        user_id: userId,
      },
    })
      .then(function () {
        loadRetroActions();
      })
      .catch(function () {
        notifications.danger($LL.updateActionItemError());
      });
  }

  function handleRetroActionAssigneeRemove(retroId: string, actionId: string, userId: string) {
    return () => {
      xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
        method: 'DELETE',
        body: {
          user_id: userId,
        },
      })
        .then(function () {
          loadRetroActions();
        })
        .catch(function () {
          notifications.danger($LL.updateActionItemError());
        });
    };
  }

  function getRetroActionMenuActions(action: RetroAction) {
    return [
      {
        label: $LL.edit(),
        icon: Pencil,
        onclick: () => toggleRetroActionEdit(action),
      },
      {
        label: $LL.assignees(),
        icon: UsersRound,
        onclick: () => toggleRetroActionAssignees(action),
      },
    ];
  }

  const {} = AppConfig;

  onMount(() => {
    if (!$user?.id) {
      router.route('/login');
      return;
    }

    void (async () => {
      await resolveVerifiedUser();
      loadDashboardData();
    })();
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
                ? `${AppConfig.PathPrefix}/organization/${selectedTeam.organization_id}/department/${selectedTeam.department_id}/team/${selectedTeam.id}`
                : selectedTeam.organization_id
                  ? `${AppConfig.PathPrefix}/organization/${selectedTeam.organization_id}/team/${selectedTeam.id}`
                  : `${AppConfig.PathPrefix}/team/${selectedTeam.id}`}
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
        {#if validateUserIsRegistered($user) && isVerifiedUser && !selectedTeam && pendingInviteItems.length > 0}
          <section class="mb-8">
            <div
              class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 shadow-sm"
            >
              <div class="flex items-start justify-between gap-4 mb-6">
                <div class="flex items-center gap-3">
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-orange-500 shadow-lg shadow-amber-500/20"
                  >
                    <Mail class="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">
                      {$LL.invites()} ({pendingInviteItems.length})
                    </h2>
                  </div>
                </div>
              </div>

              <div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
                {#each pendingInviteItems as invite}
                  <div
                    class="group rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900/60 p-4 hover:border-slate-300 dark:hover:border-slate-600 hover:shadow-md transition-all"
                  >
                    <div class="flex items-start justify-between gap-4">
                      <div class="min-w-0">
                        <div class="flex items-center gap-2 mb-2">
                          <span
                            class={`inline-flex items-center rounded-full px-3 py-1 text-[11px] font-bold uppercase tracking-[0.12em] shadow-sm ${getInviteTypeBadgeClasses(
                              invite.inviteType,
                            )}`}
                          >
                            {getInviteTypeLabel(invite.inviteType)}
                          </span>
                        </div>
                        <div class="text-base font-semibold text-slate-900 dark:text-white">
                          {invite.name}
                        </div>
                        <div class="mt-1 text-sm text-slate-500 dark:text-slate-400">
                          {$LL.expireDate()}: {formatInviteDate(invite.expireDate)}
                        </div>
                      </div>

                      <div class="flex items-center justify-end gap-2 shrink-0 self-center">
                        <HollowButton
                          onClick={() => rejectInvite(invite)}
                          disabled={rejectingInviteIds.includes(invite.inviteId)}
                          color="red"
                        >
                          {#if rejectingInviteIds.includes(invite.inviteId)}
                            {$LL.rejectingInvite()}
                          {:else}
                            {$LL.rejectInvite()}
                          {/if}
                        </HollowButton>
                        <HollowButton
                          href={invite.href}
                          color="green"
                          additionalClasses="inline-flex items-center gap-2"
                        >
                          {$LL.acceptInvite()}
                          <ArrowRight class="h-4 w-4 transition-transform group-hover:translate-x-1" />
                        </HollowButton>
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          </section>
        {/if}

        {#if AppConfig.FeatureRetro && validateUserIsRegistered($user)}
          <section class="mb-8">
            <div
              class="rounded-2xl bg-white/70 dark:bg-slate-800/70 backdrop-blur-sm p-6 ring-1 ring-slate-200/50 dark:ring-slate-700/50 shadow-sm"
            >
              <div class="flex items-start justify-between gap-4 mb-6">
                <div class="flex items-center gap-3">
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 shadow-lg shadow-emerald-500/20"
                  >
                    <CheckCircle class="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <h2 class="text-2xl font-bold text-slate-900 dark:text-white">
                      Open action items assigned to you ({retroActionCount})
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400">Outstanding commitments from your team retros</p>
                  </div>
                </div>
              </div>

              {#if retroActionCount > 0}
                <div class="grid grid-cols-1 xl:grid-cols-2 gap-4">
                  {#each retroActions as action}
                    <div
                      class="rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900/60 p-5 shadow-sm transition hover:border-slate-300 hover:shadow-md dark:hover:border-slate-600"
                    >
                      <div class="flex items-start gap-4">
                        <div class="min-w-0 flex-1">
                          <div class="mb-3 flex flex-wrap items-start justify-between gap-3">
                            <div class="flex flex-wrap items-center gap-2">
                              <span
                                class="inline-flex items-center rounded-full bg-emerald-100 px-3 py-1 text-[11px] font-bold uppercase tracking-[0.12em] text-emerald-900 ring-1 ring-emerald-300 dark:bg-emerald-500/15 dark:text-emerald-200 dark:ring-emerald-400/30"
                              >
                                {action.teamName || selectedTeam?.name || 'Team'}
                              </span>
                            </div>
                            <div class="flex items-center gap-2">
                              <button
                                type="button"
                                class="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-slate-50 px-3 py-1.5 text-xs font-medium text-slate-600 transition hover:border-blue-200 hover:bg-blue-50 hover:text-blue-700 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-300 dark:hover:border-sky-700 dark:hover:bg-sky-950/40 dark:hover:text-sky-300"
                                onclick={() => toggleRetroActionComments(action.id)}
                                aria-label={`Open comments for action item: ${action.content}`}
                              >
                                <MessageSquareMore class="h-4 w-4" />
                                <span>{action.comments.length}</span>
                              </button>
                              <ActionsMenu
                                actions={getRetroActionMenuActions(action)}
                                ariaLabel="Action item actions"
                              />
                            </div>
                          </div>
                          <p
                            class="text-base font-semibold leading-6 text-slate-900 dark:text-white whitespace-pre-wrap break-words"
                          >
                            {action.content}
                          </p>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              {:else}
                <div
                  class="rounded-xl border border-dashed border-slate-300 dark:border-slate-600 bg-slate-50/80 dark:bg-slate-900/40 p-8 text-center"
                >
                  <p class="text-lg font-semibold text-slate-900 dark:text-white">
                    {selectedTeam
                      ? `No open action items are assigned to you in ${selectedTeam.name}.`
                      : 'No open action items are assigned to you right now.'}
                  </p>
                </div>
              {/if}
            </div>
          </section>
        {/if}

        {#if showRetroActionComments}
          <ActionComments
            toggleComments={() => toggleRetroActionComments(null)}
            actions={retroActions}
            selectedActionId={selectedRetroActionId}
            getRetrosActions={loadRetroActions}
            {xfetch}
            {notifications}
          />
        {/if}

        {#if showRetroActionEdit}
          <ActionItemEdit
            toggleEdit={() => toggleRetroActionEdit(null)}
            handleEdit={handleRetroActionEdit}
            handleDelete={handleRetroActionDelete}
            action={selectedRetroAction}
          />
        {/if}

        {#if showRetroActionAssignees}
          <ActionItemAssignees
            {xfetch}
            {notifications}
            toggleAssignees={() => toggleRetroActionAssignees(null)}
            assignableUsers={[]}
            action={selectedRetroAction}
            handleAssigneeAdd={handleRetroActionAssigneeAdd}
            handleAssigneeRemove={handleRetroActionAssigneeRemove}
          />
        {/if}

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
                        class="group w-full flex items-center justify-between p-4 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 hover:border-gray-400 dark:hover:border-gray-500 shadow-sm hover:shadow-md transition-all"
                        href={appRoutes.games}
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

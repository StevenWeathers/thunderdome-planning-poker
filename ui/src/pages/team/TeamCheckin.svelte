<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Sockette from 'sockette';

  import PageLayout from '../../components/PageLayout.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import Checkin from '../../components/checkin/Checkin.svelte';
  import TeamCheckinCard from '../../components/checkin/TeamCheckinCard.svelte';
  import { Calendar, ChevronRight } from '@lucide/svelte';
  import Gauge from '../../components/Gauge.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { formatDayForInput, getTimezoneName, subtractDays } from '../../dateUtils';
  import Picker from '../../components/timezone-picker/Picker.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';
  import type { TeamUser, TeamCheckin as BaseTeamCheckin } from '../../types/team';
  import type { UserDisplay } from '../../types/user';
  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  // Local type override - the API returns user as a single object, not an array
  interface TeamCheckin extends Omit<BaseTeamCheckin, 'user'> {
    user: TeamUser;
  }

  interface GaugeDetailItem {
    id: string;
    name: string;
    avatar?: string;
    gravatarHash?: string;
    pictureUrl?: string;
    met: boolean;
  }

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    teamId: string;
    organizationId?: string;
    departmentId?: string;
  }

  // Props using Svelte 5 syntax
  let { xfetch, router, notifications, teamId, organizationId, departmentId }: Props = $props();

  // State variables using $state()
  let timezone = $state(getTimezoneName());
  let showCheckin = $state(false);
  let now = $state(new Date());
  let maxNegativeDate = $state<string>('');
  let selectedDate = $state<string>('');
  let selectedCheckin = $state<TeamCheckin | null>(null);
  let stats = $state({
    participants: 0,
    pPerc: 0,
    goals: 0,
    gPerc: 0,
    blocked: 0,
    bPerc: 0,
  });

  let team = $state({
    id: '',
    name: '',
    subscribed: false,
  });
  let organization = $state({
    id: '',
    name: '',
  });
  let department = $state({
    id: '',
    name: '',
  });

  $effect(() => {
    if (teamId) {
      team.id = teamId;
    }
    if (organizationId) {
      organization.id = organizationId;
    }
    if (departmentId) {
      department.id = departmentId;
    }
  });
  let users = $state<TeamUser[]>([]);
  let userCount = $state(1);

  let organizationRole = $state('');
  let departmentRole = $state('');
  let teamRole = $state('');

  let checkins = $state<TeamCheckin[]>([]);
  let checkinColumns = $state<Array<{ checkins: TeamCheckin[] }>>([]);
  let expandedCheckins = $state<Record<string, boolean>>({});
  let showOnlyDiscussionItems = $state(false);
  let userMap: Map<string, UserDisplay> = $state(new Map());
  let checkinDateInput: HTMLInputElement | null = $state(null);

  let ws: any;

  function updateTimezone(tz: string): void {
    timezone = tz;
    getCheckins();
  }

  function openCheckinDatePicker(): void {
    const input = checkinDateInput;

    if (!input) {
      return;
    }

    const pickerInput = input as HTMLInputElement & { showPicker?: () => void };

    if (pickerInput.showPicker) {
      pickerInput.showPicker();
      return;
    }

    input.focus();
    input.click();
  }

  function formatCheckinDateLabel(dateValue: string): string {
    if (!dateValue) {
      return '';
    }

    const [year, month, day] = dateValue.split('-').map(Number);

    if (!year || !month || !day) {
      return dateValue;
    }

    return new Intl.DateTimeFormat(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    }).format(new Date(year, month - 1, day));
  }

  function divideCheckins(checkins: TeamCheckin[]): void {
    const half = Math.ceil(checkins.length / 2);

    const checkins1 = checkins.slice(0, half);
    const checkins2 = checkins.slice(half);
    checkinColumns = [{ checkins: checkins1 }, { checkins: checkins2 }];
  }

  function filterCheckins() {
    let filteredCheckins = [...checkins];

    if (showOnlyDiscussionItems === true) {
      filteredCheckins = filteredCheckins.filter(c => {
        return c.discuss !== '' || c.blockers !== '';
      });
    }

    divideCheckins(filteredCheckins);
  }

  function toggleDiscussionFilter(): void {
    showOnlyDiscussionItems = !showOnlyDiscussionItems;
    filterCheckins();
  }

  function hasDiscussionItems(checkin: TeamCheckin): boolean {
    return checkin.blockers !== '' || checkin.discuss !== '';
  }

  function syncExpandedCheckins(nextCheckins: TeamCheckin[]): Record<string, boolean> {
    const previousDiscussionState = new Map(checkins.map(checkin => [checkin.id, hasDiscussionItems(checkin)]));

    return nextCheckins.reduce((nextExpandedCheckins: Record<string, boolean>, checkin: TeamCheckin) => {
      const previousExpanded = expandedCheckins[checkin.id];

      if (previousExpanded === undefined) {
        return nextExpandedCheckins;
      }

      const previousHadDiscussion = previousDiscussionState.get(checkin.id);
      const nextHasDiscussion = hasDiscussionItems(checkin);

      if (previousHadDiscussion !== nextHasDiscussion) {
        return nextExpandedCheckins;
      }

      nextExpandedCheckins[checkin.id] = previousExpanded;
      return nextExpandedCheckins;
    }, {});
  }

  function isCheckinExpanded(checkin: TeamCheckin): boolean {
    return expandedCheckins[checkin.id] ?? !hasDiscussionItems(checkin);
  }

  function toggleCheckinSummary(checkin: TeamCheckin): void {
    expandedCheckins = {
      ...expandedCheckins,
      [checkin.id]: !isCheckinExpanded(checkin),
    };
  }

  function getUniqueCheckinsByUser(sourceCheckins: TeamCheckin[]): Map<string, TeamCheckin> {
    return sourceCheckins.reduce((entries: Map<string, TeamCheckin>, checkin: TeamCheckin) => {
      if (!entries.has(checkin.user.id)) {
        entries.set(checkin.user.id, checkin);
      }

      return entries;
    }, new Map<string, TeamCheckin>());
  }

  function buildGaugeDetails(
    teamUsers: TeamUser[],
    uniqueCheckins: Map<string, TeamCheckin>,
    matcher: (checkin: TeamCheckin | undefined) => boolean,
  ): GaugeDetailItem[] {
    return [...teamUsers]
      .map((teamUser: TeamUser) => {
        const checkin = uniqueCheckins.get(teamUser.id);

        return {
          id: teamUser.id,
          name: teamUser.name,
          avatar: teamUser.avatar,
          gravatarHash: teamUser.gravatarHash,
          pictureUrl: teamUser.pictureUrl || '',
          met: matcher(checkin),
        };
      })
      .sort((left: GaugeDetailItem, right: GaugeDetailItem) => {
        if (left.met !== right.met) {
          return left.met ? -1 : 1;
        }

        return left.name.localeCompare(right.name);
      });
  }

  // Derived values using $derived()
  const apiPrefix = '/api';
  const orgPrefix = $derived(
    departmentId
      ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
      : `${apiPrefix}/organizations/${organizationId}`,
  );
  const teamPrefix = $derived(organizationId ? `${orgPrefix}/teams/${teamId}` : `${apiPrefix}/teams/${teamId}`);

  function getTeam() {
    xfetch(teamPrefix)
      .then((res: Response) => res.json())
      .then(function (result: any) {
        team = result.data.team;
        teamRole = result.data.teamRole;

        if (departmentId) {
          department = result.data.department;
          departmentRole = result.data.departmentRole;
        }
        if (organizationId) {
          organization = result.data.organization;
          organizationRole = result.data.organizationRole;
        }
      })
      .catch(function () {
        notifications.danger($LL.teamGetError());
      });
  }

  function getCheckins() {
    xfetch(`${teamPrefix}/checkins?date=${selectedDate}&tz=${timezone}`)
      .then((res: Response) => res.json())
      .then(function (result: any) {
        expandedCheckins = syncExpandedCheckins(result.data);
        checkins = result.data;
        calculateCheckinStats();
        filterCheckins();
      })
      .catch(function () {
        notifications.danger($LL.getCheckinsError());
      });
  }

  function getUsers() {
    xfetch(`${teamPrefix}/users?limit=1000&offset=0`)
      .then((res: Response) => res.json())
      .then(function (result: any) {
        users = result.data;
        userCount = result.meta.count;
        userMap = users.reduce((prev: Map<string, UserDisplay>, cur: TeamUser) => {
          prev.set(cur.id, {
            id: cur.id,
            name: cur.name,
            avatar: cur.avatar,
            gravatarHash: cur.gravatarHash,
            pictureUrl: cur.pictureUrl || '',
          });
          return prev;
        }, new Map<string, UserDisplay>());
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function toggleCheckin(checkin?: TeamCheckin) {
    showCheckin = !showCheckin;
    if (checkin) {
      selectedCheckin = checkin;
    } else {
      selectedCheckin = null;
    }
  }

  function handleCheckin(checkin: any) {
    const body = {
      ...checkin,
      checkinDate: selectedDate,
      timeZone: timezone,
    };

    xfetch(`${teamPrefix}/checkins`, { body })
      .then((res: Response) => res.json())
      .then(function () {
        toggleCheckin();
      })
      .catch(function (error: any) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            if (result.error === 'REQUIRES_TEAM_USER') {
              notifications.danger($LL.teamUserRequiredToCheckin());
            } else {
              notifications.danger($LL.checkinError());
            }
          });
        } else {
          notifications.danger($LL.checkinError());
        }
      });
  }

  function handleCheckinEdit(checkinId: string, checkin: any) {
    xfetch(`${teamPrefix}/checkins/${checkinId}`, {
      body: checkin,
      method: 'PUT',
    })
      .then((res: Response) => res.json())
      .then(function () {
        toggleCheckin();
      })
      .catch(function () {
        notifications.danger($LL.updateCheckinError());
      });
  }

  function handleCheckinDelete(checkinId: string) {
    xfetch(`${teamPrefix}/checkins/${checkinId}`, { method: 'DELETE' })
      .then((res: Response) => res.json())
      .catch(function () {
        notifications.danger($LL.deleteCheckinError());
      });
  }

  function handleCheckinComment(checkinId: string, comment: any) {
    const body = {
      ...comment,
    };

    xfetch(`${teamPrefix}/checkins/${checkinId}/comments`, { body })
      .then((res: Response) => res.json())
      .catch(function (error: any) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            if (result.error === 'REQUIRES_TEAM_USER') {
              notifications.danger($LL.teamUserRequiredToComment());
            } else {
              notifications.danger($LL.checkinCommentError());
            }
          });
        } else {
          notifications.danger($LL.checkinCommentError());
        }
      });
  }

  function handleCheckinCommentEdit(checkinId: string, commentId: string, comment: any) {
    const body = {
      ...comment,
    };

    xfetch(`${teamPrefix}/checkins/${checkinId}/comments/${commentId}`, {
      body,
      method: 'PUT',
    })
      .then((res: Response) => res.json())
      .catch(function (error: any) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            if (result.error === 'REQUIRES_TEAM_USER') {
              notifications.danger($LL.teamUserRequiredToComment());
            } else {
              notifications.danger($LL.checkinCommentError());
            }
          });
        } else {
          notifications.danger($LL.checkinCommentError());
        }
      });
  }

  const handleCommentDelete = (checkinId: string, commentId: string) => {
    xfetch(`${teamPrefix}/checkins/${checkinId}/comments/${commentId}`, {
      method: 'DELETE',
    })
      .then((res: Response) => res.json())
      .catch(function () {
        notifications.danger($LL.checkinCommentDeleteError());
      });
  };

  const onSocketMessage = function (evt: MessageEvent) {
    const parsedEvent = JSON.parse(evt.data);

    switch (parsedEvent.type) {
      case 'init':
      case 'checkin_added':
      case 'checkin_updated':
      case 'checkin_deleted':
      case 'comment_added':
      case 'comment_updated':
      case 'comment_deleted':
        getCheckins();
        break;
      default:
        break;
    }
  };

  const sendSocketEvent = (type: string, value: any) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

  function calculateCheckinStats() {
    const uniqueCheckins = getUniqueCheckinsByUser(checkins);

    stats.blocked = 0;
    stats.goals = 0;

    uniqueCheckins.forEach((checkin: TeamCheckin) => {
      if (checkin.blockers !== '') {
        ++stats.blocked;
      }
      if (checkin.goalsMet) {
        ++stats.goals;
      }
    });
    stats.participants = uniqueCheckins.size;

    stats.pPerc = Math.round((100 * stats.participants) / (userCount || 1));
    stats.gPerc = Math.round((100 * stats.goals) / (stats.participants || 1));
    stats.bPerc = Math.round((100 * stats.blocked) / (stats.participants || 1));

    stats = stats;

    return stats;
  }

  // Derived reactive values
  const isAdmin = $derived(organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole === 'ADMIN');

  const isTeamMember = $derived(organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole !== '');

  const alreadyCheckedIn = $derived(
    checkins && checkins.find((c: TeamCheckin) => c.user.id === $user.id) !== undefined,
  );

  const teamCheckinLocked = $derived(
    AppConfig.SubscriptionsEnabled && selectedDate !== formatDayForInput(now) && !team.subscribed,
  );

  const discussionCheckinCount = $derived(checkins.filter(checkin => hasDiscussionItems(checkin)).length);
  const uniqueCheckinsByUser = $derived(getUniqueCheckinsByUser(checkins));
  const participationGaugeDetails = $derived(
    buildGaugeDetails(users, uniqueCheckinsByUser, checkin => checkin !== undefined),
  );
  const goalsGaugeDetails = $derived(
    buildGaugeDetails(users, uniqueCheckinsByUser, checkin => checkin?.goalsMet === true),
  );
  const blockedGaugeDetails = $derived(
    buildGaugeDetails(users, uniqueCheckinsByUser, checkin => (checkin?.blockers || '') !== ''),
  );
  let participationGaugeProps: any = $derived({
    text: $LL.participation(),
    percentage: stats.pPerc,
    stat: stats.pPerc,
    count: `${stats.participants} / ${userCount}`,
    details: participationGaugeDetails,
  });
  let goalsGaugeProps: any = $derived({
    text: $LL.goalsMet(),
    percentage: stats.gPerc,
    color: 'green',
    stat: stats.gPerc,
    count: `${stats.goals} / ${stats.participants}`,
    details: goalsGaugeDetails,
  });
  let blockedGaugeProps: any = $derived({
    text: $LL.blocked(),
    percentage: stats.bPerc,
    color: 'red',
    stat: stats.bPerc,
    count: `${stats.blocked} / ${stats.participants}`,
    details: blockedGaugeDetails,
  });

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    selectedDate = formatDayForInput(now);
    maxNegativeDate = formatDayForInput(subtractDays(now, 60));

    getTeam();
    getUsers();

    ws = new Sockette(`${getWebsocketAddress()}/api/teams/${teamId}/checkin`, {
      timeout: 2e3,
      maxAttempts: 15,
      onmessage: onSocketMessage,
      onclose: e => {
        if (e.code === 4005) {
          ws.close();
        } else if (e.code === 4004) {
          router.route(appRoutes.teams);
        } else if (e.code === 4001) {
          user.delete();
          router.route(appRoutes.login);
        }
      },
    });
  });

  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
</script>

<svelte:head>
  <title>{$LL.team()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
    <div class="min-w-0 md:grow">
      <div class="text-2xl font-semibold font-rajdhani leading-none dark:text-white sm:text-3xl">
        <span class="uppercase">{$LL.team()}</span>
        <ChevronRight class="inline-block h-7 w-7 align-middle sm:h-8 sm:w-8" />
        <a
          class="align-middle text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
          href={organizationId
            ? departmentId
              ? `${appRoutes.organization}/${organization.id}/department/${department.id}/team/${team.id}`
              : `${appRoutes.organization}/${organization.id}/team/${team.id}`
            : `${appRoutes.team}/${team.id}`}
        >
          {team.name}
        </a>
        <ChevronRight class="inline-block h-7 w-7 align-middle sm:h-8 sm:w-8" />
        <h1 class="inline-block text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white">
          {$LL.checkIn()}
        </h1>
      </div>

      {#if organizationId}
        <div class="mt-3 text-lg font-semibold font-rajdhani dark:text-white sm:text-xl">
          <span class="uppercase">{$LL.organization()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.organization}/{organization.id}">{organization.name}</a
          >
          {#if departmentId}
            &nbsp;
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.department()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}">{department.name}</a
            >
          {/if}
        </div>
      {/if}
    </div>

    <div class="flex flex-col gap-3 md:items-end md:ps-2">
      <div class="flex flex-wrap items-center gap-3 md:justify-end">
        <div class="md:text-right">
          <Picker timezone={timezone || ''} onUpdate={(tz: string | null) => tz && updateTimezone(tz)} />
        </div>
        <div class="relative me-2 md:me-4">
          <button
            type="button"
            class="flex items-center gap-2 bg-transparent text-2xl font-rajdhani font-semibold leading-none uppercase dark:text-white sm:gap-2.5 sm:text-3xl"
            aria-label="Open check-in date picker"
            onclick={openCheckinDatePicker}
          >
            <span>{formatCheckinDateLabel(selectedDate)}</span>
            <Calendar class="h-5 w-5 text-emerald-500 dark:text-emerald-400 sm:h-6 sm:w-6" />
          </button>
          <input
            bind:this={checkinDateInput}
            type="date"
            id="checkindate"
            bind:value={selectedDate}
            min={maxNegativeDate}
            onchange={getCheckins}
            class="pointer-events-none absolute inset-0 h-full w-full cursor-pointer opacity-0"
            tabindex="-1"
            aria-hidden="true"
          />
        </div>
        <SolidButton
          additionalClasses="font-rajdhani uppercase text-2xl"
          onClick={toggleCheckin}
          testid="check-in"
          disabled={alreadyCheckedIn || teamCheckinLocked}
          >{$LL.checkIn()}
        </SolidButton>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-2 lg:grid-cols-4 gap-8 my-4">
    <div class="px-2 md:px-4">
      <Gauge {...participationGaugeProps} />
    </div>
    <div class="px-2 md:px-4">
      <Gauge {...goalsGaugeProps} />
    </div>
    <div class="px-2 md:px-4">
      <Gauge {...blockedGaugeProps} />
    </div>
  </div>

  <div
    class="mt-8 mb-5 flex justify-end rounded-2xl border border-slate-200 bg-white px-3 py-2 shadow-sm dark:border-gray-700 dark:bg-gray-800/90 sm:px-4"
  >
    <button
      type="button"
      class={`group inline-flex w-full items-center justify-between gap-3 rounded-xl border px-3 py-2.5 text-left transition-colors duration-200 sm:w-auto sm:min-w-[21rem] ${
        showOnlyDiscussionItems
          ? 'border-emerald-200 bg-emerald-50/80 text-emerald-900 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-100'
          : 'border-slate-200 bg-slate-50 text-slate-700 hover:border-slate-300 dark:border-gray-600 dark:bg-gray-800 dark:text-slate-200 dark:hover:border-gray-500'
      }`}
      role="switch"
      aria-checked={showOnlyDiscussionItems}
      aria-label={$LL.showBlockedCheckins()}
      onclick={toggleDiscussionFilter}
    >
      <div class="min-w-0 flex-1">
        <div class="truncate text-sm font-semibold">
          {$LL.showBlockedCheckins()}
        </div>
      </div>

      <span
        class={`shrink-0 rounded-full px-2.5 py-1 text-xs font-semibold leading-none ring-1 ${
          showOnlyDiscussionItems
            ? 'bg-emerald-100 text-emerald-800 ring-emerald-200 dark:bg-emerald-500/20 dark:text-emerald-200 dark:ring-emerald-500/30'
            : 'bg-slate-200 text-slate-700 ring-slate-300 dark:bg-gray-700 dark:text-slate-200 dark:ring-gray-600'
        }`}
      >
        {discussionCheckinCount}
      </span>

      <span
        class={`relative inline-flex h-7 w-12 shrink-0 items-center rounded-full p-1 transition-colors duration-200 ${
          showOnlyDiscussionItems
            ? 'bg-emerald-500 dark:bg-emerald-400'
            : 'bg-slate-300 group-hover:bg-slate-400 dark:bg-gray-600 dark:group-hover:bg-gray-500'
        }`}
        aria-hidden="true"
      >
        <span
          class={`h-5 w-5 rounded-full bg-white shadow-sm transition duration-200 ${showOnlyDiscussionItems ? 'translate-x-5' : 'translate-x-0'}`}
        ></span>
      </span>
    </button>
  </div>

  <div class="w-full grid gap-6 md:grid-cols-2 xl:grid-cols-2">
    {#each checkinColumns as col}
      <div class="space-y-6">
        {#each col.checkins as checkin, i}
          <TeamCheckinCard
            {checkin}
            currentUserId={$user.id}
            {isAdmin}
            isExpanded={isCheckinExpanded(checkin)}
            {userMap}
            onToggleSummary={toggleCheckinSummary}
            onEdit={toggleCheckin}
            onDelete={handleCheckinDelete}
            onCreateComment={handleCheckinComment}
            onEditComment={handleCheckinCommentEdit}
            onDeleteComment={handleCommentDelete}
          />
        {/each}
      </div>
    {/each}
  </div>

  {#if showCheckin}
    {#if selectedCheckin}
      <Checkin
        userId={$user.id}
        checkinId={selectedCheckin.id}
        yesterday={selectedCheckin.yesterday}
        today={selectedCheckin.today}
        blockers={selectedCheckin.blockers}
        discuss={selectedCheckin.discuss}
        goalsMet={selectedCheckin.goalsMet}
        {toggleCheckin}
        {handleCheckin}
        {handleCheckinEdit}
        {xfetch}
        {notifications}
        {teamPrefix}
        {selectedDate}
        timeZone={timezone}
      />
    {:else}
      <Checkin
        userId={$user.id}
        {toggleCheckin}
        {handleCheckin}
        {handleCheckinEdit}
        {teamPrefix}
        {xfetch}
        {notifications}
        {selectedDate}
        timeZone={timezone}
      />
    {/if}
  {/if}
</PageLayout>

<style>
  #checkindate::-webkit-calendar-picker-indicator {
    opacity: 0;
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }
</style>

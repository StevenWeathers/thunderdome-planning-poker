<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Sockette from 'sockette';

  import PageLayout from '../../components/PageLayout.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import Checkin from '../../components/checkin/Checkin.svelte';
  import { ChevronRight, Pencil, Trash2 } from 'lucide-svelte';
  import Comments from '../../components/checkin/Comments.svelte';
  import Gauge from '../../components/Gauge.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { formatDayForInput, getTimezoneName, subtractDays } from '../../dateUtils';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import BlockedPing from '../../components/checkin/BlockedPing.svelte';
  import Picker from '../../components/timezone-picker/Picker.svelte';
  import Toggle from '../../components/forms/Toggle.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';
  import type { TeamUser } from '../../types/team';

  // Props using Svelte 5 syntax
  let { xfetch, router, notifications, organizationId, departmentId, teamId } = $props();

  // State variables using $state()
  let timezone = $state(getTimezoneName());
  let showCheckin = $state(false);
  let now = $state(new Date());
  let maxNegativeDate = $state();
  let selectedDate = $state();
  let selectedCheckin = $state();
  let stats = $state({
    participants: 0,
    pPerc: 0,
    goals: 0,
    gPerc: 0,
    blocked: 0,
    bPerc: 0,
  });

  let team = $state({
    id: teamId,
    name: '',
  });
  let organization = $state({
    id: organizationId,
    name: '',
  });
  let department = $state({
    id: departmentId,
    name: '',
  });
  let users = $state([]);
  let userCount = $state(1);

  let organizationRole = $state('');
  let departmentRole = $state('');
  let teamRole = $state('');

  let checkins = $state([]);
  let checkinColumns = $state([]);
  let showOnlyDiscussionItems = $state(false);
  let userMap: Map<string, TeamUser> = $state(new Map());

  function updateTimezone(tz: string): void {
    timezone = tz;
    getCheckins();
  }

  function divideCheckins(checkins): void {
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
      .then(res => res.json())
      .then(function (result) {
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
      .then(res => res.json())
      .then(function (result) {
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
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
        userCount = result.meta.count;
        userMap = users.reduce((prev, cur) => {
          prev[cur.id] = cur;
          return prev;
        }, {});
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function toggleCheckin(checkin) {
    showCheckin = !showCheckin;
    if (checkin) {
      selectedCheckin = checkin;
    } else {
      selectedCheckin = null;
    }
  }

  function handleCheckin(checkin) {
    const body = {
      ...checkin,
    };

    xfetch(`${teamPrefix}/checkins`, { body })
      .then(res => res.json())
      .then(function () {
        toggleCheckin();
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
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

  function handleCheckinEdit(checkinId, checkin) {
    xfetch(`${teamPrefix}/checkins/${checkinId}`, {
      body: checkin,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function () {
        toggleCheckin();
      })
      .catch(function () {
        notifications.danger($LL.updateCheckinError());
      });
  }

  function handleCheckinDelete(checkinId) {
    xfetch(`${teamPrefix}/checkins/${checkinId}`, { method: 'DELETE' })
      .then(res => res.json())
      .catch(function () {
        notifications.danger($LL.deleteCheckinError());
      });
  }

  function handleCheckinComment(checkinId, comment) {
    const body = {
      ...comment,
    };

    xfetch(`${teamPrefix}/checkins/${checkinId}/comments`, { body })
      .then(res => res.json())
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
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

  function handleCheckinCommentEdit(checkinId, commentId, comment) {
    const body = {
      ...comment,
    };

    xfetch(`${teamPrefix}/checkins/${checkinId}/comments/${commentId}`, {
      body,
      method: 'PUT',
    })
      .then(res => res.json())
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
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

  const handleCommentDelete = (checkinId, commentId) => () => {
    xfetch(`${teamPrefix}/checkins/${checkinId}/comments/${commentId}`, {
      method: 'DELETE',
    })
      .then(res => res.json())
      .catch(function () {
        notifications.danger($LL.checkinCommentDeleteError());
      });
  };

  const onSocketMessage = function (evt) {
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

  const ws = new Sockette(`${getWebsocketAddress()}/api/teams/${teamId}/checkin`, {
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

  const sendSocketEvent = (type, value) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

  function calculateCheckinStats() {
    const ucs = [];
    stats.blocked = 0;
    stats.goals = 0;

    checkins.map(c => {
      // @todo - remove once multiple same day checkins are prevented
      if (!ucs.includes(c.user.id)) {
        ucs.push(c.user.id);

        if (c.blockers !== '') {
          ++stats.blocked;
        }
        if (c.goalsMet) {
          ++stats.goals;
        }
      }
    });
    stats.participants = ucs.length;

    stats.pPerc = Math.round((100 * stats.participants) / (userCount || 1));
    stats.gPerc = Math.round((100 * stats.goals) / (stats.participants || 1));
    stats.bPerc = Math.round((100 * stats.blocked) / (stats.participants || 1));

    stats = stats;

    return stats;
  }

  // Derived reactive values
  const isAdmin = $derived(organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole === 'ADMIN');

  const isTeamMember = $derived(organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole !== '');

  const alreadyCheckedIn = $derived(checkins && checkins.find(c => c.user.id === $user.id) !== undefined);

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    selectedDate = formatDayForInput(now);
    maxNegativeDate = formatDayForInput(subtractDays(now, 60));

    getTeam();
    getUsers();
  });

  onDestroy(() => {
    ws.close();
  });
</script>

<svelte:head>
  <title>{$LL.team()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex sm:flex-wrap">
    <div class="md:grow">
      <div class="mb-8">
        <div class="inline-block align-top text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white">
          <h1>
            {$LL.checkIn()}
          </h1>
        </div>
        <div class="inline-block align-top text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white">
          <ChevronRight class="w-8 h-8 inline-block" />
        </div>
        <div class="inline-block">
          <input
            type="date"
            id="checkindate"
            bind:value={selectedDate}
            min={maxNegativeDate}
            max={formatDayForInput(now)}
            onchange={getCheckins}
            class="bg-transparent text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white cursor-pointer"
          />
          <div>
            <Picker {timezone} onUpdate={updateTimezone} />
          </div>
        </div>
      </div>

      {#if organizationId}
        <div class="text-xl font-semibold font-rajdhani dark:text-white">
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
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.team()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}"
            >
              {team.name}
            </a>
          {:else}
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.team()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/team/{team.id}"
            >
              {team.name}
            </a>
          {/if}
        </div>
      {:else}
        <div class="text-2xl font-semibold font-rajdhani dark:text-white">
          <span class="uppercase">{$LL.team()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.team}/{team.id}"
          >
            {team.name}
          </a>
        </div>
      {/if}
    </div>
    <div class="md:ps-2 md:shrink text-right">
      <SolidButton
        additionalClasses="font-rajdhani uppercase text-2xl"
        onClick={toggleCheckin}
        testid="check-in"
        disabled={selectedDate !== formatDayForInput(now) || alreadyCheckedIn}
        >{$LL.checkIn()}
      </SolidButton>
    </div>
  </div>

  <div class="grid grid-cols-2 lg:grid-cols-4 gap-8 my-4">
    <div class="px-2 md:px-4">
      <Gauge
        text={$LL.participation()}
        percentage={stats.pPerc}
        stat={stats.pPerc}
        count="{stats.participants} / {userCount}"
      />
    </div>
    <div class="px-2 md:px-4">
      <Gauge
        text={$LL.goalsMet()}
        percentage={stats.gPerc}
        color="green"
        stat={stats.gPerc}
        count="{stats.goals} / {stats.participants}"
      />
    </div>
    <div class="px-2 md:px-4">
      <Gauge
        text={$LL.blocked()}
        percentage={stats.bPerc}
        color="red"
        stat={stats.bPerc}
        count="{stats.blocked} / {stats.participants}"
      />
    </div>
  </div>

  <div class="mt-8 mb-4 w-full text-right bg-white dark:bg-gray-800 p-3 shadow-lg rounded-lg">
    <Toggle
      name="showOnlyDiscussionItems"
      id="showOnlyDiscussionItems"
      bind:checked={showOnlyDiscussionItems}
      changeHandler={filterCheckins}
      label={$LL.showBlockedCheckins()}
    />
  </div>

  <div class="w-full grid gap-6 md:grid-cols-2 xl:grid-cols-2">
    {#each checkinColumns as col}
      <div class="space-y-6">
        {#each col.checkins as checkin, i}
          <article
            class="group relative overflow-hidden rounded-2xl border border-gray-200/60 dark:border-gray-700/60 bg-white dark:bg-gray-800 shadow-sm hover:shadow-lg dark:shadow-gray-900/10"
            data-testid="checkin"
            aria-labelledby="checkin-user-{checkin.user.id}"
          >
            <!-- Content wrapper -->
            <div class="relative p-6 sm:p-8">
              <!-- Header with avatar and user info -->
              <header class="flex items-start gap-4 sm:gap-6 mb-6">
                <!-- Avatar section -->
                <div class="flex-shrink-0">
                  <div class="relative">
                    <!-- Main avatar container -->
                    <div class="relative w-16 h-16 sm:w-20 sm:h-20">
                      <div
                        class="w-full h-full rounded-full overflow-hidden ring-3 ring-white dark:ring-gray-700 shadow-lg"
                      >
                        <UserAvatar
                          width="80"
                          warriorId={checkin.user.id}
                          avatar={checkin.user.avatar}
                          gravatarHash={checkin.user.gravatarHash}
                          userName={checkin.user.name}
                          options={{
                            class:
                              'w-full h-full object-cover rounded-full transition-transform duration-300 group-hover:scale-110',
                          }}
                        />
                      </div>

                      <!-- Status indicators -->
                      {#if checkin.goalsMet}
                        <div
                          class="absolute -bottom-1 -end-1 rtl:-start-1 rtl:end-auto w-6 h-6 sm:w-7 sm:h-7 rounded-full bg-white dark:bg-gray-800 p-1 shadow-lg ring-2 ring-white dark:ring-gray-700"
                        >
                          <div class="w-full h-full rounded-full bg-emerald-500 flex items-center justify-center">
                            <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                              <path
                                fillRule="evenodd"
                                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                                clipRule="evenodd"
                              />
                            </svg>
                          </div>
                        </div>
                      {/if}

                      {#if checkin.blockers !== ''}
                        <div class="absolute -top-1 -end-1 rtl:-start-1 rtl:end-auto">
                          <BlockedPing />
                        </div>
                      {/if}
                    </div>
                  </div>
                </div>

                <!-- User info and actions -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-start justify-between gap-3">
                    <div class="flex-1 min-w-0">
                      <h3
                        id="checkin-user-{checkin.user.id}"
                        class="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white truncate mb-1"
                        data-testid="checkin-username"
                      >
                        {checkin.user.name}
                      </h3>
                    </div>

                    <!-- Action buttons -->
                    {#if checkin.user.id === $user.id || isAdmin}
                      <div class="flex items-center gap-2" role="group" aria-label="Check-in actions">
                        <button
                          onclick={() => toggleCheckin(checkin)}
                          class="group/btn relative p-2.5 rounded-xl bg-blue-50 hover:bg-blue-100 dark:bg-blue-900/20 dark:hover:bg-blue-900/40 text-blue-600 dark:text-blue-400 transition-all duration-200 hover:scale-110 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
                          title={$LL.edit()}
                          data-testid="checkin-edit"
                          aria-label={`Edit check-in for ${checkin.user.name}`}
                        >
                          <Pencil class="w-4 h-4 transition-transform group-hover/btn:rotate-12" />
                          <span
                            class="absolute inset-0 rounded-xl bg-blue-600/10 scale-0 group-hover/btn:scale-100 transition-transform duration-200"
                          ></span>
                        </button>

                        <button
                          onclick={() => handleCheckinDelete(checkin.id)}
                          class="group/btn relative p-2.5 rounded-xl bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:hover:bg-red-900/40 text-red-600 dark:text-red-400 transition-all duration-200 hover:scale-110 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 dark:focus:ring-offset-gray-800"
                          title="Delete"
                          data-testid="checkin-delete"
                          aria-label={`Delete check-in for ${checkin.user.name}`}
                        >
                          <Trash2 class="w-4 h-4 transition-transform group-hover/btn:scale-110" />
                          <span
                            class="absolute inset-0 rounded-xl bg-red-600/10 scale-0 group-hover/btn:scale-100 transition-transform duration-200"
                          ></span>
                        </button>
                      </div>
                    {/if}
                  </div>
                </div>
              </header>

              <!-- Check-in content -->
              <div class="space-y-6">
                <!-- Yesterday section -->
                <section class="space-y-2">
                  <h4
                    class="flex items-center gap-2 text-sm font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wide"
                  >
                    <div class="w-2 h-2 rounded-full bg-blue-500"></div>
                    {$LL.yesterday()}
                  </h4>
                  <div
                    class="text-gray-800 dark:text-gray-200 leading-relaxed unreset whitespace-pre-wrap pl-4 border-s-2 border-blue-100 dark:border-blue-900/50"
                    data-testid="checkin-yesterday"
                  >
                    {@html checkin.yesterday}
                  </div>
                </section>

                <!-- Today section -->
                <section class="space-y-2">
                  <h4
                    class="flex items-center gap-2 text-sm font-semibold text-gray-600 dark:text-gray-300 uppercase tracking-wide"
                  >
                    <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
                    {$LL.today()}
                  </h4>
                  <div
                    class="text-gray-800 dark:text-gray-200 leading-relaxed unreset whitespace-pre-wrap pl-4 border-s-2 border-emerald-100 dark:border-emerald-900/50"
                    data-testid="checkin-today"
                  >
                    {@html checkin.today}
                  </div>
                </section>

                <!-- Blockers section -->
                {#if checkin.blockers !== ''}
                  <section class="space-y-2">
                    <h4
                      class="flex items-center gap-2 text-sm font-semibold text-red-600 dark:text-red-400 uppercase tracking-wide"
                    >
                      <div class="w-2 h-2 rounded-full bg-red-500 animate-pulse"></div>
                      <span class="flex items-center gap-1">
                        {$LL.blockers()}
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path
                            fillRule="evenodd"
                            d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </span>
                    </h4>
                    <div
                      class="text-gray-800 dark:text-gray-200 leading-relaxed unreset whitespace-pre-wrap pl-4 border-s-2 border-red-200 dark:border-red-900/50 bg-red-50/50 dark:bg-red-900/10 rounded-e-lg py-3"
                      data-testid="checkin-blockers"
                    >
                      {@html checkin.blockers}
                    </div>
                  </section>
                {/if}

                <!-- Discuss section -->
                {#if checkin.discuss !== ''}
                  <section class="space-y-2">
                    <h4
                      class="flex items-center gap-2 text-sm font-semibold text-amber-600 dark:text-amber-400 uppercase tracking-wide"
                    >
                      <div class="w-2 h-2 rounded-full bg-amber-500"></div>
                      <span class="flex items-center gap-1">
                        {$LL.discuss()}
                        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                          <path
                            fillRule="evenodd"
                            d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </span>
                    </h4>
                    <div
                      class="text-gray-800 dark:text-gray-200 leading-relaxed unreset whitespace-pre-wrap pl-4 border-s-2 border-amber-200 dark:border-amber-900/50 bg-amber-50/50 dark:bg-amber-900/10 rounded-e-lg py-3"
                      data-testid="checkin-discuss"
                    >
                      {@html checkin.discuss}
                    </div>
                  </section>
                {/if}

                <!-- Comments section -->
                <section class="pt-4">
                  <div
                    class="bg-gray-50 dark:bg-gray-700/30 rounded-xl p-4 border border-gray-200/50 dark:border-gray-600/30"
                  >
                    <Comments
                      {checkin}
                      {userMap}
                      {isAdmin}
                      handleCreate={handleCheckinComment}
                      handleEdit={handleCheckinCommentEdit}
                      handleDelete={handleCommentDelete}
                    />
                  </div>
                </section>
              </div>
            </div>
          </article>
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
      />
    {/if}
  {/if}
</PageLayout>

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
  import {
    formatDayForInput,
    getTimezoneName,
    subtractDays,
  } from '../../dateUtils';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import BlockedPing from '../../components/checkin/BlockedPing.svelte';
  import Picker from '../../components/timezone-picker/Picker.svelte';
  import Toggle from '../../components/forms/Toggle.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
    organizationId: any;
    departmentId: any;
    teamId: any;
  }

  let {
    xfetch,
    router,
    notifications,
    organizationId,
    departmentId,
    teamId
  }: Props = $props();

  const { AllowRegistration, AllowGuests } = AppConfig;

  let timezone = $state(getTimezoneName());
  let showCheckin = $state(false);
  let now = new Date();
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
  let users = [];
  let userCount = $state(1);

  let organizationRole = $state('');
  let departmentRole = $state('');
  let teamRole = $state('');

  let checkins = $state([]);
  let checkinColumns = $state([]);
  let showOnlyDiscussionItems = $state(false);

  function updateTimezone(ev) {
    timezone = ev.detail.timezone;
    getCheckins();
  }

  function divideCheckins(checkins) {
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

  const apiPrefix = '/api';
  let orgPrefix = $derived(departmentId
    ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
    : `${apiPrefix}/organizations/${organizationId}`);
  let teamPrefix = $derived(organizationId
    ? `${orgPrefix}/teams/${teamId}`
    : `${apiPrefix}/teams/${teamId}`);

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
        filterCheckins();
      })
      .catch(function () {
        notifications.danger($LL.getCheckinsError());
      });
  }

  let userMap = $state({});

  function getUsers() {
    xfetch(`${teamPrefix}/users?limit=1000&offset=0`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
        userCount = result.meta.count;
        userMap = users.reduce((prev, cur) => {
          prev[cur.id] = cur.name;
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
      .then(function () {})
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

  const ws = new Sockette(
    `${getWebsocketAddress()}/api/teams/${teamId}/checkin`,
    {
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
    },
  );

  const sendSocketEvent = (type, value) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

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

  let isAdmin =
    $derived(organizationRole === 'ADMIN' ||
    departmentRole === 'ADMIN' ||
    teamRole === 'ADMIN');
  let isTeamMember =
    $derived(organizationRole === 'ADMIN' ||
    departmentRole === 'ADMIN' ||
    teamRole !== '');

  let checkStats = $derived(checkins && userCount && calculateCheckinStats());
  let alreadyCheckedIn =
    $derived(checkins && checkins.find(c => c.user.id === $user.id) !== undefined);
</script>

<svelte:head>
  <title>{$LL.team()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex sm:flex-wrap">
    <div class="md:grow">
      <div class="mb-8">
        <div
          class="inline-block align-top text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white"
        >
          <h1>
            {$LL.checkIn()}
          </h1>
        </div>
        <div
          class="inline-block align-top text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white"
        >
          <ChevronRight class="w-8 h-8 inline-block" />
        </div>
        <div class="inline-block">
          <input
            type="date"
            id="checkindate"
            bind:value="{selectedDate}"
            min="{maxNegativeDate}"
            max="{formatDayForInput(now)}"
            onchange={getCheckins}
            class="bg-transparent text-3xl font-rajdhani font-semibold leading-none uppercase dark:text-white"
          />
          <Picker
            timezone="{timezone}"
            on:update="{updateTimezone}"
            btnClasses="text-lg text-blue-500 dark:text-sky-400"
          />
        </div>
      </div>

      {#if organizationId}
        <div class="text-xl font-semibold font-rajdhani dark:text-white">
          <span class="uppercase">{$LL.organization()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.organization}/{organization.id}"
            >{organization.name}</a
          >
          {#if departmentId}
            &nbsp;
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.department()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}"
              >{department.name}</a
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
        onClick="{toggleCheckin}"
        testid="check-in"
        disabled="{selectedDate !== formatDayForInput(now) || alreadyCheckedIn}"
        >{$LL.checkIn()}
      </SolidButton>
    </div>
  </div>

  <div class="grid grid-cols-2 lg:grid-cols-4 gap-8 my-4">
    <div class="px-2 md:px-4">
      <Gauge
        text="{$LL.participation()}"
        percentage="{stats.pPerc}"
        stat="{stats.pPerc}"
        count="{stats.participants} / {userCount}"
      />
    </div>
    <div class="px-2 md:px-4">
      <Gauge
        text="{$LL.goalsMet()}"
        percentage="{stats.gPerc}"
        color="green"
        stat="{stats.gPerc}"
        count="{stats.goals} / {stats.participants}"
      />
    </div>
    <div class="px-2 md:px-4">
      <Gauge
        text="{$LL.blocked()}"
        percentage="{stats.bPerc}"
        color="red"
        stat="{stats.bPerc}"
        count="{stats.blocked} / {stats.participants}"
      />
    </div>
  </div>

  <div
    class="mt-8 mb-4 w-full text-right bg-white dark:bg-gray-800 p-3 shadow-lg rounded-lg"
  >
    <Toggle
      name="showOnlyDiscussionItems"
      id="showOnlyDiscussionItems"
      bind:checked="{showOnlyDiscussionItems}"
      changeHandler="{filterCheckins}"
      label="{$LL.showBlockedCheckins()}"
    />
  </div>

  <div class="relative">
    <div class="w-full md:grid md:grid-cols-2 md:gap-4">
      {#each checkinColumns as col}
        <div>
          {#each col.checkins as checkin, i}
            <div
              class="mb-4 w-full flex dark:text-gray-300 bg-white dark:bg-gray-800 p-6 shadow-lg rounded-xl border-gray-300 dark:border-gray-700 border-b"
              data-testid="checkin"
            >
              <div class="shrink me-4 lg:me-6 text-center">
                <div class="flex justify-items-center mb-4">
                  <div class="relative w-20 h-20">
                    <div class="relative w-full h-full">
                      <div
                        class="w-full h-full bg-gray-200 rounded-full shadow"
                      >
                        <UserAvatar
                          width="80"
                          warriorId="{checkin.user.id}"
                          avatar="{checkin.user.avatar}"
                          gravatarHash="{checkin.user.gravatarHash}"
                          userName="{checkin.user.name}"
                          options="{{
                            class: 'w-full h-full rounded-full',
                          }}"
                        />
                      </div>
                      {#if checkin.goalsMet}
                        <div
                          class="absolute bottom-0 w-1/4 h-1/4 rounded-full shadow-md"
                        >
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            viewBox="0 0 1024 1024"
                          >
                            <circle
                              class="text-white fill-current"
                              cx="512"
                              cy="512"
                              r="512"></circle>
                            <circle
                              fill="green"
                              class="fill-current text-green-500"
                              cx="512"
                              cy="512"
                              r="384"></circle>
                            <path
                              class="text-white fill-current"
                              d="M456.4 576.1L334.8 454.4l-81.1 81.1 121.6 121.7 81.1 81.1 81.1-81.1 243.3-243.3-81.1-81.1z"
                            ></path>
                          </svg>
                        </div>
                      {/if}
                      {#if checkin.blockers !== ''}
                        <BlockedPing />
                      {/if}
                      <div
                        class="hidden absolute top-0 end-0 w-1/4 h-1/4 bg-white rounded-full shadow-md"
                      >
                        <!-- emoji -->
                      </div>
                    </div>
                  </div>
                </div>
                <div class="w-20">
                  <div
                    class="font-bold text-blue-500 dark:text-sky-400 mb-4"
                    data-testid="checkin-username"
                  >
                    {checkin.user.name}
                  </div>
                  {#if checkin.user.id === $user.id || isAdmin}
                    <div>
                      <button
                        onclick={() => {
                          toggleCheckin(checkin);
                        }}
                        class="text-blue-500"
                        title="{$LL.edit()}"
                        data-testid="checkin-edit"
                      >
                        <span class="sr-only">{$LL.edit()}</span>
                        <Pencil />
                      </button>
                      <button
                        onclick={() => {
                          handleCheckinDelete(checkin.id);
                        }}
                        class="text-red-500"
                        title="Delete"
                        data-testid="checkin-delete"
                      >
                        <span class="sr-only">{$LL.delete()}</span>
                        <Trash2 />
                      </button>
                    </div>
                  {/if}
                </div>
              </div>
              <div class="grow">
                <div>
                  <div class="font-bold text-gray-400">
                    {$LL.yesterday()}:
                  </div>
                  <div
                    class="unreset whitespace-pre-wrap"
                    data-testid="checkin-yesterday"
                  >
                    {@html checkin.yesterday}
                  </div>
                </div>
                <div>
                  <div class="font-bold text-gray-400">
                    {$LL.today()}:
                  </div>
                  <div
                    class="unreset whitespace-pre-wrap"
                    data-testid="checkin-today"
                  >
                    {@html checkin.today}
                  </div>
                </div>
                {#if checkin.blockers !== ''}
                  <div>
                    <div class="font-bold text-lg text-red-500">
                      {$LL.blockers()}:
                    </div>
                    <div
                      class="unreset whitespace-pre-wrap"
                      data-testid="checkin-blockers"
                    >
                      {@html checkin.blockers}
                    </div>
                  </div>
                {/if}
                {#if checkin.discuss !== ''}
                  <div>
                    <div class="font-bold text-lg text-green-500">
                      {$LL.discuss()}:
                    </div>
                    <div
                      class="unreset whitespace-pre-wrap"
                      data-testid="checkin-discuss"
                    >
                      {@html checkin.discuss}
                    </div>
                  </div>
                {/if}
                <div class="bg-gray-200 dark:bg-gray-600 rounded py-2 px-4">
                  <Comments
                    checkin="{checkin}"
                    userMap="{userMap}"
                    isAdmin="{isAdmin}"
                    handleCreate="{handleCheckinComment}"
                    handleEdit="{handleCheckinCommentEdit}"
                    handleDelete="{handleCommentDelete}"
                  />
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/each}
    </div>
  </div>

  {#if showCheckin}
    {#if selectedCheckin}
      <Checkin
        teamId="{team.id}"
        userId="{$user.id}"
        checkinId="{selectedCheckin.id}"
        yesterday="{selectedCheckin.yesterday}"
        today="{selectedCheckin.today}"
        blockers="{selectedCheckin.blockers}"
        discuss="{selectedCheckin.discuss}"
        goalsMet="{selectedCheckin.goalsMet}"
        toggleCheckin="{toggleCheckin}"
        handleCheckin="{handleCheckin}"
        handleCheckinEdit="{handleCheckinEdit}"
        xfetch="{xfetch}"
        notifications="{notifications}"
        teamPrefix="{teamPrefix}"
      />
    {:else}
      <Checkin
        teamId="{team.id}"
        userId="{$user.id}"
        toggleCheckin="{toggleCheckin}"
        handleCheckin="{handleCheckin}"
        handleCheckinEdit="{handleCheckinEdit}"
      />
    {/if}
  {/if}
</PageLayout>

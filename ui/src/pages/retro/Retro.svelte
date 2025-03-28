<script lang="ts">
  import Sockette from 'sockette';
  import { onDestroy, onMount } from 'svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import {
    Check,
    ChevronRight,
    ExternalLink,
    Pencil,
    SquareCheckBig,
  } from 'lucide-svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import EditRetro from '../../components/retro/EditRetro.svelte';
  import EditActionItem from '../../components/retro/EditActionItem.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import BecomeFacilitator from '../../components/BecomeFacilitator.svelte';
  import LL from '../../i18n/i18n-svelte';
  import Export from '../../components/retro/Export.svelte';
  import GroupPhase from '../../components/retro/GroupPhase.svelte';
  import VotePhase from '../../components/retro/VotePhase.svelte';
  import GroupedItems from '../../components/retro/GroupedItems.svelte';
  import UserCard from '../../components/retro/UserCard.svelte';
  import InviteUser from '../../components/retro/InviteUser.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import PhaseTimer from '../../components/retro/PhaseTimer.svelte';
  import BrainstormPhase from '../../components/retro/BrainstormPhase.svelte';
  import JoinCodeForm from '../../components/global/JoinCodeForm.svelte';
  import FullpageLoader from '../../components/global/FullpageLoader.svelte';
  import RetroActionItemReview from '../../components/retro/RetroActionItemReview.svelte';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';

  export let retroId;
  export let notifications;
  export let router;
  export let xfetch;

  const { AllowRegistration, AllowGuests } = AppConfig;
  const loginOrRegister = AllowGuests ? appRoutes.register : appRoutes.login;

  const hostname = window.location.origin;

  let isLoading = true;
  let socketError = false;
  let socketReconnecting = false;
  let retro = {
    name: '',
    ownerId: '',
    teamId: '',
    phase: 'intro',
    phase_time_limit_min: 0,
    phase_time_start: new Date(),
    phase_auto_advance: false,
    users: [],
    items: [],
    groups: [],
    actionItems: [],
    votes: [],
    facilitators: [],
    maxVotes: 3,
    brainstormVisibility: 'visible',
    facilitatorCode: '',
    joinCode: '',
    readyUsers: [],
    template: {
      id: '',
      name: '',
      description: '',
      format: {
        columns: [],
      },
    },
    allowCumulativeVoting: false,
  };
  let showDeleteRetro = false;
  let actionItem = '';
  let showExport = false;
  let groupedItems = [];
  let JoinPassRequired = false;
  let voteLimitReached = false;
  let allUsersVoted = false;
  let showEditRetro = false;
  let phaseTimeStart = new Date();
  let phaseTimeLimitMin = 0;
  let team = null;
  let columnColors = {};

  function getAssociatedTeam() {
    if (retro.teamId) {
      xfetch(`/api/teams/${retro.teamId}`)
        .then(r => r.json())
        .then(res => {
          team = res.data.team;
        })
        .catch(err => {
          console.error(err);
        });
    }
  }

  function organizeItemsByGroup() {
    const groupMap = retro.groups.reduce((prev, g) => {
      prev[g.id] = {
        id: g.id,
        name: g.name,
        items: [],
        votes: [],
        voteCount: 0,
        userVoted: false,
      };
      return prev;
    }, {});
    let userVoteCount = 0;
    let result = [];
    let voteCount = 0;
    const playerCount = retro.users.filter(u => u.active).length;

    retro.items.map(item => {
      groupMap[item.groupId].items.push(item);
    });

    retro.votes.map(vote => {
      voteCount = voteCount + vote.count;
      groupMap[vote.groupId].voteCount =
        groupMap[vote.groupId].voteCount + vote.count;
      groupMap[vote.groupId].votes.push(vote);
      if (vote.userId === $user.id) {
        userVoteCount = userVoteCount + vote.count;
        groupMap[vote.groupId].userVoted = true;
      }
    });

    voteLimitReached = userVoteCount === retro.maxVotes;
    allUsersVoted = voteCount === playerCount * retro.maxVotes;
    phaseReadyCheck();

    result = Object.values(groupMap);
    if (retro.phase === 'action' || retro.phase === 'completed') {
      result.sort((a, b) => {
        return b.voteCount - a.voteCount;
      });
    }
    if (retro.phase === 'vote') {
      result.sort((a, b) => {
        return b.items.length - a.items.length;
      });
    }

    return result;
  }

  const onSocketMessage = function (evt) {
    isLoading = false;
    const parsedEvent = JSON.parse(evt.data);

    switch (parsedEvent.type) {
      case 'join_code_required':
        JoinPassRequired = true;
        break;
      case 'join_code_incorrect':
        notifications.danger($LL.incorrectPassCode());
        break;
      case 'init':
        JoinPassRequired = false;
        retro = JSON.parse(parsedEvent.value);
        columnColors = retro.template.format.columns.reduce((p, c) => {
          p[c.name] = c.color;
          return p;
        }, {});
        if (retro.phase != 'brainstorm') {
          groupedItems = organizeItemsByGroup();
        }
        phaseTimeStart = new Date(retro.phase_time_start);
        phaseTimeLimitMin = retro.phase_time_limit_min;
        getAssociatedTeam();
        break;
      case 'user_joined': {
        retro.users = JSON.parse(parsedEvent.value) || [];
        const joinedUser = retro.users.find(u => u.id === parsedEvent.userId);
        notifications.success(`${joinedUser.name} joined.`);
        break;
      }
      case 'user_left': {
        const leftUser = retro.users.find(w => w.id === parsedEvent.userId);
        retro.users = JSON.parse(parsedEvent.value);

        notifications.danger(`${leftUser.name} left.`);
        break;
      }
      case 'phase_updated':
        let r = JSON.parse(parsedEvent.value);
        retro.items = r.items;
        retro.groups = r.groups;
        retro.votes = r.votes;
        retro.actionItems = r.actionItems;
        retro.phase = r.phase;
        retro.phase_time_start = new Date(r.phase_time_start);
        retro.phase_time_limit_min = r.phase_time_limit_min;
        retro.readyUsers = [];
        phaseTimeStart = new Date(r.phase_time_start);

        groupedItems = organizeItemsByGroup();

        if (retro.phase !== 'completed') {
          showExport = false;
        }
        break;
      case 'items_updated': {
        const parsedValue = JSON.parse(parsedEvent.value);
        retro.items = parsedValue;
        if (retro.phase !== 'brainstorm') {
          groupedItems = organizeItemsByGroup();
        }
        break;
      }
      case 'user_marked_ready': {
        const readyUser = retro.users.find(w => w.id === parsedEvent.userId);
        retro.readyUsers = JSON.parse(parsedEvent.value);
        phaseReadyCheck();
        notifications.success(`${readyUser.name} is done brainstorming.`);
        break;
      }
      case 'user_marked_unready': {
        const unreadyUser = retro.users.find(w => w.id === parsedEvent.userId);
        retro.readyUsers = JSON.parse(parsedEvent.value);

        notifications.warning(
          `${unreadyUser.name} is no longer done brainstorming.`,
        );
        break;
      }
      case 'item_moved': {
        const parsedValue = JSON.parse(parsedEvent.value);
        const updatedItems = [...retro.items];
        const idx = updatedItems.findIndex(item => {
          return item.id === parsedValue.id;
        });
        updatedItems[idx].groupId = parsedValue.groupId;
        retro.items = updatedItems;
        groupedItems = organizeItemsByGroup();
        break;
      }
      case 'group_name_updated': {
        const parsedValue = JSON.parse(parsedEvent.value);
        const updatedGroups = [...retro.groups];
        const idx = updatedGroups.findIndex(group => {
          return group.id === parsedValue.id;
        });
        updatedGroups[idx].name = parsedValue.name;
        retro.groups = updatedGroups;
        groupedItems = organizeItemsByGroup();
        break;
      }
      case 'groups_updated': {
        const parsedValue = JSON.parse(parsedEvent.value);
        retro.groups = parsedValue;
        groupedItems = organizeItemsByGroup();
        break;
      }
      case 'votes_updated': {
        const parsedValue = JSON.parse(parsedEvent.value);
        retro.votes = parsedValue;
        if (retro.phase === 'vote') {
          groupedItems = organizeItemsByGroup();
        }
        break;
      }
      case 'action_updated':
        retro.actionItems = JSON.parse(parsedEvent.value);
        selectedAction =
          selectedAction !== null
            ? retro.actionItems.find(a => a.id === selectedAction.id)
            : null;
        break;
      case 'facilitators_updated':
        retro.facilitators = JSON.parse(parsedEvent.value);
        break;
      case 'retro_edited':
        const revisedRetro = JSON.parse(parsedEvent.value);
        retro.name = revisedRetro.retroName;
        retro.joinCode = revisedRetro.joinCode;
        retro.brainstormVisibility = revisedRetro.brainstormVisibility;
        retro.maxVotes = revisedRetro.maxVotes;
        retro.phase_auto_advance = revisedRetro.phase_auto_advance;
        break;
      case 'conceded':
        // retro over, goodbye.
        notifications.warning($LL.retroDeleted());
        router.route(appRoutes.retros);
        break;
      default:
        break;
    }
  };

  const ws = new Sockette(`${getWebsocketAddress()}/api/retro/${retroId}`, {
    timeout: 2e3,
    maxAttempts: 15,
    onmessage: onSocketMessage,
    onerror: () => {
      socketError = true;
    },
    onclose: e => {
      if (e.code === 4004) {
        router.route(appRoutes.retros);
      } else if (e.code === 4001) {
        user.delete();
        router.route(`${loginOrRegister}/retro/${retroId}`);
      } else if (e.code === 4003) {
        notifications.danger($LL.duplicateRetroSession());
        router.route(`${appRoutes.retros}`);
      } else if (e.code === 4002) {
        router.route(appRoutes.retros);
      } else {
        socketReconnecting = true;
      }
    },
    onopen: () => {
      isLoading = false;
      socketError = false;
      socketReconnecting = false;
    },
    onmaximum: () => {
      socketReconnecting = false;
    },
  });

  onDestroy(() => {
    ws.close();
  });

  $: isFacilitator =
    retro.facilitators && retro.facilitators.includes($user.id);

  const sendSocketEvent = (type, value) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

  function concedeRetro() {
    sendSocketEvent('concede_retro', '');
  }

  function abandonRetro() {
    sendSocketEvent('abandon_retro', '');
  }

  const toggleDeleteRetro = () => {
    showDeleteRetro = !showDeleteRetro;
  };

  const toggleExport = () => {
    showExport = !showExport;
  };

  let showActionEdit = false;
  let selectedAction = null;
  const toggleActionEdit = id => () => {
    showActionEdit = !showActionEdit;
    selectedAction = retro.actionItems.find(r => r.id === id);
  };

  const handleUserReady = userId => () => {
    sendSocketEvent(`user_ready`, userId);
  };
  const handleUserUnReady = userId => () => {
    sendSocketEvent(`user_unready`, userId);
  };

  const handleItemGroupChange = (itemId, groupId) => {
    sendSocketEvent(
      `group_item`,
      JSON.stringify({
        itemId,
        groupId,
      }),
    );
  };

  const handleGroupNameChange = (groupId, name) => {
    sendSocketEvent(
      `group_name_change`,
      JSON.stringify({
        groupId,
        name,
      }),
    );
  };

  const handleActionItem = evt => {
    evt.preventDefault();

    sendSocketEvent(
      'create_action',
      JSON.stringify({
        content: actionItem,
      }),
    );
    actionItem = '';
  };

  const handleActionUpdate = (id, completed, content) => () => {
    sendSocketEvent(
      'update_action',
      JSON.stringify({
        id,
        completed: !completed,
        content,
      }),
    );
  };

  const handleActionEdit = ({ id, content, completed }) => {
    handleActionUpdate(id, !completed, content)();
    toggleActionEdit(null)();
  };

  const handleAssigneeAdd = (retroId, actionId, userId) => {
    sendSocketEvent(
      'action_assignee_add',
      JSON.stringify({
        id: actionId,
        user_id: userId,
      }),
    );
  };
  const handleAssigneeRemove = (retroId, actionId, userId) => () => {
    sendSocketEvent(
      'action_assignee_remove',
      JSON.stringify({
        id: actionId,
        user_id: userId,
      }),
    );
  };

  const handleActionDelete =
    ({ id }) =>
    () => {
      sendSocketEvent(
        'delete_action',
        JSON.stringify({
          id,
        }),
      );
      toggleActionEdit(null)();
    };

  const setPhase = phase => () => {
    if (!isFacilitator) {
      return;
    }
    const nextPhase = {
      intro: 'brainstorm',
      brainstorm: 'group',
      group: 'vote',
      vote: 'action',
      action: 'completed',
    };

    sendSocketEvent(
      'advance_phase',
      JSON.stringify({
        phase: phase || nextPhase[retro.phase],
      }),
    );
  };

  const handleVote = groupId => {
    sendSocketEvent(
      `group_vote`,
      JSON.stringify({
        groupId,
      }),
    );
  };

  const handleVoteSubtract = groupId => {
    sendSocketEvent(
      `group_vote_subtract`,
      JSON.stringify({
        groupId,
      }),
    );
  };

  const handleAddFacilitator = userId => () => {
    sendSocketEvent(
      'add_facilitator',
      JSON.stringify({
        userId,
      }),
    );
  };

  const handleRemoveFacilitator = userId => () => {
    if (retro.facilitators.length === 1) {
      notifications.danger($LL.removeOnlyFacilitatorError());
      return;
    }

    sendSocketEvent(
      'remove_facilitator',
      JSON.stringify({
        userId,
      }),
    );
  };

  function authRetro(joinPasscode) {
    sendSocketEvent('auth_retro', joinPasscode);
  }

  function handleRetroEdit(revisedRetro) {
    sendSocketEvent('edit_retro', JSON.stringify(revisedRetro));
    toggleEditRetro();
  }

  function toggleEditRetro() {
    showEditRetro = !showEditRetro;
  }

  let showBecomeFacilitator = false;

  function becomeFacilitator(facilitatorCode) {
    sendSocketEvent('self_facilitator', facilitatorCode);
    toggleBecomeFacilitator();
  }

  function phaseTimeRanOut() {
    if (isFacilitator) {
      sendSocketEvent(
        'phase_time_ran_out',
        JSON.stringify({
          phase: 'group',
        }),
      );
    }
  }

  function phaseReadyCheck() {
    const activeUsers = retro.users.filter(u => u.active);
    let allReady = retro.readyUsers.length === activeUsers.length;

    if (isFacilitator && retro.phase_auto_advance) {
      if (retro.phase === 'brainstorm' && allReady) {
        sendSocketEvent(
          'phase_all_ready',
          JSON.stringify({
            phase: 'group',
          }),
        );
      }
      if (retro.phase === 'vote' && allUsersVoted) {
        sendSocketEvent(
          'phase_all_ready',
          JSON.stringify({
            phase: 'action',
          }),
        );
      }
    }
  }

  function toggleBecomeFacilitator() {
    showBecomeFacilitator = !showBecomeFacilitator;
  }

  let showOpenActionItems = false;

  function toggleReviewActionItems() {
    showOpenActionItems = !showOpenActionItems;
  }

  onMount(() => {
    if (!$user.id) {
      router.route(`${loginOrRegister}/retro/${retroId}`);
    }
  });
</script>

<style>
  :global(input:checked ~ div) {
    @apply border-green-500;
  }

  :global(input:checked ~ div svg) {
    @apply block;
  }
</style>

<svelte:head>
  <title>{$LL.retro()} {retro.name} | {$LL.appName()}</title>
</svelte:head>

<div class="flex flex-col flex-grow w-full">
  <div
    class="flex-none px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b border-t border-gray-400 dark:border-gray-700 flex
        flex-wrap"
  >
    <div class="w-1/4">
      <h1 class="text-3xl font-bold leading-tight dark:text-gray-200">
        {retro.name}
      </h1>
    </div>
    <div class="w-3/4 text-right">
      <div>
        {#if retro.phase === 'completed'}
          <SolidButton
            color="green"
            onClick="{toggleExport}"
            testid="retro-export"
          >
            {#if showExport}
              {$LL.back()}
            {:else}
              {$LL.export()}
            {/if}
          </SolidButton>
        {/if}
        {#if retro.phase === 'brainstorm' && phaseTimeLimitMin > 0}
          <PhaseTimer
            retroId="{retro.id}"
            timeLimitMin="{phaseTimeLimitMin}"
            timeStart="{phaseTimeStart}"
            on:ended="{phaseTimeRanOut}"
          />
        {/if}
        {#if isFacilitator}
          {#if retro.phase !== 'completed'}
            <SolidButton
              color="green"
              onClick="{setPhase(null)}"
              testid="retro-nextphase"
            >
              {$LL.nextPhase()}
            </SolidButton>
          {/if}

          <HollowButton
            color="blue"
            onClick="{toggleEditRetro}"
            testid="retro-edit"
          >
            {$LL.editRetro()}
          </HollowButton>

          <HollowButton
            color="red"
            onClick="{toggleDeleteRetro}"
            class="me-2"
            testid="retro-delete"
          >
            {$LL.deleteRetro()}
          </HollowButton>
        {:else}
          <HollowButton
            color="blue"
            onClick="{toggleBecomeFacilitator}"
            testid="become-facilitator"
          >
            {$LL.becomeFacilitator()}
          </HollowButton>
          <HollowButton
            color="red"
            onClick="{abandonRetro}"
            testid="retro-leave"
          >
            {$LL.leaveRetro()}
          </HollowButton>
        {/if}
      </div>
    </div>
  </div>
  <div
    class="flex-none px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-400 dark:border-gray-700 flex flex-wrap"
  >
    <div class="w-1/2">
      <div class="flex items-center text-gray-500 dark:text-gray-300">
        <div
          class="flex-initial px-1 {retro.phase === 'intro' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('intro')}">{$LL.primeDirective()}</button>
        </div>
        <div class="flex-initial px-1">
          <ChevronRight class="inline-block" />
        </div>
        <div
          class="flex-initial px-1 {retro.phase === 'brainstorm' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('brainstorm')}">{$LL.brainstorm()}</button
          >
        </div>
        <div class="flex-initial px-1">
          <ChevronRight class="inline-block" />
        </div>
        <div
          class="flex-initial px-1 {retro.phase === 'group' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('group')}">{$LL.group()}</button>
        </div>
        <div class="flex-initial px-1">
          <ChevronRight class="inline-block" />
        </div>
        <div
          class="flex-initial px-1 {retro.phase === 'vote' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('vote')}">{$LL.vote()}</button>
        </div>
        <div class="flex-initial px-1">
          <ChevronRight class="inline-block" />
        </div>
        <div
          class="flex-initial px-1 {retro.phase === 'action' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('action')}">{$LL.actionItems()}</button>
        </div>
        <div class="flex-initial px-1">
          <ChevronRight class="inline-block" />
        </div>
        <div
          class="flex-initial px-1 {retro.phase === 'completed' &&
            'border-b-2 border-blue-500 dark:border-yellow-400 text-gray-800 dark:text-gray-200'}"
        >
          <button on:click="{setPhase('completed')}">{$LL.done()}</button>
        </div>
      </div>
    </div>
    <div class="w-1/2 text-right text-gray-600 dark:text-gray-400">
      {#if retro.phase === 'brainstorm'}
        {$LL.brainstormPhaseDescription()}
      {:else if retro.phase === 'group'}
        {$LL.groupPhaseDescription()}
      {:else if retro.phase === 'vote'}
        {$LL.votePhaseDescription()}
      {:else if retro.phase === 'action'}
        {$LL.actionPhaseDescription()}
      {/if}
    </div>
  </div>
  {#if showExport}
    <Export retro="{retro}" />
  {/if}
  {#if !showExport}
    <div class="w-full p-4 flex flex-col flex-grow">
      {#if retro.phase === 'intro'}
        {#if showOpenActionItems}
          <RetroActionItemReview
            team="{team}"
            toggle="{toggleReviewActionItems}"
            xfetch="{xfetch}"
            notifications="{notifications}"
          />
          <div class="w-full text-center pt-4 md:pt-6">
            <HollowButton
              color="purple"
              onClick="{toggleReviewActionItems}"
              testid="back-to-prime-directive"
              additionalClasses="py-4 px-6 text-lg"
              >Back to Prime Directive
            </HollowButton>
          </div>
        {:else}
          <div class="m-auto w-full md:w-3/4 lg:w-2/3 dark:text-white">
            {#if team && (!AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && team.subscribed))}
              <div class="text-center pt-10">
                <HollowButton
                  color="purple"
                  onClick="{toggleReviewActionItems}"
                  testid="review-action-items"
                  additionalClasses="py-4 px-6 text-lg"
                  >Review Open Action Items
                </HollowButton>
              </div>
            {:else if team && AppConfig.SubscriptionsEnabled && !team.subscribed}
              <FeatureSubscribeBanner
                salesPitch="Review open action items from previous team retrospectives."
              />
            {/if}
            <h2
              class="md:mt-14 lg:mt-20 text-3xl md:text-4xl lg:text-5xl font-rajdhani mb-2 tracking-wide"
            >
              The Prime Directive
            </h2>
            <div class="title-line bg-yellow-thunder"></div>
            <p
              class="md:leading-loose tracking-wider text-xl md:text-2xl lg:text-3xl"
            >
              "Regardless of what we discover, we understand and truly believe
              that everyone did the best job they could, given what they knew at
              the time, their skills and abilities, the resources available, and
              the situation at hand."
            </p>
            <p class="tracking-wider md:text-lg lg:text-xl">
              &mdash;Norm Kerth, Project Retrospectives: A Handbook for Team
              Review <a
                href="https://retrospectivewiki.org/index.php?title=The_Prime_Directive"
                target="_blank"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              >
                <ExternalLink class="w-6 h-6 md:w-8 md:h-8" />
              </a>
            </p>
          </div>
        {/if}
      {/if}
      <div class="grow flex">
        {#if retro.phase === 'brainstorm'}
          <BrainstormPhase
            items="{retro.items}"
            phase="{retro.phase}"
            isFacilitator="{isFacilitator}"
            sendSocketEvent="{sendSocketEvent}"
            template="{retro.template}"
            users="{retro.users}"
            brainstormVisibility="{retro.brainstormVisibility}"
            columnColors="{columnColors}"
          />
        {/if}
        {#if retro.phase === 'group'}
          <div class="w-full grid grid-cols-2 md:grid-cols-4 gap-2 md:gap-4">
            <GroupPhase
              phase="{retro.phase}"
              groups="{groupedItems}"
              handleItemChange="{handleItemGroupChange}"
              handleGroupNameChange="{handleGroupNameChange}"
              users="{retro.users}"
              sendSocketEvent="{sendSocketEvent}"
              isFacilitator="{isFacilitator}"
              columnColors="{columnColors}"
            />
          </div>
        {/if}
        {#if retro.phase === 'vote'}
          <div class="w-full">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-2 md:gap-4">
              <VotePhase
                phase="{retro.phase}"
                groups="{groupedItems}"
                handleVote="{handleVote}"
                handleVoteSubtract="{handleVoteSubtract}"
                voteLimitReached="{voteLimitReached}"
                allowCumulativeVoting="{retro.allowCumulativeVoting}"
                users="{retro.users}"
                sendSocketEvent="{sendSocketEvent}"
                isFacilitator="{isFacilitator}"
                columnColors="{columnColors}"
              />
            </div>
          </div>
        {/if}
        {#if retro.phase === 'action' || retro.phase === 'completed'}
          <div class="w-full md:w-2/3">
            <div class="grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-4">
              <GroupedItems
                phase="{retro.phase}"
                groups="{groupedItems}"
                users="{retro.users}"
                sendSocketEvent="{sendSocketEvent}"
                isFacilitator="{isFacilitator}"
                columnColors="{columnColors}"
              />
            </div>
          </div>
          <div class="w-full md:w-1/3">
            <div class="ps-4">
              {#if retro.phase === 'action'}
                <div class="flex items-center mb-4">
                  <div class="flex-shrink pe-2">
                    <SquareCheckBig
                      class="w-8 h-8 text-indigo-500 dark:text-violet-400"
                    />
                  </div>
                  <div class="flex-grow">
                    <form on:submit="{handleActionItem}">
                      <input
                        bind:value="{actionItem}"
                        placeholder="{$LL.actionItemPlaceholder()}"
                        class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                    px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                    focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
                        id="actionItem"
                        name="actionItem"
                        type="text"
                        required
                      />
                      <button type="submit" class="hidden"></button>
                    </form>
                  </div>
                </div>
              {/if}
              {#each retro.actionItems as item, i}
                <div
                  class="mb-2 p-2 bg-white dark:bg-gray-800 shadow border-s-4 border-indigo-500 dark:border-violet-400"
                >
                  <div class="flex items-center">
                    <div class="flex-shrink">
                      <button
                        on:click="{toggleActionEdit(item.id)}"
                        class="pe-2 pt-1 text-gray-500 dark:text-gray-400
                                                hover:text-blue-500"
                      >
                        <Pencil />
                      </button>
                    </div>
                    <div class="flex-grow dark:text-white">
                      <div class="pe-2">
                        {#each item.assignees as assignee}
                          <UserAvatar
                            warriorId="{assignee.id}"
                            gravatarHash="{assignee.gravatarHash}"
                            avatar="{assignee.avatar}"
                            userName="{assignee.name}"
                            width="24"
                            class="inline-block me-2"
                          />
                        {/each}
                        {item.content}
                      </div>
                    </div>
                    <div class="flex-shrink">
                      <input
                        type="checkbox"
                        id="{i}Completed"
                        checked="{item.completed}"
                        class="opacity-0 absolute h-6 w-6"
                        on:change="{handleActionUpdate(
                          item.id,
                          item.completed,
                          item.content,
                        )}"
                      />
                      <div
                        class="bg-white dark:bg-gray-800 border-2 rounded-md
                                            border-gray-400 dark:border-gray-300 w-6 h-6 flex flex-shrink-0
                                            justify-center items-center me-2
                                            focus-within:border-blue-500 dark:focus-within:border-sky-500"
                      >
                        <Check
                          class="hidden w-4 h-4 text-green-600 pointer-events-none"
                        />
                      </div>
                      <label for="{i}Completed" class="select-none"></label>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        {/if}
      </div>

      <div class="w-full self-end mb-4 mt-8">
        <div class="flex w-full justify-center">
          {#each retro.users as usr, index (usr.id)}
            {#if usr.active}
              <UserCard
                user="{usr}"
                votes="{retro.votes}"
                maxVotes="{retro.maxVotes}"
                facilitators="{retro.facilitators}"
                readyUsers="{retro.readyUsers}"
                handleAddFacilitator="{handleAddFacilitator}"
                handleRemoveFacilitator="{handleRemoveFacilitator}"
                handleUserReady="{handleUserReady}"
                handleUserUnReady="{handleUserUnReady}"
                phase="{retro.phase}"
              />
            {/if}
          {/each}
        </div>
        {#if retro.phase === 'intro'}
          <div class="mt-4 flex w-full p-2 dark:text-white justify-center">
            <div
              class="w-full md:w-1/2 lg:w-1/3 p-4 bg-white dark:bg-gray-800 shadow-lg rounded-lg"
            >
              <InviteUser hostname="{hostname}" retroId="{retro.id}" />
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

{#if showEditRetro}
  <EditRetro
    retroName="{retro.name}"
    handleRetroEdit="{handleRetroEdit}"
    toggleEditRetro="{toggleEditRetro}"
    joinCode="{retro.joinCode}"
    facilitatorCode="{retro.facilitatorCode}"
    maxVotes="{retro.maxVotes}"
    brainstormVisibility="{retro.brainstormVisibility}"
    phaseAutoAdvance="{retro.phase_auto_advance}"
  />
{/if}

{#if showDeleteRetro}
  <DeleteConfirmation
    toggleDelete="{toggleDeleteRetro}"
    handleDelete="{concedeRetro}"
    confirmText="{$LL.confirmDeleteRetro()}"
    confirmBtnText="{$LL.deleteRetro()}"
  />
{/if}

{#if showActionEdit}
  <EditActionItem
    toggleEdit="{toggleActionEdit(null)}"
    handleEdit="{handleActionEdit}"
    handleDelete="{handleActionDelete}"
    action="{selectedAction}"
    assignableUsers="{retro.users}"
    handleAssigneeAdd="{handleAssigneeAdd}"
    handleAssigneeRemove="{handleAssigneeRemove}"
    retroId="{retro.id}"
  />
{/if}

{#if showBecomeFacilitator}
  <BecomeFacilitator
    handleBecomeFacilitator="{becomeFacilitator}"
    toggleBecomeFacilitator="{toggleBecomeFacilitator}"
  />
{/if}

{#if socketReconnecting}
  <FullpageLoader>
    {$LL.reloadingRetro()}
  </FullpageLoader>
{:else if socketError}
  <FullpageLoader>
    {$LL.retroJoinError()}
  </FullpageLoader>
{:else if isLoading}
  <FullpageLoader>
    {$LL.loadingRetro()}
  </FullpageLoader>
{:else if JoinPassRequired}
  <JoinCodeForm handleSubmit="{authRetro}" submitText="{$LL.joinRetro()}" />
{/if}

<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Sockette from 'sockette';

  import PageLayout from '../../components/PageLayout.svelte';
  import PointCard from '../../components/poker/PointCard.svelte';
  import PokerStories from '../../components/poker/PokerStories.svelte';
  import VoteResults from '../../components/poker/VoteResults.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import ExternalLinkIcon from '../../components/icons/ExternalLinkIcon.svelte';
  import EditPokerGame from '../../components/poker/EditPokerGame.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes, PathPrefix } from '../../config';
  import UserCard from '../../components/poker/UserCard.svelte';
  import VotingControls from '../../components/poker/VotingControls.svelte';
  import InviteUser from '../../components/poker/InviteUser.svelte';
  import VoteTimer from '../../components/poker/VoteTimer.svelte';
  import type { PokerGame, PokerStory } from '../../types/poker';
  import TextInput from '../../components/forms/TextInput.svelte';

  export let battleId: string;
  export let notifications;
  export let eventTag;
  export let router;
  export let xfetch;

  const { AllowRegistration, AllowGuests } = AppConfig;
  const loginOrRegister: string = AllowGuests
    ? appRoutes.register
    : appRoutes.login;

  const hostname: string = window.location.origin;
  const socketExtension: string =
    window.location.protocol === 'https:' ? 'wss' : 'ws';
  const defaultStory: PokerStory = {
    id: '',
    active: false,
    points: '',
    priority: 0,
    skipped: false,
    voteEndTime: undefined,
    voteStartTime: undefined,
    votes: undefined,
    name: '',
    type: '',
    referenceId: '',
    link: '',
    description: '',
    acceptanceCriteria: '',
    position: 0,
  };

  let JoinPassRequired: boolean = false;
  let socketError: boolean = false;
  let socketReconnecting: boolean = false;
  let points: Array<string> = [];
  let vote: string = '';
  let battle: PokerGame = {
    leaders: [],
    autoFinishVoting: false,
    createdDate: undefined,
    hideVoterIdentity: false,
    id: '',
    name: '',
    plans: undefined,
    pointAverageRounding: '',
    pointValuesAllowed: undefined,
    updatedDate: undefined,
    users: undefined,
    votingLocked: false,
    teamId: '',
  };
  let currentStory = { ...defaultStory };
  let showEditBattle: boolean = false;
  let showDeleteBattle: boolean = false;
  let isSpectator: boolean = false;
  let joinPasscode: string = '';
  let voteStartTime: Date = new Date();

  const onSocketMessage = function (evt) {
    const parsedEvent = JSON.parse(evt.data);

    switch (parsedEvent.type) {
      case 'join_code_required':
        JoinPassRequired = true;
        break;
      case 'join_code_incorrect':
        notifications.danger($LL.incorrectPassCode());
        break;
      case 'init': {
        JoinPassRequired = false;
        battle = JSON.parse(parsedEvent.value);
        points = battle.pointValuesAllowed;
        const { spectator = false } =
          battle.users.find(w => w.id === $user.id) || {};
        isSpectator = spectator;

        if (battle.activePlanId !== '') {
          const activePlan = battle.plans.find(
            p => p.id === battle.activePlanId,
          );
          const warriorVote = activePlan.votes.find(
            v => v.warriorId === $user.id,
          ) || {
            vote: '',
          };
          currentStory = activePlan;
          voteStartTime = new Date(activePlan.voteStartTime);
          vote = warriorVote.vote;
        }

        eventTag('join', 'battle', '');
        break;
      }
      case 'warrior_joined': {
        battle.users = JSON.parse(parsedEvent.value);
        const joinedWarrior = battle.users.find(
          w => w.id === parsedEvent.warriorId,
        );
        if (joinedWarrior.id === $user.id) {
          isSpectator = joinedWarrior.spectator;
        }
        if ($user.notificationsEnabled) {
          notifications.success(
            `${$LL.warriorJoined({
              name: joinedWarrior.name,
            })}
    `,
          );
        }
        break;
      }
      case 'warrior_retreated':
        const leftWarrior = battle.users.find(
          w => w.id === parsedEvent.warriorId,
        );
        battle.users = JSON.parse(parsedEvent.value);

        if ($user.notificationsEnabled) {
          notifications.danger(
            `${$LL.warriorRetreated({
              name: leftWarrior.name,
            })}
    `,
          );
        }
        break;
      case 'users_updated':
        battle.users = JSON.parse(parsedEvent.value);
        const updatedWarrior = battle.users.find(w => w.id === $user.id);
        isSpectator = updatedWarrior.spectator;
        break;
      case 'plan_added':
        battle.plans = JSON.parse(parsedEvent.value);
        break;
      case 'story_arranged':
        battle.plans = JSON.parse(parsedEvent.value);
        break;
      case 'plan_activated':
        const updatedPlans = JSON.parse(parsedEvent.value);
        const activePlan = updatedPlans.find(p => p.active);
        currentStory = activePlan;
        voteStartTime = new Date(activePlan.voteStartTime);

        battle.plans = updatedPlans;
        battle.activePlanId = activePlan.id;
        battle.votingLocked = false;
        vote = '';
        break;
      case 'plan_skipped':
        const updatedPlans2 = JSON.parse(parsedEvent.value);
        currentStory = { ...defaultStory };
        battle.plans = updatedPlans2;
        battle.activePlanId = '';
        battle.votingLocked = true;
        vote = '';
        if ($user.notificationsEnabled) {
          notifications.warning($LL.planSkipped());
        }
        break;
      case 'vote_activity':
        const votedWarrior = battle.users.find(
          w => w.id === parsedEvent.warriorId,
        );
        if ($user.notificationsEnabled) {
          notifications.success(
            `${$LL.warriorVoted({
              name: votedWarrior.name,
            })}
    `,
          );
        }

        battle.plans = JSON.parse(parsedEvent.value);
        break;
      case 'vote_retracted':
        const devotedWarrior = battle.users.find(
          w => w.id === parsedEvent.warriorId,
        );
        if ($user.notificationsEnabled) {
          notifications.warning(
            `${$LL.warriorRetractedVote({
              name: devotedWarrior.name,
            })}
    `,
          );
        }

        battle.plans = JSON.parse(parsedEvent.value);
        break;
      case 'voting_ended':
        battle.plans = JSON.parse(parsedEvent.value);
        battle.votingLocked = true;
        break;
      case 'plan_finalized':
        battle.plans = JSON.parse(parsedEvent.value);
        battle.activePlanId = '';
        currentStory = { ...defaultStory };
        vote = '';
        break;
      case 'plan_revised':
        battle.plans = JSON.parse(parsedEvent.value);
        if (battle.activePlanId !== '') {
          const activePlan = battle.plans.find(
            p => p.id === battle.activePlanId,
          );
          currentStory = activePlan;
        }
        break;
      case 'plan_burned':
        const postBurnPlans = JSON.parse(parsedEvent.value);

        if (
          battle.activePlanId !== '' &&
          postBurnPlans.filter(p => p.id === battle.activePlanId).length === 0
        ) {
          battle.activePlanId = '';
          currentStory = { ...defaultStory };
        }

        battle.plans = postBurnPlans;

        break;
      case 'leaders_updated':
        battle.leaders = parsedEvent.value;
        break;
      case 'battle_revised':
        const revisedBattle = JSON.parse(parsedEvent.value);
        battle.name = revisedBattle.battleName;
        points = revisedBattle.pointValuesAllowed;
        battle.autoFinishVoting = revisedBattle.autoFinishVoting;
        battle.pointAverageRounding = revisedBattle.pointAverageRounding;
        battle.joinCode = revisedBattle.joinCode;
        battle.hideVoterIdentity = revisedBattle.hideVoterIdentity;
        battle.teamId = revisedBattle.teamId;
        break;
      case 'battle_conceded':
        // poker over, goodbye.
        notifications.warning($LL.battleDeleted());
        router.route(appRoutes.games);
        break;
      case 'jab_warrior':
        const userToNudge = battle.users.find(w => w.id === parsedEvent.value);
        notifications.info(
          `${$LL.warriorNudgeMessage({
            name: userToNudge.name,
          })}
    `,
        );
        break;
      default:
        break;
    }
  };

  const ws = new Sockette(
    `${socketExtension}://${window.location.host}${PathPrefix}/api/arena/${battleId}`,
    {
      timeout: 2e3,
      maxAttempts: 15,
      onmessage: onSocketMessage,
      onerror: err => {
        socketError = true;
        eventTag('socket_error', 'battle', '');
      },
      onclose: e => {
        if (e.code === 4004) {
          eventTag('not_found', 'battle', '', () => {
            router.route(appRoutes.games);
          });
        } else if (e.code === 4001) {
          eventTag('socket_unauthorized', 'battle', '', () => {
            user.delete();
            router.route(`${loginOrRegister}/battle/${battleId}`);
          });
        } else if (e.code === 4003) {
          eventTag('socket_duplicate', 'battle', '', () => {
            notifications.danger($LL.sessionDuplicate());
            router.route(`${appRoutes.games}`);
          });
        } else if (e.code === 4002) {
          eventTag('battle_warrior_abandoned', 'battle', '', () => {
            router.route(appRoutes.games);
          });
        } else {
          socketReconnecting = true;
          eventTag('socket_close', 'battle', '');
        }
      },
      onopen: () => {
        socketError = false;
        socketReconnecting = false;
        eventTag('socket_open', 'battle', '');
      },
      onmaximum: () => {
        socketReconnecting = false;
        eventTag('socket_error', 'battle', 'Socket Reconnect Max Reached');
      },
    },
  );

  onDestroy(() => {
    eventTag('leave', 'battle', '', () => {
      ws.close();
    });
  });

  const sendSocketEvent = (type, value) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

  const handleVote = event => {
    vote = event.detail.point;
    const voteValue = {
      planId: battle.activePlanId,
      voteValue: vote,
      autoFinishVoting: battle.autoFinishVoting,
    };

    sendSocketEvent('vote', JSON.stringify(voteValue));
    eventTag('vote', 'battle', vote);
  };

  const handleUnvote = () => {
    vote = '';

    sendSocketEvent('retract_vote', battle.activePlanId);
    eventTag('retract_vote', 'battle', vote);
  };

  // Determine if the warrior has voted on active Plan yet
  function didVote(warriorId) {
    if (
      battle.activePlanId === '' ||
      (battle.votingLocked && battle.hideVoterIdentity)
    ) {
      return false;
    }
    const plan = battle.plans.find(p => p.id === battle.activePlanId);
    const voted = plan.votes.find(w => w.warriorId === warriorId);

    return voted !== undefined;
  }

  // Determine if we are showing users vote
  function showVote(warriorId) {
    if (
      battle.hideVoterIdentity ||
      battle.activePlanId === '' ||
      battle.votingLocked === false
    ) {
      return '';
    }
    const story = battle.plans.find(p => p.id === battle.activePlanId);
    const voted = story.votes.find(w => w.warriorId === warriorId);

    return voted !== undefined ? voted.vote : '';
  }

  // get highest vote from active story
  function getHighestVote() {
    const voteCounts = {};
    points.forEach(p => {
      voteCounts[p] = 0;
    });
    const highestVote = {
      vote: '',
      count: 0,
    };
    const activePlan = battle.plans.find(p => p.id === battle.activePlanId);

    if (activePlan.votes.length > 0) {
      const reversedPoints = [...points]
        .filter(v => v !== '?' && v !== '☕️')
        .reverse();
      reversedPoints.push('?');
      reversedPoints.push('☕️');

      // build a count of each vote
      activePlan.votes.forEach(v => {
        const voteWarrior = battle.users.find(w => w.id === v.warriorId) || {};
        const { spectator = false } = voteWarrior;

        if (typeof voteCounts[v.vote] !== 'undefined' && !spectator) {
          ++voteCounts[v.vote];
        }
      });

      // find the highest vote giving priority to higher numbers
      reversedPoints.forEach(p => {
        if (voteCounts[p] > highestVote.count) {
          highestVote.vote = p;
          highestVote.count = voteCounts[p];
        }
      });
    }

    return highestVote.vote;
  }

  $: highestVoteCount =
    battle.activePlanId !== '' && battle.votingLocked === true
      ? getHighestVote()
      : '';
  $: showVotingResults =
    battle.activePlanId !== '' && battle.votingLocked === true;

  $: isLeader = battle.leaders.includes($user.id);

  function concedeBattle() {
    eventTag('concede_battle', 'battle', '', () => {
      sendSocketEvent('concede_battle', '');
    });
  }

  function abandonBattle() {
    eventTag('abandon_battle', 'battle', '', () => {
      sendSocketEvent('abandon_battle', '');
    });
  }

  function toggleEditBattle() {
    showEditBattle = !showEditBattle;
  }

  const toggleDeleteBattle = () => {
    showDeleteBattle = !showDeleteBattle;
  };

  function handleBattleEdit(revisedBattle) {
    sendSocketEvent('revise_battle', JSON.stringify(revisedBattle));
    eventTag('revise_battle', 'battle', '');
    toggleEditBattle();
    battle.leaderCode = revisedBattle.leaderCode;
  }

  function authBattle(e) {
    e.preventDefault();

    sendSocketEvent('auth_battle', joinPasscode);
    eventTag('auth_battle', 'battle', '');
  }

  onMount(() => {
    if (!$user.id) {
      router.route(`${loginOrRegister}/battle/${battleId}`);
      return;
    }
  });
</script>

<svelte:head>
  <title
    >{$LL.battle()}
    {battle.name} | {$LL.appName()}</title
  >
</svelte:head>

<PageLayout>
  {#if battle.name && !socketReconnecting && !socketError}
    <div class="mb-6 flex flex-wrap">
      <div class="w-full text-center md:w-2/3 md:text-left">
        <h1
          class="text-4xl font-semibold font-rajdhani leading-tight dark:text-white flex items-center"
        >
          {#if currentStory.link}
            <a
              href="{currentStory.link}"
              target="_blank"
              class="text-blue-800 dark:text-sky-400 inline-block"
              data-testid="currentplan-link"
            >
              <ExternalLinkIcon class="w-8 h-8" />
            </a>
          {/if}
          {#if currentStory.type}
            &nbsp;<span
              class="inline-block text-lg text-gray-500
                            border-gray-300 border px-1 rounded dark:text-gray-300 dark:border-gray-500"
              data-testid="currentplan-type"
            >
              {currentStory.type}
            </span>
          {/if}
          {#if currentStory.referenceId}
            &nbsp;<span data-testid="currentplan-refid"
              >[{currentStory.referenceId}]</span
            >
          {/if}
          <span data-testid="currentplan-name"
            >{#if currentStory.name === ''}[{$LL.votingNotStarted()}]{:else}&nbsp;{currentStory.name}{/if}</span
          >
        </h1>
        <h2
          class="text-gray-700 dark:text-gray-300 text-3xl font-semibold font-rajdhani leading-tight"
          data-testid="battle-name"
        >
          {battle.name}
        </h2>
      </div>

      <div class="w-full md:w-1/3 text-center md:text-right">
        <VoteTimer
          currentStoryId="{currentStory.id}"
          votingLocked="{battle.votingLocked}"
          voteStartTime="{voteStartTime}"
        />
      </div>
    </div>

    <div class="flex flex-wrap mb-4 -mx-4">
      <div class="w-full lg:w-3/4 px-4">
        {#if showVotingResults}
          <VoteResults
            warriors="{battle.users}"
            plans="{battle.plans}"
            activePlanId="{battle.activePlanId}"
            points="{points}"
            highestVote="{highestVoteCount}"
            averageRounding="{battle.pointAverageRounding}"
            hideVoterIdentity="{battle.hideVoterIdentity}"
          />
        {:else}
          <div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
            {#each points as point}
              <div class="w-1/4 md:w-1/6 px-2 mb-4">
                <PointCard
                  point="{point}"
                  active="{vote === point}"
                  on:voted="{handleVote}"
                  on:voteRetraction="{handleUnvote}"
                  isLocked="{battle.votingLocked || isSpectator}"
                />
              </div>
            {/each}
          </div>
        {/if}

        <PokerStories
          plans="{battle.plans}"
          isLeader="{isLeader}"
          sendSocketEvent="{sendSocketEvent}"
          eventTag="{eventTag}"
          notifications="{notifications}"
          xfetch="{xfetch}"
          gameId="{battle.id}"
        />
      </div>

      <div class="w-full lg:w-1/4 px-4">
        <div class="bg-white dark:bg-gray-800 shadow-lg mb-4 rounded-lg">
          <div class="bg-blue-500 dark:bg-gray-700 p-4 rounded-t-lg">
            <h3
              class="text-3xl text-white leading-tight font-semibold font-rajdhani uppercase"
            >
              {$LL.warriors()}
            </h3>
          </div>

          {#each battle.users as war (war.id)}
            {#if war.active}
              <UserCard
                warrior="{war}"
                leaders="{battle.leaders}"
                isLeader="{isLeader}"
                voted="{didVote(war.id)}"
                points="{showVote(war.id)}"
                autoFinishVoting="{battle.autoFinishVoting}"
                sendSocketEvent="{sendSocketEvent}"
                eventTag="{eventTag}"
                notifications="{notifications}"
              />
            {/if}
          {/each}

          {#if isLeader}
            <VotingControls
              points="{points}"
              planId="{battle.activePlanId}"
              sendSocketEvent="{sendSocketEvent}"
              votingLocked="{battle.votingLocked}"
              highestVote="{highestVoteCount}"
              eventTag="{eventTag}"
            />
          {/if}
        </div>

        <div class="bg-white dark:bg-gray-800 shadow-lg p-4 mb-4 rounded-lg">
          <InviteUser
            hostname="{hostname}"
            battleId="{battle.id}"
            joinCode="{battle.joinCode}"
            notifications="{notifications}"
          />
          {#if isLeader}
            <div class="mt-4 text-right">
              <HollowButton
                color="blue"
                onClick="{toggleEditBattle}"
                testid="battle-edit"
              >
                {$LL.battleEdit()}
              </HollowButton>
              <HollowButton
                color="red"
                onClick="{toggleDeleteBattle}"
                testid="battle-delete"
              >
                {$LL.battleDelete()}
              </HollowButton>
            </div>
          {:else}
            <div class="mt-4 text-right">
              <HollowButton
                color="red"
                onClick="{abandonBattle}"
                testid="battle-abandon"
              >
                {$LL.battleAbandon()}
              </HollowButton>
            </div>
          {/if}
        </div>
      </div>
    </div>
    {#if showEditBattle}
      <EditPokerGame
        battleName="{battle.name}"
        points="{points}"
        votingLocked="{battle.votingLocked}"
        autoFinishVoting="{battle.autoFinishVoting}"
        pointAverageRounding="{battle.pointAverageRounding}"
        hideVoterIdentity="{battle.hideVoterIdentity}"
        handleBattleEdit="{handleBattleEdit}"
        toggleEditBattle="{toggleEditBattle}"
        joinCode="{battle.joinCode}"
        leaderCode="{battle.leaderCode}"
        teamId="{battle.teamId}"
        notifications="{notifications}"
        xfetch="{xfetch}"
      />
    {/if}
  {:else if JoinPassRequired}
    <div class="flex justify-center">
      <div class="w-full md:w-1/2 lg:w-1/3">
        <form
          on:submit="{authBattle}"
          class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
          name="authBattle"
        >
          <div class="mb-4">
            <label
              class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
              for="battleJoinCode"
            >
              {$LL.passCodeRequired()}
            </label>
            <TextInput
              bind:value="{joinPasscode}"
              placeholder="{$LL.enterPasscode()}"
              id="battleJoinCode"
              name="battleJoinCode"
              type="password"
              required
            />
          </div>

          <div class="text-right">
            <SolidButton type="submit">{$LL.battleJoin()}</SolidButton>
          </div>
        </form>
      </div>
    </div>
  {:else if socketReconnecting}
    <div class="flex items-center">
      <div class="flex-1 text-center">
        <h1 class="text-5xl text-teal-500 leading-tight font-bold">
          {$LL.battleSocketReconnecting()}
        </h1>
      </div>
    </div>
  {:else if socketError}
    <div class="flex items-center">
      <div class="flex-1 text-center">
        <h1 class="text-5xl text-red-500 leading-tight font-bold">
          {$LL.battleSocketError()}
        </h1>
      </div>
    </div>
  {:else}
    <div class="flex items-center">
      <div class="flex-1 text-center">
        <h1 class="text-5xl text-green-500 leading-tight font-bold">
          {$LL.battleLoading()}
        </h1>
      </div>
    </div>
  {/if}

  {#if showDeleteBattle}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteBattle}"
      handleDelete="{concedeBattle}"
      confirmText="{$LL.deleteBattleConfirmText()}"
      confirmBtnText="{$LL.deleteBattle()}"
    />
  {/if}
</PageLayout>

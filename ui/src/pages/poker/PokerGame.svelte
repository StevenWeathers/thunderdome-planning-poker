<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import Sockette from 'sockette';

  import PageLayout from '../../components/PageLayout.svelte';
  import PointCard from '../../components/poker/PointCard.svelte';
  import PokerStories from '../../components/poker/PokerStories.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import EditPokerGame from '../../components/poker/EditPokerGame.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import UserCard from '../../components/poker/UserCard.svelte';
  import VotingControls from '../../components/poker/VotingControls.svelte';
  import InviteUser from '../../components/poker/InviteUser.svelte';
  import VoteTimer from '../../components/poker/VoteTimer.svelte';
  import type { PokerGame, PokerStory } from '../../types/poker';
  import { ExternalLink } from 'lucide-svelte';
  import VotingMetrics from '../../components/poker/VotingMetrics.svelte';
  import FullpageLoader from '../../components/global/FullpageLoader.svelte';
  import JoinCodeForm from '../../components/global/JoinCodeForm.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    battleId: string;
    notifications: NotificationService;
    router: any;
    xfetch: ApiClient;
  }

  let {
    battleId,
    notifications,
    router,
    xfetch
  }: Props = $props();

  const { AllowRegistration, AllowGuests } = AppConfig;
  const loginOrRegister: string = AllowGuests
    ? appRoutes.register
    : appRoutes.login;

  const hostname: string = window.location.origin;

  const defaultStory: PokerStory = {
    id: '',
    active: false,
    points: '',
    priority: 0,
    skipped: false,
    voteEndTime: undefined,
    voteStartTime: undefined,
    votes: [],
    name: '',
    type: '',
    referenceId: '',
    link: '',
    description: '',
    acceptanceCriteria: '',
    position: 0,
  };

  let isLoading: boolean = $state(true);
  let JoinPassRequired: boolean = $state(false);
  let socketError: boolean = $state(false);
  let socketReconnecting: boolean = $state(false);
  let points: Array<string> = $state(['1', '2', '3', '5', '8', '13', '?']);
  let vote: string = $state('');
  let pokerGame: PokerGame = $state({
    leaders: [],
    autoFinishVoting: false,
    createdDate: undefined,
    hideVoterIdentity: false,
    id: '',
    name: '',
    plans: [],
    pointAverageRounding: '',
    pointValuesAllowed: [],
    updatedDate: undefined,
    users: [],
    votingLocked: false,
    teamId: '',
  });
  let currentStory = $state({ ...defaultStory });
  let showEditGame: boolean = $state(false);
  let showDeleteGame: boolean = $state(false);
  let isSpectator: boolean = $state(false);
  let voteStartTime: Date = $state(new Date());

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
      case 'init': {
        JoinPassRequired = false;
        pokerGame = JSON.parse(parsedEvent.value);
        points = pokerGame.pointValuesAllowed;
        const { spectator = false } =
          pokerGame.users.find(w => w.id === $user.id) || {};
        isSpectator = spectator;

        if (pokerGame.activePlanId !== '') {
          const activePlan = pokerGame.plans.find(
            p => p.id === pokerGame.activePlanId,
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

        break;
      }
      case 'user_joined': {
        pokerGame.users = JSON.parse(parsedEvent.value);
        const joinedWarrior = pokerGame.users.find(
          w => w.id === parsedEvent.userId,
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
      case 'user_left':
        const leftWarrior = pokerGame.users.find(
          w => w.id === parsedEvent.userId,
        );
        pokerGame.users = JSON.parse(parsedEvent.value);

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
        pokerGame.users = JSON.parse(parsedEvent.value);
        const updatedWarrior = pokerGame.users.find(w => w.id === $user.id);
        isSpectator = updatedWarrior.spectator;
        break;
      case 'plan_added':
        pokerGame.plans = JSON.parse(parsedEvent.value);
        break;
      case 'story_arranged':
        pokerGame.plans = JSON.parse(parsedEvent.value);
        break;
      case 'plan_activated':
        const updatedPlans = JSON.parse(parsedEvent.value);
        const activePlan = updatedPlans.find(p => p.active);
        currentStory = activePlan;
        voteStartTime = new Date(activePlan.voteStartTime);

        pokerGame.plans = updatedPlans;
        pokerGame.activePlanId = activePlan.id;
        pokerGame.votingLocked = false;
        vote = '';
        break;
      case 'plan_skipped':
        const updatedPlans2 = JSON.parse(parsedEvent.value);
        currentStory = { ...defaultStory };
        pokerGame.plans = updatedPlans2;
        pokerGame.activePlanId = '';
        pokerGame.votingLocked = true;
        vote = '';
        if ($user.notificationsEnabled) {
          notifications.warning($LL.planSkipped());
        }
        break;
      case 'vote_activity':
        const votedWarrior = pokerGame.users.find(
          w => w.id === parsedEvent.userId,
        );
        if ($user.notificationsEnabled) {
          notifications.success(
            `${$LL.warriorVoted({
              name: votedWarrior.name,
            })}
    `,
          );
        }

        pokerGame.plans = JSON.parse(parsedEvent.value);
        break;
      case 'vote_retracted':
        const devotedWarrior = pokerGame.users.find(
          w => w.id === parsedEvent.userId,
        );
        if ($user.notificationsEnabled) {
          notifications.warning(
            `${$LL.warriorRetractedVote({
              name: devotedWarrior.name,
            })}
    `,
          );
        }

        pokerGame.plans = JSON.parse(parsedEvent.value);
        break;
      case 'voting_ended':
        pokerGame.plans = JSON.parse(parsedEvent.value);
        pokerGame.votingLocked = true;
        break;
      case 'plan_finalized':
        pokerGame.plans = JSON.parse(parsedEvent.value);
        pokerGame.activePlanId = '';
        currentStory = { ...defaultStory };
        vote = '';
        break;
      case 'plan_revised':
        pokerGame.plans = JSON.parse(parsedEvent.value);
        if (pokerGame.activePlanId !== '') {
          const activePlan = pokerGame.plans.find(
            p => p.id === pokerGame.activePlanId,
          );
          currentStory = activePlan;
        }
        break;
      case 'plan_burned':
        const postBurnPlans = JSON.parse(parsedEvent.value);

        if (
          pokerGame.activePlanId !== '' &&
          postBurnPlans.filter(p => p.id === pokerGame.activePlanId).length ===
            0
        ) {
          pokerGame.activePlanId = '';
          currentStory = { ...defaultStory };
        }

        pokerGame.plans = postBurnPlans;

        break;
      case 'leaders_updated':
        pokerGame.leaders = parsedEvent.value;
        break;
      case 'battle_revised':
        const revisedBattle = JSON.parse(parsedEvent.value);
        pokerGame.name = revisedBattle.battleName;
        points = revisedBattle.pointValuesAllowed;
        pokerGame.autoFinishVoting = revisedBattle.autoFinishVoting;
        pokerGame.pointAverageRounding = revisedBattle.pointAverageRounding;
        pokerGame.joinCode = revisedBattle.joinCode;
        pokerGame.hideVoterIdentity = revisedBattle.hideVoterIdentity;
        pokerGame.teamId = revisedBattle.teamId;
        break;
      case 'battle_conceded':
        // poker over, goodbye.
        notifications.warning($LL.battleDeleted());
        router.route(appRoutes.games);
        break;
      case 'jab_warrior':
        const userToNudge = pokerGame.users.find(
          w => w.id === parsedEvent.value,
        );
        notifications.info(
          `${$LL.warriorNudgeMessage({
            name: userToNudge.name,
          })}
    `,
        );

        // msn wizz animation and sound
        if (userToNudge.id === $user.id) {
          document.querySelector('body').classList.toggle('shake');

          setTimeout(() => {
            document.querySelector('body').classList.toggle('shake');
          }, 700);
        }
        break;
      default:
        break;
    }
  };

  const ws = new Sockette(`${getWebsocketAddress()}/api/arena/${battleId}`, {
    timeout: 2e3,
    maxAttempts: 15,
    onmessage: onSocketMessage,
    onerror: err => {
      socketError = true;
    },
    onclose: e => {
      if (e.code === 4004) {
        router.route(appRoutes.games);
      } else if (e.code === 4001) {
        user.delete();
        router.route(`${loginOrRegister}/battle/${battleId}`);
      } else if (e.code === 4003) {
        notifications.danger($LL.sessionDuplicate());
        router.route(`${appRoutes.games}`);
      } else if (e.code === 4002) {
        router.route(appRoutes.games);
      } else {
        socketReconnecting = true;
      }
    },
    onopen: () => {
      socketError = false;
      socketReconnecting = false;
      isLoading = false;
    },
    onmaximum: () => {
      socketReconnecting = false;
    },
  });

  onDestroy(() => {
    ws.close();
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
      planId: pokerGame.activePlanId,
      voteValue: vote,
      autoFinishVoting: pokerGame.autoFinishVoting,
    };

    sendSocketEvent('vote', JSON.stringify(voteValue));
  };

  const handleUnvote = () => {
    vote = '';

    sendSocketEvent('retract_vote', pokerGame.activePlanId);
  };

  // Determine if the warrior has voted on active Plan yet
  function didVote(warriorId) {
    if (
      pokerGame.activePlanId === '' ||
      (pokerGame.votingLocked && pokerGame.hideVoterIdentity)
    ) {
      return false;
    }
    const plan = pokerGame.plans.find(p => p.id === pokerGame.activePlanId);
    const voted = plan.votes.find(w => w.warriorId === warriorId);

    return voted !== undefined;
  }

  // Determine if we are showing users vote
  function showVote(warriorId) {
    if (
      pokerGame.hideVoterIdentity ||
      pokerGame.activePlanId === '' ||
      pokerGame.votingLocked === false
    ) {
      return '';
    }
    const story = pokerGame.plans.find(p => p.id === pokerGame.activePlanId);
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
    const activePlan = pokerGame.plans.find(
      p => p.id === pokerGame.activePlanId,
    );

    if (activePlan.votes.length > 0) {
      const reversedPoints = [...points]
        .filter(v => v !== '?' && v !== '☕️')
        .reverse();
      reversedPoints.push('?');
      reversedPoints.push('☕️');

      // build a count of each vote
      activePlan.votes.forEach(v => {
        const voteWarrior =
          pokerGame.users.find(w => w.id === v.warriorId) || {};
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

  let highestVoteCount =
    $derived(pokerGame.activePlanId !== '' && pokerGame.votingLocked === true
      ? getHighestVote()
      : '');
  let showVotingResults =
    $derived(pokerGame.activePlanId !== '' && pokerGame.votingLocked === true);

  let isFacilitator = $derived(pokerGame.leaders.includes($user.id));

  function concedeGame() {
    sendSocketEvent('concede_battle', '');
  }

  function abandonBattle() {
    sendSocketEvent('abandon_battle', '');
  }

  function toggleEditGame() {
    showEditGame = !showEditGame;
  }

  const toggleDeleteGame = () => {
    showDeleteGame = !showDeleteGame;
  };

  function handleGameEdit(revisedBattle) {
    sendSocketEvent('revise_battle', JSON.stringify(revisedBattle));
    toggleEditGame();
    pokerGame.leaderCode = revisedBattle.leaderCode;
  }

  function authBattle(joinPasscode) {
    sendSocketEvent('auth_game', joinPasscode);
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
    {pokerGame.name} | {$LL.appName()}</title
  >
</svelte:head>

<PageLayout>
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
            <ExternalLink class="w-8 h-8" />
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
        {pokerGame.name}
      </h2>
    </div>

    <div class="w-full md:w-1/3 text-center md:text-right">
      <VoteTimer
        currentStoryId={currentStory.id}
        votingLocked={pokerGame.votingLocked}
        voteStartTime={voteStartTime}
      />
    </div>
  </div>

  <div class="flex flex-wrap mb-4 -mx-4">
    <div class="w-full lg:w-3/4 px-4">
      {#if showVotingResults}
        <div class=" mb-2 md:mb-4">
          <VotingMetrics
            pointValues={points}
            votes={pokerGame.plans.find(p => p.id === pokerGame.activePlanId)
              .votes}
            users={pokerGame.users}
            averageRounding={pokerGame.pointAverageRounding}
          />
        </div>
      {:else}
        <div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
          {#each points as point}
            <div class="w-1/4 md:w-1/6 px-2 mb-4">
              <PointCard
                point={point}
                active={vote === point}
                on:voted="{handleVote}"
                on:voteRetraction="{handleUnvote}"
                isLocked={pokerGame.votingLocked || isSpectator}
              />
            </div>
          {/each}
        </div>
      {/if}

      <PokerStories
        plans={pokerGame.plans}
        isFacilitator={isFacilitator}
        sendSocketEvent={sendSocketEvent}
        notifications={notifications}
        xfetch={xfetch}
        gameId={pokerGame.id}
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

        {#each pokerGame.users as war (war.id)}
          {#if war.active}
            <UserCard
              warrior={war}
              leaders={pokerGame.leaders}
              isFacilitator={isFacilitator}
              voted={didVote(war.id)}
              points={showVote(war.id)}
              autoFinishVoting={pokerGame.autoFinishVoting}
              sendSocketEvent={sendSocketEvent}
              notifications={notifications}
            />
          {/if}
        {/each}

        {#if isFacilitator}
          <VotingControls
            points={points}
            planId={pokerGame.activePlanId}
            sendSocketEvent={sendSocketEvent}
            votingLocked={pokerGame.votingLocked}
            highestVote={highestVoteCount}
          />
        {/if}
      </div>

      <div class="bg-white dark:bg-gray-800 shadow-lg p-4 mb-4 rounded-lg">
        <InviteUser
          hostname={hostname}
          battleId={pokerGame.id}
          joinCode={pokerGame.joinCode}
          notifications={notifications}
        />
        {#if isFacilitator}
          <div class="mt-4 text-right">
            <HollowButton
              color="blue"
              onClick={toggleEditGame}
              testid="battle-edit"
            >
              {$LL.battleEdit()}
            </HollowButton>
            <HollowButton
              color="red"
              onClick={toggleDeleteGame}
              testid="battle-delete"
            >
              {$LL.battleDelete()}
            </HollowButton>
          </div>
        {:else}
          <div class="mt-4 text-right">
            <HollowButton
              color="red"
              onClick={abandonBattle}
              testid="battle-abandon"
            >
              {$LL.battleAbandon()}
            </HollowButton>
          </div>
        {/if}
      </div>
    </div>
  </div>

  {#if showEditGame}
    <EditPokerGame
      battleName={pokerGame.name}
      points={points}
      votingLocked={pokerGame.votingLocked}
      autoFinishVoting={pokerGame.autoFinishVoting}
      pointAverageRounding={pokerGame.pointAverageRounding}
      hideVoterIdentity={pokerGame.hideVoterIdentity}
      handleBattleEdit={handleGameEdit}
      toggleEditBattle={toggleEditGame}
      joinCode={pokerGame.joinCode}
      leaderCode={pokerGame.leaderCode}
      teamId={pokerGame.teamId}
      notifications={notifications}
      xfetch={xfetch}
    />
  {/if}

  {#if showDeleteGame}
    <DeleteConfirmation
      toggleDelete={toggleDeleteGame}
      handleDelete={concedeGame}
      confirmText={$LL.deleteBattleConfirmText()}
      confirmBtnText={$LL.deleteBattle()}
    />
  {/if}

  {#if socketReconnecting}
    <FullpageLoader>
      {$LL.battleSocketReconnecting()}
    </FullpageLoader>
  {:else if socketError}
    <FullpageLoader>
      {$LL.battleSocketError()}
    </FullpageLoader>
  {:else if isLoading}
    <FullpageLoader>
      {$LL.battleLoading()}
    </FullpageLoader>
  {:else if JoinPassRequired}
    <JoinCodeForm handleSubmit={authBattle} submitText={$LL.battleJoin()} />
  {/if}
</PageLayout>

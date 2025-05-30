<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import Sockette from 'sockette';
  import { onDestroy, onMount } from 'svelte';

  import AddGoal from '../../components/storyboard/AddGoal.svelte';
  import UserCard from '../../components/storyboard/UserCard.svelte';
  import InviteUser from '../../components/storyboard/InviteUser.svelte';
  import ColumnForm from '../../components/storyboard/ColumnForm.svelte';
  import StoryForm from '../../components/storyboard/StoryForm.svelte';
  import ColorLegendForm from '../../components/storyboard/ColorLegendForm.svelte';
  import PersonasForm from '../../components/storyboard/PersonasForm.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import DeleteStoryboard from '../../components/storyboard/DeleteStoryboard.svelte';
  import EditStoryboard from '../../components/storyboard/EditStoryboard.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import BecomeFacilitator from '../../components/BecomeFacilitator.svelte';
  import GoalEstimate from '../../components/storyboard/GoalEstimate.svelte';
  import {
    ChevronDown,
    ChevronUp,
    MessageSquareMore,
    Pencil,
    User,
    Users,
  } from 'lucide-svelte';
  import JoinCodeForm from '../../components/global/JoinCodeForm.svelte';
  import FullpageLoader from '../../components/global/FullpageLoader.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';

  interface Props {
    storyboardId: any;
    notifications: any;
    router: any;
  }

  let { storyboardId, notifications, router }: Props = $props();

  const { AllowRegistration, AllowGuests } = AppConfig;
  const loginOrRegister = AllowGuests ? appRoutes.register : appRoutes.login;

  const hostname = window.location.origin;

  let isLoading = $state(true);
  let JoinPassRequired = $state(false);
  let socketError = $state(false);
  let socketReconnecting = $state(false);
  let storyboard = $state({
    goals: [],
    users: [],
    colorLegend: [],
    personas: [],
    facilitators: [],
    facilitatorCode: '',
    joinCode: '',
  });
  let showUsers = $state(false);
  let showColorLegend = $state(false);
  let showColorLegendForm = $state(false);
  let showPersonas = $state(false);
  let showPersonasForm = $state(null);
  let editColumn = $state(null);
  let activeStory = $state(null);
  let showDeleteStoryboard = $state(false);
  let showEditStoryboard = $state(false);
  let collapseGoals = $state([]);

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
        storyboard = JSON.parse(parsedEvent.value);
        break;
      case 'user_joined':
        storyboard.users = JSON.parse(parsedEvent.value);
        const joinedUser = storyboard.users.find(
          w => w.id === parsedEvent.userId,
        );
        notifications.success(`${joinedUser.name} joined.`);
        break;
      case 'user_retreated':
        const leftUser = storyboard.users.find(
          w => w.id === parsedEvent.userId,
        );
        storyboard.users = JSON.parse(parsedEvent.value);

        notifications.danger(`${leftUser.name} left.`);
        break;
      case 'storyboard_updated':
        storyboard = JSON.parse(parsedEvent.value);
        break;
      case 'goal_added':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'goal_revised':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'goal_deleted':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'column_added':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'column_updated':
        storyboard.goals = JSON.parse(parsedEvent.value);
        if (editColumn !== null) {
          storyboard.goals.map(goal => {
            goal.columns.map(column => {
              if (column.id === editColumn.id) {
                editColumn = column;
              }
            });
          });
        }
        break;
      case 'story_added':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'story_updated':
        storyboard.goals = JSON.parse(parsedEvent.value);
        if (activeStory) {
          let activeStoryFound = false;
          for (let goal of storyboard.goals) {
            for (let column of goal.columns) {
              for (let story of column.stories) {
                if (story.id === activeStory.id) {
                  activeStory = story;
                  activeStoryFound = true;
                  break;
                }
              }
              if (activeStoryFound) {
                break;
              }
            }
            if (activeStoryFound) {
              break;
            }
          }
        }
        break;
      case 'story_moved':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'story_deleted':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'personas_updated':
        storyboard.personas = JSON.parse(parsedEvent.value);
        break;
      case 'storyboard_edited':
        const revisedStoryboard = JSON.parse(parsedEvent.value);
        storyboard.name = revisedStoryboard.storyboardName;
        storyboard.joinCode = revisedStoryboard.joinCode;
        break;
      case 'storyboard_conceded':
        // storyboard over, goodbye.
        notifications.warning($LL.storyboardDeleted());
        router.route(appRoutes.storyboards);
        break;
      default:
        break;
    }
  };

  const ws = new Sockette(
    `${getWebsocketAddress()}/api/storyboard/${storyboardId}`,
    {
      timeout: 2e3,
      maxAttempts: 15,
      onmessage: onSocketMessage,
      onerror: () => {
        socketError = true;
      },
      onclose: e => {
        if (e.code === 4004) {
          router.route(appRoutes.storyboards);
        } else if (e.code === 4001) {
          user.delete();
          router.route(`${loginOrRegister}/storyboard/${storyboardId}`);
        } else if (e.code === 4003) {
          notifications.danger($LL.duplicateStoryboardSession());
          router.route(`${appRoutes.storyboards}`);
        } else if (e.code === 4002) {
          router.route(appRoutes.storyboards);
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
    },
  );

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

  // event handlers
  function handleDndConsider(e) {
    const goalIndex = e.target.dataset.goalindex;
    const columnIndex = e.target.dataset.columnindex;

    storyboard.goals[goalIndex].columns[columnIndex].stories = e.detail.items;
    storyboard.goals = storyboard.goals;
  }

  function handleDndFinalize(e) {
    const goalIndex = e.target.dataset.goalindex;
    const columnIndex = e.target.dataset.columnindex;
    const storyId = e.detail.info.id;

    storyboard.goals[goalIndex].columns[columnIndex].stories = e.detail.items;
    storyboard.goals = storyboard.goals;

    const matchedStory = storyboard.goals[goalIndex].columns[
      columnIndex
    ].stories.find(i => i.id === storyId);

    if (matchedStory) {
      const goalId = storyboard.goals[goalIndex].id;
      const columnId = storyboard.goals[goalIndex].columns[columnIndex].id;

      // determine what story to place story before in target column
      const matchedStoryIndex =
        storyboard.goals[goalIndex].columns[columnIndex].stories.indexOf(
          matchedStory,
        );
      const sibling =
        storyboard.goals[goalIndex].columns[columnIndex].stories[
          matchedStoryIndex + 1
        ];
      const placeBefore = sibling ? sibling.id : '';

      sendSocketEvent(
        'move_story',
        JSON.stringify({
          storyId,
          goalId,
          columnId,
          placeBefore,
        }),
      );
    }
  }

  function authStoryboard(joinPasscode) {
    sendSocketEvent('auth_storyboard', joinPasscode);
  }

  const addStory = (goalId, columnId) => () => {
    sendSocketEvent(
      'add_story',
      JSON.stringify({
        goalId,
        columnId,
      }),
    );
  };

  const addStoryColumn = goalId => () => {
    sendSocketEvent(
      'add_column',
      JSON.stringify({
        goalId,
      }),
    );
  };

  const deleteColumn = columnId => () => {
    sendSocketEvent('delete_column', columnId);
    toggleColumnEdit()();
  };

  const handleAddFacilitator = userId => () => {
    sendSocketEvent(
      'facilitator_add',
      JSON.stringify({
        userId,
      }),
    );
  };

  const handleRemoveFacilitator = userId => () => {
    if (storyboard.facilitators.length === 1) {
      notifications.danger($LL.removeOnlyFacilitatorError());
      return;
    }

    sendSocketEvent(
      'facilitator_remove',
      JSON.stringify({
        userId,
      }),
    );
  };

  function concedeStoryboard() {
    sendSocketEvent('concede_storyboard', '');
  }

  function abandonStoryboard() {
    sendSocketEvent('abandon_storyboard', '');
  }

  function toggleUsersPanel() {
    showColorLegend = false;
    showPersonas = false;
    showUsers = !showUsers;
  }

  function toggleColorLegend() {
    showUsers = false;
    showPersonas = false;
    showColorLegend = !showColorLegend;
  }

  function togglePersonas() {
    showUsers = false;
    showColorLegend = false;
    showPersonas = !showPersonas;
  }

  function toggleColumnEdit(column) {
    return () => {
      editColumn = editColumn != null ? null : column;
    };
  }

  function toggleEditLegend() {
    showColorLegend = false;
    showColorLegendForm = !showColorLegendForm;
  }

  const toggleEditPersona = persona => () => {
    showPersonas = false;
    showPersonasForm = showPersonasForm != null ? null : persona;
  };

  const toggleDeleteStoryboard = () => {
    showDeleteStoryboard = !showDeleteStoryboard;
  };

  let showAddGoal = $state(false);
  let reviseGoalId = $state('');
  let reviseGoalName = $state('');

  const toggleAddGoal = goalId => () => {
    if (goalId) {
      const goalName = storyboard.goals.find(p => p.id === goalId).name;
      reviseGoalId = goalId;
      reviseGoalName = goalName;
    } else {
      reviseGoalId = '';
      reviseGoalName = '';
    }
    showAddGoal = !showAddGoal;
  };

  const handleGoalAdd = goalName => {
    sendSocketEvent('add_goal', goalName);
  };

  const handleGoalRevision = updatedGoal => {
    sendSocketEvent('revise_goal', JSON.stringify(updatedGoal));
  };

  const handleGoalDeletion = goalId => () => {
    sendSocketEvent('delete_goal', goalId);
  };

  const handleColumnRevision = column => {
    sendSocketEvent('revise_column', JSON.stringify(column));
  };

  const handleLegendRevision = legend => {
    sendSocketEvent('revise_color_legend', JSON.stringify(legend));
  };

  const handlePersonaAdd = persona => {
    sendSocketEvent('add_persona', JSON.stringify(persona));
  };

  const handleColumnPersonaAdd = column_persona => {
    sendSocketEvent('column_persona_add', JSON.stringify(column_persona));
  };
  const handleColumnPersonaRemove = column_persona => () => {
    sendSocketEvent('column_persona_remove', JSON.stringify(column_persona));
  };

  const handlePersonaRevision = persona => {
    sendSocketEvent('revise_persona', JSON.stringify(persona));
  };

  const handleDeletePersona = personaId => () => {
    sendSocketEvent('delete_persona', personaId);
  };

  function handleStoryboardEdit(revisedStoryboard) {
    sendSocketEvent('edit_storyboard', JSON.stringify(revisedStoryboard));
    toggleEditStoryboard();
  }

  function toggleEditStoryboard() {
    showEditStoryboard = !showEditStoryboard;
  }

  const toggleStoryForm = story => () => {
    activeStory = activeStory != null ? null : story;
  };

  function toggleGoalCollapse(goalId) {    
    return () => {
      if (collapseGoals.includes(goalId)) {
        collapseGoals = collapseGoals.filter(g => g !== goalId);
      } else {
        collapseGoals.push(goalId);
      }
    };
  }

  let showBecomeFacilitator = $state(false);

  function becomeFacilitator(facilitatorCode) {
    sendSocketEvent('facilitator_self', facilitatorCode);
    toggleBecomeFacilitator();
  }

  function toggleBecomeFacilitator() {
    showBecomeFacilitator = !showBecomeFacilitator;
  }

  let isFacilitator =
    $derived(storyboard.facilitators && storyboard.facilitators.includes($user.id));

  onMount(() => {
    if (!$user.id) {
      router.route(`${loginOrRegister}/storyboard/${storyboardId}`);
    }
  });
</script>

<style>
  .story-gray {
    @apply border-gray-400;
  }

  .story-gray:hover {
    @apply border-gray-800;
  }

  .story-red {
    @apply border-red-400;
  }

  .story-red:hover {
    @apply border-red-800;
  }

  .story-orange {
    @apply border-orange-400;
  }

  .story-orange:hover {
    @apply border-orange-800;
  }

  .story-yellow {
    @apply border-yellow-400;
  }

  .story-yellow:hover {
    @apply border-yellow-800;
  }

  .story-green {
    @apply border-green-400;
  }

  .story-green:hover {
    @apply border-green-800;
  }

  .story-teal {
    @apply border-teal-400;
  }

  .story-teal:hover {
    @apply border-teal-800;
  }

  .story-blue {
    @apply border-blue-400;
  }

  .story-blue:hover {
    @apply border-blue-800;
  }

  .story-indigo {
    @apply border-indigo-400;
  }

  .story-indigo:hover {
    @apply border-indigo-800;
  }

  .story-purple {
    @apply border-purple-400;
  }

  .story-purple:hover {
    @apply border-purple-800;
  }

  .story-pink {
    @apply border-pink-400;
  }

  .story-pink:hover {
    @apply border-pink-800;
  }

  .colorcard-gray {
    @apply bg-gray-400;
  }

  .colorcard-red {
    @apply bg-red-400;
  }

  .colorcard-orange {
    @apply bg-orange-400;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
  }

  .colorcard-green {
    @apply bg-green-400;
  }

  .colorcard-teal {
    @apply bg-teal-400;
  }

  .colorcard-blue {
    @apply bg-blue-400;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
  }

  .colorcard-purple {
    @apply bg-purple-400;
  }

  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>

<svelte:head>
  <title>{$LL.storyboard()} {storyboard.name} | {$LL.appName()}</title>
</svelte:head>

<div class="w-full">
  <div
    class="px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b border-t border-gray-400 dark:border-gray-700 flex
        flex-wrap"
  >
    <div class="w-1/3">
      <h1 class="text-3xl font-bold leading-tight dark:text-gray-200">
        {storyboard.name}
      </h1>
    </div>
    <div class="w-2/3 text-right">
      <div>
        {#if isFacilitator}
          <HollowButton
            color="green"
            onClick={toggleAddGoal()}
            additionalClasses="me-2"
            testid="goal-add"
          >
            {$LL.storyboardAddGoal()}
          </HollowButton>
          <HollowButton
            color="blue"
            onClick={toggleEditStoryboard}
            testid="storyboard-edit"
          >
            {$LL.editStoryboard()}
          </HollowButton>
          <HollowButton
            color="red"
            onClick={toggleDeleteStoryboard}
            additionalClasses="me-2"
            testid="storyboard-delete"
          >
            {$LL.deleteStoryboard()}
          </HollowButton>
        {:else}
          <HollowButton
            color="blue"
            onClick={toggleBecomeFacilitator}
            testid="become-facilitator"
          >
            {$LL.becomeFacilitator()}
          </HollowButton>
          <HollowButton
            color="red"
            onClick={abandonStoryboard}
            testid="storyboard-leave"
          >
            {$LL.leaveStoryboard()}
          </HollowButton>
        {/if}
        <div class="inline-block relative">
          <HollowButton
            color="indigo"
            additionalClasses="transition ease-in-out duration-150"
            onClick={togglePersonas}
            testid="personas-toggle"
          >
            {$LL.personas()}
            <ChevronDown class="ms-1 inline-block" />
          </HollowButton>
          {#if showPersonas}
            <div
              class="origin-top-right absolute end-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
            >
              <div
                class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
              >
                <div class="p-2">
                  {#each storyboard.personas as persona}
                    <div class="mb-1 w-full">
                      <div>
                        <span class="font-bold">
                          {persona.name}
                        </span>
                        {#if isFacilitator}
                          &nbsp;|&nbsp;
                          <button
                            onclick={toggleEditPersona(persona)}
                            class="text-orange-500
                                                        hover:text-orange-800"
                            data-testid="persona-edit"
                          >
                            {$LL.edit()}
                          </button>
                          &nbsp;|&nbsp;
                          <button
                            onclick={handleDeletePersona(persona.id)}
                            class="text-red-500
                                                        hover:text-red-800"
                            data-testid="persona-delete"
                          >
                            {$LL.delete()}
                          </button>
                        {/if}
                      </div>
                      <span class="text-sm">
                        {persona.role}
                      </span>
                    </div>
                  {/each}
                </div>

                {#if isFacilitator}
                  <div class="p-2 text-right">
                    <HollowButton
                      color="green"
                      onClick={toggleEditPersona({
                        id: '',
                        name: '',
                        role: '',
                        description: '',
                      })}
                      testid="persona-add"
                    >
                      {$LL.addPersona()}
                    </HollowButton>
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        </div>
        <div class="inline-block relative">
          <HollowButton
            color="teal"
            additionalClasses="transition ease-in-out duration-150"
            onClick={toggleColorLegend}
            testid="colorlegend-toggle"
          >
            {$LL.colorLegend()}
            <ChevronDown class="ms-1 inline-block" />
          </HollowButton>
          {#if showColorLegend}
            <div
              class="origin-top-right absolute end-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
            >
              <div
                class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
              >
                <div class="p-2">
                  {#each storyboard.color_legend as color}
                    <div class="mb-1 flex w-full">
                      <span
                        class="p-4 me-2 inline-block
                                                colorcard-{color.color}"></span>
                      <span
                        class="inline-block align-middle
                                                {color.legend === ''
                          ? 'text-gray-300 dark:text-gray-500'
                          : 'text-gray-600 dark:text-gray-200'}"
                      >
                        {color.legend || $LL.colorLegendNotSpecified()}
                      </span>
                    </div>
                  {/each}
                </div>

                {#if isFacilitator}
                  <div class="p-2 text-right">
                    <HollowButton
                      color="orange"
                      onClick={toggleEditLegend}
                      testid="colorlegend-edit"
                    >
                      {$LL.editColorLegend()}
                    </HollowButton>
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        </div>
        <div class="inline-block relative">
          <HollowButton
            color="orange"
            additionalClasses="transition ease-in-out duration-150"
            onClick={toggleUsersPanel}
            testid="users-toggle"
          >
            <Users class="me-1 inline-block" height="18" width="18" />
            {$LL.users()}
            <ChevronDown class="ms-1 inline-block" />
          </HollowButton>
          {#if showUsers}
            <div
              class="origin-top-right absolute end-0 mt-1 w-64
                            rounded-md shadow-lg text-left z-10"
            >
              <div
                class="rounded-md bg-white dark:bg-gray-700 dark:text-white shadow-xs"
              >
                {#each storyboard.users as usr, index (usr.id)}
                  {#if usr.active}
                    <UserCard
                      user={usr}
                      showBorder={index !== storyboard.users.length - 1}
                      facilitators={storyboard.facilitators}
                      handleAddFacilitator={handleAddFacilitator}
                      handleRemoveFacilitator={handleRemoveFacilitator}
                    />
                  {/if}
                {/each}

                <div class="p-2">
                  <InviteUser
                    hostname={hostname}
                    storyboardId={storyboard.id}
                  />
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
  {#each storyboard.goals as goal, goalIndex (goal.id)}
    <div data-goalid="{goal.id}" data-testid="storyboard-goal">
      <div
        class="flex px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b-2 border-gray-400 dark:border-gray-700 {goalIndex >
        0
          ? 'border-t-2'
          : ''}"
      >
        <div class="w-3/4 relative">
          <div class="font-bold dark:text-gray-200 text-xl">
            <h2 class="inline-block align-middle pt-1">
              <button
                onclick={toggleGoalCollapse(goal.id)}
                data-testid="goal-expand"
                data-collapsed="{collapseGoals.includes(goal.id)}"
              >
                {#if collapseGoals.includes(goal.id)}
                  <ChevronDown class="me-1 inline-block" />
                {:else}
                  <ChevronUp class="me-1 inline-block" />
                {/if}
              </button>{goal.name}&nbsp;<GoalEstimate
                columns={goal.columns}
              />
            </h2>
          </div>
        </div>
        <div class="w-1/4 text-right">
          {#if isFacilitator}
            <HollowButton
              color="green"
              onClick={addStoryColumn(goal.id)}
              btnSize="small"
              testid="column-add"
            >
              {$LL.storyboardAddColumn()}
            </HollowButton>
            <HollowButton
              color="orange"
              onClick={toggleAddGoal(goal.id)}
              btnSize="small"
              additionalClasses="ms-2"
              testid="goal-edit"
            >
              {$LL.edit()}
            </HollowButton>
            <HollowButton
              color="red"
              onClick={handleGoalDeletion(goal.id)}
              btnSize="small"
              additionalClasses="ms-2"
              testid="goal-delete"
            >
              {$LL.delete()}
            </HollowButton>
          {/if}
        </div>
      </div>
      {#if !collapseGoals.includes(goal.id)}
        <section class="px-2" style="overflow-x: scroll">
          <div class="flex">
            {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
              <div class="flex-none mx-2 w-40" data-testid="goal-personas">
                <div class="w-full mb-2">
                  {#each goalColumn.personas as persona}
                    <div
                      class="mt-4 dark:text-gray-300 text-right"
                      data-testid="goal-persona"
                    >
                      <div class="font-bold" data-testid="persona-name">
                        <User class="inline-block h-4 w-4" />
                        {persona.name}
                      </div>
                      <div class="text-sm" data-testid="persona-role">
                        {persona.role}
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            {/each}
          </div>
          <div class="flex">
            {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
              <div class="flex-none my-4 mx-2 w-40" data-testid="goal-column">
                <div class="flex-none">
                  <div class="w-full mb-2">
                    <div class="flex">
                      <span
                        class="font-bold flex-grow truncate dark:text-gray-300"
                        title="{goalColumn.name}"
                        data-testid="column-name"
                      >
                        {goalColumn.name}
                      </span>
                      <button
                        onclick={toggleColumnEdit(goalColumn)}
                        class="flex-none font-bold text-xl
                                        border-dashed border-2 border-gray-400 dark:border-gray-600
                                        hover:border-green-500 text-gray-600 dark:text-gray-400
                                        hover:text-green-500 py-1 px-2"
                        title="{$LL.storyboardEditColumn()}"
                        data-testid="column-edit"
                      >
                        <Pencil />
                      </button>
                    </div>
                  </div>
                  <div class="w-full">
                    <div class="flex">
                      <button
                        onclick={addStory(goal.id, goalColumn.id)}
                        class="flex-grow font-bold text-xl py-1
                                        px-2 border-dashed border-2
                                        border-gray-400 dark:border-gray-600 hover:border-green-500
                                        text-gray-600 dark:text-gray-400 hover:text-green-500"
                        title="{$LL.storyboardAddStoryToColumn()}"
                        data-testid="story-add"
                      >
                        +
                      </button>
                    </div>
                  </div>
                </div>
                <div
                  class="w-full relative"
                  data-testid="column-dropzone"
                  style="min-height: 160px;"
                  data-goalid="{goal.id}"
                  data-columnid="{goalColumn.id}"
                  data-goalIndex="{goalIndex}"
                  data-columnindex="{columnIndex}"
                  use:dndzone="{{
                    items: goalColumn.stories,
                    type: 'story',
                    dropTargetStyle: '',
                    dropTargetClasses: [
                      'outline',
                      'outline-2',
                      'outline-indigo-500',
                      'dark:outline-yellow-400',
                    ],
                  }}"
                  onconsider={handleDndConsider}
                  onfinalize={handleDndFinalize}
                >
                  {#each goalColumn.stories as story (story.id)}
                    <div
                      class="relative max-w-xs shadow bg-white dark:bg-gray-700 dark:text-white border-s-4
                                    story-{story.color} border my-4
                                    cursor-pointer"
                      style="list-style: none;"
                      role="button"
                      tabindex="0"
                      data-goalid="{goal.id}"
                      data-columnid="{goalColumn.id}"
                      data-storyid="{story.id}"
                      data-testid="column-story"
                      onclick={toggleStoryForm(story)}
                      onkeypress={toggleStoryForm(story)}
                    >
                      <div>
                        <div>
                          <div
                            class="h-20 p-1 text-sm
                                                overflow-hidden {story.closed
                              ? 'line-through'
                              : ''}"
                            title="{story.name}"
                            data-testid="story-name"
                          >
                            {story.name}
                          </div>
                          <div class="h-8">
                            <div
                              class="flex content-center
                                                    p-1 text-sm"
                            >
                              <div
                                class="w-1/2
                                                        text-gray-600 dark:text-gray-300"
                              >
                                {#if story.comments.length > 0}
                                  <span
                                    class="inline-block
                                                                align-middle"
                                    data-testid="story-comments"
                                  >
                                    {story.comments.length}
                                    <MessageSquareMore class="inline-block" />
                                  </span>
                                {/if}
                              </div>
                              <div class="w-1/2 text-right">
                                {#if story.points > 0}
                                  <span
                                    class="px-2
                                                                bg-gray-300 dark:bg-gray-500
                                                                inline-block
                                                                align-middle"
                                    data-testid="story-points"
                                  >
                                    {story.points}
                                  </span>
                                {/if}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                      {#if story[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
                        <div
                          class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible opacity-50 max-w-xs shadow bg-white dark:bg-gray-700 dark:text-white border-s-4
                                    story-{story.color} border
                                    cursor-pointer"
                          style="list-style: none;"
                          role="button"
                          tabindex="0"
                          data-goalid="{goal.id}"
                          data-columnid="{goalColumn.id}"
                          data-storyid="{story.id}"
                          data-testid="column-story-shadowitem"
                          onclick={toggleStoryForm(story)}
                          onkeypress={toggleStoryForm(story)}
                        >
                          <div>
                            <div>
                              <div
                                class="h-20 p-1 text-sm
                                                overflow-hidden {story.closed
                                  ? 'line-through'
                                  : ''}"
                                title="{story.name}"
                                data-testid="shadow-story-name"
                              >
                                {story.name}
                              </div>
                              <div class="h-8">
                                <div
                                  class="flex content-center
                                                    p-1 text-sm"
                                >
                                  <div
                                    class="w-1/2
                                                        text-gray-600"
                                  >
                                    {#if story.comments.length > 0}
                                      <span
                                        class="inline-block
                                                                align-middle"
                                        data-testid="shadow-story-comments"
                                      >
                                        {story.comments.length}
                                        <MessageSquareMore
                                          class="inline-block"
                                        />
                                      </span>
                                    {/if}
                                  </div>
                                  <div class="w-1/2 text-right">
                                    {#if story.points > 0}
                                      <span
                                        class="px-2
                                                                bg-gray-300
                                                                inline-block
                                                                align-middle"
                                        data-testid="shadow-story-points"
                                      >
                                        {story.points}
                                      </span>
                                    {/if}
                                  </div>
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>
                      {/if}
                    </div>
                  {/each}
                </div>
              </div>
            {/each}
          </div>
        </section>
      {/if}
    </div>
  {/each}
</div>

{#if showAddGoal}
  <AddGoal
    handleGoalAdd={handleGoalAdd}
    toggleAddGoal={toggleAddGoal()}
    handleGoalRevision={handleGoalRevision}
    goalId={reviseGoalId}
    goalName={reviseGoalName}
  />
{/if}

{#if editColumn}
  <ColumnForm
    handleColumnRevision={handleColumnRevision}
    toggleColumnEdit={toggleColumnEdit()}
    column={editColumn}
    personas={storyboard.personas}
    handlePersonaAdd={handleColumnPersonaAdd}
    handlePersonaRemove={handleColumnPersonaRemove}
    deleteColumn={deleteColumn}
  />
{/if}

{#if activeStory}
  <StoryForm
    toggleStoryForm={toggleStoryForm()}
    story={activeStory}
    sendSocketEvent={sendSocketEvent}
    notifications={notifications}
    colorLegend={storyboard.color_legend}
    users={storyboard.users}
  />
{/if}

{#if showColorLegendForm}
  <ColorLegendForm
    handleLegendRevision={handleLegendRevision}
    toggleEditLegend={toggleEditLegend}
    colorLegend={storyboard.color_legend}
  />
{/if}

{#if showPersonasForm}
  <PersonasForm
    toggleEditPersona={toggleEditPersona()}
    persona={showPersonasForm}
    handlePersonaAdd={handlePersonaAdd}
    handlePersonaRevision={handlePersonaRevision}
  />
{/if}

{#if showEditStoryboard}
  <EditStoryboard
    storyboardName={storyboard.name}
    handleStoryboardEdit={handleStoryboardEdit}
    toggleEditStoryboard={toggleEditStoryboard}
    joinCode={storyboard.joinCode}
    facilitatorCode={storyboard.facilitatorCode}
  />
{/if}

{#if showDeleteStoryboard}
  <DeleteStoryboard
    toggleDelete={toggleDeleteStoryboard}
    handleDelete={concedeStoryboard}
  />
{/if}

{#if showBecomeFacilitator}
  <BecomeFacilitator
    handleBecomeFacilitator={becomeFacilitator}
    toggleBecomeFacilitator={toggleBecomeFacilitator}
  />
{/if}

{#if socketReconnecting}
  <FullpageLoader>
    {$LL.reloadingStoryboard()}
  </FullpageLoader>
{:else if socketError}
  <FullpageLoader>
    {$LL.joinStoryboardError()}
  </FullpageLoader>
{:else if isLoading}
  <FullpageLoader>
    {$LL.loadingStoryboard()}
  </FullpageLoader>
{:else if JoinPassRequired}
  <JoinCodeForm
    handleSubmit={authStoryboard}
    submitText={$LL.joinStoryboard()}
  />
{/if}

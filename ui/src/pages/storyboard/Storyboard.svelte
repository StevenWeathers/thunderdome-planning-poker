<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import Sockette from 'sockette';
  import { onDestroy, onMount } from 'svelte';

  import AddGoal from '../../components/storyboard/AddGoal.svelte';
  import ColumnForm from '../../components/storyboard/ColumnForm.svelte';
  import StoryForm from '../../components/storyboard/StoryForm.svelte';
  import ColorLegendForm from '../../components/storyboard/ColorLegendForm.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import EditStoryboard from '../../components/storyboard/EditStoryboard.svelte';
  import ExportStoryboard from '../../components/storyboard/ExportStoryboard.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import BecomeFacilitator from '../../components/BecomeFacilitator.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import {
    ChevronDown,
    Crown,
    Download,
    Link,
    LogOut,
    MessageSquareMore,
    Pencil,
    Plus,
    Settings,
    SwatchBook,
    Trash,
    User,
    Users,
  } from 'lucide-svelte';
  import JoinCodeForm from '../../components/global/JoinCodeForm.svelte';
  import FullpageLoader from '../../components/global/FullpageLoader.svelte';
  import { getWebsocketAddress } from '../../websocketUtil';
  import SubMenu from '../../components/global/SubMenu.svelte';
  import SubMenuItem from '../../components/global/SubMenuItem.svelte';
  import Personas from '../../components/storyboard/Personas.svelte';
  import GoalSection from '../../components/storyboard/GoalSection.svelte';
  import type {
    StoryboardPersona,
    StoryboardUser,
    Storyboard,
    ColorLegend,
    StoryboardGoal,
  } from '../../types/storyboard';
  import ActiveUsers from '../../components/storyboard/ActiveUsers.svelte';
  import type { NotificationService } from '../../types/notifications';
  import StoryCard from '../../components/storyboard/StoryCard.svelte';

  interface Props {
    router: any;
    storyboardId: string;
    notifications: NotificationService;
  }

  let { storyboardId, notifications, router }: Props = $props();

  const { AllowGuests } = AppConfig;
  const loginOrRegister = AllowGuests ? appRoutes.register : appRoutes.login;

  const hostname = window.location.origin;

  let isLoading = $state(true);
  let JoinPassRequired = $state(false);
  let socketError = $state(false);
  let socketReconnecting = $state(false);
  let storyboard: Storyboard = $state({
    id: '',
    name: '',
    goals: [] as StoryboardGoal[],
    users: [] as StoryboardUser[],
    color_legend: [] as ColorLegend[],
    personas: [] as StoryboardPersona[],
    facilitators: [] as string[],
    facilitatorCode: '',
    joinCode: '',
    owner_id: '',
  });
  let showUsers = $state(false);
  let showColorLegend = $state(false);
  let showColorLegendForm = $state(false);
  let showPersonas = $state(false);
  let editColumn = $state(null);
  let activeStory = $state(null);
  let showDeleteStoryboard = $state(false);
  let showEditStoryboard = $state(false);
  let showExportStoryboard = $state(false);
  let activeUserCount = $state(0);
  let columnOrderEditMode = $state(false);

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
        updateActiveUserCount();
        const joinedUser = storyboard.users.find(w => w.id === parsedEvent.userId);
        notifications.success(`${joinedUser.name} joined.`);
        break;
      case 'user_left':
        const leftUser = storyboard.users.find(w => w.id === parsedEvent.userId);
        storyboard.users = JSON.parse(parsedEvent.value);
        updateActiveUserCount();

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

  const ws = new Sockette(`${getWebsocketAddress()}/api/storyboard/${storyboardId}`, {
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

    const matchedStory = storyboard.goals[goalIndex].columns[columnIndex].stories.find(i => i.id === storyId);

    if (matchedStory) {
      const goalId = storyboard.goals[goalIndex].id;
      const columnId = storyboard.goals[goalIndex].columns[columnIndex].id;

      // determine what story to place story before in target column
      const matchedStoryIndex = storyboard.goals[goalIndex].columns[columnIndex].stories.indexOf(matchedStory);
      const sibling = storyboard.goals[goalIndex].columns[columnIndex].stories[matchedStoryIndex + 1];
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

  const addStoryColumn = (goalId: String) => {
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

  const handleAddFacilitator = userId => {
    sendSocketEvent(
      'facilitator_add',
      JSON.stringify({
        userId,
      }),
    );
  };

  const handleRemoveFacilitator = userId => {
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

  function abandonStoryboard(toggleSubmenu?: () => void) {
    return () => {
      sendSocketEvent('abandon_storyboard', '');
      toggleSubmenu?.();
    };
  }

  function toggleUsersPanel() {
    showColorLegend = false;
    showPersonas = false;
    showUsers = !showUsers;
  }

  function togglePersonas(toggleSubmenu?: () => void) {
    return () => {
      showUsers = false;
      showColorLegend = false;
      showPersonas = !showPersonas;
      toggleSubmenu?.();
    };
  }

  function toggleColumnEdit(column) {
    return () => {
      editColumn = editColumn != null ? null : column;
    };
  }

  function toggleEditLegend(toggleSubmenu?: () => void) {
    return () => {
      showColorLegend = false;
      showColorLegendForm = !showColorLegendForm;
      toggleSubmenu?.();
    };
  }

  const toggleDeleteStoryboard = (toggleSubmenu?: () => void) => {
    return () => {
      showDeleteStoryboard = !showDeleteStoryboard;
      toggleSubmenu?.();
    };
  };

  const toggleExportStoryboard = (toggleSubmenu?: () => void) => {
    return () => {
      showExportStoryboard = !showExportStoryboard;
      toggleSubmenu?.();
    };
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

  const handleGoalDeletion = (goalId: String) => {
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

  const handlePersonaRevision = (persona: StoryboardPersona) => {
    sendSocketEvent('update_persona', JSON.stringify(persona));
  };

  const handleDeletePersona = personaId => () => {
    sendSocketEvent('delete_persona', personaId);
  };

  function handleStoryboardEdit(revisedStoryboard) {
    sendSocketEvent('edit_storyboard', JSON.stringify(revisedStoryboard));
    toggleEditStoryboard()();
  }

  function toggleEditStoryboard(toggleSubmenu?: () => void) {
    return () => {
      showEditStoryboard = !showEditStoryboard;
      toggleSubmenu?.();
    };
  }

  const toggleStoryForm = story => () => {
    if (columnOrderEditMode) {
      return;
    }
    activeStory = activeStory != null ? null : story;
  };

  const toggleColumnOrderEdit = () => {
    columnOrderEditMode = !columnOrderEditMode;
  };

  let showBecomeFacilitator = $state(false);

  function becomeFacilitator(facilitatorCode) {
    sendSocketEvent('facilitator_self', facilitatorCode);
    toggleBecomeFacilitator();
  }

  function toggleBecomeFacilitator(toggleSubmenu?: () => void) {
    return () => {
      showBecomeFacilitator = !showBecomeFacilitator;
      toggleSubmenu?.();
    };
  }

  function updateActiveUserCount() {
    activeUserCount = storyboard.users.filter((u: StoryboardUser) => u.active).length;
  }

  let isFacilitator = $derived(storyboard.facilitators.length && storyboard.facilitators.includes($user.id));

  onMount(() => {
    if (!$user.id) {
      router.route(`${loginOrRegister}/storyboard/${storyboardId}`);
    }
  });
</script>

<svelte:head>
  <title>{$LL.storyboard()} {storyboard.name} | {$LL.appName()}</title>
</svelte:head>

<div class="w-full">
  <div
    class="px-6 py-2 bg-gray-100 dark:bg-gray-800 border-b border-t border-gray-400 dark:border-gray-700 flex
        flex-wrap gap-y-2"
  >
    <div class="grow">
      <h1 class="text-3xl font-bold leading-tight dark:text-gray-200">
        <span class="text-2xl text-gray-700 dark:text-gray-400">{$LL.storyboard()}</span>
        {storyboard.name}
      </h1>
    </div>
    <div class="flex justify-end space-x-2">
      <SolidButton color="green" onClick={toggleAddGoal()} testid="goal-add">
        <Plus class="inline-block w-4 h-4" />&nbsp;{$LL.storyboardAddGoal()}
      </SolidButton>
      <SubMenu label="Storyboard Settings" icon={Settings} testId="storyboard-settings">
        {#snippet children({ toggleSubmenu })}
          <SubMenuItem
            onClickHandler={togglePersonas(toggleSubmenu)}
            testId="personas-toggle"
            icon={Users}
            label={$LL.personas()}
          />
          <SubMenuItem
            onClickHandler={toggleEditLegend(toggleSubmenu)}
            testId="colorlegend"
            icon={SwatchBook}
            label={$LL.colorLegend()}
          />
          <SubMenuItem
            onClickHandler={toggleExportStoryboard(toggleSubmenu)}
            testId="storyboard-export"
            icon={Download}
            label={$LL.export()}
          />
          {#if isFacilitator}
            <SubMenuItem
              onClickHandler={toggleEditStoryboard(toggleSubmenu)}
              testId="storyboard-edit"
              icon={Pencil}
              label={$LL.editStoryboard()}
            />
            <SubMenuItem
              onClickHandler={toggleDeleteStoryboard(toggleSubmenu)}
              testId="storyboard-delete"
              icon={Trash}
              label={$LL.deleteStoryboard()}
            />
          {:else}
            <SubMenuItem
              onClickHandler={toggleBecomeFacilitator(toggleSubmenu)}
              testId="become-facilitator"
              icon={Crown}
              label={$LL.becomeFacilitator()}
            />
            <SubMenuItem
              onClickHandler={abandonStoryboard(toggleSubmenu)}
              testId="storyboard-leave"
              icon={LogOut}
              label={$LL.leaveStoryboard()}
            />
          {/if}
        {/snippet}
      </SubMenu>
      <SolidButton color="gray" onClick={toggleUsersPanel} testid="users-toggle">
        <Users class="inline-block w-4 h-4 me-2" />
        {$LL.users()}&nbsp;<span
          class="rounded-full bg-gray-500 text-white dark:bg-gray-300 dark:text-gray-800 text-sm {activeUserCount > 9
            ? 'px-1'
            : 'px-2'}">{activeUserCount}</span
        >
        <ChevronDown class="ms-1 inline-block w-4 h-4" />
      </SolidButton>
    </div>
  </div>
  {#if showUsers}
    <ActiveUsers
      users={storyboard.users}
      facilitatorIds={storyboard.facilitators}
      {isFacilitator}
      inviteUrl="{hostname}{appRoutes.storyboard}/{storyboard.id}"
      onAddFacilitator={handleAddFacilitator}
      onRemoveFacilitator={handleRemoveFacilitator}
    />
  {/if}
  {#each storyboard.goals as goal, goalIndex (goal.id)}
    <GoalSection
      {goal}
      handleDelete={handleGoalDeletion}
      handleColumnAdd={addStoryColumn}
      toggleEdit={toggleAddGoal}
      {goalIndex}
      {isFacilitator}
      {columnOrderEditMode}
      {toggleColumnOrderEdit}
    >
      <div class="flex">
        {#each goal.columns as goalColumn, columnIndex (goalColumn.id)}
          <div class="flex-none mx-2 w-40" data-testid="goal-personas">
            <div class="w-full mb-2">
              {#each goalColumn.personas as persona}
                <div class="mt-4 dark:text-gray-300 text-right" data-testid="goal-persona">
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
                    title={goalColumn.name}
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
                    class:cursor-not-allowed={columnOrderEditMode}
                    disabled={columnOrderEditMode}
                    title={$LL.storyboardEditColumn()}
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
                    class:cursor-not-allowed={columnOrderEditMode}
                    disabled={columnOrderEditMode}
                    title={$LL.storyboardAddStoryToColumn()}
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
              data-goalid={goal.id}
              data-columnid={goalColumn.id}
              data-goalIndex={goalIndex}
              data-columnindex={columnIndex}
              use:dndzone={{
                items: goalColumn.stories,
                type: 'story',
                dropTargetStyle: '',
                dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
                dragDisabled: columnOrderEditMode,
              }}
              onconsider={handleDndConsider}
              onfinalize={handleDndFinalize}
            >
              {#each goalColumn.stories as story (story.id)}
                <StoryCard {story} {goalColumn} {goal} {columnOrderEditMode} {toggleStoryForm} />
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </GoalSection>
  {/each}
</div>

{#if showAddGoal}
  <AddGoal
    {handleGoalAdd}
    toggleAddGoal={toggleAddGoal()}
    {handleGoalRevision}
    goalId={reviseGoalId}
    goalName={reviseGoalName}
  />
{/if}

{#if editColumn}
  <ColumnForm
    {handleColumnRevision}
    toggleColumnEdit={toggleColumnEdit()}
    column={editColumn}
    personas={storyboard.personas}
    handlePersonaAdd={handleColumnPersonaAdd}
    handlePersonaRemove={handleColumnPersonaRemove}
    {deleteColumn}
  />
{/if}

{#if activeStory}
  <StoryForm
    toggleStoryForm={toggleStoryForm()}
    story={activeStory}
    {sendSocketEvent}
    {notifications}
    colorLegend={storyboard.color_legend}
    users={storyboard.users}
  />
{/if}

{#if showColorLegendForm}
  <ColorLegendForm
    {handleLegendRevision}
    toggleEditLegend={toggleEditLegend()}
    colorLegend={storyboard.color_legend}
    {isFacilitator}
  />
{/if}

{#if showEditStoryboard}
  <EditStoryboard
    storyboardName={storyboard.name}
    {handleStoryboardEdit}
    toggleEditStoryboard={toggleEditStoryboard()}
    joinCode={storyboard.joinCode}
    facilitatorCode={storyboard.facilitatorCode}
  />
{/if}

{#if showExportStoryboard}
  <ExportStoryboard {storyboard} closeModal={toggleExportStoryboard()} />
{/if}

{#if showDeleteStoryboard}
  <DeleteConfirmation
    toggleDelete={toggleDeleteStoryboard()}
    handleDelete={concedeStoryboard}
    confirmText={'Are you sure you want to delete this Storyboard?'}
    confirmBtnText={'Delete Storyboard'}
  />
{/if}

{#if showBecomeFacilitator}
  <BecomeFacilitator handleBecomeFacilitator={becomeFacilitator} toggleBecomeFacilitator={toggleBecomeFacilitator()} />
{/if}

{#if showPersonas}
  <Personas
    personas={storyboard.personas}
    closeModal={togglePersonas()}
    onDelete={handleDeletePersona}
    onAdd={handlePersonaAdd}
    onUpdate={handlePersonaRevision}
    {isFacilitator}
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
  <JoinCodeForm handleSubmit={authStoryboard} submitText={$LL.joinStoryboard()} />
{/if}

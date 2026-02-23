<script lang="ts">
  import Sockette from 'sockette';
  import { onDestroy, onMount } from 'svelte';

  import AddGoal from '../../components/storyboard/AddGoal.svelte';
  import ColumnForm from '../../components/storyboard/ColumnForm.svelte';
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
    GoalIcon,
    KanbanIcon,
    LayoutDashboardIcon,
    LogOut,
    MinusIcon,
    Pencil,
    Plus,
    PlusIcon,
    Settings,
    SwatchBook,
    Trash,
    Users,
  } from '@lucide/svelte';
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
    StoryboardColumn,
    StoryboardStory,
  } from '../../types/storyboard';
  import ActiveUsers from '../../components/storyboard/ActiveUsers.svelte';
  import type { NotificationService } from '../../types/notifications';
  import GoalColumns from '../../components/storyboard/GoalColumns.svelte';

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
  let editColumn: StoryboardColumn | null = $state(null);
  let showDeleteStoryboard = $state(false);
  let showEditStoryboard = $state(false);
  let showExportStoryboard = $state(false);
  let activeUserCount = $state(0);
  let columnOrderEditMode = $state(false);
  let scale = $state(1);

  const ZOOM_MIN = 1;
  const ZOOM_MAX = 1.5;

  let ws: any;

  const onSocketMessage = function (evt: MessageEvent) {
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
        if (joinedUser) {
          notifications.success(`${joinedUser.name} joined.`);
        }
        break;
      case 'user_left':
        const leftUser = storyboard.users.find(w => w.id === parsedEvent.userId);
        storyboard.users = JSON.parse(parsedEvent.value);
        updateActiveUserCount();
        if (leftUser) {
          notifications.danger(`${leftUser.name} left.`);
        }
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
              if (editColumn && column.id === editColumn.id) {
                editColumn = column;
              }
            });
          });
        }
        break;
      case 'column_moved':
        const updatedGoal = JSON.parse(parsedEvent.value);
        storyboard.goals = storyboard.goals.map(goal => (goal.id === updatedGoal.id ? updatedGoal : goal));
        break;
      case 'story_added':
        storyboard.goals = JSON.parse(parsedEvent.value);
        break;
      case 'story_updated':
        storyboard.goals = JSON.parse(parsedEvent.value);

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

  const sendSocketEvent = (type: string, value: any) => {
    ws.send(
      JSON.stringify({
        type,
        value,
      }),
    );
  };

  function authStoryboard(joinPasscode: string) {
    sendSocketEvent('auth_storyboard', joinPasscode);
  }

  const addStory = (goalId: string, columnId: string) => () => {
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

  const handleAddFacilitator = (userId: string) => {
    sendSocketEvent(
      'facilitator_add',
      JSON.stringify({
        userId,
      }),
    );
  };

  const handleRemoveFacilitator = (userId: string) => {
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

  function toggleColumnEdit(column: StoryboardColumn | null = null) {
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

  const toggleAddGoal = (goalId?: string) => () => {
    if (goalId) {
      const goal = storyboard.goals.find(p => p.id === goalId);
      if (goal) {
        reviseGoalId = goalId;
        reviseGoalName = goal.name;
      }
    } else {
      reviseGoalId = '';
      reviseGoalName = '';
    }
    showAddGoal = !showAddGoal;
  };

  const handleGoalAdd = (goalName: string) => {
    sendSocketEvent('add_goal', goalName);
  };

  const handleGoalRevision = (updatedGoal: any) => {
    sendSocketEvent('revise_goal', JSON.stringify(updatedGoal));
  };

  const handleGoalDeletion = (goalId: String) => {
    sendSocketEvent('delete_goal', goalId);
  };

  const handleColumnRevision = (column: StoryboardColumn) => {
    sendSocketEvent('revise_column', JSON.stringify(column));
  };

  const handleLegendRevision = (legend: ColorLegend[]) => {
    sendSocketEvent('revise_color_legend', JSON.stringify(legend));
  };

  const handlePersonaAdd = (persona: Omit<StoryboardPersona, 'id'>) => {
    sendSocketEvent('add_persona', JSON.stringify(persona));
  };

  const handlePersonaRevision = (persona: StoryboardPersona) => {
    sendSocketEvent('update_persona', JSON.stringify(persona));
  };

  const handleDeletePersona = (personaId: string) => () => {
    sendSocketEvent('delete_persona', personaId);
  };

  function handleStoryboardEdit(revisedStoryboard: any) {
    sendSocketEvent('edit_storyboard', JSON.stringify(revisedStoryboard));
    toggleEditStoryboard()();
  }

  function toggleEditStoryboard(toggleSubmenu?: () => void) {
    return () => {
      showEditStoryboard = !showEditStoryboard;
      toggleSubmenu?.();
    };
  }

  const toggleColumnOrderEdit = () => {
    columnOrderEditMode = !columnOrderEditMode;
  };

  let showBecomeFacilitator = $state(false);

  function becomeFacilitator(facilitatorCode: string) {
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

  function increaseScale() {
    scale = Math.min(scale + 0.1, ZOOM_MAX);
  }
  function decreaseScale() {
    scale = Math.max(scale - 0.1, ZOOM_MIN);
  }

  let isFacilitator = $derived(storyboard.facilitators.length > 0 && storyboard.facilitators.includes($user.id));

  onMount(() => {
    if (!$user.id) {
      router.route(`${loginOrRegister}/storyboard/${storyboardId}`);
      return;
    }

    ws = new Sockette(`${getWebsocketAddress()}/api/storyboard/${storyboardId}`, {
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
  });

  onDestroy(() => {
    if (ws) {
      ws.close();
    }
  });
</script>

<svelte:head>
  <title>{$LL.storyboard()} {storyboard.name} | {$LL.appName()}</title>
</svelte:head>

<div class="w-full">
  <div
    class="px-6 py-2 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-gray-800 dark:to-gray-900 border-b border-t border-blue-200 dark:border-gray-700 flex
        flex-wrap gap-y-2 shadow-sm"
  >
    <div class="grow flex items-center gap-2.5">
      <div
        class="flex items-center justify-center w-9 h-9 bg-blue-600 dark:bg-blue-500 rounded-lg shadow-md"
        title={$LL.storyboard()}
      >
        <LayoutDashboardIcon class="w-6 h-6 text-white" />
      </div>
      <h1 class="text-3xl font-bold leading-tight text-gray-900 dark:text-gray-100">
        {storyboard.name}
      </h1>
    </div>
    <div class="flex justify-end space-x-2">
      {#if !columnOrderEditMode}
        <div
          class="flex items-center gap-1 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg p-1 shadow-sm"
        >
          <button
            onclick={decreaseScale}
            disabled={scale <= ZOOM_MIN}
            class="p-2 rounded-md transition-colors hover:bg-gray-100 dark:hover:bg-gray-600 disabled:cursor-not-allowed disabled:hover:bg-transparent dark:disabled:hover:bg-transparent text-gray-800 dark:text-gray-100 disabled:text-gray-500 dark:disabled:text-gray-400"
            title="Zoom Storyboard out"
          >
            <MinusIcon class="w-4 h-4" />
          </button>
          <button
            onclick={increaseScale}
            disabled={scale >= ZOOM_MAX}
            class="p-2 rounded-md transition-colors hover:bg-gray-100 dark:hover:bg-gray-600 disabled:cursor-not-allowed disabled:hover:bg-transparent dark:disabled:hover:bg-transparent text-gray-800 dark:text-gray-100 disabled:text-gray-500 dark:disabled:text-gray-400"
            title="Zoom Storyboard in"
          >
            <PlusIcon class="w-4 h-4" />
          </button>
          <div class="w-px h-6 bg-gray-300 dark:bg-gray-600"></div>
          <span class="sr-only">Zoom Board</span>
          <div class="px-1 flex items-center text-gray-600 dark:text-gray-400">
            <LayoutDashboardIcon class="w-6 h-6" />
          </div>
        </div>
        <div class="w-px h-8 bg-gray-300 dark:bg-gray-600 self-center"></div>
        <SolidButton color="green" onClick={() => toggleAddGoal(undefined)()} testid="goal-add">
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
      {/if}
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
  {#if storyboard.goals && storyboard.goals.length === 0}
    <div
      class="m-6 p-6 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center text-gray-500 dark:text-gray-400"
    >
      <div class="mx-auto h-24 w-24 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-6">
        <GoalIcon class="h-12 w-12 text-slate-400" />
      </div>
      <p class="text-xl text-slate-600 dark:text-slate-400 mb-2">No goals created yet.</p>
      <p class="text-sm text-slate-500 dark:text-slate-500 max-w-2xl mx-auto">
        Story mapping goals are the big outcomes or user journeys you want to achieve. Each goal becomes a section where
        you organize stories into columns to map how work supports that outcome.
      </p>
      <div class="mt-6 flex justify-center">
        <SolidButton color="green" onClick={() => toggleAddGoal(undefined)()} testid="goal-add-empty">
          <Plus class="inline-block w-4 h-4" />&nbsp;{$LL.storyboardAddGoal()}
        </SolidButton>
      </div>
    </div>
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
      {#if goal.columns.length === 0}
        <div
          class="m-6 p-6 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg text-center text-gray-500 dark:text-gray-400"
        >
          <div
            class="mx-auto h-24 w-24 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mb-6"
          >
            <KanbanIcon class="h-12 w-12 text-slate-400" />
          </div>
          <p class="text-xl text-slate-600 dark:text-slate-400 mb-2">No columns created yet.</p>
          <p class="text-sm text-slate-500 dark:text-slate-500 max-w-2xl mx-auto">
            Columns help you organize the stories for this goal. You can create columns for different stages of work,
            different types of work, or however else you'd like to group your stories.
          </p>
          <div class="mt-6 flex justify-center">
            <SolidButton color="green" onClick={() => addStoryColumn(goal.id)} testid="goal-col-add-empty">
              <Plus class="inline-block w-4 h-4" />&nbsp;{$LL.storyboardAddColumn()}
            </SolidButton>
          </div>
        </div>
      {:else}
        <GoalColumns
          bind:goals={storyboard.goals}
          {goal}
          {goalIndex}
          {columnOrderEditMode}
          {addStory}
          {toggleColumnEdit}
          {sendSocketEvent}
          {notifications}
          colorLegend={storyboard.color_legend}
          users={storyboard.users}
          {scale}
          personas={storyboard.personas}
        />
      {/if}
    </GoalSection>
  {/each}
</div>

{#if showAddGoal}
  <AddGoal
    {handleGoalAdd}
    toggleAddGoal={toggleAddGoal(undefined)}
    {handleGoalRevision}
    goalId={reviseGoalId}
    goalName={reviseGoalName}
  />
{/if}

{#if editColumn}
  <ColumnForm {handleColumnRevision} toggleColumnEdit={toggleColumnEdit(null)} column={editColumn} />
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

<script lang="ts">
  import { ChevronDown, ChevronUp, GoalIcon, Pencil, Plus, SendToBack, Settings, Trash } from '@lucide/svelte';
  import GoalEstimate from '../../components/storyboard/GoalEstimate.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import SubMenu from '../../components/global/SubMenu.svelte';
  import SubMenuItem from '../../components/global/SubMenuItem.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    children?: import('svelte').Snippet;
    goal: any;
    goalIndex: number;
    isFacilitator: boolean;
    handleDelete: any;
    handleColumnAdd: any;
    toggleEdit: any;
    columnOrderEditMode: any;
    toggleColumnOrderEdit: any;
  }

  let {
    children,
    toggleEdit = (goalId: String) => () => {},
    handleDelete = (goalId: String) => {},
    handleColumnAdd = (goalId: String) => {},
    goal = { id: '', columns: [] },
    goalIndex = 0,
    isFacilitator = false,
    columnOrderEditMode = false,
    toggleColumnOrderEdit = () => {},
  }: Props = $props();

  let collapsed = $state(false);
  let showDeleteConfirmation = $state(false);

  const toggleCollapse = () => {
    collapsed = !collapsed;
  };

  const toggleDeleteConfirmation = (toggleSubmenu?: () => void) => () => {
    showDeleteConfirmation = !showDeleteConfirmation;
    toggleSubmenu?.();
  };

  const handleDeletion = () => {
    handleDelete(goal.id);
  };

  const handleToggleEdit = (toggleSubmenu?: () => void) => () => {
    toggleEdit(goal.id)();
    toggleSubmenu?.();
  };

  const handleColAdd = () => {
    handleColumnAdd(goal.id);
  };

  const handleToggleColumnOrderEdit = (toggleSubmenu?: () => void) => () => {
    toggleColumnOrderEdit();
    toggleSubmenu?.();
  };
</script>

<div data-goalid={goal.id} data-testid="storyboard-goal">
  <div
    class="flex flex-wrap gap-y-2 px-6 py-2 bg-gradient-to-r from-red-50 to-orange-50 dark:from-gray-800 dark:to-gray-900 border-b-2 border-red-200 dark:border-gray-700 shadow-sm {goalIndex >
    0
      ? 'border-t-2'
      : ''}"
  >
    <div class="grow flex items-center gap-2.5">
      <button
        onclick={toggleCollapse}
        data-testid="goal-expand"
        data-collapsed={collapsed}
        class="group text-gray-700 dark:text-gray-300"
      >
        {#if collapsed}
          <ChevronUp class="w-6 h-6 transition-transform duration-200 group-hover:rotate-180" />
        {:else}
          <ChevronDown class="w-6 h-6 transition-transform duration-200 group-hover:rotate-180" />
        {/if}
      </button>
      <div
        class="flex items-center justify-center w-9 h-9 bg-red-600 dark:bg-red-500 rounded-full shadow-md"
        title="Goal"
      >
        <GoalIcon class="w-6 h-6 text-white" />
      </div>
      <h2 class="inline-block font-bold text-2xl leading-tight text-gray-900 dark:text-gray-100">
        <button onclick={toggleCollapse} class="inline-flex items-center gap-3 group">
          <span>{goal.name}</span>
          <GoalEstimate columns={goal.columns} />
        </button>
      </h2>
    </div>
    <div class="flex justify-end space-x-2">
      {#if isFacilitator}
        {#if columnOrderEditMode}
          <SolidButton color="yellow" onClick={toggleColumnOrderEdit} testid="column-order-done">
            <SendToBack class="inline-block w-4 h-4 me-1" />End Column Order Edit Mode
          </SolidButton>
        {:else}
          <SolidButton color="green" onClick={handleColAdd} testid="column-add">
            <Plus class="inline-block w-4 h-4" />&nbsp;{$LL.storyboardAddColumn()}
          </SolidButton>

          <SubMenu label="Goal Settings" icon={Settings} testId="goal-settings">
            {#snippet children({ toggleSubmenu })}
              <SubMenuItem
                onClickHandler={handleToggleEdit(toggleSubmenu)}
                testId="goal-edit"
                icon={Pencil}
                label={$LL.edit()}
              />
              <SubMenuItem
                onClickHandler={handleToggleColumnOrderEdit(toggleSubmenu)}
                testId="toggle-edit-goal-column-order"
                icon={SendToBack}
                label={`Edit Columns Order`}
              />
              <SubMenuItem
                onClickHandler={toggleDeleteConfirmation(toggleSubmenu)}
                testId="goal-delete"
                icon={Trash}
                label={$LL.delete()}
              />
            {/snippet}
          </SubMenu>
        {/if}
      {/if}
    </div>
  </div>
  {#if !collapsed}
    <section class="px-2" style="overflow-x: scroll">
      {@render children?.()}
    </section>
  {/if}

  {#if showDeleteConfirmation}
    <DeleteConfirmation
      toggleDelete={toggleDeleteConfirmation()}
      handleDelete={handleDeletion}
      confirmText={'Are you sure you want to delete this goal?'}
      confirmBtnText={'Delete Goal'}
    />
  {/if}
</div>

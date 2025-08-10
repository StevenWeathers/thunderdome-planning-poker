<script lang="ts">
    import {
    ChevronDown,
    ChevronUp,
    Pencil,
    Plus,
    Settings,
    Trash,
  } from 'lucide-svelte';
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
  }

  let { children, toggleEdit = (goalId: String) => () => {}, handleDelete = (goalId: String) => {}, handleColumnAdd = (goalId: String) => {}, goal = { id: "", columns: [] }, goalIndex = 0, isFacilitator = false  }: Props = $props();

  let collapsed = $state(false);
  let showSettings = $state(false);
  let showDeleteConfirmation = $state(false);

  const toggleCollapse = () => {
    collapsed = !collapsed;
  };

  const toggleSettings = () => {
    showSettings = !showSettings;
  };

  const toggleDeleteConfirmation = () => {
    showDeleteConfirmation = !showDeleteConfirmation;
  };

  const handleDeletion = () => {
    handleDelete(goal.id);
  };

  const handleToggleEdit = () => {
    toggleEdit(goal.id)();
    toggleSettings();
  };

  const handleColAdd = () => {
    handleColumnAdd(goal.id);
  };
</script>

<div data-goalid={goal.id} data-testid="storyboard-goal">
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
            onclick={toggleCollapse}
            data-testid="goal-expand"
            data-collapsed={collapsed}
          >
            {#if collapsed}
              <ChevronDown class="me-1 inline-block" />
            {:else}
              <ChevronUp class="me-1 inline-block" />
            {/if}
          </button>{goal.name}&nbsp;<GoalEstimate columns={goal.columns} />
        </h2>
      </div>
    </div>
    <div class="w-1/4 flex justify-end space-x-2">
      {#if isFacilitator}
        <SolidButton
          color="green"
          onClick={handleColAdd}
          testid="column-add"
        >
          <Plus class="inline-block w-4 h-4" />&nbsp;{$LL.storyboardAddColumn()}
        </SolidButton>
        <div class="inline-block relative">
          <SolidButton onClick={toggleSettings} color="blue">
            <Settings class="inline-block h-3.5 w-3.5 me-1.5 -ms-1" />Settings
          </SolidButton>
          {#if showSettings}
            <SubMenu>
              <SubMenuItem
                onClickHandler={handleToggleEdit}
                testId="goal-edit"
              >
                <Pencil class="h-5 w-5 me-2 inline-block" />{$LL.edit()}
              </SubMenuItem>
              <SubMenuItem
                onClickHandler={toggleDeleteConfirmation}
                testId="goal-delete"
              >
                <Trash
                  class="h-5 w-5 me-2 inline-block"
                />{$LL.delete()}
              </SubMenuItem>
            </SubMenu>
          {/if}
        </div>
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
        toggleDelete={toggleDeleteConfirmation}
        handleDelete={handleDeletion}
        confirmText={"Are you sure you want to delete this goal?"}
        confirmBtnText={"Delete Goal"}
        />
    {/if}
</div>

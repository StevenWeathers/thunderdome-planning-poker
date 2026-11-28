<script lang="ts">
  import { dndzone } from 'svelte-dnd-action';
  import { GoalIcon, Grip } from '@lucide/svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { StoryboardGoal } from '../../types/storyboard';

  interface Props {
    goals: StoryboardGoal[];
    sendSocketEvent: (event: string, data: string) => void;
  }

  let { goals = $bindable([]), sendSocketEvent }: Props = $props();

  function handleDndConsider(e: CustomEvent) {
    goals = e.detail.items;
  }

  function handleDndFinalize(e: CustomEvent) {
    const items: StoryboardGoal[] = e.detail.items;
    const goalId: string = e.detail.info.id;
    goals = items;

    const goalIndex = items.findIndex(goal => goal.id === goalId);

    if (goalIndex !== -1) {
      sendSocketEvent('move_goal', JSON.stringify({ goalId, placeBefore: items[goalIndex + 1]?.id ?? '' }));
    }
  }
</script>

<div class="p-4">
  <p class="mb-3 text-gray-600 dark:text-gray-400">{$LL.storyboardDragToReorderGoals()}</p>
  <!-- only goal names are listed, as whole goal sections are too tall to drag -->
  <div
    class="flex flex-col gap-2"
    aria-label={$LL.storyboardEditGoalsOrder()}
    data-testid="goal-order-dropzone"
    use:dndzone={{
      items: goals,
      type: 'goal' as const,
      dropTargetStyle: {},
      dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
    }}
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
  >
    {#each goals as goal (goal.id)}
      <div
        class="flex items-center gap-2.5 px-4 py-2 rounded-lg border-2 border-red-200 dark:border-gray-700 bg-red-50 dark:bg-gray-800 shadow-sm cursor-move"
        aria-label={goal.name}
        data-testid="goal-order-item"
      >
        <Grip class="w-5 h-5 shrink-0 text-gray-500 dark:text-gray-400" />
        <GoalIcon class="w-5 h-5 shrink-0 text-red-600 dark:text-red-400" />
        <span class="font-bold text-xl min-w-0 break-words text-gray-900 dark:text-gray-100">{goal.name}</span>
      </div>
    {/each}
  </div>
</div>

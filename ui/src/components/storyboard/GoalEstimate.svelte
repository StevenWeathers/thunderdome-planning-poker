<script lang="ts">
  import type { StoryboardColumn } from '../../types/storyboard';

  interface Props {
    columns?: StoryboardColumn[];
  }

  let { columns = [] }: Props = $props();

  type Story = {
    points: number;
  };
  type Column = {
    stories: Array<Story>;
    sort_order?: string;
  };

  function calculateGoalEstimate(goalColumns: Array<Column>) {
    let estimate: number = 0;
    for (let column of goalColumns) {
      for (let story of column.stories) {
        estimate += story.points;
      }
    }

    return estimate;
  }

  let totalPoints = $derived(calculateGoalEstimate(columns));
</script>

<span
  class="inline-flex items-center gap-1.5 px-3 py-1 bg-gradient-to-r from-purple-100 to-indigo-100 dark:from-purple-900/40 dark:to-indigo-900/40 text-purple-800 dark:text-purple-200 rounded-full text-sm font-semibold shadow-sm border border-purple-200 dark:border-purple-700"
  title="Estimated Total Story Points"
>
  <span>{totalPoints}</span>
  <span class="text-xs font-normal opacity-75">pts</span>
</span>

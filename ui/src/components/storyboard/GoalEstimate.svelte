<script lang="ts">
  interface Props {
    columns?: any;
  }

  let { columns = [] }: Props = $props();

  type story = {
    points: number;
  };
  type column = {
    stories: story;
  };

  function calculateGoalEstimate(goalColumns: Array<column>) {
    let estimate: number = 0;
    for (let column of goalColumns) {
      for (let story of column.stories) {
        estimate += story.points;
      }
    }

    return estimate;
  }
</script>

<span>({calculateGoalEstimate(columns)} Story Points)</span>

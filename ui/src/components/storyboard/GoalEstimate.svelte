<script lang="ts">
  interface Props {
    columns?: any;
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
</script>

<span>({calculateGoalEstimate(columns)} Story Points)</span>

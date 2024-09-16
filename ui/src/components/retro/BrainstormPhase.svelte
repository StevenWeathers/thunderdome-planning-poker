<script lang="ts">
  import ItemForm from './ItemForm.svelte';

  export let phase: string = '';
  export let sendSocketEvent: (event: string, data: any) => void;
  export let isFacilitator: boolean = false;
  export let items: any = [];
  export let template: any = {
    format: {
      columns: [],
    },
  };
  export let users: any = [];
  export let brainstormVisibility: boolean = false;
  export let columnColors: any = {};

  $: numCols = template.format.columns.length;
</script>

<div
  class="w-full grid gap-4"
  class:grid-cols-5="{numCols === 5}"
  class:grid-cols-4="{numCols === 4}"
  class:grid-cols-3="{numCols === 3}"
>
  {#each template.format.columns as column}
    <ItemForm
      sendSocketEvent="{sendSocketEvent}"
      itemType="{column.name}"
      newItemPlaceholder="{column.label}..."
      phase="{phase}"
      isFacilitator="{isFacilitator}"
      items="{items}"
      users="{users}"
      feedbackVisibility="{brainstormVisibility}"
      color="{column.color}"
      icon="{column.icon}"
      columnColors="{columnColors}"
    />
  {/each}
</div>

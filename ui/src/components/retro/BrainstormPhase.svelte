<script lang="ts">
  import ItemForm from './ItemForm.svelte';

  interface Props {
    phase?: string;
    sendSocketEvent: (event: string, data: any) => void;
    isFacilitator?: boolean;
    items?: any;
    template?: any;
    users?: any;
    brainstormVisibility?: boolean;
    columnColors?: any;
  }

  let {
    phase = '',
    sendSocketEvent,
    isFacilitator = false,
    items = [],
    template = {
      format: {
        columns: [],
      },
    },
    users = [],
    brainstormVisibility = false,
    columnColors = {},
  }: Props = $props();

  let numCols = $derived(template.format.columns.length);
</script>

<div
  class="w-full grid gap-4"
  class:grid-cols-5={numCols === 5}
  class:grid-cols-4={numCols === 4}
  class:grid-cols-3={numCols === 3}
>
  {#each template.format.columns as column}
    <ItemForm
      {sendSocketEvent}
      itemType={column.name}
      newItemPlaceholder="{column.label}..."
      {phase}
      {isFacilitator}
      {items}
      {users}
      feedbackVisibility={brainstormVisibility}
      color={column.color}
      icon={column.icon}
      {columnColors}
    />
  {/each}
</div>

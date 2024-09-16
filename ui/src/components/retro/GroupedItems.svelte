<script lang="ts">
  import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

  export let groups = [];
  export let columns = [];
  export let phase: string = '';
  export let isFacilitator = false;
  export let users = [];
  export let sendSocketEvent = (event: string, value: any) => {};
  export let columnColors: any = {};
</script>

{#each groups as group, i (group.id)}
  {#if group.items.length > 0}
    <RetroFeedbackGroup
      phase="{phase}"
      group="{group}"
      users="{users}"
      isFacilitator="{isFacilitator}"
      sendSocketEvent="{sendSocketEvent}"
      columnColors="{columnColors}"
    />
  {/if}
{/each}

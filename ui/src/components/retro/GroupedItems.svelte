<script lang="ts">
  import RetroFeedbackGroup from './RetroFeedbackGroup.svelte';

  interface Props {
    groups?: any;
    columns?: any;
    phase?: string;
    isFacilitator?: boolean;
    users?: any;
    sendSocketEvent?: any;
    columnColors?: any;
  }

  let {
    groups = [],
    columns = [],
    phase = '',
    isFacilitator = false,
    users = [],
    sendSocketEvent = (event: string, value: any) => {},
    columnColors = {}
  }: Props = $props();
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

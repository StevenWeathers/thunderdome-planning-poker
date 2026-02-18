<script lang="ts">
  import { Check, X } from '@lucide/svelte';

  interface Props {
    boolValue?: boolean;
  }

  let { boolValue = false }: Props = $props();
</script>

{#if boolValue}
  <span class="sr-only">True</span><Check class="text-green-500 dark:text-lime-400" />
{:else}
  <span class="sr-only">False</span><X class="text-red-500" />
{/if}

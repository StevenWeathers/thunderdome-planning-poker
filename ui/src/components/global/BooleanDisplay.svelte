<script lang="ts">
  import { Check, X } from 'lucide-svelte';

  interface Props {
    boolValue?: boolean;
  }

  let { boolValue = false }: Props = $props();
</script>

{#if boolValue}
  <Check class="text-green-500 dark:text-lime-400" />
{:else}
  <X class="text-red-500" />
{/if}

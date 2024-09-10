<script lang="ts">
  import { Check, X } from 'lucide-svelte';

  export let boolValue: boolean = false;
</script>

{#if boolValue}
  <Check class="text-green-500 dark:text-lime-400" />
{:else}
  <X class="text-red-500" />
{/if}

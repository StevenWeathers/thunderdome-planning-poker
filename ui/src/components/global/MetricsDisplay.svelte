<script lang="ts">
  import type { ComponentType } from 'svelte';
  import { HelpCircle } from 'lucide-svelte';

  interface Props {
    metrics: Array<{
    key: string;
    name: string;
    value: number | string;
    icon: ComponentType;
  }>;
  }

  let { metrics }: Props = $props();

  const getMetricValue = (value: number | string) => {
    if (typeof value === 'number') {
      return value || 0; // Return 0 if value is falsy (0, null, undefined)
    }
    return value || ''; // Return empty string for non-numeric falsy values
  };
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {#each metrics as { key, name, value, icon }}
    {@const SvelteComponent = icon || HelpCircle}
    <div class="bg-gray-100 dark:bg-gray-700 rounded-lg p-4 flex items-center">
      <div class="me-4">
        <SvelteComponent
          class="w-8 h-8 text-blue-500 dark:text-blue-400"
        />
      </div>
      <div>
        <p class="text-sm font-medium text-gray-500 dark:text-gray-400">
          {name}
        </p>
        <p class="text-lg font-semibold text-gray-800 dark:text-white">
          {getMetricValue(value)}
        </p>
      </div>
    </div>
  {/each}
</div>

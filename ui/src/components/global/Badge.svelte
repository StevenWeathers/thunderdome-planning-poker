<script lang="ts">
  interface Props {
    label: string;
    color?: 'gray' | 'blue';
  }

  let { label, color = 'gray' }: Props = $props();
</script>

{#if color === 'blue'}
  <span class="inline-flex items-center px-3 py-1 rounded-full font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
      {label}
  </span>
{:else}
  <span class="inline-flex items-center px-3 py-1 rounded-full font-medium bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200">
    {label}
  </span>
{/if}
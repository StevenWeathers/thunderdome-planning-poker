<script lang="ts">
  interface Props {
    label: string;
    color?: 'gray' | 'blue' | 'indigo';
    testId?: string;
    class?: string;
  }

  let { class: klass, label, testId = "", color = 'gray' }: Props = $props();

  const badgeClasses = 'inline-flex items-center px-3 py-1 rounded-full font-medium';
</script>

{#if color === 'blue'}
  <span class="{badgeClasses} bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 {klass}" data-testid={testId}>
      {label}
  </span>
{:else if color === 'indigo'}
  <span class="{badgeClasses} bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200 {klass}" data-testid={testId}>
      {label}
  </span>
{:else}
  <span class="{badgeClasses} bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 {klass}" data-testid={testId}>
    {label}
  </span>
{/if}
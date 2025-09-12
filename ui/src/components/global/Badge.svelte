<script lang="ts">
  interface Props {
    label: string;
    color?: 'gray' | 'blue' | 'indigo' | 'green' | 'red' | 'orange';
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
{:else if color === 'green'}
  <span class="{badgeClasses} bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200 {klass}" data-testid={testId}>
      {label}
  </span>
{:else if color === 'red'}
  <span class="{badgeClasses} bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200 {klass}" data-testid={testId}>
      {label}
  </span>
{:else if color === 'orange'}
  <span class="{badgeClasses} bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200 {klass}" data-testid={testId}>
      {label}
  </span>
{:else}
  <span class="{badgeClasses} bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200 {klass}" data-testid={testId}>
    {label}
  </span>
{/if}
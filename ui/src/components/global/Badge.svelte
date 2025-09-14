<script lang="ts">
  interface Props {
    label: string;
    color?: 'gray' | 'blue' | 'indigo' | 'red' | 'green';
    testId?: string;
    class?: string;
    title?: string;
  }

  let { class: klass, label, testId = "", color = 'gray', title = label }: Props = $props();

  const badgeClasses = 'inline-flex items-center px-3 py-1 rounded-full font-medium';

  function getColorClasses(color: string) {
    switch (color) {
      case 'blue':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'indigo':
        return 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200';
      case 'red':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      case 'green':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
    }
  }
</script>

<span class="{badgeClasses} {getColorClasses(color)} {klass}" data-testid={testId} title={title}>
  {label}
</span>
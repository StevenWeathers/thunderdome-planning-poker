<script lang="ts">
  interface Props {
    disabled?: boolean;
    color?: string;
    additionalClasses?: string;
    type?: string;
    testid?: string;
    onClick?: any;
    options?: any;
    href?: string;
    children?: import('svelte').Snippet;
  }

  let {
    disabled = false,
    color = 'green',
    additionalClasses = '',
    type = 'button',
    testid = '',
    onClick = () => {},
    options = {},
    href = '',
    children,
  }: Props = $props();
</script>

{#if href != ''}
  <a
    {href}
    class="btn btn-{color}
       {disabled ? 'disabled' : ''}
       {additionalClasses}"
    data-testid={testid}
    {...options}
  >
    {@render children?.()}
  </a>
{:else}
  <button
    class="btn btn-{color}
        {disabled ? 'disabled' : ''}
        {additionalClasses}"
    onclick={onClick}
    {type}
    {disabled}
    data-testid={testid}
    {...options}
  >
    {@render children?.()}
  </button>
{/if}

<style lang="postcss">
  .btn {
    @apply text-white;
    @apply inline-flex;
    @apply items-center;
    @apply justify-center;
    @apply focus:ring-4;
    @apply font-medium;
    @apply rounded-lg;
    @apply px-4;
    @apply py-2;
    @apply focus:outline-none;
  }
  /* dark:focus:ring-blue-800*/

  .btn.disabled {
    @apply opacity-50;
    @apply cursor-not-allowed;
  }

  .btn-gray {
    @apply bg-gray-400;
    @apply focus:ring-gray-300;
    @apply hover:bg-gray-500;
  }

  :global(.dark) .btn-gray {
    @apply bg-gray-400;
    @apply text-gray-800;
    @apply hover:bg-gray-300;
    @apply focus:ring-gray-500;
  }

  .btn-green {
    @apply bg-green-500;
    @apply focus:ring-green-300;
    @apply hover:bg-green-600;
  }

  :global(.dark) .btn-green {
    @apply bg-lime-500;
    @apply text-gray-800;
    @apply hover:bg-lime-400;
    @apply focus:ring-lime-600;
  }

  .btn-blue {
    @apply bg-blue-500;
    @apply hover:bg-blue-600;
  }

  :global(.dark) .btn-blue {
    @apply bg-sky-400;
    @apply text-gray-800;
    @apply hover:bg-sky-300;
    @apply focus:ring-sky-500;
  }

  .btn-red {
    @apply bg-red-500;
    @apply hover:bg-red-600;
  }

  :global(.dark) .btn-red {
    @apply bg-red-400;
    @apply text-gray-800;
    @apply hover:bg-red-300;
    @apply focus:ring-red-500;
  }

  .btn-purple {
    @apply bg-purple-500;
    @apply hover:bg-purple-600;
  }

  :global(.dark) .btn-red {
    @apply bg-purple-400;
    @apply text-gray-800;
    @apply hover:bg-purple-300;
    @apply focus:ring-purple-500;
  }

  .btn-teal {
    @apply bg-teal-500;
    @apply hover:bg-teal-600;
  }

  :global(.dark) .btn-teal {
    @apply bg-teal-400;
    @apply text-gray-800;
    @apply hover:bg-teal-300;
    @apply focus:ring-teal-500;
  }

  .btn-blue-copy {
    @apply bg-blue-500;
    @apply border-blue-800;
    @apply hover:bg-blue-600;
  }

  :global(.dark) .btn-blue-copy {
    @apply bg-sky-400;
    @apply text-gray-800;
    @apply hover:bg-sky-300;
    @apply focus:ring-sky-500;
  }

  :global([dir='ltr']) .btn-blue-copy {
    @apply rounded-l-none;
  }

  :global([dir='ltr']) .btn-blue-copy {
    @apply rounded-r-none;
  }

  .btn-orange {
    @apply bg-orange-500;
    @apply hover:bg-orange-600;
  }

  :global(.dark) .btn-orange {
    @apply bg-orange-400;
    @apply text-gray-800;
    @apply hover:bg-orange-300;
    @apply focus:ring-orange-500;
  }

  .btn-indigo {
    @apply bg-indigo-500;
    @apply hover:bg-indigo-600;
  }

  :global(.dark) .btn-indigo {
    @apply bg-indigo-400;
    @apply text-gray-800;
    @apply hover:bg-indigo-300;
    @apply focus:ring-indigo-500;
  }

  .btn-yellow {
    @apply bg-yellow-400;
    @apply hover:bg-yellow-500;
    @apply focus:ring-yellow-300;
  }

  :global(.dark) .btn-yellow {
    @apply bg-yellow-300;
    @apply text-gray-800;
    @apply hover:bg-yellow-200;
    @apply focus:ring-yellow-400;
  }
</style>

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
    children
  }: Props = $props();
</script>

<style>
  .btn {
    @apply inline-block;
    @apply leading-tight;
    @apply py-2;
    @apply px-3;
    @apply rounded;
    @apply font-bold;
    @apply text-white;
  }

  .btn.disabled {
    @apply opacity-50;
    @apply cursor-not-allowed;
  }

  .btn-green {
    @apply bg-green-500;
  }

  .btn-green:hover {
    @apply bg-green-600;
  }

  .btn-blue {
    @apply bg-blue-500;
  }

  .btn-blue:hover {
    @apply bg-blue-600;
  }

  .btn-red {
    @apply bg-red-500;
  }

  .btn-red:hover {
    @apply bg-red-600;
  }

  .btn-purple {
    @apply bg-purple-500;
  }

  .btn-purple:hover {
    @apply bg-purple-600;
  }

  .btn-teal {
    @apply bg-teal-500;
  }

  .btn-teal:hover {
    @apply bg-teal-600;
  }

  .btn-blue-copy {
    @apply bg-blue-500;
    @apply border-blue-800;
  }

  [dir='ltr'] .btn-blue-copy {
    @apply rounded-l-none;
  }

  [dir='rtl'] .btn-blue-copy {
    @apply rounded-r-none;
  }

  .btn-blue-copy:hover {
    @apply bg-blue-600;
  }

  .btn-orange {
    @apply bg-orange-500;
  }

  .btn-orange:hover {
    @apply bg-orange-600;
  }

  .btn-indigo {
    @apply bg-indigo-500;
  }

  .btn-indigo:hover {
    @apply bg-indigo-600;
  }
</style>

{#if href != ''}
  <a
    href="{href}"
    class="btn btn-{color}
       {disabled ? 'disabled' : ''}
       {additionalClasses}"
    data-testid="{testid}"
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
    type="{type}"
    disabled="{disabled}"
    data-testid="{testid}"
    {...options}
  >
    {@render children?.()}
  </button>
{/if}

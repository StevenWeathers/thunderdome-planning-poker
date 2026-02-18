<script lang="ts">
  interface Props {
    disabled?: boolean;
    color?: string;
    additionalClasses?: string;
    class?: string;
    type?: string;
    onClick?: any;
    href?: string;
    testid?: string;
    options?: any;
    labelFor?: string;
    size?: 'small' | 'medium' | 'large';
    fullWidth?: boolean;
    children?: import('svelte').Snippet;
  }

  let {
    disabled = false,
    color = 'green',
    additionalClasses = '',
    class: className = '',
    type = 'button',
    onClick = () => {},
    href = '',
    testid = '',
    options = {},
    labelFor = '',
    size = 'medium',
    fullWidth = false,
    children,
  }: Props = $props();

  // Combine additionalClasses and class prop
  const combinedClasses = $derived(`${additionalClasses} ${className}`.trim());

  /**
   * Keyboard interaction for label-as-button scenario.
   * When using a <label> to trigger a hidden input (e.g. file input),
   * add expected button keyboard activation: Enter / Space.
   */
  function handleLabelKeydown(e: KeyboardEvent) {
    if (disabled) return;
    if (e.key === ' ' || e.key === 'Enter') {
      // Prevent page scroll (Space) or form submission (Enter) default behaviors
      e.preventDefault();
      // If a target input id is provided, click it so native interaction occurs
      if (labelFor) {
        const target = document.getElementById(labelFor) as HTMLInputElement | null;
        target?.click();
      } else {
        // Fallback: dispatch click so any onClick in options still fires
        (e.currentTarget as HTMLElement).click();
      }
    }
  }
</script>

{#if type === 'label'}
  <label
    class="btn-hollow btn-hollow-{color} cursor-pointer inline-block {combinedClasses}"
    class:btn-hollow-large={size === 'large'}
    class:btn-hollow-full={fullWidth}
    class:disabled
    data-testid={testid}
    role="button"
    aria-disabled={disabled}
    onkeydown={handleLabelKeydown}
    tabindex="0"
    {...options}
    for={labelFor}
  >
    {@render children?.()}
  </label>
{:else if href === ''}
  <button
    class="btn-hollow btn-hollow-{color}
        {disabled ? 'disabled' : ''}
        {combinedClasses}"
    class:btn-hollow-large={size === 'large'}
    class:btn-hollow-full={fullWidth}
    onclick={onClick}
    {type}
    {disabled}
    data-testid={testid}
    {...options}
  >
    {@render children?.()}
  </button>
{:else}
  <a
    {href}
    class="btn-hollow btn-hollow-{color} inline-block no-underline {combinedClasses}"
    class:btn-hollow-large={size === 'large'}
    class:btn-hollow-full={fullWidth}
    class:disabled
    data-testid={testid}
    {...options}
  >
    {@render children?.()}
  </a>
{/if}

<style lang="postcss">
  .btn-hollow {
    @apply leading-tight;
    @apply font-semibold;
    @apply py-2;
    @apply px-3;
    @apply border-2;
    @apply rounded;
    @apply transition-all;
    @apply duration-200;
    @apply touch-manipulation;
    @apply relative;
    @apply overflow-hidden;
  }

  /* Light mode base */
  .btn-hollow {
    background-color: white;
  }

  /* Dark mode base */
  :global(.dark) .btn-hollow {
    background-color: rgb(31 41 55); /* gray-800 */
  }

  .btn-hollow.disabled {
    @apply opacity-50;
    @apply cursor-not-allowed;
  }

  /* Green (Lime) variant */
  .btn-hollow-green {
    border-color: rgb(217 249 157); /* lime-200 */
    color: rgb(77 124 15); /* lime-700 */
  }

  :global(.dark) .btn-hollow-green {
    border-color: rgb(77 124 15); /* lime-700 */
    color: rgb(190 242 100); /* lime-300 */
  }

  .btn-hollow-green:hover:not(.disabled) {
    background-color: rgb(247 254 231); /* lime-50 */
    border-color: rgb(190 242 100); /* lime-300 */
    color: rgb(63 98 18); /* lime-800 */
  }

  :global(.dark) .btn-hollow-green:hover:not(.disabled) {
    background-color: rgba(77, 124, 15, 0.2); /* lime-900/20 */
    border-color: rgb(101 163 13); /* lime-600 */
    color: rgb(217 249 157); /* lime-200 */
  }

  /* Blue variant */
  .btn-hollow-blue {
    border-color: rgb(191 219 254); /* blue-200 */
    color: rgb(29 78 216); /* blue-700 */
  }

  :global(.dark) .btn-hollow-blue {
    border-color: rgb(29 78 216); /* blue-700 */
    color: rgb(147 197 253); /* blue-300 */
  }

  .btn-hollow-blue:hover:not(.disabled) {
    background-color: rgb(239 246 255); /* blue-50 */
    border-color: rgb(147 197 253); /* blue-300 */
    color: rgb(30 64 175); /* blue-800 */
  }

  :global(.dark) .btn-hollow-blue:hover:not(.disabled) {
    background-color: rgba(29, 78, 216, 0.2); /* blue-900/20 */
    border-color: rgb(37 99 235); /* blue-600 */
    color: rgb(191 219 254); /* blue-200 */
  }

  /* Red variant */
  .btn-hollow-red {
    border-color: rgb(254 202 202); /* red-200 */
    color: rgb(185 28 28); /* red-700 */
  }

  :global(.dark) .btn-hollow-red {
    border-color: rgb(185 28 28); /* red-700 */
    color: rgb(252 165 165); /* red-300 */
  }

  .btn-hollow-red:hover:not(.disabled) {
    background-color: rgb(254 242 242); /* red-50 */
    border-color: rgb(252 165 165); /* red-300 */
    color: rgb(153 27 27); /* red-800 */
  }

  :global(.dark) .btn-hollow-red:hover:not(.disabled) {
    background-color: rgba(185, 28, 28, 0.2); /* red-900/20 */
    border-color: rgb(220 38 38); /* red-600 */
    color: rgb(254 202 202); /* red-200 */
  }

  /* Purple variant */
  .btn-hollow-purple {
    border-color: rgb(221 214 254); /* purple-200 */
    color: rgb(109 40 217); /* purple-700 */
  }

  :global(.dark) .btn-hollow-purple {
    border-color: rgb(109 40 217); /* purple-700 */
    color: rgb(196 181 253); /* purple-300 */
  }

  .btn-hollow-purple:hover:not(.disabled) {
    background-color: rgb(245 243 255); /* purple-50 */
    border-color: rgb(196 181 253); /* purple-300 */
    color: rgb(91 33 182); /* purple-800 */
  }

  :global(.dark) .btn-hollow-purple:hover:not(.disabled) {
    background-color: rgba(109, 40, 217, 0.2); /* purple-900/20 */
    border-color: rgb(124 58 237); /* purple-600 */
    color: rgb(221 214 254); /* purple-200 */
  }

  /* Teal variant */
  .btn-hollow-teal {
    border-color: rgb(153 246 228); /* teal-200 */
    color: rgb(15 118 110); /* teal-700 */
  }

  :global(.dark) .btn-hollow-teal {
    border-color: rgb(15 118 110); /* teal-700 */
    color: rgb(94 234 212); /* teal-300 */
  }

  .btn-hollow-teal:hover:not(.disabled) {
    background-color: rgb(240 253 250); /* teal-50 */
    border-color: rgb(94 234 212); /* teal-300 */
    color: rgb(17 94 89); /* teal-800 */
  }

  :global(.dark) .btn-hollow-teal:hover:not(.disabled) {
    background-color: rgba(15, 118, 110, 0.2); /* teal-900/20 */
    border-color: rgb(20 184 166); /* teal-600 */
    color: rgb(153 246 228); /* teal-200 */
  }

  /* Orange variant */
  .btn-hollow-orange {
    border-color: rgb(254 215 170); /* orange-200 */
    color: rgb(194 65 12); /* orange-700 */
  }

  :global(.dark) .btn-hollow-orange {
    border-color: rgb(194 65 12); /* orange-700 */
    color: rgb(253 186 116); /* orange-300 */
  }

  .btn-hollow-orange:hover:not(.disabled) {
    background-color: rgb(255 247 237); /* orange-50 */
    border-color: rgb(253 186 116); /* orange-300 */
    color: rgb(154 52 18); /* orange-800 */
  }

  :global(.dark) .btn-hollow-orange:hover:not(.disabled) {
    background-color: rgba(194, 65, 12, 0.2); /* orange-900/20 */
    border-color: rgb(234 88 12); /* orange-600 */
    color: rgb(254 215 170); /* orange-200 */
  }

  /* Indigo variant */
  .btn-hollow-indigo {
    border-color: rgb(199 210 254); /* indigo-200 */
    color: rgb(67 56 202); /* indigo-700 */
  }

  :global(.dark) .btn-hollow-indigo {
    border-color: rgb(67 56 202); /* indigo-700 */
    color: rgb(165 180 252); /* indigo-300 */
  }

  .btn-hollow-indigo:hover:not(.disabled) {
    background-color: rgb(238 242 255); /* indigo-50 */
    border-color: rgb(165 180 252); /* indigo-300 */
    color: rgb(55 48 163); /* indigo-800 */
  }

  :global(.dark) .btn-hollow-indigo:hover:not(.disabled) {
    background-color: rgba(67, 56, 202, 0.2); /* indigo-900/20 */
    border-color: rgb(79 70 229); /* indigo-600 */
    color: rgb(199 210 254); /* indigo-200 */
  }

  .btn-hollow-large {
    @apply text-lg;
    @apply py-3;
    @apply px-4;
    @apply rounded-lg;
  }

  .btn-hollow-full {
    @apply w-full;
    @apply text-center;
  }

  /* Focus ring for keyboard users */
  .btn-hollow:focus-visible {
    outline: 2px solid currentColor;
    outline-offset: 2px;
  }
</style>

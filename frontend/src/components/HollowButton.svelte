<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
    export let href = ''
    export let testid = ''
    export let options = {}
</script>

<style>
    :global(.btn-hollow) {
        @apply leading-tight;
        @apply font-semibold;
        @apply bg-transparent;
        @apply py-2;
        @apply px-3;
        @apply border;
        @apply rounded;
    }

    :global(.btn-hollow.disabled) {
        @apply opacity-50;
        @apply cursor-not-allowed;
    }

    :global(.btn-hollow:hover) {
        @apply border-transparent;
        @apply text-white;
    }

    :global(.btn-hollow-green) {
        @apply border-green-500;
        @apply text-green-500;
    }

    :global(.btn-hollow-green:hover) {
        @apply bg-green-600;
    }

    :global(.btn-hollow-blue) {
        @apply border-blue-500;
        @apply text-blue-500;
    }

    :global(.btn-hollow-blue:hover) {
        @apply bg-blue-600;
    }

    :global(.btn-hollow-red) {
        @apply border-red-500;
        @apply text-red-500;
    }

    :global(.btn-hollow-red:hover) {
        @apply bg-red-600;
    }

    :global(.btn-hollow-purple) {
        @apply border-purple-500;
        @apply text-purple-500;
    }

    :global(.btn-hollow-purple:hover) {
        @apply bg-purple-600;
    }

    :global(.btn-hollow-teal) {
        @apply border-teal-500;
        @apply text-teal-500;
    }

    :global(.btn-hollow-teal:hover) {
        @apply bg-teal-600;
    }
</style>

{#if type === 'label'}
    <label
        class="btn-hollow btn-hollow-{color} cursor-pointer {additionalClasses}"
        data-testid="{testid}"
        {...options}
    >
        <slot />
    </label>
{:else if href === ''}
    <button
        class="btn-hollow btn-hollow-{color}
        {disabled ? 'disabled' : ''}
        {additionalClasses}"
        on:click="{onClick}"
        type="{type}"
        disabled="{disabled}"
        data-testid="{testid}"
        {...options}
    >
        <slot />
    </button>
{:else}
    <a
        href="{href}"
        class="btn-hollow btn-hollow-{color} inline-block no-underline {additionalClasses}"
        data-testid="{testid}"
        {...options}
    >
        <slot />
    </a>
{/if}

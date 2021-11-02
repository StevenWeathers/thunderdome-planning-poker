<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
</script>

<style>
    :global(.btn) {
        @apply leading-tight;
        @apply py-2;
        @apply px-3;
        @apply rounded;
        @apply font-bold;
        @apply text-white;
    }
    :global(.btn.disabled) {
        @apply opacity-50;
        @apply cursor-not-allowed;
    }
    :global(.btn-green) {
        @apply bg-green-500;
    }
    :global(.btn-green:hover) {
        @apply bg-green-600;
    }
    :global(.btn-blue) {
        @apply bg-blue-500;
    }
    :global(.btn-blue:hover) {
        @apply bg-blue-600;
    }
    :global(.btn-red) {
        @apply bg-red-500;
    }
    :global(.btn-red:hover) {
        @apply bg-red-600;
    }
    :global(.btn-purple) {
        @apply bg-purple-500;
    }
    :global(.btn-purple:hover) {
        @apply bg-purple-600;
    }
    :global(.btn-teal) {
        @apply bg-teal-500;
    }
    :global(.btn-teal:hover) {
        @apply bg-teal-600;
    }
    :global(.btn-blue-copy) {
        @apply bg-blue-500;
        @apply rounded-l-none;
        @apply border-blue-800;
    }
    :global(.btn-blue-copy:hover) {
        @apply bg-blue-600;
    }
</style>

<button
    class="btn btn-{color}
    {disabled ? 'disabled' : ''}
    {additionalClasses}"
    on:click="{onClick}"
    type="{type}"
    disabled="{disabled}"
>
    <slot />
</button>

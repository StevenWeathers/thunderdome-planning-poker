<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
    export let href = ''
</script>

<style>
    .btn-hollow {
        @apply leading-tight font-semibold bg-transparent py-2 px-3 border rounded;
    }
    .btn-hollow.disabled {
        @apply opacity-50 cursor-not-allowed;
    }
    .btn-hollow:hover {
        @apply border-transparent text-white;
    }
    .btn-hollow-green {
        @apply border-green-500 text-green-500;
    }
    .btn-hollow-green:hover {
        @apply bg-green-600;
    }
    .btn-hollow-blue {
        @apply border-blue-500 text-blue-500;
    }
    .btn-hollow-blue:hover {
        @apply bg-blue-600;
    }
    .btn-hollow-red {
        @apply border-red-500 text-red-500;
    }
    .btn-hollow-red:hover {
        @apply bg-red-600;
    }
    .btn-hollow-purple {
        @apply border-purple-500 text-purple-500;
    }
    .btn-hollow-purple:hover {
        @apply bg-purple-600;
    }
    .btn-hollow-teal {
        @apply border-teal-500 text-teal-500;
    }
    .btn-hollow-teal:hover {
        @apply bg-teal-600;
    }
</style>

{#if href === ''}
    <button
        class="btn-hollow btn-hollow-{color}
        {disabled ? 'disabled' : ''}
        {additionalClasses}"
        on:click="{onClick}"
        {type}
        {disabled}>
        <slot />
    </button>
{:else}
    <a
        {href}
        class="btn-hollow btn-hollow-{color} inline-block no-underline {additionalClasses}">
        <slot />
    </a>
{/if}

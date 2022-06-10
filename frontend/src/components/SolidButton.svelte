<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let testid = ''
    export let onClick = () => {}
    export let options = {}
    export let href = ''
</script>

<style>
    .btn {
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
        @apply rounded-l-none;
        @apply border-blue-800;
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
        <slot />
    </a>
{:else}
    <button
        class="btn btn-{color}
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
{/if}

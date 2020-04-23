<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
</script>

<style>
    .btn {
        @apply leading-tight py-2 px-3 rounded font-bold text-white;
    }
    .btn.disabled {
        @apply opacity-50 cursor-not-allowed;
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
        @apply bg-blue-500 rounded-l-none border-blue-800;
    }
    .btn-blue-copy:hover {
        @apply bg-blue-600;
    }
</style>

<button
    class="btn btn-{color}
    {disabled ? 'disabled' : ''}
    {additionalClasses}"
    on:click="{onClick}"
    {type}
    {disabled}>
    <slot />
</button>

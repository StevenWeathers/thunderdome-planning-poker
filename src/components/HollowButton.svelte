<script>
    export let disabled = false
    export let color = 'green'
    export let textcolor = 'white'
    export let additionalClasses = ''
    export let type = "button"
    export let onClick = () => {}
    export let href = ''

    const disabledButtonClass = 'opacity-50 cursor-not-allowed'
</script>

{#if href === ''}
    <button
        class="bg-transparent hover:bg-{color} text-{color}-dark font-semibold hover:text-{textcolor} py-2 px-2 border border-{color} hover:border-transparent rounded {disabled ? disabledButtonClass : `hover:bg-${color}-dark`} {additionalClasses}"
        on:click={onClick}
        {type}
        {disabled}
    >
        <slot></slot>
    </button>
{:else}
    <a
        {href}
        class="inline-block no-underline bg-transparent hover:bg-{color} text-{color}-dark font-semibold hover:text-{textcolor} py-2 px-2 border border-{color} hover:border-transparent rounded {additionalClasses}"
    >
        <slot></slot>
    </a>
{/if}
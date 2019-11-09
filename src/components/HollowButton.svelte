<script>
    export let disabled = false
    export let color = 'green'
    export let textcolor = 'white'
    export let additionalClasses = ''
    export let type = "button"
    export let onClick = () => {}
    export let href = ''

    const disabledButtonClass = 'opacity-50 cursor-not-allowed'
    const commonClasses = 'leading-tight font-semibold bg-transparent py-2 px-3 border hover:border-transparent rounded'

    $: colorScaled = color !== 'white' ? `${color}-500` : color
    $: colorScaledDark = color !== 'white' ? `${color}-600` : color
    $: textColorScaled = textcolor !== 'white' ? `${textcolor}-500` : textcolor
</script>

{#if href === ''}
    <button
        class="{commonClasses} hover:bg-{colorScaled} text-{colorScaledDark} hover:text-{textColorScaled} border-{colorScaled} {disabled ? disabledButtonClass : `hover:bg-${colorScaledDark}`} {additionalClasses}"
        on:click={onClick}
        {type}
        {disabled}
    >
        <slot></slot>
    </button>
{:else}
    <a
        {href}
        class="{commonClasses} inline-block no-underline hover:bg-{colorScaled} text-{colorScaledDark} hover:text-{textColorScaled} border-{colorScaled} {additionalClasses}"
    >
        <slot></slot>
    </a>
{/if}
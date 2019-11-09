<script>
    export let disabled = false
    export let color = 'green'
    export let textcolor = 'white'
    export let additionalClasses = ''
    export let type = "button"
    export let onClick = () => {}

    const disabledButtonClass = 'opacity-50 cursor-not-allowed'

    $: colorScaled = color !== 'white' ? `${color}-500` : color
    $: colorScaledDark = color !== 'white' ? `${color}-600` : color
    $: textColorScaled = textcolor !== 'white' ? `${textcolor}-500` : textcolor
</script>

<button
    class="leading-tight bg-{colorScaled} text-{textColorScaled} font-bold py-2 px-3 rounded {disabled ? disabledButtonClass : `hover:bg-${colorScaledDark}`} {additionalClasses}"
    on:click={onClick}
    {type}
    {disabled}
>
    <slot></slot>
</button>
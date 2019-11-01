<script>
    export let disabled = false
    export let color = 'green'
    export let textcolor = 'white'
    export let additionalClasses = ''
    export let type = "button"
    export let onClick = () => {}

    const disabledButtonClass = 'opacity-50 cursor-not-allowed'
</script>

<button
    class="bg-{color} text-{textcolor} font-bold py-2 px-4 rounded {disabled ? disabledButtonClass : `hover:bg-${color}-dark`} {additionalClasses}"
    on:click={onClick}
    {type}
    {disabled}
>
    <slot></slot>
</button>
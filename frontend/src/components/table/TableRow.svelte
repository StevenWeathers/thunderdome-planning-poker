<script>
    export let itemIndex = 0
</script>

<tr
    class:bg-slate-100="{itemIndex % 2 !== 0}"
    class:dark:bg-gray-800="{itemIndex % 2 !== 0}"
>
    <slot />
</tr>

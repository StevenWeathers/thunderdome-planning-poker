<script>
    export let type = 'default'
</script>

{#if type === 'default'}
    <td class="px-6 py-4 whitespace-nowrap">
        <slot />
    </td>
{:else if type === 'action'}
    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
        <slot />
    </td>
{/if}

<script>
    export let type = 'default'
</script>

{#if type === 'default'}
    <td class="px-6 py-4 whitespace-nowrap">
        <slot />
    </td>
{:else if type === 'action'}
    <td
        class="px-6 py-4 whitespace-nowrap ltr:text-right rtl:text-left text-sm font-medium"
    >
        <slot />
    </td>
{/if}

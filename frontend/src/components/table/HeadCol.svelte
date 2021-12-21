<script>
    export let type = 'default'
</script>

{#if type === 'default'}
    <th
        scope="col"
        class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
    >
        <slot />
    </th>
{:else if type === 'action'}
    <th scope="col" class="relative px-6 py-3">
        <slot />
    </th>
{/if}

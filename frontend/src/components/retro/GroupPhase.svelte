<script>
    import {
        dndzone,
        SHADOW_ITEM_MARKER_PROPERTY_NAME,
    } from 'svelte-dnd-action'
    import GroupNameForm from './GroupNameForm.svelte'

    export let groups = []
    export let handleItemChange = () => {}
    export let handleGroupNameChange = () => {}

    function handleDndConsider(e) {
        const groupIndex = e.target.dataset.groupindex

        groups[groupIndex].items = e.detail.items
        groups = groups
    }

    function handleDndFinalize(e) {
        const groupIndex = e.target.dataset.groupindex
        const itemId = e.detail.info.id
        const groupId = groups[groupIndex].id

        groups[groupIndex].items = e.detail.items
        groups = groups

        if (groups[groupIndex].items.find(i => i.id === itemId)) {
            handleItemChange(itemId, groupId)
        }
    }
</script>

{#each groups as group, i (group.id)}
    <div
        class="border-2 p-2 dark:border-gray-800 rounded flex flex-col flex-wrap"
    >
        <div class="mb-2">
            <GroupNameForm
                groupName="{group.name}"
                groupId="{group.id}"
                handleGroupNameChange="{handleGroupNameChange}"
            />
        </div>
        <div
            use:dndzone="{{
                items: group.items,
                type: 'item',
                dropTargetStyle: '',
                dropTargetClasses: [
                    'outline',
                    'outline-2',
                    'outline-indigo-500',
                    'dark:outline-yellow-400',
                ],
            }}"
            on:consider="{handleDndConsider}"
            on:finalize="{handleDndFinalize}"
            data-groupindex="{i}"
            class="flex-1 grow"
            style="min-height: 40px;"
        >
            {#each group.items as item, ii (item.id)}
                <div
                    class="relative p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-l-4 dark:text-white"
                    class:border-green-400="{item.type === 'worked'}"
                    class:dark:border-lime-400="{item.type === 'worked'}"
                    class:border-red-500="{item.type === 'improve'}"
                    class:border-blue-400="{item.type === 'question'}"
                    class:dark:border-sky-400="{item.type === 'question'}"
                    data-itemid="{item.id}"
                >
                    {item.content}
                    {#if item[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
                        <div
                            class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-l-4 dark:text-white"
                            class:border-green-400="{item.type === 'worked'}"
                            class:dark:border-lime-400="{item.type ===
                                'worked'}"
                            class:border-red-500="{item.type === 'improve'}"
                            class:border-blue-400="{item.type === 'question'}"
                            class:dark:border-sky-400="{item.type ===
                                'question'}"
                            style="min-height: 40px;"
                        >
                            {item.content}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
{/each}

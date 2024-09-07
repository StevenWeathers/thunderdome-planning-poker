<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import GroupNameForm from './GroupNameForm.svelte';

  export let groups = [];
  export let handleItemChange = () => {};
  export let handleGroupNameChange = () => {};
  export let columns = [];

  function handleDndConsider(e) {
    const groupIndex = e.target.dataset.groupindex;

    groups[groupIndex].items = e.detail.items;
    groups = groups;
  }

  function handleDndFinalize(e) {
    const groupIndex = e.target.dataset.groupindex;
    const itemId = e.detail.info.id;
    const groupId = groups[groupIndex].id;

    groups[groupIndex].items = e.detail.items;
    groups = groups;

    if (groups[groupIndex].items.find(i => i.id === itemId)) {
      handleItemChange(itemId, groupId);
    }
  }

  $: columnColors =
    columns &&
    columns.reduce((p, c) => {
      p[c.name] = c.color;
      return p;
    }, {});
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
          class="relative p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-s-4 dark:text-white"
          class:border-green-400="{columnColors[item.type] === 'green'}"
          class:dark:border-lime-400="{columnColors[item.type] === 'green'}"
          class:border-red-500="{columnColors[item.type] === 'red'}"
          class:border-blue-400="{columnColors[item.type] === 'blue'}"
          class:dark:border-sky-400="{columnColors[item.type] === 'blue'}"
          class:border-yellow-500="{columnColors[item.type] === 'yellow'}"
          class:dark:border-yellow-400="{columnColors[item.type] === 'yellow'}"
          class:border-orange-500="{columnColors[item.type] === 'orange'}"
          class:dark:border-orange-400="{columnColors[item.type] === 'orange'}"
          class:border-teal-500="{columnColors[item.type] === 'teal'}"
          class:dark:border-teal-400="{columnColors[item.type] === 'teal'}"
          class:border-indigo-500="{columnColors[item.type] === 'purple'}"
          class:dark:border-indigo-400="{columnColors[item.type] === 'purple'}"
          data-itemid="{item.id}"
        >
          {item.content}
          {#if item[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
            <div
              class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-s-4 dark:text-white"
              class:border-green-400="{columnColors[item.type] === 'green'}"
              class:dark:border-lime-400="{columnColors[item.type] === 'green'}"
              class:border-red-500="{columnColors[item.type] === 'red'}"
              class:border-blue-400="{columnColors[item.type] === 'blue'}"
              class:dark:border-sky-400="{columnColors[item.type] === 'blue'}"
              class:border-yellow-500="{columnColors[item.type] === 'yellow'}"
              class:dark:border-yellow-400="{columnColors[item.type] ===
                'yellow'}"
              class:border-orange-500="{columnColors[item.type] === 'orange'}"
              class:dark:border-orange-400="{columnColors[item.type] ===
                'orange'}"
              class:border-teal-500="{columnColors[item.type] === 'teal'}"
              class:dark:border-teal-400="{columnColors[item.type] === 'teal'}"
              class:border-indigo-500="{columnColors[item.type] === 'purple'}"
              class:dark:border-indigo-400="{columnColors[item.type] ===
                'purple'}"
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

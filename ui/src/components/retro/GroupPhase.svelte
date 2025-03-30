<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import GroupNameForm from './GroupNameForm.svelte';
  import RetroFeedbackItem from './RetroFeedbackItem.svelte';

  interface Props {
    phase?: string;
    groups?: any;
    handleItemChange?: any;
    handleGroupNameChange?: any;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
  }

  let {
    phase = 'group',
    groups = $bindable([]),
    handleItemChange = (itemId: string, groupId: string) => {},
    handleGroupNameChange = () => {},
    isFacilitator = false,
    users = [],
    columnColors = {},
    sendSocketEvent = (event: string, value: any) => {}
  }: Props = $props();

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
</script>

{#each groups as group, i (group.id)}
  <div
    class="p-3 bg-white dark:bg-gray-800 rounded-lg shadow-lg flex flex-col flex-wrap text-gray-800 dark:text-white"
  >
    <div class="mb-4">
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
      onconsider={handleDndConsider}
      onfinalize={handleDndFinalize}
      data-groupindex="{i}"
      class="flex-1 grow"
      style="min-height: 40px;"
    >
      {#each group.items as item, ii (item.id)}
        <RetroFeedbackItem
          item="{item}"
          class="relative"
          phase="{phase}"
          users="{users}"
          sendSocketEvent="{sendSocketEvent}"
          isFacilitator="{isFacilitator}"
          columnColors="{columnColors}"
        >
          {#if item[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
            <RetroFeedbackItem
              phase="{phase}"
              item="{item}"
              class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible min-h-[40px]"
              columnColors="{columnColors}"
            />
          {/if}
        </RetroFeedbackItem>
      {/each}
    </div>
  </div>
{/each}

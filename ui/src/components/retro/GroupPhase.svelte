<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME, SOURCES, TRIGGERS } from 'svelte-dnd-action';
  import { Plus } from '@lucide/svelte';
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
    sendSocketEvent = (event: string, value: any) => {},
  }: Props = $props();

  // group ids on screen and the item being dragged, null when not dragging
  let idsAtDragStart = $state<string[] | null>(null);
  let draggedItemId = $state<string | null>(null);

  // empty unnamed groups are hidden
  const isInUse = (group: any) => group.items.length > 0 || !!group.name;

  // a hidden group serves as the [+] drop target, shown only while dragging, to create a new group.
  // If someone else fills it mid drag, another hidden group takes its place.
  const newGroup = $derived(
    groups.findLast(
      (group: any) =>
        !group.name &&
        !idsAtDragStart?.includes(group.id) &&
        group.items.every((item: any) => item.id === draggedItemId),
    ),
  );

  // while dragging nothing is removed, so drop targets don't move or disappear from under the pointer
  const visibleGroups = $derived.by(() => {
    const inUse = groups.filter((g: any) => isInUse(g) && g.id !== newGroup?.id);

    if (idsAtDragStart === null) {
      return inUse;
    }

    return [
      ...idsAtDragStart.map(id => groups.find((g: any) => g.id === id)).filter(Boolean),
      ...inUse.filter((g: any) => !idsAtDragStart?.includes(g.id)),
    ];
  });

  const tiles = $derived(newGroup ? [...visibleGroups, newGroup] : visibleGroups);

  function findGroup(e: CustomEvent) {
    const groupId = (e.target as HTMLElement).dataset.groupid;

    return groups.find((g: any) => g.id === groupId);
  }

  function endDrag() {
    idsAtDragStart = null;
    draggedItemId = null;
  }

  function handleDndConsider(e: CustomEvent) {
    const group = findGroup(e);

    // keyboard drags end with a consider event rather than a finalize
    if (e.detail.info.trigger === TRIGGERS.DRAG_STOPPED) {
      endDrag();
    } else {
      idsAtDragStart ??= visibleGroups.map((g: any) => g.id);
      draggedItemId = e.detail.info.id;
    }

    if (group) {
      group.items = e.detail.items;
    }
  }

  function handleDndFinalize(e: CustomEvent) {
    const group = findGroup(e);
    const itemId = e.detail.info.id;

    // keyboard drags finalize on every move and carry on afterwards
    if (e.detail.info.source !== SOURCES.KEYBOARD) {
      endDrag();
    }

    if (!group) {
      return;
    }

    group.items = e.detail.items;

    if (group.items.find((i: any) => i.id === itemId)) {
      handleItemChange(itemId, group.id);
    }
  }
</script>

{#each tiles as group (group.id)}
  {@const isNewGroup = group.id === newGroup?.id}
  <div
    class="relative p-3 rounded-lg flex flex-col flex-wrap text-gray-800 dark:text-white {isNewGroup
      ? 'min-h-24 border-2 border-dashed border-gray-400 dark:border-gray-500'
      : 'bg-white dark:bg-gray-800 shadow-lg'}"
    class:invisible={isNewGroup && idsAtDragStart === null}
    data-testid={isNewGroup ? 'retro-new-group' : 'retro-group'}
  >
    {#if !isNewGroup}
      <div class="mb-4">
        <GroupNameForm groupName={group.name} groupId={group.id} {handleGroupNameChange} />
      </div>
    {:else if group.items.length === 0}
      <div class="absolute inset-0 flex items-center justify-center pointer-events-none text-gray-400">
        <Plus class="w-8 h-8" aria-hidden="true" />
      </div>
    {/if}
    <div
      use:dndzone={{
        items: group.items,
        type: 'item',
        dropTargetStyle: {},
        dropTargetClasses: ['outline', 'outline-2', 'outline-indigo-500', 'dark:outline-yellow-400'],
      }}
      onconsider={handleDndConsider}
      onfinalize={handleDndFinalize}
      data-groupid={group.id}
      class="flex-1 grow"
      style="min-height: 40px;"
    >
      {#each group.items as item (item.id)}
        <RetroFeedbackItem {item} class="relative" {phase} {users} {sendSocketEvent} {isFacilitator} {columnColors}>
          {#if item[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
            <RetroFeedbackItem
              {phase}
              {item}
              class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible min-h-[40px]"
              {columnColors}
            />
          {/if}
        </RetroFeedbackItem>
      {/each}
    </div>
  </div>
{/each}

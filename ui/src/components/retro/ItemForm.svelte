<script lang="ts">
  import { Angry, CircleHelp, Frown, Smile } from '@lucide/svelte';
  import RetroFeedbackItem from './RetroFeedbackItem.svelte';

  interface Props {
    sendSocketEvent?: any;
    itemType?: string;
    content?: string;
    newItemPlaceholder?: string;
    phase?: string;
    isFacilitator?: boolean;
    items?: any;
    users?: any;
    feedbackVisibility?: string;
    icon?: string;
    color?: string;
    columnColors?: any;
  }

  let {
    sendSocketEvent = (event: string, value: any) => {},
    itemType = '',
    content = $bindable(''),
    newItemPlaceholder = '',
    phase = 'brainstorm',
    isFacilitator = false,
    items = [],
    users = [],
    feedbackVisibility = 'visible',
    icon = '',
    color = 'blue',
    columnColors = {},
  }: Props = $props();

  const handleFormSubmit = (evt: Event) => {
    evt.preventDefault();

    sendSocketEvent(
      `create_item`,
      JSON.stringify({
        type: itemType,
        content,
        phase: phase,
      }),
    );
    content = '';
  };
</script>

<div class="p-3 bg-white dark:bg-gray-800 rounded-lg shadow-lg flex flex-col flex-wrap text-gray-800 dark:text-white">
  <div class="flex items-center mb-4">
    {#if icon !== ''}
      <div
        class="flex-shrink pe-2"
        class:text-green-400={color === 'green'}
        class:dark:text-lime-400={color === 'green'}
        class:text-red-500={color === 'red'}
        class:text-blue-400={color === 'blue'}
        class:dark:text-sky-400={color === 'blue'}
        class:text-yellow-500={color === 'yellow'}
        class:dark:text-yellow-400={color === 'yellow'}
        class:text-orange-500={color === 'orange'}
        class:dark:text-orange-400={color === 'orange'}
        class:text-teal-500={color === 'teal'}
        class:dark:text-teal-400={color === 'teal'}
        class:text-indigo-500={color === 'purple'}
        class:dark:text-indigo-400={color === 'purple'}
      >
        {#if icon === 'smiley'}
          <Smile class="w-8 h-8" />
        {:else if icon === 'frown'}
          <Frown class="w-8 h-8" />
        {:else if icon === 'question'}
          <CircleHelp class="w-8 h-8" />
        {:else if icon === 'angry'}
          <Angry class="w-8 h-8" />
        {/if}
      </div>
    {/if}
    <div class="flex-grow">
      <form onsubmit={handleFormSubmit} class="flex">
        <input
          bind:value={content}
          placeholder={newItemPlaceholder}
          class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                    px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                    focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
          id="new{itemType}"
          name="new{itemType}"
          type="text"
          required
        />
        <button type="submit" class="hidden">submit</button>
      </form>
    </div>
  </div>
  <div>
    {#each items.filter(i => i.type === itemType) as item}
      <RetroFeedbackItem {item} {phase} {users} {isFacilitator} {sendSocketEvent} {columnColors} {feedbackVisibility} />
    {/each}
  </div>
</div>

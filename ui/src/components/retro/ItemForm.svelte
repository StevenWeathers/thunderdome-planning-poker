<script lang="ts">
  import { Angry, CircleQuestionMark, Frown, Smile } from '@lucide/svelte';
  import GrowingTextArea from '../global/GrowingTextArea.svelte';
  import RetroFeedbackItem from './RetroFeedbackItem.svelte';
  import type { RetroItem, RetroUser } from '../../types/retro';

  interface Props {
    sendSocketEvent?: any;
    itemType?: string;
    content?: string;
    newItemPlaceholder?: string;
    phase?: string;
    isFacilitator?: boolean;
    items?: RetroItem[];
    users?: RetroUser[];
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

  let textareaComponent: any;

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
    textareaComponent?.resetHeight();
  };
</script>

<div class="p-3 bg-white dark:bg-gray-800 rounded-lg shadow-lg flex flex-col flex-wrap text-gray-800 dark:text-white">
  <div class="flex items-start mb-4">
    {#if icon !== ''}
      <div
        class="flex-shrink pt-1 pe-2"
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
          <CircleQuestionMark class="w-8 h-8" />
        {:else if icon === 'angry'}
          <Angry class="w-8 h-8" />
        {/if}
      </div>
    {/if}
    <div class="flex-grow">
      <form onsubmit={handleFormSubmit} class="flex flex-col w-full min-w-0">
        <GrowingTextArea
          bind:this={textareaComponent}
          bind:value={content}
          placeholder={newItemPlaceholder}
          id="new{itemType}"
          name="new{itemType}"
          required
          onkeydown={(e: KeyboardEvent) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              handleFormSubmit(e as unknown as Event);
            }
          }}
        />
        <button type="submit" class="hidden">submit</button>
      </form>
    </div>
  </div>
  <div class="min-w-0">
    {#each items.filter(i => i.type === itemType) as item}
      <RetroFeedbackItem {item} {phase} {users} {isFacilitator} {sendSocketEvent} {columnColors} {feedbackVisibility} />
    {/each}
  </div>
</div>

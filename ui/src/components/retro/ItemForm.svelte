<script lang="ts">
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import ItemComments from './ItemComments.svelte';
  import {
    Angry,
    CircleHelp,
    Frown,
    MessageSquareMore,
    Smile,
    Trash2,
  } from 'lucide-svelte';

  export let sendSocketEvent = (event: string, value: any) => {};
  export let itemType = '';
  export let content = '';
  export let newItemPlaceholder = '';
  export let phase = 'brainstorm';
  export let isFacilitator = false;
  export let items = [];
  export let users = [];
  export let feedbackVisibility = 'visible';
  export let icon = '';
  export let color = 'blue';

  let showComments = false;
  let selectedItemId = null;

  const toggleComments = itemId => () => {
    showComments = !showComments;
    selectedItemId = itemId;
  };

  const handleDelete = (type, id) => () => {
    sendSocketEvent(
      `delete_item`,
      JSON.stringify({
        id,
        type,
        phase,
      }),
    );
  };

  const handleFormSubmit = evt => {
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

<div class="">
  <div class="flex items-center mb-4">
    {#if icon !== ''}
      <div
        class="flex-shrink pe-2"
        class:text-green-400="{color === 'green'}"
        class:dark:text-lime-400="{color === 'green'}"
        class:text-red-500="{color === 'red'}"
        class:text-blue-400="{color === 'blue'}"
        class:dark:text-sky-400="{color === 'blue'}"
        class:text-yellow-500="{color === 'yellow'}"
        class:dark:text-yellow-400="{color === 'yellow'}"
        class:text-orange-500="{color === 'orange'}"
        class:dark:text-orange-400="{color === 'orange'}"
        class:text-teal-500="{color === 'teal'}"
        class:dark:text-teal-400="{color === 'teal'}"
        class:text-indigo-500="{color === 'purple'}"
        class:dark:text-indigo-400="{color === 'purple'}"
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
      <form on:submit="{handleFormSubmit}" class="flex">
        <input
          bind:value="{content}"
          placeholder="{newItemPlaceholder}"
          class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                    px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                    focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
          id="new{itemType}"
          name="new{itemType}"
          type="text"
          required
          disabled="{phase !== 'brainstorm' && !isFacilitator}"
        />
        <button type="submit" class="hidden"></button>
      </form>
    </div>
  </div>
  <div>
    {#each items.filter(i => i.type === itemType) as item}
      <div
        class="p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-s-4"
        class:border-green-400="{color === 'green'}"
        class:dark:border-lime-400="{color === 'green'}"
        class:border-red-500="{color === 'red'}"
        class:border-blue-400="{color === 'blue'}"
        class:dark:border-sky-400="{color === 'blue'}"
        class:border-yellow-500="{color === 'yellow'}"
        class:dark:border-yellow-400="{color === 'yellow'}"
        class:border-orange-500="{color === 'orange'}"
        class:dark:border-orange-400="{color === 'orange'}"
        class:border-teal-500="{color === 'teal'}"
        class:dark:border-teal-400="{color === 'teal'}"
        class:border-indigo-500="{color === 'purple'}"
        class:dark:border-indigo-400="{color === 'purple'}"
        data-itemType="{itemType}"
        data-itemId="{item.id}"
      >
        <div class="flex items-center">
          <div class="flex-grow">
            <div class="flex items-center">
              <div class="flex-grow dark:text-gray-200">
                {#if feedbackVisibility === 'hidden' && item.userId !== $user.id}
                  <span class="italic">{$LL.retroFeedbackHidden()}</span>
                {:else if feedbackVisibility === 'concealed' && item.userId !== $user.id}
                  <span class="italic"
                    >{$LL.retroFeedbackConcealed()}&nbsp;&nbsp;</span
                  ><span class="text-white dark:text-gray-800"
                    >{item.content}</span
                  >
                {:else}
                  {item.content}
                {/if}
              </div>
            </div>
          </div>
          <div class="flex-none flex gap-x-2 ps-2">
            <div>
              <button
                class="inline-block align-middle text-blue-400 dark:text-sky-400"
                on:click="{toggleComments(item.id)}"
              >
                {item.comments.length}&nbsp;
                <MessageSquareMore
                  width="14"
                  height="14"
                  class="inline-block"
                />
              </button>
            </div>
            {#if phase === 'brainstorm'}
              <div>
                <button
                  on:click="{handleDelete(itemType, item.id)}"
                  class="inline-block align-middle {item.userId !== $user.id
                    ? 'text-gray-300 dark:text-gray-600 cursor-not-allowed'
                    : 'text-gray-500 dark:text-gray-400 hover:text-red-500'}"
                  disabled="{item.userId !== $user.id}"
                >
                  <Trash2 />
                </button>
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>

  {#if showComments}
    <ItemComments
      toggleComments="{toggleComments()}"
      selectedItemId="{selectedItemId}"
      items="{items}"
      users="{users}"
      isFacilitator="{isFacilitator}"
      sendSocketEvent="{sendSocketEvent}"
    />
  {/if}
</div>

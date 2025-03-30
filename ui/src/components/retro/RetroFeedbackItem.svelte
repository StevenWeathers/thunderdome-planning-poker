<script lang="ts">
  import { MessageSquare, Trash2 } from 'lucide-svelte';
  import ItemComments from './ItemComments.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';

  
  interface Props {
    class?: string;
    phase?: string;
    item?: any;
    feedbackVisibility?: string;
    isFacilitator?: boolean;
    users?: any;
    columnColors?: any;
    sendSocketEvent?: any;
    children?: import('svelte').Snippet;
  }

  let {
    class: klass = '',
    phase = '',
    item = {
    id: '',
    type: '',
    content: '',
    comments: [],
  },
    feedbackVisibility = 'visible',
    isFacilitator = false,
    users = [],
    columnColors = {},
    sendSocketEvent = (event: string, value: any) => {},
    children
  }: Props = $props();

  let showComments = $state(false);
  let selectedItem = $state(null);

  const toggleComments = item => () => {
    showComments = !showComments;
    selectedItem = item;
  };

  const handleDelete = () => {
    sendSocketEvent(
      `delete_item`,
      JSON.stringify({
        id: item.id,
        type: item.type,
        phase,
      }),
    );
  };

  const getTypeTagColors = (type: string) => {
    const typeColor = columnColors[type];
    switch (typeColor) {
      case 'green':
        return 'bg-green-300 text-green-800 dark:bg-green-600 dark:text-green-200';
      case 'red':
        return 'bg-red-300 text-red-800 dark:bg-red-600 dark:text-red-200';
      case 'blue':
        return 'bg-blue-300 text-blue-800 dark:bg-blue-600 dark:text-blue-200';
      case 'yellow':
        return 'bg-yellow-300 text-yellow-800 dark:bg-yellow-600 dark:text-yellow-200';
      case 'orange':
        return 'bg-orange-300 text-orange-800 dark:bg-orange-600 dark:text-orange-200';
      case 'teal':
        return 'bg-teal-300 text-teal-800 dark:bg-teal-600 dark:text-teal-200';
      case 'purple':
        return 'bg-purple-300 text-purple-800 dark:bg-purple-600 dark:text-purple-200';
      default:
        return 'bg-gray-300 text-gray-800 dark:bg-gray-600 dark:text-gray-200';
    }
  };
</script>

<div
  class="{klass} p-2 mb-3 border-s-4 border-gray-300 dark:border-gray-600 transition-all duration-200 ease-in-out hover:border-blue-500 dark:hover:border-sky-400"
  data-testid="retro-feedback-item"
  data-itemid="{item.id}"
>
  <div class="flex items-center justify-between mb-1">
    <div class="flex items-center space-x-2">
      {#if phase !== 'brainstorm'}
        <span
          class="text-xs font-medium px-2 py-1 rounded-full {getTypeTagColors(
            item.type,
          )}"
          data-testid="retro-feedback-item-type"
        >
          {item.type}
        </span>
      {/if}
      <button
        class="inline-block leading-none text-gray-700 dark:text-gray-300 hover:text-blue-500 dark:hover:text-sky-400 transition-colors duration-200"
        class:cursor-not-allowed="{phase === 'brainstorm' &&
          feedbackVisibility === 'hidden'}"
        onclick={toggleComments(item)}
        disabled="{phase === 'brainstorm' && feedbackVisibility === 'hidden'}"
      >
        <MessageSquare class="inline w-4 h-4" />
        <span data-testid="retro-feedback-item-comments"
          >{item.comments.length}</span
        >
      </button>
    </div>
    {#if phase === 'brainstorm' && item.userId === $user.id}
      <button
        aria-label="Delete feedback"
        onclick={handleDelete}
        class="float-right inline-block leading-none text-gray-700 dark:text-gray-300 hover:text-red-500 dark:hover:text-red-400 transition-colors duration-200"
      >
        <Trash2 class="w-5 h-5" />
      </button>
    {/if}
  </div>
  <p data-testid="retro-feedback-item-content">
    {#if phase === 'brainstorm' && feedbackVisibility === 'hidden' && item.userId !== $user.id}
      <span class="italic">{$LL.retroFeedbackHidden()}</span>
    {:else if phase === 'brainstorm' && feedbackVisibility === 'concealed' && item.userId !== $user.id}
      <span class="italic">{$LL.retroFeedbackConcealed()}&nbsp;&nbsp;</span
      ><span class="text-white dark:text-gray-800">{item.content}</span>
    {:else}
      {item.content}
    {/if}
  </p>
  {@render children?.()}

  {#if showComments}
    <ItemComments
      toggleComments="{toggleComments(null)}"
      selectedItem="{selectedItem}"
      item="{item}"
      users="{users}"
      isFacilitator="{isFacilitator}"
      sendSocketEvent="{sendSocketEvent}"
    />
  {/if}
</div>

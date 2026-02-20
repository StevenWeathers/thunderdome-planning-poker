<script lang="ts">
  import GrowingTextArea from '../global/GrowingTextArea.svelte';
  import { SquareCheckBig } from '@lucide/svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    onsubmit?: (content: string) => void;
  }

  let { onsubmit }: Props = $props();

  let actionItem = $state('');
  let textareaComponent: any;
  let formElement: HTMLFormElement;

  const handleSubmit = (evt: Event) => {
    evt.preventDefault();
    if (actionItem.trim()) {
      onsubmit?.(actionItem);
      actionItem = '';
      textareaComponent?.resetHeight();
    }
  };

  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      formElement?.requestSubmit();
    }
  };
</script>

<div class="flex items-start mb-4">
  <div class="flex-shrink pe-2 pt-1">
    <SquareCheckBig class="w-8 h-8 text-indigo-500 dark:text-violet-400" />
  </div>
  <div class="flex-grow">
    <form bind:this={formElement} onsubmit={handleSubmit}>
      <GrowingTextArea
        bind:this={textareaComponent}
        bind:value={actionItem}
        placeholder={$LL.actionItemPlaceholder()}
        id="actionItem"
        name="actionItem"
        required
        onkeydown={handleKeydown}
      />
      <button type="submit" class="hidden">submit</button>
    </form>
  </div>
</div>

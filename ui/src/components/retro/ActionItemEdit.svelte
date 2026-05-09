<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import GrowingTextArea from '../global/GrowingTextArea.svelte';
  import { onMount } from 'svelte';

  import type { RetroAction } from '../../types/retro';

  interface Props {
    toggleEdit?: any;
    handleEdit?: any;
    handleDelete?: any;
    action?: RetroAction | null;
  }

  let {
    toggleEdit = () => {},
    handleEdit = (action: any) => {},
    handleDelete = () => {},
    action = null,
  }: Props = $props();

  const resolvedAction = $derived(
    action ?? {
      comments: [],
      id: '',
      retroId: '',
      content: '',
      completed: false,
      assignees: [],
    },
  );

  let editAction = $state({
    id: '',
    retroId: '',
    content: '',
    completed: false,
  });

  let textareaComponent: any;

  $effect(() => {
    editAction.id = resolvedAction.id;
    editAction.retroId = resolvedAction.retroId;
    editAction.content = resolvedAction.content;
    editAction.completed = resolvedAction.completed;
  });

  const handleSubmit = (e: Event) => {
    e.preventDefault();
    handleEdit(editAction);
  };

  onMount(() => {
    textareaComponent?.focus();
  });
</script>

<Modal closeModal={toggleEdit} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabel={$LL.modalEditRetroActionItem()}>
  <form onsubmit={handleSubmit}>
    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="actionItem">
        {$LL.actionItem()}
      </label>
      <div class="control">
        <GrowingTextArea
          bind:this={textareaComponent}
          bind:value={editAction.content}
          placeholder={$LL.actionItemPlaceholder()}
          id="actionItem"
          name="actionItem"
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <div class="control">
        <div class="flex-shrink">
          <Checkbox id="Completed" bind:checked={editAction.completed} label={$LL.completed()} />
        </div>
      </div>
    </div>

    <div class="flex w-full pt-4">
      <div class="w-1/2">
        <HollowButton color="red" onClick={handleDelete(editAction)}>{$LL.delete()}</HollowButton>
      </div>
      <div class="w-1/2 text-right">
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>

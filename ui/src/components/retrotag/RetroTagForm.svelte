<script lang="ts">
  import { onMount, untrack } from 'svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { RETRO_TAG_MAX_LENGTH } from '../retro/retroTags';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { RetroTag } from '../../types/retro';

  interface Props {
    tag?: RetroTag | null;
    toggleClose?: () => void;
    handleSaved?: () => void;
    apiPrefix?: string;
    xfetch: ApiClient;
    notifications: NotificationService;
  }

  let {
    tag = null,
    toggleClose = () => {},
    handleSaved = () => {},
    apiPrefix = '/api',
    xfetch,
    notifications,
  }: Props = $props();

  let name = $state(untrack(() => tag?.name ?? ''));
  let focusInput: any;

  const title = $derived(tag ? 'Edit retro tag' : 'Create retro tag');

  function onSubmit(event: Event) {
    event.preventDefault();

    xfetch(tag ? `${apiPrefix}/retro-tags/${tag.id}` : `${apiPrefix}/retro-tags`, {
      method: tag ? 'PUT' : 'POST',
      body: { name },
    })
      .then(() => {
        notifications.success('Retro tag successfully saved');
        handleSaved();
      })
      .catch(error => {
        // apiclient rejects with [error, response]
        if (error?.[1]?.status === 409) {
          notifications.danger('A retro tag with that name already exists');
        } else {
          notifications.danger('Error saving retro tag');
        }
      });
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleClose} ariaLabel={title}>
  <form onsubmit={onSubmit} name="retroTag" class="space-y-4">
    <div>
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="retroTagName">Name</label>
      <TextInput
        bind:value={name}
        bind:this={focusInput}
        id="retroTagName"
        name="retroTagName"
        placeholder="Enter tag name"
        maxlength={RETRO_TAG_MAX_LENGTH}
        required
      />
    </div>

    <div class="text-right">
      <SolidButton type="submit" disabled={name.trim() === ''}>Save Tag</SolidButton>
    </div>
  </form>
</Modal>

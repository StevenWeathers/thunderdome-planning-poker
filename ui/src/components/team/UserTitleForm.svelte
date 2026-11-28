<script lang="ts">
  import { onMount, untrack } from 'svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { UserTitle } from '../../types/team';

  interface Props {
    userTitle?: UserTitle | null;
    apiPrefix: string;
    toggleClose?: () => void;
    handleSaved?: () => void;
    xfetch: ApiClient;
    notifications: NotificationService;
  }

  let {
    userTitle = null,
    apiPrefix,
    toggleClose = () => {},
    handleSaved = () => {},
    xfetch,
    notifications,
  }: Props = $props();

  let name = $state(untrack(() => userTitle?.name ?? ''));
  let focusInput: any;

  const heading = $derived(userTitle ? $LL.userTitleEdit() : $LL.userTitleAdd());

  function onSubmit(event: Event) {
    event.preventDefault();

    xfetch(userTitle ? `${apiPrefix}/user-titles/${userTitle.id}` : `${apiPrefix}/user-titles`, {
      method: userTitle ? 'PUT' : 'POST',
      body: { name },
    })
      .then(() => {
        notifications.success($LL.userTitleSaveSuccess());
        handleSaved();
      })
      .catch(error => {
        // apiclient rejects with [error, response]
        if (error?.[1]?.status === 409) {
          notifications.danger($LL.userTitleDuplicateError());
        } else {
          notifications.danger($LL.userTitleSaveError());
        }
      });
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleClose} ariaLabel={heading}>
  <form onsubmit={onSubmit} name="userTitleForm" class="space-y-4">
    <div>
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="userTitleName">
        {$LL.userTitleLabel()}
      </label>
      <TextInput
        bind:value={name}
        bind:this={focusInput}
        id="userTitleName"
        name="userTitleName"
        placeholder={$LL.userTitlePlaceholder()}
        maxlength={128}
        autocomplete="off"
        required
      />
    </div>

    <div class="text-right">
      <SolidButton type="submit" disabled={name.trim() === ''} testid="user-title-save">
        {$LL.userTitleSave()}
      </SolidButton>
    </div>
  </form>
</Modal>

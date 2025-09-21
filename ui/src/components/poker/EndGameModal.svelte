<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { ChevronDown, Crown, Lock } from 'lucide-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  const allowedReasons = [
    {
      value: 'Completed',
      label: $LL.completed(),
    },
    {
      value: 'Cancelled',
      label: $LL.cancelled(),
    },
  ];

  interface Props {
    toggleModal?: any;
    handleSubmit?: any;
    notifications: NotificationService;
    xfetch: ApiClient;
  }

  let { toggleModal = () => {}, handleSubmit = (battle: any) => {}, notifications, xfetch }: Props = $props();

  let endGameReason = $state('');
  let submitDisabled = $derived(() => {
    return endGameReason === '';
  });

  let focusInput: any;
  onMount(() => {
    focusInput?.focus();
  });

  const onSubmit = (e: Event) => {
    e.preventDefault();

    handleSubmit({ endGameReason });
  };
</script>

<Modal closeModal={toggleModal} ariaLabel={$LL.modalEditPokerGame()}>
  <form onsubmit={onSubmit} name="createBattle">
    <div class="text-lg font-bold mb-4 dark:text-white">{$LL.endGame()}</div>
    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2" for="reason">
        {$LL.endGameReason()}
      </label>
      <SelectInput bind:this={focusInput} bind:value={endGameReason} id="endGameReason" name="endGameReason">
        <option value="" disabled>{$LL.endGameReasonPlaceholder()}</option>
        {#each allowedReasons as reason}
          <option value={reason.value}>
            {reason.label}
          </option>
        {/each}
      </SelectInput>
    </div>

    <div class="text-right">
      <SolidButton type="submit" disabled={submitDisabled()}>{$LL.endGame()}</SolidButton>
    </div>
  </form>
</Modal>

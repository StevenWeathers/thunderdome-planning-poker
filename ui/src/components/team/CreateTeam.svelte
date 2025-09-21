<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { onMount } from 'svelte';

  interface Props {
    teamName?: string;
    toggleCreate?: any;
    handleCreate?: any;
  }

  let { teamName = $bindable(''), toggleCreate = () => {}, handleCreate = () => {} }: Props = $props();

  function onSubmit(e) {
    e.preventDefault();

    handleCreate(teamName);
  }

  let createDisabled = $derived(teamName === '');

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleCreate} ariaLabel={$LL.modalCreateTeam()}>
  <form onsubmit={onSubmit} name="createTeam">
    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="teamName">
        {$LL.teamName()}
      </label>
      <TextInput
        bind:value={teamName}
        bind:this={focusInput}
        placeholder={$LL.teamNamePlaceholder()}
        id="teamName"
        name="teamName"
        required
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {$LL.teamSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

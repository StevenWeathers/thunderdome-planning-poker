<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../global/TextInput.svelte';

  export let teamName = '';

  export let toggleCreate = () => {};
  export let handleCreate = () => {};

  function onSubmit(e) {
    e.preventDefault();

    handleCreate(teamName);
  }

  $: createDisabled = teamName === '';
</script>

<Modal closeModal="{toggleCreate}">
  <form on:submit="{onSubmit}" name="createTeam">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="teamName"
      >
        {$LL.teamName()}
      </label>
      <TextInput
        bind:value="{teamName}"
        placeholder="{$LL.teamNamePlaceholder()}"
        id="teamName"
        name="teamName"
        required
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled="{createDisabled}">
          {$LL.teamSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

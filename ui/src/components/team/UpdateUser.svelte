<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../global/TextInput.svelte';
  import SelectInput from '../global/SelectInput.svelte';

  export let toggleUpdate = () => {};
  export let handleUpdate = () => {};
  export let userId = '';
  export let userEmail = '';
  export let role = '';

  const roles = ['ADMIN', 'MEMBER'];

  function onSubmit(e) {
    e.preventDefault();

    handleUpdate(userId, role);
  }

  $: updateDisabled = role === '';
</script>

<Modal closeModal="{toggleUpdate}">
  <form on:submit="{onSubmit}" name="teamUpdateUser">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2 disabled"
        for="userEmail"
      >
        {$LL.userEmail()}
      </label>
      <TextInput value="{userEmail}" id="userEmail" name="userEmail" disabled />
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="userRole"
      >
        {$LL.role()}
      </label>
      <SelectInput bind:value="{role}" id="userRole" name="userRole">
        {#each roles as userRole}
          <option value="{userRole}">{userRole}</option>
        {/each}
      </SelectInput>
    </div>

    <div>
      <div class="text-right">
        <SolidButton
          type="submit"
          disabled="{updateDisabled}"
          testid="userupdate-confirm"
        >
          {$LL.userUpdate()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

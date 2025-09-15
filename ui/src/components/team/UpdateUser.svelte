<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { onMount } from 'svelte';

  interface Props {
    toggleUpdate?: any;
    handleUpdate?: any;
    userId?: string;
    userEmail?: string;
    role?: string;
  }

  let {
    toggleUpdate = () => {},
    handleUpdate = () => {},
    userId = '',
    userEmail = '',
    role = $bindable('')
  }: Props = $props();

  const roles = ['ADMIN', 'MEMBER'];

  function onSubmit(e) {
    e.preventDefault();

    handleUpdate(userId, role);
  }

  let updateDisabled = $derived(role === '');

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleUpdate} ariaLabel={$LL.modalTeamUpdateUser()}>
  <form onsubmit={onSubmit} name="teamUpdateUser">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2 disabled"
        for="userEmail"
      >
        {$LL.userEmail()}
      </label>
      <TextInput bind:this={focusInput} value={userEmail} id="userEmail" name="userEmail" disabled />
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
          disabled={updateDisabled}
          testid="userupdate-confirm"
        >
          {$LL.userUpdate()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

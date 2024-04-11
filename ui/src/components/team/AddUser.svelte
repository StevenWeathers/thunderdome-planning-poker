<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { AppConfig } from '../../config';

  export let toggleAdd = () => {};
  export let handleAdd = () => {};
  export let pageType = '';
  export let requiresOrgMember = false;

  const roles = ['ADMIN', 'MEMBER'];
  let userEmail = '';
  let role = '';

  function onSubmit(e) {
    e.preventDefault();

    handleAdd(userEmail, role);
  }

  $: createDisabled = userEmail === '' || role === '';
</script>

<Modal closeModal="{toggleAdd}">
  <form on:submit="{onSubmit}" name="teamAddUser">
    <div class="mb-2">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="userEmail"
      >
        {$LL.userEmail()}
      </label>
      <TextInput
        bind:value="{userEmail}"
        placeholder="{$LL.userEmailPlaceholder()}"
        id="userEmail"
        name="userEmail"
        required
      />
    </div>

    <div class="mb-4 text-gray-700 dark:text-gray-400 text-sm">
      {#if AppConfig.LdapEnabled || AppConfig.HeaderAuthEnabled}
        {$LL.addUserWillInviteNotFoundFieldNote({ pageType })}
      {:else if !requiresOrgMember}
        {$LL.inviteUserFieldNote({ pageType })}
      {:else}
        {$LL.requiresOrganizationMember()}
      {/if}
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="userRole"
      >
        {$LL.role()}
      </label>
      <SelectInput bind:value="{role}" id="userRole" name="userRole">
        <option value="">{$LL.rolePlaceholder()}</option>
        {#each roles as userRole}
          <option value="{userRole}">{userRole}</option>
        {/each}
      </SelectInput>
    </div>

    <div>
      <div class="text-right">
        <SolidButton
          type="submit"
          disabled="{createDisabled}"
          testid="useradd-confirm"
        >
          {$LL.userAdd()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

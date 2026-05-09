<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import MultiSelectInput from '../forms/MultiSelectInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { AppConfig } from '../../config';
  import { onMount } from 'svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { AddUserRequest, InviteUserRequest } from '../../types/team';

  interface Props {
    toggleAdd?: any;
    handleAdd?: any;
    handleInvite?: any;
    pageType?: string;
    orgId?: string;
    deptId?: string;
    notifications: NotificationService;
    xfetch: ApiClient;
  }

  let {
    toggleAdd = () => {},
    handleAdd = () => {},
    handleInvite = () => {},
    pageType = '',
    orgId = '',
    deptId = '',
    notifications,
    xfetch,
  }: Props = $props();

  const roles = ['ADMIN', 'MEMBER'];
  const showDeptUsers = $derived(pageType === 'team' && deptId !== '');
  const showOrgUsers = $derived(pageType === 'department' || (pageType === 'team' && orgId !== ''));
  let userEmail = $state('');
  let selectedOrgUsers = $state<string[]>([]);
  let selectedDeptUsers = $state<string[]>([]);
  let role = $state('');

  let organizationUsers = $state<any[]>([]);
  let departmentUsers = $state<any[]>([]);

  const fallbackEmailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const selectedUserIds = $derived([...new Set([...selectedOrgUsers, ...selectedDeptUsers])]);
  const inviteEmailTokens = $derived(
    userEmail
      .split(/[\n,;]+/)
      .map(email => email.trim())
      .filter(Boolean),
  );
  const invalidInviteEmails = $derived.by(() => [...new Set(inviteEmailTokens.filter(email => !isValidEmail(email)))]);
  const inviteEmails = $derived.by(() => [...new Set(inviteEmailTokens.filter(email => isValidEmail(email)))]);

  function isValidEmail(email: string) {
    if (typeof document === 'undefined') {
      return fallbackEmailPattern.test(email);
    }

    const input = document.createElement('input');
    input.type = 'email';
    input.value = email;

    return input.validity.valid;
  }

  function getOrganizationUsers() {
    if (orgId !== '') {
      xfetch(`/api/organizations/${orgId}/users?limit=${9999}&offset=${0}`)
        .then(res => res.json())
        .then(function (result) {
          organizationUsers = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetUsersError());
        });
    }
  }

  function getDepartmentUsers() {
    if (orgId !== '' && deptId !== '') {
      xfetch(`/api/organizations/${orgId}/departments/${deptId}/users?limit=${9999}&offset=${0}`)
        .then(res => res.json())
        .then(function (result) {
          departmentUsers = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetUsersError());
        });
    }
  }

  function clearInviteEmails() {
    userEmail = '';
  }

  function clearSelectedUsers() {
    selectedOrgUsers = [];
    selectedDeptUsers = [];
  }

  function onInviteInput() {
    if (userEmail.trim() !== '') {
      clearSelectedUsers();
    }
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    if (selectedUserIds.length === 0 && invalidInviteEmails.length > 0) {
      focusInput?.reportValidity();
      return;
    }

    if (selectedUserIds.length > 0) {
      const usersToAdd: AddUserRequest[] = selectedUserIds.map(userId => ({
        user_id: userId,
        role,
      }));
      handleAdd(usersToAdd);
    } else {
      const usersToAdd: InviteUserRequest[] = inviteEmails.map(email => ({
        email,
        role,
      }));
      handleInvite(usersToAdd);
    }
  }

  let createDisabled = $derived((inviteEmails.length === 0 && selectedUserIds.length === 0) || role === '');

  let focusInput: HTMLTextAreaElement | undefined = $state();

  $effect(() => {
    if (!focusInput) {
      return;
    }

    focusInput.setCustomValidity(invalidInviteEmails.length > 0 ? 'Please enter valid email addresses only.' : '');
  });

  onMount(() => {
    getOrganizationUsers();
    getDepartmentUsers();
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleAdd} ariaLabel={$LL.modalTeamAddUser()}>
  <form onsubmit={onSubmit} name="teamAddUser">
    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 font-bold mb-2" for="userRole">
        {$LL.role()}
      </label>
      <SelectInput bind:value={role} id="userRole" name="userRole">
        <option value="">{$LL.rolePlaceholder()}</option>
        {#each roles as userRole}
          <option value={userRole}>{userRole}</option>
        {/each}
      </SelectInput>
    </div>

    {#if showOrgUsers}
      <div class="mb-4">
        <label class="text-gray-700 dark:text-gray-400 font-bold mb-2" for="orgUser">
          Select existing users from Organization
        </label>
        <MultiSelectInput
          bind:value={selectedOrgUsers}
          onchange={clearInviteEmails}
          id="orgUser"
          name="orgUser"
          size="6"
        >
          {#each organizationUsers as orgUser}
            <option value={orgUser.id}>{orgUser.name} ({orgUser.email}) </option>
          {/each}
        </MultiSelectInput>
        <div class="mt-2 text-sm text-gray-700 dark:text-gray-400">Select one or more users to add.</div>
      </div>
    {/if}
    {#if showDeptUsers}
      <div class="mb-4">
        <label class="text-gray-700 dark:text-gray-400 font-bold mb-2" for="deptUser">
          Select existing users from Department
        </label>
        <MultiSelectInput
          bind:value={selectedDeptUsers}
          onchange={clearInviteEmails}
          id="deptUser"
          name="deptUser"
          size="6"
        >
          {#each departmentUsers as deptUser}
            <option value={deptUser.id}>{deptUser.name} ({deptUser.email}) </option>
          {/each}
        </MultiSelectInput>
        <div class="mt-2 text-sm text-gray-700 dark:text-gray-400">Select one or more users to add.</div>
      </div>
    {/if}

    {#if showOrgUsers || showDeptUsers}
      <div class="relative">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-gray-300"></div>
        </div>
        <div class="relative flex justify-center text-sm">
          <span class="bg-white dark:bg-gray-800 px-2 text-gray-500 dark:text-white">Or invite a user by email</span>
        </div>
      </div>
    {/if}

    <div class="mb-2">
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="userEmail">
        {$LL.userEmail()}
      </label>
      <textarea
        bind:this={focusInput}
        bind:value={userEmail}
        placeholder={$LL.userEmailPlaceholder()}
        id="userEmail"
        name="userEmail"
        required={selectedUserIds.length === 0}
        disabled={selectedUserIds.length !== 0}
        aria-invalid={invalidInviteEmails.length > 0}
        rows="4"
        class="block w-full rounded-lg outline-none transition-all duration-300 bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-400 disabled:cursor-not-allowed px-5 py-3"
        oninput={onInviteInput}
      ></textarea>
      <div class="mt-2 text-sm text-gray-700 dark:text-gray-400">
        Enter one email per line, or separate multiple emails with commas.
      </div>
      {#if invalidInviteEmails.length > 0}
        <div class="mt-2 text-sm text-red-600 dark:text-red-400" data-testid="invalid-email-message">
          Invalid email{invalidInviteEmails.length === 1 ? '' : 's'}: {invalidInviteEmails.join(', ')}
        </div>
      {/if}
    </div>

    <div class="mb-4 text-gray-700 dark:text-gray-400 text-sm">
      {#if AppConfig.LdapEnabled || AppConfig.HeaderAuthEnabled}
        {$LL.addUserWillInviteNotFoundFieldNote({ pageType })}
      {:else}
        {$LL.inviteUserFieldNote({ pageType })}
      {/if}
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled} testid="useradd-confirm">
          {$LL.userAdd()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

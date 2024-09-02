<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { AppConfig } from '../../config';
  import { onMount } from 'svelte';

  export let toggleAdd = () => {};
  export let handleAdd = () => {};
  export let handleInvite = () => {};
  export let pageType = '';
  export let orgId = '';
  export let deptId = '';
  export let notifications;
  export let xfetch;

  const roles = ['ADMIN', 'MEMBER'];
  const showDeptUsers = pageType === 'team' && deptId !== '';
  const showOrgUsers =
    pageType === 'department' || (pageType === 'team' && orgId !== '');
  let userEmail = '';
  let selectedUser = '';
  let role = '';

  let organizationUsers = [];
  let departmentUsers = [];

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
      xfetch(
        `/api/organizations/${orgId}/departments/${deptId}/users?limit=${9999}&offset=${0}`,
      )
        .then(res => res.json())
        .then(function (result) {
          departmentUsers = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetUsersError());
        });
    }
  }

  function clearEmail() {
    userEmail = '';
  }

  function onSubmit(e) {
    e.preventDefault();

    const user = {
      id: selectedUser,
      email: userEmail,
      role,
    };

    if (selectedUser !== '') {
      handleAdd(user);
    } else {
      handleInvite(user);
    }
  }

  $: createDisabled = (userEmail === '' && selectedUser === '') || role === '';

  onMount(() => {
    getOrganizationUsers();
    getDepartmentUsers();
  });
</script>

<Modal closeModal="{toggleAdd}">
  <form on:submit="{onSubmit}" name="teamAddUser">
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

    {#if showOrgUsers}
      <div class="mb-4">
        <label
          class="text-gray-700 dark:text-gray-400 font-bold mb-2"
          for="orgUser"
        >
          Select an existing user from Organization
        </label>
        <SelectInput
          bind:value="{selectedUser}"
          on:change="{clearEmail}"
          id="orgUser"
          name="orgUser"
        >
          <option value="">Organization Users...</option>
          {#each organizationUsers as orgUser}
            <option value="{orgUser.id}"
              >{orgUser.name} ({orgUser.email})
            </option>
          {/each}
        </SelectInput>
      </div>
    {/if}
    {#if showDeptUsers}
      <div class="mb-4">
        <label
          class="text-gray-700 dark:text-gray-400 font-bold mb-2"
          for="deptUser"
        >
          Select an existing user from Department
        </label>
        <SelectInput
          bind:value="{selectedUser}"
          on:change="{clearEmail}"
          id="deptUser"
          name="deptUser"
        >
          <option value="">Department Users...</option>
          {#each departmentUsers as deptUser}
            <option value="{deptUser.id}"
              >{deptUser.name} ({deptUser.email})
            </option>
          {/each}
        </SelectInput>
      </div>
    {/if}

    {#if showOrgUsers || showDeptUsers}
      <div class="relative">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-gray-300"></div>
        </div>
        <div class="relative flex justify-center text-sm">
          <span
            class="bg-white dark:bg-gray-800 px-2 text-gray-500 dark:text-white"
            >Or invite a user by email</span
          >
        </div>
      </div>
    {/if}

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
        required="{selectedUser === ''}"
        disabled="{selectedUser !== ''}"
      />
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

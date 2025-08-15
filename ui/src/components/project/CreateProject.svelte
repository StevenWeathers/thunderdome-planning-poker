<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { onMount } from 'svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    apiPrefix?: string;
    toggleCreate?: any;
    handleCreate?: any;
    toggleUpdate?: any;
    handleUpdate?: any;
    projectId?: string;
    projectKey?: string;
    projectName?: string;
    description?: string;
    organizationId?: string;
    departmentId?: string;
    teamId?: string;
  }

  let {
    xfetch,
    notifications,
    apiPrefix = '/api',
    toggleCreate = () => {},
    handleCreate = () => {},
    toggleUpdate = () => {},
    handleUpdate = () => {},
    projectId = '',
    projectKey = $bindable(''),
    projectName = $bindable(''),
    description = $bindable(''),
    organizationId = '',
    departmentId = '',
    teamId = ''
  }: Props = $props();

  function toggleClose() {
    if (projectId != '') {
      toggleUpdate();
    } else {
      toggleCreate();
    }
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    const body: any = {
      projectKey: projectKey.toUpperCase(),
      name: projectName,
    };

    // Only include description if it has a value
    if (description && description.trim() !== '') {
      body.description = description.trim();
    }

    // For admin endpoint, include association IDs
    if (!teamId && !departmentId && !organizationId) {
      if (organizationId) body.organizationId = organizationId;
      if (departmentId) body.departmentId = departmentId;
      if (teamId) body.teamId = teamId;
    }

    const url = projectId !== '' 
      ? `${apiPrefix}/projects/${projectId}` 
      : `${apiPrefix}/projects`;
    
    const method = projectId !== '' ? 'PUT' : 'POST';

    xfetch(url, {
      method,
      body,
    })
      .then(res => {
        if (!res.ok) {
          throw new Error(`HTTP error! status: ${res.status}`);
        }
        return res.json();
      })
      .then(function (result) {
        if (projectId !== '') {
          notifications.success($LL.projectUpdateSuccess());
          handleUpdate();
        } else {
          notifications.success($LL.projectCreateSuccess());
          handleCreate();
        }
      })
      .catch(function (error) {
        console.error('Project API error:', error);
        if (projectId !== '') {
          notifications.danger($LL.projectUpdateError());
        } else {
          notifications.danger($LL.projectCreateError());
        }
      });
  }

  // Validation for project key format (2-10 chars, alphanumeric uppercase)
  let projectKeyValid = $derived(() => {
    const keyRegex = /^[A-Z0-9]{2,10}$/;
    return projectKey === '' || keyRegex.test(projectKey.toUpperCase());
  });

  let createDisabled = $derived(
    projectName === '' || 
    projectKey === '' || 
    !projectKeyValid() ||
    projectKey.length < 2 || 
    projectKey.length > 10
  );

  function handleProjectKeyInput(event: Event) {
    const target = event.target as HTMLInputElement;
    projectKey = target.value.toUpperCase();
    }

  let focusInput: any;
  onMount(() => {
    focusInput?.focus();
  });

  // Auto-format project key to uppercase
  $effect(() => {
    if (projectKey) {
      projectKey = projectKey.toUpperCase();
    }
  });
</script>

<Modal closeModal={toggleClose} ariaLabel={projectId ? $LL.modalUpdateProject() : $LL.modalCreateProject()}>
  <form onsubmit={onSubmit} name="createProject">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="projectKey"
      >
        {$LL.projectKey()}
      </label>
      <TextInput
        bind:value={projectKey}
        bind:this={focusInput}
        placeholder={$LL.projectKeyPlaceholder()}
        id="projectKey"
        name="projectKey"
        required
        maxlength="10"
        pattern={`[A-Z0-9]{2,10}`}
        oninput={handleProjectKeyInput}
        class={!projectKeyValid() ? 'border-red-500' : ''}
      />
      {#if !projectKeyValid() && projectKey !== ''}
        <p class="text-red-500 text-sm mt-1">
          {$LL.projectKeyValidationError()}
        </p>
      {/if}
      <p class="text-gray-500 text-sm mt-1">
        {$LL.projectKeyHelp()}
      </p>
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="projectName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value={projectName}
        placeholder={$LL.projectNamePlaceholder()}
        id="projectName"
        name="projectName"
        required
        maxlength="255"
      />
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="projectDescription"
      >
        {$LL.description()}
      </label>
      <textarea
        bind:value={description}
        placeholder={$LL.projectDescriptionPlaceholder()}
        id="projectDescription"
        name="projectDescription"
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:bg-gray-700 dark:border-gray-600 dark:text-white dark:placeholder-gray-400"
      ></textarea>
    </div>

    <!-- Show scope information -->
    {#if teamId || departmentId || organizationId}
      <div class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
        <p class="text-sm text-blue-700 dark:text-blue-300">
          <strong>{$LL.projectScope()}:</strong>
          {#if teamId}
            {$LL.team()}
          {:else if departmentId}
            {$LL.department()}
          {:else if organizationId}
            {$LL.organization()}
          {/if}
        </p>
      </div>
    {/if}

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {projectId ? $LL.projectUpdate() : $LL.projectCreate()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>
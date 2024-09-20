<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import { RetroTemplateFormat } from '../../types/retro';
  import ColumnForm from './ColumnForm.svelte';

  export let toggleCreate = () => {};
  export let handleCreate = () => {};
  export let organizationId;
  export let teamId;
  export let departmentId;
  export let apiPrefix: string = '/api';
  export let xfetch: any;
  export let eventTag: any;
  export let notifications: any;

  let name = '';
  let description = '';
  let isPublic = false;
  let defaultTemplate = false;
  let format: RetroTemplateFormat = {
    columns: [],
  };

  function toggleClose() {
    toggleCreate();
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    const body = {
      name,
      description,
      defaultTemplate,
      format,
    };

    if (isAdmin) {
      body['isPublic'] = isPublic;
    }

    xfetch(`${apiPrefix}/retro-templates`, {
      body,
    })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.retroTemplateCreateSuccess());
        handleCreate();
      })
      .catch(() => {
        notifications.danger($LL.retroTemplateCreateError());
      });
  }

  $: createDisabled =
    name === '' || format.columns.length < 2 || format.columns.length > 5;
  $: isAdmin = validateUserIsAdmin($user);
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{onSubmit}" name="createRetroTemplate">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="templateName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value="{name}"
        placeholder="{$LL.retroTemplateNamePlaceholder()}"
        id="templateName"
        name="templateName"
        required
      />
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="templateDescription"
      >
        {$LL.description()}
      </label>
      <TextInput
        bind:value="{description}"
        placeholder="{$LL.retroTemplateDescriptionPlaceholder()}"
        id="templateDescription"
        name="templateDescription"
      />
    </div>

    <ColumnForm bind:format="{format}" />

    {#if isAdmin && !organizationId && !teamId}
      <div class="mb-4">
        <Checkbox
          bind:checked="{isPublic}"
          id="isPublic"
          name="isPublic"
          label="{$LL.retroTemplateIsPublic()}"
        />
      </div>
    {/if}

    <div class="mb-4">
      <Checkbox
        bind:checked="{defaultTemplate}"
        id="defaultTemplate"
        name="defaultTemplate"
        label="{$LL.retroTemplateDefault()}"
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled="{createDisabled}">
          {$LL.retroTemplateSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

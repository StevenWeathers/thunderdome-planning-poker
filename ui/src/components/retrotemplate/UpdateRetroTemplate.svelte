<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import ColumnForm from './ColumnForm.svelte';

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    toggleUpdate?: any;
    handleUpdate?: any;
    organizationId: any;
    teamId: any;
    departmentId: any;
    apiPrefix: any;
    xfetch: any;
    notifications: NotificationService;
    templateId?: string;
    name?: string;
    description?: string;
    format: any;
    isPublic?: boolean;
    defaultTemplate?: boolean;
  }

  let {
    toggleUpdate = () => {},
    handleUpdate = () => {},
    organizationId,
    teamId,
    departmentId,
    apiPrefix,
    xfetch,
    notifications,
    templateId = '',
    name = $bindable(''),
    description = $bindable(''),
    format = $bindable(),
    isPublic = $bindable(false),
    defaultTemplate = $bindable(false)
  }: Props = $props();

  function toggleClose() {
    toggleUpdate();
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

    xfetch(`${apiPrefix}/retro-templates/${templateId}`, {
      body,
      method: 'put',
    })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.retroTemplateUpdateSuccess());
        handleUpdate();
      })
      .catch(() => {
        notifications.danger($LL.retroTemplateUpdateError());
      });
  }

  let updateDisabled =
    $derived(name === '' || format.columns.length < 2 || format.columns.length > 5);
  let isAdmin = $derived(validateUserIsAdmin($user));
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={onSubmit} name="updateRetroTemplate">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="templateName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value="{name}"
        placeholder={$LL.retroTemplateNamePlaceholder()}
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
        placeholder={$LL.retroTemplateDescriptionPlaceholder()}
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
          label={$LL.retroTemplateIsPublic()}
        />
      </div>
    {/if}

    <div class="mb-4">
      <Checkbox
        bind:checked="{defaultTemplate}"
        id="defaultTemplate"
        name="defaultTemplate"
        label={$LL.retroTemplateDefault()}
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={updateDisabled}>
          {$LL.retroTemplateSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

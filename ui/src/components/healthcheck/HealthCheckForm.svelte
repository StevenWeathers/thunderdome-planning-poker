<script lang="ts">
  import { onMount, untrack } from 'svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { HealthCheck, HealthCheckTemplate } from '../../types/healthcheck';

  interface Props {
    // the health check being edited, null to create one
    healthCheck?: HealthCheck | null;
    apiPrefix: string;
    toggleClose?: () => void;
    handleSaved?: (healthCheck: HealthCheck | null) => void;
    xfetch: ApiClient;
    notifications: NotificationService;
  }

  let {
    healthCheck = null,
    apiPrefix,
    toggleClose = () => {},
    handleSaved = () => {},
    xfetch,
    notifications,
  }: Props = $props();

  // datetime-local inputs take local time without a timezone
  const toInputValue = (date: Date) =>
    new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);

  const defaultDeadline = () => new Date(Date.now() + 7 * 24 * 60 * 60 * 1000);

  let name = $state(untrack(() => healthCheck?.name ?? ''));
  let deadline = $state(
    untrack(() => toInputValue(healthCheck ? new Date(healthCheck.deadline) : defaultDeadline())),
  );
  let templates = $state<HealthCheckTemplate[]>([]);
  let templateId = $state('');
  let sendEmail = $state(false);
  let saving = $state(false);

  const title = $derived(healthCheck ? $LL.healthCheckEdit() : $LL.healthCheckCreate());
  const canSave = $derived(!saving && name.trim() !== '' && deadline !== '' && (!!healthCheck || templateId !== ''));

  function onSubmit(event: Event) {
    event.preventDefault();
    if (!canSave) return;

    const deadlineDate = new Date(deadline);
    if (deadlineDate.getTime() <= Date.now()) {
      notifications.danger($LL.healthCheckDeadlineFuture());
      return;
    }

    const body = healthCheck
      ? { name, deadline: deadlineDate.toISOString() }
      : { name, deadline: deadlineDate.toISOString(), templateId, sendEmail };

    saving = true;
    xfetch(healthCheck ? `${apiPrefix}/health-checks/${healthCheck.id}` : `${apiPrefix}/health-checks`, {
      method: healthCheck ? 'PUT' : 'POST',
      body,
    })
      .then(res => res.json())
      .then(result => {
        notifications.success($LL.healthCheckSaveSuccess());
        handleSaved(healthCheck ? null : result.data);
      })
      .catch(error => {
        saving = false;
        // apiclient rejects with [error, response]
        if (error?.[1]?.status === 409) {
          notifications.danger($LL.healthCheckClosedError());
        } else {
          notifications.danger($LL.healthCheckSaveError());
        }
      });
  }

  onMount(() => {
    if (healthCheck) return;

    xfetch(`${apiPrefix}/health-check-templates/available`)
      .then(res => res.json())
      .then(result => {
        templates = result.data;
        templateId = templates[0]?.id ?? '';
      })
      .catch(() => {
        notifications.danger($LL.healthCheckTemplatesGetError());
      });
  });

  const labelClass = 'block text-gray-700 font-bold mb-2 dark:text-gray-400';
</script>

<Modal closeModal={toggleClose} ariaLabel={title}>
  <form onsubmit={onSubmit} name="healthCheck" class="space-y-4 pt-6">
    <div>
      <label class={labelClass} for="healthCheckName">{$LL.healthCheckNameLabel()}</label>
      <TextInput
        bind:value={name}
        id="healthCheckName"
        name="healthCheckName"
        placeholder={$LL.healthCheckNamePlaceholder()}
        maxlength={255}
        autocomplete="off"
        required
      />
    </div>

    {#if !healthCheck}
      <div>
        <label class={labelClass} for="healthCheckTemplate">{$LL.healthCheckTemplate()}</label>
        <SelectInput bind:value={templateId} id="healthCheckTemplate" name="healthCheckTemplate" required>
          {#each templates as template (template.id)}
            <option value={template.id}>{template.name}</option>
          {/each}
        </SelectInput>
      </div>
    {/if}

    <div>
      <label class={labelClass} for="healthCheckDeadline">{$LL.healthCheckDeadline()}</label>
      <TextInput
        bind:value={deadline}
        type="datetime-local"
        id="healthCheckDeadline"
        name="healthCheckDeadline"
        required
      />
    </div>

    {#if !healthCheck}
      <Checkbox
        bind:checked={sendEmail}
        id="healthCheckSendEmail"
        name="healthCheckSendEmail"
        label={$LL.healthCheckSendEmail()}
      />
    {/if}

    <div class="text-right">
      <SolidButton type="submit" disabled={!canSave} testid="health-check-save">
        {$LL.healthCheckSave()}
      </SolidButton>
    </div>
  </form>
</Modal>

<script lang="ts">
  import { untrack } from 'svelte';
  import { Plus, Trash2 } from '@lucide/svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { HealthCheckArea, HealthCheckTemplate } from '../../types/healthcheck';

  interface Props {
    template?: HealthCheckTemplate | null;
    toggleClose?: () => void;
    handleSaved?: () => void;
    apiPrefix: string;
    xfetch: ApiClient;
    notifications: NotificationService;
  }

  let {
    template = null,
    toggleClose = () => {},
    handleSaved = () => {},
    apiPrefix,
    xfetch,
    notifications,
  }: Props = $props();

  const MAX_AREAS = 20;

  let nextKey = 0;
  const newArea = (area?: HealthCheckArea) => ({
    key: nextKey++,
    name: area?.name ?? '',
    healthy: area?.healthy ?? '',
    unhealthy: area?.unhealthy ?? '',
  });

  let name = $state(untrack(() => template?.name ?? ''));
  let description = $state(untrack(() => template?.description ?? ''));
  let areas = $state(untrack(() => (template ? template.areas.map(a => newArea(a)) : [newArea()])));

  const title = $derived(template ? $LL.healthCheckTemplateEdit() : $LL.healthCheckTemplateCreate());
  const names = $derived(areas.map(a => a.name.trim()));
  const hasDuplicates = $derived(new Set(names.filter(Boolean)).size !== names.filter(Boolean).length);
  const canSave = $derived(name.trim() !== '' && names.every(Boolean) && !hasDuplicates);

  function onSubmit(event: Event) {
    event.preventDefault();
    if (!canSave) return;

    xfetch(
      template ? `${apiPrefix}/health-check-templates/${template.id}` : `${apiPrefix}/health-check-templates`,
      {
        method: template ? 'PUT' : 'POST',
        body: {
          name,
          description,
          areas: areas.map(({ name, healthy, unhealthy }) => ({ name, healthy, unhealthy })),
        },
      },
    )
      .then(() => {
        notifications.success($LL.healthCheckTemplateSaveSuccess());
        handleSaved();
      })
      .catch(() => {
        notifications.danger($LL.healthCheckTemplateSaveError());
      });
  }

  const textareaClass =
    'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-gray-100 resize-none';
  const labelClass = 'block text-gray-700 font-bold mb-2 dark:text-gray-400';
</script>

<Modal closeModal={toggleClose} ariaLabel={title}>
  <form onsubmit={onSubmit} name="healthCheckTemplate" class="space-y-4 pt-6">
    <div>
      <label class={labelClass} for="healthCheckTemplateName">{$LL.healthCheckTemplateName()}</label>
      <TextInput
        bind:value={name}
        id="healthCheckTemplateName"
        name="healthCheckTemplateName"
        maxlength={255}
        autocomplete="off"
        required
      />
    </div>
    <div>
      <label class={labelClass} for="healthCheckTemplateDescription">{$LL.healthCheckTemplateDescription()}</label>
      <textarea
        id="healthCheckTemplateDescription"
        bind:value={description}
        maxlength="2000"
        rows="2"
        class={textareaClass}
      ></textarea>
    </div>

    <h3 class="font-bold text-gray-700 dark:text-gray-400">{$LL.healthCheckAreas()}</h3>
    {#each areas as area, i (area.key)}
      <fieldset class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2">
        <legend class="sr-only">{$LL.healthCheckAreaLabel({ number: i + 1 })}</legend>
        <div class="flex items-center gap-2">
          <TextInput
            bind:value={area.name}
            aria-label="{$LL.healthCheckAreaLabel({ number: i + 1 })}: {$LL.healthCheckAreaName()}"
            placeholder={$LL.healthCheckAreaName()}
            maxlength={128}
            required
          />
          <button
            type="button"
            class="p-1.5 text-gray-500 hover:text-red-600 dark:hover:text-red-400 disabled:opacity-30"
            aria-label={$LL.healthCheckAreaRemove({ number: i + 1 })}
            disabled={areas.length === 1}
            onclick={() => areas.splice(i, 1)}
          >
            <Trash2 class="w-5 h-5" />
          </button>
        </div>
        <textarea
          bind:value={area.healthy}
          aria-label="{$LL.healthCheckAreaLabel({ number: i + 1 })}: {$LL.healthCheckHealthy()}"
          placeholder={$LL.healthCheckHealthy()}
          maxlength="500"
          rows="2"
          class={textareaClass}
        ></textarea>
        <textarea
          bind:value={area.unhealthy}
          aria-label="{$LL.healthCheckAreaLabel({ number: i + 1 })}: {$LL.healthCheckUnhealthy()}"
          placeholder={$LL.healthCheckUnhealthy()}
          maxlength="500"
          rows="2"
          class={textareaClass}
        ></textarea>
      </fieldset>
    {/each}

    {#if hasDuplicates}
      <p class="text-red-600 dark:text-red-400" role="alert">{$LL.healthCheckAreaNamesUnique()}</p>
    {/if}

    <div class="flex justify-between items-center">
      <HollowButton onClick={() => areas.push(newArea())} disabled={areas.length >= MAX_AREAS}>
        <Plus class="w-4 h-4 me-1 inline" />
        {$LL.healthCheckAreaAdd()}
      </HollowButton>
      <SolidButton type="submit" disabled={!canSave} testid="health-check-template-save">
        {$LL.healthCheckTemplateSave()}
      </SolidButton>
    </div>
  </form>
</Modal>

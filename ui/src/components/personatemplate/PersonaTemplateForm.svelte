<script lang="ts">
  import { untrack } from 'svelte';
  import { Plus, Trash2 } from '@lucide/svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import TextInput from '../forms/TextInput.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { PersonaTemplate } from '../../types/storyboard';

  interface Props {
    template?: PersonaTemplate | null;
    toggleClose?: () => void;
    handleSaved?: () => void;
    apiPrefix?: string;
    xfetch: ApiClient;
    notifications: NotificationService;
  }

  let {
    template = null,
    toggleClose = () => {},
    handleSaved = () => {},
    apiPrefix = '/api',
    xfetch,
    notifications,
  }: Props = $props();

  const MAX_PERSONAS = 32;

  let nextKey = 0;
  const newPersona = (persona?: PersonaTemplate['personas'][number]) => ({
    key: nextKey++,
    name: persona?.name ?? '',
    role: persona?.role ?? '',
    description: persona?.description ?? '',
  });

  let name = $state(untrack(() => template?.name ?? ''));
  let description = $state(untrack(() => template?.description ?? ''));
  let personas = $state(untrack(() => (template ? template.personas.map(p => newPersona(p)) : [newPersona()])));

  const title = $derived(template ? 'Edit persona template' : 'Create persona template');
  const names = $derived(personas.map(p => p.name.trim()));
  const hasDuplicates = $derived(new Set(names.filter(Boolean)).size !== names.filter(Boolean).length);
  const canSave = $derived(name.trim() !== '' && names.length > 0 && names.every(Boolean) && !hasDuplicates);

  function onSubmit(event: Event) {
    event.preventDefault();
    if (!canSave) return;

    xfetch(template ? `${apiPrefix}/persona-templates/${template.id}` : `${apiPrefix}/persona-templates`, {
      method: template ? 'PUT' : 'POST',
      body: {
        name,
        description,
        personas: personas.map(({ name, role, description }) => ({ name, role, description })),
      },
    })
      .then(() => {
        notifications.success('Persona template successfully saved');
        handleSaved();
      })
      .catch(() => {
        notifications.danger('Error saving persona template');
      });
  }

  const textareaClass =
    'w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-gray-100 resize-none';
  const labelClass = 'block text-gray-700 font-bold mb-2 dark:text-gray-400';
</script>

<Modal closeModal={toggleClose} ariaLabel={title}>
  <form onsubmit={onSubmit} name="personaTemplate" class="space-y-4 pt-6">
    <div>
      <label class={labelClass} for="personaTemplateName">Template Name</label>
      <TextInput
        bind:value={name}
        id="personaTemplateName"
        name="personaTemplateName"
        placeholder="Enter template name"
        maxlength={255}
        required
        autofocus
      />
    </div>
    <div>
      <label class={labelClass} for="personaTemplateDescription">Description</label>
      <textarea
        id="personaTemplateDescription"
        bind:value={description}
        placeholder="Enter description (optional)"
        maxlength="2000"
        rows="2"
        class={textareaClass}
      ></textarea>
    </div>

    <h3 class="font-bold text-gray-700 dark:text-gray-400">Personas</h3>
    {#each personas as persona, i (persona.key)}
      <fieldset class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2">
        <legend class="sr-only">Persona {i + 1}</legend>
        <div class="flex items-center gap-2">
          <TextInput
            bind:value={persona.name}
            aria-label="Persona {i + 1} name"
            placeholder="Name"
            maxlength={256}
            required
          />
          <TextInput
            bind:value={persona.role}
            aria-label="Persona {i + 1} role"
            placeholder="Role (optional)"
            maxlength={256}
          />
          <button
            type="button"
            class="p-1.5 text-gray-500 hover:text-red-600 dark:hover:text-red-400 disabled:opacity-30"
            aria-label="Remove persona {i + 1}"
            disabled={personas.length === 1}
            onclick={() => personas.splice(i, 1)}
          >
            <Trash2 class="w-5 h-5" />
          </button>
        </div>
        <textarea
          bind:value={persona.description}
          aria-label="Persona {i + 1} description"
          placeholder="Description (optional)"
          maxlength="2000"
          rows="2"
          class={textareaClass}
        ></textarea>
      </fieldset>
    {/each}

    {#if hasDuplicates}
      <p class="text-red-600 dark:text-red-400" role="alert">Persona names must be unique.</p>
    {/if}

    <div class="flex justify-between items-center">
      <HollowButton onClick={() => personas.push(newPersona())} disabled={personas.length >= MAX_PERSONAS}>
        <Plus class="w-4 h-4 me-1 inline" />
        Add Another Persona
      </HollowButton>
      <SolidButton type="submit" disabled={!canSave}>Save Template</SolidButton>
    </div>
  </form>
</Modal>

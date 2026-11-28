<script lang="ts">
  import { onMount } from 'svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import type { PersonaTemplate, StoryboardPersona } from '../../types/storyboard';

  interface Props {
    personas?: StoryboardPersona[];
    loadTemplates: () => Promise<PersonaTemplate[]>;
    onSelect: (templateId: string) => void;
  }

  let { personas = [], loadTemplates, onSelect }: Props = $props();

  let templates = $state<PersonaTemplate[]>([]);
  let loading = $state(true);
  let failed = $state(false);

  // persona names are unique per storyboard, ones already there aren't added again
  const existingNames = $derived(new Set(personas.map(p => p.name)));

  const newCount = (template: PersonaTemplate) => template.personas.filter(p => !existingNames.has(p.name)).length;

  const scopeLabel = (template: PersonaTemplate) => {
    if (template.teamId) return 'Team';
    if (template.organizationId) return 'Organization';
    return 'Global';
  };

  const buttonLabel = (count: number) => {
    if (count === 0) return 'Already added';
    return count === 1 ? 'Add 1 persona' : `Add ${count} personas`;
  };

  onMount(() => {
    loadTemplates()
      .then(result => {
        templates = result;
      })
      .catch(() => {
        failed = true;
      })
      .finally(() => {
        loading = false;
      });
  });
</script>

<div class="bg-gray-100 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 rounded-lg p-3 mb-4 space-y-3">
  {#if loading}
    <p class="text-gray-600 dark:text-gray-400">Loading templates…</p>
  {:else if failed}
    <p class="text-red-600 dark:text-red-400" role="alert">Failed to get persona templates</p>
  {:else}
    {#each templates as template (template.id)}
      {@const count = newCount(template)}
      <div
        class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-3"
        data-testid="persona-template-option"
      >
        <div class="flex justify-between items-start gap-3">
          <div>
            <h4 class="font-semibold text-gray-900 dark:text-gray-100">
              {template.name}
              <span class="ms-1 text-xs font-normal text-gray-500 dark:text-gray-400">{scopeLabel(template)}</span>
            </h4>
            {#if template.description}
              <p class="text-gray-600 dark:text-gray-400">{template.description}</p>
            {/if}
          </div>
          <SolidButton
            type="button"
            disabled={count === 0}
            additionalClasses="shrink-0"
            onClick={() => onSelect(template.id)}
          >
            {buttonLabel(count)}
          </SolidButton>
        </div>
        <ul class="flex flex-wrap gap-2 mt-2">
          {#each template.personas as persona (persona.name)}
            {@const exists = existingNames.has(persona.name)}
            <li
              class="px-2 py-0.5 rounded text-xs {exists
                ? 'line-through bg-gray-100 dark:bg-gray-700 text-gray-500'
                : 'bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-200'}"
              title={exists ? 'Already on this storyboard' : persona.description}
            >
              {persona.name}{persona.role ? ` · ${persona.role}` : ''}
            </li>
          {/each}
        </ul>
      </div>
    {:else}
      <p class="text-gray-600 dark:text-gray-400">No persona templates available</p>
    {/each}
  {/if}
</div>

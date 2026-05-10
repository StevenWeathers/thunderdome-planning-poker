<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectWithSubtext from '../forms/SelectWithSubtext.svelte';
  import { onMount } from 'svelte';
  import { AppConfig } from '../../config';

  import type { ApiClient } from '../../types/apiclient';
  import type { NotificationService } from '../../types/notifications';
  import type { ColorLegend, ColorLegendTemplate } from '../../types/storyboard';

  interface Props {
    handleLegendRevision?: any;
    toggleEditLegend?: any;
    colorLegend?: Array<ColorLegend>;
    isFacilitator?: boolean;
    teamId?: string;
    xfetch?: ApiClient;
    notifications?: NotificationService;
  }

  let {
    handleLegendRevision = () => {},
    toggleEditLegend = () => {},
    colorLegend = $bindable([]),
    isFacilitator = false,
    teamId = '',
    xfetch,
    notifications,
  }: Props = $props();

  const { SubscriptionsEnabled } = AppConfig;

  let canCopyFromTemplate = $state(!SubscriptionsEnabled);
  let teamSubscriptionLoaded = $state(!SubscriptionsEnabled);
  let showTemplateSelector = $state(false);
  let selectedTemplateId = $state('');
  let availableTemplates: Array<ColorLegendTemplate> = $state([]);
  let organizationId = $state('');
  let isLoadingTemplates = $state(false);
  let templatesLoaded = $state(false);

  type ColorLegendTemplateOption = ColorLegendTemplate & {
    displayName: string;
  };

  function cloneLegend(legend: Array<ColorLegend>): Array<ColorLegend> {
    return legend.map(color => ({ ...color }));
  }

  async function loadTemplateAccess() {
    if (!SubscriptionsEnabled) {
      canCopyFromTemplate = true;
      teamSubscriptionLoaded = true;
      return '';
    }

    if (!xfetch || !teamId) {
      canCopyFromTemplate = false;
      teamSubscriptionLoaded = true;
      return '';
    }

    try {
      const response = await xfetch(`/api/teams/${teamId}`);
      const result = await response.json();
      const nextOrganizationId = result.data.team?.organization_id ?? '';
      organizationId = nextOrganizationId;
      canCopyFromTemplate = Boolean(result.data.team?.subscribed);
      return nextOrganizationId;
    } catch {
      canCopyFromTemplate = false;
      notifications?.danger('Failed to get storyboard template access');
      return '';
    } finally {
      teamSubscriptionLoaded = true;
    }
  }

  async function loadTemplates() {
    if (!xfetch || !teamId || !canCopyFromTemplate || templatesLoaded) {
      return;
    }

    isLoadingTemplates = true;
    try {
      let currentOrganizationId = organizationId;
      if (SubscriptionsEnabled) {
        const response = await xfetch(`/api/teams/${teamId}`);
        const result = await response.json();
        currentOrganizationId = result.data.team?.organization_id ?? '';
        organizationId = currentOrganizationId;
      }

      const teamTemplateRequest = xfetch(`/api/teams/${teamId}/color-legend-templates`)
        .then(res => res.json())
        .then(result => result.data as Array<ColorLegendTemplate>);

      const requests: Array<Promise<Array<ColorLegendTemplate>>> = [teamTemplateRequest];

      if (currentOrganizationId) {
        requests.push(
          xfetch(`/api/organizations/${currentOrganizationId}/color-legend-templates`)
            .then(res => res.json())
            .then(result => result.data as Array<ColorLegendTemplate>),
        );
      }

      const results = await Promise.all(requests);
      availableTemplates = results.flatMap(result => result);
      templatesLoaded = true;
    } catch {
      notifications?.danger('Failed to load color legend templates');
    } finally {
      isLoadingTemplates = false;
    }
  }

  async function handleTemplateToggle() {
    showTemplateSelector = !showTemplateSelector;

    if (showTemplateSelector) {
      await loadTemplates();
    }
  }

  function applyTemplate(templateId: string) {
    const template = availableTemplates.find(item => item.id === templateId);
    if (!template) {
      return;
    }

    const legendMap = new Map(template.colorLegend.map(item => [item.color, item.legend]));
    colorLegend = colorLegend.map(item => ({
      ...item,
      legend: legendMap.get(item.color) ?? '',
    }));
  }

  function handleTemplateSelect(event: CustomEvent<ColorLegendTemplateOption>) {
    selectedTemplateId = event.detail.id;
    applyTemplate(selectedTemplateId);
  }

  const getTemplateLabel = (template: ColorLegendTemplate) => {
    const scope = template.teamId ? 'Team' : 'Organization';
    return `${scope}: ${template.name}`;
  };

  const templateOptions = $derived(
    availableTemplates.map(template => ({
      ...template,
      displayName: getTemplateLabel(template),
    })),
  );

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    handleLegendRevision(colorLegend);
    toggleEditLegend();
  }

  let focusInput: any = $state();
  onMount(async () => {
    colorLegend = cloneLegend(colorLegend);
    focusInput?.focus();
    await loadTemplateAccess();
  });
</script>

<Modal
  closeModal={toggleEditLegend}
  ariaLabel={$LL.modalStoryboardColorLegend()}
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
  <form onsubmit={handleSubmit} name="colorLegend" class="space-y-4 pt-6">
    <h2 class="text-xl font-bold dark:text-gray-300">Story Color Legend</h2>
    {#if teamSubscriptionLoaded && canCopyFromTemplate}
      <div class="space-y-3 rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900/60">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h3 class="font-semibold text-gray-900 dark:text-gray-100">Copy From Template</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              Use a saved team or organization legend as a starting point.
            </p>
          </div>
          <SolidButton type="button" onClick={handleTemplateToggle} disabled={isLoadingTemplates}>
            Copy From Template
          </SolidButton>
        </div>

        {#if showTemplateSelector}
          {#if availableTemplates.length > 0}
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="templateSelector">
                Choose Template
              </label>
              <SelectWithSubtext
                items={templateOptions}
                label="Select a color legend template..."
                selectedItemId={selectedTemplateId}
                itemType="color_legend_template"
                nameField="displayName"
                on:change={handleTemplateSelect}
              />
            </div>
          {:else if isLoadingTemplates}
            <p class="text-sm text-gray-600 dark:text-gray-400">Loading color legend templates...</p>
          {:else}
            <p class="text-sm text-gray-600 dark:text-gray-400">
              No color legend templates are available for this storyboard yet.
            </p>
          {/if}
        {/if}
      </div>
    {/if}
    <div class="space-y-2">
      {#each colorLegend as color, i}
        <div class="group">
          <label class="flex-1 min-w-0">
            <span class="sr-only">Color legend for {color.color}</span>
            {#if i === 0}
              <TextInput
                bind:this={focusInput}
                placeholder={$LL.legendRetroPlaceholder()}
                name="legend-{color.color}"
                disabled={!isFacilitator}
                value={colorLegend[i].legend}
              >
                {#snippet startElement()}
                  <div class="w-6 h-6 rounded bg-gray-400 colorcard-{color.color}"></div>
                {/snippet}
              </TextInput>
            {:else}
              <TextInput
                placeholder={$LL.legendRetroPlaceholder()}
                name="legend-{color.color}"
                disabled={!isFacilitator}
                value={colorLegend[i].legend}
              >
                {#snippet startElement()}
                  <div class="w-6 h-6 rounded bg-gray-400 colorcard-{color.color}"></div>
                {/snippet}
              </TextInput>
            {/if}
          </label>
        </div>
      {/each}
    </div>
    <div class="flex justify-end">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

<style lang="postcss">
  .colorcard-gray {
    @apply bg-gray-400;
  }

  .colorcard-red {
    @apply bg-red-400;
  }

  .colorcard-orange {
    @apply bg-orange-400;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
  }

  .colorcard-green {
    @apply bg-green-400;
  }

  .colorcard-teal {
    @apply bg-teal-400;
  }

  .colorcard-blue {
    @apply bg-blue-400;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
  }

  .colorcard-purple {
    @apply bg-purple-400;
  }

  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>

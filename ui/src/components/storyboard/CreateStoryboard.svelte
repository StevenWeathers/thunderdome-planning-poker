<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import SelectWithSubtext from '../forms/SelectWithSubtext.svelte';
  import { Crown, Lock } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { ColorLegend, ColorLegendTemplate } from '../../types/storyboard';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
    apiPrefix?: string;
    scope?: 'user' | 'team' | 'project';
    colorLegend?: Array<ColorLegend>;
    teamId?: string;
    organizationId?: string;
    subscribed?: boolean;
  }

  let {
    xfetch,
    notifications,
    router,
    apiPrefix = '/api',
    scope = 'user',
    colorLegend = [],
    teamId = '',
    organizationId = '',
    subscribed = false,
  }: Props = $props();

  const { SubscriptionsEnabled } = AppConfig;

  type TeamOption = {
    id: string;
    name: string;
  };

  type ColorLegendTemplateOption = ColorLegendTemplate & {
    displayName: string;
  };

  let storyboardName = $state('');
  let joinCode = $state('');
  let facilitatorCode = $state('');
  let selectedTeam = $state('');
  let teams: Array<TeamOption> = $state([]);
  let availableTemplates: Array<ColorLegendTemplate> = $state([]);
  let selectedTemplateId = $state('');
  let selectedColorLegend: Array<ColorLegend> = $state([]);
  let templateScopeTeamId = $state('');
  let templateScopeOrganizationId = $state('');
  let canCopyFromTemplate = $state(!SubscriptionsEnabled);
  let templateAccessLoaded = $state(false);
  let templatesLoaded = $state(false);
  let isLoadingTemplates = $state(false);

  /** @type {TextInput} */
  let storyboardNameTextInput: { focus: () => void } | undefined = $state();

  function cloneLegend(legend: Array<ColorLegend>): Array<ColorLegend> {
    return legend.map(item => ({ ...item }));
  }

  function resetTemplateSelection() {
    availableTemplates = [];
    selectedTemplateId = '';
    selectedColorLegend = cloneLegend(colorLegend);
    templatesLoaded = false;
  }

  function getTemplateLabel(template: ColorLegendTemplate) {
    return `${template.teamId ? 'Team' : 'Organization'}: ${template.name}`;
  }

  function applyTemplate(templateId: string) {
    const template = availableTemplates.find(item => item.id === templateId);
    if (!template) {
      return;
    }

    selectedColorLegend = cloneLegend(template.colorLegend);
  }

  function handleTemplateSelect(event: CustomEvent<ColorLegendTemplateOption>) {
    selectedTemplateId = event.detail.id;
    applyTemplate(selectedTemplateId);
  }

  async function loadSelectedTeamTemplateAccess(teamSelectionId: string) {
    if (!teamSelectionId) {
      templateScopeTeamId = '';
      templateScopeOrganizationId = '';
      canCopyFromTemplate = false;
      templateAccessLoaded = true;
      return;
    }

    try {
      const response = await xfetch(`/api/teams/${teamSelectionId}`);
      const result = await response.json();

      templateScopeTeamId = teamSelectionId;
      templateScopeOrganizationId = result.data.team?.organization_id ?? '';
      canCopyFromTemplate = !SubscriptionsEnabled || Boolean(result.data.team?.subscribed);
    } catch {
      canCopyFromTemplate = false;
      templateScopeTeamId = '';
      templateScopeOrganizationId = '';
      notifications.danger('Failed to get color legend template access');
    } finally {
      templateAccessLoaded = true;
    }
  }

  async function loadTemplates() {
    if ((!templateScopeTeamId && !templateScopeOrganizationId) || !canCopyFromTemplate || templatesLoaded) {
      return;
    }

    isLoadingTemplates = true;
    try {
      const requests: Array<Promise<Array<ColorLegendTemplate>>> = [];

      if (templateScopeTeamId) {
        requests.push(
          xfetch(`/api/teams/${templateScopeTeamId}/color-legend-templates`)
            .then(res => res.json())
            .then(result => result.data as Array<ColorLegendTemplate>),
        );
      }

      if (templateScopeOrganizationId) {
        requests.push(
          xfetch(`/api/organizations/${templateScopeOrganizationId}/color-legend-templates`)
            .then(res => res.json())
            .then(result => result.data as Array<ColorLegendTemplate>),
        );
      }

      const results = await Promise.all(requests);
      availableTemplates = results.flatMap(result => result);
      templatesLoaded = true;
    } catch {
      notifications.danger('Failed to load color legend templates');
    } finally {
      isLoadingTemplates = false;
    }
  }

  const templateOptions = $derived(
    availableTemplates.map(template => ({
      ...template,
      displayName: getTemplateLabel(template),
    })),
  );

  $effect(() => {
    if (scope !== 'user') {
      return;
    }

    resetTemplateSelection();
    templateAccessLoaded = false;
    void loadSelectedTeamTemplateAccess(selectedTeam);
  });

  $effect(() => {
    if (scope === 'user') {
      return;
    }

    templateScopeTeamId = teamId;
    templateScopeOrganizationId = organizationId;
    canCopyFromTemplate = !SubscriptionsEnabled || subscribed;
    templateAccessLoaded = true;
  });

  $effect(() => {
    if (!templateAccessLoaded || !canCopyFromTemplate) {
      return;
    }

    if (!templateScopeTeamId && !templateScopeOrganizationId) {
      return;
    }

    void loadTemplates();
  });

  function createStoryboard(e: SubmitEvent) {
    e.preventDefault();
    let endpoint = scope === 'project' ? `${apiPrefix}/storyboards` : `${apiPrefix}/users/${$user.id}/storyboards`;
    const body = {
      storyboardName,
      joinCode,
      facilitatorCode,
      ...(selectedColorLegend.length > 0 ? { colorLegend: selectedColorLegend } : {}),
    };

    if (scope !== 'project' && selectedTeam !== '') {
      endpoint = `/api/teams/${selectedTeam}/users/${$user.id}/storyboards`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function ({ data }) {
        router.route(`${appRoutes.storyboard}/${data.id}`);
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            notifications.danger(`Error encountered creating storyboard : ${result.error}`);
          });
        } else {
          notifications.danger(`Error encountered creating storyboard`);
        }
      });
  }

  function getTeams() {
    xfetch(`/api/users/${$user.id}/teams?limit=100`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getTeamsError());
      });
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.register);
    }

    selectedColorLegend = cloneLegend(colorLegend);

    getTeams();

    // Focus the storyboard name input field
    storyboardNameTextInput?.focus();
  });
</script>

<form onsubmit={createStoryboard} name="createStoryboard">
  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="storyboardName">
      {$LL.storyboardName()}
    </label>
    <div class="control">
      <TextInput
        name="storyboardName"
        bind:value={storyboardName}
        bind:this={storyboardNameTextInput}
        placeholder={$LL.storyboardNamePlaceholder()}
        id="storyboardName"
        required
      />
    </div>
  </div>

  {#if scope === 'user' && apiPrefix === '/api' && $user.rank !== 'GUEST'}
    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2" for="selectedTeam">
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}
        {/if}
      </label>
      <SelectInput bind:value={selectedTeam} id="selectedTeam" name="selectedTeam">
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value={team.id}>
            {team.name}
          </option>
        {/each}
      </SelectInput>
    </div>
  {/if}

  {#if templateAccessLoaded && canCopyFromTemplate && (templateScopeTeamId || templateScopeOrganizationId)}
    <div
      class="mb-4 space-y-3 rounded-lg border border-gray-200 bg-gray-50 p-4 dark:border-gray-700 dark:bg-gray-900/60"
    >
      <p class="text-sm text-gray-600 dark:text-gray-400">
        Start this storyboard with a saved team or organization color legend.
      </p>

      {#if availableTemplates.length > 0}
        <div class="space-y-2">
          <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="templateSelector">
            Color Legend
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
          No color legend templates are available for this scope yet.
        </p>
      {/if}
    </div>
  {/if}

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="joinCode">
      {$LL.passCode()}
    </label>
    <div class="control">
      <TextInput
        name="joinCode"
        bind:value={joinCode}
        placeholder={$LL.optionalPasscodePlaceholder()}
        id="joinCode"
        icon={Lock}
      />
    </div>
  </div>

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="facilitatorCode">
      {$LL.facilitatorCodeOptional()}
    </label>
    <div class="control">
      <TextInput
        name="facilitatorCode"
        bind:value={facilitatorCode}
        placeholder={$LL.facilitatorCodePlaceholder()}
        id="facilitatorCode"
        icon={Crown}
      />
    </div>
  </div>

  <div class="text-right">
    <SolidButton type="submit">{$LL.createStoryboard()}</SolidButton>
  </div>
</form>

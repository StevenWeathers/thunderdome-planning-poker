<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import SelectWithSubtext from '../forms/SelectWithSubtext.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { Crown, Lock } from 'lucide-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
    apiPrefix?: string;
    scope?: 'user' | 'project';
  }

  let {
    xfetch,
    notifications,
    router,
    apiPrefix = '/api',
    scope = 'user',
  }: Props = $props();

  const maxPhaseTimeLimitMin = 59;

  let teams = $state([]);
  let retroTemplates = $state([]);
  let publicTemplates = $state([]);
  let teamRetroTemplates = $state([]);
  let organizationRetroTemplates = $state([]);
  let defaultRetroSettings = {
    retroName: '',
    maxVotes: 3,
    brainstormVisibility: 'visible',
    phaseTimeLimit: 0,
    facilitatorCode: '',
    joinCode: '',
    selectedTeam: '',
    templateId: '',
    phaseAutoAdvance: true,
    allowCumulativeVoting: false,
    hideVotesDuringVoting: false
  };
  let retroSettings = $state({ ...defaultRetroSettings });
  let orgRetroSettings = $state({});
  let departmentRetroSettings = $state({});
  let teamRetroSettings = $state({});

  /** @type {TextInput} */
  let retroNameTextInput = $state();

  const brainstormVisibilityOptions = [
    {
      label: $LL.brainstormVisibilityLabelVisible(),
      value: 'visible',
    },
    {
      label: $LL.brainstormVisibilityLabelConcealed(),
      value: 'concealed',
    },
    {
      label: $LL.brainstormVisibilityLabelHidden(),
      value: 'hidden',
    },
  ];

  function createRetro(e) {
    e.preventDefault();
    let endpoint = scope === 'project' ? `${apiPrefix}/retros` : `${apiPrefix}/users/${$user.id}/retros`;

    if (retroSettings.templateId === '') {
      notifications.danger('Must select a retrospective template.');
      return;
    }

    if (retroSettings.phaseTimeLimit > maxPhaseTimeLimitMin || retroSettings.phaseTimeLimit < 0) {
      notifications.danger('Phase Time Limit minutes must be between 0-59');
      return;
    }

    const body = {
      retroName: retroSettings.retroName,
      joinCode: retroSettings.joinCode,
      facilitatorCode: retroSettings.facilitatorCode,
      maxVotes: retroSettings.maxVotes,
      brainstormVisibility: retroSettings.brainstormVisibility,
      phaseTimeLimitMin: retroSettings.phaseTimeLimit,
      phaseAutoAdvance: retroSettings.phaseAutoAdvance,
      allowCumulativeVoting: retroSettings.allowCumulativeVoting,
      templateId: retroSettings.templateId,
      hideVotesDuringVoting: retroSettings.hideVotesDuringVoting
    };

    if (scope !== 'project' && retroSettings.selectedTeam !== '') {
      endpoint = `/api/teams/${retroSettings.selectedTeam}/users/${$user.id}/retros`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function ({ data }) {
        router.route(`${appRoutes.retro}/${data.id}`);
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
            notifications.danger(
              `${$LL.createRetroErrorMessage()} : ${result.error}`,
            );
          });
        } else {
          notifications.danger($LL.createRetroErrorMessage());
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

  function getTemplatesPublic() {
    xfetch(`/api/retro-templates/public`)
      .then(res => res.json())
      .then(function (result) {
        publicTemplates = result.data;
        combineRetroTemplates();
      })
      .catch(function () {
        notifications.danger('error getting public templates');
      });
  }

  function getCustomRetroSettings() {
    teamRetroSettings = {};
    departmentRetroSettings = {};
    orgRetroSettings = {};
    combineRetroSettings();

    // don't get custom retro settings if a team isn't selected
    if (retroSettings.selectedTeam === '') {
      return;
    }
    const team = teams.find(t => t.id === retroSettings.selectedTeam);
    // if subscriptions are enabled and the team (or its parent org) isn't subscribed
    // don't attempt to get retro settings
    if (
      AppConfig.SubscriptionsEnabled &&
      !validateUserIsAdmin($user) &&
      !team.subscribed
    ) {
      return;
    }

    xfetch(`/api/teams/${team.id}/retro-settings`)
      .then(res => res.json())
      .then(function (result) {
        if (!result.data) {
          return;
        }
        teamRetroSettings = result.data;
        combineRetroSettings();
      })
      .catch(function () {
        notifications.danger('Failed to get team retro settings');
      });

    if (team.organization_id !== '') {
      xfetch(`/api/organizations/${team.organization_id}/retro-settings`)
        .then(res => res.json())
        .then(function (result) {
          if (!result.data) {
            return;
          }
          orgRetroSettings = result.data;
          combineRetroSettings();
        })
        .catch(function () {
          notifications.danger('Failed to get organization retro settings');
        });
    }

    if (team.department_id !== '') {
      xfetch(
        `/api/organizations/${team.organization_id}/departments/${team.department_id}/retro-settings`,
      )
        .then(res => res.json())
        .then(function (result) {
          if (!result.data) {
            return;
          }
          departmentRetroSettings = result.data;
          combineRetroSettings();
        })
        .catch(function () {
          notifications.danger('Failed to get department retro settings');
        });
    }
  }

  const combineRetroSettings = () => {
    // settings priority order (Team -> Department -> Organization -> Default)
    retroSettings = {
      ...retroSettings,
      ...orgRetroSettings,
      ...departmentRetroSettings,
      ...teamRetroSettings,
    };
  };

  function getPrivateRetroTemplates() {
    teamRetroTemplates = [];
    organizationRetroTemplates = [];
    combineRetroTemplates();

    // don't get private templates if a team isn't selected
    if (retroSettings.selectedTeam === '') {
      return;
    }
    const team = teams.find(t => t.id === retroSettings.selectedTeam);
    // if subscriptions are enabled and the team (or its parent org) isn't subscribed
    // don't attempt to get private templates
    if (
      AppConfig.SubscriptionsEnabled &&
      !validateUserIsAdmin($user) &&
      !team.subscribed
    ) {
      return;
    }
    const orgPrefix =
      team.organization_id !== ''
        ? `/api/organizations/${team.organization_id}`
        : `/api`;
    const teamPrefix =
      team.department_id !== ''
        ? `${orgPrefix}/departments/${team.department_id}`
        : `/api`;

    xfetch(`${teamPrefix}/teams/${team.id}/retro-templates`)
      .then(res => res.json())
      .then(function (result) {
        teamRetroTemplates = result.data;
        combineRetroTemplates();
      })
      .catch(function () {
        notifications.danger('Failed to get team retro templates');
      });

    if (team.organization_id !== '') {
      xfetch(`${orgPrefix}/retro-templates`)
        .then(res => res.json())
        .then(function (result) {
          organizationRetroTemplates = result.data;
          combineRetroTemplates();
        })
        .catch(function () {
          notifications.danger('Failed to get organization retro templates');
        });
    }
  }

  const combineRetroTemplates = () => {
    let defaultFound = false;
    // templates priority order (Team -> Organization -> Public)
    retroTemplates = [
      ...teamRetroTemplates,
      ...organizationRetroTemplates,
      ...publicTemplates,
    ];

    retroTemplates.map(temp => {
      // Find default template with priority order (Team -> Organization -> Public)
      if (!defaultFound && temp.defaultTemplate) {
        retroSettings.templateId = temp.id;
        defaultFound = true;
      }
    });
  };

  const updateSelectedTemplate = event => {
    retroSettings.templateId = event.detail.id;
  };

  function teamSelected() {
    // apply settings before templates
    getCustomRetroSettings();
    getPrivateRetroTemplates();
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.register);
    }
    getTeams();
    getTemplatesPublic();

    // Focus the retro name input field
    retroNameTextInput.focus();
  });
</script>

<form onsubmit={createRetro} name="createRetro">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="retroName"
    >
      {$LL.retroName()}
    </label>
    <div class="control">
      <TextInput
        name="retroName"
        bind:value={retroSettings.retroName}
        bind:this={retroNameTextInput}
        placeholder={$LL.retroNamePlaceholder()}
        id="retroName"
        required
      />
    </div>
  </div>

  {#if apiPrefix === '/api' && $user.rank !== 'GUEST'}
    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
        for="selectedTeam"
      >
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}
        {/if}
      </label>
      <SelectInput
        bind:value={retroSettings.selectedTeam}
        on:change={teamSelected}
        id="selectedTeam"
        name="selectedTeam"
      >
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value={team.id}>
            {team.name}
          </option>
        {/each}
      </SelectInput>
    </div>
  {/if}

  <div class="mb-4">
    <div
      class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
    >
      Retro Template
    </div>
    <SelectWithSubtext
      on:change={updateSelectedTemplate}
      items={retroTemplates}
      label="Select a retro template..."
      selectedItemId={retroSettings.templateId}
      itemType="retro_template"
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="joinCode"
    >
      {$LL.joinCodeLabelOptional()}
    </label>
    <div class="control">
      <TextInput
        name="joinCode"
        bind:value={retroSettings.joinCode}
        placeholder={$LL.joinCodePlaceholder()}
        id="joinCode"
        icon={Lock}
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="facilitatorCode"
    >
      {$LL.facilitatorCodeOptional()}
    </label>
    <div class="control">
      <TextInput
        name="facilitatorCode"
        bind:value={retroSettings.facilitatorCode}
        placeholder={$LL.facilitatorCodePlaceholder()}
        id="facilitatorCode"
        icon={Crown}
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="maxVotes"
    >
      {$LL.retroMaxVotesPerUserLabel()}
    </label>
    <div class="control">
      <TextInput
        name="maxVotes"
        bind:value={retroSettings.maxVotes}
        id="maxVotes"
        type="number"
        min="1"
        max="10"
        required
      />
    </div>
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked={retroSettings.allowCumulativeVoting}
      id="allowCumulativeVoting"
      name="allowCumulativeVoting"
      label={$LL.allowCumulativeVotingLabel()}
    />
  </div>

  <div class="mb-4">
    <label
      class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="brainstormVisibility"
    >
      {$LL.brainstormPhaseFeedbackVisibility()}
    </label>
    <SelectInput
      bind:value={retroSettings.brainstormVisibility}
      id="brainstormVisibility"
      name="brainstormVisibility"
    >
      {#each brainstormVisibilityOptions as item}
        <option value={item.value}>
          {item.label}
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="phaseTimeLimitMin"
    >
      {$LL.retroPhaseTimeLimitMinLabel()}
    </label>
    <div class="control">
      <TextInput
        name="phaseTimeLimitMin"
        bind:value={retroSettings.phaseTimeLimit}
        id="phaseTimeLimitMin"
        type="number"
        min="0"
        max={maxPhaseTimeLimitMin}
        required
      />
    </div>
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked={retroSettings.phaseAutoAdvance}
      id="phaseAutoAdvance"
      name="phaseAutoAdvance"
      label={$LL.phaseAutoAdvanceLabel()}
    />
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked={retroSettings.hideVotesDuringVoting}
      id="hideVotesDuringVoting"
      name="hideVotesDuringVoting"
      label={`Hide Votes During Voting Phase`}
    />
  </div>

  <div class="text-right">
    <SolidButton type="submit">{$LL.createRetro()}</SolidButton>
  </div>
</form>

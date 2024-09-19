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

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let router;
  export let apiPrefix = '/api';

  const maxPhaseTimeLimitMin = 59;

  let retroName = '';
  let joinCode = '';
  let facilitatorCode = '';
  let maxVotes = '3';
  let brainstormVisibility = 'visible';
  let teams = [];
  let selectedTeam = '';
  let phaseTimeLimitMin = 0;
  let templateId = '';
  let phaseAutoAdvance = true;
  let allowCumulativeVoting = false;
  let retroTemplates = [];
  let publicTemplates = [];
  let teamRetroTemplates = [];
  let organizationRetroTemplates = [];

  /** @type {TextInput} */
  let retroNameTextInput;

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
    let endpoint = `${apiPrefix}/users/${$user.id}/retros`;

    if (templateId === '') {
      notifications.danger('Must select a retrospective template.');
      return;
    }

    if (phaseTimeLimitMin > maxPhaseTimeLimitMin || phaseTimeLimitMin < 0) {
      notifications.danger('Phase Time Limit minutes must be between 0-59');
      return;
    }

    const body = {
      retroName,
      joinCode,
      facilitatorCode,
      maxVotes: parseInt(maxVotes, 10),
      brainstormVisibility,
      phaseTimeLimitMin: parseInt(`${phaseTimeLimitMin}`, 10),
      phaseAutoAdvance,
      allowCumulativeVoting,
      templateId,
    };

    if (selectedTeam !== '') {
      endpoint = `/api/teams/${selectedTeam}/users/${$user.id}/retros`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function ({ data }) {
        eventTag('create_retro', 'engagement', 'success', () => {
          router.route(`${appRoutes.retro}/${data.id}`);
        });
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
        eventTag('create_retro', 'engagement', 'failure');
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

  function getPrivateRetroTemplates() {
    teamRetroTemplates = [];
    organizationRetroTemplates = [];
    combineRetroTemplates();

    // don't get private templates if a team isn't selected
    if (selectedTeam === '') {
      return;
    }
    const team = teams.find(t => t.id === selectedTeam);
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
        templateId = temp.id;
        defaultFound = true;
      }
    });
  };

  const updateSelectedTemplate = event => {
    templateId = event.detail.id;
  };

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

<form on:submit="{createRetro}" name="createRetro">
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
        bind:value="{retroName}"
        bind:this="{retroNameTextInput}"
        placeholder="{$LL.retroNamePlaceholder()}"
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
        bind:value="{selectedTeam}"
        on:change="{getPrivateRetroTemplates}"
        id="selectedTeam"
        name="selectedTeam"
      >
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value="{team.id}">
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
      on:change="{updateSelectedTemplate}"
      items="{retroTemplates}"
      label="Select a retro template..."
      selectedItemId="{templateId}"
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
        bind:value="{joinCode}"
        placeholder="{$LL.joinCodePlaceholder()}"
        id="joinCode"
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
        bind:value="{facilitatorCode}"
        placeholder="{$LL.facilitatorCodePlaceholder()}"
        id="facilitatorCode"
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
        bind:value="{maxVotes}"
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
      bind:checked="{allowCumulativeVoting}"
      id="allowCumulativeVoting"
      name="allowCumulativeVoting"
      label="{$LL.allowCumulativeVotingLabel()}"
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
      bind:value="{brainstormVisibility}"
      id="brainstormVisibility"
      name="brainstormVisibility"
    >
      {#each brainstormVisibilityOptions as item}
        <option value="{item.value}">
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
        bind:value="{phaseTimeLimitMin}"
        id="phaseTimeLimitMin"
        type="number"
        min="0"
        max="{maxPhaseTimeLimitMin}"
        required
      />
    </div>
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked="{phaseAutoAdvance}"
      id="phaseAutoAdvance"
      name="phaseAutoAdvance"
      label="{$LL.phaseAutoAdvanceLabel()}"
    />
  </div>

  <div class="text-right">
    <SolidButton type="submit">{$LL.createRetro()}</SolidButton>
  </div>
</form>

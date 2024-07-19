<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let router;
  export let apiPrefix = '/api';

  let retroName = '';
  let joinCode = '';
  let facilitatorCode = '';
  let maxVotes = '3';
  let brainstormVisibility = 'visible';
  let teams = [];
  let selectedTeam = '';

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
    const body = {
      retroName,
      format: 'worked_improve_question',
      joinCode,
      facilitatorCode,
      maxVotes: parseInt(maxVotes, 10),
      brainstormVisibility,
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

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.register);
    }
    getTeams();

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

  <div class="text-right">
    <SolidButton type="submit">{$LL.createRetro()}</SolidButton>
  </div>
</form>

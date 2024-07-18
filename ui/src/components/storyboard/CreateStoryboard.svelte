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

  let storyboardName = '';
  let joinCode = '';
  let facilitatorCode = '';
  let selectedTeam = '';
  let teams = [];

  /** @type {TextInput} */
  let storyboardNameTextInput;

  function createStoryboard(e) {
    e.preventDefault();
    let endpoint = `${apiPrefix}/users/${$user.id}/storyboards`;
    const body = {
      storyboardName,
      joinCode,
      facilitatorCode,
    };

    if (selectedTeam !== '') {
      endpoint = `/api/teams/${selectedTeam}/users/${$user.id}/storyboards`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function ({ data }) {
        eventTag('create_storyboard', 'engagement', 'success', () => {
          router.route(`${appRoutes.storyboard}/${data.id}`);
        });
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
            notifications.danger(
              `Error encountered creating storyboard : ${result.error}`,
            );
          });
        } else {
          notifications.danger(`Error encountered creating storyboard`);
        }
        eventTag('create_storyboard', 'engagement', 'failure');
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

    // Focus the storyboard name input field
    storyboardNameTextInput.focus();
  });
</script>

<form on:submit="{createStoryboard}" name="createStoryboard">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="storyboardName"
    >
      {$LL.storyboardName()}
    </label>
    <div class="control">
      <TextInput
        name="storyboardName"
        bind:value="{storyboardName}"
        bind:this="{storyboardNameTextInput}"
        placeholder="{$LL.storyboardNamePlaceholder()}"
        id="storyboardName"
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
      {$LL.passCode()}
    </label>
    <div class="control">
      <TextInput
        name="joinCode"
        bind:value="{joinCode}"
        placeholder="{$LL.optionalPasscodePlaceholder()}"
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

  <div class="text-right">
    <SolidButton type="submit">{$LL.createStoryboard()}</SolidButton>
  </div>
</form>

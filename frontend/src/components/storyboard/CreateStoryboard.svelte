<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../SolidButton.svelte';
  import { warrior as user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import ChevronDown from '../icons/ChevronDown.svelte';

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
  });
</script>

<form on:submit="{createStoryboard}" name="createStoryboard">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="storyboardName"
    >
      Storyboard Name
    </label>
    <div class="control">
      <input
        name="storyboardName"
        bind:value="{storyboardName}"
        placeholder="Enter a storyboard name"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                    rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                    focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="storyboardName"
        required
      />
    </div>
  </div>

  {#if apiPrefix === '/api'}
    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
        for="selectedTeam"
      >
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}
        {/if}
      </label>
      <div class="relative">
        <select
          bind:value="{selectedTeam}"
          class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pe-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
          id="selectedTeam"
          name="selectedTeam"
        >
          <option value="" disabled>{$LL.selectTeam()}</option>
          {#each teams as team}
            <option value="{team.id}">
              {team.name}
            </option>
          {/each}
        </select>
        <div
          class="pointer-events-none absolute inset-y-0 end-0 flex
                items-center px-2 text-gray-700 dark:text-gray-400"
        >
          <ChevronDown />
        </div>
      </div>
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
      <input
        name="joinCode"
        bind:value="{joinCode}"
        placeholder="{$LL.optionalPasscodePlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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
      <input
        name="facilitatorCode"
        bind:value="{facilitatorCode}"
        placeholder="{$LL.facilitatorCodePlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2
            appearance-none
            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500
            dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="facilitatorCode"
      />
    </div>
  </div>

  <div class="text-right">
    <SolidButton type="submit">Create Storyboard</SolidButton>
  </div>
</form>

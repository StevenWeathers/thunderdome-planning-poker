<script lang="ts">
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import { onMount } from 'svelte';

  export let eventTag;
  export let notifications;
  export let xfetch;
  export let handleImport;
  export let gameId = '';

  let selectedGameIdx = '';
  let games = [];
  let game = {
    plans: [],
  };

  function getGames() {
    xfetch(`/api/users/${$user.id}/battles?limit=${9999}&offset=${0}`)
      .then(res => res.json())
      .then(function (result) {
        games = result.data;
      })
      .catch(function () {
        notifications.danger($LL.myBattlesError());
        eventTag('fetch_games', 'engagement', 'failure');
      });
  }

  function getGameStories() {
    if (selectedGameIdx === null) {
      notifications.danger('Game not selected');
      return;
    }
    const gameId = games[selectedGameIdx].id;
    xfetch(`/api/battles/${gameId}`)
      .then(res => res.json())
      .then(function (result) {
        game = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getBattleError());
        eventTag('fetch_game', 'engagement', 'failure');
      });
  }

  const importStory = idx => () => {
    handleImport(game.plans[idx]);
  };

  onMount(() => {
    getGames();
  });
</script>

<div class="mb-4">
  <SelectInput
    id="jirainstance"
    bind:value="{selectedGameIdx}"
    on:change="{getGameStories}"
  >
    >
    <option value="" disabled>Select game to import from</option>
    {#each games as game, idx}
      {#if game.id !== gameId}
        <option value="{idx}">{game.name}</option>
      {/if}
    {/each}
  </SelectInput>
</div>

{#if selectedGameIdx !== ''}
  <div class="mb-4">
    {#each game.plans as story, idx}
      <div
        class="p-2 w-full flex flex-wrap justify-between bg-gray-200 dark:bg-gray-900 dark:text-white rounded shadow mb-2 items-center"
      >
        <div>
          <span
            class="inline-block text-sm text-gray-500 dark:text-gray-300
                  border-gray-300 border px-1 rounded mr-1.5"
            data-testid="plan-type"
          >
            {story.type}</span
          >[{story.referenceId}] {story.name}
        </div>
        <div>
          <SolidButton onClick="{importStory(idx)}">Import</SolidButton>
        </div>
      </div>
    {/each}
  </div>
{/if}

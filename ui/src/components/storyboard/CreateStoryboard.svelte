<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { Crown, Lock } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
    apiPrefix?: string;
    scope?: 'user' | 'project';
  }

  let { xfetch, notifications, router, apiPrefix = '/api', scope = 'user' }: Props = $props();

  let storyboardName = $state('');
  let joinCode = $state('');
  let facilitatorCode = $state('');
  let selectedTeam = $state('');
  let teams = $state([]);

  /** @type {TextInput} */
  let storyboardNameTextInput = $state();

  function createStoryboard(e: SubmitEvent) {
    e.preventDefault();
    let endpoint = scope === 'project' ? `${apiPrefix}/storyboards` : `${apiPrefix}/users/${$user.id}/storyboards`;
    const body = {
      storyboardName,
      joinCode,
      facilitatorCode,
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
    getTeams();

    // Focus the storyboard name input field
    storyboardNameTextInput.focus();
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

  {#if apiPrefix === '/api' && $user.rank !== 'GUEST'}
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

<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import DownCarrotIcon from '../icons/ChevronDown.svelte';
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';

  const allowedPointValues = AppConfig.AllowedPointValues;
  const allowedPointAverages = ['ceil', 'round', 'floor'];

  export let toggleEditBattle = () => {};
  export let handleBattleEdit = (battle: any) => {};
  export let points = [];
  export let battleName = '';
  export let votingLocked = false;
  export let autoFinishVoting = true;
  export let pointAverageRounding = 'ceil';
  export let joinCode = '';
  export let leaderCode = '';
  export let hideVoterIdentity = false;
  export let teamId = '';
  export let notifications: any;
  export let xfetch: any;

  let checkedPointColor =
    'border-green-500 bg-green-100 text-green-600 dark:bg-gray-900 dark:text-lime-500 dark:border-lime-500';
  let uncheckedPointColor =
    'border-gray-300 bg-white dark:bg-gray-900 dark:border-gray-600 dark:text-gray-300';

  let teams = [];

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

  function saveBattle(e) {
    e.preventDefault();

    const pointValuesAllowed = allowedPointValues.filter(pv => {
      return points.includes(pv);
    });

    const battle = {
      battleName,
      pointValuesAllowed,
      autoFinishVoting,
      pointAverageRounding,
      hideVoterIdentity,
      joinCode,
      leaderCode,
      teamId,
    };

    handleBattleEdit(battle);
  }

  onMount(() => {
    getTeams();
  });
</script>

<Modal
  closeModal="{toggleEditBattle}"
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
  <form on:submit="{saveBattle}" name="createBattle">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="battleName"
      >
        {$LL.battleName({ friendly: AppConfig.FriendlyUIVerbs })}
      </label>
      <div class="control">
        <TextInput
          name="battleName"
          bind:value="{battleName}"
          placeholder="{$LL.battleNamePlaceholder({
            friendly: AppConfig.FriendlyUIVerbs,
          })}"
          id="battleName"
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <h3 class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
        {$LL.pointValuesAllowed()}
      </h3>
      <div class="control relative -me-2 md:-me-1">
        {#if !votingLocked}
          <div class="font-bold text-red-500">
            {$LL.battleEditPointsDisabled()}
          </div>
        {/if}
        {#each allowedPointValues as point}
          <label
            class="
                        {points.includes(point)
              ? checkedPointColor
              : uncheckedPointColor}
                        cursor-pointer font-bold border p-2 me-2 xl:me-1 mb-2
                        xl:mb-0 rounded inline-block {!votingLocked
              ? 'opacity-25 cursor-not-allowed'
              : 'cursor-pointer'}"
          >
            <input
              type="checkbox"
              bind:group="{points}"
              value="{point}"
              class="hidden"
              disabled="{!votingLocked}"
            />
            {point}
          </label>
        {/each}
      </div>
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="averageRounding"
      >
        {$LL.pointAverageRounding()}
      </label>
      <div class="relative">
        <SelectInput
          bind:value="{pointAverageRounding}"
          id="averageRounding"
          name="averageRounding"
        >
          {#each allowedPointAverages as item}
            <option value="{item}">
              {$LL.averageRoundingOptions[item]()}
            </option>
          {/each}
        </SelectInput>
        <div
          class="pointer-events-none absolute inset-y-0 end-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-400"
        >
          <DownCarrotIcon />
        </div>
      </div>
    </div>

    <div class="mb-4">
      <Checkbox
        bind:checked="{autoFinishVoting}"
        id="autoFinishVoting"
        name="autoFinishVoting"
        disabled="{!votingLocked}"
        label="{$LL.autoFinishVotingLabel({
          friendly: AppConfig.FriendlyUIVerbs,
        })}"
      />
    </div>

    <div class="mb-4">
      <Checkbox
        bind:checked="{hideVoterIdentity}"
        id="hideVoterIdentity"
        name="hideVoterIdentity"
        label="{$LL.hideVoterIdentity()}"
      />
    </div>

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
        for="leaderCode"
      >
        {$LL.leaderPasscode()}
      </label>
      <div class="control">
        <TextInput
          name="leaderCode"
          bind:value="{leaderCode}"
          placeholder="{$LL.optionalLeadercodePlaceholder()}"
          id="leaderCode"
        />
      </div>
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
        for="selectedTeam"
      >
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}{/if}
      </label>
      <SelectInput bind:value="{teamId}" id="selectedTeam" name="selectedTeam">
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value="{team.id}">
            {team.name}
          </option>
        {/each}
      </SelectInput>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

<script lang="ts">
  import SolidButton from '../SolidButton.svelte';
  import Modal from '../Modal.svelte';
  import DownCarrotIcon from '../icons/ChevronDown.svelte';
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  const allowedPointValues = AppConfig.AllowedPointValues;
  const allowedPointAverages = ['ceil', 'round', 'floor'];

  export let toggleEditBattle = () => {};
  export let handleBattleEdit = () => {};
  export let points = [];
  export let battleName = '';
  export let votingLocked = false;
  export let autoFinishVoting = true;
  export let pointAverageRounding = 'ceil';
  export let joinCode = '';
  export let leaderCode = '';
  export let hideVoterIdentity = false;

  let checkedPointColor =
    'border-green-500 bg-green-100 text-green-600 dark:bg-gray-900 dark:text-lime-500 dark:border-lime-500';
  let uncheckedPointColor =
    'border-gray-300 bg-white dark:bg-gray-900 dark:border-gray-600 dark:text-gray-300';

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
    };

    handleBattleEdit(battle);
  }
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
        <input
          name="battleName"
          bind:value="{battleName}"
          placeholder="{$LL.battleNamePlaceholder({
            friendly: AppConfig.FriendlyUIVerbs,
          })}"
          class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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
        <select
          bind:value="{pointAverageRounding}"
          class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pe-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
          id="averageRounding"
          name="averageRounding"
        >
          {#each allowedPointAverages as item}
            <option value="{item}">
              {$LL.averageRoundingOptions[item]()}
            </option>
          {/each}
        </select>
        <div
          class="pointer-events-none absolute inset-y-0 end-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-400"
        >
          <DownCarrotIcon />
        </div>
      </div>
    </div>

    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
        <input
          type="checkbox"
          bind:checked="{autoFinishVoting}"
          id="autoFinishVoting"
          name="autoFinishVoting"
          disabled="{!votingLocked}"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        {$LL.autoFinishVotingLabel({
          friendly: AppConfig.FriendlyUIVerbs,
        })}
      </label>
    </div>

    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
        <input
          type="checkbox"
          bind:checked="{hideVoterIdentity}"
          id="hideVoterIdentity"
          name="hideVoterIdentity"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        Hide Voter Identity
      </label>
    </div>

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
          class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500
                dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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
        <input
          name="leaderCode"
          bind:value="{leaderCode}"
          placeholder="{$LL.optionalLeadercodePlaceholder()}"
          class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500
                dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          id="leaderCode"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

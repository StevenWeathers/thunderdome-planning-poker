<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../SolidButton.svelte';
  import HollowButton from '../HollowButton.svelte';
  import JiraImport from './JiraImport.svelte';
  import DownCarrotIcon from '../icons/ChevronDown.svelte';
  import { warrior as user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import CsvImport from './CsvImport.svelte';

  export let notifications;
  export let eventTag;
  export let router;
  export let xfetch;
  export let apiPrefix = '/api';

  const allowedPointValues = AppConfig.AllowedPointValues;
  const allowedPointAverages = ['ceil', 'round', 'floor'];

  let points = AppConfig.DefaultPointValues;
  let battleName = '';
  let plans = [];
  let autoFinishVoting = true;
  let pointAverageRounding = 'ceil';
  let joinCode = '';
  let leaderCode = '';
  let selectedTeam = '';
  let teams = [];
  let hideVoterIdentity = false;

  let checkedPointColor =
    'border-green-500 bg-green-100 text-green-600 dark:bg-gray-900 dark:text-lime-500 dark:border-lime-500';
  let uncheckedPointColor =
    'border-gray-300 bg-white dark:bg-gray-900 dark:border-gray-600 dark:text-gray-300';

  function addPlan() {
    plans.unshift({
      name: '',
      type: $LL.planTypeStory(),
      referenceId: '',
      link: '',
      description: '',
      acceptanceCriteria: '',
    });
    plans = plans;
  }

  function handlePlanImport(newPlan) {
    const plan = {
      name: newPlan.planName,
      type: newPlan.type,
      referenceId: newPlan.referenceId,
      link: newPlan.link,
      description: newPlan.description,
      acceptanceCriteria: newPlan.acceptanceCriteria,
    };
    plans.unshift(plan);
    plans = plans;
  }

  function removePlan(i) {
    return function remove() {
      plans.splice(i, 1);
      plans = plans;
    };
  }

  function createBattle(e) {
    e.preventDefault();
    let endpoint = `${apiPrefix}/users/${$user.id}/battles`;

    const pointValuesAllowed = allowedPointValues.filter(pv => {
      return points.includes(pv);
    });

    const body = {
      name: battleName,
      pointValuesAllowed,
      plans,
      autoFinishVoting,
      pointAverageRounding,
      hideVoterIdentity,
      joinCode,
      leaderCode,
    };

    if (selectedTeam !== '') {
      endpoint = `/api/teams/${selectedTeam}/users/${$user.id}/battles`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function (result) {
        const battle = result.data;
        eventTag('create_battle', 'engagement', 'success', () => {
          router.route(`${appRoutes.game}/${battle.id}`);
        });
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
            notifications.danger(
              `${$LL.createBattleError({
                friendly: AppConfig.FriendlyUIVerbs,
              })} : ${result.error}`,
            );
          });
        } else {
          notifications.danger(
            $LL.createBattleError({
              friendly: AppConfig.FriendlyUIVerbs,
            }),
          );
        }
        eventTag('create_battle', 'engagement', 'failure');
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

<form on:submit="{createBattle}" name="createBattle">
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

  {#if apiPrefix === '/api'}
    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
        for="selectedTeam"
      >
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}{/if}
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
          <DownCarrotIcon />
        </div>
      </div>
    </div>
  {/if}

  <div class="mb-4">
    <h3 class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
      {$LL.pointValuesAllowed()}
    </h3>
    <div class="control relative -me-2 md:-me-1">
      {#each allowedPointValues as point, pi}
        <label
          class="{points.includes(point)
            ? checkedPointColor
            : uncheckedPointColor}
                    cursor-pointer font-bold border p-2 me-2 xl:me-1 mb-2
                    xl:mb-0 rounded inline-block"
        >
          <input
            type="checkbox"
            bind:group="{points}"
            value="{point}"
            class="hidden"
          />
          {point}
        </label>
      {/each}
    </div>
  </div>

  <div class="mb-4">
    <h3 class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
      {$LL.plans({ friendly: AppConfig.FriendlyUIVerbs })}
    </h3>
    <div class="control mb-4">
      <JiraImport
        handlePlanAdd="{handlePlanImport}"
        notifications="{notifications}"
        eventTag="{eventTag}"
      />
      <HollowButton onClick="{addPlan}">
        {$LL.addPlan({ friendly: AppConfig.FriendlyUIVerbs })}
      </HollowButton>
    </div>
    <div class="control mb-4">
      <CsvImport
        handlePlanAdd="{handlePlanImport}"
        notifications="{notifications}"
        eventTag="{eventTag}"
      />
    </div>

    {#each plans as plan, i}
      <div class="flex flex-wrap mb-2">
        <div class="w-3/4">
          <input
            type="text"
            bind:value="{plan.name}"
            placeholder="{$LL.planNamePlaceholder({
              friendly: AppConfig.FriendlyUIVerbs,
            })}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            required
          />
        </div>
        <div class="w-1/4">
          <div class="ps-2">
            <HollowButton onClick="{removePlan(i)}" color="red">
              {$LL.remove()}
            </HollowButton>
          </div>
        </div>
      </div>
    {/each}
  </div>

  <div class="mb-4">
    <label
      class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
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
        class="w-4 h-4 dark:accent-lime-400 me-1"
      />
      {$LL.autoFinishVotingLabel({ friendly: AppConfig.FriendlyUIVerbs })}
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
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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
    <SolidButton type="submit"
      >{$LL.battleCreate({
        friendly: AppConfig.FriendlyUIVerbs,
      })}</SolidButton
    >
  </div>
</form>

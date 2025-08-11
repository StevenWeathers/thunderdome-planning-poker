<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import ImportModal from './ImportModal.svelte';
  import SelectWithSubtext from '../forms/SelectWithSubtext.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { Crown, Lock } from 'lucide-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    notifications: NotificationService;
    router: any;
    xfetch: ApiClient;
    apiPrefix?: string;
  }

  let {
    notifications,
    router,
    xfetch,
    apiPrefix = '/api'
  }: Props = $props();

  const allowedPointAverages = ['ceil', 'round', 'floor'];

  let allowedPointValues = $state([]);
  let points = $state([]);
  let plans = $state([]);
  let teams = $state([]);
  let publicEstimationScales = [];
  let teamEstimationScales = [];
  let organizationEstimationScales = [];
  let estimateScales = $state([]);
  let selectedEstimationScale = $state('');
  let defaultSettings = {
    battleName: '',
    autoFinishVoting: true,
    pointAverageRounding: AppConfig.DefaultPointAverageRounding || 'ceil',
    hideVoterIdentity: false,
    joinCode: '',
    leaderCode: '',
    selectedTeam: '',
  };
  let pokerSettings = $state({ ...defaultSettings });
  let teamPokerSettings = {};
  let departmentPokerSettings = {};
  let orgPokerSettings = {};

  /** @type {TextInput} */
  let battleNameTextInput = $state();

  let checkedPointColor =
    'border-green-500 bg-green-50 text-green-700 dark:bg-lime-50 dark:text-lime-700 dark:border-lime-500';
  let uncheckedPointColor =
    'border-gray-300 bg-white text-gray-700 dark:bg-gray-900 dark:border-gray-500 dark:text-gray-300';

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

    if (selectedEstimationScale === '' || allowedPointValues.length === 0) {
      notifications.danger(
        'Must select an estimation scale and allowed point values.',
      );
      return;
    }

    const pointValuesAllowed = allowedPointValues.filter(pv => {
      return points.includes(pv);
    });

    const body = {
      name: pokerSettings.battleName,
      pointValuesAllowed,
      plans,
      autoFinishVoting: pokerSettings.autoFinishVoting,
      pointAverageRounding: pokerSettings.pointAverageRounding,
      hideVoterIdentity: pokerSettings.hideVoterIdentity,
      joinCode: pokerSettings.joinCode,
      leaderCode: pokerSettings.leaderCode,
      estimationScaleId: selectedEstimationScale,
    };

    if (pokerSettings.selectedTeam !== '') {
      endpoint = `/api/teams/${pokerSettings.selectedTeam}/users/${$user.id}/battles`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function (result) {
        const battle = result.data;
        router.route(`${appRoutes.game}/${battle.id}`);
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
            notifications.danger(
              `${$LL.createBattleError()} : ${result.error}`,
            );
          });
        } else {
          notifications.danger($LL.createBattleError());
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

  function getPublicEstimateScales() {
    xfetch(`/api/estimation-scales/public`)
      .then(res => res.json())
      .then(function (result) {
        publicEstimationScales = result.data;
        combineEstimationScales();
      })
      .catch(function () {
        notifications.danger('Failed to get public estimation scales');
      });
  }

  function getPrivateEstimationScales() {
    teamEstimationScales = [];
    organizationEstimationScales = [];
    combineEstimationScales();

    // don't get private scales if a team isn't selected
    if (pokerSettings.selectedTeam === '') {
      return;
    }
    const team = teams.find(t => t.id === pokerSettings.selectedTeam);
    // if subscriptions are enabled and the team (or its parent org) isn't subscribed
    // don't attempt to get private scales
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

    xfetch(`${teamPrefix}/teams/${team.id}/estimation-scales`)
      .then(res => res.json())
      .then(function (result) {
        teamEstimationScales = result.data;
        combineEstimationScales();
      })
      .catch(function () {
        notifications.danger('Failed to get team estimation scales');
      });

    if (team.organization_id !== '') {
      xfetch(`${orgPrefix}/estimation-scales`)
        .then(res => res.json())
        .then(function (result) {
          organizationEstimationScales = result.data;
          combineEstimationScales();
        })
        .catch(function () {
          notifications.danger('Failed to get organization estimation scales');
        });
    }
  }

  const combineEstimationScales = () => {
    // scales priority order (Team -> Organization -> Public)
    let defaultFound = false;
    estimateScales = [
      ...teamEstimationScales,
      ...organizationEstimationScales,
      ...publicEstimationScales,
    ];

    estimateScales.map(scale => {
      // Find default scale with priority order (Team -> Organization -> Public)
      if (!defaultFound && scale.defaultScale) {
        allowedPointValues = scale.values;
        points = scale.values;
        selectedEstimationScale = scale.id;
        defaultFound = true;
      }
    });
  };

  function getCustomPokerSettings() {
    teamPokerSettings = {};
    departmentPokerSettings = {};
    orgPokerSettings = {};
    combinePokerSettings();

    // don't get custom poker settings if a team isn't selected
    if (pokerSettings.selectedTeam === '') {
      return;
    }
    const team = teams.find(t => t.id === pokerSettings.selectedTeam);
    // if subscriptions are enabled and the team (or its parent org) isn't subscribed
    // don't attempt to get poker settings
    if (
      AppConfig.SubscriptionsEnabled &&
      !validateUserIsAdmin($user) &&
      !team.subscribed
    ) {
      return;
    }

    xfetch(`/api/teams/${team.id}/poker-settings`)
      .then(res => res.json())
      .then(function (result) {
        if (!result.data) {
          return;
        }
        teamPokerSettings = result.data;
        combinePokerSettings();
      })
      .catch(function () {
        notifications.danger('Failed to get team poker settings');
      });

    if (team.organization_id !== '') {
      xfetch(`/api/organizations/${team.organization_id}/poker-settings`)
        .then(res => res.json())
        .then(function (result) {
          if (!result.data) {
            return;
          }
          orgPokerSettings = result.data;
          combinePokerSettings();
        })
        .catch(function () {
          notifications.danger('Failed to get organization poker settings');
        });
    }

    if (team.department_id !== '') {
      xfetch(
        `/api/organizations/${team.organization_id}/departments/${team.department_id}/poker-settings`,
      )
        .then(res => res.json())
        .then(function (result) {
          if (!result.data) {
            return;
          }
          departmentPokerSettings = result.data;
          combinePokerSettings();
        })
        .catch(function () {
          notifications.danger('Failed to get department poker settings');
        });
    }
  }

  const combinePokerSettings = () => {
    // settings priority order (Team -> Department -> Organization -> Default)
    pokerSettings = {
      ...pokerSettings,
      ...orgPokerSettings,
      ...departmentPokerSettings,
      ...teamPokerSettings,
    };
  };

  const updatePointValues = event => {
    const scale = event.detail;
    selectedEstimationScale = scale.id;
    allowedPointValues = scale.values;
    points = scale.values;
  };

  let showImport = $state(false);

  const toggleImport = () => {
    showImport = !showImport;
  };

  function handleTeamChange() {
    // get custom poker settings before estimation scales
    getCustomPokerSettings();
    getPrivateEstimationScales();
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.register);
    }
    getTeams();
    getPublicEstimateScales();

    // Focus the battle name input field
    battleNameTextInput.focus();
  });
</script>

<form onsubmit={createBattle} name="createBattle">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="battleName"
    >
      {$LL.battleName()}
    </label>
    <div class="control">
      <TextInput
        name="battleName"
        bind:this={battleNameTextInput}
        bind:value="{pokerSettings.battleName}"
        placeholder={$LL.battleNamePlaceholder()}
        id="battleName"
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
        {#if !AppConfig.RequireTeams}{$LL.optional()}{/if}
      </label>
      <SelectInput
        bind:value="{pokerSettings.selectedTeam}"
        on:change="{handleTeamChange}"
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
      Estimation Scale
    </div>
    <SelectWithSubtext
      on:change={updatePointValues}
      items={estimateScales}
      label="Select an estimation scale..."
      selectedItemId={selectedEstimationScale}
      itemType="estimation_scale"
    />
  </div>

  <div class="mb-4">
    <h3 class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2">
      {$LL.pointValuesAllowed()}
    </h3>
    <div class="control relative flex flex-wrap gap-1">
      {#each allowedPointValues as point, pi}
        <label
          class="{points.includes(point)
            ? checkedPointColor
            : uncheckedPointColor}
                    cursor-pointer font-bold border p-2 rounded inline-block"
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
      {$LL.plans()}
    </h3>
    <div class="control mb-4">
      <HollowButton onClick={toggleImport} color="blue">
        {$LL.importPlans()}
      </HollowButton>
      <HollowButton onClick={addPlan}>
        {$LL.addPlan()}
      </HollowButton>
      {#if showImport}
        <ImportModal
          notifications={notifications}
          toggleImport={toggleImport}
          handlePlanAdd={handlePlanImport}
          xfetch={xfetch}
        />
      {/if}
    </div>

    {#each plans as plan, i}
      <div class="flex flex-wrap mb-2">
        <div class="w-3/4">
          <TextInput
            bind:value="{plan.name}"
            placeholder={$LL.planNamePlaceholder()}
            required
          />
        </div>
        <div class="w-1/4">
          <div class="ps-2">
            <HollowButton onClick={removePlan(i)} color="red">
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
    <SelectInput
      bind:value="{pokerSettings.pointAverageRounding}"
      id="averageRounding"
      name="averageRounding"
    >
      {#each allowedPointAverages as item}
        <option value="{item}">
          {$LL.averageRoundingOptions[item]()}
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked="{pokerSettings.autoFinishVoting}"
      id="autoFinishVoting"
      name="autoFinishVoting"
      label={$LL.autoFinishVotingLabel()}
    />
  </div>

  <div class="mb-4">
    <Checkbox
      bind:checked="{pokerSettings.hideVoterIdentity}"
      id="hideVoterIdentity"
      name="hideVoterIdentity"
      label={$LL.hideVoterIdentity()}
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
        bind:value="{pokerSettings.joinCode}"
        placeholder={$LL.optionalPasscodePlaceholder()}
        id="joinCode"
        icon={Lock}
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="leaderCode"
    >
      {$LL.facilitatorCodeOptional()}
    </label>
    <div class="control">
      <TextInput
        name="leaderCode"
        bind:value="{pokerSettings.leaderCode}"
        placeholder={$LL.facilitatorCodePlaceholder()}
        id="leaderCode"
        icon={Crown}
      />
    </div>
  </div>

  <div class="text-right">
    <SolidButton type="submit">{$LL.battleCreate()}</SolidButton>
  </div>
</form>

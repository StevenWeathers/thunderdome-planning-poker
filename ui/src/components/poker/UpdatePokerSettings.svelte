<script lang="ts">
  import { preventDefault } from 'svelte/legacy';

  import Modal from '../global/Modal.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import LL from '../../i18n/i18n-svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { Crown, Lock } from 'lucide-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import { createEventDispatcher } from 'svelte';

  interface Props {
    toggleClose?: any;
    xfetch: any;
    notifications: any;
    organizationId: any;
    teamId: any;
    departmentId: any;
    apiPrefix?: string;
    isEntityAdmin?: boolean;
    pokerSettings?: any;
  }

  let {
    toggleClose = () => {},
    xfetch,
    notifications,
    organizationId,
    teamId,
    departmentId,
    apiPrefix = '/api',
    isEntityAdmin = false,
    pokerSettings = $bindable({
    id: '',
    autoFinishVoting: true,
    pointAverageRounding: 'ceil',
    hideVoterIdentity: false,
    estimationScaleId: null,
    joinCode: '',
    facilitatorCode: '',
  })
  }: Props = $props();

  const dispatch = createEventDispatcher();
  const allowedPointAverages = ['ceil', 'round', 'floor'];

  // let estimationScales = []
  //
  // onMount(async () => {
  //   // Fetch estimation scales from the server
  //   // This is a placeholder and should be replaced with actual API call
  //   estimationScales = await fetchEstimationScales()
  // })
  //
  // async function fetchEstimationScales () {
  //   // Placeholder function to fetch estimation scales
  //   return [
  //     { id: '1', name: 'Fibonacci' },
  //     { id: '2', name: 'T-Shirt Sizes' },
  //     { id: '3', name: 'Powers of 2' }
  //   ]
  // }

  async function handleSubmit() {
    const method = pokerSettings.id !== '' ? 'PUT' : 'POST';

    const response = await xfetch(`${apiPrefix}/poker-settings`, {
      method,
      body: pokerSettings,
    });
    if (response.ok) {
      const res = await response.json();
      pokerSettings = res.data;
      dispatch('updatePokerSettings', { settings: pokerSettings });
      notifications.success(
        `Default poker settings ${method === 'PUT' ? 'updated' : 'created'}`,
      );
      toggleClose();
    } else {
      notifications.error(
        `Failed to ${
          method === 'PUT' ? 'update' : 'create'
        } default poker settings`,
      );
    }
  }
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={preventDefault(handleSubmit)} class="mt-6 space-y-6">
    <div>
      <label
        for="pointAverageRounding"
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
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

    <div>
      <Checkbox
        bind:checked="{pokerSettings.autoFinishVoting}"
        id="autoFinishVoting"
        name="autoFinishVoting"
        label={$LL.autoFinishVotingLabel()}
      />
    </div>

    <div>
      <Checkbox
        bind:checked="{pokerSettings.hideVoterIdentity}"
        id="hideVoterIdentity"
        name="hideVoterIdentity"
        label={$LL.hideVoterIdentity()}
      />
    </div>

    <!--        <div>-->
    <!--            <label for="estimationScaleId" class="block text-sm font-medium text-gray-700">-->
    <!--                Estimation Scale-->
    <!--            </label>-->
    <!--            <select-->
    <!--                    id="estimationScaleId"-->
    <!--                    bind:value={pokerSettings.estimationScaleId}-->
    <!--                    class="mt-1 block w-full ps-3 pe-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"-->
    <!--            >-->
    <!--                <option value="">Select a scale</option>-->
    <!--                {#each estimationScales as scale}-->
    <!--                    <option value={scale.id}>{scale.name}</option>-->
    <!--                {/each}-->
    <!--            </select>-->
    <!--        </div>-->

    <div>
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

    <div>
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="leaderCode"
      >
        {$LL.facilitatorCodeOptional()}
      </label>
      <div class="control">
        <TextInput
          name="leaderCode"
          bind:value="{pokerSettings.facilitatorCode}"
          placeholder={$LL.facilitatorCodePlaceholder()}
          id="leaderCode"
          icon={Crown}
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">Save Settings</SolidButton>
    </div>
  </form>
</Modal>

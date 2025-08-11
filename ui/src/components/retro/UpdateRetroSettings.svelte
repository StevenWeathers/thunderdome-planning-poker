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

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    toggleClose?: any;
    xfetch: any;
    notifications: NotificationService;
    organizationId: any;
    teamId: any;
    departmentId: any;
    apiPrefix?: string;
    isEntityAdmin?: boolean;
    retroSettings?: any;
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
    retroSettings = $bindable({
    id: '',
    maxVotes: 3,
    allowMultipleVotes: false,
    brainstormVisibility: 'visible',
    phaseTimeLimit: 0,
    phaseAutoAdvance: true,
    allowCumulativeVoting: false,
    templateId: null,
    joinCode: '',
    facilitatorCode: '',
  })
  }: Props = $props();

  const dispatch = createEventDispatcher();
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
  const maxPhaseTimeLimitMin = 59;

  async function handleSubmit() {
    const method = retroSettings.id !== '' ? 'PUT' : 'POST';

    // hack because of svelte bug where binding to dynamic number input doesn't keep type
    retroSettings.maxVotes = parseInt(`${retroSettings.maxVotes}`, 10);
    retroSettings.phaseTimeLimit = parseInt(
      `${retroSettings.phaseTimeLimit}`,
      10,
    );

    if (
      retroSettings.phaseTimeLimit > maxPhaseTimeLimitMin ||
      retroSettings.phaseTimeLimit < 0
    ) {
      notifications.danger('Phase Time Limit minutes must be between 0-59');
      return;
    }

    const response = await xfetch(`${apiPrefix}/retro-settings`, {
      method,
      body: retroSettings,
    });
    if (response.ok) {
      const res = await response.json();
      retroSettings = res.data;
      dispatch('updateRetroSettings', { settings: retroSettings });
      notifications.success(
        `Default retro settings ${method === 'PUT' ? 'updated' : 'created'}`,
      );
      toggleClose();
    } else {
      notifications.error(
        `Failed to ${
          method === 'PUT' ? 'update' : 'create'
        } default retro settings`,
      );
    }
  }
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={preventDefault(handleSubmit)} class="mt-6 space-y-6">
    <div>
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="maxVotes"
      >
        {$LL.retroMaxVotesPerUserLabel()}
      </label>
      <div class="control">
        <TextInput
          name="maxVotes"
          bind:value={retroSettings.maxVotes}
          id="maxVotes"
          type="number"
          min="1"
          max="10"
          required
        />
      </div>
    </div>

    <div>
      <Checkbox
        bind:checked={retroSettings.allowCumulativeVoting}
        id="allowCumulativeVoting"
        name="allowCumulativeVoting"
        label={$LL.allowCumulativeVotingLabel()}
      />
    </div>

    <div>
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="brainstormVisibility"
      >
        {$LL.brainstormPhaseFeedbackVisibility()}
      </label>
      <SelectInput
        bind:value="{retroSettings.brainstormVisibility}"
        id="brainstormVisibility"
        name="brainstormVisibility"
      >
        {#each brainstormVisibilityOptions as item}
          <option value={item.value}>
            {item.label}
          </option>
        {/each}
      </SelectInput>
    </div>

    <div>
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="phaseTimeLimitMin"
      >
        {$LL.retroPhaseTimeLimitMinLabel()}
      </label>
      <div class="control">
        <TextInput
          name="phaseTimeLimitMin"
          bind:value={retroSettings.phaseTimeLimit}
          id="phaseTimeLimitMin"
          type="number"
          min="0"
          max={maxPhaseTimeLimitMin}
          required
        />
      </div>
    </div>

    <div>
      <Checkbox
        bind:checked="{retroSettings.phaseAutoAdvance}"
        id="phaseAutoAdvance"
        name="phaseAutoAdvance"
        label={$LL.phaseAutoAdvanceLabel()}
      />
    </div>

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
          bind:value={retroSettings.joinCode}
          placeholder={$LL.joinCodePlaceholder()}
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
          bind:value={retroSettings.facilitatorCode}
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

<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { Crown, Lock } from 'lucide-svelte';
  import { onMount } from 'svelte';

  interface Props {
    toggleEditRetro?: any;
    handleRetroEdit?: any;
    retroName?: string;
    joinCode?: string;
    facilitatorCode?: string;
    maxVotes?: string;
    brainstormVisibility?: string;
    phaseAutoAdvance?: boolean;
    hideVotesDuringVoting?: boolean;
  }

  let {
    toggleEditRetro = () => {},
    handleRetroEdit = () => {},
    retroName = $bindable(''),
    joinCode = $bindable(''),
    facilitatorCode = $bindable(''),
    maxVotes = $bindable('3'),
    brainstormVisibility = $bindable('visible'),
    phaseAutoAdvance = $bindable(true),
    hideVotesDuringVoting = $bindable(false),
  }: Props = $props();

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

  function saveRetro(e) {
    e.preventDefault();

    const retro = {
      retroName,
      joinCode,
      facilitatorCode,
      maxVotes: parseInt(maxVotes, 10),
      brainstormVisibility,
      phase_auto_advance: phaseAutoAdvance,
      hideVotesDuringVoting,
    };

    handleRetroEdit(retro);
  }

  let focusInput: any;
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleEditRetro} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabel={$LL.modalEditRetro()}>
  <form onsubmit={saveRetro} name="createRetro">
    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="retroName">
        {$LL.retroName()}
      </label>
      <div class="control">
        <TextInput
          name="retroName"
          bind:value={retroName}
          bind:this={focusInput}
          placeholder={$LL.retroNamePlaceholder()}
          id="retroName"
          required
        />
      </div>
    </div>

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
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="facilitatorCode">
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

    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="maxVotes">
        {$LL.retroMaxVotesPerUserLabel()}
      </label>
      <div class="control">
        <TextInput name="retroName" bind:value={maxVotes} id="maxVotes" type="number" min="1" max="10" required />
      </div>
    </div>

    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="brainstormVisibility">
        {$LL.brainstormPhaseFeedbackVisibility()}
      </label>
      <SelectInput bind:value={brainstormVisibility} id="brainstormVisibility" name="brainstormVisibility">
        {#each brainstormVisibilityOptions as item}
          <option value={item.value}>
            {item.label}
          </option>
        {/each}
      </SelectInput>
    </div>

    <div class="mb-4">
      <Checkbox
        bind:checked={phaseAutoAdvance}
        id="phaseAutoAdvance"
        name="phaseAutoAdvance"
        label={$LL.phaseAutoAdvanceLabel()}
      />
    </div>

    <div class="mb-4">
      <Checkbox
        bind:checked={hideVotesDuringVoting}
        id="hideVotesDuringVoting"
        name="hideVotesDuringVoting"
        label={`Hide Votes During Voting Phase`}
      />
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

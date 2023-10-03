<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../global/TextInput.svelte';
  import SelectInput from '../global/SelectInput.svelte';

  export let toggleEditRetro = () => {};
  export let handleRetroEdit = () => {};
  export let retroName = '';
  export let joinCode = '';
  export let facilitatorCode = '';
  export let maxVotes = 3;
  export let brainstormVisibility = 'visible';

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
      maxVotes,
      brainstormVisibility,
    };

    handleRetroEdit(retro);
  }
</script>

<Modal closeModal="{toggleEditRetro}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <form on:submit="{saveRetro}" name="createRetro">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="retroName"
      >
        {$LL.retroName()}
      </label>
      <div class="control">
        <TextInput
          name="retroName"
          bind:value="{retroName}"
          placeholder="{$LL.retroNamePlaceholder()}"
          id="retroName"
          required
        />
      </div>
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
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="facilitatorCode"
      >
        {$LL.facilitatorCodeOptional()}
      </label>
      <div class="control">
        <TextInput
          name="facilitatorCode"
          bind:value="{facilitatorCode}"
          placeholder="{$LL.facilitatorCodePlaceholder()}"
          id="facilitatorCode"
        />
      </div>
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="maxVotes"
      >
        {$LL.retroMaxVotesPerUserLabel()}
      </label>
      <div class="control">
        <TextInput
          name="retroName"
          bind:value="{maxVotes}"
          id="maxVotes"
          type="number"
          min="1"
          max="10"
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="brainstormVisibility"
      >
        {$LL.brainstormPhaseFeedbackVisibility()}
      </label>
      <SelectInput
        bind:value="{brainstormVisibility}"
        id="brainstormVisibility"
        name="brainstormVisibility"
      >
        {#each brainstormVisibilityOptions as item}
          <option value="{item.value}">
            {item.label}
          </option>
        {/each}
      </SelectInput>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

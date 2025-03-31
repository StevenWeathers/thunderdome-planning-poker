<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  interface Props {
    toggleCreate?: any;
    handleCreate?: any;
    organizationId: any;
    teamId: any;
    departmentId: any;
    apiPrefix?: string;
    xfetch: any;
    notifications: any;
  }

  let {
    toggleCreate = () => {},
    handleCreate = () => {},
    organizationId,
    teamId,
    departmentId,
    apiPrefix = '/api',
    xfetch,
    notifications
  }: Props = $props();

  let name = $state('');
  let description = $state('');
  let scaleType = $state('custom');
  let values: string[] = $state([]);
  let isPublic = $state(false);
  let defaultScale = $state(false);

  const scaleTypes = [
    'fibonacci',
    'modified_fibonacci',
    't_shirt',
    'powers_of_two',
    'thunderdome_default',
    'custom',
  ];

  function toggleClose() {
    toggleCreate();
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    const body = {
      name,
      description,
      values,
      defaultScale,
    };

    if (isAdmin) {
      body['isPublic'] = isPublic;
      body['scaleType'] = scaleType;
    }

    xfetch(`${apiPrefix}/estimation-scales`, {
      body,
    })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.estimationScaleCreateSuccess());
        handleCreate();
      })
      .catch(() => {
        notifications.danger($LL.estimationScaleCreateError());
      });
  }

  function handleValuesChange(e: Event) {
    const target = e.target as HTMLInputElement;
    const inputValue = target.value;

    // Regex to validate the entire input string
    const validationRegex =
      /^[\p{L}\p{N}\p{Emoji}\p{Emoji_Component}_?]*[\p{L}\p{N}\p{Emoji}\p{Emoji_Component}?]+( *, *[\p{L}\p{N}\p{Emoji}\p{Emoji_Component}_?]*[\p{L}\p{N}\p{Emoji}\p{Emoji_Component}?]+)*$/u;

    if (validationRegex.test(inputValue)) {
      // If the input is valid, split it into values
      values = inputValue
        .split(',')
        .map(v => v.trim())
        .filter(v => v !== '');

      // Update the input field with the cleaned values
      target.value = values.join(', ');
    } else {
      // If the input is not valid, don't update the values array
      // Optionally, you could provide feedback to the user here
      notifications.danger(
        'Invalid input. Please use only letters, numbers, and commas.',
      );
    }
  }

  let createDisabled = $derived(name === '' || scaleType === '' || values.length === 0);
  let isAdmin = $derived(validateUserIsAdmin($user));
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={onSubmit} name="createEstimationScale">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="scaleName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value={name}
        placeholder={$LL.estimationScaleNamePlaceholder()}
        id="scaleName"
        name="scaleName"
        required
      />
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="scaleDescription"
      >
        {$LL.description()}
      </label>
      <TextInput
        bind:value={description}
        placeholder={$LL.estimationScaleDescriptionPlaceholder()}
        id="scaleDescription"
        name="scaleDescription"
      />
    </div>

    {#if isAdmin && !organizationId && !teamId}
      <div class="mb-4">
        <label class="block font-bold mb-2 dark:text-gray-400" for="scaleType">
          {$LL.scaleType()}
        </label>
        <SelectInput
          name="scaleType"
          id="scaleType"
          bind:value="{scaleType}"
          required
        >
          <option value="" disabled>
            {$LL.estimationScaleTypePlaceholder()}
          </option>
          {#each scaleTypes as type}
            <option value={type}>{type}</option>
          {/each}
        </SelectInput>
      </div>
    {/if}

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="scaleValues"
      >
        {$LL.scaleValues()}
      </label>
      <TextInput
        value={values.join(', ')}
        on:input={handleValuesChange}
        placeholder={$LL.estimationScaleValuesPlaceholder()}
        id="scaleValues"
        name="scaleValues"
        required
      />
      <p class="text-sm text-gray-500 mt-1">
        {$LL.estimationScaleValuesHelp()}
      </p>
    </div>

    {#if isAdmin && !organizationId && !teamId}
      <div class="mb-4">
        <Checkbox
          bind:checked={isPublic}
          id="isPublic"
          name="isPublic"
          label={$LL.estimationScaleIsPublic()}
        />
      </div>
    {/if}

    <div class="mb-4">
      <Checkbox
        bind:checked={defaultScale}
        id="defaultScale"
        name="defaultScale"
        label={$LL.estimationScaleDefault()}
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {$LL.estimationScaleSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

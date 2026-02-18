<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import { onMount } from 'svelte';

  interface Props {
    toggleUpdate?: any;
    handleUpdate?: any;
    organizationId?: any;
    teamId?: any;
    departmentId?: any;
    apiPrefix: any;
    xfetch: ApiClient;
    notifications: NotificationService;
    scaleId?: string;
    name?: string;
    description?: string;
    scaleType?: string;
    values?: string[];
    isPublic?: boolean;
    defaultScale?: boolean;
  }

  let {
    toggleUpdate = () => {},
    handleUpdate = () => {},
    organizationId = null,
    teamId = null,
    departmentId = null,
    apiPrefix,
    xfetch,
    notifications,
    scaleId = '',
    name = $bindable(''),
    description = $bindable(''),
    scaleType = $bindable('custom'),
    values = $bindable([]),
    isPublic = $bindable(false),
    defaultScale = $bindable(false),
  }: Props = $props();

  const scaleTypes = ['fibonacci', 'modified_fibonacci', 't_shirt', 'powers_of_two', 'thunderdome_default', 'custom'];

  function toggleClose() {
    toggleUpdate();
  }

  function onSubmit(e: Event) {
    e.preventDefault();

    const body: {
      name: string;
      description: string;
      values: string[];
      defaultScale: boolean;
      isPublic?: boolean;
      scaleType?: string;
    } = {
      name,
      description,
      values,
      defaultScale,
    };

    if (isAdmin) {
      body.isPublic = isPublic;
      body.scaleType = scaleType;
    }

    xfetch(`${apiPrefix}/estimation-scales/${scaleId}`, {
      body,
      method: 'put',
    })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.estimationScaleUpdateSuccess());
        handleUpdate();
      })
      .catch(() => {
        notifications.danger($LL.estimationScaleUpdateError());
      });
  }

  function handleValuesChange(e: Event) {
    const target = e.target as HTMLInputElement;
    const inputValue = target.value;

    // Regex to allow letters, numbers, commas, spaces, underscores, and emoji
    const allowedCharsRegex = /^[\p{L}\p{N}, _\p{Emoji}\p{Emoji_Component}]*$/u;

    if (allowedCharsRegex.test(inputValue)) {
      values = inputValue
        .split(',')
        .map(v => v.trim())
        .filter(v => v !== '');
    } else {
      notifications.danger('Invalid input. Please use only letters, numbers, emojis, commas, spaces, and underscores.');
    }
  }

  let updateDisabled = $derived(name === '' || scaleType === '' || values.length === 0);
  let isAdmin = $derived(validateUserIsAdmin($user));

  let focusInput: any;
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleClose} ariaLabel={$LL.modalUpdateEstimationScale()}>
  <form onsubmit={onSubmit} name="updateEstimationScale">
    <div class="mb-4">
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="scaleName">
        {$LL.name()}
      </label>
      <TextInput
        bind:value={name}
        bind:this={focusInput}
        placeholder={$LL.estimationScaleNamePlaceholder()}
        id="scaleName"
        name="scaleName"
        required
      />
    </div>

    <div class="mb-4">
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="scaleDescription">
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
        <SelectInput name="scaleType" id="scaleType" bind:value={scaleType} required>
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
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="scaleValues">
        {$LL.scaleValues()}
      </label>
      <TextInput
        value={values.join(', ')}
        oninput={handleValuesChange}
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
        <Checkbox bind:checked={isPublic} id="isPublic" name="isPublic" label={$LL.estimationScaleIsPublic()} />
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
        <SolidButton type="submit" disabled={updateDisabled}>
          {$LL.estimationScaleSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

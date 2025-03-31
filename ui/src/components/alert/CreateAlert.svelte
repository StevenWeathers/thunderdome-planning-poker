<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';

  interface Props {
    toggleCreate?: any;
    handleCreate?: any;
    toggleUpdate?: any;
    handleUpdate?: any;
    alertId?: string;
    alertName?: string;
    alertType?: string;
    content?: string;
    active?: boolean;
    registeredOnly?: boolean;
    allowDismiss?: boolean;
  }

  let {
    toggleCreate = () => {},
    handleCreate = () => {},
    toggleUpdate = () => {},
    handleUpdate = () => {},
    alertId = '',
    alertName = $bindable(''),
    alertType = $bindable(''),
    content = $bindable(''),
    active = $bindable(true),
    registeredOnly = $bindable(false),
    allowDismiss = $bindable(true)
  }: Props = $props();

  const alertTypes = ['ERROR', 'INFO', 'NEW', 'SUCCESS', 'WARNING'];

  function toggleClose() {
    if (alertId != '') {
      toggleUpdate();
    } else {
      toggleCreate();
    }
  }

  function onSubmit(e) {
    e.preventDefault();

    const body = {
      name: alertName,
      type: alertType,
      content,
      active,
      registeredOnly,
      allowDismiss,
    };

    if (alertId !== '') {
      handleUpdate(alertId, body);
    } else {
      handleCreate(body);
    }
  }

  let createDisabled = $derived(alertName === '' || alertType === '' || content === '');
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={onSubmit} name="createAlert">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="alertName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value={alertName}
        placeholder={$LL.alertNamePlaceholder()}
        id="alertName"
        name="alertName"
        required
      />
    </div>

    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="alertType">
        {$LL.type()}
      </label>
      <SelectInput
        name="alertType"
        id="alertType"
        bind:value={alertType}
        required
      >
        <option value="" disabled>
          {$LL.alertTypePlaceholder()}
        </option>
        {#each alertTypes as aType}
          <option value={aType}>{aType}</option>
        {/each}
      </SelectInput>
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="alertContent"
      >
        {$LL.alertContent()}
      </label>
      <TextInput
        bind:value={content}
        placeholder={$LL.alertContentPlaceholder()}
        id="alertContent"
        name="alertContent"
        required
      />
    </div>

    <div class="mb-4">
      <Checkbox
        bind:checked={active}
        id="active"
        name="active"
        label={$LL.active()}
      />
    </div>
    <div class="mb-4">
      <Checkbox
        bind:checked={registeredOnly}
        id="registeredOnly"
        name="registeredOnly"
        label={$LL.alertRegisteredOnly()}
      />
    </div>
    <div class="mb-4">
      <Checkbox
        bind:checked={allowDismiss}
        id="allowDismiss"
        name="allowDismiss"
        label={$LL.alertAllowDismiss()}
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {$LL.alertSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

<script lang="ts">
  import Modal from '../Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../SolidButton.svelte';
  import TextInput from '../TextInput.svelte';
  import SelectInput from '../SelectInput.svelte';

  export let toggleCreate = () => {};
  export let handleCreate = () => {};
  export let toggleUpdate = () => {};
  export let handleUpdate = () => {};
  export let alertId = '';
  export let alertName = '';
  export let alertType = '';
  export let content = '';
  export let active = true;
  export let registeredOnly = false;
  export let allowDismiss = true;

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

  $: createDisabled = alertName === '' || alertType === '' || content === '';
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{onSubmit}" name="createAlert">
    <div class="mb-4">
      <label
        class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
        for="alertName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value="{alertName}"
        placeholder="{$LL.alertNamePlaceholder()}"
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
        bind:value="{alertType}"
        required
      >
        <option value="" disabled>
          {$LL.alertTypePlaceholder()}
        </option>
        {#each alertTypes as aType}
          <option value="{aType}">{aType}</option>
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
        bind:value="{content}"
        placeholder="{$LL.alertContentPlaceholder()}"
        id="alertContent"
        name="alertContent"
        required
      />
    </div>

    <div class="mb-4">
      <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
        <input
          type="checkbox"
          bind:checked="{active}"
          id="active"
          name="active"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        {$LL.active()}
      </label>
    </div>
    <div class="mb-4">
      <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
        <input
          type="checkbox"
          bind:checked="{registeredOnly}"
          id="registeredOnly"
          name="registeredOnly"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        {$LL.alertRegisteredOnly()}
      </label>
    </div>
    <div class="mb-4">
      <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
        <input
          type="checkbox"
          bind:checked="{allowDismiss}"
          id="allowDismiss"
          name="allowDismiss"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        {$LL.alertAllowDismiss()}
      </label>
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled="{createDisabled}">
          {$LL.alertSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

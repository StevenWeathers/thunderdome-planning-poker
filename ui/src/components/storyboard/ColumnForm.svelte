<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { onMount } from 'svelte';

  interface Props {
    toggleColumnEdit?: any;
    handleColumnRevision?: any;
    column?: any;
  }

  let {
    toggleColumnEdit = () => {},
    handleColumnRevision = () => {},
    column = $bindable({
      id: '',
      name: '',
      personas: [],
    }),
  }: Props = $props();

  let focusInput: any = $state();
  let columnName = $state(column.name);

  $effect(() => {
    columnName = column.name;
  });

  function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    const c = {
      id: column.id,
      name: columnName,
    };

    handleColumnRevision(c);
    toggleColumnEdit();
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleColumnEdit} ariaLabel={$LL.modalStoryboardColumnSettings()}>
  <form onsubmit={handleSubmit} name="addColumn">
    <div class="mb-4">
      <label class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2" for="columnName">
        Column Name
      </label>
      <TextInput
        id="columnName"
        bind:value={columnName}
        placeholder="Enter a column name"
        name="columnName"
        bind:this={focusInput}
      />
    </div>
    <div class="flex justify-end gap-2">
      <SolidButton type="submit">Save</SolidButton>
    </div>
  </form>
</Modal>

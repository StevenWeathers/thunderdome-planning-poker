<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import TextInput from '../global/TextInput.svelte';

  export let toggleColumnEdit = () => {};
  export let handleColumnRevision = () => {};
  export let deleteColumn = () => () => {};

  export let column = {
    id: '',
    name: '',
  };

  function handleSubmit(event) {
    event.preventDefault();

    handleColumnRevision(column);
    toggleColumnEdit();
  }
</script>

<Modal closeModal="{toggleColumnEdit}">
  <form on:submit="{handleSubmit}" name="addColumn">
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="columnName"
      >
        Column Name
      </label>
      <TextInput
        id="columnName"
        bind:value="{column.name}"
        placeholder="Enter a column name"
        name="columnName"
      />
    </div>
    <div class="flex">
      <div class="md:w-1/2 text-left">
        <HollowButton color="red" onClick="{deleteColumn(column.id)}">
          Delete Column
        </HollowButton>
      </div>
      <div class="md:w-1/2 text-right">
        <SolidButton type="submit">Save</SolidButton>
      </div>
    </div>
  </form>
</Modal>

<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';


  interface Props {
    toggleCreate?: any;
    handleCreate?: any;
    departmentName?: string;
  }

  let { toggleCreate = () => {}, handleCreate = () => {}, departmentName = $bindable('') }: Props = $props();

  function onSubmit(e) {
    e.preventDefault();

    handleCreate(departmentName);
  }

  let createDisabled = $derived(departmentName === '');
</script>

<Modal closeModal={toggleCreate}>
  <form onsubmit={onSubmit} name="createDepartment">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="departmentName"
      >
        {$LL.departmentName()}
      </label>
      <TextInput
        bind:value="{departmentName}"
        placeholder={$LL.departmentNamePlaceholder()}
        id="departmentName"
        name="departmentName"
        required
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {$LL.departmentSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

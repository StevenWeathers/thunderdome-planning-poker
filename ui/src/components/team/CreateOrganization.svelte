<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';


  interface Props {
    toggleCreate?: any;
    handleCreate?: any;
    organizationName?: string;
  }

  let { toggleCreate = () => {}, handleCreate = () => {}, organizationName = $bindable('') }: Props = $props();

  function onSubmit(e) {
    e.preventDefault();

    handleCreate(organizationName);
  }

  let createDisabled = $derived(organizationName === '');
</script>

<Modal closeModal={toggleCreate}>
  <form onsubmit={onSubmit} name="createOrganization">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="organizationName"
      >
        {$LL.organizationName()}
      </label>
      <TextInput
        bind:value="{organizationName}"
        placeholder={$LL.organizationNamePlaceholder()}
        id="organizationName"
        name="organizationName"
        required
      />
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={createDisabled}>
          {$LL.organizationSave()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

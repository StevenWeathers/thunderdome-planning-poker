<script lang="ts">
  import Modal from './global/Modal.svelte';
  import LL from '../i18n/i18n-svelte';
  import SolidButton from './global/SolidButton.svelte';
  import TextInput from './forms/TextInput.svelte';

  export let toggleBecomeFacilitator = () => {};
  export let handleBecomeFacilitator = () => {};

  let facilitatorCode = '';

  function handleSubmit(e) {
    e.preventDefault();

    handleBecomeFacilitator(facilitatorCode);
  }
</script>

<Modal
  closeModal="{toggleBecomeFacilitator}"
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
  <form on:submit="{handleSubmit}" name="becomeLeader">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="facilitatorCode"
      >
        {$LL.facilitatorCode()}
      </label>
      <div class="control">
        <TextInput
          name="facilitatorCode"
          bind:value="{facilitatorCode}"
          id="facilitatorCode"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

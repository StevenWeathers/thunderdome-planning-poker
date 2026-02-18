<script lang="ts">
  import Modal from './global/Modal.svelte';
  import LL from '../i18n/i18n-svelte';
  import SolidButton from './global/SolidButton.svelte';
  import TextInput from './forms/TextInput.svelte';
  import { onMount } from 'svelte';
  import { Crown } from '@lucide/svelte';

  interface Props {
    toggleBecomeFacilitator?: any;
    handleBecomeFacilitator?: any;
  }

  let { toggleBecomeFacilitator = () => {}, handleBecomeFacilitator = () => {} }: Props = $props();

  let facilitatorCode = $state('');

  function handleSubmit(e) {
    e.preventDefault();

    handleBecomeFacilitator(facilitatorCode);
  }

  let focusInput: any;
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleBecomeFacilitator} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabelledby="modalTitle">
  <form onsubmit={handleSubmit} name="becomeFacilitator" class="space-y-4">
    <h2 class="text-xl font-bold dark:text-gray-300" id="modalTitle">{$LL.becomeFacilitator()}</h2>

    <div>
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="facilitatorCode">
        {$LL.facilitatorCode()}
      </label>
      <div class="control">
        <TextInput
          name="facilitatorCode"
          bind:value={facilitatorCode}
          placeholder={$LL.enterFacilitatorCode()}
          id="facilitatorCode"
          bind:this={focusInput}
          icon={Crown}
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

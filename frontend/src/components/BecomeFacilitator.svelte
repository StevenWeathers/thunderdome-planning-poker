<script lang="ts">
  import Modal from './Modal.svelte';
  import LL from '../i18n/i18n-svelte';
  import SolidButton from './SolidButton.svelte';

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
        <input
          name="facilitatorCode"
          bind:value="{facilitatorCode}"
          class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          id="facilitatorCode"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

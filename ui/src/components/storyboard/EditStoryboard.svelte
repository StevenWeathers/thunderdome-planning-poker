<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';

  export let toggleEditStoryboard = () => {};
  export let handleStoryboardEdit = () => {};
  export let storyboardName = '';
  export let joinCode = '';
  export let facilitatorCode = '';

  function saveStoryboard(e) {
    e.preventDefault();

    const storyboard = {
      storyboardName,
      joinCode,
      facilitatorCode,
    };

    handleStoryboardEdit(storyboard);
  }
</script>

<Modal
  closeModal="{toggleEditStoryboard}"
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
  <form on:submit="{saveStoryboard}" name="createStoryboard">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="storyboardName"
      >
        {$LL.storyboardName()}
      </label>
      <div class="control">
        <TextInput
          name="storyboardName"
          bind:value="{storyboardName}"
          placeholder="{$LL.storyboardNamePlaceholder()}"
          id="storyboardName"
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="joinCode"
      >
        {$LL.passCode()}
      </label>
      <div class="control">
        <TextInput
          name="joinCode"
          bind:value="{joinCode}"
          placeholder="{$LL.optionalPasscodePlaceholder()}"
          id="joinCode"
        />
      </div>
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="facilitatorCode"
      >
        {$LL.facilitatorCodeOptional()}
      </label>
      <div class="control">
        <TextInput
          name="facilitatorCode"
          bind:value="{facilitatorCode}"
          placeholder="{$LL.facilitatorCodePlaceholder()}"
          id="facilitatorCode"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

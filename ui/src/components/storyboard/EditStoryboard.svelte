<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { Crown, Lock } from 'lucide-svelte';

  interface Props {
    toggleEditStoryboard?: any;
    handleStoryboardEdit?: any;
    storyboardName?: string;
    joinCode?: string;
    facilitatorCode?: string;
  }

  let {
    toggleEditStoryboard = () => {},
    handleStoryboardEdit = () => {},
    storyboardName = $bindable(''),
    joinCode = $bindable(''),
    facilitatorCode = $bindable('')
  }: Props = $props();

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
  <form onsubmit={saveStoryboard} name="createStoryboard">
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
          icon="{Lock}"
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
          icon="{Crown}"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

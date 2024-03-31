<script lang="ts">
  import { appRoutes } from '../../config';
  import SolidButton from '../global/SolidButton.svelte';
  import ClipboardIcon from '../icons/ClipboardIcon.svelte';

  export let hostname = '';
  export let storyboardId = '';

  function copyStoryboardLink() {
    const storyboardLink = document.getElementById('StoryboardLink');
    storyboardLink.select();
    document.execCommand('copy');
  }
</script>

<h4 class="text-xl mb-2 leading-tight font-bold">Invite a user</h4>

<div class="flex flex-wrap items-stretch w-full">
  <input
    class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
        border-2 h-10 bg-gray-200 border-gray-200 rounded rounded-e-none px-3
        appearance-none text-gray-700 focus:outline-none focus:bg-white
        focus:border-orange-500"
    type="text"
    value="{hostname}{appRoutes.storyboard}/{storyboardId}"
    id="StoryboardLink"
    readonly
  />
  <div class="flex -mr-px">
    <SolidButton
      color="blue-copy"
      onClick="{copyStoryboardLink}"
      additionalClasses="flex items-center leading-normal
            whitespace-no-wrap text-sm"
    >
      <ClipboardIcon />
    </SolidButton>
  </div>
</div>

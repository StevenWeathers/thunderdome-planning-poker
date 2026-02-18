<script lang="ts">
  import { appRoutes } from '../../config';
  import SolidButton from '../global/SolidButton.svelte';
  import { ClipboardCopy } from '@lucide/svelte';

  interface Props {
    hostname?: string;
    retroId?: string;
  }

  let { hostname = '', retroId = '' }: Props = $props();

  function copyRetrospectiveLink() {
    const retrospectiveLink = document.getElementById('RetrospectiveLink');
    retrospectiveLink.select();
    document.execCommand('copy');
  }
</script>

<h4 class="text-xl mb-2 leading-tight font-bold text-center">Invite a user</h4>

<div class="flex flex-wrap items-stretch w-full">
  <input
    class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
        border-2 h-10 bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-900 rounded rounded-e-none px-3
        appearance-none text-gray-700 dark:text-gray-400 focus:outline-none focus:bg-white
        focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
    type="text"
    value="{hostname}{appRoutes.retro}/{retroId}"
    id="RetrospectiveLink"
    readonly
  />
  <div class="flex -me-px">
    <SolidButton
      color="blue-copy"
      onClick={copyRetrospectiveLink}
      additionalClasses="flex items-center leading-normal whitespace-no-wrap text-sm"
    >
      <ClipboardCopy />
    </SolidButton>
  </div>
</div>

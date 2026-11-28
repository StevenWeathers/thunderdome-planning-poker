<script lang="ts">
  import SolidButton from './SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { newVersion } from '../../stores';

  interface Props {
    reload?: () => void;
  }

  let { reload = () => window.location.reload() }: Props = $props();
</script>

<!-- the live region stays in the page so that screen readers announce the message when it appears -->
<div role="status">
  {#if $newVersion}
    <div class="bg-blue-100 dark:bg-sky-900">
      <div class="max-w-7xl mx-auto py-3 px-3 sm:px-6 lg:px-8 flex items-center justify-between flex-wrap gap-3">
        <p class="font-medium text-blue-800 dark:text-sky-100">{$LL.newVersionAvailable()}</p>
        <SolidButton color="blue" onClick={reload} testid="new-version-reload">{$LL.pageLoadReload()}</SolidButton>
      </div>
    </div>
  {/if}
</div>

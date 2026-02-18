<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { ApiClient } from '../../types/apiclient';
  import { onMount } from 'svelte';

  interface Props {
    xfetch: ApiClient;
    verifyId: any;
  }

  let { xfetch, verifyId }: Props = $props();

  let accountVerified = $state(false);
  let verficationError = $state(false);

  onMount(() => {
    xfetch('/api/auth/verify', { body: { verifyId }, method: 'PATCH' })
      .then(function () {
        accountVerified = true;
      })
      .catch(function () {
        verficationError = true;
      });
  });
</script>

<svelte:head>
  <title>{$LL.verifyAccount()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex justify-center">
    <div class="w-full md:w-1/2 xl:w-1/3 py-4">
      {#if accountVerified}
        <div
          class="bg-green-100 border border-green-400 text-green-700
                    px-4 py-3 rounded relative"
          role="alert"
        >
          <strong class="font-bold">
            {$LL.verifyAccountVerifiedTitle()}
          </strong>
          <p>{$LL.verifyAccountVerifiedThanks()}</p>
        </div>
      {:else if verficationError}
        <div
          class="bg-red-100 border border-red-400 text-red-700 px-4
                    py-3 rounded relative"
          role="alert"
        >
          <strong class="font-bold">
            {$LL.verifyAccountFailedTitle()}
          </strong>
          <p>{$LL.verifyAccountFailedError()}</p>
        </div>
      {:else}
        <div class="text-center">
          <h1 class="text-4xl text-teal-500 leading-tight font-bold">
            {$LL.verifyAccountLoading()}
          </h1>
        </div>
      {/if}
    </div>
  </div>
</PageLayout>

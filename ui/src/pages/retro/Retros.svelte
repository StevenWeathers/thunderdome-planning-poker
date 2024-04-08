<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateRetro from '../../components/retro/CreateRetro.svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import BoxList from '../../components/BoxList.svelte';
  import Pagination from '../../components/global/Pagination.svelte';

  export let xfetch;
  export let notifications;
  export let router;
  export let eventTag;

  let retros = [];
  const retrosPageLimit = 10;
  let retroCount = 0;
  let retrosPage = 1;

  function getRetros() {
    const retrosOffset = (retrosPage - 1) * retrosPageLimit;

    xfetch(
      `/api/users/${$user.id}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        retros = result.data;
        retroCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getRetrosErrorMessage());
        eventTag('fetch_retros', 'engagement', 'failure');
      });
  }

  const changePage = evt => {
    retrosPage = evt.detail;
    getRetros();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
    }
    getRetros();
  });
</script>

<svelte:head>
  <title>{$LL.yourRetros()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
  >
    {$LL.myRetros()}
  </h1>

  <div class="flex flex-wrap">
    <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
      <BoxList
        items="{retros}"
        itemType="retro"
        pageRoute="{appRoutes.retro}"
        ownerField="ownerId"
        showOwnerName="{true}"
        ownerNameField="teamName"
        joinBtnText="{$LL.joinRetro()}"
      />
      {#if retroCount > retrosPageLimit}
        <div class="mt-6 pt-1 flex justify-center">
          <Pagination
            bind:current="{retrosPage}"
            num_items="{retroCount}"
            per_page="{retrosPageLimit}"
            on:navigate="{changePage}"
          />
        </div>
      {/if}
    </div>

    <div class="w-full md:w-1/2 lg:w-2/5 md:ps-2 xl:ps-4">
      <div
        class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
      >
        <h2
          class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
        >
          {$LL.createARetro()}
        </h2>
        <CreateRetro
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
        />
      </div>
    </div>
  </div>
</PageLayout>

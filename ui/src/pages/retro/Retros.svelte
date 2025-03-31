<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateRetro from '../../components/retro/CreateRetro.svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import BoxList from '../../components/BoxList.svelte';
  import Pagination from '../../components/global/Pagination.svelte';

  interface Props {
    xfetch: any;
    notifications: any;
    router: any;
  }

  let { xfetch, notifications, router }: Props = $props();

  let retros = $state([]);
  const retrosPageLimit = 10;
  let retroCount = $state(0);
  let retrosPage = $state(1);
  let loading = $state(true);

  function getRetros() {
    const retrosOffset = (retrosPage - 1) * retrosPageLimit;

    xfetch(
      `/api/users/${$user.id}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        retros = result.data;
        retroCount = result.meta.count;
        loading = false;
      })
      .catch(function () {
        notifications.danger($LL.getRetrosErrorMessage());
        loading = false;
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
      {#if retroCount > 0}
        <BoxList
          items="{retros}"
          itemType="retro"
          pageRoute="{appRoutes.retro}"
          ownerField="ownerId"
          showOwnerName="{true}"
          ownerNameField="teamName"
          joinBtnText="{$LL.joinRetro()}"
        />
      {:else if loading === false}
        <div
          class="w-full my-10 text-lg md:text-xl dark:text-white text-center"
        >
          {$LL.noRetrosFound()}
        </div>
      {/if}
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
          notifications={notifications}
          router="{router}"
          xfetch={xfetch}
        />
      </div>
    </div>
  </div>

  <div class="w-full text-gray-600 dark:text-gray-300">
    <p class="py-8 md:text-lg italic">{$LL.retroDescription()}</p>
  </div>
</PageLayout>

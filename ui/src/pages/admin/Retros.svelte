<script lang="ts">
  import { onMount } from 'svelte';
  import Pagination from '../../components/global/Pagination.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  // export let eventTag

  const retrosPageLimit = 100;
  let retroCount = 0;
  let retros = [];
  let retrosPage = 1;
  let activeRetros = false;

  function getRetros() {
    const retrosOffset = (retrosPage - 1) * retrosPageLimit;
    xfetch(
      `/api/retros?limit=${retrosPageLimit}&offset=${retrosOffset}&active=${activeRetros}`,
    )
      .then(res => res.json())
      .then(function (result) {
        retros = result.data;
        retroCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger(`${$LL.getRetrosErrorMessage()}`);
      });
  }

  const changePage = evt => {
    retrosPage = evt.detail;
    getRetros();
  };

  const changeActiveRetrosToggle = () => {
    retrosPage = 1;
    getRetros();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getRetros();
  });
</script>

<svelte:head>
  <title>{$LL.retros()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="retros">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.retros()}
    </h1>
  </div>

  <div class="w-full">
    <div class="text-right mb-4">
      <div
        class="relative inline-block w-10 me-2 align-middle select-none transition duration-200 ease-in"
      >
        <input
          type="checkbox"
          name="activeRetros"
          id="activeRetros"
          bind:checked="{activeRetros}"
          on:change="{changeActiveRetrosToggle}"
          class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
        />
        <label
          for="activeRetros"
          class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
        >
        </label>
      </div>
      <label for="activeRetros" class="dark:text-gray-300"
        >{$LL.showActiveRetros()}</label
      >
    </div>

    <Table>
      <tr slot="header">
        <HeadCol>
          {$LL.name()}
        </HeadCol>
        <HeadCol>
          {$LL.dateCreated()}
        </HeadCol>
        <HeadCol>
          {$LL.dateUpdated()}
        </HeadCol>
        <HeadCol type="action">
          <span class="sr-only">{$LL.actions()}</span>
        </HeadCol>
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each retros as retro, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <a
                href="{appRoutes.admin}/retros/{retro.id}"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >{retro.name}</a
              >
            </RowCol>
            <RowCol>
              {new Date(retro.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(retro.updatedDate).toLocaleString()}
            </RowCol>
            <RowCol type="action">
              <HollowButton href="{appRoutes.retro}/{retro.id}">
                {$LL.joinRetro()}
              </HollowButton>
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if retroCount > retrosPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{retrosPage}"
          num_items="{retroCount}"
          per_page="{retrosPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}
  </div>
</AdminPageLayout>

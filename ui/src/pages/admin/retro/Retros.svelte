<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import Table from '../../../components/table/Table.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import Toggle from '../../../components/forms/Toggle.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const retrosPageLimit = 100;
  let retroCount = $state(0);
  let retros = $state([]);
  let retrosPage = $state(1);
  let activeRetros = $state(false);

  function getRetros() {
    const retrosOffset = (retrosPage - 1) * retrosPageLimit;
    xfetch(`/api/retros?limit=${retrosPageLimit}&offset=${retrosOffset}&active=${activeRetros}`)
      .then(res => res.json())
      .then(function (result) {
        retros = result.data;
        retroCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getRetrosErrorMessage());
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
  <TableContainer>
    <TableNav title={$LL.retros()} createBtnEnabled={false}>
      <Toggle
        name="activeRetros"
        id="activeRetros"
        bind:checked={activeRetros}
        changeHandler={changeActiveRetrosToggle}
        label={$LL.showActiveRetros()}
      />
    </TableNav>
    <Table>
      {#snippet header()}
        <tr>
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
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each retros as retro, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <a
                  href="{appRoutes.admin}/retros/{retro.id}"
                  class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">{retro.name}</a
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
      {/snippet}
    </Table>
    <TableFooter bind:current={retrosPage} num_items={retroCount} per_page={retrosPageLimit} on:navigate={changePage} />
  </TableContainer>
</AdminPageLayout>

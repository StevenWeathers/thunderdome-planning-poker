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
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import Toggle from '../../../components/forms/Toggle.svelte';

  import type { NotificationService } from '../../../types/notifications';

  interface Props {
    xfetch: any;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const storyboardsPageLimit = 100;
  let storyboardCount = $state(0);
  let storyboards = $state([]);
  let storyboardsPage = $state(1);
  let activeStoryboards = $state(false);

  function getStoryboards() {
    const storyboardsOffset = (storyboardsPage - 1) * storyboardsPageLimit;
    xfetch(
      `/api/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}&active=${activeStoryboards}`,
    )
      .then(res => res.json())
      .then(function (result) {
        storyboards = result.data;
        storyboardCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardsErrorMessage());
      });
  }

  const changePage = evt => {
    storyboardsPage = evt.detail;
    getStoryboards();
  };

  const changeActiveStoryboardsToggle = () => {
    storyboardsPage = 1;
    getStoryboards();
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

    getStoryboards();
  });
</script>

<svelte:head>
  <title>{$LL.storyboards()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="storyboards">
  <TableContainer>
    <TableNav title={$LL.storyboards()} createBtnEnabled={false}>
      <Toggle
        name="activeStoryboards"
        id="activeStoryboards"
        bind:checked="{activeStoryboards}"
        changeHandler={changeActiveStoryboardsToggle}
        label={$LL.showActiveStoryboards()}
      />
    </TableNav>
    <Table>
      {#snippet header()}
            <tr >
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
            <tbody   class="{className}">
          {#each storyboards as storyboard, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <a
                  href="{appRoutes.admin}/storyboards/{storyboard.id}"
                  class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                  >{storyboard.name}</a
                >
              </RowCol>
              <RowCol>
                {new Date(storyboard.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(storyboard.updatedDate).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                <HollowButton href="{appRoutes.storyboard}/{storyboard.id}">
                  {$LL.joinStoryboard()}
                </HollowButton>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
    <TableFooter
      bind:current={storyboardsPage}
      num_items={storyboardCount}
      per_page={storyboardsPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>
</AdminPageLayout>

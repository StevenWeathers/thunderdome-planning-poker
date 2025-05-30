<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import Table from '../../../components/table/Table.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import CrudActions from '../../../components/table/CrudActions.svelte';
  import BooleanDisplay from '../../../components/global/BooleanDisplay.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
  }

  let { xfetch, router, notifications }: Props = $props();

  const apikeysPageLimit = 100;

  let appStats = $state({
    unregisteredUserCount: 0,
    registeredUserCount: 0,
    battleCount: 0,
    planCount: 0,
    organizationCount: 0,
    departmentCount: 0,
    teamCount: 0,
    apikeyCount: 0,
  });
  let apikeys = $state([]);
  let apikeysPage = $state(1);

  function getAppStats() {
    xfetch('/api/admin/stats')
      .then(res => res.json())
      .then(function (result) {
        appStats = result.data;
      })
      .catch(function () {
        notifications.danger($LL.applicationStatsError());
      });
  }

  function getApiKeys() {
    const apikeysOffset = (apikeysPage - 1) * apikeysPageLimit;
    xfetch(
      `/api/admin/apikeys?limit=${apikeysPageLimit}&offset=${apikeysOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        apikeys = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getApikeysError());
      });
  }

  function deleteApiKey(userId, apk) {
    return function () {
      xfetch(`/api/users/${userId}/apikeys/${apk}`, {
        method: 'DELETE',
      })
        .then(res => res.json())
        .then(function () {
          notifications.success($LL.apiKeyDeleteSuccess());
          getAppStats();
          getApiKeys();
        })
        .catch(function () {
          notifications.danger($LL.apiKeyDeleteFailed());
        });
    };
  }

  function toggleApiKeyActiveStatus(userId, apk, active) {
    return function () {
      const body = {
        active: !active,
      };

      xfetch(`/api/users/${userId}/apikeys/${apk}`, {
        body,
        method: 'PUT',
      })
        .then(res => res.json())
        .then(function () {
          notifications.success($LL.apiKeyUpdateSuccess());
          getApiKeys();
        })
        .catch(function () {
          notifications.danger($LL.apiKeyUpdateFailed());
        });
    };
  }

  const changePage = evt => {
    apikeysPage = evt.detail;
    getApiKeys();
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

    getAppStats();
    getApiKeys();
  });
</script>

<svelte:head>
  <title>{$LL.apiKeys()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="apikeys">
  <TableContainer>
    <TableNav title={$LL.apiKeys()} createBtnEnabled={false} />
    <Table>
      {#snippet header()}
            <tr >
          <HeadCol>
            {$LL.apiKeyName()}
          </HeadCol>
          <HeadCol>
            {$LL.apiKeyPrefix()}
          </HeadCol>
          <HeadCol>
            {$LL.userName()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.dateCreated()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
          <HeadCol>
            {$LL.actions()}
          </HeadCol>
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
          {#each apikeys as apikey, i}
            <TableRow itemIndex={i}>
              <RowCol>
                {apikey.name}
              </RowCol>
              <RowCol>
                {apikey.prefix}
              </RowCol>
              <RowCol>
                <a
                  href="{appRoutes.adminUsers}/{apikey.userId}"
                  class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >
                  {apikey.userName}
                </a>
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue={apikey.active} />
              </RowCol>
              <RowCol>
                {new Date(apikey.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(apikey.updatedDate).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                <CrudActions
                  editBtnEnabled={false}
                  deleteBtnClickHandler={deleteApiKey(apikey.userId, apikey.id)}
                  deleteBtnTestId="apikey-delete"
                >
                  <HollowButton
                    onClick={toggleApiKeyActiveStatus(
                      apikey.userId,
                      apikey.id,
                      apikey.active,
                    )}
                    testid="apikey-activetoggle"
                  >
                    {#if !apikey.active}
                      {$LL.activate()}
                    {:else}
                      {$LL.deactivate()}
                    {/if}
                  </HollowButton>
                </CrudActions>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
    <TableFooter
      bind:current={apikeysPage}
      num_items={appStats.apikeyCount}
      per_page={apikeysPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>
</AdminPageLayout>

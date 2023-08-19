<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import Pagination from '../../components/Pagination.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  // export let eventTag

  const apikeysPageLimit = 100;

  let appStats = {
    unregisteredUserCount: 0,
    registeredUserCount: 0,
    battleCount: 0,
    planCount: 0,
    organizationCount: 0,
    departmentCount: 0,
    teamCount: 0,
    apikeyCount: 0,
  };
  let apikeys = [];
  let apikeysPage = 1;

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
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.apiKeys()}
    </h1>
  </div>

  <div class="w-full">
    <Table>
      <tr slot="header">
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
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each apikeys as apikey, i}
          <TableRow itemIndex="{i}">
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
              {#if apikey.active}
                <span class="text-green-600"><CheckIcon /></span>
              {/if}
            </RowCol>
            <RowCol>
              {new Date(apikey.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(apikey.updatedDate).toLocaleString()}
            </RowCol>
            <RowCol>
              <HollowButton
                onClick="{toggleApiKeyActiveStatus(
                  apikey.userId,
                  apikey.id,
                  apikey.active,
                )}"
                testid="apikey-activetoggle"
              >
                {#if !apikey.active}
                  {$LL.activate()}
                {:else}
                  {$LL.deactivate()}
                {/if}
              </HollowButton>
              <HollowButton
                color="red"
                onClick="{deleteApiKey(apikey.userId, apikey.id)}"
                testid="apikey-delete"
              >
                {$LL.delete()}
              </HollowButton>
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if appStats.apikeyCount > apikeysPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{apikeysPage}"
          num_items="{appStats.apikeyCount}"
          per_page="{apikeysPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}
  </div>
</AdminPageLayout>

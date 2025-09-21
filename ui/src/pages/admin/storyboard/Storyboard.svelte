<script lang="ts">
  import { onMount } from 'svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import Table from '../../../components/table/Table.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import BooleanDisplay from '../../../components/global/BooleanDisplay.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    storyboardId: any;
  }

  let { xfetch, router, notifications, storyboardId }: Props = $props();

  let storyboard = $state({
    name: '',
    users: [],
    owner_id: '',
    createdDate: '',
    updatedDate: '',
  });

  let showDeleteStoryboard = $state(false);

  function getStoryboard() {
    xfetch(`/api/storyboards/${storyboardId}`)
      .then(res => res.json())
      .then(function (result) {
        storyboard = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardErrorMessage());
      });
  }

  function deleteStoryboard() {
    xfetch(`/api/storyboards/${storyboardId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        router.route(appRoutes.adminStoryboards);
      })
      .catch(function () {
        notifications.danger($LL.deleteStoryboardErrorMessage());
      });
  }

  function toggleDeleteStoryboard() {
    showDeleteStoryboard = !showDeleteStoryboard;
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getStoryboard();
  });
</script>

<svelte:head>
  <title>{$LL.storyboard()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="storyboards">
  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title={storyboard.name} createBtnEnabled={false} />
      <Table>
        {#snippet header()}
          <tr>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
          </tr>
        {/snippet}
        {#snippet body({ class: className })}
          <tbody class={className}>
            <TableRow itemIndex={0}>
              <RowCol>
                {new Date(storyboard.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(storyboard.updatedDate).toLocaleString()}
              </RowCol>
            </TableRow>
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>
  </div>

  <TableContainer>
    <TableNav title={$LL.users()} createBtnEnabled={false} />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.abandoned()}
          </HeadCol>
          <HeadCol>
            {$LL.facilitator()}
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each storyboard.users as user, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId={user.id}
                      avatar={user.avatar}
                      gravatarHash={user.gravatarHash}
                      userName={user.name}
                      width={48}
                      class="h-10 w-10 rounded-full"
                    />
                  </div>
                  <div class="ms-4">
                    <div class="text-sm font-medium text-gray-900">
                      <a
                        href="{appRoutes.adminUsers}/{user.id}"
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        >{user.name}</a
                      >
                      {#if user.country}
                        &nbsp;
                        <CountryFlag country={user.country} additionalClass="inline-block" width="32" height="24" />
                      {/if}
                    </div>
                  </div>
                </div>
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue={user.active} />
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue={user.abandoned} />
              </RowCol>
              <RowCol>
                <RowCol>
                  <BooleanDisplay boolValue={storyboard.owner_id === user.id} />
                </RowCol>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  <div class="text-center mt-4">
    <HollowButton color="red" onClick={toggleDeleteStoryboard} testid="storyboard-delete">
      {$LL.deleteStoryboard()}
    </HollowButton>
  </div>

  {#if showDeleteStoryboard}
    <DeleteConfirmation
      toggleDelete={toggleDeleteStoryboard}
      handleDelete={deleteStoryboard}
      confirmText={'Are you sure you want to delete this Storyboard?'}
      confirmBtnText={'Delete Storyboard'}
    />
  {/if}
</AdminPageLayout>

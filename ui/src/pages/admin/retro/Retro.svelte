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
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import BooleanDisplay from '../../../components/global/BooleanDisplay.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
    retroId: any;
  }

  let {
    xfetch,
    router,
    notifications,
    retroId
  }: Props = $props();

  let showDeleteRetro = $state(false);

  let retro = $state({
    name: '',
    users: [],
    owner_id: '',
    createdDate: '',
    updatedDate: '',
  });

  function getRetro() {
    xfetch(`/api/retros/${retroId}`)
      .then(res => res.json())
      .then(function (result) {
        retro = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getRetroErrorMessage());
      });
  }

  function deleteRetro() {
    xfetch(`/api/retros/${retroId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        router.route(appRoutes.adminRetros);
      })
      .catch(function () {
        notifications.danger($LL.deleteRetroErrorMessage());
      });
  }

  function toggleDeleteRetro() {
    showDeleteRetro = !showDeleteRetro;
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

    getRetro();
  });
</script>

<svelte:head>
  <title>{$LL.retro()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="retros">
  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title={retro.name} createBtnEnabled={false} />
      <Table>
        {#snippet header()}
                <tr >
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
          </tr>
              {/snippet}
        {#snippet body({ class: className })}
                <tbody   class="{className}">
            <TableRow itemIndex={0}>
              <RowCol>
                {new Date(retro.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(retro.updatedDate).toLocaleString()}
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
            <tr >
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
            {$LL.leader()}
          </HeadCol>
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
          {#each retro.users as user, i}
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
                        <CountryFlag
                          country={user.country}
                          additionalClass="inline-block"
                          width="32"
                          height="24"
                        />
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
                <BooleanDisplay boolValue={retro.owner_id === user.id} />
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
  </TableContainer>

  <div class="text-center mt-4">
    <HollowButton
      color="red"
      onClick={toggleDeleteRetro}
      testid="retro-delete"
    >
      {$LL.deleteRetro()}
    </HollowButton>
  </div>

  {#if showDeleteRetro}
    <DeleteConfirmation
      toggleDelete={toggleDeleteRetro}
      handleDelete={deleteRetro}
      confirmText={$LL.confirmDeleteRetro()}
      confirmBtnText={$LL.deleteRetro()}
    />
  {/if}
</AdminPageLayout>

<script lang="ts">
  import { onMount } from 'svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import Table from '../../../components/table/Table.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import ProfileForm from '../../../components/user/ProfileForm.svelte';
  import Modal from '../../../components/global/Modal.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import CreateUser from '../../../components/user/CreateUser.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import CrudActions from '../../../components/table/CrudActions.svelte';
  import { BadgeCheck, ToggleLeft, ToggleRight } from 'lucide-svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const usersPageLimit = 100;

  let totalUsers = $state(0);
  let users = $state([]);
  let showCreateUser = $state(false);
  let usersPage = $state(1);
  let userDeleteId = null;
  let showUserDeletion = $state(false);
  let searchEmail = '';
  let showUserEdit = $state(false);
  let selectedUserProfile = $state({});

  const toggleDeleteUser = id => () => {
    showUserDeletion = !showUserDeletion;
    userDeleteId = id;
  };

  function toggleCreateUser() {
    showCreateUser = !showCreateUser;
  }

  const toggleUserEdit = profile => () => {
    showUserEdit = !showUserEdit;
    selectedUserProfile = profile;
  };

  function createUser(warriorName, warriorEmail, warriorPassword1, warriorPassword2) {
    const body = {
      name: warriorName,
      email: warriorEmail,
      password1: warriorPassword1,
      password2: warriorPassword2,
    };

    xfetch('/api/admin/users', { body })
      .then(function () {
        getUsers();
        toggleCreateUser();
      })
      .catch(function () {
        notifications.danger($LL.createUserError());
      });
  }

  function getUsers() {
    const offset = (usersPage - 1) * usersPageLimit;
    const isSearch = searchEmail !== '';
    const apiPrefix = isSearch ? `/api/admin/search/users/email?search=${searchEmail}&` : '/api/admin/users?';

    if (isSearch && searchEmail.length < 3) {
      notifications.danger($LL.searchLengthError());
      return;
    }

    xfetch(`${apiPrefix}limit=${usersPageLimit}&offset=${offset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
        totalUsers = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getUsersError());
      });
  }

  function handleUserEdit(p) {
    xfetch(`/api/users/${selectedUserProfile.id}`, {
      body: p,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.profileUpdateSuccess());
        getUsers();
        toggleUserEdit({})();
      })
      .catch(function () {
        notifications.danger($LL.profileErrorUpdating());
      });
  }

  function promoteUser(userId) {
    return function () {
      xfetch(`/api/admin/users/${userId}/promote`, { method: 'PATCH' })
        .then(function () {
          getUsers();
        })
        .catch(function () {
          notifications.danger($LL.promoteUserError());
        });
    };
  }

  function demoteUser(userId) {
    return function () {
      xfetch(`/api/admin/users/${userId}/demote`, { method: 'PATCH' })
        .then(function () {
          getUsers();
        })
        .catch(function () {
          notifications.danger($LL.demoteUserError());
        });
    };
  }

  function disableUser(userId) {
    return function () {
      xfetch(`/api/admin/users/${userId}/disable`, { method: 'PATCH' })
        .then(function () {
          getUsers();
        })
        .catch(function () {
          notifications.danger('Error disabling user');
        });
    };
  }

  function enableUser(userId) {
    return function () {
      xfetch(`/api/admin/users/${userId}/enable`, { method: 'PATCH' })
        .then(function () {
          getUsers();
        })
        .catch(function () {
          notifications.danger('Error enabling user');
        });
    };
  }

  function handleDeleteUser() {
    xfetch(`/api/users/${userDeleteId}`, { method: 'DELETE' })
      .then(function () {
        getUsers();
        toggleDeleteUser(null)();
      })
      .catch(function () {
        notifications.danger('deleteUserError');
      });
  }

  function onSearchSubmit(term) {
    searchEmail = term;
    usersPage = 1;
    getUsers();
  }

  const changePage = evt => {
    usersPage = evt.detail;
    getUsers();
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

    getUsers();
  });
</script>

<svelte:head>
  <title>{$LL.users()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="users">
  <TableContainer>
    <TableNav
      title={$LL.registeredUsers()}
      createBtnText={$LL.warriorCreate()}
      createButtonHandler={toggleCreateUser}
      createBtnTestId="user-create"
      searchEnabled={true}
      searchPlaceholder={$LL.email()}
      searchHandler={onSearchSubmit}
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.email()}
          </HeadCol>
          <HeadCol>
            {$LL.company()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>Enabled</HeadCol>
          <HeadCol type="action">
            <span class="sr-only">{$LL.actions()}</span>
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each users as user, i}
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
                {user.email}
                {#if user.verified}
                  <span class="text-green-600" title={$LL.verified()}>
                    <BadgeCheck class="inline-block" />
                  </span>
                {/if}
              </RowCol>
              <RowCol>
                <div class="text-sm text-gray-900 dark:text-gray-400">
                  {user.company}
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-300">
                  {user.jobTitle}
                </div>
              </RowCol>
              <RowCol>
                <span class="text-gray-500 dark:text-gray-300">{user.rank}</span>
              </RowCol>
              <RowCol>
                <button
                  onclick={!user.disabled ? disableUser(user.id) : enableUser(user.id)}
                  class:text-red-500={user.disabled}
                  class:text-green-500={!user.disabled}
                  title={!user.disabled ? 'enabled' : 'disabled'}
                >
                  {#if user.disabled}
                    <ToggleLeft class="h-10 w-10" />
                  {:else}
                    <ToggleRight class="h-10 w-10" />
                  {/if}
                </button>
              </RowCol>
              <RowCol type="action">
                <CrudActions
                  editBtnClickHandler={toggleUserEdit(user)}
                  deleteBtnClickHandler={toggleDeleteUser(user.id)}
                >
                  {#if user.rank !== 'ADMIN'}
                    <HollowButton onClick={promoteUser(user.id)} color="blue">
                      {$LL.promote()}
                    </HollowButton>
                  {:else}
                    <HollowButton onClick={demoteUser(user.id)} color="blue">
                      {$LL.demote()}
                    </HollowButton>
                  {/if}
                </CrudActions>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
    <TableFooter bind:current={usersPage} num_items={totalUsers} per_page={usersPageLimit} on:navigate={changePage} />
  </TableContainer>

  {#if showCreateUser}
    <CreateUser toggleCreate={toggleCreateUser} handleCreate={createUser} notifications />
  {/if}

  {#if showUserEdit}
    <Modal closeModal={toggleUserEdit({})}>
      <ProfileForm profile={selectedUserProfile} handleUpdate={handleUserEdit} {xfetch} {notifications} />
    </Modal>
  {/if}

  {#if showUserDeletion}
    <DeleteConfirmation
      toggleDelete={toggleDeleteUser(null)}
      handleDelete={handleDeleteUser}
      confirmText={$LL.deleteAccountWarningStatement()}
      confirmBtnText={$LL.deleteConfirmButton()}
    />
  {/if}
</AdminPageLayout>

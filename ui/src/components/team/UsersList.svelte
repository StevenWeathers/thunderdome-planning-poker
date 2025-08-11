<script lang="ts">
  import Table from '../table/Table.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';

  import LL from '../../i18n/i18n-svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import AddUser from '../../components/team/AddUser.svelte';
  import UpdateUser from '../../components/team/UpdateUser.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import { createEventDispatcher } from 'svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    orgId?: string;
    deptId?: string;
    teamPrefix?: string;
    isAdmin?: boolean;
    pageType?: string;
    users?: any;
    getUsers?: any;
  }

  let {
    xfetch,
    notifications,
    orgId = '',
    deptId = '',
    teamPrefix = '',
    isAdmin = false,
    pageType = '',
    users = [],
    getUsers = () => {}
  }: Props = $props();

  const dispatch = createEventDispatcher();

  let showAddUser = $state(false);
  let showUpdateUser = $state(false);
  let updateUser = $state({});
  let showRemoveUser = $state(false);
  let removeUserId = null;

  function handleUserAdd({ id, email, role }) {
    const body = {
      user_id: id,
      role,
    };

    xfetch(`${teamPrefix}/users`, { body })
      .then(result => result.json())
      .then(function () {
        toggleAddUser();
        dispatch('user-added');
        notifications.success($LL.userAddSuccess());

        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userAddError());
      });
  }

  function handleUserInvite({ id, email, role }) {
    const body = {
      email,
      role,
    };

    xfetch(`${teamPrefix}/invites`, { body })
      .then(result => result.json())
      .then(function (result) {
        toggleAddUser();
        dispatch('user-invited');
        notifications.success($LL.userInviteSent());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userAddError());
      });
  }

  function handleUserUpdate(userId, role) {
    const body = {
      role,
    };

    xfetch(`${teamPrefix}/users/${userId}`, { body, method: 'PUT' })
      .then(function () {
        toggleUpdateUser({})();
        notifications.success($LL.userUpdateSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userUpdateError());
      });
  }

  function handleUserRemove() {
    xfetch(`${teamPrefix}/users/${removeUserId}`, { method: 'DELETE' })
      .then(function () {
        toggleRemoveUser(null)();
        notifications.success($LL.userRemoveSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userRemoveError());
      });
  }

  function toggleAddUser() {
    showAddUser = !showAddUser;
  }

  const toggleUpdateUser = user => () => {
    updateUser = user;
    showUpdateUser = !showUpdateUser;
  };

  const toggleRemoveUser = userId => () => {
    showRemoveUser = !showRemoveUser;
    removeUserId = userId;
  };
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title={$LL.users()}
      createBtnEnabled={isAdmin}
      createBtnText={$LL.userAdd()}
      createButtonHandler={toggleAddUser}
      createBtnTestId="user-add"
    />
    <Table>
      {#snippet header()}
            <tr >
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.email()}
          </HeadCol>
          <HeadCol>
            {$LL.role()}
          </HeadCol>
          <HeadCol type="action">
            <span class="sr-only">{$LL.actions()}</span>
          </HeadCol>
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
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
                    <div class="font-medium text-gray-900 dark:text-gray-200">
                      <span data-testid="user-name">{user.name}</span>
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
                <span data-testid="user-email">{user.email}</span>
              </RowCol>
              <RowCol>
                <div class="text-sm text-gray-500 dark:text-gray-300">
                  {user.role}
                </div>
              </RowCol>
              <RowCol type="action">
                {#if isAdmin}
                  <CrudActions
                    editBtnClickHandler={toggleUpdateUser(user)}
                    deleteBtnClickHandler={toggleRemoveUser(user.id)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
  </TableContainer>

  {#if showAddUser}
    <AddUser
      toggleAdd={toggleAddUser}
      handleAdd={handleUserAdd}
      handleInvite={handleUserInvite}
      pageType={pageType}
      orgId={orgId}
      deptId={deptId}
      xfetch={xfetch}
      notifications={notifications}
    />
  {/if}

  {#if showUpdateUser}
    <UpdateUser
      toggleUpdate={toggleUpdateUser({})}
      handleUpdate={handleUserUpdate}
      userId={updateUser.id}
      userEmail={updateUser.email}
      role={updateUser.role}
    />
  {/if}

  {#if showRemoveUser}
    <DeleteConfirmation
      toggleDelete={toggleRemoveUser(null)}
      handleDelete={handleUserRemove}
      permanent={false}
      confirmText={$LL.removeUserConfirmText()}
      confirmBtnText={$LL.removeUser()}
    />
  {/if}
</div>

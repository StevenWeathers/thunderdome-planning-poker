<script lang="ts">
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import Table from '../table/Table.svelte';
  import LL from '../../i18n/i18n-svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import { onMount } from 'svelte';
  import CrudActions from '../table/CrudActions.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
    teamPrefix?: String;
    pageType?: String;
  }

  let {
    xfetch,
    router,
    notifications,
    teamPrefix = '',
    pageType = ''
  }: Props = $props();
  export const f = event => {
    if (event === 'user-invited') {
      getInvites();
    }
  };

  let invites = $state([]);
  let showDeleteInvite = $state(false);
  let deleteInviteId = '';

  const toggleDeleteInvite = inviteId => () => {
    showDeleteInvite = !showDeleteInvite;
    deleteInviteId = inviteId;
  };

  function getInvites() {
    xfetch(`${teamPrefix}/invites`)
      .then(res => res.json())
      .then(function (result) {
        invites = result.data;
      })
      .catch(function () {
        notifications.danger(`error getting ${pageType} invites`);
      });
  }

  function handleDeleteInvite() {
    xfetch(`${teamPrefix}/invites/${deleteInviteId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteInvite(null)();
        notifications.success('Successfully deleted user invite');
        getInvites();
      })
      .catch(function () {
        notifications.danger('Error deleting user invite');
      });
  }

  onMount(() => {
    getInvites();
  });
</script>

<TableContainer>
  <TableNav title={$LL.userInvites()} createBtnEnabled={false} />
  <Table>
    {#snippet header()}
        <tr >
        <HeadCol>{$LL.email()}</HeadCol>
        <HeadCol>{$LL.role()}</HeadCol>
        <HeadCol>{$LL.dateCreated()}</HeadCol>
        <HeadCol>{$LL.expireDate()}</HeadCol>
        <HeadCol />
      </tr>
      {/snippet}
    {#snippet body({ class: className })}
        <tbody   class="{className}">
        {#each invites as item, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <span data-testid="invite-user-email">{item.email}</span>
            </RowCol>
            <RowCol>
              {item.role}
            </RowCol>
            <RowCol>
              {new Date(item.created_date).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(item.expire_date).toLocaleString()}
            </RowCol>
            <RowCol type="action">
              <CrudActions
                editBtnEnabled={false}
                deleteBtnClickHandler={toggleDeleteInvite(item.invite_id)}
              />
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
      {/snippet}
  </Table>
</TableContainer>

{#if showDeleteInvite}
  <DeleteConfirmation
    toggleDelete={toggleDeleteInvite(null)}
    handleDelete={handleDeleteInvite}
    confirmText={$LL.userInviteConfirmDelete()}
    confirmBtnText={$LL.userInviteDelete()}
  />
{/if}

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

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;
  export let teamPrefix: String = '';
  export let pageType: String = '';
  export const f = event => {
    if (event === 'user-invited') {
      getInvites();
    }
  };

  let invites = [];
  let showDeleteInvite = false;
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
        eventTag(`${pageType}_get_invites`, 'engagement', 'success');
      })
      .catch(function () {
        eventTag(`${pageType}_get_invites`, 'engagement', 'failure');
        notifications.danger(`error getting ${pageType} invites`);
      });
  }

  function handleDeleteInvite() {
    xfetch(`${teamPrefix}/invites/${deleteInviteId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag(`${pageType}_delete_invite`, 'engagement', 'success');
        toggleDeleteInvite(null)();
        notifications.success('Successfully deleted user invite');
        getInvites();
      })
      .catch(function () {
        notifications.danger('Error deleting user invite');
        eventTag(`${pageType}_delete_invite`, 'engagement', 'failure');
      });
  }

  onMount(() => {
    getInvites();
  });
</script>

<TableContainer>
  <TableNav title="{$LL.userInvites()}" createBtnEnabled="{false}" />
  <Table>
    <tr slot="header">
      <HeadCol>{$LL.email()}</HeadCol>
      <HeadCol>{$LL.role()}</HeadCol>
      <HeadCol>{$LL.dateCreated()}</HeadCol>
      <HeadCol>{$LL.expireDate()}</HeadCol>
      <HeadCol />
    </tr>
    <tbody slot="body" let:class="{className}" class="{className}">
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
              editBtnEnabled="{false}"
              deleteBtnClickHandler="{toggleDeleteInvite(item.invite_id)}"
            />
          </RowCol>
        </TableRow>
      {/each}
    </tbody>
  </Table>
</TableContainer>

{#if showDeleteInvite}
  <DeleteConfirmation
    toggleDelete="{toggleDeleteInvite(null)}"
    handleDelete="{handleDeleteInvite}"
    confirmText="{$LL.userInviteConfirmDelete()}"
    confirmBtnText="{$LL.userInviteDelete()}"
  />
{/if}

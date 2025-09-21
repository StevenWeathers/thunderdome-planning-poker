<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import Table from '../../../components/table/Table.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import CrudActions from '../../../components/table/CrudActions.svelte';
  import BooleanDisplay from '../../../components/global/BooleanDisplay.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';
  import UpdateSupportTicket from '../../../components/admin/UpdateSupportTicket.svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const ticketsPageLimit = 25;
  let ticketCount = $state(0);

  const defaultTicket = {
    id: '',
    name: '',
    type: '',
    content: '',
    active: false,
    registeredOnly: false,
    allowDismiss: true,
  };

  let tickets = $state([]);
  let ticketsPage = $state(1);
  let showTicketUpdate = $state(false);
  let showDeleteTicket = $state(false);
  let selectedTicket = $state({ ...defaultTicket });
  let deleteTicketId = $state(null);
  let adminUsers = $state([]);
  let loaded = $state(false);

  function getAdminUsers() {
    return xfetch('/api/admin/admin-users?limit=1000&offset=0')
      .then(res => res.json())
      .then(function (result) {
        adminUsers = result.data;
      })
      .catch(function () {
        notifications.danger('Error fetching admin users');
      });
  }

  const toggleUpdateTicket = ticket => () => {
    showTicketUpdate = !showTicketUpdate;
    selectedTicket = ticket;
  };

  const toggleDeleteTicket = ticketId => () => {
    showDeleteTicket = !showDeleteTicket;
    deleteTicketId = ticketId;
  };

  function updateTicket(id, body) {
    xfetch(`/api/admin/support-tickets/${id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.updateSupportTicketSuccess());
        getTickets();
        toggleUpdateTicket({ ...defaultTicket })();
      })
      .catch(function () {
        notifications.danger($LL.updateSupportTicketError());
      });
  }

  function getTickets() {
    const ticketsOffset = (ticketsPage - 1) * ticketsPageLimit;
    xfetch(`/api/admin/support-tickets?limit=${ticketsPageLimit}&offset=${ticketsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        tickets = result.data;
        ticketCount = result.meta.count || 0;
      })
      .catch(function () {
        notifications.danger($LL.getSupportTicketsError());
      });
  }

  function handleDeleteTicket() {
    xfetch(`/api/admin/support-tickets/${deleteTicketId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        getTickets();
        toggleDeleteTicket(null)();
        notifications.success($LL.deleteSupportTicketSuccess());
      })
      .catch(function () {
        notifications.danger($LL.deleteSupportTicketError());
      });
  }

  const changePage = (evt: CustomEvent) => {
    ticketsPage = evt.detail;
    getTickets();
  };

  async function loadData() {
    await getAdminUsers();
    await getTickets();
    loaded = true;
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

    loadData();
  });
</script>

<svelte:head>
  <title>{$LL.supportTickets()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="support-tickets">
  <TableContainer>
    <TableNav title={$LL.supportTickets()} createBtnEnabled={false} />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>Full Name</HeadCol>
          <HeadCol>Email</HeadCol>
          <HeadCol>Inquiry</HeadCol>
          <HeadCol>Assigned To</HeadCol>
          <HeadCol>Resolved</HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
          <HeadCol type="action">
            <span class="sr-only">Actions</span>
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#if !loaded}
            <TableRow>
              <RowCol colspan="7" class="text-center">
                {$LL.loading()}
              </RowCol>
            </TableRow>
          {:else}
            {#each tickets as ticket, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  {#if ticket.userId === null}
                    {ticket.fullName}
                  {:else}
                    <a
                      href="{appRoutes.adminUsers}/{ticket.userId}"
                      class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >
                      {ticket.fullName}
                    </a>
                  {/if}
                </RowCol>
                <RowCol>
                  {ticket.email}
                </RowCol>
                <RowCol>
                  {ticket.inquiry.substring(0, 50)}{ticket.inquiry.length > 50 ? '...' : ''}
                </RowCol>
                <RowCol>
                  {#if ticket.assignedTo !== null && ticket.assignedTo !== ''}
                    {@const admin = adminUsers.find(u => u.id === ticket.assignedTo)}
                    <a
                      href="{appRoutes.adminUsers}/{admin.id}"
                      class="flex items-center text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >
                      <UserAvatar
                        warriorId={admin.id}
                        pictureUrl={admin.picture}
                        gravatarHash={admin.gravatarHash}
                        avatar={admin.avatar}
                        userName={admin.name}
                      />
                      {admin?.name || 'Unknown Admin'}
                    </a>
                  {:else}
                    <span class="text-gray-500 dark:text-gray-400 italic">
                      {$LL.unassigned()}
                    </span>
                  {/if}
                </RowCol>
                <RowCol>
                  <BooleanDisplay boolValue={ticket.resolvedAt !== null} />
                </RowCol>
                <RowCol>
                  {new Date(ticket.updatedAt).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  <CrudActions
                    detailsLink={`${appRoutes.adminSupportTickets}/${ticket.id}`}
                    editBtnClickHandler={toggleUpdateTicket(ticket)}
                    deleteBtnClickHandler={toggleDeleteTicket(ticket.id)}
                  />
                </RowCol>
              </TableRow>
            {/each}
          {/if}
        </tbody>
      {/snippet}
    </Table>
    <TableFooter
      bind:current={ticketsPage}
      num_items={ticketCount}
      per_page={ticketsPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>

  {#if showTicketUpdate}
    <UpdateSupportTicket
      toggleUpdate={toggleUpdateTicket({ ...defaultTicket })}
      handleUpdate={updateTicket}
      ticket={selectedTicket}
      {adminUsers}
      {xfetch}
      {notifications}
    />
  {/if}

  {#if showDeleteTicket}
    <DeleteConfirmation
      toggleDelete={toggleDeleteTicket(null)}
      handleDelete={handleDeleteTicket}
      confirmText={$LL.supportTicketDeleteConfirmation()}
      confirmBtnText={$LL.supportTicketDelete()}
    />
  {/if}
</AdminPageLayout>

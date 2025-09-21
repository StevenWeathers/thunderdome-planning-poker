<script lang="ts">
  import { onMount } from 'svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import LL from '../../../i18n/i18n-svelte';
  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';
  import type { SupportTicket, User } from '../../../types/user';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import { appRoutes } from '../../../config';
  import { user } from '../../../stores';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import { MessageCircleMore, SquarePen, UserIcon, UserPlus, ArrowLeft, Check, Clock } from 'lucide-svelte';
  import UpdateSupportTicket from '../../../components/admin/UpdateSupportTicket.svelte';

  interface Props {
    ticketId: string;
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
  }

  const { ticketId, xfetch, notifications, router }: Props = $props();

  let adminUsers = $state([]);
  let loaded = $state(false);
  let ticket: SupportTicket | null = $state({
    id: '',
    userId: null,
    fullName: '',
    email: '',
    inquiry: '',
    assignedTo: null,
    notes: '',
    resolvedAt: null,
    resolvedBy: null,
    createdAt: '',
    updatedAt: '',
  });
  let assignedAdmin: User | null = $state(null);
  let resolvedByAdmin: User | null = $state(null);
  let showTicketUpdate = $state(false);

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

  function getSupportTicket() {
    return xfetch(`/api/admin/support-tickets/${ticketId}`)
      .then(res => {
        if (res.status === 404) {
          notifications.danger('Support ticket not found');
          router.navigate('/admin/support-tickets');
          return;
        }
        return res.json();
      })
      .then(function (result) {
        ticket = result.data as SupportTicket;
        assignedAdmin = adminUsers.find((user: User) => user.id === ticket?.assignedTo) || null;
        resolvedByAdmin = adminUsers.find((user: User) => user.id === ticket?.resolvedBy) || null;
      })
      .catch(function () {
        notifications.danger('Error fetching support ticket');
      });
  }

  const toggleUpdateTicket = () => {
    showTicketUpdate = !showTicketUpdate;
  };

  function updateTicket(id: string, body: any) {
    xfetch(`/api/admin/support-tickets/${id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function () {
        notifications.success($LL.updateSupportTicketSuccess());
        getSupportTicket();
        toggleUpdateTicket();
      })
      .catch(function () {
        notifications.danger($LL.updateSupportTicketError());
      });
  }

  onMount(async () => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    await getAdminUsers();
    await getSupportTicket();
    loaded = true;
  });
</script>

<svelte:head>
  <title>{$LL.supportTickets()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="support-tickets">
  <!-- Hero Header Section with Gradient Background -->
  <div
    class="relative overflow-hidden bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 dark:from-gray-900 dark:via-blue-900/20 dark:to-indigo-900/30 rounded-2xl border border-slate-200/60 dark:border-gray-700/60 mb-8"
  >
    <div
      class="absolute inset-0 bg-grid-slate-100/50 dark:bg-grid-slate-700/25 [mask-image:linear-gradient(0deg,white,rgba(255,255,255,0.6))] dark:[mask-image:linear-gradient(0deg,rgba(255,255,255,0.1),rgba(255,255,255,0.5))]"
    ></div>
    <div class="relative px-6 py-8 sm:px-8">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6">
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-blue-100 dark:bg-blue-900/30 rounded-xl">
              <MessageCircleMore class="w-6 h-6 text-blue-600 dark:text-blue-400" />
            </div>
            <h1
              class="text-3xl sm:text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-indigo-800 dark:from-white dark:via-blue-100 dark:to-indigo-100 bg-clip-text text-transparent"
            >
              Support Ticket Details
            </h1>
          </div>
          <div class="flex items-center gap-2">
            <span
              class="inline-flex items-center gap-1.5 px-3 py-1 text-sm font-medium text-blue-700 bg-blue-100/80 dark:text-blue-300 dark:bg-blue-900/40 rounded-full border border-blue-200 dark:border-blue-700"
            >
              <span class="w-1.5 h-1.5 bg-blue-500 rounded-full"></span>
              Ticket ID: {ticketId}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            onclick={() => router.route(appRoutes.adminSupportTickets)}
            class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-white dark:bg-gray-800 dark:text-gray-300 rounded-lg border border-gray-200 dark:border-gray-700 shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <ArrowLeft class="w-4 h-4" />
            Back to Tickets
          </button>
          <button
            onclick={toggleUpdateTicket}
            class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 rounded-lg border border-blue-700 dark:border-blue-400 shadow-sm hover:shadow-md transition-all duration-200"
          >
            <SquarePen class="w-4 h-4" />
            Update Ticket
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Loading State with Enhanced Animation -->
  {#if !loaded}
    <div class="flex items-center justify-center min-h-96" role="status" aria-label="Loading support ticket details">
      <div class="flex flex-col items-center gap-6">
        <div class="relative">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-blue-100 dark:border-blue-900"></div>
          <div
            class="animate-spin rounded-full h-12 w-12 border-4 border-blue-600 dark:border-blue-400 border-t-transparent absolute top-0 left-0"
          ></div>
        </div>
        <div class="text-center space-y-2">
          <p class="text-lg font-medium text-gray-900 dark:text-white">Loading ticket details</p>
          <p class="text-sm text-gray-500 dark:text-gray-400">Please wait while we fetch the information...</p>
        </div>
      </div>
    </div>
  {:else}
    <!-- Main Content with Enhanced Layout -->
    <div class="max-w-6xl mx-auto space-y-8">
      <!-- Enhanced Status Badge Section -->
      <div class="flex items-center justify-center sm:justify-start">
        {#if ticket?.resolvedAt}
          <div
            class="group inline-flex items-center gap-3 px-6 py-3 text-sm font-semibold text-emerald-800 bg-gradient-to-r from-emerald-50 to-green-50 dark:from-emerald-900/20 dark:to-green-900/20 dark:text-emerald-200 rounded-2xl border border-emerald-200 dark:border-emerald-700/50 shadow-sm hover:shadow-md transition-all duration-200"
          >
            <div class="relative">
              <Check class="w-5 h-5" />
            </div>
            <span class="tracking-wide">Resolved</span>
          </div>
        {:else}
          <div
            class="group inline-flex items-center gap-3 px-6 py-3 text-sm font-semibold text-amber-800 bg-gradient-to-r from-amber-50 to-orange-50 dark:from-amber-900/20 dark:to-orange-900/20 dark:text-amber-200 rounded-2xl border border-amber-200 dark:border-amber-700/50 shadow-sm hover:shadow-md transition-all duration-200"
          >
            <div class="relative">
              <Clock class="w-5 h-5" />
            </div>
            <span class="tracking-wide">Open</span>
          </div>
        {/if}
      </div>

      <!-- Enhanced Information Cards Grid -->
      <div class="grid gap-8 lg:grid-cols-2">
        <!-- Customer Information Card -->
        <div
          class="group bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200/60 dark:border-gray-700/60 overflow-hidden hover:shadow-xl transition-all duration-300"
        >
          <div
            class="px-8 py-6 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 border-b border-gray-200/60 dark:border-gray-700/60"
          >
            <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
              <div class="p-2.5 bg-blue-100 dark:bg-blue-900/50 rounded-xl">
                <UserIcon class="w-6 h-6 text-blue-600 dark:text-blue-400" />
              </div>
              Customer Information
            </h2>
          </div>
          <div class="p-8 space-y-6">
            <div class="space-y-3">
              <div class="block text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Full Name
              </div>
              <div
                class="p-4 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-gray-200/60 dark:border-gray-700/60"
              >
                {#if ticket?.userId === null}
                  <p class="text-lg font-medium text-gray-900 dark:text-white">{ticket?.fullName}</p>
                {:else}
                  <a
                    href="{appRoutes.adminUsers}/{ticket?.userId}"
                    class="text-lg font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors duration-200 hover:underline decoration-2 underline-offset-2"
                  >
                    {ticket?.fullName}
                  </a>
                {/if}
              </div>
            </div>
            <div class="space-y-3">
              <div class="block text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Email Address
              </div>
              <div
                class="p-4 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-gray-200/60 dark:border-gray-700/60"
              >
                <a
                  href="mailto:{ticket?.email}"
                  class="text-lg font-medium text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 transition-colors duration-200 break-all hover:underline decoration-2 underline-offset-2"
                  aria-label="Send email to {ticket?.email}"
                >
                  {ticket?.email}
                </a>
              </div>
            </div>
          </div>
        </div>

        <!-- Assignment Information Card -->
        <div
          class="group bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200/60 dark:border-gray-700/60 overflow-hidden hover:shadow-xl transition-all duration-300"
        >
          <div
            class="px-8 py-6 bg-gradient-to-r from-purple-50 to-pink-50 dark:from-purple-900/30 dark:to-pink-900/30 border-b border-gray-200/60 dark:border-gray-700/60"
          >
            <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
              <div class="p-2.5 bg-purple-100 dark:bg-purple-900/50 rounded-xl">
                <UserPlus class="w-6 h-6 text-purple-600 dark:text-purple-400" />
              </div>
              Assignment & Status
            </h2>
          </div>
          <div class="p-8 space-y-6">
            <div class="space-y-3">
              <div class="block text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Assigned To
              </div>
              <div
                class="p-4 bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-gray-200/60 dark:border-gray-700/60"
              >
                {#if ticket?.assignedTo && assignedAdmin}
                  <div class="flex items-center gap-4">
                    <div class="w-14 h-14 ring-2 ring-purple-200 dark:ring-purple-700 rounded-full p-0.5">
                      <UserAvatar
                        warriorId={assignedAdmin.id}
                        pictureUrl={assignedAdmin.picture}
                        gravatarHash={assignedAdmin.gravatarHash}
                        avatar={assignedAdmin.avatar}
                        userName={assignedAdmin.name}
                      />
                    </div>
                    <div>
                      <a
                        href="{appRoutes.adminUsers}/{assignedAdmin.id}"
                        class="text-lg font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors duration-200 hover:underline decoration-2 underline-offset-2"
                      >
                        {assignedAdmin.name}
                      </a>
                    </div>
                  </div>
                {:else}
                  <p class="text-lg font-medium text-gray-500 dark:text-gray-400 italic">Unassigned</p>
                {/if}
              </div>
            </div>

            {#if ticket?.resolvedAt}
              <div class="space-y-3">
                <div class="block text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Resolved By
                </div>
                <div
                  class="p-4 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200/60 dark:border-emerald-700/60"
                >
                  <div class="flex items-center gap-4">
                    <div class="w-14 h-14 ring-2 ring-emerald-200 dark:ring-emerald-700 rounded-full p-0.5">
                      <UserAvatar
                        warriorId={resolvedByAdmin?.id}
                        pictureUrl={resolvedByAdmin?.picture}
                        gravatarHash={resolvedByAdmin?.gravatarHash}
                        avatar={resolvedByAdmin?.avatar}
                        userName={resolvedByAdmin?.name}
                      />
                    </div>
                    <div>
                      <a
                        href="{appRoutes.adminUsers}/{resolvedByAdmin?.id}"
                        class="text-lg font-medium text-emerald-600 hover:text-emerald-800 dark:text-emerald-400 dark:hover:text-emerald-300 transition-colors duration-200 hover:underline decoration-2 underline-offset-2"
                      >
                        {resolvedByAdmin?.name || 'Unknown'}
                      </a>
                    </div>
                  </div>
                </div>
              </div>
              <div class="space-y-3">
                <div class="block text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                  Resolved At
                </div>
                <div
                  class="p-4 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200/60 dark:border-emerald-700/60"
                >
                  <p class="text-lg font-medium text-gray-900 dark:text-white">
                    {new Date(ticket.resolvedAt).toLocaleString()}
                  </p>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Enhanced Inquiry Section -->
      <div
        class="group bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200/60 dark:border-gray-700/60 overflow-hidden hover:shadow-xl transition-all duration-300"
      >
        <div
          class="px-8 py-6 bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-900/30 dark:to-emerald-900/30 border-b border-gray-200/60 dark:border-gray-700/60"
        >
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
            <div class="p-2.5 bg-green-100 dark:bg-green-900/50 rounded-xl">
              <MessageCircleMore class="w-6 h-6 text-green-600 dark:text-green-400" />
            </div>
            Customer Inquiry
          </h2>
        </div>
        <div class="p-8">
          <div
            class="bg-gradient-to-br from-gray-50 to-blue-50/30 dark:from-gray-900/50 dark:to-blue-900/10 rounded-xl border border-gray-200/60 dark:border-gray-700/60 p-6"
          >
            <div class="prose prose-gray dark:prose-invert max-w-none">
              <p class="text-gray-900 dark:text-white whitespace-pre-line leading-relaxed text-base">
                {ticket?.inquiry}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Enhanced Notes Section -->
      <div
        class="group bg-white dark:bg-gray-800 rounded-2xl shadow-lg border border-gray-200/60 dark:border-gray-700/60 overflow-hidden hover:shadow-xl transition-all duration-300"
      >
        <div
          class="px-8 py-6 bg-gradient-to-r from-orange-50 to-red-50 dark:from-orange-900/30 dark:to-red-900/30 border-b border-gray-200/60 dark:border-gray-700/60"
        >
          <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-3">
            <div class="p-2.5 bg-orange-100 dark:bg-orange-900/50 rounded-xl">
              <SquarePen class="w-6 h-6 text-orange-600 dark:text-orange-400" />
            </div>
            Internal Notes
          </h2>
        </div>
        <div class="p-8">
          {#if ticket?.notes && ticket?.notes.trim()}
            <div
              class="bg-gradient-to-br from-amber-50 to-orange-50/30 dark:from-amber-900/20 dark:to-orange-900/10 rounded-xl border border-amber-200/60 dark:border-amber-700/60 p-6"
            >
              <div class="prose prose-gray dark:prose-invert max-w-none">
                <p class="text-gray-900 dark:text-white whitespace-pre-line leading-relaxed text-base">
                  {ticket.notes}
                </p>
              </div>
            </div>
          {:else}
            <div
              class="bg-gray-50 dark:bg-gray-900/50 rounded-xl border border-gray-200/60 dark:border-gray-700/60 p-8 text-center"
            >
              <div
                class="inline-flex items-center justify-center w-16 h-16 bg-gray-100 dark:bg-gray-800 rounded-full mb-4"
              >
                <SquarePen class="w-8 h-8 text-gray-400" />
              </div>
              <p class="text-lg font-medium text-gray-500 dark:text-gray-400">No notes available</p>
              <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">Internal notes will appear here when added</p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}

  {#if showTicketUpdate}
    <UpdateSupportTicket
      toggleUpdate={toggleUpdateTicket}
      handleUpdate={updateTicket}
      ticket
      {adminUsers}
      {xfetch}
      {notifications}
    />
  {/if}
</AdminPageLayout>

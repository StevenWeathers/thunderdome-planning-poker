<script lang="ts">
  import { onMount } from 'svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import LL from '../../../i18n/i18n-svelte';
  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';
  import type { User } from '../../../types/user';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import { appRoutes } from '../../../config';
  import { user } from '../../../stores';
  import { validateUserIsAdmin } from '../../../validationUtils';

  interface Props {
    ticketId: string;
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
  }

  const { ticketId, xfetch, notifications, router }: Props = $props();

  let adminUsers = $state([]);
  let loaded = $state(false);
  let ticket = $state({
    id: '',
    fullName: '',
    email: '',
    inquiry: '',
    assignedTo: null,
    notes: '',
    resolvedAt: null,
    resolvedBy: null,
  });
  let assignedAdmin: User | null = $state(null);
  let resolvedByAdmin: User | null = $state(null);

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
        ticket = result.data;
        assignedAdmin = adminUsers.find(user => user.id === ticket.assignedTo) || null;
        resolvedByAdmin = adminUsers.find(user => user.id === ticket.resolvedBy) || null;
      })
      .catch(function () {
        notifications.danger('Error fetching support ticket');
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
  <!-- Header Section -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
    <div>
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">Support Ticket Details</h1>
      <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
        Ticket ID: #{ticketId}
      </p>
    </div>

    <!-- Back Button -->
    <a
      href={appRoutes.adminSupportTickets}
      class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 dark:hover:bg-gray-700 dark:focus:ring-blue-400 transition-colors duration-200"
      aria-label="Go back to support tickets list"
    >
      <svg class="w-4 h-4 rtl:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
      </svg>
      Back to Tickets
    </a>
  </div>

  <!-- Loading State -->
  {#if !loaded}
    <div class="flex items-center justify-center min-h-64" role="status" aria-label="Loading support ticket details">
      <div class="flex flex-col items-center gap-4">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 dark:border-blue-400"></div>
        <p class="text-gray-500 dark:text-gray-400">Loading ticket details...</p>
      </div>
    </div>
  {:else}
    <!-- Main Content -->
    <div class="max-w-4xl mx-auto">
      <!-- Status Badge Section -->
      <div class="mb-6">
        {#if ticket.resolvedAt}
          <div
            class="inline-flex items-center gap-2 px-3 py-1 text-sm font-medium text-green-800 bg-green-100 rounded-full dark:text-green-200 dark:bg-green-900/30"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clip-rule="evenodd"
              />
            </svg>
            Resolved
          </div>
        {:else}
          <div
            class="inline-flex items-center gap-2 px-3 py-1 text-sm font-medium text-amber-800 bg-amber-100 rounded-full dark:text-amber-200 dark:bg-amber-900/30"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                clip-rule="evenodd"
              />
            </svg>
            Open
          </div>
        {/if}
      </div>

      <!-- Ticket Information Cards -->
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- Customer Information Card -->
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden"
        >
          <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <svg
                class="w-5 h-5 text-blue-600 dark:text-blue-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                />
              </svg>
              Customer Information
            </h2>
          </div>
          <div class="p-6 space-y-4">
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-500 dark:text-gray-400">Full Name</label>
              <p class="text-base text-gray-900 dark:text-white font-medium">
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
              </p>
            </div>
            <div class="space-y-1">
              <label class="block text-sm font-medium text-gray-500 dark:text-gray-400">Email Address</label>
              <p class="text-base text-gray-900 dark:text-white">
                <a
                  href="mailto:{ticket.email}"
                  class="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 transition-colors duration-200 underline decoration-dotted underline-offset-2"
                  aria-label="Send email to {ticket.email}"
                >
                  {ticket.email}
                </a>
              </p>
            </div>
          </div>
        </div>

        <!-- Assignment Information Card -->
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden"
        >
          <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <svg
                class="w-5 h-5 text-purple-600 dark:text-purple-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
                />
              </svg>
              Assignment & Status
            </h2>
          </div>
          <div class="p-6 space-y-4">
            <div class="space-y-1">
              <div class="block text-sm font-medium text-gray-500 dark:text-gray-400">Assigned To</div>
              {#if ticket.assignedTo && assignedAdmin}
                <div class="flex items-center gap-2">
                  <div class="w-12 h-12">
                    <UserAvatar
                      warriorId={assignedAdmin.id}
                      pictureUrl={assignedAdmin.picture}
                      gravatarHash={assignedAdmin.gravatarHash}
                      avatar={assignedAdmin.avatar}
                      userName={assignedAdmin.name}
                    />
                  </div>
                  <div class="flex text-base text-gray-900 dark:text-white font-medium">
                    <a
                      href="{appRoutes.adminUsers}/{assignedAdmin.id}"
                      class="flex items-center text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                      >{assignedAdmin.name}</a
                    >
                  </div>
                </div>
              {:else}
                <p class="text-base text-gray-500 dark:text-gray-400 italic">Unassigned</p>
              {/if}
            </div>

            {#if ticket.resolvedAt}
              <div class="space-y-1">
                <label class="block text-sm font-medium text-gray-500 dark:text-gray-400">Resolved By</label>
                <div class="flex items-center gap-2">
                  <div class="w-12 h-12">
                    <UserAvatar
                      warriorId={resolvedByAdmin.id}
                      pictureUrl={resolvedByAdmin.picture}
                      gravatarHash={resolvedByAdmin.gravatarHash}
                      avatar={resolvedByAdmin.avatar}
                      userName={resolvedByAdmin.name}
                    />
                  </div>
                  <div class="flex text-base text-gray-900 dark:text-white font-medium">
                    <a
                      href="{appRoutes.adminUsers}/{resolvedByAdmin.id}"
                      class="flex items-center text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                      >{resolvedByAdmin?.name || 'Unknown'}</a
                    >
                  </div>
                </div>
              </div>
              <div class="space-y-1">
                <label class="block text-sm font-medium text-gray-500 dark:text-gray-400">Resolved At</label>
                <p class="text-base text-gray-900 dark:text-white">
                  {new Date(ticket.resolvedAt).toLocaleString()}
                </p>
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Inquiry Section -->
      <div class="mt-6">
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden"
        >
          <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <svg
                class="w-5 h-5 text-green-600 dark:text-green-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
              Customer Inquiry
            </h2>
          </div>
          <div class="p-6">
            <div class="prose prose-gray dark:prose-invert max-w-none">
              <p class="text-gray-900 dark:text-white whitespace-pre-line leading-relaxed">
                {ticket.inquiry}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Notes Section -->
      <div class="mt-6">
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden"
        >
          <div class="px-6 py-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <svg
                class="w-5 h-5 text-orange-600 dark:text-orange-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                />
              </svg>
              Internal Notes
            </h2>
          </div>
          <div class="p-6">
            {#if ticket.notes && ticket.notes.trim()}
              <div class="prose prose-gray dark:prose-invert max-w-none">
                <p class="text-gray-900 dark:text-white whitespace-pre-line leading-relaxed">
                  {ticket.notes}
                </p>
              </div>
            {:else}
              <p class="text-gray-500 dark:text-gray-400 italic">No notes available</p>
            {/if}
          </div>
        </div>
      </div>
    </div>
  {/if}
</AdminPageLayout>

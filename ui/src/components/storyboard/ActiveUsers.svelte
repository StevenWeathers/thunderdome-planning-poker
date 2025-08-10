<script lang="ts">
  import { onMount } from 'svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import { Check, CopyIcon, Minus, Plus, Users } from 'lucide-svelte';

  // TypeScript interfaces
  interface User {
    id: string;
    name: string;
    avatar: string;
    gravatarHash: string;
    active: boolean;
  }

  // Props
  let { 
    users = [] as User[],
    facilitatorIds = [] as string[],
    inviteUrl = '',
    onAddFacilitator,
    onRemoveFacilitator,
    isFacilitator = false
  }: {
    users?: User[];
    facilitatorIds?: string[];
    inviteUrl?: string;
    onAddFacilitator?: (userId: string) => void;
    onRemoveFacilitator?: (userId: string) => void;
    currentUserId?: string | null;
    isFacilitator?: boolean;
  } = $props();

  // State
  let showInviteForm = $state(false);
  let copied = $state(false);

  let displayUsers = users.filter(user => user.active);

  // Helper function to check if a user is a facilitator
  function isUserFacilitator(userId: string): boolean {
    return facilitatorIds.includes(userId);
  }

  async function copyInviteLink(): Promise<void> {
    try {
      await navigator.clipboard.writeText(inviteUrl);
      copied = true;
      setTimeout(() => copied = false, 2000);
    } catch (err) {
      console.error('Failed to copy invite link:', err);
    }
  }

  function handleFacilitatorAction(user: User): void {
    const currentlyFacilitator = isUserFacilitator(user.id);
    if (currentlyFacilitator) {
      onRemoveFacilitator?.(user.id);
    } else {
      onAddFacilitator?.(user.id);
    }
  }
</script>

<div class="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700">
  <!-- Header -->
  <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="p-2 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <Users class="inline-block w-5 h-5 text-blue-600 dark:text-blue-400" />
        </div>
        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            Active Users
          </h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            {displayUsers.length} {displayUsers.length === 1 ? 'participant' : 'participants'} online
          </p>
        </div>
      </div>
      
      <button
        onclick={() => showInviteForm = !showInviteForm}
        class="inline-flex items-center px-3 py-2 border border-transparent leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:focus:ring-offset-gray-900 transition-colors duration-200"
      >
        <Plus class="inline-block w-4 h-4 me-2" />
        Invite User
      </button>
    </div>
  </div>

  <!-- Invite Form (Collapsible) -->
  {#if showInviteForm}
    <div class="px-6 py-4 bg-gray-50 dark:bg-gray-800/50 border-b border-gray-200 dark:border-gray-700">
      <!-- Copy Link Section -->
      <div>
        <label for="invite-link" class="block font-medium text-gray-700 dark:text-gray-300 mb-2">
          Share Invite Link
        </label>
        <div class="flex space-x-1">
          <input
            id="invite-link"
            type="text"
            readonly
            value={inviteUrl}
            class="flex-1 min-w-0 block w-full px-3 py-2 rounded-none rounded-s-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
          />
          <button
            onclick={copyInviteLink}
            class="inline-flex items-center px-3 py-2 border border-s-0 border-gray-300 dark:border-gray-600 rounded-e-md bg-white dark:bg-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200"
          >
            {#if copied}
              <Check class="w-4 h-4 text-green-500" />
            {:else}
              <CopyIcon class="w-4 h-4" />
            {/if}
          </button>
        </div>
        {#if copied}
          <p class="mt-1 text-green-600 dark:text-green-400">Link copied to clipboard!</p>
        {/if}
      </div>
    </div>
  {/if}

  <!-- User Grid -->
  <div class="px-6 py-4">
    <div class="flex flex-wrap gap-4 justify-center">
      {#each displayUsers as user (user.id)}
        <div class="group">
          <!-- User Card -->
          <div class="flex flex-col items-center p-3 rounded-lg border border-gray-300 dark:border-gray-800 hover:bg-gray-100 dark:hover:bg-gray-800/50 hover:shadow-md transition-all duration-200 min-w-0 w-full max-w-[140px]">
            <!-- Avatar -->
            <div class="mb-2">
              <UserAvatar
                    warriorId={user.id}
                    avatar={user.avatar}
                    gravatarHash={user.gravatarHash}
                    userName={user.name}
                />
            </div>

            <!-- User Info -->
            <div class="text-center min-w-0 w-full">
              <p class="text-base font-medium text-gray-900 dark:text-white truncate" title={user.name}>
                {user.name}
              </p>
            </div>

            <!-- Role Badge -->
            <div class="mt-2">
              {#if isUserFacilitator(user.id)}
                <span class="inline-flex items-center px-3 py-1 rounded-full font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                  Facilitator
                </span>
              {:else}
                <span class="inline-flex items-center px-3 py-1 rounded-full font-medium bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-200">
                  Participant
                </span>
              {/if}
            </div>

            <!-- Action Button -->
            {#if isFacilitator}
              <div class="mt-3">
                <button
                  onclick={() => handleFacilitatorAction(user)}
                  class="inline-flex items-center px-3 py-1.5 font-medium rounded-md transition-colors duration-200 {isUserFacilitator(user.id) ? 'text-red-700 bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:text-red-400 dark:hover:bg-red-900/30 border border-red-200 dark:border-red-800' : 'text-blue-700 bg-blue-50 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-400 dark:hover:bg-blue-900/30 border border-blue-200 dark:border-blue-800'}"
                  title={isUserFacilitator(user.id) ? 'Remove facilitator status' : 'Make facilitator'}
                >
                  {#if isUserFacilitator(user.id)}
                    <Minus class="inline-block w-4 h-4 me-1.5"/>
                    <span class="sr-only">Remove</span> Facilitator
                  {:else}
                    <Plus class="w-4 h-4 me-1.5" />
                    <span class="sr-only">Make</span> Facilitator
                  {/if}
                </button>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
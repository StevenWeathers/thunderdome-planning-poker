<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import { Trash2 } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { RetroAction } from '../../types/retro';

  type AssignableUser = {
    id: string;
    name: string;
    avatar?: string;
    gravatarHash?: string;
    pictureUrl?: string;
  };

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    toggleAssignees?: () => void;
    handleAssigneeAdd?: (retroId: string, actionId: string, userId: string) => void;
    handleAssigneeRemove?: (retroId: string, actionId: string, userId: string) => () => void;
    assignableUsers?: AssignableUser[];
    action?: RetroAction | null;
  }

  let {
    xfetch,
    notifications,
    toggleAssignees = () => {},
    handleAssigneeAdd = (retroId: string, actionId: string, userId: string) => {},
    handleAssigneeRemove = (retroId: string, actionId: string, userId: string) => () => {},
    assignableUsers = [] as AssignableUser[],
    action = null,
  }: Props = $props();

  const resolvedAction = $derived(
    action ?? {
      comments: [],
      id: '',
      retroId: '',
      teamId: '',
      content: '',
      completed: false,
      assignees: [],
    },
  );

  let fallbackUsers = $state<AssignableUser[]>([]);
  let fallbackUsersTeamId = $state('');
  let selectedAssignee = $state('');

  $effect(() => {
    if (assignableUsers.length > 0) {
      fallbackUsers = [];
      fallbackUsersTeamId = '';
      return;
    }

    const teamId = resolvedAction.teamId ?? '';
    if (teamId === '' || fallbackUsersTeamId === teamId) {
      return;
    }

    fallbackUsersTeamId = teamId;

    xfetch(`/api/teams/${teamId}/users?limit=1000&offset=0`)
      .then(res => res.json())
      .then(function (result) {
        fallbackUsers = result.data ?? [];
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  });

  const availableAssignableUsers = $derived(assignableUsers.length > 0 ? assignableUsers : fallbackUsers);

  const addAssignee = () => {
    handleAssigneeAdd(resolvedAction.retroId, resolvedAction.id, selectedAssignee);
    selectedAssignee = '';
  };
</script>

<Modal closeModal={toggleAssignees} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2" ariaLabel={$LL.assignees()}>
  <div class="mt-6 flex flex-col gap-4">
    <div>
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">{$LL.assignees()}</h2>
      <p class="mt-2 whitespace-pre-wrap break-words text-sm text-gray-600 dark:text-gray-400">
        {resolvedAction.content}
      </p>
    </div>

    {#if availableAssignableUsers.length}
      <div class="flex w-full gap-4">
        <div class="w-2/3">
          <SelectInput bind:value={selectedAssignee} id="assignee" name="assignee">
            <option value="" disabled>{$LL.assigneeSelectPlaceholder()}</option>
            {#each availableAssignableUsers as user}
              <option value={user.id}>
                {user.name}
              </option>
            {/each}
          </SelectInput>
        </div>
        <div class="w-1/3">
          <HollowButton onClick={addAssignee} disabled={selectedAssignee === ''}>
            {$LL.assigneeAdd()}
          </HollowButton>
        </div>
      </div>
    {/if}

    {#if resolvedAction.assignees.length}
      <div class="grid grid-cols-1 gap-3 border-t border-gray-200 pt-4 dark:border-gray-700 md:grid-cols-2">
        {#each resolvedAction.assignees as assignee}
          <div class="flex items-center text-gray-700 dark:text-gray-300">
            <div class="w-1/4">
              <UserAvatar
                warriorId={assignee.id}
                gravatarHash={assignee.gravatarHash}
                avatar={assignee.avatar}
                userName={assignee.name}
                class="inline-block me-2"
              />
            </div>
            <div class="w-2/4 text-lg">{assignee.name}</div>
            <div class="w-1/4 text-right">
              <HollowButton
                color="red"
                onClick={handleAssigneeRemove(resolvedAction.retroId, resolvedAction.id, assignee.id)}
              >
                <Trash2 />
              </HollowButton>
            </div>
          </div>
        {/each}
      </div>
    {:else}
      <div
        class="rounded-xl border border-dashed border-gray-300 bg-gray-50 p-6 text-sm text-gray-500 dark:border-gray-600 dark:bg-gray-800/40 dark:text-gray-400"
      >
        No assignees yet.
      </div>
    {/if}
  </div>
</Modal>

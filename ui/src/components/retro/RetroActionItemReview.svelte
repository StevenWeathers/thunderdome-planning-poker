<script lang="ts">
  import { onMount } from 'svelte';
  import TableNav from '../table/TableNav.svelte';
  import LL from '../../i18n/i18n-svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import RowCol from '../table/RowCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import BooleanDisplay from '../global/BooleanDisplay.svelte';
  import { MessageSquareMore } from 'lucide-svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import TableFooter from '../table/TableFooter.svelte';
  import EditActionItem from './EditActionItem.svelte';
  import ActionComments from './ActionComments.svelte';
  import Table from '../table/Table.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    toggle?: any;
    team?: any;
  }

  let {
    xfetch,
    notifications,
    toggle = () => {},
    team = {
    id: '',
    name: '',
    organization_id: '',
    department_id: '',
    role: '',
  }
  }: Props = $props();

  const orgPrefix =
    team.organization_id !== ''
      ? `/api/organizations/${team.organization_id}`
      : '/api';
  const deptPrefix =
    team.department_id !== ''
      ? `${orgPrefix}/departments/${team.department_id}`
      : orgPrefix;
  const teamPrefix = `${deptPrefix}/teams/${team.id}`;

  const retroActionsPageLimit = 10;
  const usersPageLimit = 1000;

  let totalRetroActions = $state(0);
  let retroActionsPage = $state(1);
  let actionItems = $state([]);
  let users = $state([]);
  let usersPage = 1;

  let organizationRole = '';
  let departmentRole = '';
  let teamRole = '';
  let isAdmin = false;
  let isTeamMember = false;

  function getTeamOpenActionItems() {
    const offset = (retroActionsPage - 1) * retroActionsPageLimit;

    xfetch(
      `${teamPrefix}/retro-actions?limit=${retroActionsPageLimit}&offset=${offset}&completed=false`,
    )
      .then(response => response.json())
      .then(res => {
        actionItems = res.data;
        totalRetroActions = res.meta.count;
      })
      .catch(err => {
        console.error(err);
        notifications.danger("Failed to fetch team's action items");
      });
  }

  let showRetroActionComments = $state(false);
  let selectedRetroAction = $state(null);
  const toggleRetroActionComments = id => () => {
    showRetroActionComments = !showRetroActionComments;
    selectedRetroAction = id;
  };

  let showRetroActionEdit = $state(false);
  let selectedAction = $state(null);
  const toggleRetroActionEdit = (retroId, id) => () => {
    showRetroActionEdit = !showRetroActionEdit;
    selectedAction =
      retroId !== null ? actionItems.find(r => r.id === id) : null;
  };

  function handleRetroActionEdit(action) {
    xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
      method: 'PUT',
      body: {
        content: action.content,
        completed: action.completed,
      },
    })
      .then(function () {
        getTeamOpenActionItems();
        toggleRetroActionEdit(null, null)();
        notifications.success($LL.updateActionItemSuccess());
      })
      .catch(function () {
        notifications.danger($LL.updateActionItemError());
      });
  }

  function handleRetroActionDelete(action) {
    return () => {
      xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
        method: 'DELETE',
      })
        .then(function () {
          getTeamOpenActionItems();
          toggleRetroActionEdit(null)();
          notifications.success($LL.deleteActionItemSuccess());
        })
        .catch(function () {
          notifications.danger($LL.deleteActionItemError());
        });
    };
  }

  function handleRetroActionAssigneeAdd(retroId, actionId, userId) {
    xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
      method: 'POST',
      body: {
        user_id: userId,
      },
    })
      .then(function () {
        getTeamOpenActionItems();
      })
      .catch(function () {});
  }

  function handleRetroActionAssigneeRemove(retroId, actionId, userId) {
    return () => {
      xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
        method: 'DELETE',
        body: {
          user_id: userId,
        },
      }).then(function () {
        getTeamOpenActionItems();
      });
    };
  }

  const changeRetroActionPage = evt => {
    retroActionsPage = evt.detail;
    getTeamOpenActionItems();
  };

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${teamPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  onMount(() => {
    getUsers();
    getTeamOpenActionItems();
  });
</script>

<div class="w-full flex justify-center">
  <div class="w-full md:w-4/5 lg:w-2/3">
    <TableContainer>
      <TableNav title="Open Action Items" createBtnEnabled={false} />
      <Table>
        {#snippet header()}
                <tr >
            <HeadCol>{$LL.actionItem()}</HeadCol>
            <HeadCol>{$LL.completed()}</HeadCol>
            <HeadCol>{$LL.comments()}</HeadCol>
            <HeadCol />
          </tr>
              {/snippet}
        {#snippet body({ class: className })}
                <tbody   class="{className}">
            {#each actionItems as item, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <div class="whitespace-pre-wrap">
                    {#each item.assignees as assignee}
                      <UserAvatar
                        warriorId={assignee.id}
                        gravatarHash={assignee.gravatarHash}
                        avatar={assignee.avatar}
                        userName={assignee.name}
                        width={24}
                        class="inline-block me-2"
                      />
                    {/each}{item.content}
                  </div>
                </RowCol>
                <RowCol>
                  <BooleanDisplay boolValue={item.completed} />
                </RowCol>
                <RowCol>
                  <MessageSquareMore
                    width="22"
                    height="22"
                    class="inline-block"
                  />
                  <button
                    class="text-lg text-blue-400 dark:text-sky-400"
                    onclick={toggleRetroActionComments(item.id)}
                  >
                    &nbsp;{item.comments.length}
                  </button>
                </RowCol>
                <RowCol type="action">
                  <CrudActions
                    editBtnClickHandler={toggleRetroActionEdit(
                      item.retroId,
                      item.id,
                    )}
                    deleteBtnEnabled={false}
                  />
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
              {/snippet}
      </Table>
      <TableFooter
        bind:current={retroActionsPage}
        num_items={totalRetroActions}
        per_page={retroActionsPageLimit}
        on:navigate={changeRetroActionPage}
      />
    </TableContainer>
  </div>
</div>

{#if showRetroActionEdit}
  <EditActionItem
    toggleEdit={toggleRetroActionEdit(null, null)}
    handleEdit={handleRetroActionEdit}
    handleDelete={handleRetroActionDelete}
    assignableUsers={users}
    action={selectedAction}
    handleAssigneeAdd={handleRetroActionAssigneeAdd}
    handleAssigneeRemove={handleRetroActionAssigneeRemove}
  />
{/if}

{#if showRetroActionComments}
  <ActionComments
    toggleComments={toggleRetroActionComments(null)}
    actions={actionItems}
    users={users}
    selectedActionId={selectedRetroAction}
    getRetrosActions={getTeamOpenActionItems}
    xfetch={xfetch}
    notifications={notifications}
    isAdmin={isAdmin}
  />
{/if}

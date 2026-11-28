<script lang="ts">
  import { onMount } from 'svelte';

  import Table from '../table/Table.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import UserTitleForm from './UserTitleForm.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { UserTitle } from '../../types/team';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    apiPrefix: string;
  }

  let { xfetch, notifications, apiPrefix }: Props = $props();

  let titles = $state<UserTitle[]>([]);
  let showForm = $state(false);
  let editTitle = $state<UserTitle | null>(null);
  let removeTitleId = $state<string | null>(null);

  function getTitles() {
    xfetch(`${apiPrefix}/user-titles`)
      .then(res => res.json())
      .then(result => {
        titles = result.data;
      })
      .catch(() => {
        notifications.danger($LL.userTitlesGetError());
      });
  }

  function openForm(userTitle: UserTitle | null) {
    editTitle = userTitle;
    showForm = true;
  }

  function closeForm() {
    showForm = false;
    editTitle = null;
  }

  function handleSaved() {
    closeForm();
    getTitles();
  }

  function handleRemove() {
    xfetch(`${apiPrefix}/user-titles/${removeTitleId}`, { method: 'DELETE' })
      .then(() => {
        removeTitleId = null;
        notifications.success($LL.userTitleDeleteSuccess());
        getTitles();
      })
      .catch(() => {
        notifications.danger($LL.userTitleDeleteError());
      });
  }

  onMount(getTitles);
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title={$LL.userTitles()}
      createBtnText={$LL.userTitleAdd()}
      createButtonHandler={() => openForm(null)}
      createBtnTestId="user-title-add"
    />
    <p class="px-4 pb-4 text-sm text-gray-500 dark:text-gray-400">{$LL.userTitlesHelp()}</p>
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>{$LL.userTitleLabel()}</HeadCol>
          <HeadCol type="action"><span class="sr-only">{$LL.actions()}</span></HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each titles as userTitle, i (userTitle.id)}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200" data-testid="user-title-name">
                  {userTitle.name}
                </div>
              </RowCol>
              <RowCol type="action">
                <CrudActions
                  editBtnClickHandler={() => openForm(userTitle)}
                  deleteBtnClickHandler={() => (removeTitleId = userTitle.id)}
                />
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  {#if showForm}
    <UserTitleForm userTitle={editTitle} toggleClose={closeForm} {handleSaved} {apiPrefix} {xfetch} {notifications} />
  {/if}

  {#if removeTitleId}
    <DeleteConfirmation
      toggleDelete={() => (removeTitleId = null)}
      handleDelete={handleRemove}
      permanent={true}
      confirmText={$LL.userTitleDeleteConfirm()}
      confirmBtnText={$LL.userTitleDelete()}
    />
  {/if}
</div>

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
  import RetroTagForm from './RetroTagForm.svelte';

  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { RetroTag } from '../../types/retro';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    isEntityAdmin?: boolean;
    apiPrefix?: string;
  }

  let { xfetch, notifications, isEntityAdmin = false, apiPrefix = '/api' }: Props = $props();

  let tags = $state<RetroTag[]>([]);
  let showForm = $state(false);
  let editTag = $state<RetroTag | null>(null);
  let removeTagId = $state<string | null>(null);

  const canManage = $derived(validateUserIsAdmin($user) || isEntityAdmin);

  function getTags() {
    xfetch(`${apiPrefix}/retro-tags`)
      .then(res => res.json())
      .then(result => {
        tags = result.data;
      })
      .catch(() => {
        notifications.danger('Failed to get retro tags');
      });
  }

  function openForm(tag: RetroTag | null) {
    editTag = tag;
    showForm = true;
  }

  function closeForm() {
    showForm = false;
    editTag = null;
  }

  function handleSaved() {
    closeForm();
    getTags();
  }

  function handleRemove() {
    xfetch(`${apiPrefix}/retro-tags/${removeTagId}`, { method: 'DELETE' })
      .then(() => {
        removeTagId = null;
        notifications.success('Retro tag successfully removed');
        getTags();
      })
      .catch(() => {
        notifications.danger('Error removing retro tag');
      });
  }

  onMount(getTags);
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Retro Tags"
      createBtnEnabled={canManage}
      createBtnText="Create Retro Tag"
      createButtonHandler={() => openForm(null)}
      createBtnTestId="retro-tag-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>Name</HeadCol>
          <HeadCol type="action"><span class="sr-only">Actions</span></HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each tags as tag, i (tag.id)}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200" data-testid="retro-tag-name">{tag.name}</div>
              </RowCol>
              <RowCol type="action">
                {#if canManage}
                  <CrudActions
                    editBtnClickHandler={() => openForm(tag)}
                    deleteBtnClickHandler={() => (removeTagId = tag.id)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  {#if showForm}
    <RetroTagForm tag={editTag} toggleClose={closeForm} {handleSaved} {apiPrefix} {xfetch} {notifications} />
  {/if}

  {#if removeTagId}
    <DeleteConfirmation
      toggleDelete={() => (removeTagId = null)}
      handleDelete={handleRemove}
      permanent={true}
      confirmText="Are you sure you want to remove this retro tag? Items already tagged will keep it."
      confirmBtnText="Remove Retro Tag"
    />
  {/if}
</div>

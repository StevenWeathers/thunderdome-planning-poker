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
  import HealthCheckTemplateForm from './HealthCheckTemplateForm.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { HealthCheckTemplate } from '../../types/healthcheck';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    isEntityAdmin?: boolean;
    apiPrefix: string;
  }

  let { xfetch, notifications, isEntityAdmin = false, apiPrefix }: Props = $props();

  let templates = $state<HealthCheckTemplate[]>([]);
  let showForm = $state(false);
  let editTemplate = $state<HealthCheckTemplate | null>(null);
  let removeTemplateId = $state<string | null>(null);

  const canManage = $derived(validateUserIsAdmin($user) || isEntityAdmin);

  function getTemplates() {
    xfetch(`${apiPrefix}/health-check-templates`)
      .then(res => res.json())
      .then(result => {
        templates = result.data;
      })
      .catch(() => {
        notifications.danger($LL.healthCheckTemplatesGetError());
      });
  }

  function openForm(template: HealthCheckTemplate | null) {
    editTemplate = template;
    showForm = true;
  }

  function closeForm() {
    showForm = false;
    editTemplate = null;
  }

  function handleSaved() {
    closeForm();
    getTemplates();
  }

  function handleRemove() {
    xfetch(`${apiPrefix}/health-check-templates/${removeTemplateId}`, { method: 'DELETE' })
      .then(() => {
        removeTemplateId = null;
        notifications.success($LL.healthCheckTemplateDeleteSuccess());
        getTemplates();
      })
      .catch(() => {
        notifications.danger($LL.healthCheckTemplateDeleteError());
      });
  }

  onMount(getTemplates);
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title={$LL.healthCheckTemplates()}
      createBtnEnabled={canManage}
      createBtnText={$LL.healthCheckTemplateCreate()}
      createButtonHandler={() => openForm(null)}
      createBtnTestId="health-check-template-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>{$LL.healthCheckTemplateName()}</HeadCol>
          <HeadCol>{$LL.healthCheckTemplateDescription()}</HeadCol>
          <HeadCol>{$LL.healthCheckAreas()}</HeadCol>
          <HeadCol type="action"><span class="sr-only">{$LL.actions()}</span></HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each templates as template, i (template.id)}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200" data-testid="health-check-template-name">
                  {template.name}
                </div>
              </RowCol>
              <RowCol>{template.description}</RowCol>
              <RowCol>{template.areas.map(a => a.name).join(', ')}</RowCol>
              <RowCol type="action">
                {#if canManage}
                  <CrudActions
                    editBtnClickHandler={() => openForm(template)}
                    deleteBtnClickHandler={() => (removeTemplateId = template.id)}
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
    <HealthCheckTemplateForm
      template={editTemplate}
      toggleClose={closeForm}
      {handleSaved}
      {apiPrefix}
      {xfetch}
      {notifications}
    />
  {/if}

  {#if removeTemplateId}
    <DeleteConfirmation
      toggleDelete={() => (removeTemplateId = null)}
      handleDelete={handleRemove}
      permanent={true}
      confirmText={$LL.healthCheckTemplateDeleteConfirm()}
      confirmBtnText={$LL.healthCheckTemplateDelete()}
    />
  {/if}
</div>

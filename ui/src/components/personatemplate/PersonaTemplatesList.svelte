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
  import PersonaTemplateForm from './PersonaTemplateForm.svelte';

  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { PersonaTemplate } from '../../types/storyboard';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    isEntityAdmin?: boolean;
    apiPrefix?: string;
  }

  let { xfetch, notifications, isEntityAdmin = false, apiPrefix = '/api' }: Props = $props();

  let templates = $state<PersonaTemplate[]>([]);
  let showForm = $state(false);
  let editTemplate = $state<PersonaTemplate | null>(null);
  let removeTemplateId = $state<string | null>(null);

  const canManage = $derived(validateUserIsAdmin($user) || isEntityAdmin);

  function getTemplates() {
    xfetch(`${apiPrefix}/persona-templates`)
      .then(res => res.json())
      .then(result => {
        templates = result.data;
      })
      .catch(() => {
        notifications.danger('Failed to get persona templates');
      });
  }

  function openForm(template: PersonaTemplate | null) {
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
    xfetch(`${apiPrefix}/persona-templates/${removeTemplateId}`, { method: 'DELETE' })
      .then(() => {
        removeTemplateId = null;
        notifications.success('Persona template successfully removed');
        getTemplates();
      })
      .catch(() => {
        notifications.danger('Error removing persona template');
      });
  }

  onMount(getTemplates);
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Persona Templates"
      createBtnEnabled={canManage}
      createBtnText="Create Persona Template"
      createButtonHandler={() => openForm(null)}
      createBtnTestId="persona-template-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>Name</HeadCol>
          <HeadCol>Description</HeadCol>
          <HeadCol>Personas</HeadCol>
          <HeadCol type="action"><span class="sr-only">Actions</span></HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each templates as template, i (template.id)}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200" data-testid="persona-template-name">
                  {template.name}
                </div>
              </RowCol>
              <RowCol>{template.description}</RowCol>
              <RowCol>{template.personas.map(p => p.name).join(', ')}</RowCol>
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
    <PersonaTemplateForm
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
      confirmText="Are you sure you want to remove this persona template? Storyboards that loaded it keep their personas."
      confirmBtnText="Remove Persona Template"
    />
  {/if}
</div>

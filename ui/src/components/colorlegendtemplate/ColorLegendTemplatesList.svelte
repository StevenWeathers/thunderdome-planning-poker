<script lang="ts">
  import Table from '../table/Table.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';

  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import CreateColorLegendTemplate from './CreateColorLegendTemplate.svelte';
  import UpdateColorLegendTemplate from './UpdateColorLegendTemplate.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { ColorLegendTemplate } from '../../types/storyboard';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    isEntityAdmin?: boolean;
    templates?: Array<ColorLegendTemplate>;
    apiPrefix?: string;
    getTemplates?: () => void;
  }

  let {
    xfetch,
    notifications,
    isEntityAdmin = false,
    templates = [],
    apiPrefix = '/api',
    getTemplates = () => {},
  }: Props = $props();

  let showAddTemplate = $state(false);
  let showUpdateTemplate = $state(false);
  let updateTemplate = $state<ColorLegendTemplate | null>(null);
  let showRemoveTemplate = $state(false);
  let removeTemplateId = $state<string | null>(null);

  function handleCreateTemplate() {
    getTemplates();
    toggleCreateTemplate();
  }

  function handleTemplateUpdate() {
    getTemplates();
    toggleUpdateTemplate(null)();
  }

  function handleTemplateRemove() {
    xfetch(`${apiPrefix}/color-legend-templates/${removeTemplateId}`, {
      method: 'DELETE',
    })
      .then(() => {
        toggleRemoveTemplate(null)();
        notifications.success('Color legend template successfully removed');
        getTemplates();
      })
      .catch(() => {
        notifications.danger('Error removing color legend template');
      });
  }

  function toggleCreateTemplate() {
    showAddTemplate = !showAddTemplate;
  }

  function cloneTemplate(template: ColorLegendTemplate): ColorLegendTemplate {
    return {
      ...template,
      colorLegend: template.colorLegend.map(legend => ({ ...legend })),
    };
  }

  const toggleUpdateTemplate = (template: ColorLegendTemplate | null) => () => {
    updateTemplate = template ? cloneTemplate(template) : null;
    showUpdateTemplate = !showUpdateTemplate;
  };

  const toggleRemoveTemplate = (templateId: string | null) => () => {
    showRemoveTemplate = !showRemoveTemplate;
    removeTemplateId = templateId;
  };

  let isAdmin = $derived(validateUserIsAdmin($user));
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Color Legend Templates"
      createBtnEnabled={isAdmin || isEntityAdmin}
      createBtnText="Create Color Legend Template"
      createButtonHandler={toggleCreateTemplate}
      createBtnTestId="color-legend-template-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>Name</HeadCol>
          <HeadCol>Description</HeadCol>
          <HeadCol>Legend</HeadCol>
          <HeadCol type="action"><span class="sr-only">Actions</span></HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each templates as template, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200">{template.name}</div>
              </RowCol>
              <RowCol>{template.description}</RowCol>
              <RowCol>
                <div class="flex flex-wrap gap-2">
                  {#each template.colorLegend as legend}
                    <div
                      class="inline-flex items-center gap-2 rounded-full border border-gray-200 px-2 py-1 text-xs text-gray-700 dark:border-gray-700 dark:text-gray-200"
                      title={legend.legend ? `${legend.color}: ${legend.legend}` : legend.color}
                    >
                      <span class="inline-block h-3 w-3 rounded-full colorcard-{legend.color}"></span>
                      <span>{legend.legend || legend.color}</span>
                    </div>
                  {/each}
                </div>
              </RowCol>
              <RowCol type="action">
                {#if isAdmin || isEntityAdmin}
                  <CrudActions
                    editBtnClickHandler={toggleUpdateTemplate(template)}
                    deleteBtnClickHandler={toggleRemoveTemplate(template.id)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  {#if showAddTemplate}
    <CreateColorLegendTemplate
      toggleCreate={toggleCreateTemplate}
      handleCreate={handleCreateTemplate}
      {apiPrefix}
      {xfetch}
      {notifications}
    />
  {/if}

  {#if showUpdateTemplate && updateTemplate}
    <UpdateColorLegendTemplate
      toggleUpdate={toggleUpdateTemplate(null)}
      handleUpdate={handleTemplateUpdate}
      templateId={updateTemplate.id}
      name={updateTemplate.name}
      description={updateTemplate.description}
      colorLegend={updateTemplate.colorLegend}
      {apiPrefix}
      {xfetch}
      {notifications}
    />
  {/if}

  {#if showRemoveTemplate}
    <DeleteConfirmation
      toggleDelete={toggleRemoveTemplate(null)}
      handleDelete={handleTemplateRemove}
      permanent={true}
      confirmText="Are you sure you want to remove this color legend template?"
      confirmBtnText="Remove Color Legend Template"
    />
  {/if}
</div>

<style lang="postcss">
  .colorcard-gray {
    @apply bg-gray-400;
  }
  .colorcard-red {
    @apply bg-red-400;
  }
  .colorcard-orange {
    @apply bg-orange-400;
  }
  .colorcard-yellow {
    @apply bg-yellow-400;
  }
  .colorcard-green {
    @apply bg-green-400;
  }
  .colorcard-teal {
    @apply bg-teal-400;
  }
  .colorcard-blue {
    @apply bg-blue-400;
  }
  .colorcard-indigo {
    @apply bg-indigo-400;
  }
  .colorcard-purple {
    @apply bg-purple-400;
  }
  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>

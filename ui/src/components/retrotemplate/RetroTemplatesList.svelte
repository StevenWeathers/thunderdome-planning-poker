<script lang="ts">
  import Table from '../table/Table.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import RowCol from '../table/RowCol.svelte';

  import LL from '../../i18n/i18n-svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import TableNav from '../table/TableNav.svelte';
  import CrudActions from '../table/CrudActions.svelte';
  import { createEventDispatcher } from 'svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import CreateRetroTemplate from './CreateRetroTemplate.svelte';
  import UpdateRetroTemplate from './UpdateRetroTemplate.svelte';
  import BooleanDisplay from '../global/BooleanDisplay.svelte';
  import TableFooter from '../table/TableFooter.svelte';
  import { Eye } from 'lucide-svelte';
  import ViewFormat from './ViewFormat.svelte';

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let organizationId;
  export let teamId;
  export let departmentId;
  export let isEntityAdmin = false;
  export let templates = [];
  export let apiPrefix = '/api';
  export let templateCount = 0;
  export let templatesPage = 1;
  export let templatesPageLimit = 10;
  export let changePage = () => {};
  export let getTemplates = () => {};

  const dispatch = createEventDispatcher();

  let showAddTemplate = false;
  let showUpdateTemplate = false;
  let updateTemplate = {};
  let showRemoveTemplate = false;
  let removeTemplateId = null;

  function handleCreateTemplate() {
    getTemplates();
    toggleCreateTemplate();
  }

  function handleTemplateUpdate() {
    getTemplates();
    toggleUpdateTemplate({})();
  }

  function handleTemplateRemove() {
    xfetch(`${apiPrefix}/retro-templates/${removeTemplateId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('remove_retro_template', 'engagement', 'success');
        toggleRemoveTemplate(null)();
        notifications.success($LL.retroTemplateRemoveSuccess());
        getTemplates();
      })
      .catch(function () {
        notifications.danger($LL.retroTemplateRemoveError());
        eventTag('remove_retro_template', 'engagement', 'failure');
      });
  }

  function toggleCreateTemplate() {
    showAddTemplate = !showAddTemplate;
  }

  const toggleUpdateTemplate = template => () => {
    updateTemplate = template;
    showUpdateTemplate = !showUpdateTemplate;
  };

  const toggleRemoveTemplate = templateId => () => {
    showRemoveTemplate = !showRemoveTemplate;
    removeTemplateId = templateId;
  };

  $: isAdmin = validateUserIsAdmin($user);

  let showFormat = false;
  let selectedTemplate = null;

  let toggleViewFormat = template => {
    showFormat = !showFormat;
    selectedTemplate = template;
  };
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="{$LL.retroTemplates()}"
      createBtnEnabled="{isAdmin || isEntityAdmin}"
      createBtnText="{$LL.retroTemplateCreate()}"
      createButtonHandler="{toggleCreateTemplate}"
      createBtnTestId="template-create"
    />
    <Table>
      <tr slot="header">
        <HeadCol>{$LL.name()}</HeadCol>
        <HeadCol>{$LL.description()}</HeadCol>
        <HeadCol>{$LL.format()}</HeadCol>
        {#if isAdmin && !teamId && !organizationId}
          <HeadCol>{$LL.isPublic()}</HeadCol>
        {/if}
        <HeadCol>{$LL.default()}</HeadCol>
        <HeadCol type="action">
          <span class="sr-only">{$LL.actions()}</span>
        </HeadCol>
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each templates as template, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <div class="font-medium text-gray-900 dark:text-gray-200">
                <span data-testid="template-name">{template.name}</span>
              </div>
            </RowCol>
            <RowCol>
              <span data-testid="template-description"
                >{template.description}</span
              >
            </RowCol>
            <RowCol>
              <span data-testid="template-format">
                <button
                  on:click="{() => toggleViewFormat(template)}"
                  class="text-blue-500 dark:text-sky-400"
                >
                  <Eye class="w-5 h-5" />
                </button>
              </span>
            </RowCol>
            {#if isAdmin && !teamId && !organizationId}
              <RowCol>
                <span data-testid="template-is-public"
                  ><BooleanDisplay boolValue="{template.isPublic}" /></span
                >
              </RowCol>
            {/if}
            <RowCol>
              <span data-testid="template-is-default"
                ><BooleanDisplay boolValue="{template.defaultTemplate}" />
              </span>
            </RowCol>
            <RowCol type="action">
              {#if isAdmin || isEntityAdmin}
                <CrudActions
                  editBtnClickHandler="{toggleUpdateTemplate(template)}"
                  deleteBtnClickHandler="{toggleRemoveTemplate(template.id)}"
                />
              {/if}
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>
    <TableFooter
      bind:current="{templatesPage}"
      num_items="{templateCount}"
      per_page="{templatesPageLimit}"
      on:navigate="{changePage}"
    />
  </TableContainer>

  {#if showAddTemplate}
    <CreateRetroTemplate
      toggleCreate="{toggleCreateTemplate}"
      handleCreate="{handleCreateTemplate}"
      organizationId="{organizationId}"
      departmentId="{departmentId}"
      teamId="{teamId}"
      apiPrefix="{apiPrefix}"
      xfetch="{xfetch}"
      notifications="{notifications}"
      eventTag="{eventTag}"
    />
  {/if}

  {#if showUpdateTemplate}
    <UpdateRetroTemplate
      toggleUpdate="{toggleUpdateTemplate({})}"
      handleUpdate="{handleTemplateUpdate}"
      templateId="{updateTemplate.id}"
      name="{updateTemplate.name}"
      description="{updateTemplate.description}"
      format="{updateTemplate.format}"
      isPublic="{updateTemplate.isPublic}"
      defaultTemplate="{updateTemplate.defaultTemplate}"
      organizationId="{organizationId}"
      departmentId="{departmentId}"
      teamId="{teamId}"
      apiPrefix="{apiPrefix}"
      xfetch="{xfetch}"
      notifications="{notifications}"
      eventTag="{eventTag}"
    />
  {/if}

  {#if showRemoveTemplate}
    <DeleteConfirmation
      toggleDelete="{toggleRemoveTemplate(null)}"
      handleDelete="{handleTemplateRemove}"
      permanent="{true}"
      confirmText="{$LL.removeRetroTemplateConfirmText()}"
      confirmBtnText="{$LL.removeRetroTemplate()}"
    />
  {/if}

  {#if showFormat}
    <ViewFormat
      format="{selectedTemplate.format}"
      toggleClose="{toggleViewFormat}"
    />
  {/if}
</div>

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
//   import CreateProject from './CreateProject.svelte';
//   import UpdateProject from './UpdateProject.svelte';
  import TableFooter from '../table/TableFooter.svelte';
  import { ExternalLink } from 'lucide-svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    organizationId?: any;
    teamId?: any;
    departmentId?: any;
    isEntityAdmin?: boolean;
    projects?: any;
    apiPrefix?: string;
    projectCount?: number;
    projectsPage?: number;
    projectsPageLimit?: number;
    changePage?: any;
    getProjects?: any;
  }

  let {
    xfetch,
    notifications,
    organizationId,
    teamId,
    departmentId,
    isEntityAdmin = false,
    projects = [],
    apiPrefix = '/api',
    projectCount = 0,
    projectsPage = $bindable(1),
    projectsPageLimit = 10,
    changePage = () => {},
    getProjects = () => {}
  }: Props = $props();

  const dispatch = createEventDispatcher();

  let showAddProject = $state(false);
  let showUpdateProject = $state(false);
  let updateProject = $state({});
  let showRemoveProject = $state(false);
  let removeProjectId = null;

  function handleCreateProject() {
    getProjects();
    toggleCreateProject();
  }

  function handleProjectUpdate() {
    getProjects();
    toggleUpdateProject({})();
  }

  function handleProjectRemove() {
    xfetch(`${apiPrefix}/projects/${removeProjectId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleRemoveProject(null)();
        notifications.success($LL.projectRemoveSuccess());
        getProjects();
      })
      .catch(function () {
        notifications.danger($LL.projectRemoveError());
      });
  }

  function toggleCreateProject() {
    showAddProject = !showAddProject;
  }

  const toggleUpdateProject = project => () => {
    updateProject = project;
    showUpdateProject = !showUpdateProject;
  };

  const toggleRemoveProject = projectId => () => {
    showRemoveProject = !showRemoveProject;
    removeProjectId = projectId;
  };

  let isAdmin = $derived(validateUserIsAdmin($user));

  function getProjectScope(project) {
    if (project.teamId) return 'Team';
    if (project.departmentId) return 'Department';
    if (project.organizationId) return 'Organization';
    return 'Global';
  }

  function formatDate(dateString) {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString();
  }
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title={$LL.projects()}
      createBtnEnabled={isAdmin || isEntityAdmin}
      createBtnText={$LL.projectCreate()}
      createButtonHandler={toggleCreateProject}
      createBtnTestId="project-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>{$LL.projectKey()}</HeadCol>
          <HeadCol>{$LL.name()}</HeadCol>
          <HeadCol>{$LL.description()}</HeadCol>
          {#if isAdmin && !teamId && !organizationId && !departmentId}
            <HeadCol>{$LL.scope()}</HeadCol>
          {/if}
          <HeadCol>{$LL.createdAt()}</HeadCol>
          <HeadCol type="action">
            <span class="sr-only">{$LL.actions()}</span>
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class="{className}">
          {#each projects as project, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200">
                  <span data-testid="project-key" class="font-mono text-sm bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">
                    {project.projectKey}
                  </span>
                </div>
              </RowCol>
              <RowCol>
                <div class="font-medium text-gray-900 dark:text-gray-200">
                  <span data-testid="project-name">{project.name}</span>
                </div>
              </RowCol>
              <RowCol>
                <span data-testid="project-description" class="text-gray-600 dark:text-gray-400">
                  {project.description || '-'}
                </span>
              </RowCol>
              {#if isAdmin && !teamId && !organizationId && !departmentId}
                <RowCol>
                  <span data-testid="project-scope" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                    {getProjectScope(project)}
                  </span>
                </RowCol>
              {/if}
              <RowCol>
                <span data-testid="project-created-at" class="text-sm text-gray-500 dark:text-gray-400">
                  {formatDate(project.createdAt)}
                </span>
              </RowCol>
              <RowCol type="action">
                {#if isAdmin || isEntityAdmin}
                  <CrudActions
                    editBtnClickHandler={toggleUpdateProject(project)}
                    deleteBtnClickHandler={toggleRemoveProject(project.id)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
    <TableFooter
      bind:current="{projectsPage}"
      num_items={projectCount}
      per_page={projectsPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>

  <!-- {#if showAddProject}
    <CreateProject
      toggleCreate={toggleCreateProject}
      handleCreate={handleCreateProject}
      organizationId={organizationId}
      departmentId={departmentId}
      teamId={teamId}
      apiPrefix={apiPrefix}
      xfetch={xfetch}
      notifications={notifications}
    />
  {/if}

  {#if showUpdateProject}
    <UpdateProject
      toggleUpdate={toggleUpdateProject({})}
      handleUpdate={handleProjectUpdate}
      projectId={updateProject.id}
      projectKey={updateProject.projectKey}
      name={updateProject.name}
      description={updateProject.description}
      organizationId={updateProject.organizationId || organizationId}
      departmentId={updateProject.departmentId || departmentId}
      teamId={updateProject.teamId || teamId}
      apiPrefix={apiPrefix}
      xfetch={xfetch}
      notifications={notifications}
    />
  {/if}
 -->
  {#if showRemoveProject}
    <DeleteConfirmation
      toggleDelete={toggleRemoveProject(null)}
      handleDelete={handleProjectRemove}
      permanent={true}
      confirmText={$LL.removeProjectConfirmText()}
      confirmBtnText={$LL.removeProject()}
    />
  {/if}
</div>
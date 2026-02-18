<script lang="ts">
  import { onMount } from 'svelte';

  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import { ChevronRight } from '@lucide/svelte';
  import CreateDepartment from '../../components/team/CreateDepartment.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { Department } from '../../types/organization';
  import OrgPageLayout from '../../components/organization/OrgPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const departmentsPageLimit = 1000;
  const orgPrefix = $derived(`/api/organizations/${organizationId}`);

  let organization = $state({
    id: '',
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });

  $effect(() => {
    organization.id = organizationId;
  });
  let role = $state('MEMBER');
  let departments: Department[] = $state([]);
  let showCreateDepartment = $state(false);
  let showDeleteDepartment = $state(false);
  let deleteDeptId = $state<string | null>(null);
  let departmentsPage = $state(1);

  function toggleCreateDepartment() {
    showCreateDepartment = !showCreateDepartment;
  }

  const toggleDeleteDepartment = (deptId: string | null) => () => {
    showDeleteDepartment = !showDeleteDepartment;
    deleteDeptId = deptId;
  };

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getDepartments();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  function getDepartments() {
    const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit;
    xfetch(`${orgPrefix}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        departments = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetDepartmentsError());
      });
  }

  function createDepartmentHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/departments`, { body })
      .then(res => res.json())
      .then(function (result) {
        router.route(`${appRoutes.organization}/${organizationId}/department/${result.data.id}`);
      })
      .catch(function () {
        notifications.danger($LL.departmentCreateError());
      });
  }

  function handleDeleteDepartment() {
    xfetch(`${orgPrefix}/departments/${deleteDeptId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteDepartment(null)();
        notifications.success($LL.departmentDeleteSuccess());
        getDepartments();
      })
      .catch(function () {
        notifications.danger($LL.departmentDeleteError());
      });
  }

  let defaultDepartment: Department = {
    id: '',
    name: '',
    createdDate: '',
    updatedDate: '',
  };
  let selectedDepartment = $state({ ...defaultDepartment });
  let showDepartmentUpdate = $state(false);

  function toggleUpdateDepartment(dept: Department) {
    return () => {
      selectedDepartment = dept;
      showDepartmentUpdate = !showDepartmentUpdate;
    };
  }

  function updateDepartmentHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${organizationId}/departments/${selectedDepartment.id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function (result) {
        getDepartments();
        toggleUpdateDepartment(defaultDepartment)();
        notifications.success(`${$LL.deptUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.deptUpdateError()}`);
      });
  }

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
  });

  let isAdmin = $derived(role === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.departments()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<OrgPageLayout activePage="departments" {organizationId}>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
    <ChevronRight class="w-8 h-8 inline-block" />
    {$LL.departments()}
  </h1>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title={$LL.departments()}
        createBtnEnabled={isAdmin}
        createBtnText={$LL.departmentCreate()}
        createButtonHandler={toggleCreateDepartment}
        createBtnTestId="department-create"
      />
      <Table>
        {#snippet header()}
          <tr>
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
            <HeadCol type="action">
              <span class="sr-only">{$LL.actions()}</span>
            </HeadCol>
          </tr>
        {/snippet}
        {#snippet body({ class: className })}
          <tbody class={className}>
            {#each departments as department, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <a
                    href="{appRoutes.organization}/{organizationId}/department/{department.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                  >
                    {department.name}
                  </a>
                </RowCol>
                <RowCol>
                  {new Date(department.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(department.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  {#if isAdmin}
                    <CrudActions
                      editBtnClickHandler={toggleUpdateDepartment(department)}
                      deleteBtnClickHandler={toggleDeleteDepartment(department.id)}
                    />
                  {/if}
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>
  </div>

  {#if showCreateDepartment}
    <CreateDepartment toggleCreate={toggleCreateDepartment} handleCreate={createDepartmentHandler} />
  {/if}

  {#if showDepartmentUpdate}
    <CreateDepartment
      departmentName={selectedDepartment.name}
      toggleCreate={toggleUpdateDepartment(defaultDepartment)}
      handleCreate={updateDepartmentHandler}
    />
  {/if}

  {#if showDeleteDepartment}
    <DeleteConfirmation
      toggleDelete={toggleDeleteDepartment(null)}
      handleDelete={handleDeleteDepartment}
      confirmText={$LL.deleteDepartmentConfirmText()}
      confirmBtnText={$LL.deleteDepartment()}
    />
  {/if}
</OrgPageLayout>

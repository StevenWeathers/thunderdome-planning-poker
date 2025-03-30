<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import {
    BarChart2,
    CheckSquare,
    ChevronRight,
    LayoutDashboard,
    Network,
    RefreshCcw,
    SquareDashedKanban,
    User,
    Users,
    Vote,
  } from 'lucide-svelte';
  import CreateDepartment from '../../components/team/CreateDepartment.svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import UsersList from '../../components/team/UsersList.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import InvitesList from '../../components/team/InvitesList.svelte';
  import EstimationScalesList from '../../components/estimationscale/EstimationScalesList.svelte';
  import MetricsDisplay from '../../components/global/MetricsDisplay.svelte';
  import {
    fetchAndUpdateMetrics,
    MetricItem,
  } from '../../components/team/metrics';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
  import RetroTemplatesList from '../../components/retrotemplate/RetroTemplatesList.svelte';
  import PokerSettings from '../../components/poker/PokerSettings.svelte';
  import RetroSettings from '../../components/retro/RetroSettings.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
    organizationId: any;
  }

  let {
    xfetch,
    router,
    notifications,
    organizationId
  }: Props = $props();

  const departmentsPageLimit = 1000;
  const teamsPageLimit = 1000;
  const usersPageLimit = 1000;
  const orgPrefix = `/api/organizations/${organizationId}`;

  let invitesList = $state();
  let organization = $state({
    id: organizationId,
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });
  let role = $state('MEMBER');
  let users = $state([]);
  let departments = $state([]);
  let teams = $state([]);
  let invites = [];
  let showCreateDepartment = $state(false);
  let showCreateTeam = $state(false);
  let showDeleteTeam = $state(false);
  let showDeleteDepartment = $state(false);
  let showDeleteOrganization = $state(false);
  let deleteTeamId = null;
  let deleteDeptId = null;
  let teamsPage = 1;
  let departmentsPage = 1;
  let usersPage = 1;

  function toggleCreateDepartment() {
    showCreateDepartment = !showCreateDepartment;
  }

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  const toggleDeleteTeam = teamId => () => {
    showDeleteTeam = !showDeleteTeam;
    deleteTeamId = teamId;
  };

  const toggleDeleteDepartment = deptId => () => {
    showDeleteDepartment = !showDeleteDepartment;
    deleteDeptId = deptId;
  };

  const toggleDeleteOrganization = () => {
    showDeleteOrganization = !showDeleteOrganization;
  };

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getDepartments();
        getTeams();
        getUsers();
        getEstimationScales();
        getRetroTemplates();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  let organizationMetrics: MetricItem[] = $state([
    {
      key: 'department_count',
      name: 'Departments',
      value: 0,
      icon: Network,
    },
    { key: 'team_count', name: 'Team Count', value: 0, icon: Users },
    {
      key: 'team_checkin_count',
      name: 'Team Check-ins',
      value: 0,
      icon: CheckSquare,
    },
    { key: 'user_count', name: 'Users', value: 0, icon: User },
    { key: 'poker_count', name: 'Poker Games', value: 0, icon: Vote },
    { key: 'retro_count', name: 'Retros', value: 0, icon: RefreshCcw },
    {
      key: 'storyboard_count',
      name: 'Storyboards',
      value: 0,
      icon: LayoutDashboard,
    },
    {
      key: 'estimation_scale_count',
      name: 'Custom Estimation Scales',
      value: 0,
      icon: BarChart2,
    },
    {
      key: 'retro_template_count',
      name: 'Retro Templates',
      value: 0,
      icon: SquareDashedKanban,
    },
  ]);

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${orgPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function getDepartments() {
    const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit;
    xfetch(
      `${orgPrefix}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        departments = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetDepartmentsError());
      });
  }

  function getTeams() {
    const teamsOffset = (teamsPage - 1) * teamsPageLimit;
    xfetch(`${orgPrefix}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetTeamsError());
      });
  }

  function createDepartmentHandler(name) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/departments`, { body })
      .then(res => res.json())
      .then(function (result) {
        router.route(
          `${appRoutes.organization}/${organizationId}/department/${result.data.id}`,
        );
      })
      .catch(function () {
        notifications.danger($LL.departmentCreateError());
      });
  }

  function createTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/teams`, { body })
      .then(res => res.json())
      .then(function () {
        toggleCreateTeam();
        notifications.success($LL.teamCreateSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamCreateError());
      });
  }

  function handleDeleteTeam() {
    xfetch(`${orgPrefix}/teams/${deleteTeamId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteTeam(null)();
        notifications.success($LL.teamDeleteSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
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

  function handleDeleteOrganization() {
    xfetch(`${orgPrefix}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteTeam();
        notifications.success($LL.organizationDeleteSuccess());
        router.route(appRoutes.teams);
      })
      .catch(function () {
        notifications.danger($LL.organizationDeleteError());
      });
  }

  let defaultDepartment = {
    id: '',
    name: '',
  };
  let selectedDepartment = $state({ ...defaultDepartment });
  let showDepartmentUpdate = $state(false);

  function toggleUpdateDepartment(dept) {
    return () => {
      selectedDepartment = dept;
      showDepartmentUpdate = !showDepartmentUpdate;
    };
  }

  let defaultTeam = {
    id: '',
    name: '',
  };
  let selectedTeam = $state({ ...defaultTeam });
  let showTeamUpdate = $state(false);

  function toggleUpdateTeam(team) {
    return () => {
      selectedTeam = team;
      showTeamUpdate = !showTeamUpdate;
    };
  }

  function updateDepartmentHandler(name) {
    const body = {
      name,
    };

    xfetch(
      `/api/organizations/${organizationId}/departments/${selectedDepartment.id}`,
      { body, method: 'PUT' },
    )
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

  function updateTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${organizationId}/teams/${selectedTeam.id}`, {
      body,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function () {
        toggleUpdateTeam(defaultTeam)();
        getTeams();
        notifications.success(`${$LL.teamUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.teamUpdateError()}`);
      });
  }

  const scalesPageLimit = 20;
  let estimationScales = $state([]);
  let scaleCount = 0;
  let scalesPage = $state(1);

  const changeScalesPage = evt => {
    scalesPage = evt.detail;
    getEstimationScales();
  };

  function getEstimationScales() {
    const scalesOffset = (scalesPage - 1) * scalesPageLimit;
    if (
      AppConfig.FeaturePoker &&
      (!AppConfig.SubscriptionsEnabled ||
        (AppConfig.SubscriptionsEnabled && organization.subscribed))
    ) {
      xfetch(
        `${orgPrefix}/estimation-scales?limit=${scalesPageLimit}&offset=${scalesOffset}`,
      )
        .then(res => res.json())
        .then(function (result) {
          estimationScales = result.data;
        })
        .catch(function () {
          notifications.danger('Failed to get estimation scales');
        });
    }
  }

  const retroTemplatePageLimit = 20;
  let retroTemplates = $state([]);
  let retroTemplateCount = 0;
  let retroTemplatesPage = $state(1);

  const changeRetroTemplatesPage = evt => {
    retroTemplatesPage = evt.detail;
    getRetroTemplates();
  };

  function getRetroTemplates() {
    const offset = (retroTemplatesPage - 1) * retroTemplatePageLimit;
    if (
      AppConfig.FeatureRetro &&
      (!AppConfig.SubscriptionsEnabled ||
        (AppConfig.SubscriptionsEnabled && organization.subscribed))
    ) {
      xfetch(
        `${orgPrefix}/retro-templates?limit=${retroTemplatePageLimit}&offset=${offset}`,
      )
        .then(res => res.json())
        .then(function (result) {
          retroTemplates = result.data;
        })
        .catch(function () {
          notifications.danger('Failed to get retro templates');
        });
    }
  }

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
    try {
      organizationMetrics = await fetchAndUpdateMetrics(
        orgPrefix,
        organizationMetrics,
      );
    } catch (e) {
      notifications.danger('Failed to get organization metrics');
    }
  });

  let isAdmin = $derived(role === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.organization()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
  </h1>

  <div class="mb-8">
    <MetricsDisplay metrics="{organizationMetrics}" />
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title="{$LL.departments()}"
        createBtnEnabled="{isAdmin}"
        createBtnText="{$LL.departmentCreate()}"
        createButtonHandler="{toggleCreateDepartment}"
        createBtnTestId="department-create"
      />
      <Table>
        {#snippet header()}
                <tr >
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
                <tbody   class="{className}">
            {#each departments as department, i}
              <TableRow itemIndex="{i}">
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
                      editBtnClickHandler="{toggleUpdateDepartment(department)}"
                      deleteBtnClickHandler="{toggleDeleteDepartment(
                        department.id,
                      )}"
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

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title="{$LL.teams()}"
        createBtnEnabled="{isAdmin}"
        createBtnText="{$LL.teamCreate()}"
        createButtonHandler="{toggleCreateTeam}"
        createBtnTestId="team-create"
      />
      <Table>
        {#snippet header()}
                <tr >
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
                <tbody   class="{className}">
            {#each teams as team, i}
              <TableRow itemIndex="{i}">
                <RowCol>
                  <a
                    href="{appRoutes.organization}/{organizationId}/team/{team.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                  >
                    {team.name}
                  </a>
                </RowCol>
                <RowCol>
                  {new Date(team.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(team.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  {#if isAdmin}
                    <CrudActions
                      editBtnClickHandler="{toggleUpdateTeam(team)}"
                      deleteBtnClickHandler="{toggleDeleteTeam(team.id)}"
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

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <InvitesList
        xfetch="{xfetch}"
        notifications="{notifications}"
        pageType="organization"
        teamPrefix="{orgPrefix}"
        bind:this="{invitesList}"
      />
    </div>
  {/if}

  <UsersList
    users="{users}"
    getUsers="{getUsers}"
    xfetch="{xfetch}"
    notifications="{notifications}"
    isAdmin="{isAdmin}"
    pageType="organization"
    orgId="{organizationId}"
    teamPrefix="/api/organizations/{organizationId}"
    on:user-invited="{() => {
      invitesList.f('user-invited');
    }}"
  />

  {#if AppConfig.FeaturePoker}
    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <PokerSettings
          xfetch="{xfetch}"
          notifications="{notifications}"
          isEntityAdmin="{isAdmin}"
          apiPrefix="{orgPrefix}"
          organizationId="{organizationId}"
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Optimize your Organization's estimation workflow with customized default Planning Poker settings."
        />
      {/if}
    </div>

    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <EstimationScalesList
          xfetch="{xfetch}"
          notifications="{notifications}"
          isEntityAdmin="{isAdmin}"
          apiPrefix="{orgPrefix}"
          organizationId="{organizationId}"
          scales="{estimationScales}"
          getScales="{getEstimationScales}"
          scaleCount="{scaleCount}"
          scalesPage="{scalesPage}"
          scalesPageLimit="{scalesPageLimit}"
          changePage="{changeScalesPage}"
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Create custom poker point scales to match your Organization's estimation style."
        />
      {/if}
    </div>
  {/if}

  {#if AppConfig.FeatureRetro}
    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <RetroSettings
          xfetch="{xfetch}"
          notifications="{notifications}"
          isEntityAdmin="{isAdmin}"
          apiPrefix="{orgPrefix}"
          organizationId="{organizationId}"
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Enhance your Organization's reflection process with customized default Retrospective settings."
        />
      {/if}
    </div>

    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <RetroTemplatesList
          xfetch="{xfetch}"
          notifications="{notifications}"
          isEntityAdmin="{isAdmin}"
          apiPrefix="{orgPrefix}"
          organizationId="{organizationId}"
          templates="{retroTemplates}"
          getTemplates="{getRetroTemplates}"
          templateCount="{retroTemplateCount}"
          templatesPage="{retroTemplatesPage}"
          templatesPageLimit="{retroTemplatePageLimit}"
          changePage="{changeRetroTemplatesPage}"
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Tailor your Organization's reflection process with custom retrospective templates."
        />
      {/if}
    </div>
  {/if}

  {#if isAdmin}
    <div class="w-full text-center mt-8">
      <HollowButton onClick="{toggleDeleteOrganization}" color="red">
        {$LL.deleteOrganization()}
      </HollowButton>
    </div>
  {/if}

  {#if showCreateDepartment}
    <CreateDepartment
      toggleCreate="{toggleCreateDepartment}"
      handleCreate="{createDepartmentHandler}"
    />
  {/if}

  {#if showDepartmentUpdate}
    <CreateDepartment
      departmentName="{selectedDepartment.name}"
      toggleCreate="{toggleUpdateDepartment(defaultDepartment)}"
      handleCreate="{updateDepartmentHandler}"
    />
  {/if}

  {#if showCreateTeam}
    <CreateTeam
      toggleCreate="{toggleCreateTeam}"
      handleCreate="{createTeamHandler}"
    />
  {/if}

  {#if showTeamUpdate}
    <CreateTeam
      teamName="{selectedTeam.name}"
      toggleCreate="{toggleUpdateTeam(defaultTeam)}"
      handleCreate="{updateTeamHandler}"
    />
  {/if}

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteTeam(null)}"
      handleDelete="{handleDeleteTeam}"
      confirmText="{$LL.deleteTeamConfirmText()}"
      confirmBtnText="{$LL.deleteTeam()}"
    />
  {/if}

  {#if showDeleteDepartment}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteDepartment(null)}"
      handleDelete="{handleDeleteDepartment}"
      confirmText="{$LL.deleteDepartmentConfirmText()}"
      confirmBtnText="{$LL.deleteDepartment()}"
    />
  {/if}

  {#if showDeleteOrganization}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteOrganization}"
      handleDelete="{handleDeleteOrganization}"
      confirmText="{$LL.deleteOrganizationConfirmText()}"
      confirmBtnText="{$LL.deleteOrganization()}"
    />
  {/if}
</PageLayout>

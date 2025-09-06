<script lang="ts">
  import { onMount } from 'svelte';

  import { ChevronRight } from 'lucide-svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import ProjectsList from '../../components/project/ProjectsList.svelte';
  import TeamPageLayout from '../../components/team/TeamPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
    departmentId: any;
    teamId: any;
  }

  let {
    xfetch,
    router,
    notifications,
    organizationId,
    departmentId,
    teamId
  }: Props = $props();

  const projectsPageLimit = 10;

  let team = $state({
    id: teamId,
    name: '',
    subscribed: false,
  });
  let organization = $state({
    id: organizationId,
    name: '',
    subscribed: false,
  });
  let department = $state({
    id: departmentId,
    name: '',
  });

  let projectCount = $state(0);
  let projectsPage = $state(1);
  let projects = $state([]);

  let organizationRole = $state('');
  let departmentRole = $state('');
  let teamRole = $state('');
  let isEntityAdmin = $state(false);
  let isTeamMember = $state(false);

  const apiPrefix = '/api';
  let orgPrefix = $derived(departmentId
    ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
    : `${apiPrefix}/organizations/${organizationId}`);
  let teamPrefix = $derived(organizationId
    ? `${orgPrefix}/teams/${teamId}`
    : `${apiPrefix}/teams/${teamId}`);

  function getTeam() {
    xfetch(teamPrefix)
      .then(res => res.json())
      .then(function (result) {
        team = result.data.team;
        teamRole = result.data.teamRole;

        if (departmentId) {
          department = result.data.department;
          departmentRole = result.data.departmentRole;
        }
        if (organizationId) {
          organization = result.data.organization;
          organizationRole = result.data.organizationRole;
        }

        isEntityAdmin =
          organizationRole === 'ADMIN' ||
          departmentRole === 'ADMIN' ||
          teamRole === 'ADMIN';
        isTeamMember =
          organizationRole === 'ADMIN' ||
          departmentRole === 'ADMIN' ||
          teamRole !== '';

        getProjects();
      })
      .catch(function () {
        notifications.danger($LL.teamGetError());
      });
  }

  function getProjects() {
    const projectsOffset = (projectsPage - 1) * projectsPageLimit;
    xfetch(
      `${teamPrefix}/projects?limit=${projectsPageLimit}&offset=${projectsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        projects = result.data;
        projectCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger('Failed to fetch projects');
      });
  }

  const changeProjectsPage = evt => {
    projectsPage = evt.detail;
    getProjects();
  };

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getTeam();
  });
</script>

<svelte:head>
  <title>{$LL.projects()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<TeamPageLayout activePage="projects" {teamId}>
  <div class="flex mb-6 lg:mb-8">
    <div class="flex-1">
      <h1 class="text-3xl font-semibold font-rajdhani dark:text-white">
        <span class="uppercase">{$LL.team()}</span>
        <ChevronRight class="w-8 h-8 inline-block" />
        {team.name} <ChevronRight class="w-8 h-8 inline-block" /> {$LL.projects()}
      </h1>

      {#if organizationId}
        <div class="text-xl font-semibold font-rajdhani dark:text-white">
          <span class="uppercase">{$LL.organization()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.organization}/{organization.id}"
          >
            {organization.name}
          </a>
          {#if departmentId}
            &nbsp;
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.department()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}"
            >
              {department.name}
            </a>
          {/if}
        </div>
      {/if}
    </div>
  </div>

  {#if AppConfig.FeatureProject}
    <div class="mt-8">
        {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <ProjectsList
          {xfetch}
          {notifications}
          {projects}
          apiPrefix={teamPrefix}
          getProjects={getProjects}
          changePage={changeProjectsPage}
          {projectCount}
          {projectsPage}
          projectsPageLimit={projectsPageLimit}
          {organizationId}
          {departmentId}
          {teamId}
          {isEntityAdmin}
        />
      {:else}
        <FeatureSubscribeBanner
          isNew={true}
          salesPitch="Streamline your team's workflow with organized Projects that keep everything connected."
        />
      {/if}
    </div>
  {/if}
</TeamPageLayout>

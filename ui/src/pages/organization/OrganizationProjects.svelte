<script lang="ts">
  import { onMount } from 'svelte';

  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { ChevronRight } from 'lucide-svelte';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import ProjectsList from '../../components/project/ProjectsList.svelte';
  import OrgPageLayout from '../../components/organization/OrgPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const projectsPageLimit = 10;
  const orgPrefix = `/api/organizations/${organizationId}`;

  let organization = $state({
    id: organizationId,
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });
  let projectCount = $state(0);
  let projectsPage = $state(1);
  let projects = $state([]);
  let role = $state('MEMBER');
  let isEntityAdmin = $state(false);

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        if (role === 'ADMIN') {
          isEntityAdmin = true;
        }

        getProjects();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  function getProjects() {
    const projectsOffset = (projectsPage - 1) * projectsPageLimit;
    xfetch(`${orgPrefix}/projects?limit=${projectsPageLimit}&offset=${projectsOffset}`)
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

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
  });
</script>

<svelte:head>
  <title>{$LL.projects()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<OrgPageLayout activePage="projects" {organizationId}>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
    <ChevronRight class="w-8 h-8 inline-block" />
    {$LL.projects()}
  </h1>

  {#if AppConfig.FeatureProject}
    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <ProjectsList
          {xfetch}
          {notifications}
          {projects}
          apiPrefix={orgPrefix}
          {getProjects}
          changePage={changeProjectsPage}
          {projectCount}
          {projectsPage}
          {projectsPageLimit}
          {organizationId}
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
</OrgPageLayout>

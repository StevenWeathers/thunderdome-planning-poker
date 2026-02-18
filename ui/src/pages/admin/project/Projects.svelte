<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import ProjectsList from '../../../components/project/ProjectsList.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const projectsPageLimit = 100;
  let projectCount = $state(0);
  let projectsPage = $state(1);
  let projects = $state([]);

  function getProjects() {
    const projectsOffset = (projectsPage - 1) * projectsPageLimit;
    xfetch(`/api/admin/projects?limit=${projectsPageLimit}&offset=${projectsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        projects = result.data;
        projectCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger('Failed to fetch projects');
      });
  }

  const changePage = (evt: CustomEvent) => {
    projectsPage = evt.detail;
    getProjects();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getProjects();
  });
</script>

<svelte:head>
  <title
    >{$LL.projects()}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="projects">
  <ProjectsList
    {xfetch}
    {notifications}
    {projects}
    apiPrefix="/api/admin"
    {getProjects}
    {changePage}
    {projectCount}
    {projectsPage}
    {projectsPageLimit}
    isAdminPage={true}
  />
</AdminPageLayout>

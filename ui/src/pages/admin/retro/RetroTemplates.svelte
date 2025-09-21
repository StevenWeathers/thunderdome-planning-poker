<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import RetroTemplatesList from '../../../components/retrotemplate/RetroTemplatesList.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const templatesPageLimit = 100;
  let templateCount = $state(0);
  let templatesPage = $state(1);
  let templates = $state([]);

  function getTemplates() {
    const templatesOffset = (templatesPage - 1) * templatesPageLimit;
    xfetch(`/api/admin/retro-templates?limit=${templatesPageLimit}&offset=${templatesOffset}`)
      .then(res => res.json())
      .then(function (result) {
        templates = result.data;
        templateCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger('Failed to fetch templates');
      });
  }

  const changePage = evt => {
    templatesPage = evt.detail;
    getTemplates();
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

    getTemplates();
  });
</script>

<svelte:head>
  <title
    >{$LL.retroTemplates()}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="retro-templates">
  <RetroTemplatesList
    {xfetch}
    {notifications}
    {templates}
    apiPrefix="/api/admin"
    {getTemplates}
    {changePage}
    {templateCount}
    {templatesPage}
    {templatesPageLimit}
  />
</AdminPageLayout>

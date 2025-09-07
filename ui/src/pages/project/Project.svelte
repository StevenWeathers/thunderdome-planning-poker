<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '../../stores';
    import { validateUserIsRegistered } from '../../validationUtils';
    import { appRoutes, AppConfig } from '../../config';
    import ProjectPageLayout from '../../components/project/ProjectPageLayout.svelte';

    import type { ApiClient } from '../../types/apiclient';
    import type { NotificationService } from '../../types/notifications';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';

    interface Props {
        xfetch: ApiClient;
        router: any;
        notifications: NotificationService;
        projectId: string;
    }

    let {
        xfetch,
        router,
        notifications,
        projectId,
    }: Props = $props();

    // Basic project state
    let project = $state({
        id: projectId,
        projectKey: '',
        name: '',
        description: ''
    });

    const apiPrefix = '/api';
    // Determine correct scoped prefix
    let projectPrefix = `${apiPrefix}/projects/${projectId}`;

    function getProject() {
        xfetch(projectPrefix)
            .then(res => res.json())
            .then(result => {
                const data = result?.data || result;
                if (data) {
                    project = { ...project, ...data };
                }
            })
            .catch(() => {
            notifications.danger('Failed to get project');
            });
    }

    onMount(() => {
        if (!$user.id || !validateUserIsRegistered($user)) {
            router.route(appRoutes.login);
            return;
        }
        getProject();
    });
</script>

<svelte:head>
    <title>Project {project.name}</title>
</svelte:head>

<ProjectPageLayout activePage="project" {projectId}>
<div class="container mx-auto px-4 py-4 md:py-6 lg:py-8">
    <h1 class="text-3xl font-semibold font-rajdhani dark:text-white" data-testid="project-title">{project.name}</h1>

    <!-- Subscription Required if enabled -->
    {#if AppConfig.SubscriptionsEnabled && $user && !$user.isSubscriber}
        <FeatureSubscribeBanner
          salesPitch="Active subscription required to access this project."
        />
    {/if}
</div>
</ProjectPageLayout>
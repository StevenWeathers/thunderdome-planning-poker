<script lang="ts">
  import SidenavPageLayout, { type PageItem } from '../SidenavPageLayout.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { Users, User, Package } from '@lucide/svelte';

  interface Props {
    activePage?: string;
    children?: import('svelte').Snippet;
    projectId: string;
  }

  let { activePage = 'Team', children, projectId }: Props = $props();

  const { FeatureProject } = AppConfig;

  // Project pages configuration
  let pages: PageItem[] = $derived(
    $LL
      ? [
          {
            name: 'Project',
            label: $LL.project(),
            path: `${appRoutes.projects}/${projectId}`,
            icon: Package,
            enabled: true,
          },
          {
            name: 'Users',
            label: $LL.users(),
            path: `${appRoutes.projects}/${projectId}/users`,
            icon: User,
            enabled: true,
          },
        ]
      : [],
  );
</script>

<SidenavPageLayout {pages} {activePage} menuType="project" expanded={true}>
  {@render children?.()}
</SidenavPageLayout>

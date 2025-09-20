<script lang="ts">
  import SidenavPageLayout, { type PageItem } from '../SidenavPageLayout.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import {
    Users,
    User,
    Package,
  } from 'lucide-svelte';

  interface Props {
    activePage?: string;
    children?: import('svelte').Snippet;
    teamId: string;
    organizationId?: string;
    departmentId?: string;
  }

  let { activePage = 'Team', children, teamId, organizationId, departmentId }: Props = $props();

  const {
    FeatureProject,
  } = AppConfig;

  let routePrefix = departmentId
    ? `/organization/${organizationId}/department/${departmentId}/team/${teamId}`
    : organizationId
      ? `/organization/${organizationId}/team/${teamId}`
      : `/${appRoutes.team}/${teamId}`;

   // Team pages configuration
  let pages: PageItem[] = $derived($LL ? [
    {
      name: 'Team',
      label: $LL.team(),
      path: routePrefix,
      icon: Users,
      enabled: true,
    },
    {
      name: 'Users',
      label: $LL.users(),
      path: `${routePrefix}/users`,
      icon: User,
      enabled: true,
    },
    {
      name: 'Projects',
      label: $LL.projects(),
      path: `${routePrefix}/projects`,
      icon: Package,
      enabled: FeatureProject,
    },
  ] : []);
</script>

<SidenavPageLayout {pages} activePage={activePage} menuType="team" expanded={true}>
  {@render children?.()}
</SidenavPageLayout>
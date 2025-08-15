<script lang="ts">
  import SidenavPageLayout, { type PageItem } from '../SidenavPageLayout.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import {
    Building,
    User,
    Users,
    Package,
    Network,
  } from 'lucide-svelte';

  interface Props {
    activePage?: string;
    children?: import('svelte').Snippet;
    organizationId: string;
  }

  let { activePage = 'Organization', children, organizationId }: Props = $props();

  const {
    FeatureProject,
  } = AppConfig;

   // Organization pages configuration
  let organizationPages: PageItem[] = $derived($LL ? [
    {
      name: 'Organization',
      label: $LL.organization(),
      path: `${appRoutes.organization}/${organizationId}`,
      icon: Building,
      enabled: true,
    },
    {
      name: 'Departments',
      label: $LL.departments(),
      path: `${appRoutes.organization}/${organizationId}/departments`,
      icon: Network,
      enabled: true,
    },
    {
      name: 'Teams',
      label: $LL.teams(),
      path: `${appRoutes.organization}/${organizationId}/teams`,
      icon: Users,
      enabled: true,
    },
    {
      name: 'Users',
      label: $LL.users(),
      path: `${appRoutes.organization}/${organizationId}/users`,
      icon: User,
      enabled: true,
    },
    {
      name: 'Projects',
      label: $LL.projects(),
      path: `${appRoutes.organization}/${organizationId}/projects`,
      icon: Package,
      enabled: FeatureProject,
    },
  ] : []);
</script>

<SidenavPageLayout pages={organizationPages} activePage={activePage} menuType="organization" expanded={true}>
  {@render children?.()}
</SidenavPageLayout>
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
  }

  let { activePage = 'Team', children, teamId }: Props = $props();

  const {
    FeatureProject,
  } = AppConfig;

   // Team pages configuration
  let pages: PageItem[] = $derived($LL ? [
    {
      name: 'Team',
      label: $LL.team(),
      path: `${appRoutes.team}/${teamId}`,
      icon: Users,
      enabled: true,
    },
    {
      name: 'Users',
      label: $LL.users(),
      path: `${appRoutes.team}/${teamId}/users`,
      icon: User,
      enabled: true,
    },
    {
      name: 'Projects',
      label: $LL.projects(),
      path: `${appRoutes.team}/${teamId}/projects`,
      icon: Package,
      enabled: FeatureProject,
    },
  ] : []);
</script>

<SidenavPageLayout {pages} activePage={activePage} menuType="team" expanded={true}>
  {@render children?.()}
</SidenavPageLayout>
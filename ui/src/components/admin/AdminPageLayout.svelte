<script lang="ts">
  import SidenavPageLayout, { type PageItem } from '../SidenavPageLayout.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import {
    ChartNoAxesColumn,
    Bell,
    Building,
    CreditCard,
    Gamepad,
    Key,
    LayoutDashboard,
    RefreshCcw,
    Settings,
    SquareDashedKanban,
    User,
    Users,
    Package,
    InboxIcon,
  } from 'lucide-svelte';

  interface Props {
    activePage?: string;
    children?: import('svelte').Snippet;
  }

  let { activePage = 'admin', children }: Props = $props();

  const {
    ExternalAPIEnabled,
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    FeatureProject,
    OrganizationsEnabled,
    SubscriptionsEnabled,
  } = AppConfig;

  let adminPages: PageItem[] = $derived($LL ? [
    {
      name: 'Admin',
      label: $LL.adminPageAdmin(),
      path: appRoutes.admin,
      icon: Settings,
      enabled: true,
    },
    {
      name: 'Support Tickets',
      icon: InboxIcon,
      label: $LL.supportTickets(),
      path: appRoutes.adminSupportTickets,
      enabled: true,
    },
    {
      name: 'Alerts',
      icon: Bell,
      label: $LL.adminPageAlerts(),
      path: appRoutes.adminAlerts,
      enabled: true,
    },
    {
      name: 'Projects',
      icon: Package,
      label: $LL.projects(),
      path: appRoutes.adminProjects,
      enabled: FeatureProject,
    },
    {
      name: 'Battles',
      label: $LL.battles(),
      path: appRoutes.adminPokerGames,
      icon: Gamepad,
      enabled: FeaturePoker,
    },
    {
      name: 'Retros',
      label: $LL.retros(),
      path: appRoutes.adminRetros,
      icon: RefreshCcw,
      enabled: FeatureRetro,
    },
    {
      name: 'Storyboards',
      label: $LL.storyboards(),
      path: appRoutes.adminStoryboards,
      icon: LayoutDashboard,
      enabled: FeatureStoryboard,
    },
    {
      name: 'Organizations',
      label: $LL.adminPageOrganizations(),
      path: appRoutes.adminOrganizations,
      icon: Building,
      enabled: OrganizationsEnabled,
    },
    {
      name: 'Teams',
      label: $LL.adminPageTeams(),
      path: appRoutes.adminTeams,
      icon: Users,
      enabled: true,
    },
    {
      name: 'Users',
      label: $LL.adminPageUsers(),
      path: appRoutes.adminUsers,
      icon: User,
      enabled: true,
    },
    {
      name: 'API Keys',
      label: $LL.adminPageApi(),
      path: appRoutes.adminApiKeys,
      icon: Key,
      enabled: ExternalAPIEnabled,
    },
    {
      name: 'Subscriptions',
      label: $LL.adminPageSubscriptions(),
      path: appRoutes.adminSubscriptions,
      icon: CreditCard,
      enabled: SubscriptionsEnabled,
    },
    {
      name: 'Estimation Scales',
      label: $LL.estimationScales(),
      path: appRoutes.adminEstimationScales,
      icon: ChartNoAxesColumn,
      enabled: FeaturePoker,
    },
    {
      name: 'Retro Templates',
      label: $LL.retroTemplates(),
      path: appRoutes.adminRetroTemplates,
      icon: SquareDashedKanban,
      enabled: FeatureRetro,
    },
  ] : []);
</script>

<SidenavPageLayout pages={adminPages} activePage={activePage} menuType="admin">
  {@render children?.()}
</SidenavPageLayout>
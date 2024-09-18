<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import SideNavigation from '../global/SideNavigation.svelte';
  import {
    BarChart2,
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
  } from 'lucide-svelte';

  export let activePage = 'admin';

  const {
    ExternalAPIEnabled,
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    OrganizationsEnabled,
    SubscriptionsEnabled,
  } = AppConfig;

  $: pages = $LL && [
    {
      name: 'Admin',
      label: $LL.adminPageAdmin(),
      path: appRoutes.admin,
      icon: Settings,
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
      icon: BarChart2,
      enabled: FeaturePoker,
    },
    {
      name: 'Retro Templates',
      label: $LL.retroTemplates(),
      path: appRoutes.adminRetroTemplates,
      icon: SquareDashedKanban,
      enabled: FeatureRetro,
    },
  ];
</script>

<section class="flex min-h-screen">
  <SideNavigation
    menuItems="{pages}"
    activePage="{activePage}"
    menuType="admin"
  />
  <div class="flex-1 px-4 py-4 md:py-6 md:px-6 lg:py-8 lg:px-8">
    <slot />
  </div>
</section>

<script lang="ts">
  import { onMount } from 'svelte';

  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { ChevronRight } from 'lucide-svelte';
  import UsersList from '../../components/team/UsersList.svelte';
  import InvitesList from '../../components/team/InvitesList.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import OrgPageLayout from '../../components/organization/OrgPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const usersPageLimit = 1000;
  const orgPrefix = $derived(`/api/organizations/${organizationId}`);

  let invitesList = $state();
  let organization = $state({
    id: '',
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });

  $effect(() => {
    organization.id = organizationId;
  });
  let role = $state('MEMBER');
  let users = $state([]);
  let usersPage = $state(1);

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${orgPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
  });

  let isAdmin = $derived(role === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.users()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<OrgPageLayout activePage="users" {organizationId}>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
    <ChevronRight class="w-8 h-8 inline-block" />
    {$LL.users()}
  </h1>

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <InvitesList {xfetch} {notifications} pageType="organization" teamPrefix={orgPrefix} bind:this={invitesList} />
    </div>
  {/if}

  <UsersList
    {users}
    {getUsers}
    {xfetch}
    {notifications}
    {isAdmin}
    pageType="organization"
    orgId={organizationId}
    teamPrefix="/api/organizations/{organizationId}"
    on:user-invited={() => {
      invitesList.f('user-invited');
    }}
  />
</OrgPageLayout>

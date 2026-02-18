<script lang="ts">
  import { onMount } from 'svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import { Check, ChevronRight, MessageSquareMore } from '@lucide/svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import Table from '../../../components/table/Table.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import Toggle from '../../../components/forms/Toggle.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';
  import type { Team } from '../../../types/team';
  import type { PokerGame } from '../../../types/poker';
  import type { Retro, RetroAction } from '../../../types/retro';
  import type { Storyboard } from '../../../types/storyboard';
  import type { TeamUser } from '../../../types/team';

  const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
    departmentId: any;
    teamId: any;
  }

  let { xfetch, router, notifications, organizationId, departmentId, teamId }: Props = $props();

  const battlesPageLimit = 1000;
  const retrosPageLimit = 1000;
  const retroActionsPageLimit = 5;
  const storyboardsPageLimit = 1000;
  const usersPageLimit = 1000;

  let team = $state<Team>({
    id: '',
    name: '',
    createdDate: '',
    updatedDate: '',
  });
  let organization = $state({
    id: '',
    name: '',
  });
  let department = $state({
    id: '',
    name: '',
  });
  let users = $state<TeamUser[]>([]);
  let battles = $state<PokerGame[]>([]);
  let retros = $state<Retro[]>([]);
  let retroActions = $state<RetroAction[]>([]);
  let storyboards = $state<Storyboard[]>([]);

  $effect(() => {
    team.id = teamId;
    organization.id = organizationId;
    department.id = departmentId;
  });
  let usersPage = 1;
  let battlesPage = 1;
  let retrosPage = 1;
  let retroActionsPage = $state(1);
  let storyboardsPage = 1;
  let totalRetroActions = $state(0);
  let completedActionItems = $state(false);

  let showDeleteTeam = $state(false);

  const apiPrefix = '/api';
  let orgPrefix = $derived(
    departmentId
      ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
      : `${apiPrefix}/organizations/${organizationId}`,
  );
  let teamPrefix = $derived(organizationId ? `${orgPrefix}/teams/${teamId}` : `${apiPrefix}/teams/${teamId}`);

  let showRetroActionComments = false;
  let selectedRetroAction = null;
  const toggleRetroActionComments = (id: string) => () => {
    showRetroActionComments = !showRetroActionComments;
    selectedRetroAction = id;
  };

  function getTeam() {
    xfetch(teamPrefix)
      .then(res => res.json())
      .then(function (result) {
        team = result.data.team;

        if (departmentId) {
          department = result.data.department;
        }
        if (organizationId) {
          organization = result.data.organization;
        }

        getBattles();
        getRetros();
        getRetrosActions();
        getStoryboards();
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.teamGetError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${teamPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function getBattles() {
    if (FeaturePoker) {
      const battlesOffset = (battlesPage - 1) * battlesPageLimit;
      xfetch(`${teamPrefix}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`)
        .then(res => res.json())
        .then(function (result) {
          battles = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetBattlesError());
        });
    }
  }

  function getRetros() {
    if (FeatureRetro) {
      const retrosOffset = (retrosPage - 1) * retrosPageLimit;
      xfetch(`${teamPrefix}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`)
        .then(res => res.json())
        .then(function (result) {
          retros = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetRetrosError());
        });
    }
  }

  function getRetrosActions() {
    if (FeatureRetro) {
      const offset = (retroActionsPage - 1) * retroActionsPageLimit;
      xfetch(
        `${teamPrefix}/retro-actions?limit=${retroActionsPageLimit}&offset=${offset}&completed=${completedActionItems}`,
      )
        .then(res => res.json())
        .then(function (result) {
          retroActions = result.data;
          totalRetroActions = result.meta.count;
        })
        .catch(function () {
          notifications.danger($LL.teamGetRetroActionsError());
        });
    }
  }

  function getStoryboards() {
    if (FeatureStoryboard) {
      const storyboardsOffset = (storyboardsPage - 1) * storyboardsPageLimit;
      xfetch(`${teamPrefix}/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}`)
        .then(res => res.json())
        .then(function (result) {
          storyboards = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetStoryboardsError());
        });
    }
  }

  function deleteTeam() {
    xfetch(`${teamPrefix}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        if (departmentId) {
          router.route(`${appRoutes.adminOrganizations}/${organizationId}/department/${departmentId}`);
        } else if (organizationId) {
          router.route(`${appRoutes.adminOrganizations}/${organizationId}`);
        } else {
          router.route(appRoutes.adminTeams);
        }
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
      });
  }

  const changeRetroActionPage = (evt: CustomEvent) => {
    retroActionsPage = evt.detail;
    getRetrosActions();
  };

  const changeRetroActionCompletedToggle = () => {
    retroActionsPage = 1;
    getRetrosActions();
  };

  function toggleDeleteTeam() {
    showDeleteTeam = !showDeleteTeam;
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getTeam();
  });
</script>

<svelte:head>
  <title>{$LL.team()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="teams">
  <div class="px-2 mb-4">
    <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white">
      {team.name}
    </h1>

    {#if organizationId}
      <div class="text-xl font-semibold font-rajdhani dark:text-white">
        <span class="uppercase">{$LL.organization()}</span>
        <ChevronRight class="inline-block" />
        <a
          class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
          href="{appRoutes.adminOrganizations}/{organization.id}"
        >
          {organization.name}
        </a>
        {#if departmentId}
          &nbsp;
          <ChevronRight class="inline-block" />
          <span class="uppercase">{$LL.department()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.adminOrganizations}/{organization.id}/department/{department.id}"
          >
            {department.name}
          </a>
        {/if}
      </div>
    {/if}
  </div>

  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title={team.name} createBtnEnabled={false} />
      <Table>
        {#snippet header()}
          <tr>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
          </tr>
        {/snippet}
        {#snippet body({ class: className })}
          <tbody class={className}>
            <TableRow itemIndex={0}>
              <RowCol>
                {new Date(team.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(team.updatedDate).toLocaleString()}
              </RowCol>
            </TableRow>
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>
  </div>

  <div>
    {#if FeaturePoker}
      <div class="w-full mb-6 lg:mb-8">
        <div class="flex w-full">
          <div class="flex-1">
            <h2 class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
              {$LL.battles()}
            </h2>
          </div>
        </div>

        <div class="flex flex-wrap">
          {#each battles as battle}
            <div
              class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
            >
              <div class="flex flex-wrap items-center p-4">
                <div
                  class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                >
                  <span data-testid="battle-name">{battle.name}</span>
                </div>
                <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                  <HollowButton href="{appRoutes.game}/{battle.id}">
                    {$LL.battleJoin()}
                  </HollowButton>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    {#if FeatureRetro}
      <div class="w-full mb-6 lg:mb-8">
        <div class="flex w-full">
          <div class="flex-1">
            <h2 class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
              {$LL.retros()}
            </h2>
          </div>
        </div>

        <div class="flex flex-wrap">
          {#each retros as retro}
            <div
              class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
            >
              <div class="flex flex-wrap items-center p-4">
                <div
                  class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                >
                  <span data-testid="retro-name">{retro.name}</span>
                </div>
                <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                  <HollowButton href="{appRoutes.retro}/{retro.id}">
                    {$LL.joinRetro()}
                  </HollowButton>
                </div>
              </div>
            </div>
          {/each}
        </div>

        {#if retros.length}
          <div class="w-full pt-4 px-4">
            <div class="w-full">
              <h3 class="text-xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
                {$LL.retroActionItems()}
              </h3>

              <div class="text-right mb-4">
                <Toggle
                  name="completedActionItems"
                  id="completedActionItems"
                  bind:checked={completedActionItems}
                  changeHandler={changeRetroActionCompletedToggle}
                  label={$LL.showCompletedActionItems()}
                />
              </div>
            </div>

            <Table>
              {#snippet header()}
                <tr>
                  <HeadCol>{$LL.actionItem()}</HeadCol>
                  <HeadCol>{$LL.completed()}</HeadCol>
                  <HeadCol>{$LL.comments()}</HeadCol>
                </tr>
              {/snippet}
              {#snippet body({ class: className })}
                <tbody class={className}>
                  {#each retroActions as item, i}
                    <TableRow itemIndex={i}>
                      <RowCol>
                        <div class="whitespace-pre-wrap">
                          {item.content}
                        </div>
                      </RowCol>
                      <RowCol>
                        <input
                          type="checkbox"
                          id="{i}Completed"
                          checked={item.completed}
                          class="opacity-0 absolute h-6 w-6"
                          disabled
                        />
                        <div
                          class="bg-white dark:bg-gray-800 border-2 rounded-md
                                              border-gray-400 dark:border-gray-300 w-6 h-6 flex flex-shrink-0
                                              justify-center items-center me-2
                                              focus-within:border-blue-500 dark:focus-within:border-sky-500"
                        >
                          <Check class="hidden w-4 h-4 text-green-600 pointer-events-none" />
                        </div>
                        <label for="{i}Completed" class="select-none"></label>
                      </RowCol>
                      <RowCol>
                        <MessageSquareMore width="22" height="22" class="inline-block" />
                        <button
                          class="text-lg text-blue-400 dark:text-sky-400"
                          onclick={toggleRetroActionComments(item.id)}
                        >
                          &nbsp;{item.comments.length}
                        </button>
                      </RowCol>
                    </TableRow>
                  {/each}
                </tbody>
              {/snippet}
            </Table>
            <TableFooter
              bind:current={retroActionsPage}
              num_items={totalRetroActions}
              per_page={retroActionsPageLimit}
              on:navigate={changeRetroActionPage}
            />
          </div>
        {/if}
      </div>
    {/if}

    {#if FeatureStoryboard}
      <div class="w-full mb-6 lg:mb-8">
        <div class="flex w-full">
          <div class="flex-1">
            <h2 class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
              {$LL.storyboards()}
            </h2>
          </div>
        </div>

        <div class="flex flex-wrap">
          {#each storyboards as storyboard}
            <div
              class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
            >
              <div class="flex flex-wrap items-center p-4">
                <div
                  class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                >
                  <span data-testid="storyboard-name">{storyboard.name}</span>
                </div>
                <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                  <HollowButton href="{appRoutes.storyboard}/{storyboard.id}">
                    {$LL.joinStoryboard()}
                  </HollowButton>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <TableContainer>
      <TableNav title={$LL.users()} createBtnEnabled={false} />
      <Table>
        {#snippet header()}
          <tr>
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.email()}
            </HeadCol>
            <HeadCol>
              {$LL.role()}
            </HeadCol>
          </tr>
        {/snippet}
        {#snippet body({ class: className })}
          <tbody class={className}>
            {#each users as user, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10">
                      <UserAvatar
                        warriorId={user.id}
                        avatar={user.avatar}
                        gravatarHash={user.gravatarHash}
                        userName={user.name}
                        width={48}
                        class="h-10 w-10 rounded-full"
                      />
                    </div>
                    <div class="ms-4">
                      <div class="font-medium text-gray-900 dark:text-gray-200">
                        <a
                          data-testid="user-name"
                          href="{appRoutes.adminUsers}/{user.id}"
                          class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                          >{user.name}</a
                        >
                        {#if user.country}
                          &nbsp;
                          <CountryFlag country={user.country} additionalClass="inline-block" width="32" height="24" />
                        {/if}
                      </div>
                    </div>
                  </div>
                </RowCol>
                <RowCol>
                  <span data-testid="user-email">{user.email}</span>
                </RowCol>
                <RowCol>
                  <div class="text-sm text-gray-500 dark:text-gray-300">
                    {user.role}
                  </div>
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>

    {#if !organizationId}
      <div class="text-center mt-4">
        <HollowButton color="red" onClick={toggleDeleteTeam} testid="team-delete">
          {$LL.deleteTeam()}
        </HollowButton>
      </div>
    {/if}

    {#if showDeleteTeam}
      <DeleteConfirmation
        toggleDelete={toggleDeleteTeam}
        handleDelete={deleteTeam}
        confirmText={$LL.deleteTeamConfirmText()}
        confirmBtnText={$LL.deleteTeam()}
      />
    {/if}
  </div>
</AdminPageLayout>

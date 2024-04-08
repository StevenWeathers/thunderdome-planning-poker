<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import ChevronRight from '../../components/icons/ChevronRight.svelte';
  import CreateBattle from '../../components/poker/CreatePokerGame.svelte';
  import CreateRetro from '../../components/retro/CreateRetro.svelte';
  import CreateStoryboard from '../../components/storyboard/CreateStoryboard.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import ActionComments from '../../components/retro/ActionComments.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import Modal from '../../components/global/Modal.svelte';
  import EditActionItem from '../../components/retro/EditActionItem.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import CheckboxIcon from '../../components/icons/CheckboxIcon.svelte';
  import CommentIcon from '../../components/icons/CommentIcon.svelte';
  import BoxList from '../../components/BoxList.svelte';
  import UsersList from '../../components/team/UsersList.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;
  export let organizationId;
  export let departmentId;
  export let teamId;

  const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

  const battlesPageLimit = 1000;
  const retrosPageLimit = 1000;
  const retroActionsPageLimit = 5;
  const storyboardsPageLimit = 1000;
  const usersPageLimit = 1000;

  let team = {
    id: teamId,
    name: '',
  };
  let organization = {
    id: organizationId,
    name: '',
  };
  let department = {
    id: departmentId,
    name: '',
  };

  let users = [];
  let battles = [];
  let retros = [];
  let retroActions = [];
  let storyboards = [];
  let invites = [];
  let showCreateBattle = false;
  let showCreateRetro = false;
  let showCreateStoryboard = false;
  let showRemoveBattle = false;
  let showRemoveRetro = false;
  let showRemoveStoryboard = false;
  let showDeleteTeam = false;
  let removeBattleId = null;
  let removeRetroId = null;
  let removeStoryboardId = null;
  let usersPage = 1;
  let battlesPage = 1;
  let retrosPage = 1;
  let retroActionsPage = 1;
  let storyboardsPage = 1;
  let totalRetroActions = 0;
  let completedActionItems = false;
  let showInviteUser = false;
  let showDeleteInvite = false;
  let deleteInviteId = '';

  let organizationRole = '';
  let departmentRole = '';
  let teamRole = '';
  let isAdmin = false;
  let isTeamMember = false;

  const apiPrefix = '/api';
  $: orgPrefix = departmentId
    ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
    : `${apiPrefix}/organizations/${organizationId}`;
  $: teamPrefix = organizationId
    ? `${orgPrefix}/teams/${teamId}`
    : `${apiPrefix}/teams/${teamId}`;

  $: currentPageUrl = teamPrefix
    .replace('/api', '')
    .replace('organizations', 'organization')
    .replace('departments', 'department')
    .replace('teams', 'team');

  function toggleCreateBattle() {
    showCreateBattle = !showCreateBattle;
  }

  function toggleCreateRetro() {
    showCreateRetro = !showCreateRetro;
  }

  function toggleCreateStoryboard() {
    showCreateStoryboard = !showCreateStoryboard;
  }

  const toggleRemoveBattle = battleId => () => {
    showRemoveBattle = !showRemoveBattle;
    removeBattleId = battleId;
  };

  const toggleRemoveRetro = retroId => () => {
    showRemoveRetro = !showRemoveRetro;
    removeRetroId = retroId;
  };

  const toggleRemoveStoryboard = storyboardId => () => {
    showRemoveStoryboard = !showRemoveStoryboard;
    removeStoryboardId = storyboardId;
  };

  const toggleDeleteTeam = () => {
    showDeleteTeam = !showDeleteTeam;
  };

  const toggleInviteUser = () => {
    showInviteUser = !showInviteUser;
  };

  const toggleDeleteInvite = inviteId => () => {
    showDeleteInvite = !showDeleteInvite;
    deleteInviteId = inviteId;
  };

  let showRetroActionComments = false;
  let selectedRetroAction = null;
  const toggleRetroActionComments = id => () => {
    showRetroActionComments = !showRetroActionComments;
    selectedRetroAction = id;
  };

  function getTeam() {
    xfetch(teamPrefix)
      .then(res => res.json())
      .then(function (result) {
        team = result.data.team;
        teamRole = result.data.teamRole;

        if (departmentId) {
          department = result.data.department;
          departmentRole = result.data.departmentRole;
        }
        if (organizationId) {
          organization = result.data.organization;
          organizationRole = result.data.organizationRole;
        }

        isAdmin =
          organizationRole === 'ADMIN' ||
          departmentRole === 'ADMIN' ||
          teamRole === 'ADMIN';
        isTeamMember =
          organizationRole === 'ADMIN' ||
          departmentRole === 'ADMIN' ||
          teamRole !== '';

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

  function getInvites() {
    xfetch(`${teamPrefix}/invites`)
      .then(res => res.json())
      .then(function (result) {
        invites = result.data;
      })
      .catch(function () {
        notifications.danger('error getting team invites');
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${teamPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
        if (isAdmin) {
          getInvites();
        }
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function getBattles() {
    if (FeaturePoker) {
      const battlesOffset = (battlesPage - 1) * battlesPageLimit;
      xfetch(
        `${teamPrefix}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
      )
        .then(res => res.json())
        .then(function (result) {
          battles = result.data;
        })
        .catch(function () {
          notifications.danger(
            $LL.teamGetBattlesError({
              friendly: AppConfig.FriendlyUIVerbs,
            }),
          );
        });
    }
  }

  function getRetros() {
    if (FeatureRetro) {
      const retrosOffset = (retrosPage - 1) * retrosPageLimit;
      xfetch(
        `${teamPrefix}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`,
      )
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
          selectedAction =
            selectedAction !== null
              ? retroActions.find(r => r.id === selectedAction.id)
              : null;
        })
        .catch(function () {
          notifications.danger($LL.teamGetRetroActionsError());
        });
    }
  }

  function getStoryboards() {
    if (FeatureStoryboard) {
      const storyboardsOffset = (storyboardsPage - 1) * storyboardsPageLimit;
      xfetch(
        `${teamPrefix}/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}`,
      )
        .then(res => res.json())
        .then(function (result) {
          storyboards = result.data;
        })
        .catch(function () {
          notifications.danger($LL.teamGetStoryboardsError());
        });
    }
  }

  function handleBattleRemove() {
    xfetch(`${teamPrefix}/battles/${removeBattleId}`, { method: 'DELETE' })
      .then(function () {
        eventTag('team_remove_battle', 'engagement', 'success');
        toggleRemoveBattle(null)();
        notifications.success(
          $LL.battleRemoveSuccess({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
        getBattles();
      })
      .catch(function () {
        notifications.danger(
          $LL.battleRemoveError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
        eventTag('team_remove_battle', 'engagement', 'failure');
      });
  }

  function handleRetroRemove() {
    xfetch(`${teamPrefix}/retros/${removeRetroId}`, { method: 'DELETE' })
      .then(function () {
        eventTag('team_remove_retro', 'engagement', 'success');
        toggleRemoveRetro(null)();
        notifications.success($LL.retroRemoveSuccess());
        getRetros();
      })
      .catch(function () {
        notifications.danger($LL.retroRemoveError());
        eventTag('team_remove_retro', 'engagement', 'failure');
      });
  }

  function handleStoryboardRemove() {
    xfetch(`${teamPrefix}/storyboards/${removeStoryboardId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('team_remove_storyboard', 'engagement', 'success');
        toggleRemoveStoryboard(null)();
        notifications.success($LL.storyboardRemoveSuccess());
        getStoryboards();
      })
      .catch(function () {
        notifications.danger($LL.storyboardRemoveError());
        eventTag('team_remove_storyboard', 'engagement', 'failure');
      });
  }

  function handleDeleteInvite() {
    xfetch(`${teamPrefix}/invites/${deleteInviteId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('team_delete_invite', 'engagement', 'success');
        toggleDeleteInvite(null)();
        notifications.success('Successfully deleted user invite');
        getInvites();
      })
      .catch(function () {
        notifications.danger('Error deleting user invite');
        eventTag('team_delete_invite', 'engagement', 'failure');
      });
  }

  function handleDeleteTeam() {
    xfetch(`${teamPrefix}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('team_delete', 'engagement', 'success');
        toggleDeleteTeam();
        notifications.success($LL.teamDeleteSuccess());
        router.route(appRoutes.teams);
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
        eventTag('team_delete', 'engagement', 'failure');
      });
  }

  const changeRetroActionPage = evt => {
    retroActionsPage = evt.detail;
    getRetrosActions();
  };

  const changeRetroActionCompletedToggle = () => {
    retroActionsPage = 1;
    getRetrosActions();
  };

  let showRetroActionEdit = false;
  let selectedAction = null;
  let selectedActionRetroId = null;
  const toggleRetroActionEdit = (retroId, id) => () => {
    showRetroActionEdit = !showRetroActionEdit;
    selectedAction = retroActions.find(r => r.id === id);
    selectedActionRetroId = retroId;
  };

  function handleRetroActionEdit(action) {
    xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
      method: 'PUT',
      body: {
        content: action.content,
        completed: action.completed,
      },
    })
      .then(function () {
        getRetrosActions();
        toggleRetroActionEdit(null)();
        notifications.success($LL.updateActionItemSuccess());
        eventTag('team_action_update', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger($LL.updateActionItemError());
        eventTag('team_action_update', 'engagement', 'failure');
      });
  }

  function handleRetroActionDelete(action) {
    return () => {
      xfetch(`/api/retros/${action.retroId}/actions/${action.id}`, {
        method: 'DELETE',
      })
        .then(function () {
          getRetrosActions();
          toggleRetroActionEdit(null)();
          notifications.success($LL.deleteActionItemSuccess());
          eventTag('team_action_delete', 'engagement', 'success');
        })
        .catch(function () {
          notifications.danger($LL.deleteActionItemError());
          eventTag('team_action_delete', 'engagement', 'failure');
        });
    };
  }

  function handleRetroActionAssigneeAdd(retroId, actionId, userId) {
    xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
      method: 'POST',
      body: {
        user_id: userId,
      },
    })
      .then(function () {
        getRetrosActions();
        eventTag('team_action_assignee_add', 'engagement', 'success');
      })
      .catch(function () {
        eventTag('team_action_assignee_add', 'engagement', 'failure');
      });
  }

  function handleRetroActionAssigneeRemove(retroId, actionId, userId) {
    return () => {
      xfetch(`/api/retros/${retroId}/actions/${actionId}/assignees`, {
        method: 'DELETE',
        body: {
          user_id: userId,
        },
      })
        .then(function () {
          getRetrosActions();
          eventTag('team_action_assignee_remove', 'engagement', 'success');
        })
        .catch(function () {
          eventTag('team_action_assignee_remove', 'engagement', 'failure');
        });
    };
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getTeam();
  });
</script>

<svelte:head>
  <title>{$LL.team()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex mb-6 lg:mb-8">
    <div class="flex-1">
      <h1 class="text-3xl font-semibold font-rajdhani dark:text-white">
        <span class="uppercase">{$LL.team()}</span>
        <ChevronRight class="w-8 h-8" />
        {team.name}
      </h1>

      {#if organizationId}
        <div class="text-xl font-semibold font-rajdhani dark:text-white">
          <span class="uppercase">{$LL.organization()}</span>
          <ChevronRight />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.organization}/{organization.id}"
          >
            {organization.name}
          </a>
          {#if departmentId}
            &nbsp;
            <ChevronRight />
            <span class="uppercase">{$LL.department()}</span>
            <ChevronRight />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}"
            >
              {department.name}
            </a>
          {/if}
        </div>
      {/if}
    </div>
    <div class="flex-1 text-right">
      <SolidButton
        additionalClasses="font-rajdhani uppercase text-2xl"
        href="{`${currentPageUrl}/checkin`}"
        >{$LL.checkins()}
      </SolidButton>
    </div>
  </div>

  {#if FeaturePoker}
    <div class="w-full mb-6 lg:mb-8">
      <div class="flex w-full">
        <div class="flex-1">
          <h2
            class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
          >
            {$LL.battles({ friendly: AppConfig.FriendlyUIVerbs })}
          </h2>
        </div>
        <div class="flex-1 text-right">
          {#if isTeamMember}
            <SolidButton onClick="{toggleCreateBattle}"
              >{$LL.battleCreate({
                friendly: AppConfig.FriendlyUIVerbs,
              })}
            </SolidButton>
          {/if}
        </div>
      </div>

      <div class="flex flex-wrap">
        <BoxList
          items="{battles}"
          itemType="battle"
          pageRoute="{appRoutes.game}"
          joinBtnText="{$LL.battleJoin({
            friendly: AppConfig.FriendlyUIVerbs,
          })}"
          isAdmin="{isAdmin}"
          toggleRemove="{toggleRemoveBattle}"
        />
      </div>
    </div>

    {#if showCreateBattle}
      <Modal closeModal="{toggleCreateBattle}">
        <CreateBattle
          apiPrefix="{teamPrefix}"
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
          showOwner="{false}"
        />
      </Modal>
    {/if}
  {/if}

  {#if FeatureRetro}
    <div class="w-full mb-6 lg:mb-8">
      <div class="flex w-full">
        <div class="flex-1">
          <h2
            class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
          >
            {$LL.retros()}
          </h2>
        </div>
        <div class="flex-1 text-right">
          {#if isTeamMember}
            <SolidButton onClick="{toggleCreateRetro}"
              >{$LL.createRetro()}</SolidButton
            >
          {/if}
        </div>
      </div>

      <div class="flex flex-wrap">
        <BoxList
          items="{retros}"
          itemType="retro"
          pageRoute="{appRoutes.retro}"
          joinBtnText="{$LL.joinRetro()}"
          isAdmin="{isAdmin}"
          toggleRemove="{toggleRemoveRetro}"
          showOwner="{false}"
        />
      </div>

      {#if retros.length}
        <div class="w-full pt-4 px-4">
          <TableContainer>
            <TableNav
              title="{$LL.retroActionItems()}"
              createBtnEnabled="{false}"
            >
              <div
                class="relative inline-block w-10 me-2 align-middle select-none transition duration-200 ease-in"
              >
                <input
                  type="checkbox"
                  name="completedActionItems"
                  id="completedActionItems"
                  bind:checked="{completedActionItems}"
                  on:change="{changeRetroActionCompletedToggle}"
                  class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                  for="completedActionItems"
                  class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
              </div>
              <label for="completedActionItems" class="dark:text-gray-300"
                >{$LL.showCompletedActionItems()}</label
              >
            </TableNav>
            <Table>
              <tr slot="header">
                <HeadCol>{$LL.actionItem()}</HeadCol>
                <HeadCol>{$LL.completed()}</HeadCol>
                <HeadCol>{$LL.comments()}</HeadCol>
                <HeadCol />
              </tr>
              <tbody slot="body" let:class="{className}" class="{className}">
                {#each retroActions as item, i}
                  <TableRow itemIndex="{i}">
                    <RowCol>
                      <div class="whitespace-pre-wrap">
                        {#each item.assignees as assignee}
                          <UserAvatar
                            warriorId="{assignee.id}"
                            gravatarHash="{assignee.gravatarHash}"
                            avatar="{assignee.avatar}"
                            userName="{assignee.name}"
                            width="24"
                            class="inline-block me-2"
                          />
                        {/each}{item.content}
                      </div>
                    </RowCol>
                    <RowCol>
                      <input
                        type="checkbox"
                        id="{i}Completed"
                        checked="{item.completed}"
                        class="opacity-0 absolute h-6 w-6"
                        disabled
                      />
                      <div
                        class="bg-white dark:bg-gray-800 border-2 rounded-md
                                            border-gray-400 dark:border-gray-300 w-6 h-6 flex flex-shrink-0
                                            justify-center items-center me-2
                                            focus-within:border-blue-500 dark:focus-within:border-sky-500"
                      >
                        <CheckboxIcon />
                      </div>
                      <label for="{i}Completed" class="select-none"></label>
                    </RowCol>
                    <RowCol>
                      <CommentIcon width="22" height="22" />
                      <button
                        class="text-lg text-blue-400 dark:text-sky-400"
                        on:click="{toggleRetroActionComments(item.id)}"
                      >
                        &nbsp;{item.comments.length}
                      </button>
                    </RowCol>
                    <RowCol>
                      <div class="text-right">
                        <HollowButton
                          color="teal"
                          onClick="{toggleRetroActionEdit(
                            item.retroId,
                            item.id,
                          )}"
                          >{$LL.edit()}
                        </HollowButton>
                      </div>
                    </RowCol>
                  </TableRow>
                {/each}
              </tbody>
            </Table>
            <TableFooter
              bind:current="{retroActionsPage}"
              num_items="{totalRetroActions}"
              per_page="{retroActionsPageLimit}"
              on:navigate="{changeRetroActionPage}"
            />
          </TableContainer>
        </div>
      {/if}
    </div>

    {#if showCreateRetro}
      <Modal closeModal="{toggleCreateRetro}">
        <CreateRetro
          apiPrefix="{teamPrefix}"
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
        />
      </Modal>
    {/if}
  {/if}

  {#if FeatureStoryboard}
    <div class="w-full mb-6 lg:mb-8">
      <div class="flex w-full">
        <div class="flex-1">
          <h2
            class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
          >
            {$LL.storyboards()}
          </h2>
        </div>
        <div class="flex-1 text-right">
          {#if isTeamMember}
            <SolidButton onClick="{toggleCreateStoryboard}"
              >{$LL.createStoryboard()}
            </SolidButton>
          {/if}
        </div>
      </div>

      <div class="flex flex-wrap">
        <BoxList
          items="{storyboards}"
          itemType="storyboard"
          pageRoute="{appRoutes.storyboard}"
          joinBtnText="{$LL.joinStoryboard()}"
          isAdmin="{isAdmin}"
          toggleRemove="{toggleRemoveStoryboard}"
          showOwner="{false}"
        />
      </div>
    </div>

    {#if showCreateStoryboard}
      <Modal closeModal="{toggleCreateStoryboard}">
        <CreateStoryboard
          apiPrefix="{teamPrefix}"
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
        />
      </Modal>
    {/if}
  {/if}

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <TableContainer>
        <TableNav title="{$LL.userInvites()}" createBtnEnabled="{false}" />
        <Table>
          <tr slot="header">
            <HeadCol>{$LL.email()}</HeadCol>
            <HeadCol>{$LL.role()}</HeadCol>
            <HeadCol>{$LL.dateCreated()}</HeadCol>
            <HeadCol>{$LL.expireDate()}</HeadCol>
            <HeadCol />
          </tr>
          <tbody slot="body" let:class="{className}" class="{className}">
            {#each invites as item, i}
              <TableRow itemIndex="{i}">
                <RowCol>
                  {item.email}
                </RowCol>
                <RowCol>
                  {item.role}
                </RowCol>
                <RowCol>
                  {new Date(item.created_date).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(item.expire_date).toLocaleString()}
                </RowCol>
                <RowCol>
                  <div class="text-right">
                    <HollowButton
                      onClick="{toggleDeleteInvite(item.invite_id)}"
                      color="red"
                    >
                      {$LL.delete()}
                    </HollowButton>
                  </div>
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        </Table>
      </TableContainer>
    </div>
  {/if}

  <UsersList
    users="{users}"
    getUsers="{getUsers}"
    xfetch="{xfetch}"
    eventTag="{eventTag}"
    notifications="{notifications}"
    isAdmin="{isAdmin}"
    pageType="team"
    teamPrefix="{teamPrefix}"
  />

  {#if isAdmin && !organizationId && !departmentId}
    <div class="w-full text-center mt-8">
      <HollowButton onClick="{toggleDeleteTeam}" color="red">
        {$LL.deleteTeam()}
      </HollowButton>
    </div>
  {/if}

  {#if showRemoveBattle}
    <DeleteConfirmation
      toggleDelete="{toggleRemoveBattle(null)}"
      handleDelete="{handleBattleRemove}"
      permanent="{false}"
      confirmText="{$LL.removeBattleConfirmText({
        friendly: AppConfig.FriendlyUIVerbs,
      })}"
      confirmBtnText="{$LL.removeBattle({
        friendly: AppConfig.FriendlyUIVerbs,
      })}"
    />
  {/if}

  {#if showRemoveRetro}
    <DeleteConfirmation
      toggleDelete="{toggleRemoveRetro(null)}"
      handleDelete="{handleRetroRemove}"
      permanent="{false}"
      confirmText="{$LL.removeRetroConfirmText()}"
      confirmBtnText="{$LL.removeRetro()}"
    />
  {/if}

  {#if showRemoveStoryboard}
    <DeleteConfirmation
      toggleDelete="{toggleRemoveStoryboard(null)}"
      handleDelete="{handleStoryboardRemove}"
      permanent="{false}"
      confirmText="{$LL.removeStoryboardConfirmText()}"
      confirmBtnText="{$LL.removeStoryboard()}"
    />
  {/if}

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteTeam}"
      handleDelete="{handleDeleteTeam}"
      confirmText="{$LL.deleteTeamConfirmText()}"
      confirmBtnText="{$LL.deleteTeam()}"
    />
  {/if}

  {#if showRetroActionEdit}
    <EditActionItem
      toggleEdit="{toggleRetroActionEdit(null)}"
      handleEdit="{handleRetroActionEdit}"
      handleDelete="{handleRetroActionDelete}"
      assignableUsers="{users}"
      action="{selectedAction}"
      handleAssigneeAdd="{handleRetroActionAssigneeAdd}"
      handleAssigneeRemove="{handleRetroActionAssigneeRemove}"
      retroId="{selectedActionRetroId}"
    />
  {/if}

  {#if showRetroActionComments}
    <ActionComments
      toggleComments="{toggleRetroActionComments(null)}"
      actions="{retroActions}"
      users="{users}"
      selectedActionId="{selectedRetroAction}"
      getRetrosActions="{getRetrosActions}"
      xfetch="{xfetch}"
      eventTag="{eventTag}"
      notifications="{notifications}"
      isAdmin="{isAdmin}"
    />
  {/if}

  {#if showDeleteInvite}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteInvite(null)}"
      handleDelete="{handleDeleteInvite}"
      confirmText="{$LL.userInviteConfirmDelete()}"
      confirmBtnText="{$LL.userInviteDelete()}"
    />
  {/if}
</PageLayout>

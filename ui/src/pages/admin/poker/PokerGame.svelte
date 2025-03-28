<script lang="ts">
  import { onMount } from 'svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import Table from '../../../components/table/Table.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import BooleanDisplay from '../../../components/global/BooleanDisplay.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let battleId;

  let showDeleteBattle = false;

  let battle = {
    name: '',
    votingLocked: false,
    autoFinishVoting: false,
    activePlanId: '',
    pointValuesAllowed: [],
    pointAverageRounding: '',
    users: [],
    plans: [],
    createdDate: '',
    updatedDate: '',
  };

  function getBattle() {
    xfetch(`/api/battles/${battleId}`)
      .then(res => res.json())
      .then(function (result) {
        battle = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getBattleError());
      });
  }

  function deleteBattle() {
    xfetch(`/api/battles/${battleId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        router.route(appRoutes.adminPokerGames);
      })
      .catch(function () {
        notifications.danger($LL.deleteBattleError());
      });
  }

  const toggleDeleteBattle = () => {
    showDeleteBattle = !showDeleteBattle;
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getBattle();
  });
</script>

<svelte:head>
  <title
    >{$LL.battles()}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="battles">
  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title="{battle.name}" createBtnEnabled="{false}" />
      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.votingLocked()}
          </HeadCol>
          <HeadCol>
            {$LL.autoFinishVoting()}
          </HeadCol>
          <HeadCol>
            {$LL.pointValuesAllowed()}
          </HeadCol>
          <HeadCol>
            {$LL.pointAverageRounding()}
          </HeadCol>
          <HeadCol>
            {$LL.dateCreated()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          <TableRow itemIndex="{0}">
            <RowCol>
              <BooleanDisplay boolValue="{battle.votingLocked}" />
            </RowCol>
            <RowCol>
              <BooleanDisplay boolValue="{battle.autoFinishVoting}" />
            </RowCol>
            <RowCol>
              {battle.pointValuesAllowed.join(', ')}
            </RowCol>
            <RowCol>
              {battle.pointAverageRounding}
            </RowCol>
            <RowCol>
              {new Date(battle.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(battle.updatedDate).toLocaleString()}
            </RowCol>
          </TableRow>
        </tbody>
      </Table>
    </TableContainer>
  </div>

  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title="{$LL.users()}" createBtnEnabled="{false}" />
      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.abandoned()}
          </HeadCol>
          <HeadCol>
            {$LL.spectator()}
          </HeadCol>
          <HeadCol>
            {$LL.leader()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each battle.users as user, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId="{user.id}"
                      avatar="{user.avatar}"
                      gravatarHash="{user.gravatarHash}"
                      userName="{user.name}"
                      width="48"
                      class="h-10 w-10 rounded-full"
                    />
                  </div>
                  <div class="ms-4">
                    <div class="text-sm font-medium text-gray-900">
                      <a
                        href="{appRoutes.adminUsers}/{user.id}"
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        >{user.name}</a
                      >
                      {#if user.country}
                        &nbsp;
                        <CountryFlag
                          country="{user.country}"
                          additionalClass="inline-block"
                          width="32"
                          height="24"
                        />
                      {/if}
                    </div>
                  </div>
                </div>
              </RowCol>
              <RowCol>
                {user.rank}
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue="{user.active}" />
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue="{user.abandoned}" />
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue="{user.spectator}" />
              </RowCol>
              <RowCol>
                <BooleanDisplay
                  boolValue="{battle.leaders.includes(user.id)}"
                />
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </TableContainer>
  </div>

  <TableContainer>
    <TableNav title="{$LL.plans()}" createBtnEnabled="{false}" />
    <Table>
      <tr slot="header">
        <HeadCol>
          {$LL.name()}
        </HeadCol>
        <HeadCol>
          {$LL.type()}
        </HeadCol>
        <HeadCol>
          {$LL.planReferenceId()}
        </HeadCol>
        <HeadCol>
          {$LL.voteCount()}
        </HeadCol>
        <HeadCol>
          {$LL.points()}
        </HeadCol>
        <HeadCol>
          {$LL.active()}
        </HeadCol>
        <HeadCol>
          {$LL.skipped()}
        </HeadCol>
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each battle.plans as plan, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              {plan.name}
            </RowCol>
            <RowCol>
              {plan.type}
            </RowCol>
            <RowCol>
              {plan.referenceId}
            </RowCol>
            <RowCol>
              {plan.votes.length}
            </RowCol>
            <RowCol>
              {plan.points}
            </RowCol>
            <RowCol>
              <BooleanDisplay boolValue="{plan.active}" />
            </RowCol>
            <RowCol>
              <BooleanDisplay boolValue="{plan.skipped}" />
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>
  </TableContainer>

  <div class="text-center mt-4">
    <HollowButton
      color="red"
      onClick="{toggleDeleteBattle}"
      testid="battle-delete"
    >
      {$LL.battleDelete()}
    </HollowButton>
  </div>

  {#if showDeleteBattle}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteBattle}"
      handleDelete="{deleteBattle}"
      confirmText="{$LL.deleteBattleConfirmText()}"
      confirmBtnText="{$LL.deleteBattle()}"
    />
  {/if}
</AdminPageLayout>

<script lang="ts">
  import { onMount } from 'svelte';
  import VerifiedIcon from '../../components/icons/VerifiedIcon.svelte';
  import UpdatePasswordForm from '../../components/user/UpdatePasswordForm.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import Modal from '../../components/global/Modal.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;
  export let userId;

  const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

  let showUpdatePassword = false;

  let userDetails = {
    id: '',
    name: '',
    email: '',
    rank: '',
    avatar: '',
    verified: false,
    notificationsEnabled: true,
    country: '',
    locale: '',
    company: '',
    jobTitle: '',
    createdDate: '',
    updatedDate: '',
    lastActive: '',
  };

  function toggleUpdatePassword() {
    showUpdatePassword = !showUpdatePassword;
  }

  function getUser() {
    xfetch(`/api/users/${userId}`)
      .then(res => res.json())
      .then(function (result) {
        userDetails = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getUserError());
      });
  }

  const battlesPageLimit = 100;
  let battleCount = 0;
  let battles = [];
  let battlesPage = 1;

  function getBattles() {
    const battlesOffset = (battlesPage - 1) * battlesPageLimit;
    xfetch(
      `/api/users/${userId}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        battles = result.data;
        battleCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger(
          $LL.getBattlesError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
      });
  }

  const changeBattlesPage = evt => {
    battlesPage = evt.detail;
    getBattles();
  };

  const retrosPageLimit = 100;
  let retroCount = 0;
  let retros = [];
  let retrosPage = 1;

  function getRetros() {
    const offset = (retrosPage - 1) * retrosPageLimit;
    xfetch(
      `/api/users/${userId}/retros?limit=${retrosPageLimit}&offset=${offset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        retros = result.data;
        retroCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getRetrosError());
      });
  }

  const changeRetrosPage = evt => {
    retrosPage = evt.detail;
    getRetros();
  };

  const storyboardsPageLimit = 100;
  let storyboardCount = 0;
  let storyboards = [];
  let storyboardsPage = 1;

  function getStoryboards() {
    const offset = (storyboardsPage - 1) * storyboardsPageLimit;
    xfetch(
      `/api/users/${userId}/storyboards?limit=${storyboardsPageLimit}&offset=${offset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        storyboards = result.data;
        storyboardCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardsError());
      });
  }

  const changeStoryboardsPage = evt => {
    storyboardsPage = evt.detail;
    getStoryboards();
  };

  function updatePassword(password1, password2) {
    const body = {
      password1,
      password2,
    };

    xfetch(`/api/admin/users/${userId}/password`, { body, method: 'PATCH' })
      .then(function () {
        notifications.success($LL.passwordUpdated(), 1500);
        toggleUpdatePassword();
        eventTag('update_password', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger($LL.passwordUpdateError());
        eventTag('update_password', 'engagement', 'failure');
      });
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

    getUser();
    FeaturePoker && getBattles();
    FeatureRetro && getRetros();
    FeatureStoryboard && getStoryboards();
  });
</script>

<svelte:head>
  <title>{$LL.users()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="users">
  <div class="mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title="{userDetails.name}" createBtnEnabled="{false}">
        <SolidButton onClick="{toggleUpdatePassword}"
          >{$LL.updatePassword()}</SolidButton
        >
      </TableNav>
      <Table>
        <tr slot="header">
          <HeadCol />
          <HeadCol>
            {$LL.country()}
          </HeadCol>
          <HeadCol>
            {$LL.email()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.dateCreated()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
          <HeadCol>
            {$LL.lastActive()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          <TableRow itemIndex="{0}">
            <RowCol>
              <div class="flex items-center flex-nowrap">
                <div class="flex-shrink-0 h-10 w-10">
                  <UserAvatar
                    warriorId="{userDetails.id}"
                    avatar="{userDetails.avatar}"
                    gravatarHash="{userDetails.gravatarHash}"
                    userName="{user.name}"
                    width="48"
                    class="h-10 w-10 rounded-full"
                  />
                </div>
              </div>
            </RowCol>
            <RowCol>
              {#if userDetails.country}
                &nbsp;
                <CountryFlag
                  country="{userDetails.country}"
                  additionalClass="inline-block"
                  width="32"
                  height="24"
                />
              {/if}
            </RowCol>
            <RowCol>
              {userDetails.email}
              {#if userDetails.verified}
                <span class="text-green-600" title="{$LL.verified()}">
                  <VerifiedIcon />
                </span>
              {/if}
            </RowCol>
            <RowCol>
              {userDetails.rank}
            </RowCol>
            <RowCol>
              {new Date(userDetails.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(userDetails.updatedDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(userDetails.lastActive).toLocaleString()}
            </RowCol>
          </TableRow>
        </tbody>
      </Table>
    </TableContainer>
  </div>

  {#if FeaturePoker}
    <div class="mb-6 lg:mb-8">
      <TableContainer>
        <TableNav
          title="{$LL.battles({ friendly: AppConfig.FriendlyUIVerbs })}"
          createBtnEnabled="{false}"
        />
        <Table>
          <tr slot="header">
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
            <HeadCol type="action">
              <span class="sr-only">{$LL.actions()}</span>
            </HeadCol>
          </tr>
          <tbody slot="body" let:class="{className}" class="{className}">
            {#each battles as battle, i}
              <TableRow itemIndex="{i}">
                <RowCol>
                  <a
                    href="{appRoutes.admin}/battles/{battle.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >{battle.name}</a
                  >
                </RowCol>
                <RowCol>
                  {new Date(battle.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(battle.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  <HollowButton href="{appRoutes.game}/{battle.id}">
                    {$LL.battleJoin({
                      friendly: AppConfig.FriendlyUIVerbs,
                    })}
                  </HollowButton>
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        </Table>
        <TableFooter
          bind:current="{battlesPage}"
          num_items="{battleCount}"
          per_page="{battlesPageLimit}"
          on:navigate="{changeBattlesPage}"
        />
      </TableContainer>
    </div>
  {/if}

  {#if FeatureRetro}
    <div class="mb-6 lg:mb-8">
      <TableContainer>
        <TableNav title="{$LL.retros()}" createBtnEnabled="{false}" />
        <Table>
          <tr slot="header">
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
            <HeadCol type="action">
              <span class="sr-only">{$LL.actions()}</span>
            </HeadCol>
          </tr>
          <tbody slot="body" let:class="{className}" class="{className}">
            {#each retros as retro, i}
              <TableRow itemIndex="{i}">
                <RowCol>
                  <a
                    href="{appRoutes.admin}/retros/{retro.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >{retro.name}</a
                  >
                </RowCol>
                <RowCol>
                  {new Date(retro.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(retro.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  <HollowButton href="{appRoutes.retro}/{retro.id}">
                    {$LL.joinRetro()}
                  </HollowButton>
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        </Table>
        <TableFooter
          bind:current="{retrosPage}"
          num_items="{retroCount}"
          per_page="{retrosPageLimit}"
          on:navigate="{changeRetrosPage}"
        />
      </TableContainer>
    </div>
  {/if}

  {#if FeatureStoryboard}
    <div class="mb-6 lg:mb-8">
      <TableContainer>
        <TableNav title="{$LL.storyboards()}" createBtnEnabled="{false}" />
        <Table>
          <tr slot="header">
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
            <HeadCol type="action">
              <span class="sr-only">{$LL.actions()}</span>
            </HeadCol>
          </tr>
          <tbody slot="body" let:class="{className}" class="{className}">
            {#each storyboards as storyboard, i}
              <TableRow itemIndex="{i}">
                <RowCol>
                  <a
                    href="{appRoutes.admin}/storyboards/{storyboard.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >{storyboard.name}</a
                  >
                </RowCol>
                <RowCol>
                  {new Date(storyboard.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(storyboard.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  <HollowButton href="{appRoutes.storyboard}/{storyboard.id}">
                    {$LL.joinStoryboard()}
                  </HollowButton>
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        </Table>
        <TableFooter
          bind:current="{storyboardsPage}"
          num_items="{storyboardCount}"
          per_page="{storyboardsPageLimit}"
          on:navigate="{changeStoryboardsPage}"
        />
      </TableContainer>
    </div>
  {/if}

  {#if showUpdatePassword}
    <Modal closeModal="{toggleUpdatePassword}">
      <UpdatePasswordForm
        handleUpdate="{updatePassword}"
        toggleForm="{toggleUpdatePassword}"
        notifications="{notifications}"
      />
    </Modal>
  {/if}
</AdminPageLayout>

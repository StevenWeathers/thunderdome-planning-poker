<script lang="ts">
  import { onMount } from 'svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import DeleteStoryboard from '../../components/storyboard/DeleteStoryboard.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let storyboardId;

  let storyboard = {
    name: '',
    users: [],
    owner_id: '',
    createdDate: '',
    updatedDate: '',
  };

  let showDeleteStoryboard = false;

  function getStoryboard() {
    xfetch(`/api/storyboards/${storyboardId}`)
      .then(res => res.json())
      .then(function (result) {
        storyboard = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardErrorMessage());
      });
  }

  function deleteStoryboard() {
    xfetch(`/api/storyboards/${storyboardId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function () {
        router.route(appRoutes.adminStoryboards);
      })
      .catch(function () {
        notifications.danger($LL.deleteStoryboardErrorMessage());
      });
  }

  function toggleDeleteStoryboard() {
    showDeleteStoryboard = !showDeleteStoryboard;
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

    getStoryboard();
  });
</script>

<svelte:head>
  <title>{$LL.storyboard()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="storyboards">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
    >
      {storyboard.name}
    </h1>
  </div>

  <div class="w-full">
    <div class="p-4 md:p-6">
      <Table>
        <tr slot="header">
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
              {new Date(storyboard.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(storyboard.updatedDate).toLocaleString()}
            </RowCol>
          </TableRow>
        </tbody>
      </Table>
    </div>
    <div class="p-4 md:p-6">
      <h3
        class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
      >
        {$LL.users()}
      </h3>

      <Table>
        <tr slot="header">
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.abandoned()}
          </HeadCol>
          <HeadCol>
            {$LL.leader()}
          </HeadCol>
        </tr>
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each storyboard.users as user, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId="{user.id}"
                      avatar="{user.avatar}"
                      gravatarHash="{user.gravatarHash}"
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
                {#if user.active}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                {#if user.abandoned}
                  <span class="text-green-600"><CheckIcon /></span>
                {/if}
              </RowCol>
              <RowCol>
                <RowCol>
                  {#if storyboard.owner_id === user.id}
                    <span class="text-green-600"><CheckIcon /></span>
                  {/if}
                </RowCol>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>

      <div class="text-center mt-4">
        <HollowButton
          color="red"
          onClick="{toggleDeleteStoryboard}"
          testid="storyboard-delete"
        >
          {$LL.deleteStoryboard()}
        </HollowButton>
      </div>

      {#if showDeleteStoryboard}
        <DeleteStoryboard
          toggleDelete="{toggleDeleteStoryboard}"
          handleDelete="{deleteStoryboard}"
        />
      {/if}
    </div>
  </div>
</AdminPageLayout>

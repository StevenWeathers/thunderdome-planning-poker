<script lang="ts">
    import { onMount } from 'svelte';
    import { user } from '../../stores';
    import { validateUserIsAdmin, validateUserIsRegistered } from '../../validationUtils';
    import { appRoutes, AppConfig } from '../../config';
    import ProjectPageLayout from '../../components/project/ProjectPageLayout.svelte';
    import CreateBattle from '../../components/poker/CreatePokerGame.svelte';
    import CreateRetro from '../../components/retro/CreateRetro.svelte';
    import CreateStoryboard from '../../components/storyboard/CreateStoryboard.svelte';
    import type { ApiClient } from '../../types/apiclient';
    import type { NotificationService } from '../../types/notifications';
    import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
    import SolidButton from '../../components/global/SolidButton.svelte';
    import BoxList from '../../components/BoxList.svelte';
    import LL from '../../i18n/i18n-svelte';
    import Modal from '../../components/global/Modal.svelte';
    import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import { get } from 'svelte/store';

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

    interface Props {
        xfetch: ApiClient;
        router: any;
        notifications: NotificationService;
        projectId: string;
    }

    let {
        xfetch,
        router,
        notifications,
        projectId,
    }: Props = $props();

    // Basic project state
    let project = $state({
        id: projectId,
        projectKey: '',
        name: '',
        description: ''
    });

    const apiPrefix = '/api';

    const battlesPageLimit = 1000;
    const retrosPageLimit = 1000;
    const storyboardsPageLimit = 1000;

    let projectPrefix = `${apiPrefix}/projects/${projectId}`;
    let isProjectAdmin = $state(false);
    let showCreateBattle = $state(false);
    let showCreateRetro = $state(false);
    let showCreateStoryboard = $state(false);
    let showRemoveBattle = $state(false);
    let showRemoveRetro = $state(false);
    let showRemoveStoryboard = $state(false);
    let removeBattleId = $state(null);
    let removeRetroId = $state(null);
    let removeStoryboardId = $state(null);
    let battles = $state([]);
    let retros = $state([]);
    let storyboards = $state([]);
    let battlesPage = $state(1);
    let retrosPage = $state(1);
    let storyboardsPage = $state(1);

    function getProject() {
        xfetch(projectPrefix)
            .then(res => res.json())
            .then(result => {
                const data = result?.data || result;
                if (data) {
                    project = { ...project, ...data };
                    isProjectAdmin = result?.meta?.role.toLowerCase() === 'admin' || validateUserIsAdmin($user);
                }
                getBattles();
                getRetros();
                getStoryboards();
            })
            .catch(() => {
            notifications.danger('Failed to get project');
            });
    }

    function getBattles() {
    if (FeaturePoker) {
      const battlesOffset = (battlesPage - 1) * battlesPageLimit;
      xfetch(
        `${projectPrefix}/poker?limit=${battlesPageLimit}&offset=${battlesOffset}`,
      )
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
      xfetch(
        `${projectPrefix}/retros?limit=${retrosPageLimit}&offset=${retrosOffset}`,
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

  function getStoryboards() {
    if (FeatureStoryboard) {
      const storyboardsOffset = (storyboardsPage - 1) * storyboardsPageLimit;
      xfetch(
        `${projectPrefix}/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}`,
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
    xfetch(`${projectPrefix}/poker/${removeBattleId}`, { method: 'DELETE' })
      .then(function () {
        toggleRemoveBattle(null)();
        notifications.success($LL.battleRemoveSuccess());
        getBattles();
      })
      .catch(function () {
        notifications.danger($LL.battleRemoveError());
      });
  }

  function handleRetroRemove() {
    xfetch(`${projectPrefix}/retros/${removeRetroId}`, { method: 'DELETE' })
      .then(function () {
        toggleRemoveRetro(null)();
        notifications.success($LL.retroRemoveSuccess());
        getRetros();
      })
      .catch(function () {
        notifications.danger($LL.retroRemoveError());
      });
  }

  function handleStoryboardRemove() {
    xfetch(`${projectPrefix}/storyboards/${removeStoryboardId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleRemoveStoryboard(null)();
        notifications.success($LL.storyboardRemoveSuccess());
        getStoryboards();
      })
      .catch(function () {
        notifications.danger($LL.storyboardRemoveError());
      });
  }

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

    onMount(() => {
        if (!$user.id || !validateUserIsRegistered($user)) {
            router.route(appRoutes.login);
            return;
        }
        getProject();
    });
</script>

<svelte:head>
    <title>Project {project.name}</title>
</svelte:head>

<ProjectPageLayout activePage="project" {projectId}>
  <div class="container mx-auto px-4 py-4 md:py-6 lg:py-8">
      <h1 class="text-3xl font-semibold font-rajdhani dark:text-white mb-4" data-testid="project-title">Project {project.name}</h1>

      <!-- Subscription Required if enabled -->
      {#if AppConfig.SubscriptionsEnabled && $user && !$user.isSubscriber}
          <FeatureSubscribeBanner
            salesPitch="Active subscription required to access this project."
          />
      {:else}
          {#if FeaturePoker}
              <div class="w-full mb-6 lg:mb-8">
                <div class="flex w-full">
                    <div class="flex-1">
                    <h2
                        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                    >
                        {$LL.battles()}
                    </h2>
                    </div>
                    <div class="flex-1 text-right">
                        <SolidButton onClick={toggleCreateBattle}
                        >{$LL.battleCreate()}
                        </SolidButton>
                    </div>
                </div>

                <div class="flex flex-wrap">
                    <BoxList
                    items={battles}
                    itemType="battle"
                    pageRoute={appRoutes.game}
                    joinBtnText={$LL.battleJoin()}
                    isAdmin={isProjectAdmin}
                    toggleRemove={toggleRemoveBattle}
                    />
                </div>
              </div>

              {#if showCreateBattle}
                <Modal closeModal={toggleCreateBattle}>
                    <CreateBattle
                      apiPrefix={projectPrefix}
                      scope={"project"}
                      notifications={notifications}
                      router={router}
                      xfetch={xfetch}
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
                        <SolidButton onClick={toggleCreateRetro}
                        >{$LL.createRetro()}</SolidButton
                        >
                    </div>
                </div>

                <div class="flex flex-wrap">
                    <BoxList
                    items={retros}
                    itemType="retro"
                    pageRoute={appRoutes.retro}
                    joinBtnText={$LL.joinRetro()}
                    isAdmin={isProjectAdmin}
                    toggleRemove={toggleRemoveRetro}
                    showOwner={false}
                    />
                </div>
              </div>
              {#if showCreateRetro}
                  <Modal closeModal={toggleCreateRetro}>
                      <CreateRetro
                        apiPrefix={projectPrefix}
                        scope={"project"}
                        notifications={notifications}
                        router={router}
                        xfetch={xfetch}
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
                        <SolidButton onClick={toggleCreateStoryboard}
                        >{$LL.createStoryboard()}
                        </SolidButton>
                    </div>
                </div>

                <div class="flex flex-wrap">
                    <BoxList
                    items={storyboards}
                    itemType="storyboard"
                    pageRoute={appRoutes.storyboard}
                    joinBtnText={$LL.joinStoryboard()}
                    isAdmin={isProjectAdmin}
                    toggleRemove={toggleRemoveStoryboard}
                    showOwner={false}
                    />
                </div>
              </div>

              {#if showCreateStoryboard}
              <Modal closeModal={toggleCreateStoryboard}>
                  <CreateStoryboard
                    apiPrefix={projectPrefix}
                    scope={"project"}
                    notifications={notifications}
                    router={router}
                    xfetch={xfetch}
                  />
              </Modal>
              {/if}
          {/if}
      {/if}

      {#if showRemoveBattle}
          <DeleteConfirmation
          toggleDelete={toggleRemoveBattle(null)}
          handleDelete={handleBattleRemove}
          permanent={false}
          confirmText={$LL.removeBattleConfirmText()}
          confirmBtnText={$LL.removeBattle()}
          />
      {/if}

      {#if showRemoveRetro}
          <DeleteConfirmation
          toggleDelete={toggleRemoveRetro(null)}
          handleDelete={handleRetroRemove}
          permanent={false}
          confirmText={$LL.removeRetroConfirmText()}
          confirmBtnText={$LL.removeRetro()}
          />
      {/if}

      {#if showRemoveStoryboard}
          <DeleteConfirmation
          toggleDelete={toggleRemoveStoryboard(null)}
          handleDelete={handleStoryboardRemove}
          permanent={false}
          confirmText={$LL.removeStoryboardConfirmText()}
          confirmBtnText={$LL.removeStoryboard()}
          />
      {/if}
      
  </div>
</ProjectPageLayout>
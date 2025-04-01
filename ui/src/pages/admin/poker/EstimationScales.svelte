<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import EstimationScalesList from '../../../components/estimationscale/EstimationScalesList.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
  }

  let { xfetch, router, notifications }: Props = $props();

  const scalesPageLimit = 100;
  let scaleCount = $state(0);
  let scalesPage = $state(1);
  let scales = $state([]);

  function getScales() {
    const scalesOffset = (scalesPage - 1) * scalesPageLimit;
    xfetch(
      `/api/admin/estimation-scales?limit=${scalesPageLimit}&offset=${scalesOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        scales = result.data;
        scaleCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger('Failed to fetch scales');
      });
  }

  const changePage = evt => {
    scalesPage = evt.detail;
    getScales();
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

    getScales();
  });
</script>

<svelte:head>
  <title
    >{$LL.estimationScales()}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="estimation-scales">
  <EstimationScalesList
    xfetch={xfetch}
    notifications={notifications}
    scales={scales}
    apiPrefix="/api/admin"
    getScales={getScales}
    changePage={changePage}
    scaleCount={scaleCount}
    scalesPage={scalesPage}
    scalesPageLimit={scalesPageLimit}
  />
</AdminPageLayout>

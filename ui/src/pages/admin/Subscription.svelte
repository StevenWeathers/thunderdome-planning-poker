<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import Table from '../../components/table/Table.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import BooleanDisplay from '../../components/global/BooleanDisplay.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    subscriptionId: any;
  }

  let { xfetch, router, notifications, subscriptionId }: Props = $props();

  let defaultSubscription = $state({
    id: '',
    user_id: '',
    customer_id: '',
    subscription_id: '',
    active: false,
    type: '',
    expires: '',
    user: {
      name: '',
    },
  });

  $effect(() => {
    defaultSubscription.id = subscriptionId;
  });

  let subscription = $state({
    ...defaultSubscription,
  });

  function getSubscription() {
    xfetch(`/api/subscriptions/${subscriptionId}`)
      .then(res => res.json())
      .then(function (result) {
        subscription = result.data;
      })
      .catch(function () {
        notifications.danger('Error getting subscription');
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

    getSubscription();
  });
</script>

<svelte:head>
  <title>Subscription {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="subscriptions">
  <TableContainer>
    <TableNav title="Subscription" createBtnEnabled={false} />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>
            {$LL.email()}
          </HeadCol>
          <HeadCol>Customer ID</HeadCol>
          <HeadCol>Subscription ID</HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.dateCreated()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
          <HeadCol type="action">
            <span class="sr-only">Actions</span>
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          <TableRow>
            <RowCol>
              <a
                href="{appRoutes.adminUsers}/{subscription.user_id}"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              >
                {subscription.user.name}
              </a>
            </RowCol>
            <RowCol>
              {subscription.customer_id}
            </RowCol>
            <RowCol>
              {subscription.subscription_id}
            </RowCol>
            <RowCol>
              {subscription.type}
            </RowCol>
            <RowCol>
              <BooleanDisplay boolValue={subscription.active} />
            </RowCol>
            <RowCol>
              {new Date(subscription.created_date).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(subscription.updated_date).toLocaleString()}
            </RowCol>
            <RowCol type="action" />
          </TableRow>
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>
</AdminPageLayout>

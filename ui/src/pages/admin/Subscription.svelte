<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';

  const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig;

  export let xfetch;
  export let router;
  export let notifications;
  export let subscriptionId;
  export let eventTag;

  let defaultSubscription = {
    id: subscriptionId,
    user_id: '',
    customer_id: '',
    subscription_id: '',
    active: false,
    type: '',
    expires: '',
    user: {
      name: '',
    },
  };

  let subscription = {
    ...defaultSubscription,
  };

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
  <div class="px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
    >
      Subscription
    </h1>
  </div>

  <div class="w-full">
    <div class="mb-4">
      <Table>
        <tr slot="header">
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
        <tbody slot="body" let:class="{className}" class="{className}">
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
              {#if subscription.active}
                <span class="text-green-600">
                  <CheckIcon />
                </span>
              {/if}
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
      </Table>
    </div>
  </div>
</AdminPageLayout>

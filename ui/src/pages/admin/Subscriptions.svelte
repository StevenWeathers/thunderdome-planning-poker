<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import Table from '../../components/table/Table.svelte';
  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import SubscriptionForm from '../../components/subscription/SubscriptionForm.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import BooleanDisplay from '../../components/global/BooleanDisplay.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
  }

  let { xfetch, router, notifications }: Props = $props();

  const subscriptionsPageLimit = 25;
  let subscriptionCount = $state(0);
  let subscriptions = $state([]);
  let subscriptionsPage = $state(1);

  let showSubCreate = $state(false);
  let showSubUpdate = $state(false);
  let showSubDelete = $state(false);
  let deleteSubId = null;
  let defaultSubscription = {
    id: '',
    user_id: '',
    team_id: '',
    organization_id: '',
    customer_id: '',
    subscription_id: '',
    active: false,
    type: '',
    expires: '',
    user: {
      name: '',
    },
  };
  let selectedSub = $state({
    ...defaultSubscription,
  });

  function toggleSubCreate() {
    showSubCreate = !showSubCreate;
  }

  const toggleSubUpdate = sub => () => {
    showSubUpdate = !showSubUpdate;
    selectedSub = sub;
  };

  const toggleSubDelete = subId => () => {
    showSubDelete = !showSubDelete;
    deleteSubId = subId;
  };

  function getSubscriptions() {
    const subscriptionsOffset =
      (subscriptionsPage - 1) * subscriptionsPageLimit;
    xfetch(
      `/api/subscriptions?limit=${subscriptionsPageLimit}&offset=${subscriptionsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        subscriptions = result.data;
        subscriptionCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger('Error getting subscriptions');
      });
  }

  function deleteSub() {
    xfetch(`/api/subscriptions/${deleteSubId}`, { method: 'DELETE' })
      .then(function () {
        toggleSubDelete(null)();
        getSubscriptions();
        notifications.success('Subscription deleted');
      })
      .catch(function () {
        toggleSubDelete(null)();
        notifications.danger('Error deleting subscription');
      });
  }

  const changePage = evt => {
    subscriptionsPage = evt.detail;
    getSubscriptions();
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

    getSubscriptions();
  });
</script>

<svelte:head>
  <title>{$LL.subscriptions()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="subscriptions">
  <TableContainer>
    <TableNav
      title="Subscriptions"
      createBtnText="Create Subscription"
      createButtonHandler="{toggleSubCreate}"
      createBtnTestId="subscription-create"
    />
    <Table>
      {#snippet header()}
            <tr >
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
          <HeadCol>Expires</HeadCol>
          <HeadCol type="action">
            <span class="sr-only">Actions</span>
          </HeadCol>
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
          {#each subscriptions as subscription, i}
            <TableRow itemIndex="{i}">
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
                <BooleanDisplay boolValue="{subscription.active}" />
              </RowCol>
              <RowCol>
                {new Date(subscription.expires).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                <CrudActions
                  detailsLink="{appRoutes.adminSubscriptions}/{subscription.id}"
                  editBtnClickHandler="{toggleSubUpdate(subscription)}"
                  deleteBtnClickHandler="{toggleSubDelete(subscription.id)}"
                />
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
    <TableFooter
      bind:current="{subscriptionsPage}"
      num_items="{subscriptionCount}"
      per_page="{subscriptionsPageLimit}"
      on:navigate="{changePage}"
    />

    {#if showSubCreate}
      <SubscriptionForm
        toggleClose="{toggleSubCreate}"
        handleUpdate="{getSubscriptions}"
        xfetch={xfetch}
        notifications={notifications}
      />
    {/if}
    {#if showSubUpdate}
      <SubscriptionForm
        toggleClose="{toggleSubUpdate({ ...defaultSubscription })}"
        handleUpdate="{getSubscriptions}"
        xfetch={xfetch}
        notifications={notifications}
        subscriptionId="{selectedSub.id}"
        customer_id="{selectedSub.customer_id}"
        subscription_id="{selectedSub.subscription_id}"
        user_id="{selectedSub.user_id}"
        team_id="{selectedSub.team_id}"
        organization_id="{selectedSub.organization_id}"
        active="{selectedSub.active}"
        expires="{selectedSub.expires}"
      />
    {/if}

    {#if showSubDelete}
      <DeleteConfirmation
        toggleDelete="{toggleSubDelete(null)}"
        handleDelete="{deleteSub}"
        confirmText="{$LL.deleteSubscriptionConfirmation()}"
        confirmBtnText="{$LL.deleteSubscription()}"
      />
    {/if}
  </TableContainer>
</AdminPageLayout>

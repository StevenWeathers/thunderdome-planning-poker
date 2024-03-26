<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import Pagination from '../../components/global/Pagination.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';
  import SubscriptionForm from '../../components/subscription/SubscriptionForm.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;

  const subscriptionsPageLimit = 25;
  let subscriptionCount = 0;
  let subscriptions = [];
  let subscriptionsPage = 1;

  let showSubCreate = false;
  let showSubUpdate = false;
  let showSubDelete = false;
  let deleteSubId = null;
  let defaultSubscription = {
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
  };
  let selectedSub = {
    ...defaultSubscription,
  };

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
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      Subscriptions
    </h1>
  </div>

  <div class="w-full">
    <div class="text-right mb-4">
      <HollowButton onClick="{toggleSubCreate}">
        Create Subscription
      </HollowButton>
    </div>

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
        <HeadCol>Expires</HeadCol>
        <HeadCol type="action">
          <span class="sr-only">Actions</span>
        </HeadCol>
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
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
              {#if subscription.active}
                <span class="text-green-600">
                  <CheckIcon />
                </span>
              {/if}
            </RowCol>
            <RowCol>
              {new Date(subscription.expires).toLocaleString()}
            </RowCol>
            <RowCol type="action">
              <a
                href="{appRoutes.adminSubscriptions}/{subscription.id}"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >Details</a
              >
              <HollowButton
                onClick="{toggleSubUpdate(subscription)}"
                color="blue"
              >
                {$LL.edit()}
              </HollowButton>
              <HollowButton
                onClick="{toggleSubDelete(subscription.id)}"
                color="red"
              >
                {$LL.delete()}
              </HollowButton>
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if subscriptionCount > subscriptionsPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{subscriptionsPage}"
          num_items="{subscriptionCount}"
          per_page="{subscriptionsPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}

    {#if showSubCreate}
      <SubscriptionForm
        toggleClose="{toggleSubCreate}"
        handleUpdate="{getSubscriptions}"
        eventTag="{eventTag}"
        xfetch="{xfetch}"
        notifications="{notifications}"
      />
    {/if}
    {#if showSubUpdate}
      <SubscriptionForm
        toggleClose="{toggleSubUpdate({ ...defaultSubscription })}"
        handleUpdate="{getSubscriptions}"
        eventTag="{eventTag}"
        xfetch="{xfetch}"
        notifications="{notifications}"
        subscriptionId="{selectedSub.id}"
        customer_id="{selectedSub.customer_id}"
        subscription_id="{selectedSub.subscription_id}"
        user_id="{selectedSub.user_id}"
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
  </div>
</AdminPageLayout>

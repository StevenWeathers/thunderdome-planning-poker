<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import TextInput from '../forms/TextInput.svelte';

  export let handleUpdate = () => {};
  export let toggleClose = () => {};
  export let eventTag = (a, b, c) => {};
  export let xfetch = async (url, options) => {};
  export let notifications;

  export let subscriptionId = '';
  export let user_id = '';
  export let team_id = '';
  export let organization_id = '';
  export let customer_id = '';
  export let subscription_id = '';
  export let type = 'user';
  export let active = true;
  export let expires = new Date().toISOString();

  function handleSubmit(event) {
    event.preventDefault();

    if (user_id === '') {
      notifications.danger('user_id field required');
      eventTag('subscription_form_user_id_invalid', 'engagement', 'failure');
      return false;
    }

    if (customer_id === '') {
      notifications.danger('customer_id field required');
      eventTag(
        'subscription_form_customer_id_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    if (subscription_id === '') {
      notifications.danger('subscription_id field required');
      eventTag(
        'subscription_form_subscription_id_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    if (type === '') {
      notifications.danger('type field required');
      eventTag('subscription_form_type_invalid', 'engagement', 'failure');
      return false;
    }

    if (active === '') {
      notifications.danger('active field required');
      eventTag('subscription_form_active_invalid', 'engagement', 'failure');
      return false;
    }

    const body = {
      user_id,
      team_id,
      organization_id,
      customer_id,
      subscription_id,
      type,
      active,
      expires,
    };

    const endpoint =
      subscriptionId != ''
        ? `/api/subscriptions/${subscriptionId}`
        : `/api/subscriptions`;
    const method = subscriptionId != '' ? 'PUT' : 'POST';

    xfetch(endpoint, { body, method: method })
      .then((res: any) => res.json())
      .then(function () {
        handleUpdate();
        toggleClose();
        eventTag('subscription_form', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger('failed to update subscription');

        eventTag('subscription_form', 'engagement', 'failure');
      });
  }
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{handleSubmit}" name="subscriptionform">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="userId">
        Thunderdome User Id<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="userId"
        name="userId"
        bind:value="{user_id}"
        placeholder="Enter the User Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="teamId">
        Thunderdome Team Id
      </label>
      <TextInput
        id="teamId"
        name="teamId"
        bind:value="{team_id}"
        placeholder="Enter the associated Team Id..."
      />
    </div>
    <div class="mb-4">
      <label
        class="block dark:text-gray-400 font-bold mb-2"
        for="organizationId"
      >
        Thunderdome Organization Id
      </label>
      <TextInput
        id="organizationId"
        name="organizationId"
        bind:value="{organization_id}"
        placeholder="Enter the associated Organization Id..."
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="customerId">
        Stripe Customer Id<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="customerId"
        name="customerId"
        bind:value="{customer_id}"
        placeholder="Enter the Stripe Customer Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label
        class="block dark:text-gray-400 font-bold mb-2"
        for="subscriptionId"
      >
        Stripe Subscription Id<span class="text-red-500 dark:text-red-400"
          >*</span
        >
      </label>
      <TextInput
        id="subscriptionId"
        name="subscriptionId"
        bind:value="{subscription_id}"
        placeholder="Enter the Stripe Subscription Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="subType">
        Subscription Type<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="subType"
        name="subType"
        bind:value="{type}"
        placeholder="Enter the Subscription Type..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
        <input
          type="checkbox"
          bind:checked="{active}"
          id="active"
          name="active"
          class="w-4 h-4 dark:accent-lime-400 me-1"
        />
        {$LL.active()}
      </label>
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="expires">
        Expires<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="expires"
        name="expires"
        bind:value="{expires}"
        placeholder="Enter the expiration date..."
        required
      />
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">
          {$LL.save()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>

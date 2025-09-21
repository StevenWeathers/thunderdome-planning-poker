<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import TextInput from '../forms/TextInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';

  import type { NotificationService } from '../../types/notifications';
  import { onMount } from 'svelte';

  interface Props {
    handleUpdate?: any;
    toggleClose?: any;
    xfetch?: any;
    notifications: NotificationService;
    subscriptionId?: string;
    user_id?: string;
    team_id?: string;
    organization_id?: string;
    customer_id?: string;
    subscription_id?: string;
    type?: string;
    active?: boolean;
    expires?: any;
  }

  let {
    handleUpdate = () => {},
    toggleClose = () => {},
    xfetch = async (url, options) => {},
    notifications,
    subscriptionId = '',
    user_id = $bindable(''),
    team_id = $bindable(''),
    organization_id = $bindable(''),
    customer_id = $bindable(''),
    subscription_id = $bindable(''),
    type = $bindable('user'),
    active = $bindable(true),
    expires = $bindable(new Date().toISOString()),
  }: Props = $props();

  function handleSubmit(event) {
    event.preventDefault();

    if (user_id === '') {
      notifications.danger('user_id field required');
      return false;
    }

    if (customer_id === '') {
      notifications.danger('customer_id field required');
      return false;
    }

    if (subscription_id === '') {
      notifications.danger('subscription_id field required');
      return false;
    }

    if (type === '') {
      notifications.danger('type field required');
      return false;
    }

    if (active === '') {
      notifications.danger('active field required');
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

    const endpoint = subscriptionId != '' ? `/api/subscriptions/${subscriptionId}` : `/api/subscriptions`;
    const method = subscriptionId != '' ? 'PUT' : 'POST';

    xfetch(endpoint, { body, method: method })
      .then((res: any) => res.json())
      .then(function () {
        handleUpdate();
        toggleClose();
      })
      .catch(function () {
        notifications.danger('failed to update subscription');
      });
  }

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleClose} ariaLabel={$LL.modalSubscriptionForm()}>
  <form onsubmit={handleSubmit} name="subscriptionform">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="userId">
        Thunderdome User Id<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="userId"
        name="userId"
        bind:value={user_id}
        bind:this={focusInput}
        placeholder="Enter the User Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="teamId"> Thunderdome Team Id </label>
      <TextInput id="teamId" name="teamId" bind:value={team_id} placeholder="Enter the associated Team Id..." />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="organizationId"> Thunderdome Organization Id </label>
      <TextInput
        id="organizationId"
        name="organizationId"
        bind:value={organization_id}
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
        bind:value={customer_id}
        placeholder="Enter the Stripe Customer Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="subscriptionId">
        Stripe Subscription Id<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput
        id="subscriptionId"
        name="subscriptionId"
        bind:value={subscription_id}
        placeholder="Enter the Stripe Subscription Id..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="subType">
        Subscription Type<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput id="subType" name="subType" bind:value={type} placeholder="Enter the Subscription Type..." required />
    </div>
    <div class="mb-4">
      <Checkbox bind:checked={active} id="active" name="active" label={$LL.active()} />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="expires">
        Expires<span class="text-red-500 dark:text-red-400">*</span>
      </label>
      <TextInput id="expires" name="expires" bind:value={expires} placeholder="Enter the expiration date..." required />
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

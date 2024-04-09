<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';

  export let handleUpdate = () => {};
  export let toggleClose = () => {};
  export let eventTag = (a, b, c) => {};
  export let xfetch = async (url, ...options) => {};
  export let notifications;

  export let subscriptionId = '';

  let organizations = [];
  let selectedOrganization = '';

  xfetch(`/api/users/${$user.id}/organizations?limit=1000&offset=0`)
    .then(res => res.json())
    .then(function (result) {
      organizations = result.data;
    })
    .catch(function () {
      notifications.danger($LL.getOrganizationsError());
    });

  function handleSubmit(event) {
    event.preventDefault();

    if (selectedOrganization === '') {
      notifications.danger('Select a organization field required');
      eventTag(
        'subscription_form_associated_organization_id_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    const body = {
      organization_id: selectedOrganization,
    };

    xfetch(`/api/users/${$user.id}/subscriptions/${subscriptionId}`, {
      body,
      method: 'PATCH',
    })
      .then((res: any) => res.json())
      .then(function () {
        handleUpdate();
        toggleClose();
        eventTag(
          'subscription_form_associate_organization',
          'engagement',
          'success',
        );
      })
      .catch(function () {
        notifications.danger(
          'failed to associate organization to subscription',
        );

        eventTag(
          'subscription_form_associate_organization',
          'engagement',
          'failure',
        );
      });
  }
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{handleSubmit}" name="associateOrganizationForm">
    <div class="mb-4">
      <label
        class="block dark:text-gray-400 font-bold mb-2"
        for="selectedOrganization"
      >
        Associate Organization
      </label>
      <SelectInput
        bind:value="{selectedOrganization}"
        id="selectedOrganization"
        name="selectedOrganization"
      >
        <option value="" disabled>Select Organization</option>
        {#each organizations as organization}
          <option value="{organization.id}">
            {organization.name}
          </option>
        {/each}
      </SelectInput>
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

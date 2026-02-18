<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { Organization } from '../../types/organization';
  import { onMount } from 'svelte';

  interface Props {
    handleUpdate?: any;
    toggleClose?: any;
    xfetch?: any;
    notifications: NotificationService;
    subscriptionId?: string;
  }

  let {
    handleUpdate = () => {},
    toggleClose = () => {},
    xfetch = async (url: string, ...options: any[]) => {},
    notifications,
    subscriptionId = '',
  }: Props = $props();

  let organizations = $state<Organization[]>([]);
  let selectedOrganization = $state('');
  let focusInput: any = $state();

  onMount(() => {
    xfetch(`/api/users/${$user.id}/organizations?limit=1000&offset=0`)
      .then((res: any) => res.json())
      .then(function (result: any) {
        organizations = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getOrganizationsError());
      });

    focusInput?.focus();
  });

  function handleSubmit(event: Event) {
    event.preventDefault();

    if (selectedOrganization === '') {
      notifications.danger('Select a organization field required');
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
      })
      .catch(function () {
        notifications.danger('failed to associate organization to subscription');
      });
  }
</script>

<Modal closeModal={toggleClose} ariaLabel={$LL.modalAssociateOrganizationToSubscription()}>
  <form onsubmit={handleSubmit} name="associateOrganizationForm">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="selectedOrganization"> Associate Organization </label>
      <SelectInput
        bind:value={selectedOrganization}
        bind:this={focusInput}
        id="selectedOrganization"
        name="selectedOrganization"
      >
        <option value="" disabled>Select Organization</option>
        {#each organizations as organization}
          <option value={organization.id}>
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

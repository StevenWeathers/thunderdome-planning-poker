<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';


  interface Props {
    handleUpdate?: any;
    toggleClose?: any;
    xfetch?: any;
    notifications: any;
    subscriptionId?: string;
  }

  let {
    handleUpdate = () => {},
    toggleClose = () => {},
    xfetch = async (url, ...options) => {},
    notifications,
    subscriptionId = ''
  }: Props = $props();

  let organizations = $state([]);
  let selectedOrganization = $state('');

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
        notifications.danger(
          'failed to associate organization to subscription',
        );
      });
  }
</script>

<Modal closeModal={toggleClose}>
  <form onsubmit={handleSubmit} name="associateOrganizationForm">
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

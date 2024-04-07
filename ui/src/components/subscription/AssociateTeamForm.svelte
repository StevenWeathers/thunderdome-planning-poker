<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../global/SelectInput.svelte';
  import { user } from '../../stores';

  export let handleUpdate = () => {};
  export let toggleClose = () => {};
  export let eventTag = (a, b, c) => {};
  export let xfetch = async (url, ...options) => {};
  export let notifications;

  export let subscriptionId = '';

  let teams = [];
  let selectedTeam = '';

  xfetch(`/api/users/${$user.id}/teams?limit=1000&offset=0`)
    .then(res => res.json())
    .then(function (result) {
      teams = result.data;
    })
    .catch(function () {
      notifications.danger(`${$LL.getTeamsError()}`);
    });

  function handleSubmit(event) {
    event.preventDefault();

    if (selectedTeam === '') {
      notifications.danger('Select a team field required');
      eventTag(
        'subscription_form_associated_team_id_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    const body = {
      team_id: selectedTeam,
    };

    xfetch(`/api/users/${$user.id}/subscriptions/${subscriptionId}`, {
      body,
      method: 'PATCH',
    })
      .then((res: any) => res.json())
      .then(function () {
        handleUpdate();
        toggleClose();
        eventTag('subscription_form_associate_team', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger('failed to associate team to subscription');

        eventTag('subscription_form_associate_team', 'engagement', 'failure');
      });
  }
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{handleSubmit}" name="associateTeamForm">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="selectedTeam">
        {$LL.associateTeam()}
      </label>
      <SelectInput
        bind:value="{selectedTeam}"
        id="selectedTeam"
        name="selectedTeam"
      >
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value="{team.id}">
            {team.name}
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

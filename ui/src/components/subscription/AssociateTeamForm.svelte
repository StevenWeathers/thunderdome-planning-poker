<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';

  import type { NotificationService } from '../../types/notifications';
  import type { Team } from '../../types/team';
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

  let teams = $state<Team[]>([]);
  let selectedTeam = $state('');
  let focusInput: any = $state();

  onMount(() => {
    xfetch(`/api/users/${$user.id}/teams?limit=1000&offset=0`)
      .then((res: any) => res.json())
      .then(function (result: any) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getTeamsError());
      });

    focusInput?.focus();
  });

  function handleSubmit(event: Event) {
    event.preventDefault();

    if (selectedTeam === '') {
      notifications.danger('Select a team field required');
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
      })
      .catch(function () {
        notifications.danger('failed to associate team to subscription');
      });
  }
</script>

<Modal closeModal={toggleClose} ariaLabel={$LL.modalAssociateTeamToSubscription()}>
  <form onsubmit={handleSubmit} name="associateTeamForm">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="selectedTeam">
        {$LL.associateTeam()}
      </label>
      <SelectInput bind:value={selectedTeam} bind:this={focusInput} id="selectedTeam" name="selectedTeam">
        <option value="" disabled>{$LL.selectTeam()}</option>
        {#each teams as team}
          <option value={team.id}>
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

<script lang="ts">
  import type { NotificationService } from '../../types/notifications';
  import { onMount } from 'svelte';

  interface Props {
    teamId?: string;
    userId?: string;
    xfetch?: any;
    notifications: NotificationService;
  }

  let { teamId = '', userId = '', xfetch = async () => {}, notifications }: Props = $props();

  let team = $state({
    name: '',
  });

  onMount(() => {
    xfetch(`/api/teams/${teamId}`)
      .then((res: any) => res.json())
      .then(function (result: any) {
        team = result.data.team;
      })
      .catch(function () {
        notifications.danger('Error getting associated team');
      });
  });
</script>

{team.name}

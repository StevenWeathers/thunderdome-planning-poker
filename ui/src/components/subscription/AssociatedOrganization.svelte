<script lang="ts">
  interface Props {
    organizationId?: string;
    xfetch?: any;
    notifications: NotificationService;
  }

  import type { NotificationService } from '../../types/notifications';
  import { onMount } from 'svelte';

  let { organizationId = '', xfetch = async () => {}, notifications }: Props = $props();

  let organization = $state({
    name: '',
  });

  onMount(() => {
    xfetch(`/api/organizations/${organizationId}`)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
      })
      .catch(function () {
        notifications.danger('Error getting associated organization');
      });
  });
</script>

{organization.name}

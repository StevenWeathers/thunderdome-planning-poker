<script lang="ts">
  interface Props {
    teamId?: string;
    userId?: string;
    xfetch?: any;
    notifications: any;
  }

  let {
    teamId = '',
    userId = '',
    xfetch = async () => {},
    notifications
  }: Props = $props();

  let team = $state({
    name: '',
  });

  xfetch(`/api/teams/${teamId}`)
    .then(res => res.json())
    .then(function (result) {
      team = result.data.team;
    })
    .catch(function () {
      notifications.danger('Error getting associated team');
    });
</script>

{team.name}

<script lang="ts">
  export let teamId = '';
  export let userId = '';
  export let xfetch = async () => {};
  export let notifications;
  export let eventTag;

  let team = {
    name: '',
  };

  xfetch(`/api/teams/${teamId}`)
    .then(res => res.json())
    .then(function (result) {
      team = result.data.team;
    })
    .catch(function () {
      notifications.danger('Error getting associated team');
      eventTag('fetch_associated_team', 'engagement', 'failure');
    });
</script>

{team.name}

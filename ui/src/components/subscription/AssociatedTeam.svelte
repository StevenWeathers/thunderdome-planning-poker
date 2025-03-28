<script lang="ts">
  export let teamId = '';
  export let userId = '';
  export let xfetch = async () => {};
  export let notifications;

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
    });
</script>

{team.name}

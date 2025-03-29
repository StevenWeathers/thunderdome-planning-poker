<script lang="ts">
  export let organizationId = '';
  export let userId = '';
  export let xfetch = async () => {};
  export let notifications;

  let organization = {
    name: '',
  };

  xfetch(`/api/organizations/${organizationId}`)
    .then(res => res.json())
    .then(function (result) {
      organization = result.data.organization;
    })
    .catch(function () {
      notifications.danger('Error getting associated organization');
    });
</script>

{organization.name}

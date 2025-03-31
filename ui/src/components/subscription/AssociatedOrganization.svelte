<script lang="ts">
  interface Props {
    organizationId?: string;
    userId?: string;
    xfetch?: any;
    notifications: any;
  }

  let {
    organizationId = '',
    userId = '',
    xfetch = async () => {},
    notifications
  }: Props = $props();

  let organization = $state({
    name: '',
  });

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

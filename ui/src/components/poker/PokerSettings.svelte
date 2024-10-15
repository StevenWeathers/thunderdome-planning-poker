<script lang="ts">
  import Table from '../table/Table.svelte';
  import TableNav from '../table/TableNav.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import RowCol from '../table/RowCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let organizationId;
  export let teamId;
  export let departmentId;
  export let apiPrefix = '/api';
  export let isEntityAdmin = false;

  $: isAdmin = validateUserIsAdmin($user);

  let defaultSettings = {
    id: '',
    autoFinishVoting: false,
    pointAverageRounding: 0,
    hideVoterIdentity: false,
    estimationScaleId: '',
    joinCode: '',
    facilitatorCode: '',
  };
  let showCreateDefaultSettings = false;
  let showUpdateDefaultSettings = false;

  function toggleCreateDefaultSettings() {
    showCreateDefaultSettings = !showCreateDefaultSettings;
  }

  function toggleUpdateDefaultSettings() {
    showUpdateDefaultSettings = !showUpdateDefaultSettings;
  }

  async function getDefaultSettings() {
    const response = await xfetch(`${apiPrefix}/poker-settings`);
    if (response.ok) {
      defaultSettings = await response.json();
    } else {
      notifications.error('Failed to get default poker settings');
    }
  }

  getDefaultSettings();
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Default Poker Settings"
      createBtnEnabled="{isAdmin || isEntityAdmin}"
      createBtnText="{defaultSettings.id !== ''
        ? 'Update'
        : 'Create'} Default Settings"
      createButtonHandler="{defaultSettings.id !== ''
        ? toggleUpdateDefaultSettings
        : toggleCreateDefaultSettings}"
      createBtnTestId="poker-settings-{defaultSettings.id !== ''
        ? 'update'
        : 'create'}-btn"
    />
    <Table>
      <tr slot="header">
        <HeadCol>AutoFinishVoting</HeadCol>
        <HeadCol>PointAverageRounding</HeadCol>
        <HeadCol>HideVoterIdentity</HeadCol>
        <HeadCol>Estimation Scale ID</HeadCol>
        <HeadCol>Join Code</HeadCol>
        <HeadCol>Facilitator Code</HeadCol>
      </tr>
      <tbody slot="body">
        {#if defaultSettings.id !== ''}
          <TableRow>
            <RowCol>
              {defaultSettings.autoFinishVoting ? 'Yes' : 'No'}
            </RowCol>
            <RowCol>
              {defaultSettings.pointAverageRounding}
            </RowCol>
            <RowCol>
              {defaultSettings.hideVoterIdentity ? 'Yes' : 'No'}
            </RowCol>
            <RowCol>
              {defaultSettings.estimationScaleId}
            </RowCol>
            <RowCol>
              {defaultSettings.joinCode}
            </RowCol>
            <RowCol>
              {defaultSettings.facilitatorCode}
            </RowCol>
          </TableRow>
        {/if}
      </tbody>
    </Table>
  </TableContainer>
</div>

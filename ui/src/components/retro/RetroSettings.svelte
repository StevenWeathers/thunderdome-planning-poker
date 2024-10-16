<script lang="ts">
  import Table from '../table/Table.svelte';
  import TableNav from '../table/TableNav.svelte';
  import HeadCol from '../table/HeadCol.svelte';
  import RowCol from '../table/RowCol.svelte';
  import TableRow from '../table/TableRow.svelte';
  import TableContainer from '../table/TableContainer.svelte';
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import BooleanDisplay from '../global/BooleanDisplay.svelte';

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
    maxVotes: 0,
    allowMultipleVotes: false,
    brainstormVisibility: false,
    phaseTimeLimit: 0,
    phaseAutoAdvance: false,
    allowCumulativeVoting: false,
    templateId: '',
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
    const response = await xfetch(`${apiPrefix}/retro-settings`);
    if (response.ok) {
      defaultSettings = await response.json();
    } else {
      notifications.error('Failed to get default retro settings');
    }
  }

  getDefaultSettings();
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Default Retrospective Settings"
      createBtnEnabled="{isAdmin || isEntityAdmin}"
      createBtnText="{defaultSettings.id !== ''
        ? 'Update'
        : 'Create'} Default Settings"
      createButtonHandler="{defaultSettings.id !== ''
        ? toggleUpdateDefaultSettings
        : toggleCreateDefaultSettings}"
      createBtnTestId="retro-settings-{defaultSettings.id !== ''
        ? 'update'
        : 'create'}-btn"
    />
    <Table>
      <tr slot="header">
        <HeadCol>MaxVotes</HeadCol>
        <HeadCol>AllowMultipleVotes</HeadCol>
        <HeadCol>BrainstormVisibility</HeadCol>
        <HeadCol>PhaseTimeLimit</HeadCol>
        <HeadCol>PhaseAutoAdvance</HeadCol>
        <HeadCol>AllowCumulativeVoting</HeadCol>
        <!--                <HeadCol>TemplateId</HeadCol>-->
        <HeadCol>Join Code</HeadCol>
        <HeadCol>Facilitator Code</HeadCol>
      </tr>
      <tbody slot="body">
        {#if defaultSettings.id !== ''}
          <TableRow>
            <RowCol>
              {defaultSettings.maxVotes}
            </RowCol>
            <RowCol>
              <BooleanDisplay
                boolValue="{defaultSettings.allowMultipleVotes}"
              />
            </RowCol>
            <RowCol>
              <BooleanDisplay
                boolValue="{defaultSettings.brainstormVisibility}"
              />
            </RowCol>
            <RowCol>
              {defaultSettings.phaseTimeLimit}
            </RowCol>
            <RowCol>
              <BooleanDisplay boolValue="{defaultSettings.phaseAutoAdvance}" />
            </RowCol>
            <RowCol>
              <BooleanDisplay
                boolValue="{defaultSettings.allowCumulativeVoting}"
              />
            </RowCol>
            <!--                    <RowCol>-->
            <!--                        {defaultSettings.templateId}-->
            <!--                    </RowCol>-->
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

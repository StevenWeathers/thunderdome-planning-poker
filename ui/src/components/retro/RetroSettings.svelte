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
  import UpdateRetroSettings from './UpdateRetroSettings.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    xfetch: any;
    notifications: any;
    organizationId: any;
    teamId: any;
    departmentId: any;
    apiPrefix?: string;
    isEntityAdmin?: boolean;
  }

  let {
    xfetch,
    notifications,
    organizationId,
    teamId,
    departmentId,
    apiPrefix = '/api',
    isEntityAdmin = false
  }: Props = $props();

  let isAdmin = $derived(validateUserIsAdmin($user));

  let defaultSettings = $state({
    id: '',
    maxVotes: 0,
    brainstormVisibility: '',
    phaseTimeLimit: 0,
    phaseAutoAdvance: false,
    allowCumulativeVoting: false,
    templateId: null,
    joinCode: '',
    facilitatorCode: '',
  });
  let showCreateDefaultSettings = $state(false);
  let showUpdateDefaultSettings = $state(false);

  function toggleCreateDefaultSettings() {
    showCreateDefaultSettings = !showCreateDefaultSettings;
  }

  function toggleUpdateDefaultSettings() {
    showUpdateDefaultSettings = !showUpdateDefaultSettings;
  }

  async function getDefaultSettings() {
    const response = await xfetch(`${apiPrefix}/retro-settings`);
    if (response.ok) {
      const res = await response.json();
      if (!res.data) {
        return;
      }
      defaultSettings = res.data;
    } else {
      notifications.error('Failed to get default retro settings');
    }
  }

  function handleSettingsUpdate(event) {
    defaultSettings = event.detail.settings;
  }

  getDefaultSettings();
</script>

<div class="w-full">
  <TableContainer>
    <TableNav
      title="Default Retrospective Settings"
      createBtnEnabled={isAdmin || isEntityAdmin}
      createBtnText="{defaultSettings.id !== ''
        ? 'Update'
        : 'Create'} Default Settings"
      createButtonHandler={defaultSettings.id !== ''
        ? toggleUpdateDefaultSettings
        : toggleCreateDefaultSettings}
      createBtnTestId="retro-settings-{defaultSettings.id !== ''
        ? 'update'
        : 'create'}-btn"
    />
    <Table>
      {#snippet header()}
            <tr >
          <HeadCol>{$LL.retroMaxVotesPerUserLabel()}</HeadCol>
          <HeadCol>{$LL.brainstormPhaseFeedbackVisibility()}</HeadCol>
          <HeadCol>{$LL.retroPhaseTimeLimitMinLabel()}</HeadCol>
          <HeadCol>{$LL.phaseAutoAdvanceLabel()}</HeadCol>
          <HeadCol>{$LL.allowCumulativeVotingLabel()}</HeadCol>
          <!--                <HeadCol>TemplateId</HeadCol>-->
          <HeadCol>{$LL.joinCodeLabelOptional()}</HeadCol>
          <HeadCol>{$LL.facilitatorCodeOptional()}</HeadCol>
        </tr>
          {/snippet}
      {#snippet body()}
            <tbody >
          {#if defaultSettings.id !== ''}
            <TableRow>
              <RowCol>
                {defaultSettings.maxVotes}
              </RowCol>
              <RowCol>
                {defaultSettings.brainstormVisibility}
              </RowCol>
              <RowCol>
                {defaultSettings.phaseTimeLimit}
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue={defaultSettings.phaseAutoAdvance} />
              </RowCol>
              <RowCol>
                <BooleanDisplay
                  boolValue={defaultSettings.allowCumulativeVoting}
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
          {/snippet}
    </Table>
  </TableContainer>

  {#if showCreateDefaultSettings}
    <UpdateRetroSettings
      toggleClose={toggleCreateDefaultSettings}
      xfetch={xfetch}
      notifications={notifications}
      apiPrefix={apiPrefix}
      organizationId={organizationId}
      teamId={teamId}
      departmentId={departmentId}
      on:updateRetroSettings={handleSettingsUpdate}
    />
  {/if}

  {#if showUpdateDefaultSettings}
    <UpdateRetroSettings
      toggleClose={toggleUpdateDefaultSettings}
      xfetch={xfetch}
      notifications={notifications}
      retroSettings={defaultSettings}
      apiPrefix={apiPrefix}
      organizationId={organizationId}
      teamId={teamId}
      departmentId={departmentId}
      on:updateRetroSettings={handleSettingsUpdate}
    />
  {/if}
</div>

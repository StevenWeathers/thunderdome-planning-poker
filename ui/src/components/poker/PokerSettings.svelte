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
  import UpdatePokerSettings from './UpdatePokerSettings.svelte';
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
    autoFinishVoting: false,
    pointAverageRounding: '',
    hideVoterIdentity: false,
    estimationScaleId: null,
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
    const response = await xfetch(`${apiPrefix}/poker-settings`);
    if (response.ok) {
      const res = await response.json();
      if (!res.data) {
        return;
      }
      defaultSettings = res.data;
    } else {
      notifications.error('Failed to get default poker settings');
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
      {#snippet header()}
            <tr >
          <HeadCol>{$LL.autoFinishVotingLabel()}</HeadCol>
          <HeadCol>{$LL.pointAverageRounding()}</HeadCol>
          <HeadCol>{$LL.hideVoterIdentity()}</HeadCol>
          <!--                <HeadCol>Estimation Scale ID</HeadCol>-->
          <HeadCol>{$LL.joinCodeLabelOptional()}</HeadCol>
          <HeadCol>{$LL.facilitatorCodeOptional()}</HeadCol>
        </tr>
          {/snippet}
      {#snippet body()}
            <tbody >
          {#if defaultSettings.id !== ''}
            <TableRow>
              <RowCol>
                <BooleanDisplay boolValue="{defaultSettings.autoFinishVoting}" />
              </RowCol>
              <RowCol>
                {defaultSettings.pointAverageRounding}
              </RowCol>
              <RowCol>
                <BooleanDisplay boolValue="{defaultSettings.hideVoterIdentity}" />
              </RowCol>
              <!--                    <RowCol>-->
              <!--                        {defaultSettings.estimationScaleId}-->
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
    <UpdatePokerSettings
      toggleClose="{toggleCreateDefaultSettings}"
      xfetch="{xfetch}"
      notifications="{notifications}"
      apiPrefix="{apiPrefix}"
      organizationId="{organizationId}"
      teamId="{teamId}"
      departmentId="{departmentId}"
      on:updatePokerSettings="{handleSettingsUpdate}"
    />
  {/if}

  {#if showUpdateDefaultSettings}
    <UpdatePokerSettings
      toggleClose="{toggleUpdateDefaultSettings}"
      xfetch="{xfetch}"
      notifications="{notifications}"
      pokerSettings="{defaultSettings}"
      apiPrefix="{apiPrefix}"
      organizationId="{organizationId}"
      teamId="{teamId}"
      departmentId="{departmentId}"
      on:updatePokerSettings="{handleSettingsUpdate}"
    />
  {/if}
</div>

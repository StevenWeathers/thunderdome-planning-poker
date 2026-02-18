<script lang="ts">
  import { onMount } from 'svelte';
  import { activeAlerts, user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import CreateAlert from '../../components/alert/CreateAlert.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import Table from '../../components/table/Table.svelte';
  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import BooleanDisplay from '../../components/global/BooleanDisplay.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  interface Alert {
    id: string;
    name: string;
    type: string;
    content: string;
    active: boolean;
    registeredOnly: boolean;
    allowDismiss: boolean;
    updatedDate: string;
  }

  let { xfetch, router, notifications }: Props = $props();

  const alertsPageLimit = 25;
  let alertCount = $state(0);

  const defaultAlert: Alert = {
    id: '',
    name: '',
    type: '',
    content: '',
    active: false,
    registeredOnly: false,
    allowDismiss: true,
    updatedDate: '',
  };

  let alerts = $state<Alert[]>([]);
  let alertsPage = $state(1);
  let showAlertCreate = $state(false);
  let showAlertUpdate = $state(false);
  let showDeleteAlert = $state(false);
  let selectedAlert = $state<Alert>({ ...defaultAlert });
  let deleteAlertId = $state<string | null>(null);

  function toggleCreateAlert() {
    showAlertCreate = !showAlertCreate;
  }

  const toggleUpdateAlert = (alert: Alert) => () => {
    showAlertUpdate = !showAlertUpdate;
    selectedAlert = alert;
  };

  const toggleDeleteAlert = (alertId: string | null) => () => {
    showDeleteAlert = !showDeleteAlert;
    deleteAlertId = alertId;
  };

  function createAlert(body: any) {
    xfetch('/api/alerts', { body })
      .then(res => res.json())
      .then(function (result) {
        activeAlerts.update(result.data);
        getAlerts();
        toggleCreateAlert();
        notifications.success($LL.createAlertSuccess());
      })
      .catch(function () {
        notifications.danger('createAlertError');
      });
  }

  function updateAlert(id: string, body: any) {
    xfetch(`/api/alerts/${id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function (result) {
        activeAlerts.update(result.data);
        getAlerts();
        toggleUpdateAlert({ ...defaultAlert })();
        notifications.success($LL.updateAlertSuccess());
      })
      .catch(function () {
        notifications.danger($LL.updateAlertError());
      });
  }

  function getAlerts() {
    const alertsOffset = (alertsPage - 1) * alertsPageLimit;
    xfetch(`/api/alerts?limit=${alertsPageLimit}&offset=${alertsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        alerts = result.data;
        alertCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getAlertsError());
      });
  }

  function handleDeleteAlert() {
    xfetch(`/api/alerts/${deleteAlertId}`, { method: 'DELETE' })
      .then(res => res.json())
      .then(function (result) {
        activeAlerts.update(result.data);
        getAlerts();
        toggleDeleteAlert(null)();
        notifications.success($LL.deleteAlertSuccess());
      })
      .catch(function () {
        notifications.danger($LL.deleteAlertError());
      });
  }

  const changePage = (evt: CustomEvent) => {
    alertsPage = evt.detail;
    getAlerts();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getAlerts();
  });
</script>

<svelte:head>
  <title>{$LL.alerts()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="alerts">
  <TableContainer>
    <TableNav
      title={$LL.alerts()}
      createBtnText={$LL.alertCreate()}
      createButtonHandler={toggleCreateAlert}
      createBtnTestId="alert-create"
    />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.type()}
          </HeadCol>
          <HeadCol>
            {$LL.active()}
          </HeadCol>
          <HeadCol>
            {$LL.alertRegisteredOnly()}
          </HeadCol>
          <HeadCol>
            {$LL.alertAllowDismiss()}
          </HeadCol>
          <HeadCol>
            {$LL.dateUpdated()}
          </HeadCol>
          <HeadCol type="action">
            <span class="sr-only">Actions</span>
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each alerts as alert, i}
            <TableRow itemIndex={i}>
              <RowCol>
                {alert.name}
              </RowCol>
              <RowCol>
                {alert.type}
              </RowCol>
              <RowCol>
                {#if alert.active}
                  <BooleanDisplay boolValue={alert.active} />
                {/if}
              </RowCol>
              <RowCol>
                {#if alert.registeredOnly}
                  <BooleanDisplay boolValue={alert.registeredOnly} />
                {/if}
              </RowCol>
              <RowCol>
                {#if alert.allowDismiss}
                  <BooleanDisplay boolValue={alert.allowDismiss} />
                {/if}
              </RowCol>
              <RowCol>
                {new Date(alert.updatedDate).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                <CrudActions
                  editBtnClickHandler={toggleUpdateAlert(alert)}
                  deleteBtnClickHandler={toggleDeleteAlert(alert.id)}
                />
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
    <TableFooter bind:current={alertsPage} num_items={alertCount} per_page={alertsPageLimit} on:navigate={changePage} />
  </TableContainer>

  {#if showAlertCreate}
    <CreateAlert toggleCreate={toggleCreateAlert} handleCreate={createAlert} />
  {/if}
  {#if showAlertUpdate}
    <CreateAlert
      toggleUpdate={toggleUpdateAlert({ ...defaultAlert })}
      handleUpdate={updateAlert}
      alertId={selectedAlert.id}
      alertName={selectedAlert.name}
      alertType={selectedAlert.type}
      content={selectedAlert.content}
      active={selectedAlert.active}
      registeredOnly={selectedAlert.registeredOnly}
      allowDismiss={selectedAlert.allowDismiss}
    />
  {/if}

  {#if showDeleteAlert}
    <DeleteConfirmation
      toggleDelete={toggleDeleteAlert(null)}
      handleDelete={handleDeleteAlert}
      confirmText={$LL.alertDeleteConfirmation()}
      confirmBtnText={$LL.alertDelete()}
    />
  {/if}
</AdminPageLayout>

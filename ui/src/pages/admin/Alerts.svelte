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
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import Table from '../../components/table/Table.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;

  const alertsPageLimit = 25;
  let alertCount = 0;

  const defaultAlert = {
    id: '',
    name: '',
    type: '',
    content: '',
    active: '',
    registeredOnly: '',
    allowDismiss: '',
  };

  let alerts = [];
  let alertsPage = 1;
  let showAlertCreate = false;
  let showAlertUpdate = false;
  let showDeleteAlert = false;
  let selectedAlert = { ...defaultAlert };
  let deleteAlertId = null;

  function toggleCreateAlert() {
    showAlertCreate = !showAlertCreate;
  }

  const toggleUpdateAlert = alert => () => {
    showAlertUpdate = !showAlertUpdate;
    selectedAlert = alert;
  };

  const toggleDeleteAlert = alertId => () => {
    showDeleteAlert = !showDeleteAlert;
    deleteAlertId = alertId;
  };

  function createAlert(body) {
    xfetch('/api/alerts', { body })
      .then(res => res.json())
      .then(function (result) {
        eventTag('admin_create_alert', 'engagement', 'success');

        activeAlerts.update(result.data);
        getAlerts();
        toggleCreateAlert();
        notifications.success($LL.createAlertSuccess());
      })
      .catch(function () {
        notifications.danger('createAlertError');
        eventTag('admin_create_alert', 'engagement', 'failure');
      });
  }

  function updateAlert(id, body) {
    xfetch(`/api/alerts/${id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function (result) {
        eventTag('admin_update_alert', 'engagement', 'success');

        activeAlerts.update(result.data);
        getAlerts();
        toggleUpdateAlert({ ...defaultAlert })();
        notifications.success($LL.updateAlertSuccess());
      })
      .catch(function () {
        notifications.danger($LL.updateAlertError());
        eventTag('admin_update_alert', 'engagement', 'failure');
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
        eventTag('admin_delete_alert', 'engagement', 'success');
        activeAlerts.update(result.data);
        getAlerts();
        toggleDeleteAlert(null)();
        notifications.success($LL.deleteAlertSuccess());
      })
      .catch(function () {
        notifications.danger($LL.deleteAlertError());
        eventTag('admin_delete_alert', 'engagement', 'failure');
      });
  }

  const changePage = evt => {
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
      title="{$LL.alerts()}"
      createBtnText="{$LL.alertCreate()}"
      createButtonHandler="{toggleCreateAlert}"
      createBtnTestId="alert-create"
    />
    <Table>
      <tr slot="header">
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
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each alerts as alert, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              {alert.name}
            </RowCol>
            <RowCol>
              {alert.type}
            </RowCol>
            <RowCol>
              {#if alert.active}
                <span class="text-green-600">
                  <CheckIcon />
                </span>
              {/if}
            </RowCol>
            <RowCol>
              {#if alert.registeredOnly}
                <span class="text-green-600">
                  <CheckIcon />
                </span>
              {/if}
            </RowCol>
            <RowCol>
              {#if alert.allowDismiss}
                <span class="text-green-600">
                  <CheckIcon />
                </span>
              {/if}
            </RowCol>
            <RowCol>
              {new Date(alert.updatedDate).toLocaleString()}
            </RowCol>
            <RowCol type="action">
              <CrudActions
                editBtnClickHandler="{toggleUpdateAlert(alert)}"
                deleteBtnClickHandler="{toggleDeleteAlert(alert.id)}"
              />
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>
    <TableFooter
      bind:current="{alertsPage}"
      num_items="{alertCount}"
      per_page="{alertsPageLimit}"
      on:navigate="{changePage}"
    />
  </TableContainer>

  {#if showAlertCreate}
    <CreateAlert
      toggleCreate="{toggleCreateAlert}"
      handleCreate="{createAlert}"
      alerts="{alerts}"
    />
  {/if}
  {#if showAlertUpdate}
    <CreateAlert
      toggleUpdate="{toggleUpdateAlert({ ...defaultAlert })}"
      handleUpdate="{updateAlert}"
      alertId="{selectedAlert.id}"
      alertName="{selectedAlert.name}"
      alertType="{selectedAlert.type}"
      content="{selectedAlert.content}"
      active="{selectedAlert.active}"
      registeredOnly="{selectedAlert.registeredOnly}"
      allowDismiss="{selectedAlert.allowDismiss}"
    />
  {/if}

  {#if showDeleteAlert}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteAlert(null)}"
      handleDelete="{handleDeleteAlert}"
      confirmText="{$LL.alertDeleteConfirmation()}"
      confirmBtnText="{$LL.alertDelete()}"
    />
  {/if}
</AdminPageLayout>

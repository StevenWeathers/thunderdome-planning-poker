<script lang="ts">
  import { onMount } from 'svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { activeAlerts, user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import CreateAlert from '../../components/alert/CreateAlert.svelte';
  import Pagination from '../../components/global/Pagination.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import CheckIcon from '../../components/icons/CheckIcon.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';

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
        notifications.success(`${$LL.createAlertSuccess()}`);
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
        notifications.success(`${$LL.updateAlertSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.updateAlertError()}`);
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
        notifications.danger(`${$LL.getAlertsError()}`);
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
        notifications.success(`${$LL.deleteAlertSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.deleteAlertError()}`);
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
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.alerts()}
    </h1>
  </div>

  <div class="w-full">
    <div class="text-right mb-4">
      <HollowButton onClick="{toggleCreateAlert}">
        {$LL.alertCreate()}
      </HollowButton>
    </div>

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
              <HollowButton onClick="{toggleUpdateAlert(alert)}" color="blue">
                {$LL.edit()}
              </HollowButton>
              <HollowButton onClick="{toggleDeleteAlert(alert.id)}" color="red">
                {$LL.delete()}
              </HollowButton>
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if alertCount > alertsPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{alertsPage}"
          num_items="{alertCount}"
          per_page="{alertsPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}
  </div>

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

<script>
    import { onMount } from 'svelte'

    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CreateAlert from '../../components/alert/CreateAlert.svelte'
    import DeleteConfirmation from '../../components/DeleteConfirmation.svelte'
    import { activeAlerts, warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import RowCol from '../../components/table/RowCol.svelte'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const alertsPageLimit = 25
    let alertCount = 0

    const defaultAlert = {
        id: '',
        name: '',
        type: '',
        content: '',
        active: '',
        registeredOnly: '',
        allowDismiss: '',
    }

    let alerts = []
    let alertsPage = 1
    let showAlertCreate = false
    let showAlertUpdate = false
    let showDeleteAlert = false
    let selectedAlert = { ...defaultAlert }
    let deleteAlertId = null

    function toggleCreateAlert() {
        showAlertCreate = !showAlertCreate
    }

    const toggleUpdateAlert = alert => () => {
        showAlertUpdate = !showAlertUpdate
        selectedAlert = alert
    }

    const toggleDeleteAlert = alertId => () => {
        showDeleteAlert = !showDeleteAlert
        deleteAlertId = alertId
    }

    function createAlert(body) {
        xfetch('/api/alerts', { body })
            .then(res => res.json())
            .then(function (result) {
                eventTag('admin_create_alert', 'engagement', 'success')

                activeAlerts.update(result.data)
                getAlerts()
                toggleCreateAlert()
                notifications.success($_('createAlertSuccess'))
            })
            .catch(function () {
                notifications.danger('createAlertError')
                eventTag('admin_create_alert', 'engagement', 'failure')
            })
    }

    function updateAlert(id, body) {
        xfetch(`/api/alerts/${id}`, { body, method: 'PUT' })
            .then(res => res.json())
            .then(function (result) {
                eventTag('admin_update_alert', 'engagement', 'success')

                activeAlerts.update(result.data)
                getAlerts()
                toggleUpdateAlert({ ...defaultAlert })()
                notifications.success($_('updateAlertSuccess'))
            })
            .catch(function () {
                notifications.danger($_('updateAlertError'))
                eventTag('admin_update_alert', 'engagement', 'failure')
            })
    }

    function getAlerts() {
        const alertsOffset = (alertsPage - 1) * alertsPageLimit
        xfetch(`/api/alerts?limit=${alertsPageLimit}&offset=${alertsOffset}`)
            .then(res => res.json())
            .then(function (result) {
                alerts = result.data
                alertCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getAlertsError'))
            })
    }

    function handleDeleteAlert() {
        xfetch(`/api/alerts/${deleteAlertId}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(function (result) {
                eventTag('admin_delete_alert', 'engagement', 'success')
                activeAlerts.update(result.data)
                getAlerts()
                toggleDeleteAlert(null)()
                notifications.success($_('deleteAlertSuccess'))
            })
            .catch(function () {
                notifications.danger($_('deleteAlertError'))
                eventTag('admin_delete_alert', 'engagement', 'failure')
            })
    }

    const changePage = evt => {
        alertsPage = evt.detail
        getAlerts()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
            return
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
            return
        }

        getAlerts()
    })
</script>

<svelte:head>
    <title>{$_('alerts')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="alerts">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('alerts')}
        </h1>
    </div>

    <div class="w-full">
        <div class="text-right mb-4">
            <HollowButton onClick="{toggleCreateAlert}">
                {$_('alertCreate')}
            </HollowButton>
        </div>

        <Table>
            <tr slot="header">
                <HeadCol>
                    {$_('name')}
                </HeadCol>
                <HeadCol>
                    {$_('type')}
                </HeadCol>
                <HeadCol>
                    {$_('active')}
                </HeadCol>
                <HeadCol>
                    {$_('alertRegisteredOnly')}
                </HeadCol>
                <HeadCol>
                    {$_('alertAllowDismiss')}
                </HeadCol>
                <HeadCol>
                    {$_('dateUpdated')}
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
                            <HollowButton
                                onClick="{toggleUpdateAlert(alert)}"
                                color="blue"
                            >
                                {$_('edit')}
                            </HollowButton>
                            <HollowButton
                                onClick="{toggleDeleteAlert(alert.id)}"
                                color="red"
                            >
                                {$_('delete')}
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
            confirmText="{$_('alertDeleteConfirmation')}"
            confirmBtnText="{$_('alertDelete')}"
        />
    {/if}
</AdminPageLayout>

<script>
    import { onMount } from 'svelte'

    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CreateAlert from '../../components/CreateAlert.svelte'
    import DeleteAlert from '../../components/DeleteAlert.svelte'
    import { activeAlerts, warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils.js'

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
                notifications.success('Alert created successfully.')
            })
            .catch(function () {
                notifications.danger('Error encountered creating alert')
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
                notifications.success('Alert updating successfully.')
            })
            .catch(function () {
                notifications.danger('Error encountered updating alert')
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
                notifications.danger('Error getting alerts')
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
                notifications.success('Alert deleted successfully.')
            })
            .catch(function () {
                notifications.danger('Error attempting to delete alert')
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
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
        }

        getAlerts()
    })
</script>

<svelte:head>
    <title>{$_('alerts')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="alerts">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('alerts')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="text-right">
                <HollowButton onClick="{toggleCreateAlert}">
                    {$_('alertCreate')}
                </HollowButton>
            </div>
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('name')}</th>
                        <th class="flex-1 p-2">{$_('type')}</th>
                        <th class="flex-1 p-2">{$_('active')}</th>
                        <th class="flex-1 p-2">
                            {$_('alertRegisteredOnly')}
                        </th>
                        <th class="flex-1 p-2">
                            {$_('alertAllowDismiss')}
                        </th>
                        <th class="flex-1 p-2">{$_('dateUpdated')}</th>
                        <th class="flex-1 p-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each alerts as alert}
                        <tr>
                            <td class="border p-2">{alert.name}</td>
                            <td class="border p-2">{alert.type}</td>
                            <td class="border p-2">
                                {#if alert.active}
                                    <span class="text-green-600">
                                        <CheckIcon />
                                    </span>
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if alert.registeredOnly}
                                    <span class="text-green-600">
                                        <CheckIcon />
                                    </span>
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if alert.allowDismiss}
                                    <span class="text-green-600">
                                        <CheckIcon />
                                    </span>
                                {/if}
                            </td>
                            <td class="border p-2">
                                {new Date(alert.updatedDate).toLocaleString()}
                            </td>
                            <td class="border p-2 text-right">
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
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

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
        <DeleteAlert
            toggleDelete="{toggleDeleteAlert(null)}"
            handleDelete="{handleDeleteAlert}"
        />
    {/if}
</AdminPageLayout>

<script>
    import CloseIcon from '../icons/CloseIcon.svelte'
    import { activeAlerts, dismissedAlerts } from '../../stores.js'

    export let registered = false

    let alerts = []
    let dismissed = []

    activeAlerts.subscribe(a => {
        alerts = a
    })
    dismissedAlerts.subscribe(d => {
        dismissed = d
    })

    const dismissAlert = alertId => () => {
        dismissedAlerts.dismiss([...alerts], [...$dismissedAlerts, alertId])
    }

    const showAlert = (dismissedAlerts, isRegistered, alert) => {
        const meetsAllowDesmissed = alert.allowDismiss
            ? !dismissedAlerts.includes(alert.id)
            : true
        const meetsRegisteredOnly = alert.registeredOnly ? isRegistered : true

        return meetsAllowDesmissed && meetsRegisteredOnly
    }
</script>

<style>
    .NEWAlert {
        @apply bg-indigo-600;
    }

    .NEWAlert-body {
        @apply text-indigo-100;
    }

    .NEWAlert-type {
        @apply bg-indigo-500;
        @apply text-white;
    }

    .NEWAlert-dismiss {
        @apply text-indigo-100;
    }

    .NEWAlert-dismiss:hover {
        @apply text-white;
        @apply bg-indigo-500;
    }

    .ERRORAlert {
        @apply bg-red-100;
        @apply border-b;
        @apply border-red-200;
    }

    .ERRORAlert-body {
        @apply text-red-800;
    }

    .ERRORAlert-type {
        @apply bg-red-500;
        @apply text-red-100;
    }

    .ERRORAlert-dismiss {
        @apply text-red-700;
    }

    .ERRORAlert-dismiss:hover {
        @apply text-white;
        @apply bg-red-500;
    }

    .INFOAlert {
        @apply bg-blue-100;
    }

    .INFOAlert-body {
        @apply text-blue-800;
    }

    .INFOAlert-type {
        @apply bg-blue-500;
        @apply text-blue-100;
    }

    .INFOAlert-dismiss {
        @apply text-blue-700;
    }

    .INFOAlert-dismiss:hover {
        @apply text-white;
        @apply bg-blue-500;
    }

    .SUCCESSAlert {
        @apply bg-green-100;
        @apply border-b;
        @apply border-green-200;
    }

    .sSUCCESSAlert-body {
        @apply text-green-800;
    }

    .SUCCESSAlert-type {
        @apply bg-green-500;
        @apply text-green-100;
    }

    .SUCCESSAlert-dismiss {
        @apply text-green-700;
    }

    .SUCCESSAlert-dismiss:hover {
        @apply text-white;
        @apply bg-green-500;
    }

    .WARNINGAlert {
        @apply bg-yellow-100;
        @apply border-b;
        @apply border-yellow-200;
    }

    .WARNINGAlert-body {
        @apply text-yellow-800;
    }

    .WARNINGAlert-type {
        @apply bg-yellow-700;
        @apply text-yellow-100;
    }

    .WARNINGAlert-dismiss {
        @apply text-yellow-700;
    }

    .WARNINGAlert-dismiss:hover {
        @apply text-white;
        @apply bg-yellow-700;
    }
</style>

{#each alerts as alert}
    {#if showAlert(dismissed, registered, alert)}
        <div class="{alert.type}Alert">
            <div class="max-w-7xl mx-auto py-3 px-3 sm:px-6 lg:px-8">
                <div
                    class="flex items-center justify-between flex-wrap"
                    role="alert"
                >
                    <div class="w-0 flex-1 flex items-center">
                        <span
                            class="{alert.type}Alert-type flex rounded-lg uppercase
                        px-2 py-1 text-xs font-bold mr-3"
                        >
                            {alert.type}
                        </span>
                        <p
                            class="ml-3 font-medium {alert.type}Alert-body truncate"
                        >
                            {alert.content}
                        </p>
                    </div>
                    {#if alert.allowDismiss}
                        <div class="order-2 flex-shrink-0 sm:order-3 sm:ml-3">
                            <button
                                type="button"
                                on:click="{dismissAlert(alert.id)}"
                                class="-mr-1 flex p-2 rounded-md {alert.type}Alert-dismiss focus:outline-none focus:ring-2 focus:ring-white sm:-mr-2"
                            >
                                <CloseIcon />
                            </button>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    {/if}
{/each}

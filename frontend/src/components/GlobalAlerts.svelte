<script>
    import { activeAlerts, dismissedAlerts } from '../stores.js'

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
        @apply bg-indigo-900;
    }
    .NEWAlert-body {
        @apply bg-indigo-800;
        @apply text-indigo-100;
    }
    .NEWAlert-type {
        @apply bg-indigo-500;
    }
    .NEWAlert-dismiss {
        @apply text-indigo-300;
    }

    .ERRORAlert {
        @apply bg-red-100;
        @apply border-b;
        @apply border-red-200;
    }
    .ERRORAlert-body {
        @apply bg-red-200;
        @apply text-red-800;
    }
    .ERRORAlert-type {
        @apply bg-red-700;
        @apply text-red-100;
    }
    .ERRORAlert-dismiss {
        @apply text-red-700;
    }

    .INFOAlert {
        @apply bg-blue-100;
        @apply border-b;
        @apply border-blue-200;
    }
    .INFOAlert-body {
        @apply bg-blue-200;
        @apply text-blue-800;
    }
    .INFOAlert-type {
        @apply bg-blue-700;
        @apply text-blue-100;
    }
    .INFOAlert-dismiss {
        @apply text-blue-700;
    }

    .SUCCESSAlert {
        @apply bg-green-100;
        @apply border-b;
        @apply border-green-200;
    }
    .sSUCCESSAlert-body {
        @apply bg-green-200;
        @apply text-green-800;
    }
    .SUCCESSAlert-type {
        @apply bg-green-700;
        @apply text-green-100;
    }
    .SUCCESSAlert-dismiss {
        @apply text-green-700;
    }

    .WARNINGAlert {
        @apply bg-yellow-100;
        @apply border-b;
        @apply border-yellow-200;
    }
    .WARNINGAlert-body {
        @apply bg-yellow-200;
        @apply text-yellow-800;
    }
    .WARNINGAlert-type {
        @apply bg-yellow-700;
        @apply text-yellow-100;
    }
    .WARNINGAlert-dismiss {
        @apply text-yellow-700;
    }
</style>

{#each alerts as alert}
    {#if showAlert(dismissed, registered, alert)}
        <div class="{alert.type}Alert text-center py-4 lg:px-4 relative">
            <div
                class="{alert.type}Alert-body p-2 items-center leading-none
                lg:rounded-full flex lg:inline-flex"
                role="alert">
                <span
                    class="{alert.type}Alert-type flex rounded-full uppercase
                    px-2 py-1 text-xs font-bold mr-3">
                    {alert.type}
                </span>
                <span class="font-semibold mr-2 text-left flex-auto">
                    {alert.content}
                </span>
            </div>

            {#if alert.allowDismiss}
                <button
                    class="{alert.type}Alert-dismiss absolute right-0 px-4 py-2"
                    on:click="{dismissAlert(alert.id)}">
                    <svg
                        class="fill-current h-6 w-6"
                        role="button"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 20 20">
                        <title>Close</title>
                        <path
                            d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10
                            11.819l-2.651 3.029a1.2 1.2 0 1
                            1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1
                            1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697
                            1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"></path>
                    </svg>
                </button>
            {/if}
        </div>
    {/if}
{/each}

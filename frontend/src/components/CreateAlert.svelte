<script>
    import CloseIcon from './icons/CloseIcon.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'
    import SolidButton from './SolidButton.svelte'

    export let toggleCreate = () => {}
    export let handleCreate = () => {}
    export let toggleUpdate = () => {}
    export let handleUpdate = () => {}
    export let alertId = ''
    export let alertName = ''
    export let alertType = ''
    export let content = ''
    export let active = true
    export let registeredOnly = false
    export let allowDismiss = true

    const alertTypes = ['ERROR', 'INFO', 'NEW', 'SUCCESS', 'WARNING']

    function toggleClose() {
        if (alertId != '') {
            toggleUpdate()
        } else {
            toggleCreate()
        }
    }

    function onSubmit(e) {
        e.preventDefault()

        const body = {
            name: alertName,
            type: alertType,
            content,
            active,
            registeredOnly,
            allowDismiss,
        }

        if (alertId !== '') {
            handleUpdate(alertId, body)
        } else {
            handleCreate(body)
        }
    }

    $: createDisabled = alertName === '' || alertType === '' || content === ''
</script>

<div class="fixed inset-0 flex items-center z-40">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/3 z-50
        m-8">
        <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6">
            <div class="flex justify-end mb-2">
                <button
                    aria-label="close"
                    on:click="{toggleClose}"
                    class="text-gray-800">
                    <CloseIcon />
                </button>
            </div>

            <form on:submit="{onSubmit}" name="createAlert">
                <div class="mb-4">
                    <label
                        class="block text-gray-700 text-sm font-bold mb-2"
                        for="alertName">
                        Name
                    </label>
                    <input
                        bind:value="{alertName}"
                        placeholder="Enter an alert name"
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        id="alertName"
                        name="alertName"
                        required />
                </div>

                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="alertType">
                        Type
                    </label>
                    <div class="relative">
                        <select
                            name="alertType"
                            bind:value="{alertType}"
                            required
                            class="block appearance-none w-full border-2
                            border-gray-400 text-gray-700 py-3 px-4 pr-8 rounded
                            leading-tight focus:outline-none
                            focus:border-purple-500">
                            <option value="" disabled>
                                Choose an Alert type
                            </option>
                            {#each alertTypes as aType}
                                <option value="{aType}">{aType}</option>
                            {/each}
                        </select>
                        <div
                            class="pointer-events-none absolute inset-y-0
                            right-0 flex items-center px-2 text-gray-700">
                            <DownCarrotIcon />
                        </div>
                    </div>
                </div>

                <div class="mb-4">
                    <label
                        class="block text-gray-700 text-sm font-bold mb-2"
                        for="alertContent">
                        Content
                    </label>
                    <input
                        bind:value="{content}"
                        placeholder="Enter alert content"
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        id="alertContent"
                        name="alertContent"
                        required />
                </div>

                <div class="mb-4">
                    <label class="text-gray-700 text-sm font-bold mb-2">
                        <input
                            type="checkbox"
                            bind:checked="{active}"
                            id="active"
                            name="active" />
                        Active
                    </label>
                </div>
                <div class="mb-4">
                    <label class="text-gray-700 text-sm font-bold mb-2">
                        <input
                            type="checkbox"
                            bind:checked="{registeredOnly}"
                            id="registeredOnly"
                            name="registeredOnly" />
                        Registered Only
                    </label>
                </div>
                <div class="mb-4">
                    <label class="text-gray-700 text-sm font-bold mb-2">
                        <input
                            type="checkbox"
                            bind:checked="{allowDismiss}"
                            id="allowDismiss"
                            name="allowDismiss" />
                        Allow dismiss
                    </label>
                </div>

                <div>
                    <div class="text-right">
                        <SolidButton type="submit" disabled="{createDisabled}">
                            Save
                        </SolidButton>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>

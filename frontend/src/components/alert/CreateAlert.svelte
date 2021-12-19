<script>
    import Modal from '../Modal.svelte'
    import DownCarrotIcon from '../icons/ChevronDown.svelte'
    import SolidButton from '../SolidButton.svelte'
    import { _ } from '../../i18n.js'

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

<Modal closeModal="{toggleClose}">
    <form on:submit="{onSubmit}" name="createAlert">
        <div class="mb-4">
            <label
                class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
                for="alertName"
            >
                {$_('name')}
            </label>
            <input
                bind:value="{alertName}"
                placeholder="{$_('alertNamePlaceholder')}"
                class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="alertName"
                name="alertName"
                required
            />
        </div>

        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="alertType"
            >
                {$_('type')}
            </label>
            <div class="relative">
                <select
                    name="alertType"
                    id="alertType"
                    bind:value="{alertType}"
                    required
                    class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
                >
                    <option value="" disabled>
                        {$_('alertTypePlaceholder')}
                    </option>
                    {#each alertTypes as aType}
                        <option value="{aType}">{aType}</option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 right-0 flex
                    items-center px-2 text-gray-700"
                >
                    <DownCarrotIcon />
                </div>
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 font-bold mb-2 dark:text-gray-400"
                for="alertContent"
            >
                {$_('alertContent')}
            </label>
            <input
                bind:value="{content}"
                placeholder="{$_('alertContentPlaceholder')}"
                class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="alertContent"
                name="alertContent"
                required
            />
        </div>

        <div class="mb-4">
            <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
                <input
                    type="checkbox"
                    bind:checked="{active}"
                    id="active"
                    name="active"
                    class="w-4 h-4 dark:accent-lime-400 mr-1"
                />
                {$_('active')}
            </label>
        </div>
        <div class="mb-4">
            <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
                <input
                    type="checkbox"
                    bind:checked="{registeredOnly}"
                    id="registeredOnly"
                    name="registeredOnly"
                    class="w-4 h-4 dark:accent-lime-400 mr-1"
                />
                {$_('alertRegisteredOnly')}
            </label>
        </div>
        <div class="mb-4">
            <label class="text-gray-700 font-bold mb-2 dark:text-gray-400">
                <input
                    type="checkbox"
                    bind:checked="{allowDismiss}"
                    id="allowDismiss"
                    name="allowDismiss"
                    class="w-4 h-4 dark:accent-lime-400 mr-1"
                />
                {$_('alertAllowDismiss')}
            </label>
        </div>

        <div>
            <div class="text-right">
                <SolidButton type="submit" disabled="{createDisabled}">
                    {$_('alertSave')}
                </SolidButton>
            </div>
        </div>
    </form>
</Modal>

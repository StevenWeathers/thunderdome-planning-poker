<script>
    import Modal from '../Modal.svelte'
    import DownCarrotIcon from '../icons/DownCarrotIcon.svelte'
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
                class="block text-gray-700 text-sm font-bold mb-2"
                for="alertName"
            >
                {$_('name')}
            </label>
            <input
                bind:value="{alertName}"
                placeholder="{$_('alertNamePlaceholder')}"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="alertName"
                name="alertName"
                required
            />
        </div>

        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="alertType">
                {$_('type')}
            </label>
            <div class="relative">
                <select
                    name="alertType"
                    bind:value="{alertType}"
                    required
                    class="block appearance-none w-full border-2 border-gray-300
                    text-gray-700 py-3 px-4 pr-8 rounded leading-tight
                    focus:outline-none focus:border-purple-500"
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
                class="block text-gray-700 text-sm font-bold mb-2"
                for="alertContent"
            >
                {$_('alertContent')}
            </label>
            <input
                bind:value="{content}"
                placeholder="{$_('alertContentPlaceholder')}"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="alertContent"
                name="alertContent"
                required
            />
        </div>

        <div class="mb-4">
            <label class="text-gray-700 text-sm font-bold mb-2">
                <input
                    type="checkbox"
                    bind:checked="{active}"
                    id="active"
                    name="active"
                />
                {$_('active')}
            </label>
        </div>
        <div class="mb-4">
            <label class="text-gray-700 text-sm font-bold mb-2">
                <input
                    type="checkbox"
                    bind:checked="{registeredOnly}"
                    id="registeredOnly"
                    name="registeredOnly"
                />
                {$_('alertRegisteredOnly')}
            </label>
        </div>
        <div class="mb-4">
            <label class="text-gray-700 text-sm font-bold mb-2">
                <input
                    type="checkbox"
                    bind:checked="{allowDismiss}"
                    id="allowDismiss"
                    name="allowDismiss"
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

<script>
    import Modal from './Modal.svelte'
    import SolidButton from './SolidButton.svelte'
    import { _ } from '../i18n'

    export let teamName = ''

    export let toggleCreate = () => {}
    export let handleCreate = () => {}

    function onSubmit(e) {
        e.preventDefault()

        handleCreate(teamName)
    }

    $: createDisabled = teamName === ''
</script>

<Modal closeModal="{toggleCreate}">
    <form on:submit="{onSubmit}" name="createOrganization">
        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="teamName"
            >
                {$_('teamName')}
            </label>
            <input
                bind:value="{teamName}"
                placeholder="{$_('teamNamePlaceholder')}"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="teamName"
                name="teamName"
                required
            />
        </div>

        <div>
            <div class="text-right">
                <SolidButton type="submit" disabled="{createDisabled}">
                    {$_('teamSave')}
                </SolidButton>
            </div>
        </div>
    </form>
</Modal>

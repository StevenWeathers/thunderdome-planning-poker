<script>
    import Modal from './Modal.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'
    import SolidButton from './SolidButton.svelte'
    import { _ } from '../i18n'

    export let toggleAdd = () => {}
    export let handleAdd = () => {}

    const roles = ['ADMIN', 'MEMBER']
    let userEmail = ''
    let role = ''

    function onSubmit(e) {
        e.preventDefault()

        handleAdd(userEmail, role)
    }

    $: createDisabled = userEmail === ''
</script>

<Modal closeModal={toggleAdd}>
    <form on:submit="{onSubmit}" name="teamAddUser">
        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="userEmail">
                {$_('userEmail')}
            </label>
            <input
                bind:value="{userEmail}"
                placeholder="{$_('userEmailPlaceholder')}"
                class="bg-gray-200 border-gray-200 border-2
                appearance-none rounded w-full py-2 px-3 text-gray-700
                leading-tight focus:outline-none focus:bg-white
                focus:border-purple-500"
                id="userEmail"
                name="userEmail"
                required />
        </div>

        <div class="mb-4">
            <label
                class="text-gray-700 text-sm font-bold mb-2"
                for="userRole">
                {$_('role')}
            </label>
            <div class="relative">
                <select
                    bind:value="{role}"
                    class="block appearance-none w-full border-2
                    border-gray-400 text-gray-700 py-3 px-4 pr-8 rounded
                    leading-tight focus:outline-none
                    focus:border-purple-500"
                    id="userRole"
                    name="userRole">
                    <option value="">{$_('rolePlaceholder')}</option>
                    {#each roles as userRole}
                        <option value="{userRole}">{userRole}</option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0
                    right-0 flex items-center px-2 text-gray-700">
                    <DownCarrotIcon />
                </div>
            </div>
        </div>

        <div>
            <div class="text-right">
                <SolidButton type="submit" disabled="{createDisabled}">
                    {$_('userAdd')}
                </SolidButton>
            </div>
        </div>
    </form>
</Modal>
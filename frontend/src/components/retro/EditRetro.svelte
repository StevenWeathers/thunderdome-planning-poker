<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import DownCarrotIcon from '../icons/ChevronDown.svelte'
    import { _ } from '../../i18n.js'

    export let toggleEditRetro = () => {}
    export let handleRetroEdit = () => {}
    export let retroName = ''
    export let joinCode = ''
    export let facilitatorCode = ''
    export let maxVotes = 3
    export let brainstormVisibility = 'visible'

    const brainstormVisibilityOptions = [
        {
            label: $_('brainstormVisibilityLabelVisible'),
            value: 'visible',
        },
        {
            label: $_('brainstormVisibilityLabelConcealed'),
            value: 'concealed',
        },
        {
            label: $_('brainstormVisibilityLabelHidden'),
            value: 'hidden',
        },
    ]

    function saveRetro(e) {
        e.preventDefault()

        const retro = {
            retroName,
            joinCode,
            facilitatorCode,
            maxVotes,
            brainstormVisibility,
        }

        handleRetroEdit(retro)
    }
</script>

<Modal closeModal="{toggleEditRetro}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
    <form on:submit="{saveRetro}" name="createRetro">
        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="retroName"
            >
                {$_('retroName')}
            </label>
            <div class="control">
                <input
                    name="retroName"
                    bind:value="{retroName}"
                    placeholder="{$_('retroNamePlaceholder')}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="retroName"
                    required
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="joinCode"
            >
                {$_('passCode')}
            </label>
            <div class="control">
                <input
                    name="joinCode"
                    bind:value="{joinCode}"
                    placeholder="{$_('optionalPasscodePlaceholder')}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="joinCode"
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                for="facilitatorCode"
            >
                {$_('facilitatorCodeOptional')}
            </label>
            <div class="control">
                <input
                    name="facilitatorCode"
                    bind:value="{facilitatorCode}"
                    placeholder="{$_('facilitatorCodePlaceholder')}"
                    class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="facilitatorCode"
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="maxVotes"
            >
                {$_('retroMaxVotesPerUserLabel')}
            </label>
            <div class="control">
                <input
                    name="retroName"
                    bind:value="{maxVotes}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="maxVotes"
                    type="number"
                    min="1"
                    max="10"
                    required
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="brainstormVisibility"
            >
                {$_('brainstormPhaseFeedbackVisibility')}
            </label>
            <div class="relative">
                <select
                    bind:value="{brainstormVisibility}"
                    class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
                    id="brainstormVisibility"
                    name="brainstormVisibility"
                >
                    {#each brainstormVisibilityOptions as item}
                        <option value="{item.value}">
                            {item.label}
                        </option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 right-0 flex
                items-center px-2 text-gray-700 dark:text-gray-400"
                >
                    <DownCarrotIcon />
                </div>
            </div>
        </div>

        <div class="text-right">
            <SolidButton type="submit">{$_('save')}</SolidButton>
        </div>
    </form>
</Modal>

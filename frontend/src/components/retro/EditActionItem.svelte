<script>
    import CheckboxIcon from '../icons/CheckboxIcon.svelte'
    import SolidButton from '../SolidButton.svelte'
    import HollowButton from '../HollowButton.svelte'
    import Modal from '../Modal.svelte'
    import { _ } from '../../i18n.js'

    export let toggleEdit = () => {}
    export let handleEdit = () => {}
    export let handleDelete = () => {}
    export let action = {
        content: '',
        completed: false,
    }

    let editAction = { ...action }

    const handleSubmit = e => {
        e.preventDefault()

        handleEdit(editAction)
    }
</script>

<Modal closeModal="{toggleEdit}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
    <form on:submit="{handleSubmit}">
        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="actionItem"
            >
                {$_('actionItem')}
            </label>
            <div class="control">
                <input
                    bind:value="{editAction.content}"
                    placeholder="{$_('actionItemPlaceholder')}"
                    class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                            px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                            focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
                    id="actionItem"
                    name="actionItem"
                    type="text"
                    required
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="Completed"
            >
                {$_('completed')}
            </label>
            <div class="control">
                <div class="flex-shrink">
                    <input
                        type="checkbox"
                        id="Completed"
                        bind:checked="{editAction.completed}"
                        class="opacity-0 absolute h-6 w-6"
                    />
                    <div
                        class="bg-white dark:bg-gray-800 border-2 rounded-md
                                                border-gray-400 dark:border-gray-300 w-6 h-6 flex flex-shrink-0
                                                justify-center items-center mr-2
                                                focus-within:border-blue-500 dark:focus-within:border-sky-500"
                    >
                        <CheckboxIcon />
                    </div>
                    <label for="Completed" class="select-none"></label>
                </div>
            </div>
        </div>

        <div class="flex w-full pt-4">
            <div class="w-1/2">
                <HollowButton color="red" onClick="{handleDelete(editAction)}"
                    >{$_('delete')}</HollowButton
                >
            </div>
            <div class="w-1/2 text-right">
                <SolidButton type="submit">{$_('save')}</SolidButton>
            </div>
        </div>
    </form>
</Modal>

<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import HollowButton from '../HollowButton.svelte'

    export let toggleColumnEdit = () => {}
    export let handleColumnRevision = () => {}
    export let deleteColumn = () => () => {}

    export let column = {
        id: '',
        name: '',
    }

    function handleSubmit(event) {
        event.preventDefault()

        handleColumnRevision(column)
        toggleColumnEdit()
    }
</script>

<Modal closeModal="{toggleColumnEdit}">
    <form on:submit="{handleSubmit}" name="addColumn">
        <div class="mb-4">
            <label
                class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                for="columnName"
            >
                Column Name
            </label>
            <input
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="columnName"
                type="text"
                bind:value="{column.name}"
                placeholder="Enter a column name"
                name="columnName"
            />
        </div>
        <div class="flex">
            <div class="md:w-1/2 text-left">
                <HollowButton color="red" onClick="{deleteColumn(column.id)}">
                    Delete Column
                </HollowButton>
            </div>
            <div class="md:w-1/2 text-right">
                <SolidButton type="submit">Save</SolidButton>
            </div>
        </div>
    </form>
</Modal>

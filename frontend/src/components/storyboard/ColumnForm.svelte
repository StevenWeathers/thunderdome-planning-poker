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

  function handleSubmit (event) {
    event.preventDefault()

    handleColumnRevision(column)
    toggleColumnEdit()
  }
</script>

<Modal closeModal="{toggleColumnEdit}">
    <form on:submit="{handleSubmit}" name="addColumn">
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="columnName">
                Column Name
            </label>
            <input
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                    id="columnName"
                    type="text"
                    bind:value="{column.name}"
                    placeholder="Enter a column name"
                    name="columnName"/>
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

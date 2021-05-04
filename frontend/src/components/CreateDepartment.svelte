<script>
    import CloseIcon from './icons/CloseIcon.svelte'
    import SolidButton from './SolidButton.svelte'

    export let toggleCreate = () => {}
    export let handleCreate = () => {}

    export let departmentName = ''

    function onSubmit(e) {
        e.preventDefault()

        handleCreate(departmentName)
    }

    $: createDisabled = departmentName === ''
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
                    on:click="{toggleCreate}"
                    class="text-gray-800">
                    <CloseIcon />
                </button>
            </div>

            <form on:submit="{onSubmit}" name="createDepartment">
                <div class="mb-4">
                    <label
                        class="block text-gray-700 text-sm font-bold mb-2"
                        for="departmentName">
                        Department Name
                    </label>
                    <input
                        bind:value="{departmentName}"
                        placeholder="Enter an department name"
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        id="departmentName"
                        name="departmentName"
                        required />
                </div>

                <div>
                    <div class="text-right">
                        <SolidButton type="submit" disabled="{createDisabled}">
                            Create Department
                        </SolidButton>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>

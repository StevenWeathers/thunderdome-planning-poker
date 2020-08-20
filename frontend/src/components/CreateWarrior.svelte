<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'

    export let toggleCreate = () => {}

    function handleSubmit(event) {
        event.preventDefault()
    }
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

        </div>
    </div>
</div>

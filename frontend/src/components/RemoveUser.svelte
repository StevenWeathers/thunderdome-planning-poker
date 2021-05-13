<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import { _ } from '../i18n'

    export let handleRemove = () => {}
    export let toggleRemove = () => {}
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/3 z-50
        max-h-full">
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{toggleRemove}"
                        class="text-gray-800">
                        <CloseIcon />
                    </button>
                </div>

                <div class="mb-4">
                    <p class="font-bold text-xl text-red-600">
                        Are you sure you want to remove this user?
                    </p>
                </div>
                <div class="text-right">
                    <button
                        type="button"
                        class="inline-block align-baseline font-bold text-sm
                        text-blue-500 hover:text-blue-800 mr-4"
                        on:click="{toggleRemove}">
                        Cancel
                    </button>
                    <SolidButton onClick="{handleRemove}" color="red">
                        Remove User
                    </SolidButton>
                </div>
            </div>
        </div>
    </div>
</div>

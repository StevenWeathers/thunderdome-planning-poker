<script>
    import CloseIcon from './icons/CloseIcon.svelte'

    export let closeModal = () => {}
    export let widthClasses = ''
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll"
>
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full {widthClasses != ''
            ? widthClasses
            : 'md:w-2/3 lg:w-3/5 xl:w-1/3'}
        z-50 max-h-full"
    >
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{closeModal}"
                        class="text-gray-800"
                    >
                        <CloseIcon />
                    </button>
                </div>

                <slot />
            </div>
        </div>
    </div>
</div>

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
            <div
                class="relative shadow-xl bg-white dark:bg-gray-800 rounded-lg p-4 xl:p-6 max-h-full"
            >
                <div class="absolute top-4 right-4">
                    <button
                        aria-label="close"
                        on:click="{closeModal}"
                        class="text-gray-800 dark:text-gray-400 hover:bg-red-500 hover:text-white p-1 rounded-full"
                    >
                        <CloseIcon class="w-7 h-7" />
                    </button>
                </div>

                <div class="mt-2">
                    <slot />
                </div>
            </div>
        </div>
    </div>
</div>

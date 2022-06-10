<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'

    export let handleLegendRevision = () => {}
    export let toggleEditLegend = () => {}

    export let colorLegend = []

    function handleSubmit(event) {
        event.preventDefault()

        handleLegendRevision(colorLegend)
        toggleEditLegend()
    }
</script>

<style>
    .colorcard-gray {
        @apply bg-gray-400;
    }

    .colorcard-red {
        @apply bg-red-400;
    }

    .colorcard-orange {
        @apply bg-orange-400;
    }

    .colorcard-yellow {
        @apply bg-yellow-400;
    }

    .colorcard-green {
        @apply bg-green-400;
    }

    .colorcard-teal {
        @apply bg-teal-400;
    }

    .colorcard-blue {
        @apply bg-blue-400;
    }

    .colorcard-indigo {
        @apply bg-indigo-400;
    }

    .colorcard-purple {
        @apply bg-purple-400;
    }

    .colorcard-pink {
        @apply bg-pink-400;
    }
</style>

<Modal closeModal="{toggleEditLegend}">
    <form on:submit="{handleSubmit}" name="colorLegend">
        <div class="mt-8 mb-4">
            {#each colorLegend as color, i}
                <div class="mb-1 flex">
                    <span class="p-4 inline-block colorcard-{color.color}"
                    ></span>
                    <input
                        class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        type="text"
                        bind:value="{colorLegend[i].legend}"
                        placeholder="Enter a color legend"
                        name="legend-{color.color}"
                    />
                </div>
            {/each}
        </div>
        <div class="text-right">
            <div>
                <SolidButton type="submit">Save</SolidButton>
            </div>
        </div>
    </form>
</Modal>

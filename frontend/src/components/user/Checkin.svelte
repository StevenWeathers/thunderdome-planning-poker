<script>
    import Modal from '../Modal.svelte'
    import SolidButton from '../SolidButton.svelte'

    export let toggleCheckin = () => {}
    export let handleCheckin = () => {}
    export let userId

    let today = ''
    let yesterday = ''
    let blockers = ''
    let discuss = ''
    let goalsMet = true

    function onSubmit(e) {
        e.preventDefault()

        handleCheckin({
            userId,
            yesterday,
            today,
            blockers,
            discuss,
            goalsMet,
        })
    }
</script>

<Modal closeModal="{toggleCheckin}">
    <form on:submit="{onSubmit}" name="teamCheckin">
        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="yesterday"
            >
                Yesterday
            </label>
            <textarea
                bind:value="{yesterday}"
                placeholder="Yesterday I..."
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="yesterday"
                name="yesterday"></textarea>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="today"
            >
                Today
            </label>
            <textarea
                bind:value="{today}"
                placeholder="Yesterday I..."
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="today"
                name="today"></textarea>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="blockers"
            >
                Blockers
            </label>
            <textarea
                bind:value="{blockers}"
                placeholder="Yesterday I..."
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="blockers"
                name="blockers"></textarea>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="discuss"
            >
                Discuss
            </label>
            <textarea
                bind:value="{discuss}"
                placeholder="Yesterday I..."
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="discuss"
                name="discuss"></textarea>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="discuss"
            >
                Goals Met
            </label>
            <div
                class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in"
            >
                <input
                    type="checkbox"
                    name="activeBattles"
                    id="activeBattles"
                    bind:checked="{goalsMet}"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                    for="activeBattles"
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
            </div>
        </div>

        <div>
            <div class="text-right">
                <SolidButton type="submit">Save</SolidButton>
            </div>
        </div>
    </form>
</Modal>

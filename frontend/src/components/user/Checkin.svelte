<script>
    import { quill } from '../../quill'
    import Modal from '../Modal.svelte'
    import SolidButton from '../SolidButton.svelte'

    export let toggleCheckin = () => {}
    export let handleCheckin = () => {}
    export let handleCheckinEdit = () => {}
    export let userId
    export let checkinId
    export let today = ''
    export let yesterday = ''
    export let blockers = ''
    export let discuss = ''
    export let goalsMet = true

    function onSubmit(e) {
        e.preventDefault()

        if (checkinId) {
            handleCheckinEdit(checkinId, {
                yesterday,
                today,
                blockers,
                discuss,
                goalsMet,
            })
        } else {
            handleCheckin({
                userId,
                yesterday,
                today,
                blockers,
                discuss,
                goalsMet,
            })
        }
    }
</script>

<Modal closeModal="{toggleCheckin}" widthClasses="md:w-2/3">
    <form on:submit="{onSubmit}" name="teamCheckin">
        <div class="mb-4">
            <div class="text-blue-500 uppercase font-rajdhani text-xl mb-2">
                Yesterday
            </div>
            <div
                class="w-full"
                use:quill="{{
                    placeholder: `Yesterday I...`,
                    content: yesterday,
                }}"
                on:text-change="{e => (yesterday = e.detail.html)}"
                id="yesterday"
            ></div>
        </div>

        <div class="mb-4">
            <div class="text-green-500 uppercase font-rajdhani text-xl mb-2">
                Today
            </div>
            <div
                class="w-full"
                use:quill="{{
                    placeholder: `Today I will...`,
                    content: today,
                }}"
                on:text-change="{e => (today = e.detail.html)}"
                id="today"
            ></div>
        </div>

        <div class="mb-4">
            <div class="text-red-500 uppercase font-rajdhani text-xl mb-2">
                Blockers
            </div>
            <div
                class="w-full"
                use:quill="{{
                    placeholder: `I'm blocked by...`,
                    content: blockers,
                }}"
                on:text-change="{e => (blockers = e.detail.html)}"
                id="blockers"
            ></div>
        </div>

        <div class="mb-4">
            <div class="text-purple-500 uppercase font-rajdhani text-xl mb-2">
                Discuss
            </div>
            <div
                class="w-full"
                use:quill="{{
                    placeholder: 'I would like to discuss...',
                    content: discuss,
                }}"
                on:text-change="{e => (discuss = e.detail.html)}"
                id="discuss"
            ></div>
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

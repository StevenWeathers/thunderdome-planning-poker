<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import { _ } from '../../i18n.js'

    export let toggleBecomeLeader = () => {}
    export let handleBecomeLeader = () => {}

    let leaderCode = ''

    function handleSubmit(e) {
        e.preventDefault()

        handleBecomeLeader(leaderCode)
    }
</script>

<Modal
    closeModal="{toggleBecomeLeader}"
    widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
    <form on:submit="{handleSubmit}" name="becomeLeader">
        <div class="mb-4">
            <label
                class="block text-gray-700 text-sm font-bold mb-2"
                for="leaderCode"
            >
                {$_('leaderPasscode')}
            </label>
            <div class="control">
                <input
                    name="leaderCode"
                    bind:value="{leaderCode}"
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                    rounded w-full py-2 px-3 text-gray-700 leading-tight
                    focus:outline-none focus:bg-white focus:border-purple-500"
                    id="leaderCode"
                />
            </div>
        </div>

        <div class="text-right">
            <SolidButton type="submit">{$_('save')}</SolidButton>
        </div>
    </form>
</Modal>

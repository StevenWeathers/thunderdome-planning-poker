<script>
    import Modal from '../Modal.svelte'
    import SolidButton from '../SolidButton.svelte'
    import { warrior as user } from '../../stores.js'
    import { _ } from '../../i18n.js'

    export let toggleForm = () => {}
    export let handleComment = () => {}
    export let checkinId

    let comment = ''

    function onSubmit(e) {
        e.preventDefault()

        handleComment({
            userId: $user.id,
            comment,
        })
    }
</script>

<Modal closeModal="{toggleForm}" widthClasses="md:w-2/3">
    <form on:submit="{onSubmit}" name="teamCheckin">
        <div class="mt-12 mb-2 w-full">
            <textarea
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                placeholder="{$_('writeCommentPlaceholder')}"
                bind:value="{comment}"></textarea>
        </div>

        <div>
            <div class="text-right">
                <SolidButton type="submit">{$_('postComment')}</SolidButton>
            </div>
        </div>
    </form>
</Modal>

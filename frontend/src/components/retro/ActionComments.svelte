<script>
    import Modal from '../Modal.svelte'
    import HollowButton from '../HollowButton.svelte'
    import UserIcon from '../icons/UserIcon.svelte'
    import { _ } from '../../i18n'
    import { warrior as user } from '../../stores.js'

    export let xfetch
    export let eventTag
    export let notifications
    export let toggleComments = () => {}
    export let getRetrosActions = () => {}
    export let actions = []
    export let users = []
    export let selectedActionId
    export let isAdmin = false

    const userMap = users.reduce((prev, cur) => {
        prev[cur.id] = cur.name
        return prev
    }, {})

    let userComment = ''
    let selectedComment = null
    let selectedCommentContent = ''

    const toggleCommentEdit = comment => () => {
        selectedComment = comment
        if (comment !== null) {
            selectedCommentContent = comment.comment
        }
    }

    function handleCommentSubmit() {
        const body = {
            comment: userComment,
        }

        xfetch(
            `/api/retros/${selectedAction.retroId}/actions/${selectedAction.id}/comments`,
            { body },
        )
            .then(res => res.json())
            .then(function ({ data }) {
                userComment = ''
                getRetrosActions()
                eventTag('retro_comment_add', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('retroActionCommentAddError'))
                eventTag('retro_comment_add', 'engagement', 'failure')
            })
    }

    const handleCommentDelete = commentId => () => {
        xfetch(
            `/api/retros/${selectedAction.retroId}/actions/${selectedAction.id}/comments/${commentId}`,
            { method: 'DELETE' },
        )
            .then(res => res.json())
            .then(function ({ data }) {
                getRetrosActions()
                eventTag('retro_comment_delete', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('retroActionCommentDeleteError'))
                eventTag('retro_comment_delete', 'engagement', 'failure')
            })
    }

    const handleCommentEdit = () => {
        const body = {
            comment: selectedCommentContent,
        }

        xfetch(
            `/api/retros/${selectedAction.retroId}/actions/${selectedAction.id}/comments/${selectedComment.id}`,
            { body, method: 'PUT' },
        )
            .then(res => res.json())
            .then(function ({ data }) {
                selectedComment = null
                selectedCommentContent = ''
                getRetrosActions()
                eventTag('retro_comment_edit', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('retroActionCommentAddError'))
                eventTag('retro_comment_edit', 'engagement', 'failure')
            })
    }

    $: selectedAction = actions && actions.find(a => a.id === selectedActionId)
</script>

<Modal closeModal="{toggleComments}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
    <div class="mt-12 dark:text-gray-300">
        <h3
            class="text-xl pb-2 mb-4 border-b border-gray-600 dark:border-gray-400"
        >
            {$_('actionComments')}
        </h3>
        {#each selectedAction.comments as comment}
            <div
                class="w-full mb-4 text-gray-700 dark:text-gray-400 border-b border-gray-300 dark:border-gray-700"
                data-commentid="{comment.id}"
            >
                <div class="font-bold">
                    <UserIcon class="h-4 w-4" />&nbsp;{userMap[
                        comment.user_id
                    ] || '...'}
                </div>
                {#if selectedComment !== null && selectedComment.id === comment.id}
                    <div class="w-full my-2">
                        <textarea
                            class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                            bind:value="{selectedCommentContent}"></textarea>
                        <div class="text-right">
                            <HollowButton
                                color="blue"
                                onClick="{toggleCommentEdit(null)}"
                            >
                                {$_('cancel')}
                            </HollowButton>
                            <HollowButton
                                color="green"
                                onClick="{handleCommentEdit}"
                                disabled="{selectedCommentContent === ''}"
                            >
                                {$_('updateComment')}
                            </HollowButton>
                        </div>
                    </div>
                {:else}
                    <div class="py-2">
                        {comment.comment}
                    </div>
                {/if}
                {#if (comment.user_id === $user.id || comment.user_id === isAdmin) && !(selectedComment !== null && selectedComment.id === comment.id)}
                    <div class="mb-2 text-right">
                        <button
                            class="text-blue-500 hover:text-blue-300 mr-1"
                            on:click="{toggleCommentEdit(comment)}"
                        >
                            {$_('edit')}
                        </button>
                        <button
                            class="text-red-500"
                            on:click="{handleCommentDelete(comment.id)}"
                        >
                            {$_('delete')}
                        </button>
                    </div>
                {/if}
            </div>
        {/each}
        {#if selectedAction.comments.length === 0}
            <p class="text-lg dark:text-gray-400">{$_('noComments')}</p>
        {/if}

        <div class="w-full mt-8">
            <textarea
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                placeholder="{$_('writeCommentPlaceholder')}"
                bind:value="{userComment}"></textarea>
            <div class="text-right">
                <HollowButton
                    color="teal"
                    onClick="{handleCommentSubmit}"
                    disabled="{userComment === ''}"
                >
                    {$_('postComment')}
                </HollowButton>
            </div>
        </div>
    </div>
</Modal>

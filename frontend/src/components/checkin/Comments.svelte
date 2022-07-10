<script>
    import CommentIcon from '../icons/CommentIcon.svelte'
    import SolidButton from '../SolidButton.svelte'
    import Comment from './Comment.svelte'
    import { _ } from '../../i18n.js'
    import { warrior as user } from '../../stores.js'

    export let checkin = {}
    export let userMap = {}
    export let isAdmin = false
    export let handleCreate = () => {}
    export let handleEdit = () => {}
    export let handleDelete = () => {}

    let showComments = false
    let comment = ''

    function toggleComments() {
        showComments = !showComments
    }

    function onSubmit(e) {
        e.preventDefault()

        handleCreate(checkin.id, {
            userId: $user.id,
            comment,
        })
        comment = ''
    }
</script>

<button class="text-blue-500 dark:text-sky-400" on:click="{toggleComments}">
    <CommentIcon />&nbsp;{checkin.comments.length}
    {checkin.comments.length === 1 ? 'Comment' : $_('comments')}
</button>
{#if showComments}
    <div class="mt-2">
        {#each checkin.comments as comment}
            <Comment
                checkinId="{checkin.id}"
                comment="{comment}"
                userMap="{userMap}"
                isAdmin="{isAdmin}"
                handleEdit="{handleEdit}"
                handleDelete="{handleDelete}"
            />
        {/each}
    </div>
    <div class="text-right mb-2">
        <form on:submit="{onSubmit}" name="checkinComment">
            <div class="mb-2 w-full">
                <textarea
                    class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    placeholder="{$_('writeCommentPlaceholder')}"
                    bind:value="{comment}"></textarea>
            </div>

            <div>
                <div class="text-right">
                    <SolidButton type="submit">
                        {$_('postComment')}
                    </SolidButton>
                </div>
            </div>
        </form>
    </div>
{/if}

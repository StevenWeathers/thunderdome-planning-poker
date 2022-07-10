<script>
    import UserIcon from '../icons/UserIcon.svelte'
    import HollowButton from '../HollowButton.svelte'
    import SolidButton from '../SolidButton.svelte'
    import { _ } from '../../i18n.js'
    import { warrior as user } from '../../stores.js'

    export let checkinId = {}
    export let comment = {}
    export let userMap = {}
    export let isAdmin = false
    export let handleEdit = () => {}
    export let handleDelete = () => {}

    let showEdit = false
    let editcomment = `${comment.comment}`

    function toggleEdit() {
        showEdit = !showEdit
    }

    function onSubmit(e) {
        e.preventDefault()

        handleEdit(checkinId, comment.id, {
            userId: $user.id,
            comment: editcomment,
        })

        toggleEdit()
    }
</script>

<div
    class="w-full mb-2 text-gray-700 dark:text-gray-300 border-b border-gray-300 dark:border-gray-700"
    data-commentid="{comment.id}"
>
    <div class="font-bold">
        <UserIcon class="h-4 w-4" />&nbsp;{userMap[comment.user_id] || '...'}
    </div>
    {#if showEdit}
        <div class="w-full my-2">
            <form on:submit="{onSubmit}" name="checkinComment">
                <textarea
                    class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
    rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
    focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                    bind:value="{editcomment}"></textarea>
                <div class="text-right">
                    <HollowButton color="blue" onClick="{toggleEdit}">
                        {$_('cancel')}
                    </HollowButton>
                    <SolidButton type="submit" disabled="{editcomment === ''}">
                        {$_('updateComment')}
                    </SolidButton>
                </div>
            </form>
        </div>
    {:else}
        <div class="py-2">
            {comment.comment}
        </div>
    {/if}
    {#if (comment.user_id === $user.id || comment.user_id === isAdmin) && !showEdit}
        <div class="mb-2 text-right">
            <button
                class="text-blue-500 hover:text-blue-300 dark:text-sky-300 dark:hover:text-sky-100 mr-1"
                on:click="{toggleEdit}"
            >
                {$_('edit')}
            </button>
            <button
                class="text-red-500"
                on:click="{handleDelete(checkinId, comment.id)}"
            >
                {$_('delete')}
            </button>
        </div>
    {/if}
</div>

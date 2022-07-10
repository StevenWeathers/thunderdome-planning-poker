<script>
    import { quill } from '../../quill.js'
    import Modal from '../Modal.svelte'
    import HollowButton from '../HollowButton.svelte'
    import UserIcon from '../icons/UserIcon.svelte'
    import { warrior as user } from '../../stores.js'
    import { _ } from '../../i18n.js'

    export let toggleStoryForm = () => {}
    export let sendSocketEvent = () => {}
    export let eventTag = () => {}
    export let notifications

    export let story = {}
    export let colorLegend = []
    export let users = []

    const isAbsolute = new RegExp('^([a-z]+://|//)', 'i')

    let userComment = ''
    let selectedComment = null
    let selectedCommentContent = ''

    $: userMap = users.reduce((prev, usr) => {
        prev[usr.id] = usr.name
        return prev
    }, {})

    function handleStoryDelete() {
        sendSocketEvent('delete_story', story.id)
        eventTag('story_delete', 'storyboard', '')
        toggleStoryForm()
    }

    function markClosed() {
        sendSocketEvent(
            'update_story_closed',
            JSON.stringify({
                storyId: story.id,
                closed: true,
            }),
        )
        eventTag('story_edit_closed', 'storyboard', 'true')
    }

    function markOpen() {
        sendSocketEvent(
            'update_story_closed',
            JSON.stringify({
                storyId: story.id,
                closed: false,
            }),
        )
        eventTag('story_edit_closed', 'storyboard', 'false')
    }

    const changeColor = color => () => {
        sendSocketEvent(
            'update_story_color',
            JSON.stringify({
                storyId: story.id,
                color,
            }),
        )
        eventTag('story_edit_color', 'storyboard', color)
    }

    const updateName = evt => {
        const name = evt.target.value
        sendSocketEvent(
            'update_story_name',
            JSON.stringify({
                storyId: story.id,
                name,
            }),
        )
        eventTag('story_edit_name', 'storyboard', '')
    }

    const updateContent = () => {
        sendSocketEvent(
            'update_story_content',
            JSON.stringify({
                storyId: story.id,
                content: story.content,
            }),
        )
        eventTag('story_edit_content', 'storyboard', '')
    }

    const updatePoints = evt => {
        const points = parseInt(evt.target.value, 10)
        sendSocketEvent(
            'update_story_points',
            JSON.stringify({
                storyId: story.id,
                points,
            }),
        )
        eventTag('story_edit_points', 'storyboard', '')
    }

    const updateLink = evt => {
        const link = evt.target.value
        if (link !== '' && !isAbsolute.test(link)) {
            notifications.danger('Link must be an absolute URL')
            return
        }

        sendSocketEvent(
            'update_story_link',
            JSON.stringify({
                storyId: story.id,
                link,
            }),
        )
        eventTag('story_edit_link', 'storyboard', '')
    }

    const handleCommentSubmit = () => {
        if (userComment !== '') {
            sendSocketEvent(
                'add_story_comment',
                JSON.stringify({ storyId: story.id, comment: userComment }),
            )
            eventTag('story_add_comment', 'storyboard', '')
            userComment = ''
        }
    }

    const toggleCommentEdit = comment => () => {
        selectedComment = comment
        if (comment !== null) {
            selectedCommentContent = comment.comment
        }
    }

    const handleCommentEdit = () => {
        sendSocketEvent(
            'edit_story_comment',
            JSON.stringify({
                commentId: selectedComment.id,
                comment: selectedCommentContent,
            }),
        )
        selectedComment = null
        selectedCommentContent = ''
        eventTag('story_edit_comment', 'storyboard', '')
    }

    const handleCommentDelete = commentId => () => {
        sendSocketEvent('delete_story_comment', JSON.stringify({ commentId }))
        eventTag('story_delete_comment', 'storyboard', '')
    }
</script>

<style>
    .colorcard-gray {
        @apply bg-gray-400;
    }

    .colorcard-gray:hover {
        @apply bg-gray-800;
    }

    .colorcard-red {
        @apply bg-red-400;
    }

    .colorcard-red:hover {
        @apply bg-red-800;
    }

    .colorcard-orange {
        @apply bg-orange-400;
    }

    .colorcard-orange:hover {
        @apply bg-orange-800;
    }

    .colorcard-yellow {
        @apply bg-yellow-400;
    }

    .colorcard-yellow:hover {
        @apply bg-yellow-800;
    }

    .colorcard-green {
        @apply bg-green-400;
    }

    .colorcard-green:hover {
        @apply bg-green-800;
    }

    .colorcard-teal {
        @apply bg-teal-400;
    }

    .colorcard-teal:hover {
        @apply bg-teal-800;
    }

    .colorcard-blue {
        @apply bg-blue-400;
    }

    .colorcard-blue:hover {
        @apply bg-blue-800;
    }

    .colorcard-indigo {
        @apply bg-indigo-400;
    }

    .colorcard-indigo:hover {
        @apply bg-indigo-800;
    }

    .colorcard-purple {
        @apply bg-purple-400;
    }

    .colorcard-purple:hover {
        @apply bg-purple-800;
    }

    .colorcard-pink {
        @apply bg-pink-400;
    }

    .colorcard-pink:hover {
        @apply bg-pink-800;
    }
</style>

<Modal closeModal="{toggleStoryForm}" widthClasses="w-full md:w-2/3">
    <div class="flex w-full mt-8">
        <div class="w-3/4">
            <div class="mx-4">
                <div class="mb-4">
                    <label
                        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="storyName"
                    >
                        Story Name
                    </label>
                    <input
                        class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        id="storyName"
                        type="text"
                        on:change="{updateName}"
                        value="{story.name}"
                        placeholder="Enter a story name e.g. Ricky Bobby"
                        name="storyName"
                    />
                </div>
                <div class="mb-4">
                    <label
                        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="storyLink"
                    >
                        Story Link
                    </label>
                    <input
                        class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        id="storyLink"
                        type="text"
                        on:change="{updateLink}"
                        value="{story.link}"
                        placeholder="Enter a story link"
                        name="storyLink"
                    />
                </div>
                <div class="mb-4">
                    <label
                        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="storyDescription"
                    >
                        Story Content
                    </label>
                    <div class="bg-white">
                        <div
                            class="w-full bg-white"
                            use:quill="{{
                                placeholder: 'Enter story content',
                                content: story.content,
                            }}"
                            on:text-change="{e => {
                                story.content = e.detail.html
                                updateContent()
                            }}"
                            id="storyDescription"
                        ></div>
                    </div>
                </div>
                <div class="mb-4">
                    <div
                        class="text-gray-700 dark:text-gray-400 font-bold text-lg mb-2"
                    >
                        Discussion{story.comments
                            ? ` (${story.comments.length})`
                            : ''}
                    </div>
                    <div class="mb-2">
                        {#if story.comments}
                            {#each story.comments as comment}
                                <div
                                    class="w-full mb-4 text-gray-700 dark:text-gray-400 border-b border-gray-300 dark:border-gray-700"
                                    data-commentid="{comment.id}"
                                >
                                    <div class="font-bold">
                                        <UserIcon
                                            class="h-4 w-4"
                                        />&nbsp;{userMap[comment.user_id]}
                                    </div>
                                    {#if selectedComment !== null && selectedComment.id === comment.id}
                                        <div class="w-full my-2">
                                            <textarea
                                                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                                                bind:value="{selectedCommentContent}"
                                            ></textarea>
                                            <div class="text-right">
                                                <HollowButton
                                                    color="blue"
                                                    onClick="{toggleCommentEdit(
                                                        null,
                                                    )}"
                                                >
                                                    {$_('cancel')}
                                                </HollowButton>
                                                <HollowButton
                                                    color="green"
                                                    onClick="{handleCommentEdit}"
                                                    disabled="{selectedCommentContent ===
                                                        ''}"
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
                                    {#if comment.user_id === $user.id && !(selectedComment !== null && selectedComment.id === comment.id)}
                                        <div class="mb-2 text-right">
                                            <button
                                                class="text-blue-500 hover:text-blue-300 mr-1"
                                                on:click="{toggleCommentEdit(
                                                    comment,
                                                )}"
                                            >
                                                {$_('edit')}
                                            </button>
                                            <button
                                                class="text-red-500 hover:text-red-300"
                                                on:click="{handleCommentDelete(
                                                    comment.id,
                                                )}"
                                            >
                                                {$_('delete')}
                                            </button>
                                        </div>
                                    {/if}
                                </div>
                            {/each}
                        {/if}
                    </div>
                    <div class="w-full">
                        <textarea
                            class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
                            placeholder="Write a comment..."
                            bind:value="{userComment}"></textarea>
                        <div class="text-right">
                            <HollowButton
                                color="teal"
                                onClick="{handleCommentSubmit}"
                                disabled="{userComment === ''}"
                            >
                                Post comment
                            </HollowButton>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="w-1/4">
            <div class="mx-4">
                <div class="mb-4">
                    <label
                        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="storyPoints"
                    >
                        Story Points
                    </label>
                    <input
                        class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
        rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
        focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        id="storyPoints"
                        type="number"
                        min="0"
                        max="999"
                        bind:value="{story.points}"
                        on:change="{updatePoints}"
                        placeholder="Enter story points e.g. 1, 2, 3, 5,
                        8"
                        name="storyPoints"
                    />
                </div>
                <div class="mb-2">
                    <div class="text-gray-700 dark:text-gray-400 font-bold">
                        Storycard Color
                    </div>
                    <div>
                        {#each colorLegend as color}
                            <button
                                on:click="{changeColor(color.color)}"
                                class="p-4 mr-2 mb-2 colorcard-{color.color}
                                border-2 border-solid {story.color ===
                                color.color
                                    ? `border-${color.color}-800`
                                    : 'border-transparent'}"
                                title="{color.color}{color.legend !== ''
                                    ? ` - ${color.legend}`
                                    : ''}"></button>
                        {/each}
                    </div>
                </div>
                <div class="mb-4">
                    {#if !story.closed}
                        <HollowButton color="orange" onClick="{markClosed}">
                            Mark story as Closed
                        </HollowButton>
                    {:else}
                        <HollowButton color="green" onClick="{markOpen}">
                            Reopen story
                        </HollowButton>
                    {/if}
                </div>
                <div class="text-right">
                    <HollowButton color="red" onClick="{handleStoryDelete}">
                        Delete Story
                    </HollowButton>
                </div>
            </div>
        </div>
    </div>
</Modal>

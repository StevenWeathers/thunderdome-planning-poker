<script>
  import CloseIcon from '../icons/CloseIcon.svelte'
  import HollowButton from '../HollowButton.svelte'

  export let toggleStoryForm = () => {}
  export let updateContent = () => () => {}
  export let updateName = () => () => {}
  export let changeColor = () => () => {}
  export let updatePoints = () => () => {}
  export let updateClosed = () => () => {}
  export let deleteStory = () => () => {}
  export let addComment = () => {}

  export let story = {}
  export let colorLegend = []
  export let users = []
  let userComment = ''

  $: userMap = users.reduce((prev, usr) => {
    prev[usr.id] = usr.name
    return prev
  }, {})

  function handleStoryDelete () {
    deleteStory(story.id)()
    toggleStoryForm()
  }

  function markClosed () {
    updateClosed(story.id)(true)
  }

  function markOpen () {
    updateClosed(story.id)(false)
  }

  function handleCommentSubmit () {
    if (userComment !== '') {
      addComment(story.id, userComment)
      userComment = ''
    }
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

<div class="fixed inset-0 flex items-center z-40">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div class="relative mx-4 md:mx-auto w-full md:w-2/3 z-50 m-8">
        <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6">
            <div class="flex justify-end mb-2">
                <button
                        aria-label="close"
                        on:click="{toggleStoryForm}"
                        class="text-gray-800"
                >
                    <CloseIcon/>
                </button>
            </div>

            <div class="flex w-full">
                <div class="w-3/4">
                    <div class="mx-4">
                        <div class="mb-4">
                            <label
                                    class="block text-sm font-bold mb-2"
                                    for="storyName"
                            >
                                Story Name
                            </label>
                            <input
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="storyName"
                                    type="text"
                                    on:change="{updateName(story.id)}"
                                    value="{story.name}"
                                    placeholder="Enter a story name e.g. Ricky Bobby"
                                    name="storyName"
                            />
                        </div>
                        <div class="mb-4">
                            <label
                                    class="block text-sm font-bold mb-2"
                                    for="storyDescription"
                            >
                                Story Content
                            </label>
                            <textarea
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    placeholder="Enter story content"
                                    on:change="{updateContent(story.id)}"
                                    value="{story.content}"></textarea>
                        </div>
                        <div class="mb-4">
                            <div class="font-bold text-lg text-gray-600">
                                Discussion{story.comments
                              ? ` (${story.comments.length})`
                              : ''}
                            </div>
                            <div class="mb-2 w-full">
                                <textarea
                                        class="border-gray-200 border-2
                                    appearance-none rounded w-full py-2 px-3
                                    text-gray-700 leading-tight
                                    focus:outline-none focus:border-purple-500
                                    mb-2"
                                        placeholder="Write a comment..."
                                        bind:value="{userComment}"></textarea>
                                <div class="text-right">
                                    <HollowButton
                                            color="gray"
                                            onClick="{handleCommentSubmit}"
                                            disabled="{userComment === ''}"
                                    >
                                        Post comment
                                    </HollowButton>
                                </div>
                            </div>
                            <div>
                                {#if story.comments}
                                    {#each story.comments as comment}
                                        <div class="w-full mb-4">
                                            <div class="text-sm">
                                                <span class="font-bold">
                                                    {userMap[comment.user_id]}
                                                </span>
                                            </div>
                                            <div>{comment.comment}</div>
                                        </div>
                                    {/each}
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="w-1/4">
                    <div class="mx-4">
                        <div class="mb-4">
                            <label
                                    class="block text-sm font-bold mb-2"
                                    for="storyPoints"
                            >
                                Story Points
                            </label>
                            <input
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="storyPoints"
                                    type="number"
                                    min="0"
                                    max="999"
                                    bind:value="{story.points}"
                                    on:change="{updatePoints(story.id)}"
                                    placeholder="Enter story points e.g. 1, 2, 3, 5,
                                8"
                                    name="storyPoints"
                            />
                        </div>
                        <div class="mb-2">
                            <div class="font-bold">Storycard Color</div>
                            <div>
                                {#each colorLegend as color}
                                    <button
                                            on:click="{changeColor(
                                            story.id,
                                            color.color,
                                        )}"
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
                                <HollowButton
                                        color="orange"
                                        onClick="{markClosed}"
                                >
                                    Mark story as Closed
                                </HollowButton>
                            {:else}
                                <HollowButton
                                        color="green"
                                        onClick="{markOpen}"
                                >
                                    Reopen story
                                </HollowButton>
                            {/if}
                        </div>
                        <div class="text-right">
                            <HollowButton
                                    color="red"
                                    onClick="{handleStoryDelete}"
                            >
                                Delete Story
                            </HollowButton>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

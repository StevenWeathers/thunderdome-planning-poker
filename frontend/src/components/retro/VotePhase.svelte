<script>
    import { warrior as user } from '../../stores'
    import ThumbsUp from '../icons/ThumbsUp.svelte'

    export let groups = []
    export let handleVote = () => {}
    export let handleVoteSubtract = () => {}
    export let voteLimitReached = false

    const handleVoteAction = group => {
        const alreadyVoted = group.votes.includes($user.id)

        if (alreadyVoted) {
            handleVoteSubtract(group.id)
        } else {
            handleVote(group.id)
        }
    }
</script>

{#each groups as group, i (group.id)}
    {#if group.items.length > 0}
        <div
            class="border-2 p-2 dark:border-gray-800 rounded flex flex-col flex-wrap"
        >
            <div class="dark:text-gray-200 w-full text-center text-lg mb-2">
                <div class="flex content-center justify-center">
                    <button
                        on:click="{() => {
                            handleVoteAction(group)
                        }}"
                        disabled="{voteLimitReached && !group.userVoted}"
                        class="inline-block align-middle"
                        class:text-gray-300="{voteLimitReached &&
                            !group.userVoted}"
                        class:dark:text-gray-600="{voteLimitReached &&
                            !group.userVoted}"
                        class:cursor-not-allowed="{voteLimitReached &&
                            !group.userVoted}"
                        class:hover:text-blue-500="{!(
                            voteLimitReached && !group.userVoted
                        )}"
                        class:dark:hover:text-sky-500="{!(
                            voteLimitReached && !group.userVoted
                        )}"
                        class:text-green-500="{group.userVoted}"
                        class:dark:text-lime-500="{group.userVoted}"
                    >
                        <ThumbsUp class="w-6 h-6 inline-block" />
                    </button>
                    <div class="inline-block align-middle text-2xl ml-2">
                        {group.votes.length}
                    </div>
                </div>
                {group.name ? group.name : 'Group'}
            </div>
            <div class="flex-1 grow">
                {#each group.items as item, ii (item.id)}
                    <div
                        class="p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-l-4 dark:text-white"
                        class:border-green-400="{item.type === 'worked'}"
                        class:dark:border-lime-400="{item.type === 'worked'}"
                        class:border-red-500="{item.type === 'improve'}"
                        class:border-blue-400="{item.type === 'question'}"
                        class:dark:border-sky-400="{item.type === 'question'}"
                    >
                        {item.content}
                    </div>
                {/each}
            </div>
        </div>
    {/if}
{/each}

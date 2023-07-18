<script lang="ts">
    import { onMount } from 'svelte'
    import { addTimeLeadZero, timeUnitsBetween } from '../../dateUtils'

    export let currentStoryId: string = ''
    export let votingLocked: boolean = true
    export let voteStartTime: Date = new Date()

    let currentTime: Date = new Date()

    $: voteDuration =
        currentStoryId !== '' && votingLocked === false
            ? timeUnitsBetween(voteStartTime, currentTime)
            : {}

    onMount(() => {
        const voteCounter = setInterval(() => {
            currentTime = new Date()
        }, 1000)

        return () => {
            clearInterval(voteCounter)
        }
    })
</script>

<div
    class="font-semibold
                text-3xl md:text-4xl text-gray-700 dark:text-gray-300"
    data-testid="vote-timer"
>
    {#if voteDuration.seconds !== undefined}
        {#if voteDuration.hours !== 0}
            {addTimeLeadZero(voteDuration.hours)}:
        {/if}
        {addTimeLeadZero(voteDuration.minutes)}:{addTimeLeadZero(
            voteDuration.seconds,
        )}
    {/if}
</div>

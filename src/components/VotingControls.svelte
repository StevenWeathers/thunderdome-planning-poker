<script>
    import SolidButton from './SolidButton.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'

    export let sendSocketEvent = () => {}
    export let planId = ''
    export let points = []
    export let votingLocked = true
    export let highestVote = ''

    $: planPoints = highestVote

    const endPlanVoting = () => {
        sendSocketEvent('end_voting', planId)
    }

    const skipPlan = () => {
        sendSocketEvent('skip_plan', planId)
    }

    function handleSubmit(event) {
        event.preventDefault()

        sendSocketEvent(
            'finalize_plan',
            JSON.stringify({
                planId,
                planPoints,
            }),
        )

        planPoints = ''
    }
</script>

{#if planId != ''}
    <div class="p-4" data-testId="votingControls">
        <SolidButton
            color="blue"
            additionalClasses="mb-2 w-full"
            onClick="{skipPlan}">
            Skip Plan
        </SolidButton>
        {#if !votingLocked}
            <SolidButton additionalClasses="w-full" onClick="{endPlanVoting}">
                Finish Voting
            </SolidButton>
        {:else}
            <form on:submit="{handleSubmit}" name="savePlanPoints">
                <legend class="text-xl mb-2 font-semibold leading-tight">
                    Final Points
                </legend>
                <div class="flex -mx-2">
                    <div class="w-1/2 px-2">
                        <div class="relative">
                            <select
                                name="planPoints"
                                bind:value="{planPoints}"
                                required
                                class="block appearance-none w-full border
                                border-gray-400 text-gray-700 py-3 px-4 pr-8
                                rounded leading-tight focus:outline-none
                                focus:border-gray">
                                <option value="" disabled>Points</option>
                                {#each points as point}
                                    <option value="{point}">{point}</option>
                                {/each}
                            </select>
                            <div
                                class="pointer-events-none absolute inset-y-0
                                right-0 flex items-center px-2 text-gray-700">
                                <DownCarrotIcon />
                            </div>
                        </div>
                    </div>
                    <div class="w-1/2 text-right px-2">
                        <SolidButton
                            additionalClasses="w-full h-full"
                            type="submit">
                            Save
                        </SolidButton>
                    </div>
                </div>
            </form>
        {/if}
    </div>
{/if}

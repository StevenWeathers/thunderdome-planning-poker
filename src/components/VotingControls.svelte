<script>
    export let sendSocketEvent = () => {}
    export let planId = ''
    export let points = []
    export let votingLocked = true

    let planPoints = ''

    const endPlanVoting = () => {
        sendSocketEvent('end_voting', planId)
    }

    function handleSubmit(event) {
        event.preventDefault()

        sendSocketEvent('finalize_plan', JSON.stringify({
            planId,
            planPoints
        }))

        planPoints = ''
    }
</script>

<div class="p-4">
    {#if planId !== '' && !votingLocked}
        <button class="bg-blue hover:bg-blue-dark text-white font-bold py-2 px-4 rounded w-full" on:click={endPlanVoting} disabled={votingLocked}>Finish Voting</button>
    {:else if planId !== '' && votingLocked}
        <form on:submit={handleSubmit}>
            <legend class="text-xl mb-2">Final Points</legend>
            <div class="flex -mx-2">
                <div class="w-1/2 px-2">
                    <div class="relative">
                        <select
                            bind:value={planPoints}
                            required
                            class="block appearance-none w-full border border-grey-light text-grey-darker py-3 px-4 pr-8 rounded leading-tight focus:outline-none focus:border-grey"
                        >
                            <option value="" disabled>Points</option>
                            {#each points as point}
                                <option value="{point}">{point}</option>
                            {/each}
                        </select>
                        <div class="pointer-events-none absolute pin-y pin-r flex items-center px-2 text-grey-darker">
                            <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
                        </div>
                    </div>
                </div>
                <div class="w-1/2 text-right px-2">
                    <button class="bg-green hover:bg-green-dark text-white font-bold py-2 px-4 rounded w-full h-full" type="submit">Save</button>
                </div>
            </div>
        </form>
    {/if}
</div>

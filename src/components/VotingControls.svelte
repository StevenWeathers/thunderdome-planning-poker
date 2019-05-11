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

<div class="box">
    <button class="button is-link is-outlined" on:click={endPlanVoting} disabled={votingLocked}>Finish Voting</button>

    {#if planId !== '' && votingLocked}
    <form on:submit={handleSubmit}>
        <div class="field">
            <label class="label">Plan Name</label>
            <div class="control">
                <div class="select">
                    <select bind:value={planPoints} required>
                        <option value=""></option>
                        {#each points as point}
                            <option value="{point}">{point}</option>
                        {/each}
                    </select>
                </div>
            </div>
        </div>
        
        <div class="field">
            <div class="control">
                <button class="button is-success" type="submit">Finalize Plan</button>
            </div>
        </div>
    </form>
    {/if}
</div>
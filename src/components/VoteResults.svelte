<script>
    import VoteIcon from './VoteIcon.svelte'

    export let activePlanId = ''
    export let plans = []
    
    function compileVoteCounts() {
        const currentPlan = plans.find(p => p.id === activePlanId)
        const voteCounts = currentPlan.votes.reduce((obj, v) => {
            obj[v.vote] = (obj[v.vote] || 0) + 1;
            return obj;
        }, {})
        const voteCountsArray = Object.keys(voteCounts).map(k => ({
                vote: k,
                count: voteCounts[k]
        }))

        return voteCountsArray
    }

    let counts = compileVoteCounts()
</script>

<div class="mb-6">
    <h2 class="text-center">Vote Results</h2>
 
    <div class="w-full">
        {#each counts as {vote, count}}
            <div><VoteIcon />{vote} : ({count})</div>
        {/each}
    </div>
</div>

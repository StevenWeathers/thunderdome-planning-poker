<script>
    import VoteIcon from './VoteIcon.svelte'
    import PointCard from '../components/PointCard.svelte'

    export let activePlanId = ''
    export let plans = []
    export let points = []
    
    function compileVoteCounts() {
        const currentPlan = plans.find(p => p.id === activePlanId)
        return currentPlan.votes.reduce((obj, v) => {
            obj[v.vote] = (obj[v.vote] || 0) + 1;
            return obj;
        }, {})
    }

    let counts = compileVoteCounts()
</script>

<div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
    {#each points as point}
        <div class="w-1/4 md:w-1/6 px-2 mb-4">
            <PointCard point={point} isLocked={true} count={counts[point]} />
        </div>
    {/each}
</div>
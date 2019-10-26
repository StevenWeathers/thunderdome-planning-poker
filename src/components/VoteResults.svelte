<script>
    import VoteIcon from './VoteIcon.svelte'
    import PointCard from '../components/PointCard.svelte'

    export let activePlanId = ''
    export let plans = []
    export let points = []
    export let highestVote = ''

    let totalVotes = plans.find(p => p.id === activePlanId).votes.length

    // get vote average from active plan
    function getVoteAverage() {
        const activePlan = plans.find(p => p.id === activePlanId)
        let average = 0

        if (activePlan.votes.length > 0) {
            const votesToAverage = activePlan.votes.filter((v) => v.vote !== '?').map(v => {
                const vote = v.vote === '1/2' ? .5 : parseInt(v.vote)
                return vote;
            })
            const sum = votesToAverage.reduce((previous, current) => current += previous)
            
            average = Math.ceil(sum / votesToAverage.length)
        }

        return average
    }

    function compileVoteCounts() {
        const currentPlan = plans.find(p => p.id === activePlanId)
        return currentPlan.votes.reduce((obj, v) => {
            obj[v.vote] = (obj[v.vote] || 0) + 1;
            return obj;
        }, {})
    }

    function calculateVotePercentage(count) {
        return Math.floor(100 * (count / totalVotes))
    }
    
    let counts = compileVoteCounts()
    let average = getVoteAverage()
</script>

<div class="flex flex-wrap items-center mb-4 lg:mb-6 pt-4 pb-4 lg:pt-6 lg:pb-6 bg-white shadow-md rounded">
    <div class="w-1/2">
        <div class="text-center">
            <div class="text-2xl mb-1">
                {totalVotes} warrior{totalVotes > 1 ? 's': ''}
            </div>
            <div class="text-grey-darker mb-3">voted</div>
            <div class="text-xl">Avg: {average}</div>
        </div>
    </div>
    <div class="w-1/2">
        <ul class="list-reset">
            {#each points.filter(p => counts[p] !== undefined) as point}
                <li class="mb-2 lg:mb-4">
                    <div class="text-xl {highestVote === point ? 'font-bold text-green-dark' : ''}">{point}</div>
                    <div class="text-grey-dark">
                        {calculateVotePercentage(counts[point])} % ({counts[point]} warrior{counts[point] > 1 ? 's' : ''})
                    </div>
                </li>
            {/each}
        </ul>
    </div>
</div>
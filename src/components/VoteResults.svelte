<script>
    import VoteIcon from './VoteIcon.svelte'
    import PointCard from '../components/PointCard.svelte'
    import WarriorIcon from '../components/WarriorIcon.svelte'

    export let activePlanId = ''
    export let plans = []
    export let points = []
    export let warriors = []
    export let highestVote = ''

    let totalVotes = plans.find(p => p.id === activePlanId).votes.length
console.log('warriors >> ', warriors)
    // get vote average from active plan
    function getVoteAverage() {
        const activePlan = plans.find(p => p.id === activePlanId)
        let average = 0

        if (activePlan.votes.length > 0) {
            const votesToAverage = activePlan.votes.filter((v) => v.vote !== '?').map(v => {
                const vote = v.vote === '1/2' ? .5 : parseInt(v.vote)
                return vote
            })
            const sum = votesToAverage.reduce((previous, current) => current += previous)
            
            average = Math.ceil(sum / votesToAverage.length)
        }

        return average
    }

    function compileVoteCounts() {
        const currentPlan = plans.find(p => p.id === activePlanId)
        return currentPlan.votes.reduce((obj, v) => {
            const currentVote = obj[v.vote] || {
                count: 0,
                voters: [],
            }
            
            currentVote.count = currentVote.count + 1
            // need to get warriors array to update properly from parent
            // currentVote.voters.push(warriors.find(w => w.id === v.warriorId).name)

            obj[v.vote] = currentVote
            
            return obj
        }, {})
    }

    function calculateVotePercentage(count) {
        return Math.floor(100 * (count / totalVotes))
    }
    
    let counts = compileVoteCounts()
    let average = getVoteAverage()
</script>

<div class="flex flex-wrap items-center text-center mb-2 md:mb-4 pt-2 pb-2 md:pt-4 md:pb-4 bg-white shadow-md rounded text-xl">
    <div class="w-1/3 ">
            {totalVotes} <WarriorIcon />{totalVotes > 1 ? 's': ''} voted
    </div>
    <div class="w-1/3">
        <div class="mb-2">Average</div>
        <span class="font-bold text-green-dark border-green border p-2 rounded ml-2 inline-block">{average}</span>
    </div>
    <div class="w-1/3">
        <div class="mb-2">Highest</div>
        <span class="font-bold text-green-dark border-green border p-2 rounded ml-2 inline-block">{highestVote}</span> - {counts[highestVote].count}<WarriorIcon />
    </div>
</div>

<div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
    {#each points as point}
        <div class="w-1/4 md:w-1/6 px-2 mb-4">
            <PointCard results={counts[point]} isLocked={true} point={point} />
        </div>
    {/each}
</div>
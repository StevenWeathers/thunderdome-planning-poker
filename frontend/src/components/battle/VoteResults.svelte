<script>
    import PointCard from './PointCard.svelte'
    import WarriorIcon from '../icons/UserIcon.svelte'
    import { _ } from '../../i18n.js'

    export let activePlanId = ''
    export let plans = []
    export let points = []
    export let warriors = []
    export let highestVote = ''
    export let averageRounding = 'ceil'

    let totalVotes = plans.find(p => p.id === activePlanId).votes.length

    // get vote average from active plan
    function getVoteAverage() {
        const activePlan = plans.find(p => p.id === activePlanId)
        let average = 0

        if (activePlan.votes.length > 0) {
            const votesToAverage = activePlan.votes
                .filter(v => {
                    const { spectator = false } = warriors.find(
                        w => w.id === v.warriorId,
                    )
                    return !spectator && v.vote !== '?'
                })
                .map(v => {
                    const vote = v.vote === '1/2' ? 0.5 : parseInt(v.vote)
                    return vote
                })

            const sum = votesToAverage.length
                ? votesToAverage.reduce(
                      (previous, current) => (current += previous),
                  )
                : 0

            const preAverage = sum / votesToAverage.length || 0
            switch (averageRounding) {
                case 'round':
                    average = Math.round(preAverage)
                    break
                case 'floor':
                    average = Math.floor(preAverage)
                    break
                default:
                    average = Math.ceil(preAverage)
            }
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
            let warriorName = $_('pages.battle.voteResults.unknownWarrior')

            if (warriors.length) {
                const warrior = warriors.find(w => w.id === v.warriorId)
                warriorName = warrior ? warrior.name : warriorName
                if (warrior.spectator) {
                    return obj
                }
            }
            currentVote.voters.push(warriorName)

            currentVote.count = currentVote.count + 1

            obj[v.vote] = currentVote

            return obj
        }, {})
    }

    $: average = getVoteAverage(warriors)
    $: counts = compileVoteCounts(warriors)
    let showHighestVoters = false
</script>

<div
    class="flex flex-wrap items-center text-center mb-2 md:mb-4 pt-2 pb-2
    md:pt-4 md:pb-4 bg-white dark:bg-gray-800 shadow-lg rounded-lg text-xl"
>
    <div class="w-1/3 dark:text-white">
        <div class="mb-2">{$_('pages.battle.voteResults.totalVotes')}</div>
        <span data-testid="voteresult-total">{totalVotes}</span>
        <WarriorIcon class="h-5 w-5" />
    </div>
    <div class="w-1/3 dark:text-white">
        <div class="mb-2">{$_('pages.battle.voteResults.average')}</div>
        <span
            class="font-bold text-green-600 dark:text-lime-400 border-green-500 dark:border-lime-500 border p-2 rounded
            ml-2 inline-block"
            data-testid="voteresult-average"
        >
            {average}
        </span>
    </div>
    <div class="w-1/3 dark:text-white">
        <div class="mb-2">{$_('pages.battle.voteResults.highest')}</div>
        <div>
            <span
                class="font-bold text-green-600 dark:text-lime-400 border-green-500 dark:border-lime-500 border p-2
                rounded ml-2 inline-block"
                data-testid="voteresult-high"
            >
                {highestVote || 0}
            </span>
            -
            <span data-testid="voteresult-highcount"
                >{counts[highestVote] ? counts[highestVote].count : 0}</span
            >
            <span class="relative">
                <button
                    on:mouseenter="{() => (showHighestVoters = true)}"
                    on:mouseleave="{() => (showHighestVoters = false)}"
                    class="relative leading-none"
                    title="{$_('pages.battle.voteResults.showVoters')}"
                >
                    <WarriorIcon class="h-5 w-5" />
                    <span
                        class="text-sm text-right text-gray-900 font-normal w-48
                        absolute left-0 top-0 -mt-2 ml-4 bg-white p-2 rounded
                        shadow-lg {showHighestVoters ? '' : 'hidden'}"
                    >
                        {#if counts[highestVote]}
                            {#each counts[highestVote].voters as voter}
                                {voter}
                                <br />
                            {/each}
                        {/if}
                    </span>
                </button>
            </span>
        </div>
    </div>
</div>

<div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
    {#each points as point}
        <div class="w-1/4 md:w-1/6 px-2 mb-4">
            <PointCard
                results="{counts[point] || { count: 0 }}"
                isLocked="{true}"
                point="{point}"
            />
        </div>
    {/each}
</div>

<script lang="ts">
  import PointCard from './PointCard.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { User } from 'lucide-svelte';

  export let activePlanId = '';
  export let plans = [];
  export let points = [];
  export let warriors = [];
  export let highestVote = '';
  export let averageRounding = 'ceil';
  export let hideVoterIdentity = false;

  let totalVotes = plans.find(p => p.id === activePlanId).votes.length;

  // get vote average from active plan
  function getVoteAverage() {
    const activePlan = plans.find(p => p.id === activePlanId);
    let average = 0;

    if (activePlan.votes.length > 0) {
      let sum = 0;
      let votesToAverage = activePlan.votes.reduce((prev, v) => {
        const voteWarrior = warriors.find(w => w.id === v.warriorId) || {};
        const { spectator = false } = voteWarrior;

        if (!spectator && !isNaN(v.vote)) {
          const vote = v.vote === '1/2' ? 0.5 : parseInt(v.vote);
          prev.push(vote);
          sum += vote;
        }

        return prev;
      }, []);

      const preAverage = sum / votesToAverage.length || 0;
      if (preAverage !== 0.5) {
        switch (averageRounding) {
          case 'round':
            average = Math.round(preAverage);
            break;
          case 'floor':
            average = Math.floor(preAverage);
            break;
          default:
            average = Math.ceil(preAverage);
        }
      } else {
        average = preAverage;
      }
    }

    return average;
  }

  function compileVoteCounts() {
    const currentPlan = plans.find(p => p.id === activePlanId);
    return currentPlan.votes.reduce((obj, v) => {
      const currentVote = obj[v.vote] || {
        count: 0,
        voters: [],
      };
      let warriorName = $LL.unknownWarrior();

      if (warriors.length) {
        const warrior = warriors.find(w => w.id === v.warriorId) || {
          name: warriorName,
          spectator: false,
        };
        warriorName = warrior.name;
        if (warrior.spectator) {
          return obj;
        }
      }
      currentVote.voters.push(warriorName);

      currentVote.count = currentVote.count + 1;

      obj[v.vote] = currentVote;

      return obj;
    }, {});
  }

  $: average = getVoteAverage(warriors);
  $: counts = compileVoteCounts(warriors);
  let showHighestVoters = false;
</script>

<div
  class="flex flex-wrap items-center text-center mb-2 md:mb-4 pt-2 pb-2
    md:pt-4 md:pb-4 bg-white dark:bg-gray-800 shadow-lg rounded-lg text-xl"
>
  <div class="w-1/3 dark:text-white">
    <div class="mb-2">{$LL.totalVotes()}</div>
    <span data-testid="voteresult-total">{totalVotes}</span>
    <User class="h-5 w-5 inline-block" />
  </div>
  <div class="w-1/3 dark:text-white">
    <div class="mb-2">{$LL.voteResultsAverage()}</div>
    <span
      class="font-bold text-green-600 dark:text-lime-400 border-green-500 dark:border-lime-500 border p-2 rounded
            me-2 inline-block"
      data-testid="voteresult-average"
    >
      {average}
    </span>
  </div>
  <div class="w-1/3 dark:text-white">
    <div class="mb-2">{$LL.voteResultsHighest()}</div>
    <div>
      <span
        class="font-bold text-green-600 dark:text-lime-400 border-green-500 dark:border-lime-500 border p-2
          rounded ms-2 inline-block"
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
          on:mouseenter="{() => {
            if (!hideVoterIdentity) {
              showHighestVoters = true;
            }
          }}"
          on:mouseleave="{() => {
            if (!hideVoterIdentity) {
              showHighestVoters = false;
            }
          }}"
          class="relative leading-none"
          title="{$LL.showVoters()}"
        >
          <User class="h-5 w-5 inline-block" />
          <span
            class="text-sm text-right text-gray-900 font-normal w-48
                        absolute start-0 top-0 -mt-2 ms-4 bg-white p-2 rounded
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
        hideVoterIdentity="{hideVoterIdentity}"
      />
    </div>
  {/each}
</div>

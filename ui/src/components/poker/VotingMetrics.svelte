<script>
  import { run } from 'svelte/legacy';

  import { scaleLinear } from 'd3-scale';
  import { TriangleAlert } from '@lucide/svelte';
  import { onMount } from 'svelte';

  /** @type {{votes?: any, pointValues?: any, users?: any, averageRounding?: string}} */
  let {
    votes = [],
    pointValues = ['XS', 'S', 'M', 'L', 'XL', 'XXL', '?'],
    users = [],
    averageRounding = 'ceil',
  } = $props();
  let chartData = $state([]);
  let consensusValue = $state('');
  let consensusPercentage = $state(0);
  let isNumeric = $state(false);
  let totalVoters = $state(0);
  let userMap = $state({});

  run(() => {
    userMap = users.reduce((prev, u) => {
      prev[u.id] = u.name;
      return prev;
    }, {});
    isNumeric = pointValues.filter(v => v !== '1/2' && v !== '?' && v !== `☕️`).every(v => !isNaN(v));

    chartData = pointValues.map(value => {
      const count = votes.filter(v => v.vote === value).length;
      const users = votes
        .filter(v => v.vote === value)
        .map(v => {
          return userMap[v.warriorId] || 'Unknown';
        });
      return { value, count, users };
    });

    const modeData = chartData.reduce((a, b) => (b.count > a.count ? b : a), {
      count: 0,
    });
    consensusValue = modeData.value || 'N/A';
    totalVoters = votes.length;
    consensusPercentage = totalVoters > 0 ? Math.round((modeData.count / totalVoters) * 100) : 0;
  });

  function roundWithConfiguredAvg(middleIndex) {
    let average = 0;
    switch (averageRounding) {
      case 'round':
        average = Math.round(middleIndex);
        break;
      case 'floor':
        average = Math.floor(middleIndex);
        break;
      default:
        average = Math.ceil(middleIndex);
    }

    return average;
  }

  function getAverageOrMedian(votes, pointValues) {
    if (isNumeric) {
      const numericVotes = votes
        .filter(v => !isNaN(v.vote) || v.vote === '1/2')
        .map(v => (v.vote === '1/2' ? 0.5 : Number(v.vote)));
      const sum = numericVotes.reduce((a, b) => a + b, 0);
      let average = 0;

      if (numericVotes.length > 0) {
        const preAverage = sum / numericVotes.length || 0;
        if (preAverage !== 0.5) {
          average = roundWithConfiguredAvg(preAverage);
        } else {
          average = preAverage;
        }
      } else {
        average = 'N/A';
      }

      return average;
    } else {
      const validVotes = votes.filter(v => v.vote !== '?').map(v => v.vote);
      if (validVotes.length === 0) return 'N/A';

      const sortedVotes = validVotes.sort((a, b) => pointValues.indexOf(a) - pointValues.indexOf(b));
      const middleIndex = roundWithConfiguredAvg(sortedVotes.length / 2);

      if (sortedVotes.length % 2 === 0) {
        // If even number of votes, take the middle two and find the value between them
        const lowerMiddle = pointValues.indexOf(sortedVotes[middleIndex - 1]);
        const upperMiddle = pointValues.indexOf(sortedVotes[middleIndex]);
        const averageIndex = roundWithConfiguredAvg((lowerMiddle + upperMiddle) / 2);
        return pointValues[averageIndex];
      } else {
        // If odd number of votes, return the middle value
        return sortedVotes[middleIndex];
      }
    }
  }

  function getBarHeight(count) {
    const minHeight = 5; // Minimum height in percentage
    const scaledHeight = totalVoters > 0 ? (count / totalVoters) * 100 : 0;
    return count > 0 ? Math.max(scaledHeight, minHeight) : 0;
  }

  let averageOrMedian = $derived(getAverageOrMedian(votes, pointValues) || 'N/A');
</script>

<div class="p-4 rounded-lg shadow-md bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100">
  <div
    class="flex items-end space-x-1 h-64 mb-8 bg-gray-100 dark:bg-gray-700 rounded-lg overflow-hidden"
    data-testid="voteresult-graph"
  >
    {#each chartData as { value, count, users }}
      <div class="flex flex-col items-center flex-1 group h-full relative" data-testid="voteresult-graph-bar">
        <div class="w-full h-full flex flex-col justify-end">
          <div
            class="bg-blue-500 transition-all duration-500 ease-out absolute bottom-0 left-0 right-0"
            style="height: {getBarHeight(count)}%; min-height: 2px;"
          ></div>
          <div class="w-full relative" style="height: {getBarHeight(count)}%; min-height: 80px;">
            <div
              class="h-full flex items-center justify-center font-bold text-lg {getBarHeight(count) > 10
                ? 'text-white'
                : 'dark:text-white'}"
              data-testid="voteresult-graph-count"
            >
              {count > 0 ? count : ''}
            </div>
            {#if count > 0}
              <div class="absolute top-0 left-0 right-0 text-center pt-1 text-xs text-blue-300">
                {getBarHeight(count).toFixed(1)}%
              </div>
              <div
                data-testid="voteresult-graph-users"
                class="absolute {getBarHeight(count) === 100
                  ? 'top-1/2 -translate-y-1/2'
                  : 'bottom-full'} left-1/2 transform -translate-x-1/2 {getBarHeight(count) === 100
                  ? 'mb-0'
                  : 'mb-2'} bg-gray-600 text-white p-2 rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-sm whitespace-nowrap z-50"
              >
                {users.join(', ')}
              </div>
            {/if}
          </div>
        </div>
        <div class="absolute bottom-0 left-0 right-0 text-center pb-1">
          <span
            class="font-semibold text-sm {getBarHeight(count) > 5 ? 'text-white' : ''}"
            data-testid="voteresult-graph-value">{value}</span
          >
        </div>
      </div>
    {/each}
  </div>

  <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center">
    <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg" data-testid="voteresult-total">
      <div class="text-3xl font-bold">{totalVoters}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Total Voters</div>
    </div>
    <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg" data-testid="voteresult-consensus">
      <div class="text-3xl font-bold">{consensusValue}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Consensus</div>
    </div>
    <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg" data-testid="voteresult-agreement">
      <div class="text-3xl font-bold">{consensusPercentage}%</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Agreement</div>
    </div>
    <div class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg" data-testid="voteresult-average">
      <div class="text-3xl font-bold">{averageOrMedian}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">
        {isNumeric ? 'Average' : 'Median'}
      </div>
    </div>
  </div>

  {#if consensusPercentage < 70}
    <div class="mt-4 text-center text-yellow-800 dark:text-yellow-400" data-testid="voteresult-discussion">
      <span class="text-2xl me-2"><TriangleAlert class="inline-block" /></span> Low consensus. Further discussion recommended.
    </div>
  {/if}
</div>

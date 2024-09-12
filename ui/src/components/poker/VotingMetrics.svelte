<script>
  import { scaleLinear } from 'd3-scale';
  import { TriangleAlert } from 'lucide-svelte';
  import { onMount } from 'svelte';

  export let votes = [];
  export let pointValues = ['XS', 'S', 'M', 'L', 'XL', 'XXL', '?'];
  export let users = [];
  let chartData = [];
  let maxCount = 0;
  let consensusValue = '';
  let consensusPercentage = 0;
  let isNumeric = false;
  let totalVoters = 0;
  let userMap = {};

  $: {
    userMap = users.reduce((prev, u) => {
      prev[u.id] = u.name;
      return prev;
    }, {});
    isNumeric = pointValues.every(value => !isNaN(value) && value !== '?');

    chartData = pointValues.map(value => {
      const count = votes.filter(v => v.vote === value).length;
      const users = votes
        .filter(v => v.vote === value)
        .map(v => {
          return userMap[v.warriorId] || 'Unknown';
        });
      return { value, count, users };
    });

    maxCount = Math.max(...chartData.map(d => d.count));

    const modeData = chartData.reduce((a, b) => (b.count > a.count ? b : a), {
      count: 0,
    });
    consensusValue = modeData.value || 'N/A';
    totalVoters = votes.length;
    consensusPercentage =
      totalVoters > 0 ? Math.round((modeData.count / totalVoters) * 100) : 0;
  }

  const yScale = scaleLinear().domain([0, maxCount]).range([0, 100]);

  function getAverageOrMedian(votes, pointValues) {
    if (isNumeric) {
      const numericVotes = votes
        .filter(v => !isNaN(v.vote))
        .map(v => Number(v.vote));
      const sum = numericVotes.reduce((a, b) => a + b, 0);
      return numericVotes.length > 0
        ? (sum / numericVotes.length).toFixed(1)
        : 'N/A';
    } else {
      const validVotes = votes.filter(v => v.vote !== '?').map(v => v.vote);
      if (validVotes.length === 0) return 'N/A';

      const sortedVotes = validVotes.sort(
        (a, b) => pointValues.indexOf(a) - pointValues.indexOf(b),
      );
      const middleIndex = Math.floor(sortedVotes.length / 2);

      if (sortedVotes.length % 2 === 0) {
        // If even number of votes, take the middle two and find the value between them
        const lowerMiddle = pointValues.indexOf(sortedVotes[middleIndex - 1]);
        const upperMiddle = pointValues.indexOf(sortedVotes[middleIndex]);
        const averageIndex = Math.round((lowerMiddle + upperMiddle) / 2);
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

  $: averageOrMedian = getAverageOrMedian(votes, pointValues);
</script>

<div
  class="p-4 rounded-lg shadow-md bg-white dark:bg-gray-800 text-gray-800 dark:text-gray-100"
>
  <div
    class="flex items-end space-x-1 h-64 mb-8 bg-gray-700 rounded-lg overflow-hidden"
    data-testid="voteresult-graph"
  >
    {#each chartData as { value, count, users }}
      <div
        class="flex flex-col items-center flex-1 group h-full relative"
        data-testid="voteresult-graph-bar"
      >
        <div class="w-full h-full flex flex-col justify-end">
          <div
            class="w-full bg-blue-500 transition-all duration-500 ease-out relative"
            style="height: {getBarHeight(count)}%; min-height: 2px;"
          >
            <div
              class="h-full flex items-center justify-center text-white font-bold text-lg"
              data-testid="voteresult-graph-count"
            >
              {count > 0 ? count : ''}
            </div>
            {#if count > 0}
              <div
                class="absolute top-0 left-0 right-0 text-center pt-1 text-xs text-blue-300"
              >
                {getBarHeight(count).toFixed(1)}%
              </div>
              <div
                data-testid="voteresult-graph-users"
                class="absolute {getBarHeight(count) === 100
                  ? 'top-1/2 -translate-y-1/2'
                  : 'bottom-full'} left-1/2 transform -translate-x-1/2 {getBarHeight(
                  count,
                ) === 100
                  ? 'mb-0'
                  : 'mb-2'} bg-gray-600 p-2 rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-sm whitespace-nowrap z-50"
              >
                {users.join(', ')}
              </div>
            {/if}
          </div>
        </div>
        <div class="absolute bottom-0 left-0 right-0 text-center pb-1">
          <span
            class="font-semibold text-sm"
            data-testid="voteresult-graph-value">{value}</span
          >
        </div>
      </div>
    {/each}
  </div>

  <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center">
    <div
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg"
      data-testid="voteresult-total"
    >
      <div class="text-3xl font-bold">{totalVoters}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Total Voters</div>
    </div>
    <div
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg"
      data-testid="voteresult-consensus"
    >
      <div class="text-3xl font-bold">{consensusValue}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Consensus</div>
    </div>
    <div
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg"
      data-testid="voteresult-agreement"
    >
      <div class="text-3xl font-bold">{consensusPercentage}%</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">Agreement</div>
    </div>
    <div
      class="bg-gray-100 dark:bg-gray-700 p-4 rounded-lg"
      data-testid="voteresult-average"
    >
      <div class="text-3xl font-bold">{averageOrMedian}</div>
      <div class="text-sm text-gray-600 dark:text-gray-400">
        {isNumeric ? 'Average' : 'Median'}
      </div>
    </div>
  </div>

  {#if consensusPercentage < 70}
    <div
      class="mt-4 text-center text-yellow-800 dark:text-yellow-400"
      data-testid="voteresult-discussion"
    >
      <span class="text-2xl mr-2"><TriangleAlert class="inline-block" /></span> Low
      consensus. Further discussion recommended.
    </div>
  {/if}
</div>

<script lang="ts">
  import { ArrowDownRight, ArrowRight, ArrowUpRight } from '@lucide/svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { HealthCheckAreaResult } from '../../types/healthcheck';

  interface Props {
    results: HealthCheckAreaResult[];
  }

  let { results }: Props = $props();

  const ratings = $derived([
    { value: 'green' as const, label: $LL.healthCheckRatingGreen(), bar: 'bg-green-500', dot: 'bg-green-500' },
    { value: 'yellow' as const, label: $LL.healthCheckRatingYellow(), bar: 'bg-yellow-400', dot: 'bg-yellow-400' },
    { value: 'red' as const, label: $LL.healthCheckRatingRed(), bar: 'bg-red-500', dot: 'bg-red-500' },
  ]);

  const trends = $derived([
    { value: 'improving' as const, label: $LL.healthCheckTrendImproving(), icon: ArrowUpRight },
    { value: 'stable' as const, label: $LL.healthCheckTrendStable(), icon: ArrowRight },
    { value: 'declining' as const, label: $LL.healthCheckTrendDeclining(), icon: ArrowDownRight },
  ]);

  const total = (result: HealthCheckAreaResult) => result.ratings.green + result.ratings.yellow + result.ratings.red;
</script>

<div class="space-y-4">
  {#each results as result (result.name)}
    <section class="rounded-lg bg-white dark:bg-gray-800 shadow p-4 space-y-3" data-testid="health-check-result">
      <h3 class="text-xl font-semibold font-rajdhani dark:text-white">{result.name}</h3>

      <!-- the counts below carry the same information for those who can't see the bar -->
      <div class="flex h-3 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700" aria-hidden="true">
        {#each ratings as rating (rating.value)}
          {#if result.ratings[rating.value] > 0}
            <div class={rating.bar} style="width: {(result.ratings[rating.value] / total(result)) * 100}%"></div>
          {/if}
        {/each}
      </div>

      <ul class="flex flex-wrap gap-4 text-sm dark:text-gray-200">
        {#each ratings as rating (rating.value)}
          <li class="flex items-center gap-1" data-testid="health-check-rating-{rating.value}">
            <span class="inline-block h-3 w-3 rounded-full {rating.dot}" aria-hidden="true"></span>
            {rating.label}: {result.ratings[rating.value]}
          </li>
        {/each}
      </ul>

      <ul class="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-400">
        {#each trends as trend (trend.value)}
          {@const Icon = trend.icon}
          <li class="flex items-center gap-1">
            <Icon class="h-4 w-4" aria-hidden="true" />
            {trend.label}: {result.trends[trend.value]}
          </li>
        {/each}
      </ul>

      {#if result.comments.length}
        <div>
          <h4 class="font-bold text-gray-700 dark:text-gray-300">{$LL.comments()}</h4>
          <ul class="list-disc ps-5 text-gray-700 dark:text-gray-300">
            {#each result.comments as comment}
              <li class="whitespace-pre-wrap break-words">{comment}</li>
            {/each}
          </ul>
        </div>
      {/if}
    </section>
  {/each}
</div>

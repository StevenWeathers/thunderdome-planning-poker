<script lang="ts">
  import { untrack } from 'svelte';

  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type {
    HealthCheckAnswer,
    HealthCheckArea,
    HealthCheckRating,
    HealthCheckTrend,
  } from '../../types/healthcheck';

  interface Props {
    areas: HealthCheckArea[];
    answers?: HealthCheckAnswer[];
    handleSubmit: (answers: HealthCheckAnswer[]) => void;
  }

  let { areas, answers = [], handleSubmit }: Props = $props();

  let responses = $state(
    untrack(() =>
      areas.map(area => {
        const answer = answers.find(a => a.area === area.name);
        return {
          area: area.name,
          rating: (answer?.rating ?? '') as HealthCheckRating | '',
          trend: (answer?.trend ?? 'stable') as HealthCheckTrend,
          comment: answer?.comment ?? '',
        };
      }),
    ),
  );

  const complete = $derived(responses.every(r => r.rating !== ''));

  const ratings = $derived<{ value: HealthCheckRating; label: string; classes: string }[]>([
    {
      value: 'green',
      label: $LL.healthCheckRatingGreen(),
      classes: 'has-[:checked]:bg-green-200 dark:has-[:checked]:bg-green-800 border-green-500',
    },
    {
      value: 'yellow',
      label: $LL.healthCheckRatingYellow(),
      classes: 'has-[:checked]:bg-yellow-200 dark:has-[:checked]:bg-yellow-700 border-yellow-500',
    },
    {
      value: 'red',
      label: $LL.healthCheckRatingRed(),
      classes: 'has-[:checked]:bg-red-200 dark:has-[:checked]:bg-red-800 border-red-500',
    },
  ]);

  const trends = $derived<{ value: HealthCheckTrend; label: string }[]>([
    { value: 'improving', label: $LL.healthCheckTrendImproving() },
    { value: 'stable', label: $LL.healthCheckTrendStable() },
    { value: 'declining', label: $LL.healthCheckTrendDeclining() },
  ]);

  function onSubmit(event: Event) {
    event.preventDefault();
    if (!complete) return;

    handleSubmit(responses.map(r => ({ ...r, rating: r.rating as HealthCheckRating })));
  }
</script>

<form onsubmit={onSubmit} name="healthCheckResponse" class="space-y-4">
  {#each areas as area, i (area.name)}
    <fieldset class="rounded-lg bg-white dark:bg-gray-800 shadow p-4 space-y-3">
      <legend class="sr-only">{area.name}</legend>
      <div class="text-xl font-semibold font-rajdhani dark:text-white" aria-hidden="true">{area.name}</div>
      <dl class="grid gap-2 sm:grid-cols-2 text-sm text-gray-600 dark:text-gray-400">
        {#if area.healthy}
          <div>
            <dt class="font-bold text-green-700 dark:text-lime-400">{$LL.healthCheckHealthy()}</dt>
            <dd>{area.healthy}</dd>
          </div>
        {/if}
        {#if area.unhealthy}
          <div>
            <dt class="font-bold text-red-700 dark:text-red-400">{$LL.healthCheckUnhealthy()}</dt>
            <dd>{area.unhealthy}</dd>
          </div>
        {/if}
      </dl>

      <div class="flex flex-wrap gap-2" role="radiogroup" aria-label={$LL.healthCheckRating()}>
        {#each ratings as rating (rating.value)}
          <label
            class="cursor-pointer rounded-lg border-2 px-4 py-2 dark:text-white has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-indigo-500 {rating.classes}"
          >
            <input
              type="radio"
              class="me-1"
              name="rating-{i}"
              value={rating.value}
              bind:group={responses[i].rating}
              required
            />
            {rating.label}
          </label>
        {/each}
      </div>

      <div class="flex flex-wrap gap-4 dark:text-white" role="radiogroup" aria-label={$LL.healthCheckTrend()}>
        {#each trends as trend (trend.value)}
          <label class="cursor-pointer">
            <input
              type="radio"
              class="me-1"
              name="trend-{i}"
              value={trend.value}
              bind:group={responses[i].trend}
            />
            {trend.label}
          </label>
        {/each}
      </div>

      <textarea
        bind:value={responses[i].comment}
        aria-label="{area.name}: {$LL.healthCheckComment()}"
        placeholder={$LL.healthCheckComment()}
        maxlength="1000"
        rows="2"
        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-gray-100"
      ></textarea>
    </fieldset>
  {/each}

  <div class="text-right">
    <SolidButton type="submit" disabled={!complete} testid="health-check-submit">
      {$LL.healthCheckSubmit()}
    </SolidButton>
  </div>
</form>

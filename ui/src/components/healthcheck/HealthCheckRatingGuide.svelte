<script lang="ts">
  import LL from '../../i18n/i18n-svelte';

  const ratings = $derived([
    { label: $LL.healthCheckRatingGreen(), meaning: $LL.healthCheckAboutGreen(), dot: 'bg-green-500' },
    { label: $LL.healthCheckRatingYellow(), meaning: $LL.healthCheckAboutYellow(), dot: 'bg-yellow-400' },
    { label: $LL.healthCheckRatingRed(), meaning: $LL.healthCheckAboutRed(), dot: 'bg-red-500' },
  ]);
</script>

<div>
  <h3 class="font-bold mb-1">{$LL.healthCheckAboutRatings()}</h3>
  <ul class="space-y-1">
    {#each ratings as rating (rating.dot)}
      <li>
        <span class="inline-block h-3 w-3 rounded-full me-1 {rating.dot}" aria-hidden="true"></span>
        <strong>{rating.label}:</strong>
        {rating.meaning}
      </li>
    {/each}
  </ul>
</div>

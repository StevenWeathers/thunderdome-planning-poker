<script>
  let {
    title,
    description,
    benefits
  } = $props();
</script>

<section class="mt-24 bg-gradient-to-r from-indigo-100 to-blue-100 dark:from-indigo-900/20 dark:to-blue-900/20 rounded-2xl p-8 text-center">
  <h2 class="text-3xl font-bold mb-4 dark:text-white">
    {title}
  </h2>
  <p class="text-lg text-gray-600 dark:text-gray-400 mb-6 max-w-2xl mx-auto">
    {description}
  </p>
  <div class="grid md:grid-cols-3 gap-6 text-center">
    {#each benefits as benefit}
      <div class="p-4">
        <div class="text-2xl font-bold text-indigo-600 dark:text-indigo-400 mb-2">{benefit.title}</div>
        <div class="text-gray-600 dark:text-gray-400">{benefit.description}</div>
      </div>
    {/each}
  </div>
</section>
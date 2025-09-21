<script>
  let { title, subtitle, problemTitle, problemItems, solutionTitle, solutionItems, conclusion } = $props();
</script>

<section class="mb-20 text-gray-700 dark:text-gray-300">
  <div class="mb-6 text-center">
    <div class="title-line bg-yellow-thunder mx-auto"></div>
    <h2 class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-4">
      {title}
    </h2>
    <p class="text-xl text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
      {subtitle}
    </p>
  </div>

  <div class="grid md:grid-cols-2 gap-8 mb-8">
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm">
      <h3 class="font-semibold text-lg mb-3 text-indigo-600 dark:text-indigo-400">{problemTitle}</h3>
      <ul class="space-y-2 text-gray-600 dark:text-gray-400">
        {#each problemItems as item}
          <li>• {item}</li>
        {/each}
      </ul>
    </div>
    <div class="bg-indigo-100 dark:bg-indigo-900/20 p-6 rounded-lg shadow-sm">
      <h3 class="font-semibold text-lg mb-3 text-indigo-600 dark:text-indigo-400">{solutionTitle}</h3>
      <ul class="space-y-2 text-gray-700 dark:text-gray-300">
        {#each solutionItems as item}
          <li>• {item}</li>
        {/each}
      </ul>
    </div>
  </div>

  <p class="text-center text-lg">
    {@html conclusion}
  </p>
</section>

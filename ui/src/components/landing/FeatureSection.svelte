<script>
  import { PathPrefix } from '../../config';

  let { title, description, features, imageDark, imageLight, imageAlt, index = 0 } = $props();

  const reverse = $derived(index % 2 === 1);

  const sectionClass = $derived(reverse
    ? 'flex flex-col md:flex-row-reverse items-center'
    : 'flex flex-col md:flex-row items-center');

  const textClass = $derived(reverse ? 'md:w-1/2 md:ps-8 mb-8 md:mb-0' : 'md:w-1/2 md:pe-8 mb-8 md:mb-0');
</script>

<section class={sectionClass}>
  <div class={textClass}>
    <div class="mb-6">
      <div class="title-line bg-yellow-thunder"></div>
      <h3 class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-4">
        {title}
      </h3>
      <p class="text-lg text-gray-600 dark:text-gray-400 mb-4">
        {description}
      </p>
    </div>
    <ul class="space-y-3 text-gray-700 dark:text-gray-300">
      {#each features as feature}
        <li class="flex items-start">
          <span class="text-indigo-500 dark:text-indigo-400 me-2">✓</span>
          <span>{@html feature}</span>
        </li>
      {/each}
    </ul>
  </div>
  <div class="md:w-1/2">
    <img src="{PathPrefix}{imageDark}" alt={imageAlt} class="rounded-lg shadow-lg mx-auto hidden dark:block" />
    <img src="{PathPrefix}{imageLight}" alt={imageAlt} class="rounded-lg shadow-lg mx-auto dark:hidden" />
  </div>
</section>

<script lang="ts">
  import 'shepherd.js/dist/css/shepherd.css';
  import Shepherd from 'shepherd.js';
  import { onMount } from 'svelte';

  export let steps = [];
  export let introText =
    'Care for a tour of the features offered on this page?';
  let hidden = false;

  const tour = new Shepherd.Tour({
    defaultStepOptions: {
      cancelIcon: {
        enabled: true,
      },
      classes: 'class-1 class-2',
      scrollTo: { behavior: 'smooth', block: 'center' },
    },
  });

  const tourCTA = () => {
    hideCTA();
    tour.start();
  };

  const hideCTA = () => {
    hidden = !hidden;
  };

  onMount(() => {
    steps.map(step => {
      tour.addStep(step);
    });
  });
</script>

{#if !hidden}
  <div
    class="flex items-center justify-between p-4 mb-8 text-sm font-semibold text-indigo-100 bg-indigo-500 rounded-lg shadow-md"
  >
    <div class="flex items-center">
      <span class="text-lg lg:text-xl">{introText}</span>
    </div>
    <div class="grid grid-cols-2 gap-2">
      <button on:click="{hideCTA}" class="text-white">No thanks.</button>
      <button
        on:click="{tourCTA}"
        class="bg-white text-indigo-500 py-2 px-3 rounded"
        >Yes please!
      </button>
    </div>
  </div>
{/if}

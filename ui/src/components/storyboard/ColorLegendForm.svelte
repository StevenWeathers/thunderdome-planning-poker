<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { onMount } from 'svelte';

  interface Props {
    handleLegendRevision?: any;
    toggleEditLegend?: any;
    colorLegend?: any;
    isFacilitator?: boolean;
  }

  let {
    handleLegendRevision = () => {},
    toggleEditLegend = () => {},
    colorLegend = $bindable([]),
    isFacilitator = false,
  }: Props = $props();

  function handleSubmit(event) {
    event.preventDefault();

    handleLegendRevision(colorLegend);
    toggleEditLegend();
  }

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleEditLegend} ariaLabel={$LL.modalStoryboardColorLegend()}>
  <form onsubmit={handleSubmit} name="colorLegend" class="space-y-4 pt-6">
    <h2 class="text-xl font-bold dark:text-gray-300">Story Color Legend</h2>
    <div class="space-y-2">
      {#each colorLegend as color, i}
        <div class="group">
          <label class="flex-1 min-w-0">
            <span class="sr-only">Color legend for {color.color}</span>
            {#if i === 0}
              <TextInput
                bind:this={focusInput}
                placeholder={$LL.legendRetroPlaceholder()}
                name="legend-{color.color}"
                disabled={!isFacilitator}
                value={colorLegend[i].legend}
              >
                {#snippet startElement()}
                  <div class="w-6 h-6 rounded bg-gray-400 colorcard-{color.color}"></div>
                {/snippet}
              </TextInput>
            {:else}
              <TextInput
                placeholder={$LL.legendRetroPlaceholder()}
                name="legend-{color.color}"
                disabled={!isFacilitator}
                value={colorLegend[i].legend}
              >
                {#snippet startElement()}
                  <div class="w-6 h-6 rounded bg-gray-400 colorcard-{color.color}"></div>
                {/snippet}
              </TextInput>
            {/if}
          </label>
        </div>
      {/each}
    </div>
    <div class="flex justify-end">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

<style>
  .colorcard-gray {
    @apply bg-gray-400;
  }

  .colorcard-red {
    @apply bg-red-400;
  }

  .colorcard-orange {
    @apply bg-orange-400;
  }

  .colorcard-yellow {
    @apply bg-yellow-400;
  }

  .colorcard-green {
    @apply bg-green-400;
  }

  .colorcard-teal {
    @apply bg-teal-400;
  }

  .colorcard-blue {
    @apply bg-blue-400;
  }

  .colorcard-indigo {
    @apply bg-indigo-400;
  }

  .colorcard-purple {
    @apply bg-purple-400;
  }

  .colorcard-pink {
    @apply bg-pink-400;
  }
</style>

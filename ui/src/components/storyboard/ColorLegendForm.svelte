<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';


  interface Props {
    handleLegendRevision?: any;
    toggleEditLegend?: any;
    colorLegend?: any;
    isFacilitator?: boolean;
  }

  let { handleLegendRevision = () => {}, toggleEditLegend = () => {}, colorLegend = $bindable([]), isFacilitator = false }: Props = $props();

  function handleSubmit(event) {
    event.preventDefault();

    handleLegendRevision(colorLegend);
    toggleEditLegend();
  }
</script>

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

<Modal closeModal={toggleEditLegend}>
  <form onsubmit={handleSubmit} name="colorLegend">
    <div class="mt-8 mb-4">
      {#each colorLegend as color, i}
        <div class="mb-1 flex">
          <span class="p-4 inline-block colorcard-{color.color}"></span>
          <TextInput
            bind:value="{colorLegend[i].legend}"
            placeholder={$LL.legendRetroPlaceholder()}
            name="legend-{color.color}"
            disabled={!isFacilitator}
          />
        </div>
      {/each}
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>

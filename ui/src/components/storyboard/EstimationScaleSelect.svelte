<script lang="ts">
  import SelectInput from '../forms/SelectInput.svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { StoryboardEstimationScale } from '../../types/storyboard';

  interface Props {
    value?: string;
    scales?: StoryboardEstimationScale[];
  }

  let { value = $bindable(''), scales = [] }: Props = $props();
</script>

<label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="estimationScale">
  {$LL.storyboardPointScale()}
</label>
<SelectInput bind:value id="estimationScale" name="estimationScale">
  <option value="">{$LL.storyboardPointScaleNone()}</option>
  {#each scales as scale (scale.id)}
    <option value={scale.id}>{scale.name} ({scale.values.join(', ')})</option>
  {/each}
</SelectInput>

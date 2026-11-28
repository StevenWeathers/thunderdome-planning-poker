<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { Crown, Lock } from '@lucide/svelte';
  import { onMount, untrack } from 'svelte';
  import EstimationScaleSelect from './EstimationScaleSelect.svelte';
  import type { StoryboardEstimationScale } from '../../types/storyboard';

  interface Props {
    toggleEditStoryboard?: any;
    handleStoryboardEdit?: any;
    storyboardName?: string;
    joinCode?: string;
    facilitatorCode?: string;
    estimationScale?: StoryboardEstimationScale | null;
    getEstimationScales?: () => Promise<StoryboardEstimationScale[]>;
    onEstimationScalesError?: () => void;
  }

  let {
    toggleEditStoryboard = () => {},
    handleStoryboardEdit = () => {},
    storyboardName = $bindable(''),
    joinCode = $bindable(''),
    facilitatorCode = $bindable(''),
    estimationScale = null,
    getEstimationScales = () => Promise.resolve([]),
    onEstimationScalesError = () => {},
  }: Props = $props();

  const initialScale = untrack(() => estimationScale);
  const initialScaleId = initialScale?.id ?? '';
  let selectedScaleId = $state(initialScaleId);
  // lists the current scale while the others load
  let estimationScales: StoryboardEstimationScale[] = $state(initialScale ? [initialScale] : []);

  let focusInput: any = $state();

  function saveStoryboard(e: Event) {
    e.preventDefault();

    const storyboard = {
      storyboardName,
      joinCode,
      facilitatorCode,
      // only sent when changed, so that saving other settings leaves the scale as it is
      ...(selectedScaleId !== initialScaleId ? { estimationScaleId: selectedScaleId } : {}),
    };

    handleStoryboardEdit(storyboard);
  }

  onMount(() => {
    focusInput?.focus();
    getEstimationScales()
      .then(scales => {
        estimationScales = scales;
      })
      .catch(onEstimationScalesError);
  });
</script>

<Modal
  closeModal={toggleEditStoryboard}
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
  ariaLabel={$LL.modalEditStoryboard()}
>
  <form onsubmit={saveStoryboard} name="createStoryboard">
    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="storyboardName">
        {$LL.storyboardName()}
      </label>
      <div class="control">
        <TextInput
          name="storyboardName"
          bind:value={storyboardName}
          placeholder={$LL.storyboardNamePlaceholder()}
          id="storyboardName"
          bind:this={focusInput}
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2" for="joinCode">
        {$LL.passCode()}
      </label>
      <div class="control">
        <TextInput
          name="joinCode"
          bind:value={joinCode}
          placeholder={$LL.optionalPasscodePlaceholder()}
          id="joinCode"
          icon={Lock}
        />
      </div>
    </div>

    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="facilitatorCode">
        {$LL.facilitatorCodeOptional()}
      </label>
      <div class="control">
        <TextInput
          name="facilitatorCode"
          bind:value={facilitatorCode}
          placeholder={$LL.facilitatorCodePlaceholder()}
          id="facilitatorCode"
          icon={Crown}
        />
      </div>
    </div>

    <div class="mb-4">
      <EstimationScaleSelect bind:value={selectedScaleId} scales={estimationScales} />
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>

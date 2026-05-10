<script lang="ts">
  import { onMount } from 'svelte';

  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import ColorLegendTemplateFields from './ColorLegendTemplateFields.svelte';

  import type { ColorLegend } from '../../types/storyboard';
  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    toggleUpdate?: () => void;
    handleUpdate?: () => void;
    apiPrefix?: string;
    xfetch: ApiClient;
    notifications: NotificationService;
    templateId?: string;
    name?: string;
    description?: string;
    colorLegend?: Array<ColorLegend>;
  }

  let {
    toggleUpdate = () => {},
    handleUpdate = () => {},
    apiPrefix = '/api',
    xfetch,
    notifications,
    templateId = '',
    name = $bindable(''),
    description = $bindable(''),
    colorLegend = $bindable([]),
  }: Props = $props();

  let focusInput: any;

  function toggleClose() {
    toggleUpdate();
  }

  function onSubmit(event: Event) {
    event.preventDefault();

    xfetch(`${apiPrefix}/color-legend-templates/${templateId}`, {
      body: {
        name,
        description,
        colorLegend,
      },
      method: 'put',
    })
      .then(res => res.json())
      .then(() => {
        notifications.success('Color legend template successfully updated');
        handleUpdate();
      })
      .catch(() => {
        notifications.danger('Error updating color legend template');
      });
  }

  let updateDisabled = $derived(name.trim() === '');

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleClose} ariaLabel="Update color legend template" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <form onsubmit={onSubmit} name="updateColorLegendTemplate" class="space-y-4">
    <div>
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="templateName">Name</label>
      <TextInput
        bind:value={name}
        bind:this={focusInput}
        id="templateName"
        name="templateName"
        placeholder="Enter template name"
        required
      />
    </div>

    <div>
      <label class="block text-gray-700 font-bold mb-2 dark:text-gray-400" for="templateDescription">Description</label>
      <TextInput
        bind:value={description}
        id="templateDescription"
        name="templateDescription"
        placeholder="Enter template description"
      />
    </div>

    <ColorLegendTemplateFields bind:colorLegend />

    <div class="text-right">
      <SolidButton type="submit" disabled={updateDisabled}>Save Template</SolidButton>
    </div>
  </form>
</Modal>

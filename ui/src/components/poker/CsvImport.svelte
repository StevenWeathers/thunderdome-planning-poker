<script lang="ts">
  import HollowButton from '../global/HollowButton.svelte';
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    notifications: any;
    handlePlanAdd?: any;
  }

  let { notifications, handlePlanAdd = () => {} }: Props = $props();

  const allowCsvImport = AppConfig.AllowCsvImport;

  let plans = [];

  function uploadFile() {
    let file = this.files[0];
    if (!file) {
      return;
    }
    if (file.type !== 'text/csv') {
      notifications.danger($LL.importCsvFileBadFileTypeError());
      return;
    }

    let reader = new FileReader();

    reader.readAsText(file);

    reader.onload = () => {
      try {
        const content = reader.result;
        const items = parseCsvFile(content);
        if (items) {
          const totalItems = items.length;
          for (let i = 0; i < totalItems; i++) {
            const item = items[i];
            const plan = extractPlanData(item);
            plans.push(plan);
            handlePlanAdd(plan);
          }
        }
      } catch (e) {
        notifications.danger($LL.importCsvFileReadFileError());
      }
    };

    reader.onerror = () => {
      notifications.danger($LL.importCsvFileReadFileError());
    };
  }

  function parseCsvFile(content) {
    const lines = content.split('\n');
    const items = [];

    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim();
      if (line) {
        items.push(line);
      }
    }

    return items;
  }

  function extractPlanData(item) {
    const fields = item.split(',');
    const plan = {
      type: fields[0].trim(),
      planName: fields[1].trim(),
      referenceId: fields[2].trim(),
      link: fields[3].trim(),
      description: fields[4].trim(),
      acceptanceCriteria: fields[5] ? fields[5].trim() : '',
    };

    return plan;
  }
</script>

{#if allowCsvImport}
  <HollowButton
    type="label"
    additionalClasses="me-2"
    color="purple"
    labelFor="csvimport"
  >
    {$LL.selectFile()}
    <input
      type="file"
      onchange={uploadFile}
      class="hidden"
      id="csvimport"
      accept=".csv"
    />
  </HollowButton>
{/if}

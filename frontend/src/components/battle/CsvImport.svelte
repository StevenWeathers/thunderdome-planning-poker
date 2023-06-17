<script>
    import HollowButton from '../HollowButton.svelte'
    import { AppConfig } from '../../config.js'
    import { _ } from '../../i18n.js'

    export let notifications
    export let eventTag = () => {}
    export let handlePlanAdd = () => {}

    const allowCsvImport = AppConfig.AllowCsvImport

    let plans = []

    function uploadFile() {
      let file = this.files[0]
      if (!file) {
        return
      }
      if (file.type !== 'text/csv') {
        notifications.danger($_('importCsvFileBadFileTypeError'))
        eventTag('Csv_import_failed', 'battle', `file.type not text/csv or application/vnd.ms-excel`)
        return
      }

      let reader = new FileReader()

      reader.readAsText(file)

      reader.onload = () => {
        try {
          const content = reader.result;
          const items = parseCsvFile(content);
          if (items) {
            const totalItems = items.length
            for (let i = 0; i < totalItems; i++) {
              const item = items[i]
              const plan = extractPlanData(item);
              plans.push(plan)
              handlePlanAdd(plan);
            }
            eventTag(
              'Csv_import_success',
              'battle',
              `total stories imported: ${totalItems}`,
            )
          }
        } catch (e) {
          notifications.danger($_('importCsvFileReadFileError'))
          eventTag('Csv_import_failed', 'battle', `error reading file`)
        }
      }

      reader.onerror = () => {
        notifications.danger($_('importCsvFileReadFileError'))
        eventTag('Csv_import_failed', 'battle', `error reading file`)
      }
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
        acceptanceCriteria: fields[5] ? fields[5].trim() : ''
      };

      return plan;
    }
  </script>

  {#if allowCsvImport}
    <HollowButton
      type="label"
      additionalClasses="rtl:ml-2 ltr:mr-2"
      color="purple"
      labelFor="csvimport"
    >
      {$_('importCsv')}
      <input
        type="file"
        on:change="{uploadFile}"
        class="hidden"
        id="csvimport"
        accept=".csv"
      />
    </HollowButton>
  {/if}

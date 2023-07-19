<script lang="ts">
  import Modal from '../Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import SolidButton from '../SolidButton.svelte';
  import { AppConfig } from '../../config';

  export let notifications;
  export let eventTag;
  export let toggleImport = () => {};
  export let handlePlanAdd = handleAdd => {};

  const handleAdd = newPlan => {
    handlePlanAdd(newPlan);
    toggleImport();
  };
</script>

<Modal closeModal="{toggleImport}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <div class="mt-8 mb-4">
    <div class="mb-4">
      <h3 class="font-bold mb-2 dark:text-gray-300 text-lg">
        {$LL.importCsv({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>
      <p class="dark:text-gray-400 mb-2">
        The CSV file must include all the following fields with no header row:
      </p>
      <div class="mb-2 whitespace-nowrap dark:text-gray-300">
        Type,Title,ReferenceId,Link,Description,AcceptanceCriteria
      </div>
      <CsvImport
        handlePlanAdd="{handleAdd}"
        notifications="{notifications}"
        eventTag="{eventTag}"
        testid="plans-Csvimport"
      />
    </div>
    <div class="mb-4">
      <h3 class="font-bold mb-2 dark:text-gray-300 text-lg">
        {$LL.importJiraXML({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>
      <JiraImport
        handlePlanAdd="{handleAdd}"
        notifications="{notifications}"
        eventTag="{eventTag}"
        testid="plans-importjira"
      />
    </div>
  </div>
  <div class="text-right">
    <SolidButton onClick="{toggleImport}">
      {$LL.cancel()}
    </SolidButton>
  </div>
</Modal>

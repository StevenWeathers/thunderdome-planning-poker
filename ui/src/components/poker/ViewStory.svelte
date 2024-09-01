<script lang="ts">
  import ExternalLinkIcon from '../icons/ExternalLinkIcon.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import NoSymbolIcon from '../icons/NoSymbol.svelte';
  import DoubleChevronDown from '../icons/DoubleChevronDown.svelte';
  import DoubleChevronUp from '../icons/DoubleChevronUp.svelte';
  import ChevronDown from '../icons/ChevronDown.svelte';
  import ChevronUp from '../icons/ChevronUp.svelte';
  import Bars2 from '../icons/Bars2.svelte';

  export let togglePlanView = () => {};

  export let planName = '';
  export let planType = '';
  export let referenceId = '';
  export let planLink = '';
  export let description = '';
  export let acceptanceCriteria = '';
  export let priority = 99;

  const priorities = {
    99: {
      name: '',
      icon: false,
    },
    1: {
      name: $LL.planPriorityBlocker(),
      icon: NoSymbolIcon,
    },
    2: {
      name: $LL.planPriorityHighest(),
      icon: DoubleChevronUp,
    },
    3: {
      name: $LL.planPriorityHigh(),
      icon: ChevronUp,
    },
    4: {
      name: $LL.planPriorityMedium(),
      icon: Bars2,
    },
    5: {
      name: $LL.planPriorityLow(),
      icon: ChevronDown,
    },
    6: {
      name: $LL.planPriorityLowest(),
      icon: DoubleChevronDown,
    },
  };
</script>

<Modal closeModal="{togglePlanView}" widthClasses="md:w-2/3 lg:w-3/5">
  <div class="mb-4 dark:text-white">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planType()}
    </div>
    {planType}
  </div>
  <div class="mb-4 dark:text-white">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planName()}
    </div>
    {planName}
  </div>
  <div class="mb-4 dark:text-white">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planReferenceId()}
    </div>
    {referenceId}
  </div>
  <div class="mb-4">
    <div class="font-bold mb-2 dark:text-gray-400">{$LL.planLink()}</div>
    {#if planLink !== ''}
      <a
        href="{planLink}"
        target="_blank"
        class="text-blue-800 hover:text-blue-600 dark:text-sky-400 dark:hover:text-sky-600"
      >
        <ExternalLinkIcon />
        {planLink}
      </a>
    {/if}
  </div>
  <div class="mb-4 dark:text-white">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planPriority()}
    </div>
    <svelte:component this="{priorities[priority].icon}" />{priorities[priority]
      .name}
  </div>
  <div class="mb-4">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planDescription()}
    </div>
    <div class="unreset dark:text-white">
      {@html description}
    </div>
  </div>
  <div class="mb-4">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planAcceptanceCriteria()}
    </div>
    <div class="unreset dark:text-white">
      {@html acceptanceCriteria}
    </div>
  </div>
</Modal>

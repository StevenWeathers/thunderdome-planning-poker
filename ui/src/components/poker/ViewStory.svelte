<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import {
    Ban,
    ChevronDown,
    ChevronsDown,
    ChevronsUp,
    ChevronUp,
    ExternalLink,
  } from 'lucide-svelte';
  import Bars2 from '../icons/Bars2.svelte';


  interface Props {
    togglePlanView?: any;
    planName?: string;
    planType?: string;
    referenceId?: string;
    planLink?: string;
    description?: string;
    acceptanceCriteria?: string;
    priority?: number;
  }

  let {
    togglePlanView = () => {},
    planName = '',
    planType = '',
    referenceId = '',
    planLink = '',
    description = '',
    acceptanceCriteria = '',
    priority = 99
  }: Props = $props();

  const priorities = {
    99: {
      name: '',
      icon: false,
    },
    1: {
      name: $LL.planPriorityBlocker(),
      icon: Ban,
    },
    2: {
      name: $LL.planPriorityHighest(),
      icon: ChevronsUp,
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
      icon: ChevronsDown,
    },
  };
</script>

<Modal closeModal={togglePlanView} widthClasses="md:w-2/3 lg:w-3/5">
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
        <ExternalLink class="inline-block" />
        {planLink}
      </a>
    {/if}
  </div>
  {@const SvelteComponent = priorities[priority].icon}
  <div class="mb-4 dark:text-white">
    <div class="font-bold mb-2 dark:text-gray-400">
      {$LL.planPriority()}
    </div>
    <SvelteComponent
      class="inline-block w-6 h-6"
    />{priorities[priority].name}
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

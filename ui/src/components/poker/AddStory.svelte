<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import {
    Ban,
    ChevronDown,
    ChevronsDown,
    ChevronsUp,
    ChevronUp,
  } from 'lucide-svelte';
  import Bars2 from '../icons/Bars2.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Editor from '../forms/Editor.svelte';
  import { onMount } from 'svelte';

  import type { NotificationService } from '../../types/notifications';

  // going by common Jira issue types for now
  const planTypes = [
    $LL.planTypeStory(),
    $LL.planTypeBug(),
    $LL.planTypeSpike(),
    $LL.planTypeEpic(),
    $LL.planTypeTask(),
    $LL.planTypeSubtask(),
  ];

  // going by common Jira issue priorities for now
  const priorities = [
    { name: $LL.planPriorityBlocker(), value: 1, icon: Ban },
    {
      name: $LL.planPriorityHighest(),
      value: 2,
      icon: ChevronsUp,
    },
    { name: $LL.planPriorityHigh(), value: 3, icon: ChevronUp },
    { name: $LL.planPriorityMedium(), value: 4, icon: Bars2 },
    { name: $LL.planPriorityLow(), value: 5, icon: ChevronDown },
    {
      name: $LL.planPriorityLowest(),
      value: 6,
      icon: ChevronsDown,
    },
  ];

  interface Props {
    handlePlanAdd?: any;
    toggleAddPlan?: any;
    handlePlanRevision?: any;
    notifications: NotificationService;
    planId?: string;
    planName?: string;
    planType?: any;
    referenceId?: string;
    planLink?: string;
    description?: string;
    acceptanceCriteria?: string;
    priority?: number;
  }

  let {
    handlePlanAdd = () => {},
    toggleAddPlan = () => {},
    handlePlanRevision = () => {},
    notifications,
    planId = '',
    planName = $bindable(''),
    planType = $bindable($LL.planTypeStory()),
    referenceId = $bindable(''),
    planLink = $bindable(''),
    description = $bindable(''),
    acceptanceCriteria = $bindable(''),
    priority = $bindable(99)
  }: Props = $props();

  /** @type {TextInput} */
  let planNameTextInput = $state();

  const isAbsolute = new RegExp('^([a-z]+://|//)', 'i');
  let descriptionExpanded = $state(false);
  let acceptanceExpanded = $state(false);

  function handleSubmit(event) {
    event.preventDefault();
    let invalidPlan = false;

    if (planLink !== '' && !isAbsolute.test(planLink)) {
      invalidPlan = true;
      notifications.danger($LL.planLinkInvalid());
    }

    const plan = {
      planName,
      type: planType,
      referenceId,
      link: planLink,
      description,
      acceptanceCriteria,
      priority,
    };

    if (!invalidPlan) {
      if (planId === '') {
        handlePlanAdd(plan);
      } else {
        plan.planId = planId;
        handlePlanRevision(plan);
      }

      toggleAddPlan();
    }
  }

  // Focus the plan name input field when the modal is opened
  onMount(() => {
    planNameTextInput.focus();
  });
</script>

<Modal closeModal={toggleAddPlan} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <form onsubmit={handleSubmit} name="addPlan">
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="planType">
        {$LL.planType()}
      </label>
      <SelectInput
        name="planType"
        id="planType"
        bind:value="{planType}"
        required
      >
        <option value="" disabled>
          {$LL.planTypePlaceholder()}
        </option>
        {#each planTypes as pType}
          <option value="{pType}">{pType}</option>
        {/each}
      </SelectInput>
    </div>
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="planName">
        {$LL.planName()}
      </label>
      <TextInput
        id="planName"
        name="planName"
        bind:this="{planNameTextInput}"
        bind:value="{planName}"
        placeholder={$LL.planNamePlaceholder()}
      />
    </div>
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="referenceId">
        {$LL.planReferenceId()}
      </label>
      <TextInput
        id="referenceId"
        name="referenceId"
        bind:value="{referenceId}"
        placeholder={$LL.planReferenceIdPlaceholder()}
      />
    </div>
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="planLink">
        {$LL.planLink()}
      </label>
      <TextInput
        id="planLink"
        name="planLink"
        bind:value="{planLink}"
        placeholder={$LL.planLinkPlaceholder()}
      />
    </div>
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="priority">
        {$LL.planPriority()}
      </label>
      <SelectInput name="priority" id="priority" bind:value="{priority}">
        <option value="{99}" disabled>
          {$LL.planPriorityPlaceholder()}
        </option>
        {#each priorities as p}
          <option value="{p.value}">
            <p.icon
              class="inline-block w-6 h-6"
            />{p.name}</option
          >
        {/each}
      </SelectInput>
    </div>
    <div>
      <div class="font-bold mb-2">
        <button
          onclick={e => {
            e.preventDefault();
            descriptionExpanded = !descriptionExpanded;
          }}
          class="inline-block align-baseline text-sm
                        text-blue-700 dark:text-sky-400 hover:text-blue-800 dark:hover:text-sky-600 bg-transparent
                        border-transparent me-1 font-bold text-xl"
          type="button"
        >
          {#if descriptionExpanded}-{:else}+{/if}
        </button>
        <span class="dark:text-gray-400">{$LL.planDescription()}</span>
      </div>
      {#if descriptionExpanded}
        <div class="mb-2">
          <div class="bg-white">
            <Editor
              content={description}
              placeholder={$LL.planDescriptionPlaceholder()}
              id="storyDescription"
              handleTextChange={c => (description = c)}
            />
          </div>
        </div>
      {/if}
    </div>
    <div>
      <div class="font-bold mb-2">
        <button
          onclick={e => {
            e.preventDefault();
            acceptanceExpanded = !acceptanceExpanded;
          }}
          class="inline-block align-baseline text-sm
                        text-blue-700 dark:text-sky-400 hover:text-blue-800 dark:hover:text-sky-600 bg-transparent
                        border-transparent me-1 font-bold text-xl"
          type="button"
        >
          {#if acceptanceExpanded}-{:else}+{/if}
        </button>
        <span class="dark:text-gray-400">{$LL.planAcceptanceCriteria()}</span>
      </div>
      {#if acceptanceExpanded}
        <div class="mb-2">
          <div class="bg-white">
            <Editor
              content={acceptanceCriteria}
              placeholder={$LL.planAcceptanceCriteriaPlaceholder()}
              id="acceptanceCriteria"
              handleTextChange={c => (acceptanceCriteria = c)}
            />
          </div>
        </div>
      {/if}
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit" testid="plan-save">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>

<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import NoSymbol from '../icons/NoSymbol.svelte';
  import DoubleChevronUp from '../icons/DoubleChevronUp.svelte';
  import ChevronUp from '../icons/ChevronUp.svelte';
  import Bars2 from '../icons/Bars2.svelte';
  import ChevronDown from '../icons/ChevronDown.svelte';
  import DoubleChevronDown from '../icons/DoubleChevronDown.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Editor from '../forms/Editor.svelte';
  import { onMount } from 'svelte';

  export let handlePlanAdd = () => {};
  export let toggleAddPlan = () => {};
  export let handlePlanRevision = () => {};
  export let eventTag = () => {};
  export let notifications;

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
    { name: $LL.planPriorityBlocker(), value: 1, icon: NoSymbol },
    {
      name: $LL.planPriorityHighest(),
      value: 2,
      icon: DoubleChevronUp,
    },
    { name: $LL.planPriorityHigh(), value: 3, icon: ChevronUp },
    { name: $LL.planPriorityMedium(), value: 4, icon: Bars2 },
    { name: $LL.planPriorityLow(), value: 5, icon: ChevronDown },
    {
      name: $LL.planPriorityLowest(),
      value: 6,
      icon: DoubleChevronDown,
    },
  ];

  export let planId = '';
  export let planName = '';
  export let planType = $LL.planTypeStory();
  export let referenceId = '';
  export let planLink = '';
  export let description = '';
  export let acceptanceCriteria = '';
  export let priority = 99;

  /** @type {TextInput} */
  let planNameTextInput;

  const isAbsolute = new RegExp('^([a-z]+://|//)', 'i');
  let descriptionExpanded = false;
  let acceptanceExpanded = false;

  function handleSubmit(event) {
    event.preventDefault();
    let invalidPlan = false;

    if (planLink !== '' && !isAbsolute.test(planLink)) {
      invalidPlan = true;
      notifications.danger($LL.planLinkInvalid());
      eventTag('plan_add_invalid_link', 'battle', ``);
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

<Modal closeModal="{toggleAddPlan}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <form on:submit="{handleSubmit}" name="addPlan">
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="planType">
        {$LL.planType({ friendly: AppConfig.FriendlyUIVerbs })}
      </label>
      <SelectInput
        name="planType"
        id="planType"
        bind:value="{planType}"
        required
      >
        <option value="" disabled>
          {$LL.planTypePlaceholder({
            friendly: AppConfig.FriendlyUIVerbs,
          })}
        </option>
        {#each planTypes as pType}
          <option value="{pType}">{pType}</option>
        {/each}
      </SelectInput>
    </div>
    <div class="mb-4">
      <label class="block font-bold mb-2 dark:text-gray-400" for="planName">
        {$LL.planName({ friendly: AppConfig.FriendlyUIVerbs })}
      </label>
      <TextInput
        id="planName"
        name="planName"
        bind:this="{planNameTextInput}"
        bind:value="{planName}"
        placeholder="{$LL.planNamePlaceholder({
          friendly: AppConfig.FriendlyUIVerbs,
        })}"
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
        placeholder="{$LL.planReferenceIdPlaceholder()}"
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
        placeholder="{$LL.planLinkPlaceholder({
          friendly: AppConfig.FriendlyUIVerbs,
        })}"
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
            <svelte:component this="{p.icon}" />{p.name}</option
          >
        {/each}
      </SelectInput>
    </div>
    <div>
      <div class="font-bold mb-2">
        <button
          on:click="{e => {
            e.preventDefault();
            descriptionExpanded = !descriptionExpanded;
          }}"
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
              content="{description}"
              placeholder="{$LL.planDescriptionPlaceholder()}"
              id="storyDescription"
              handleTextChange="{c => (description = c)}"
            />
          </div>
        </div>
      {/if}
    </div>
    <div>
      <div class="font-bold mb-2">
        <button
          on:click="{e => {
            e.preventDefault();
            acceptanceExpanded = !acceptanceExpanded;
          }}"
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
              content="{acceptanceCriteria}"
              placeholder="{$LL.planAcceptanceCriteriaPlaceholder()}"
              id="acceptanceCriteria"
              handleTextChange="{c => (acceptanceCriteria = c)}"
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

<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import ExternalLinkIcon from '../icons/ExternalLinkIcon.svelte';
  import HollowButton from '../HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import NoSymbolIcon from '../icons/NoSymbol.svelte';
  import DoubleChevronUp from '../icons/DoubleChevronUp.svelte';
  import ChevronUp from '../icons/ChevronUp.svelte';
  import Bars2 from '../icons/Bars2.svelte';
  import ChevronDown from '../icons/ChevronDown.svelte';
  import DoubleChevronDown from '../icons/DoubleChevronDown.svelte';
  import AddPlan from './AddStory.svelte';
  import ViewPlan from './ViewStory.svelte';
  import { AppConfig } from '../../config';
  import ImportModal from './ImportModal.svelte';

  export let plans = [];
  export let isLeader = false;
  export let sendSocketEvent = (event: string, value: string) => {};
  export let eventTag;
  export let notifications;

  let defaultPlan = {
    id: '',
    name: '',
    type: $LL.planTypeStory(),
    referenceId: '',
    link: '',
    description: '',
    acceptanceCriteria: '',
    priority: 99,
  };

  let priorities = {
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

  let showAddPlan = false;
  let showViewPlan = false;
  let selectedPlan = { ...defaultPlan };
  let showCompleted = false;
  let showImport = false;

  const toggleImport = () => {
    showImport = !showImport;
  };

  const toggleAddPlan = planId => () => {
    if (planId) {
      selectedPlan = plans.find(p => p.id === planId);

      eventTag('plan_show_edit', 'battle', ``);
    } else {
      selectedPlan = { ...defaultPlan };

      eventTag('plan_show_add', 'battle', ``);
    }
    showAddPlan = !showAddPlan;
  };

  const togglePlanView = planId => () => {
    if (planId) {
      selectedPlan = plans.find(p => p.id === planId);
      eventTag('plan_show_view', 'battle', ``);
    } else {
      selectedPlan = { ...defaultPlan };

      eventTag('plan_unshow_view', 'battle', ``);
    }
    showViewPlan = !showViewPlan;
  };

  const handlePlanAdd = newPlan => {
    sendSocketEvent('add_plan', JSON.stringify(newPlan));
    eventTag('plan_add', 'battle', '');
  };

  const activatePlan = id => () => {
    sendSocketEvent('activate_plan', id);
    eventTag('plan_activate', 'battle', '');
  };

  const handlePlanRevision = updatedPlan => {
    sendSocketEvent('revise_plan', JSON.stringify(updatedPlan));
    eventTag('plan_revise', 'battle', '');
  };

  const handlePlanDeletion = planId => () => {
    sendSocketEvent('burn_plan', planId);
    eventTag('plan_burn', 'battle', '');
  };

  const toggleShowCompleted = show => () => {
    showCompleted = show;
    eventTag('plans_show', 'battle', `completed: ${show}`);
  };

  $: pointedPlans = plans.filter(p => p.points !== '');
  $: totalPoints = pointedPlans.reduce((previousValue, currentValue) => {
    var currentPoints =
      currentValue.points === '1/2' ? 0.5 : parseInt(currentValue.points);
    return isNaN(currentPoints) ? previousValue : previousValue + currentPoints;
  }, 0);
  $: unpointedPlans = plans.filter(p => p.points === '');

  $: plansToShow = showCompleted ? pointedPlans : unpointedPlans;

  // event handlers
  function handleDndConsider(e) {
    plans = e.detail.items;
  }

  function handleDndFinalize(e) {
    const storyId = e.detail.info.id;

    plans = e.detail.items;

    const matchedStory = plans.find(i => i.id === storyId);

    if (matchedStory) {
      // determine what story to place story before
      const matchedStoryIndex = plans.indexOf(matchedStory);
      const sibling = plans[matchedStoryIndex + 1];
      const placeBefore = sibling ? sibling.id : '';

      sendSocketEvent(
        'story_arrange',
        JSON.stringify({
          story_id: storyId,
          before_story_id: placeBefore,
        }),
      );
      eventTag('story_arrange', 'battle', '');
    }
  }
</script>

<div class="shadow-lg mb-4">
  <div class="flex items-center bg-gray-200 dark:bg-gray-700 p-4 rounded-t-lg">
    <div class="w-1/3">
      <h3
        class="text-3xl leading-tight font-semibold font-rajdhani uppercase dark:text-white"
      >
        {$LL.plans({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>
    </div>
    <div class="w-2/3 text-right">
      {#if isLeader}
        <HollowButton onClick="{toggleImport}" color="blue">
          {$LL.importPlans({ friendly: AppConfig.FriendlyUIVerbs })}
        </HollowButton>
        <HollowButton onClick="{toggleAddPlan()}" testid="plans-add">
          {$LL.planAdd({ friendly: AppConfig.FriendlyUIVerbs })}
        </HollowButton>
      {/if}
    </div>
  </div>

  <ul
    class="flex border-b border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800"
  >
    <li class="-mb-px {showCompleted ? '' : 'me-1'}">
      <button
        class="{showCompleted
          ? 'hover:text-blue-600 text-blue-400 dark:hover:text-sky-300 dark:text-sky-600'
          : 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
        on:click="{toggleShowCompleted(false)}"
        data-testid="plans-unpointed"
      >
        {$LL.unpointed({ count: unpointedPlans.length })}
      </button>
    </li>
    <li class="me-1 {showCompleted ? 'me-1' : ''}">
      <button
        class="{showCompleted
          ? 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'
          : 'hover:text-blue-600 dark:hover:text-sky-300 text-blue-400 dark:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
        on:click="{toggleShowCompleted(true)}"
        data-testid="plans-pointed"
      >
        {$LL.pointed({ count: pointedPlans.length })}
      </button>
    </li>
  </ul>

  <div
    use:dndzone="{{
      items: plans,
      type: 'story',
      dropTargetStyle: '',
      dropTargetClasses: [
        'outline',
        'outline-2',
        'outline-indigo-500',
        'dark:outline-yellow-400',
      ],
      dragDisabled: !isLeader,
    }}"
    on:consider="{handleDndConsider}"
    on:finalize="{handleDndFinalize}"
  >
    {#each plansToShow as plan (plan.id)}
      <div
        class="relative flex flex-wrap items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800{isLeader
          ? ' cursor-pointer'
          : ''}"
        data-testid="plan"
        data-storyid="{plan.id}"
      >
        <div class="w-full lg:w-1/3 mb-4 lg:mb-0">
          <div class="inline-block font-bold align-middle dark:text-white">
            {#if plan.link !== ''}
              <a
                href="{plan.link}"
                target="_blank"
                class="text-blue-800 dark:text-sky-400"
              >
                <ExternalLinkIcon />
              </a>
              &nbsp;
            {/if}
            <div
              class="inline-block text-sm text-gray-500 dark:text-gray-300
                        border-gray-300 border px-1 rounded"
              data-testid="plan-type"
            >
              {plan.type}
            </div>
            &nbsp;
            {#if plan.referenceId}[{plan.referenceId}]&nbsp;{/if}
            <svelte:component this="{priorities[plan.priority].icon}" />
            <span data-testid="plan-name">{plan.name}</span>
          </div>
          &nbsp;
          {#if plan.points !== ''}
            <div
              class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded ms-2"
              data-testid="plan-points"
            >
              {plan.points}
            </div>
          {/if}
        </div>
        <div class="w-full lg:w-2/3 text-right">
          <HollowButton
            color="blue"
            onClick="{togglePlanView(plan.id)}"
            testid="plan-view"
          >
            {$LL.view()}
          </HollowButton>
          {#if isLeader}
            {#if !plan.active}
              <HollowButton
                color="red"
                onClick="{handlePlanDeletion(plan.id)}"
                testid="plan-delete"
              >
                {$LL.delete()}
              </HollowButton>
            {/if}
            <HollowButton
              color="purple"
              onClick="{toggleAddPlan(plan.id)}"
              testid="plan-edit"
            >
              {$LL.edit()}
            </HollowButton>
            {#if !plan.active}
              <HollowButton
                onClick="{activatePlan(plan.id)}"
                testid="plan-activate"
              >
                {$LL.activate()}
              </HollowButton>
            {/if}
          {/if}
        </div>
      </div>
      {#if plan[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
        <div
          class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible opacity-50 cursor-pointer flex flex-wrap items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800"
          data-testid="plan"
          data-storyid="{plan.id}"
        >
          <div class="w-full lg:w-1/3 mb-4 lg:mb-0">
            <div class="inline-block font-bold align-middle dark:text-white">
              {#if plan.link !== ''}
                <a
                  href="{plan.link}"
                  target="_blank"
                  class="text-blue-800 dark:text-sky-400"
                >
                  <ExternalLinkIcon />
                </a>
                &nbsp;
              {/if}
              <div
                class="inline-block text-sm text-gray-500 dark:text-gray-300
                        border-gray-300 border px-1 rounded"
                data-testid="plan-type"
              >
                {plan.type}
              </div>
              &nbsp;
              {#if plan.referenceId}[{plan.referenceId}]&nbsp;{/if}
              <svelte:component this="{priorities[plan.priority].icon}" />
              <span data-testid="plan-name">{plan.name}</span>
            </div>
            &nbsp;
            {#if plan.points !== ''}
              <div
                class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded ms-2"
                data-testid="plan-points"
              >
                {plan.points}
              </div>
            {/if}
          </div>
          <div class="w-full lg:w-2/3 text-right">
            <HollowButton
              color="blue"
              onClick="{togglePlanView(plan.id)}"
              testid="plan-view"
            >
              {$LL.view()}
            </HollowButton>
            {#if isLeader}
              {#if !plan.active}
                <HollowButton
                  color="red"
                  onClick="{handlePlanDeletion(plan.id)}"
                  testid="plan-delete"
                >
                  {$LL.delete()}
                </HollowButton>
              {/if}
              <HollowButton
                color="purple"
                onClick="{toggleAddPlan(plan.id)}"
                testid="plan-edit"
              >
                {$LL.edit()}
              </HollowButton>
              {#if !plan.active}
                <HollowButton
                  onClick="{activatePlan(plan.id)}"
                  testid="plan-activate"
                >
                  {$LL.activate()}
                </HollowButton>
              {/if}
            {/if}
          </div>
        </div>
      {/if}
    {/each}
  </div>
  {#if showCompleted && totalPoints}
    <div
      class="flex flex-wrap items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800"
    >
      <div class="w-full lg:w-2/3 mb-4 lg:mb-0">
        <div class="inline-block font-bold align-middle dark:text-gray-300">
          {$LL.totalPoints()}:
        </div>
        &nbsp;
        <div
          class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded ms-2"
        >
          {totalPoints}
        </div>
      </div>
    </div>
  {/if}
</div>

{#if showAddPlan}
  <AddPlan
    handlePlanAdd="{handlePlanAdd}"
    toggleAddPlan="{toggleAddPlan()}"
    handlePlanRevision="{handlePlanRevision}"
    planId="{selectedPlan.id}"
    planName="{selectedPlan.name}"
    planType="{selectedPlan.type}"
    referenceId="{selectedPlan.referenceId}"
    planLink="{selectedPlan.link}"
    description="{selectedPlan.description}"
    acceptanceCriteria="{selectedPlan.acceptanceCriteria}"
    priority="{selectedPlan.priority}"
    notifications="{notifications}"
    eventTag="{eventTag}"
  />
{/if}

{#if showViewPlan}
  <ViewPlan
    togglePlanView="{togglePlanView()}"
    planName="{selectedPlan.name}"
    planType="{selectedPlan.type}"
    referenceId="{selectedPlan.referenceId}"
    planLink="{selectedPlan.link}"
    description="{selectedPlan.description}"
    acceptanceCriteria="{selectedPlan.acceptanceCriteria}"
    priority="{selectedPlan.priority}"
  />
{/if}

{#if showImport}
  <ImportModal
    notifications="{notifications}"
    eventTag="{eventTag}"
    toggleImport="{toggleImport}"
    handlePlanAdd="{handlePlanAdd}"
  />
{/if}

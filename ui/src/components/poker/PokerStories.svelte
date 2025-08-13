<script lang="ts">
  import { dndzone, SHADOW_ITEM_MARKER_PROPERTY_NAME } from 'svelte-dnd-action';
  import HollowButton from '../global/HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import {
    Ban,
    ChevronDown,
    ChevronsDown,
    ChevronsUp,
    ChevronUp,
    ExternalLink,
    Grip
  } from 'lucide-svelte';
  import Bars2 from '../icons/Bars2.svelte';
  import AddPlan from './AddStory.svelte';
  import ViewPlan from './ViewStory.svelte';
  import ImportModal from './ImportModal.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    plans?: any;
    isFacilitator?: boolean;
    sendSocketEvent?: any;
    notifications: NotificationService;
    xfetch: ApiClient;
    gameId?: string;
  }

  let {
    plans = $bindable([]),
    isFacilitator = false,
    sendSocketEvent = (event: string, value: string) => {},
    notifications,
    xfetch,
    gameId = ''
  }: Props = $props();

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

  let showAddPlan = $state(false);
  let showViewPlan = $state(false);
  let selectedPlan = $state({ ...defaultPlan });
  let storysShow = $state('unpointed');
  let showImport = $state(false);

  const toggleImport = () => {
    showImport = !showImport;
  };

  const toggleAddPlan = planId => () => {
    if (planId) {
      selectedPlan = plans.find(p => p.id === planId);
    } else {
      selectedPlan = { ...defaultPlan };
    }
    showAddPlan = !showAddPlan;
  };

  const togglePlanView = planId => () => {
    if (planId) {
      selectedPlan = plans.find(p => p.id === planId);
    } else {
      selectedPlan = { ...defaultPlan };
    }
    showViewPlan = !showViewPlan;
  };

  const handlePlanAdd = newPlan => {
    sendSocketEvent('add_plan', JSON.stringify(newPlan));
  };

  const activatePlan = id => () => {
    sendSocketEvent('activate_plan', id);
  };

  const handlePlanRevision = updatedPlan => {
    sendSocketEvent('revise_plan', JSON.stringify(updatedPlan));
  };

  const handlePlanDeletion = planId => () => {
    sendSocketEvent('burn_plan', planId);
  };

  const toggleShow = show => () => {
    storysShow = show;
  };

  let pointedPlans = $derived(plans.filter(p => p.points !== ''));
  let totalPoints = $derived(pointedPlans.reduce((previousValue, currentValue) => {
    var currentPoints =
      currentValue.points === '1/2' ? 0.5 : parseInt(currentValue.points);
    return isNaN(currentPoints) ? previousValue : previousValue + currentPoints;
  }, 0));
  let unpointedPlans = $derived(plans.filter(p => p.points === ''));

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
    }
  }
</script>

<div class="shadow-lg mb-4">
  <div class="flex items-center bg-gray-200 dark:bg-gray-700 p-4 rounded-t-lg">
    <div class="w-1/3">
      <h3
        class="text-3xl leading-tight font-semibold font-rajdhani uppercase dark:text-white"
      >
        {$LL.plans()}
      </h3>
    </div>
    <div class="w-2/3 text-right">
      {#if isFacilitator}
        <HollowButton onClick={toggleImport} color="blue">
          {$LL.importPlans()}
        </HollowButton>
        <HollowButton onClick={toggleAddPlan(null)} testid="plans-add">
          {$LL.planAdd()}
        </HollowButton>
      {/if}
    </div>
  </div>

  <ul
    class="flex border-b border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800"
  >
    <li class="-mb-px me-1">
      <button
        class="{storysShow === 'unpointed'
          ? 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'
          : 'hover:text-blue-600 dark:hover:text-sky-300 text-blue-400 dark:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
        onclick={toggleShow('unpointed')}
        data-testid="plans-unpointed"
      >
        {$LL.unpointed({ count: unpointedPlans.length })}
      </button>
    </li>
    <li class="me-1">
      <button
        class="{storysShow === 'pointed'
          ? 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'
          : 'hover:text-blue-600 dark:hover:text-sky-300 text-blue-400 dark:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
        onclick={toggleShow('pointed')}
        data-testid="plans-pointed"
      >
        {$LL.pointed({ count: pointedPlans.length })}
      </button>
    </li>
    <li class="me-1">
      <button
        class="{storysShow === 'all'
          ? 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'
          : 'hover:text-blue-600 dark:hover:text-sky-300 text-blue-400 dark:text-sky-600'}
                  bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
        onclick={toggleShow('all')}
        data-testid="plans-all"
      >
        {$LL.allStoryWithCount({ count: plans.length })}
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
    }}"
    onconsider={handleDndConsider}
    onfinalize={handleDndFinalize}
  >
    {#each plans as plan (plan.id)}
      <div
        class="relative flex items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800{isFacilitator &&
        storysShow === 'all'
          ? ' cursor-pointer'
          : ''}{(plan.points === '' && storysShow === 'pointed') ||
        (plan.points !== '' && storysShow === 'unpointed')
          ? ' hidden'
          : ''}"
        data-testid="plan"
        data-storyid="{plan.id}"
        draggable="true"
      >
        <div class="flex-grow font-bold align-middle dark:text-white me-1">
          <span><Grip class="inline-block"/></span>
          &nbsp;
          {#if plan.link !== ''}
            <a
              href="{plan.link}"
              target="_blank"
              class="text-blue-800 dark:text-sky-400"
            >
              <ExternalLink class="inline-block" />
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
          {#if priorities[plan.priority]}
            {@const SvelteComponent = priorities[plan.priority].icon}
            <SvelteComponent
              class="inline-block w-6 h-6"
            />
          {/if}
          <span data-testid="plan-name">{plan.name}</span>
        </div>
        <div class="lg:flex-none text-right">
          {#if plan.points !== ''}
            <div
              class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded me-1"
              data-testid="plan-points"
            >
              {plan.points}
            </div>
          {/if}
          <HollowButton
            color="blue"
            onClick={togglePlanView(plan.id)}
            testid="plan-view"
          >
            {$LL.view()}
          </HollowButton>
          {#if isFacilitator}
            {#if !plan.active}
              <HollowButton
                color="red"
                onClick={handlePlanDeletion(plan.id)}
                testid="plan-delete"
              >
                {$LL.delete()}
              </HollowButton>
            {/if}
            <HollowButton
              color="purple"
              onClick={toggleAddPlan(plan.id)}
              testid="plan-edit"
            >
              {$LL.edit()}
            </HollowButton>
            {#if !plan.active}
              <HollowButton
                onClick={activatePlan(plan.id)}
                testid="plan-activate"
              >
                {$LL.activate()}
              </HollowButton>
            {/if}
          {/if}
        </div>
      </div>
      {#if plan[SHADOW_ITEM_MARKER_PROPERTY_NAME]}
        {@const SvelteComponent_1 = priorities[plan.priority].icon}
        <div
          class="opacity-50 absolute top-0 left-0 right-0 bottom-0 visible opacity-50 cursor-pointer flex items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800"
          data-testid="plan"
          data-storyid="{plan.id}"
        >
          <div class="flex-grow font-bold align-middle dark:text-white">
            {#if plan.link !== ''}
              <a
                href="{plan.link}"
                target="_blank"
                class="text-blue-800 dark:text-sky-400"
              >
                <ExternalLink class="inline-block" />
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
            <SvelteComponent_1
              class="inline-block w-6 h-6"
            />
            <span data-testid="plan-name">{plan.name}</span>
          </div>
          <div class="lg:flex-none text-right">
            {#if plan.points !== ''}
              <div
                class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded me-1"
                data-testid="plan-points"
              >
                {plan.points}
              </div>
            {/if}
            <HollowButton
              color="blue"
              onClick={togglePlanView(plan.id)}
              testid="plan-view"
            >
              {$LL.view()}
            </HollowButton>
            {#if isFacilitator}
              {#if !plan.active}
                <HollowButton
                  color="red"
                  onClick={handlePlanDeletion(plan.id)}
                  testid="plan-delete"
                >
                  {$LL.delete()}
                </HollowButton>
              {/if}
              <HollowButton
                color="purple"
                onClick={toggleAddPlan(plan.id)}
                testid="plan-edit"
              >
                {$LL.edit()}
              </HollowButton>
              {#if !plan.active}
                <HollowButton
                  onClick={activatePlan(plan.id)}
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
  {#if storysShow === 'pointed' || storysShow === 'all'}
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
    handlePlanAdd={handlePlanAdd}
    toggleAddPlan={toggleAddPlan(null)}
    handlePlanRevision={handlePlanRevision}
    planId={selectedPlan.id}
    planName={selectedPlan.name}
    planType={selectedPlan.type}
    referenceId={selectedPlan.referenceId}
    planLink={selectedPlan.link}
    description={selectedPlan.description}
    acceptanceCriteria={selectedPlan.acceptanceCriteria}
    priority={selectedPlan.priority}
    notifications={notifications}
  />
{/if}

{#if showViewPlan}
  <ViewPlan
    togglePlanView={togglePlanView(null)}
    planName={selectedPlan.name}
    planType={selectedPlan.type}
    referenceId={selectedPlan.referenceId}
    planLink={selectedPlan.link}
    description={selectedPlan.description}
    acceptanceCriteria={selectedPlan.acceptanceCriteria}
    priority={selectedPlan.priority}
  />
{/if}

{#if showImport}
  <ImportModal
    notifications={notifications}
    toggleImport={toggleImport}
    handlePlanAdd={handlePlanAdd}
    xfetch={xfetch}
    gameId={gameId}
  />
{/if}

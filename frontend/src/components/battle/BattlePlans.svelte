<script>
    import ExternalLinkIcon from '../icons/ExternalLinkIcon.svelte'
    import AddPlan from './AddPlan.svelte'
    import HollowButton from '../HollowButton.svelte'
    import ViewPlan from './ViewPlan.svelte'
    import JiraImport from './JiraImport.svelte'
    import { _ } from '../../i18n.js'

    export let plans = []
    export let isLeader = false
    export let sendSocketEvent = () => {}
    export let eventTag
    export let notifications

    const defaultPlan = {
        id: '',
        name: '',
        type: $_('planTypeStory'),
        referenceId: '',
        link: '',
        description: '',
        acceptanceCriteria: '',
    }

    let showAddPlan = false
    let showViewPlan = false
    let selectedPlan = { ...defaultPlan }
    let showCompleted = false

    const toggleAddPlan = planId => () => {
        if (planId) {
            selectedPlan = plans.find(p => p.id === planId)

            eventTag('plan_show_edit', 'battle', ``)
        } else {
            selectedPlan = { ...defaultPlan }

            eventTag('plan_show_add', 'battle', ``)
        }
        showAddPlan = !showAddPlan
    }

    const togglePlanView = planId => () => {
        if (planId) {
            selectedPlan = plans.find(p => p.id === planId)
            eventTag('plan_show_view', 'battle', ``)
        } else {
            selectedPlan = { ...defaultPlan }

            eventTag('plan_unshow_view', 'battle', ``)
        }
        showViewPlan = !showViewPlan
    }

    const handlePlanAdd = newPlan => {
        sendSocketEvent('add_plan', JSON.stringify(newPlan))
        eventTag('plan_add', 'battle', '')
    }

    const activatePlan = id => () => {
        sendSocketEvent('activate_plan', id)
        eventTag('plan_activate', 'battle', '')
    }

    const handlePlanRevision = updatedPlan => {
        sendSocketEvent('revise_plan', JSON.stringify(updatedPlan))
        eventTag('plan_revise', 'battle', '')
    }

    const handlePlanDeletion = planId => () => {
        sendSocketEvent('burn_plan', planId)
        eventTag('plan_burn', 'battle', '')
    }

    const toggleShowCompleted = show => () => {
        showCompleted = show
        eventTag('plans_show', 'battle', `completed: ${show}`)
    }

    $: pointedPlans = plans.filter(p => p.points !== '')
    $: totalPoints = pointedPlans.reduce((previousValue, currentValue) => {
        var currentPoints =
            currentValue.points === '1/2' ? 0.5 : parseInt(currentValue.points)
        return isNaN(currentPoints)
            ? previousValue
            : previousValue + currentPoints
    }, 0)
    $: unpointedPlans = plans.filter(p => p.points === '')

    $: plansToShow = showCompleted ? pointedPlans : unpointedPlans
</script>

<div class="shadow-lg mb-4">
    <div
        class="flex items-center bg-gray-200 dark:bg-gray-700 p-4 rounded-t-lg"
    >
        <div class="w-1/2">
            <h3
                class="text-3xl leading-tight font-semibold font-rajdhani uppercase dark:text-white"
            >
                {$_('plans')}
            </h3>
        </div>
        <div class="w-1/2 text-right">
            {#if isLeader}
                <JiraImport
                    handlePlanAdd="{handlePlanAdd}"
                    notifications="{notifications}"
                    eventTag="{eventTag}"
                    testid="plans-importjira"
                />
                <HollowButton onClick="{toggleAddPlan()}" testid="plans-add">
                    {$_('planAdd')}
                </HollowButton>
            {/if}
        </div>
    </div>

    <ul
        class="flex border-b border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 "
    >
        <li class="-mb-px {showCompleted ? '' : 'mr-1'}">
            <button
                class="{showCompleted
                    ? 'hover:text-blue-600 text-blue-400 dark:hover:text-sky-300 dark:text-sky-600'
                    : 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
                on:click="{toggleShowCompleted(false)}"
                data-testid="plans-unpointed"
            >
                {$_('unpointed', { values: { count: unpointedPlans.length } })}
            </button>
        </li>
        <li class="mr-1 {showCompleted ? 'mr-1' : ''}">
            <button
                class="{showCompleted
                    ? 'border-b border-blue-500 dark:border-sky-300 text-blue-600 dark:text-sky-300 hover:text-blue-800 dark:hover:text-sky-600'
                    : 'hover:text-blue-600 dark:hover:text-sky-300 text-blue-400 dark:text-sky-600'}
                bg-white dark:bg-gray-800 inline-block py-4 px-4 font-semibold"
                on:click="{toggleShowCompleted(true)}"
                data-testid="plans-pointed"
            >
                {$_('pointed', { values: { count: pointedPlans.length } })}
            </button>
        </li>
    </ul>

    {#each plansToShow as plan (plan.id)}
        <div
            class="flex flex-wrap items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800 "
            data-testid="plan"
        >
            <div class="w-full lg:w-2/3 mb-4 lg:mb-0">
                <div
                    class="inline-block font-bold align-middle dark:text-white"
                >
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
                    <span data-testid="plan-name">{plan.name}</span>
                </div>
                &nbsp;
                {#if plan.points !== ''}
                    <div
                        class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded ml-2"
                        data-testid="plan-points"
                    >
                        {plan.points}
                    </div>
                {/if}
            </div>
            <div class="w-full lg:w-1/3 text-right">
                <HollowButton
                    color="blue"
                    onClick="{togglePlanView(plan.id)}"
                    testid="plan-view"
                >
                    {$_('view')}
                </HollowButton>
                {#if isLeader}
                    {#if !plan.active}
                        <HollowButton
                            color="red"
                            onClick="{handlePlanDeletion(plan.id)}"
                            testid="plan-delete"
                        >
                            {$_('delete')}
                        </HollowButton>
                    {/if}
                    <HollowButton
                        color="purple"
                        onClick="{toggleAddPlan(plan.id)}"
                        testid="plan-edit"
                    >
                        {$_('edit')}
                    </HollowButton>
                    {#if !plan.active}
                        <HollowButton
                            onClick="{activatePlan(plan.id)}"
                            testid="plan-activate"
                        >
                            {$_('activate')}
                        </HollowButton>
                    {/if}
                {/if}
            </div>
        </div>
    {/each}
    {#if showCompleted && totalPoints}
        <div
            class="flex flex-wrap items-center border-b border-gray-300 dark:border-gray-700 p-4 bg-white dark:bg-gray-800 "
        >
            <div class="w-full lg:w-2/3 mb-4 lg:mb-0">
                <div
                    class="inline-block font-bold align-middle dark:text-gray-300"
                >
                    {$_('totalPoints')}:
                </div>
                &nbsp;
                <div
                    class="inline-block font-bold text-green-600 dark:text-lime-400
                        border-green-500 dark:border-lime-400 border px-2 py-1 rounded ml-2"
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
    />
{/if}

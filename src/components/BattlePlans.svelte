<script>
    import AddPlan from '../components/AddPlan.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    
    export let plans = []
    export let isLeader = false
    export let sendSocketEvent = () => {}

    let showAddPlan = false
    let revisePlanId = ''
    let revisePlanName = ''
    let showCompleted = false

    const toggleAddPlan = (planId) => () => {
        if (planId) {
            const planName = plans.find(p => p.id === planId).name
            revisePlanId = planId
            revisePlanName = planName
        } else {
            revisePlanId = ''
            revisePlanName = ''
        }
        showAddPlan = !showAddPlan
    }

    const handlePlanAdd = (planName) => {
        sendSocketEvent('add_plan', planName)
    }

    const activatePlan = id => () => {
        sendSocketEvent('activate_plan', id)
    }

    const handlePlanRevision = (updatedPlan) => {
        sendSocketEvent('revise_plan', JSON.stringify(updatedPlan))
    }

    const handlePlanDeletion = (planId) => () => {
        sendSocketEvent('burn_plan', planId)
    }

    const toggleShowCompleted = (show) => () => {
        showCompleted = show
    }

    $: pointedPlans = plans.filter(p => p.points !== '')
    $: unpointedPlans = plans.filter(p => p.points === '')

    $: plansToShow = showCompleted ? pointedPlans : unpointedPlans
</script>

<div class="bg-white shadow-md mb-4 rounded">
    <div class="flex items-center bg-grey-lighter p-4 rounded-t">
        <div class="w-1/2 lg:w-3/4">
            <h3 class="text-2xl">Plans</h3>
        </div>
        <div class="w-1/2 lg:w-1/4 text-right">
            {#if isLeader}
                <HollowButton color="blue" onClick={toggleAddPlan()}>
                    Add Plan
                </HollowButton>
            {/if}
        </div>
    </div>

    <ul class="list-reset flex border-b border-grey-light">
        <li class="-mb-px {showCompleted ? '' : 'mr-1'}">
            <button
                class="{showCompleted ? 'hover:text-blue-dark text-blue-light' : 'border-b border-blue text-blue-dark hover:text-blue-darker'} bg-white inline-block py-4 px-4 font-semibold"
                on:click={toggleShowCompleted(false)}
            >
                Unpointed ({unpointedPlans.length})
            </button>
        </li>
        <li class="mr-1 {showCompleted ? 'mr-1' : ''}">
            <button
                class="{showCompleted ? 'border-b border-blue text-blue-dark hover:text-blue-darker' : 'hover:text-blue-dark text-blue-light'} bg-white inline-block py-4 px-4 font-semibold"
                on:click={toggleShowCompleted(true)}
            >
                Pointed ({pointedPlans.length})
            </button>
        </li>
    </ul>

    {#each plansToShow as plan (plan.id)}
        <div class="flex flex-wrap items-center border-b border-grey-light p-4" data-testId="battlePlan" data-planName={plan.name}>
            <div class="w-full lg:w-3/4 mb-4 lg:mb-0">
                <div class="inline-block font-bold align-middle" data-testId="battlePlanName">{plan.name}</div>
                &nbsp;
                {#if plan.points !== ''}<div class="inline-block font-bold text-green-dark border-green border px-2 py-1 rounded ml-2" data-testId="battlePlanPoints">{plan.points}</div>{/if}
            </div>
            <div class="w-full lg:w-1/4 text-right">
            {#if isLeader}
                {#if !plan.active}
                    <HollowButton color="red" onClick={handlePlanDeletion(plan.id)}>Delete</HollowButton>
                {/if}
                <HollowButton color="purple" onClick={toggleAddPlan(plan.id)}>Edit</HollowButton>
                {#if !plan.active}
                    <HollowButton onClick={activatePlan(plan.id)}>Activate</HollowButton>
                {/if}
            {/if}
            </div>
        </div>
    {/each}
</div>

{#if showAddPlan}
    <AddPlan handlePlanAdd={handlePlanAdd} toggleAddPlan={toggleAddPlan()} handlePlanRevision={handlePlanRevision} planId={revisePlanId} planName={revisePlanName} />
{/if}
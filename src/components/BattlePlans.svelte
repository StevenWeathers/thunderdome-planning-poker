<script>
    import AddPlan from '../components/AddPlan.svelte'
    
    export let plans = []
    export let isLeader = false
    export let sendSocketEvent = () => {}

    let showAddPlan = false
    let revisePlanId = ''
    let revisePlanName = ''

    const toggleAddPlan = (planId, planName) => () => {
        if (planId) {
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
</script>

<div class="bg-white shadow-md mb-4 rounded">
    <div class="flex bg-grey-lighter p-4 rounded-t">
        <div class="w-3/4">
            <h3 class="text-2xl">Plans</h3>
        </div>
        <div class="w-1/4 text-right">
            {#if isLeader}
                <button
                    class="bg-transparent hover:bg-blue text-blue-dark font-semibold hover:text-white py-2 px-4 border border-blue hover:border-transparent rounded"
                    on:click={toggleAddPlan()}
                >
                    Add Plan
                </button>
            {/if}
        </div>
    </div>

    {#each plans as plan (plan.id)}
        <div class="flex border-b border-grey-light p-4">
            <div class="w-3/4">
                <span class="font-bold">{plan.name}</span>
                &nbsp;
                {#if plan.points !== ''}<span class="font-bold text-green-dark border-green border px-2 py-1 rounded ml-2">{plan.points}</span>{/if}
            </div>
            <div class="w-1/4 text-right">
            {#if isLeader}
                <button
                    class="bg-transparent hover:bg-purple text-purple-dark font-semibold hover:text-white py-2 px-4 border border-purple hover:border-transparent rounded"
                    on:click={toggleAddPlan(plan.id, plan.name)}
                >
                    Edit
                </button>
                {#if !plan.active}
                    <button
                        class="bg-transparent hover:bg-red text-red-dark font-semibold hover:text-white py-2 px-4 border border-red hover:border-transparent rounded"
                        on:click={handlePlanDeletion(plan.id)}
                    >
                        Delete
                    </button>
    
                    {#if plan.points === ''}
                        <button
                            class="bg-transparent hover:bg-green text-green-dark font-semibold hover:text-white py-2 px-4 border border-green hover:border-transparent rounded"
                            on:click={activatePlan(plan.id)}
                        >
                            Activate
                        </button>
                    {/if}
                {/if}
            {/if}
            </div>
        </div>
    {/each}
</div>

{#if showAddPlan}
    <AddPlan handlePlanAdd={handlePlanAdd} toggleAddPlan={toggleAddPlan()} handlePlanRevision={handlePlanRevision} planId={revisePlanId} planName={revisePlanName} />
{/if}
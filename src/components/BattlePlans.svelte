<script>
    import AddPlan from '../components/AddPlan.svelte'
    
    export let plans = []
    export let isLeader = false
    export let sendSocketEvent = () => {}

    let showAddPlan = false

    const toggleAddPlan = () => {
        showAddPlan = !showAddPlan
    }

    const handlePlanAdd = (planName) => {
        sendSocketEvent('add_plan', planName)
    }

    const activatePlan = id => () => {
        sendSocketEvent('activate_plan', id)
    }

    // const handlePlanRevision = (updatedPlan) => {
    //     sendSocketEvent('revise_plan', JSON.stringify(updatedPlan))
    // }

    const handlePlanDeletion = (planId) => () => {
        sendSocketEvent('burn_plan', planId)
    }
</script>

<nav class="panel">
  <p class="panel-heading">
    Battle Plans
  </p>

    {#each plans as plan (plan.id)}
        <div class="panel-block">
            <span class="panel-icon">
                <i class="fas fa-book" aria-hidden="true"></i>
            </span>
            {plan.name}&nbsp;{#if plan.points !== ''}<span class="has-text-success has-text-weight-bold">{plan.points}</span>{/if}
            {#if isLeader}
                {#if !plan.active}
                    <div class="has-text-right" style="width: 100%">
                        <button class="button is-danger is-outlined" on:click={handlePlanDeletion(plan.id)}>Delete</button>
                    </div>
                    {#if plan.points === ''}
                        <div class="has-text-right" style="width: 100%">
                            <button class="button is-primary is-outlined" on:click={activatePlan(plan.id)}>Activate</button>
                        </div>
                    {/if}
                {/if}
            {/if}
        </div>
    {/each}

    {#if isLeader}
        <div class="panel-block">
            <button class="button is-link is-outlined is-fullwidth" on:click={toggleAddPlan}>Add Plan</button>
        </div>
    {/if}

    {#if showAddPlan}
        <AddPlan handlePlanAdd={handlePlanAdd} toggleAddPlan={toggleAddPlan} />
    {/if}
</nav>
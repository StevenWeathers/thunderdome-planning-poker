<script>
    import AddPlan from '../components/AddPlan.svelte'
    
    export let plans = []
    export let isLeader = false
    export let handlePlanAdd = () => {}
    export let handlePlanActivation = () => {}

    let showAddPlan = false

    const toggleAddPlan = () => {
        showAddPlan = !showAddPlan
    }

    const activatePlan = id => () => {
        handlePlanActivation(id)
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
            {plan.name}
            {#if isLeader && !plan.active}
                <div class="has-text-right" style="width: 100%">
                    <button class="button is-primary is-outlined" on:click={activatePlan(plan.id)}>Activate</button>
                </div>
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
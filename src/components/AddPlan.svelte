<script>
    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    
    let planName = ''

    function handleSubmit(event) {
        event.preventDefault()

        handlePlanAdd(planName)
        toggleAddPlan()
    }
</script>

<div class="modal is-active">
  <div class="modal-background"></div>
  <div class="modal-content">
    <div class="box">
        <form on:submit={handleSubmit}>
            <div class="field">
                <label class="label">Plan Name</label>
                <div class="control">
                    <input bind:value={planName} placeholder="Enter a plan name" class="input" required />
                </div>
            </div>
            
            <div class="field">
                <div class="control">
                    <button class="button is-success" type="submit">Save</button>
                </div>
            </div>
        </form>
    </div>
  </div>
  <button class="modal-close is-large" aria-label="close" on:click={toggleAddPlan}></button>
</div>
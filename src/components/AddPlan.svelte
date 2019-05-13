<script>
    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    export let handlePlanRevision = () => {}
    
    export let planId = ''
    export let planName = ''

    function handleSubmit(event) {
        event.preventDefault()

        if (planId === '') {
            handlePlanAdd(planName)
        } else {
            handlePlanRevision({
                planId,
                planName
            })
        }
        toggleAddPlan()
    }
</script>

<div class="fixed pin flex items-center">
    <div class="fixed pin bg-black opacity-75 z-10"></div>

    <div class="relative mx-6 md:mx-auto w-full md:w-1/2 lg:w-1/4 z-20 m-8">
        <div class="shadow-lg bg-white rounded-lg p-8">
            <div class="flex justify-end mb-2">
                <button aria-label="close" on:click={toggleAddPlan}>
                    <span class="mr-2">
                        <svg class="fill-current text-grey-darkest" width="24" height="24" viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg">
                            <path d="M1490 1322q0 40-28 68l-136 136q-28 28-68 28t-68-28l-294-294-294 294q-28 28-68 28t-68-28l-136-136q-28-28-28-68t28-68l294-294-294-294q-28-28-28-68t28-68l136-136q28-28 68-28t68 28l294 294 294-294q28-28 68-28t68 28l136 136q28 28 28 68t-28 68l-294 294 294 294q28 28 28 68z"/>
                        </svg>
                    </span>
                </button>
            </div>

            <form class="pb-2 my-2" on:submit={handleSubmit}>
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="planName">
                        Plan Name
                    </label>
                    <input
                        class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker"
                        id="planName"
                        type="text"
                        bind:value={planName}
                        placeholder="Enter a plan name"
                    >
                </div>
                <div class="text-right">
                    <div>
                        <button class="bg-green hover:bg-green-dark text-white font-bold py-2 px-4 rounded" type="submit">
                            Save
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>
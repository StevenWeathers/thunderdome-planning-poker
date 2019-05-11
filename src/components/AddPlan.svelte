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

<div class="fixed pin flex items-center">
    <div class="fixed pin bg-black opacity-75 z-10"></div>

    <div class="relative mx-6 md:mx-auto w-full md:w-1/2 lg:w-1/4 z-20 m-8">
        <div class="shadow-lg bg-white rounded-lg p-8">
            <div class="flex justify-end mb-2">
                <button aria-label="close" on:click={toggleAddPlan}>
                    <span class="mr-2">Close</span>
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
                            Sign In
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>
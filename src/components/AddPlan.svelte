<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'

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

<div class="fixed pin flex items-center z-40">
    <div class="fixed pin bg-black opacity-75"></div>

    <div class="relative mx-6 md:mx-auto w-full md:w-1/2 lg:w-1/4 z-50 m-8">
        <div class="shadow-lg bg-white rounded-lg p-6">
            <div class="flex justify-end mb-2">
                <button aria-label="close" on:click={toggleAddPlan}>
                    <CloseIcon />
                </button>
            </div>

            <form on:submit={handleSubmit} name="addPlan">
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
                        name="planName"
                    >
                </div>
                <div class="text-right">
                    <div>
                        <SolidButton type="submit">Save</SolidButton>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>
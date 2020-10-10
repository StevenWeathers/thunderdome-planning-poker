<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'

    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    export let handlePlanRevision = () => {}

    export let planId = ''
    export let planName = ''
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''

    function handleSubmit(event) {
        event.preventDefault()
        const plan = {
            planName,
            referenceId,
            link: planLink,
            description,
            acceptanceCriteria
        }
        if (planId === '') {
            handlePlanAdd(plan)
        } else {
            plan.planId = planId
            handlePlanRevision(plan)
        }
        
        toggleAddPlan()
    }
</script>

<div class="fixed inset-0 flex items-center z-40">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/3 z-50
        m-8">
        <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6">
            <div class="flex justify-end mb-2">
                <button
                    aria-label="close"
                    on:click="{toggleAddPlan}"
                    class="text-gray-800">
                    <CloseIcon />
                </button>
            </div>

            <form on:submit="{handleSubmit}" name="addPlan">
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="planName">
                        Plan Name
                    </label>
                    <input
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        type="text"
                        id="planName"
                        name="planName"
                        bind:value="{planName}"
                        placeholder="Enter a plan name"
                        />
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="referenceId">
                        Reference ID
                    </label>
                    <input
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        type="text"
                        id="referenceId"
                        name="referenceId"
                        bind:value="{referenceId}"
                        placeholder="Enter a reference ID"
                        />
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="planLink">
                        Link
                    </label>
                    <input
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        type="text"
                        id="planLink"
                        name="planLink"
                        bind:value="{planLink}"
                        placeholder="Enter a link to story"
                        />
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="description">
                        Description
                    </label>
                    <textarea
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        id="description"
                        name="description"
                        bind:value="{description}"
                        placeholder="Enter a plan description"
                        />
                </div>
                <div class="mb-4">
                    <label class="block text-sm font-bold mb-2" for="acceptanceCriteria">
                        Acceptance Criteria
                    </label>
                    <textarea
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        id="acceptanceCriteria"
                        name="acceptanceCriteria"
                        bind:value="{acceptanceCriteria}"
                        placeholder="Enter a plan acceptance criteria"
                        />
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

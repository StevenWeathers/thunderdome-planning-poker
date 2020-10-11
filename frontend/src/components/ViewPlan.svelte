<script>
    import ExternalLinkIcon from './icons/ExternalLinkIcon.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'

    export let togglePlanView = () => {}

    export let planName = ''
    export let planType = ''
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''
</script>

<div class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75 p-4"></div>

    <div class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 z-50 max-h-full">
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{togglePlanView}"
                        class="text-gray-800">
                        <CloseIcon />
                    </button>
                </div>
    
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Plan Type
                    </div>
                    {planType}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Plan Name
                    </div>
                    {planName}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Reference ID
                    </div>
                    {referenceId}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Link
                    </div>
                    {#if planLink !== ''}
                    <a href="{planLink}" target="_blank" class="text-blue-800">
                        <ExternalLinkIcon /> {planLink}
                    </a>
                    {/if}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Description
                    </div>
                    {description}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        Acceptance Criteria
                    </div>
                    {acceptanceCriteria}
                </div>
            </div>
        </div>
    </div>
</div>

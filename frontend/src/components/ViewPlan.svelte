<script>
    import marked from 'marked'

    import ExternalLinkIcon from './icons/ExternalLinkIcon.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import { _ } from '../i18n'

    export let togglePlanView = () => {}

    export let planName = ''
    export let planType = ''
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75 p-4"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 z-50 max-h-full">
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
                        {$_('actions.plan.fields.type.label')}
                    </div>
                    {planType}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        {$_('actions.plan.fields.name.label')}
                    </div>
                    {planName}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        {$_('actions.plan.fields.referenceId.label')}
                    </div>
                    {referenceId}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        {$_('actions.plan.fields.link.label')}
                    </div>
                    {#if planLink !== ''}
                        <a
                            href="{planLink}"
                            target="_blank"
                            class="text-blue-800">
                            <ExternalLinkIcon />
                            {planLink}
                        </a>
                    {/if}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        {$_('actions.plan.fields.description.label')}
                    </div>
                    {@html marked(description)}
                </div>
                <div class="mb-4">
                    <div class="text-sm font-bold mb-2">
                        {$_('actions.plan.fields.acceptanceCriteria.label')}
                    </div>
                    {@html marked(acceptanceCriteria)}
                </div>
            </div>
        </div>
    </div>
</div>

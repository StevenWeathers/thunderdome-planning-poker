<script>
    import ExternalLinkIcon from './icons/ExternalLinkIcon.svelte'
    import Modal from './Modal.svelte'
    import { _ } from '../i18n'

    export let togglePlanView = () => {}

    export let planName = ''
    export let planType = ''
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''
</script>

<Modal closeModal={togglePlanView} widthClasses="md:w-2/3 lg:w-3/5">
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
        <div class="unreset">
            {@html description}
        </div>
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('actions.plan.fields.acceptanceCriteria.label')}
        </div>
        <div class="unreset">
            {@html acceptanceCriteria}
        </div>
    </div>
</Modal>

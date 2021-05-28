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

<Modal closeModal="{togglePlanView}" widthClasses="md:w-2/3 lg:w-3/5">
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planType')}
        </div>
        {planType}
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planName')}
        </div>
        {planName}
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planReferenceId')}
        </div>
        {referenceId}
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planLink')}
        </div>
        {#if planLink !== ''}
            <a href="{planLink}" target="_blank" class="text-blue-800">
                <ExternalLinkIcon />
                {planLink}
            </a>
        {/if}
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planDescription')}
        </div>
        <div class="unreset">
            {@html description}
        </div>
    </div>
    <div class="mb-4">
        <div class="text-sm font-bold mb-2">
            {$_('planAcceptanceCriteria')}
        </div>
        <div class="unreset">
            {@html acceptanceCriteria}
        </div>
    </div>
</Modal>

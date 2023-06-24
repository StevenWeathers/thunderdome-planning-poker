<script lang="ts">
    import ExternalLinkIcon from '../icons/ExternalLinkIcon.svelte'
    import Modal from '../Modal.svelte'
    import { _ } from '../../i18n.js'
    import NoSymbolIcon from '../icons/NoSymbol.svelte'
    import DoubleChevronDown from '../icons/DoubleChevronDown.svelte'
    import DoubleChevronUp from '../icons/DoubleChevronUp.svelte'
    import ChevronDown from '../icons/ChevronDown.svelte'
    import ChevronUp from '../icons/ChevronUp.svelte'
    import Bars2 from '../icons/Bars2.svelte'

    export let togglePlanView = () => {}

    export let planName = ''
    export let planType = ''
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''
    export let priority = 99

    const priorities = {
        99: {
            name: '',
            icon: false,
        },
        1: {
            name: $_('planPriorityBlocker'),
            icon: NoSymbolIcon,
        },
        2: {
            name: $_('planPriorityHighest'),
            icon: DoubleChevronUp,
        },
        3: {
            name: $_('planPriorityHigh'),
            icon: ChevronUp,
        },
        4: {
            name: $_('planPriorityMedium'),
            icon: Bars2,
        },
        5: {
            name: $_('planPriorityLow'),
            icon: ChevronDown,
        },
        6: {
            name: $_('planPriorityLowest'),
            icon: DoubleChevronDown,
        },
    }
</script>

<Modal closeModal="{togglePlanView}" widthClasses="md:w-2/3 lg:w-3/5">
    <div class="mb-4 dark:text-white">
        <div class="font-bold mb-2 dark:text-gray-400">{$_('planType')}</div>
        {planType}
    </div>
    <div class="mb-4 dark:text-white">
        <div class="font-bold mb-2 dark:text-gray-400">{$_('planName')}</div>
        {planName}
    </div>
    <div class="mb-4 dark:text-white">
        <div class="font-bold mb-2 dark:text-gray-400">
            {$_('planReferenceId')}
        </div>
        {referenceId}
    </div>
    <div class="mb-4">
        <div class="font-bold mb-2 dark:text-gray-400">{$_('planLink')}</div>
        {#if planLink !== ''}
            <a
                href="{planLink}"
                target="_blank"
                class="text-blue-800 hover:text-blue-600 dark:text-sky-400 dark:hover:text-sky-600"
            >
                <ExternalLinkIcon />
                {planLink}
            </a>
        {/if}
    </div>
    <div class="mb-4 dark:text-white">
        <div class="font-bold mb-2 dark:text-gray-400">
            {$_('planPriority')}
        </div>
        <svelte:component this="{priorities[priority].icon}" />{priorities[
            priority
        ].name}
    </div>
    <div class="mb-4">
        <div class="font-bold mb-2 dark:text-gray-400">
            {$_('planDescription')}
        </div>
        <div class="unreset dark:text-white">
            {@html description}
        </div>
    </div>
    <div class="mb-4">
        <div class="font-bold mb-2 dark:text-gray-400">
            {$_('planAcceptanceCriteria')}
        </div>
        <div class="unreset dark:text-white">
            {@html acceptanceCriteria}
        </div>
    </div>
</Modal>

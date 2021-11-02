<script>
    import SolidButton from './SolidButton.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'
    import { _ } from '../i18n'

    export let sendSocketEvent = () => {}
    export let eventTag
    export let planId = ''
    export let points = []
    export let votingLocked = true
    export let highestVote = ''

    $: planPoints = highestVote

    const endPlanVoting = () => {
        sendSocketEvent('end_voting', planId)
        eventTag('vote_end', 'battle', '')
    }

    const skipPlan = () => {
        sendSocketEvent('skip_plan', planId)
        eventTag('plan_skip', 'battle', '')
    }

    const restartVoting = () => {
        sendSocketEvent('activate_plan', planId)
        eventTag('plan_restart_vote', 'battle', '')
    }

    function handleSubmit(event) {
        event.preventDefault()

        sendSocketEvent(
            'finalize_plan',
            JSON.stringify({
                planId,
                planPoints,
            }),
        )
        eventTag('plan_finalize', 'battle', planPoints)

        planPoints = ''
    }
</script>

{#if planId != ''}
    <div class="p-4" data-testId="votingControls">
        <SolidButton
            color="blue"
            additionalClasses="mb-2 w-full"
            onClick="{skipPlan}"
        >
            {$_('planSkip')}
        </SolidButton>
        {#if !votingLocked}
            <SolidButton additionalClasses="w-full" onClick="{endPlanVoting}">
                {$_('votingFinish')}
            </SolidButton>
        {:else}
            <SolidButton
                color="blue"
                additionalClasses="mb-2 w-full"
                onClick="{restartVoting}"
            >
                {$_('votingRestart')}
            </SolidButton>
            <form on:submit="{handleSubmit}" name="savePlanPoints">
                <legend class="text-xl mb-2 font-semibold leading-tight">
                    {$_('pages.battle.finalPoints')}
                </legend>
                <div class="flex -mx-2">
                    <div class="w-1/2 px-2">
                        <div class="relative">
                            <select
                                name="planPoints"
                                bind:value="{planPoints}"
                                required
                                class="block appearance-none w-full border-2
                                border-gray-400 text-gray-700 py-3 px-4 pr-8
                                rounded leading-tight focus:outline-none
                                focus:border-purple-500"
                            >
                                <option value="" disabled>
                                    {$_('pages.battle.points')}
                                </option>
                                {#each points as point}
                                    <option value="{point}">{point}</option>
                                {/each}
                            </select>
                            <div
                                class="pointer-events-none absolute inset-y-0
                                right-0 flex items-center px-2 text-gray-700"
                            >
                                <DownCarrotIcon />
                            </div>
                        </div>
                    </div>
                    <div class="w-1/2 text-right px-2">
                        <SolidButton
                            additionalClasses="w-full h-full"
                            type="submit"
                        >
                            {$_('save')}
                        </SolidButton>
                    </div>
                </div>
            </form>
        {/if}
    </div>
{/if}

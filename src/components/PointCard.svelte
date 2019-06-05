<script>
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    
    export let point = '1'
    export let active = false
    export let isLocked = true
    export let count = 0

    $: activeColor = active ? 'border-green bg-green-lightest text-green-dark' : 'border-grey-light bg-white'
    $: lockedClass = isLocked ? 'opacity-75 cursor-not-allowed' : 'cursor-pointer'

    function voteAction() {
        if (isLocked) {
            return false
        }
        if (!active) {
            dispatch('voted', {
                point
            })
        } else {
            dispatch('voteRetraction')
        }
    }
</script>

<style>
</style>

<div
    class="w-full rounded overflow-hidden shadow-md border {activeColor} {lockedClass} relative text-3xl lg:text-5xl"
    on:click={voteAction}
    data-testId="pointCard"
    data-active={active}
    data-locked={isLocked}
    data-point={point}
>
    {#if count}<span class="text-green font-semibold inline-block absolute pin-r pin-t p-2" data-testId="pointCardCount">{count}</span>{/if}
    <div class="py-12 md:py-16 text-center">{point}</div>
</div>
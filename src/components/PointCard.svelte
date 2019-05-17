<script>
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    
    export let point = '1'
    export let active = false
    export let isLocked = true

    $: activeColor = active ? 'border-green bg-green-lightest text-green-dark' : 'border-grey-light bg-white'
    $: lockedClass = isLocked ? 'opacity-50 cursor-not-allowed' : ''

    function voteAction() {
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

<button
    class="w-full py-12 md:py-16 rounded overflow-hidden shadow-md border {activeColor} {lockedClass} text-3xl lg:text-5xl"
    on:click={voteAction}
    disabled={isLocked}
>
    {point}
</button>
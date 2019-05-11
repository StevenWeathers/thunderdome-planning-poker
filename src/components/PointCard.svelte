<script>
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    
    export let point = '1'
    export let active = false
    export let isLocked = true

    $: activeColor = active ? 'border-green bg-green-lightest text-green-dark' : 'border-grey-light bg-white'
    $: lockedClass = isLocked ? 'opacity-50 cursor-not-allowed' : ''

    function vote() {
        dispatch('voted', {
            point
        });
    }
</script>

<style>
</style>

<button
    class="w-full py-12 rounded overflow-hidden shadow-lg border {activeColor} {lockedClass} text-5xl"
    on:click={vote}
    disabled={isLocked}
>
    {point}
</button>
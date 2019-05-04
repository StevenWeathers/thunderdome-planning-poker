<script>
    import { createEventDispatcher } from 'svelte';

    const dispatch = createEventDispatcher();
    
    export let point = '1'
    export let active = false

    function vote() {
        dispatch('voted', {
            point
        });
    }
</script>

<style>
</style>

<button class="button is-outlined is-large is-size-2 has-text-weight-bold {active ? 'is-success' : ''}" on:click={vote}>
    {point}
</button>
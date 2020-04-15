<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import CreateBattle from '../components/CreateBattle.svelte'
    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import LeaderIcon from '../components/icons/LeaderIcon.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { warrior } from '../stores.js'

    export let notifications
    export let eventTag
    export let router

    let battles = []

    fetch('/api/battles', {
        method: 'GET',
        credentials: 'same-origin',
    })
        .then(function(response) {
            if (!response.ok) {
                throw Error(response.statusText)
            }
            return response
        })
        .then(function(response) {
            return response.json()
        })
        .then(function(bs) {
            battles = bs
        })
        .catch(function(error) {
            notifications.danger('Error finding your battles')
        })

    onMount(() => {
        if (!$warrior.id) {
            router.route('/enlist')
        }
    })
</script>

<PageLayout>
    <h1 class="mb-4 text-3xl font-bold">My Battles</h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pr-4">
            {#each battles as battle}
                <div class="bg-white shadow-lg rounded mb-2">
                    <div
                        class="flex flex-wrap items-center p-4 border-gray-400
                        border-b">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight">
                            {#if $warrior.id === battle.leaderId}
                                <LeaderIcon />
                                &nbsp;
                            {/if}
                            {battle.name}
                            <div class="font-semibold md:text-sm text-gray-600">
                                {battle.plans.filter(p => p.points !== '').length}
                                of {battle.plans.length} plans pointed
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton href="/battle/{battle.id}">
                                Join Battle
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 pl-4">
            <div class="p-6 bg-white shadow-lg rounded">
                <h2 class="mb-4 text-2xl font-bold leading-tight">
                    Create a Battle
                </h2>
                <CreateBattle {notifications} {router} {eventTag} />
            </div>
        </div>
    </div>
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import { warrior } from '../stores.js'

    export let router
    export let notifications

    let appStats = {
        unregisteredWarriorCount: 0,
        registeredWarriorCount: 0,
        battleCount: 0,
        planCount: 0
    }

    fetch('/api/admin/stats', {
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
        .then(function(result) {
            appStats = result
        })
        .catch(function(error) {
            notifications.danger('Error getting application stats')
        })

    onMount(() => {
        if (!$warrior.id) {
            router.route('/enlist')
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route('/')
        }
    })
</script>

<PageLayout>
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">Admin</h1>
    </div>
    <div class="flex justify-center">
        <div class="w-full">
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2
                md:pt-4 md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Unregistered Warriors</div>
                    {appStats.unregisteredWarriorCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Registered Warriors</div>
                    {appStats.registeredWarriorCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Battles</div>
                    {appStats.battleCount}
                </div>
                <div class="w-1/4">
                    <div class="mb-2 font-bold">Plans</div>
                    {appStats.planCount}
                </div>
            </div>
        </div>
    </div>
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import { warrior } from '../stores.js'

    export let xfetch
    export let router
    export let notifications

    let appStats = {
        unregisteredWarriorCount: 0,
        registeredWarriorCount: 0,
        battleCount: 0,
        planCount: 0,
    }

    let warriors = []

    xfetch('/api/admin/stats')
        .then(res => res.json())
        .then(function(result) {
            appStats = result
        })
        .catch(function(error) {
            notifications.danger('Error getting application stats')
        })

    xfetch('/api/admin/warriors')
        .then(res => res.json())
        .then(function(result) {
            warriors = result
        })
        .catch(function(error) {
            notifications.danger('Error getting warriors')
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
    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
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

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <h2 class="text-2xl md:text-3xl font-bold text-center mb-4">Registered Warriors</h2>

            <table class="table-fixed w-full">
                <thead>
                <tr>
                    <th class="w-2/6 px-4 py-2">Name</th>
                    <th class="w-2/6 px-4 py-2">Email</th>
                    <th class="w-1/6 px-4 py-2">Verified</th>
                    <th class="w-1/6 px-4 py-2"></th>
                </tr>
                </thead>
                <tbody>
                {#each warriors as warrior}
                    <tr>
                        <td class="border px-4 py-2">{warrior.name}</td>
                        <td class="border px-4 py-2">{warrior.email}</td>
                        <td class="border px-4 py-2">{warrior.verified}</td>
                        <td class="border px-4 py-2"></td>
                    </tr>
                {/each}
                </tbody>
            </table>
        </div>
    </div>
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let battleId

    let battle = {
        name: '',
        votingLocked: false,
        autoFinishVoting: false,
        activePlanId: '',
        pointValuesAllowed: [],
        pointAverageRounding: '',
        users: [],
        plans: [],
        createdDate: '',
        updatedDate: '',
    }

    function getBattle() {
        xfetch(`/api/battles/${battleId}`)
            .then(res => res.json())
            .then(function (result) {
                battle = result.data
            })
            .catch(function () {
                notifications.danger('Error getting battle')
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
        }

        getBattle()
    })
</script>

<svelte:head>
    <title>{$_('battles')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="battles">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('battle')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <h3 class="text-2xl md:text-3xl font-bold mb-4 text-center">
                {battle.name}
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">Voting Locked</th>
                        <th class="flex-1 p-2">Auto Finish Voting</th>
                        <th class="flex-1 p-2">Point Values Allowed</th>
                        <th class="flex-1 p-2">Point Average Rounding</th>
                        <th class="flex-1 p-2">Created</th>
                        <th class="flex-1 p-2">Updated</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td class="border p-2">
                            {#if battle.votingLocked}
                                <span class="text-green-600"><CheckIcon /></span
                                >
                            {/if}
                        </td>
                        <td class="border p-2">
                            {#if battle.autoFinishVoting}
                                <span class="text-green-600"><CheckIcon /></span
                                >
                            {/if}
                        </td>
                        <td class="border p-2"
                            >{battle.pointValuesAllowed.join(', ')}</td
                        >
                        <td class="border p-2">{battle.pointAverageRounding}</td
                        >
                        <td class="border p-2"
                            >{new Date(battle.createdDate).toLocaleString()}</td
                        >
                        <td class="border p-2"
                            >{new Date(battle.updatedDate).toLocaleString()}</td
                        >
                    </tr>
                </tbody>
            </table>

            <h3 class="text-2xl md:text-3xl font-bold mb-4 text-center">
                Users
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">
                            {$_('name')}
                        </th>
                        <th class="flex-1 p-2">Type</th>
                        <th class="flex-1 p-2">Active</th>
                        <th class="flex-1 p-2">Abandoned</th>
                        <th class="flex-1 p-2">Spectator</th>
                        <th class="flex-1 p-2">Leader</th>
                    </tr>
                </thead>
                <tbody>
                    {#each battle.users as user}
                        <tr>
                            <td class="border p-2">
                                <a
                                    href="{appRoutes.admin}/users/{user.id}"
                                    class="no-underline text-blue-500 hover:text-blue-800"
                                    >{user.name}</a
                                >
                            </td>
                            <td class="border p-2">{user.rank}</td>
                            <td class="border p-2">
                                {#if user.active}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if user.abandoned}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if user.spectator}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if battle.leaders.includes(user.id)}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            <h3 class="text-2xl md:text-3xl font-bold mb-4 text-center">
                Plans
            </h3>
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">Name</th>
                        <th class="flex-1 p-2">Type</th>
                        <th class="flex-1 p-2">Reference Id</th>
                        <th class="flex-1 p-2">Vote Count</th>
                        <th class="flex-1 p-2">Points</th>
                        <th class="flex-1 p-2">Active</th>
                        <th class="flex-1 p-2">Skipped</th>
                    </tr>
                </thead>
                <tbody>
                    {#each battle.plans as plan}
                        <tr>
                            <td class="border p-2">{plan.name}</td>
                            <td class="border p-2">{plan.type}</td>
                            <td class="border p-2">{plan.referenceId}</td>
                            <td class="border p-2">{plan.votes.length}</td>
                            <td class="border p-2">{plan.points}</td>
                            <td class="border p-2">
                                {#if plan.active}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                            <td class="border p-2">
                                {#if plan.skipped}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</AdminPageLayout>

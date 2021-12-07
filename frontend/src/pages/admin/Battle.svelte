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
                notifications.danger($_('getBattleError'))
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
            return
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
            return
        }

        getBattle()
    })
</script>

<svelte:head>
    <title>{$_('battles')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="battles">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase">
            {$_('battle')}
        </h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 mb-4 md:mb-6 bg-white shadow-lg rounded">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {battle.name}
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('votingLocked')}</th>
                        <th class="flex-1 p-2">{$_('autoFinishVoting')}</th>
                        <th class="flex-1 p-2">{$_('pointValuesAllowed')}</th>
                        <th class="flex-1 p-2">{$_('pointAverageRounding')}</th>
                        <th class="flex-1 p-2">{$_('dateCreated')}</th>
                        <th class="flex-1 p-2">{$_('dateUpdated')}</th>
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
        </div>
        <div class="p-4 md:p-6 mb-4 md:mb-6 bg-white shadow-lg rounded">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('users')}
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">
                            {$_('name')}
                        </th>
                        <th class="flex-1 p-2">{$_('type')}</th>
                        <th class="flex-1 p-2">{$_('active')}</th>
                        <th class="flex-1 p-2">{$_('abandoned')}</th>
                        <th class="flex-1 p-2">{$_('spectator')}</th>
                        <th class="flex-1 p-2">{$_('leader')}</th>
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
        </div>
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('plans')}
            </h3>
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('name')}</th>
                        <th class="flex-1 p-2">{$_('type')}</th>
                        <th class="flex-1 p-2">{$_('planReferenceId')}</th>
                        <th class="flex-1 p-2">{$_('voteCount')}</th>
                        <th class="flex-1 p-2">{$_('points')}</th>
                        <th class="flex-1 p-2">{$_('active')}</th>
                        <th class="flex-1 p-2">{$_('skipped')}</th>
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

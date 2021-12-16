<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { AppConfig, appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

    const { AvatarService } = AppConfig

    export let xfetch
    export let router
    export let notifications
    // export let eventTag
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
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani">
            {battle.name}
        </h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6">
            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"
                        >
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead class="bg-gray-50">
                                    <tr>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('votingLocked')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('autoFinishVoting')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('pointValuesAllowed')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('pointAverageRounding')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('dateCreated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('dateUpdated')}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    <tr>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {#if battle.votingLocked}
                                                <span class="text-green-600"
                                                    ><CheckIcon /></span
                                                >
                                            {/if}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap">
                                            {#if battle.autoFinishVoting}
                                                <span class="text-green-600"
                                                    ><CheckIcon /></span
                                                >
                                            {/if}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{battle.pointValuesAllowed.join(
                                                ', ',
                                            )}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{battle.pointAverageRounding}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{new Date(
                                                battle.createdDate,
                                            ).toLocaleString()}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{new Date(
                                                battle.updatedDate,
                                            ).toLocaleString()}</td
                                        >
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="p-4 md:p-6">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('users')}
            </h3>

            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"
                        >
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead class="bg-gray-50">
                                    <th class="flex-1 p-2"></th>
                                    <tr>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('name')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('rank')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('active')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('abandoned')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('spectator')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('leader')}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    {#each battle.users as user, i}
                                        <tr class:bg-slate-100="{i % 2 !== 0}">
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                <div class="flex items-center">
                                                    <div
                                                        class="flex-shrink-0 h-10 w-10"
                                                    >
                                                        <UserAvatar
                                                            warriorId="{user.id}"
                                                            avatar="{user.avatar}"
                                                            avatarService="{AvatarService}"
                                                            gravatarHash="{user.gravatarHash}"
                                                            width="48"
                                                            class="h-10 w-10 rounded-full"
                                                        />
                                                    </div>
                                                    <div class="ml-4">
                                                        <div
                                                            class="text-sm font-medium text-gray-900"
                                                        >
                                                            <a
                                                                href="{appRoutes.admin}/users/{user.id}"
                                                                class="no-underline text-blue-500 hover:text-blue-800"
                                                                >{user.name}</a
                                                            >
                                                            {#if user.country}
                                                                &nbsp;
                                                                <CountryFlag
                                                                    country="{user.country}"
                                                                    additionalClass="inline-block"
                                                                    width="32"
                                                                    height="24"
                                                                />
                                                            {/if}
                                                        </div>
                                                    </div>
                                                </div>
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                                >{user.rank}</td
                                            >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {#if user.active}
                                                    <span class="text-green-600"
                                                        ><CheckIcon /></span
                                                    >
                                                {/if}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {#if user.abandoned}
                                                    <span class="text-green-600"
                                                        ><CheckIcon /></span
                                                    >
                                                {/if}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {#if user.spectator}
                                                    <span class="text-green-600"
                                                        ><CheckIcon /></span
                                                    >
                                                {/if}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
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
                    </div>
                </div>
            </div>
        </div>
        <div class="p-4 md:p-6">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('plans')}
            </h3>

            <div class="flex flex-col">
                <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div
                        class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                    >
                        <div
                            class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"
                        >
                            <table class="min-w-full divide-y divide-gray-200">
                                <thead class="bg-gray-50">
                                    <tr>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('name')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('type')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('planReferenceId')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('voteCount')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('points')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('active')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('skipped')}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    {#each battle.plans as plan, i}
                                        <tr class:bg-slate-100="{i % 2 !== 0}">
                                            <td class="px-6 py-4">
                                                {plan.name}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                                >{plan.type}</td
                                            >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                                >{plan.referenceId}</td
                                            >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                                >{plan.votes.length}</td
                                            >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                                >{plan.points}</td
                                            >
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {#if plan.active}
                                                    <span class="text-green-600"
                                                        ><CheckIcon /></span
                                                    >
                                                {/if}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
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
                </div>
            </div>
        </div>
    </div>
</AdminPageLayout>

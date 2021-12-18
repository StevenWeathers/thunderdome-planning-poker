<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import Checkin from '../components/checkin/Checkin.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import PencilIcon from '../components/icons/PencilIcon.svelte'
    import TrashIcon from '../components/icons/TrashIcon.svelte'
    import BlockedPing from '../components/checkin/BlockedPing.svelte'
    import Gauge from '../components/Gauge.svelte'
    import { _ } from '../i18n.js'
    import { warrior } from '../stores.js'
    import { AppConfig, appRoutes } from '../config.js'
    import { validateUserIsRegistered } from '../validationUtils.js'
    import {
        formatDayForInput,
        getTimezoneName,
        getTodaysDate,
        subtractDays,
    } from '../dateUtils.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let organizationId
    export let departmentId
    export let teamId

    const { AvatarService } = AppConfig

    let showCheckin = false
    let checkins = []
    let now = new Date()
    let maxNegativeDate
    let selectedDate
    let selectedCheckin
    let stats = {
        participants: 0,
        pPerc: 0,
        goals: 0,
        gPerc: 0,
        blocked: 0,
        bPerc: 0,
    }

    let team = {
        id: teamId,
        name: '',
    }
    let organization = {
        id: organizationId,
        name: '',
    }
    let department = {
        id: departmentId,
        name: '',
    }
    let users = []
    let userCount = 1

    let organizationRole = ''
    let departmentRole = ''
    let teamRole = ''

    const apiPrefix = '/api'
    $: orgPrefix = departmentId
        ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
        : `${apiPrefix}/organizations/${organizationId}`
    $: teamPrefix = organizationId
        ? `${orgPrefix}/teams/${teamId}`
        : `${apiPrefix}/teams/${teamId}`

    function getTeam() {
        xfetch(teamPrefix)
            .then(res => res.json())
            .then(function (result) {
                team = result.data.team
                teamRole = result.data.teamRole

                if (departmentId) {
                    department = result.data.department
                    departmentRole = result.data.departmentRole
                }
                if (organizationId) {
                    organization = result.data.organization
                    organizationRole = result.data.organizationRole
                }
            })
            .catch(function () {
                notifications.danger($_('teamGetError'))
                eventTag('team_checkin_team', 'engagement', 'failure')
            })
    }

    function getCheckins() {
        xfetch(
            `${teamPrefix}/checkins?date=${selectedDate}&tz=${getTimezoneName()}`,
        )
            .then(res => res.json())
            .then(function (result) {
                checkins = result.data
            })
            .catch(function () {
                notifications.danger('Error getting checkins')
                eventTag('team_checkin_checkins', 'engagement', 'failure')
            })
    }

    function getUsers() {
        xfetch(`${teamPrefix}/users?limit=1&offset=0`)
            .then(res => res.json())
            .then(function (result) {
                users = result.data
                userCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('teamGetUsersError'))
                eventTag('team_checkin_users', 'engagement', 'failure')
            })
    }

    function toggleCheckin(checkin) {
        showCheckin = !showCheckin
        if (checkin) {
            selectedCheckin = checkin
        } else {
            selectedCheckin = null
        }
    }

    function handleCheckin(checkin) {
        xfetch(`${teamPrefix}/checkins`, { body: checkin })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCheckin()
                eventTag('team_checkin_create', 'engagement', 'success')
            })
            .catch(function (error) {
                if (Array.isArray(error)) {
                    error[1].json().then(function (result) {
                        if (result.error == 'REQUIRES_TEAM_USER') {
                            notifications.danger(
                                `User must be in team to checkin`,
                            )
                        } else {
                            notifications.danger(`Error checking in`)
                        }
                    })
                } else {
                    notifications.danger(`Error checking in`)
                }
                eventTag('team_checkin_create', 'engagement', 'failure')
            })
    }

    function handleCheckinEdit(checkinId, checkin) {
        xfetch(`${teamPrefix}/checkins/${checkinId}`, {
            body: checkin,
            method: 'PUT',
        })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCheckin()
                eventTag('team_checkin_edit', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger(`Error updating checkin`)
                eventTag('team_checkin_edit', 'engagement', 'failure')
            })
    }

    function handleCheckinDelete(checkinId) {
        xfetch(`${teamPrefix}/checkins/${checkinId}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                eventTag('team_checkin_delete', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger(`Error deleting checkin`)
                eventTag('team_checkin_delete', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id || !validateUserIsRegistered($warrior)) {
            router.route(appRoutes.login)
            return
        }

        selectedDate = formatDayForInput(now)
        maxNegativeDate = formatDayForInput(subtractDays(now, 60))

        getTeam()
        getCheckins()
        getUsers()
    })

    function calculateCheckinStats() {
        const ucs = []
        stats.blocked = 0
        stats.goals = 0

        checkins.map(c => {
            // @todo - remove once multiple same day checkins are prevented
            if (!ucs.includes(c.user.id)) {
                ucs.push(c.user.id)

                if (c.blockers != '') {
                    ++stats.blocked
                }
                if (c.goalsMet) {
                    ++stats.goals
                }
            }
        })
        stats.participants = ucs.length

        stats.pPerc = (100 * stats.participants) / userCount
        stats.gPerc = (100 * stats.goals) / userCount
        stats.bPerc = (100 * stats.blocked) / userCount

        stats = stats

        return stats
    }

    $: isAdmin =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole === 'ADMIN'
    $: isTeamMember =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole !== ''

    $: checkStats = checkins && userCount && calculateCheckinStats()
</script>

<svelte:head>
    <title>{$_('team')} {team.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex sm:flex-wrap">
        <div class="md:grow">
            <h1
                class="text-3xl font-semibold font-rajdhani leading-none uppercase"
            >
                Checkin: <input
                    type="date"
                    id="checkindate"
                    bind:value="{selectedDate}"
                    min="{maxNegativeDate}"
                    max="{getTodaysDate()}"
                    on:change="{getCheckins}"
                    class="bg-transparent"
                />
            </h1>
            {#if organizationId}
                <div class="text-xl font-semibold font-rajdhani">
                    <span class="uppercase">{$_('organization')}</span>
                    <ChevronRight />
                    <a
                        class="text-blue-500 hover:text-blue-800"
                        href="{appRoutes.organization}/{organization.id}"
                        >{organization.name}</a
                    >
                    {#if departmentId}
                        &nbsp;
                        <ChevronRight />
                        <span class="uppercase">{$_('department')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                            >{department.name}</a
                        >
                        <ChevronRight />
                        <span class="uppercase">{$_('team')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                        >
                            {team.name}
                        </a>
                    {:else}
                        <ChevronRight />
                        <span class="uppercase">{$_('team')}</span>
                        <ChevronRight />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/team/{team.id}"
                        >
                            {team.name}
                        </a>
                    {/if}
                </div>
            {:else}
                <div class="text-2xl font-semibold font-rajdhani">
                    <span class="uppercase">{$_('team')}</span>
                    <ChevronRight />
                    <a
                        class="text-blue-500 hover:text-blue-800"
                        href="{appRoutes.team}/{team.id}"
                    >
                        {team.name}
                    </a>
                </div>
            {/if}
        </div>
        <div class="md:pl-2 md:shrink text-right">
            <SolidButton
                additionalClasses="font-rajdhani uppercase text-2xl"
                onClick="{toggleCheckin}"
                >Check In
            </SolidButton>
        </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 my-4">
        <div class="px-2 md:px-4">
            <Gauge
                text="Participation"
                percentage="{stats.pPerc}"
                stat="{stats.pPerc}"
            />
        </div>
        <div class="px-2 md:px-4">
            <Gauge
                text="Goals Met"
                percentage="{stats.gPerc}"
                color="green"
                stat="{stats.gPerc}"
            />
        </div>
        <div class="px-2 md:px-4">
            <Gauge
                text="Blocked"
                percentage="{stats.bPerc}"
                color="red"
                stat="{stats.bPerc}"
            />
        </div>
    </div>

    <div class="flex flex-col mt-8">
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
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600 uppercase tracking-wider"
                                >
                                    Name
                                </th>
                                <th
                                    scope="col"
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600 uppercase tracking-wider"
                                >
                                    Yesterday
                                </th>
                                <th
                                    scope="col"
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600 uppercase tracking-wider"
                                >
                                    Today
                                </th>
                                <th
                                    scope="col"
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600 uppercase tracking-wider"
                                >
                                    Blockers
                                </th>
                                <th
                                    scope="col"
                                    class="px-6 py-3 text-left text-sm font-medium text-gray-600 uppercase tracking-wider"
                                >
                                    Discuss
                                </th>
                            </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                            {#each checkins as checkin}
                                <tr>
                                    <td class="px-4 py-2 whitespace-nowrap">
                                        <div class="flex items-center">
                                            <div class="">
                                                <div
                                                    class="relative cursor-pointer w-14 h-14"
                                                >
                                                    <div
                                                        class="relative w-full h-full"
                                                    >
                                                        <div
                                                            class="w-full h-full bg-gray-200 rounded-full shadow"
                                                        >
                                                            <UserAvatar
                                                                warriorId="{checkin
                                                                    .user.id}"
                                                                avatar="{checkin
                                                                    .user
                                                                    .avatar}"
                                                                gravatarHash="{checkin
                                                                    .user
                                                                    .gravatarHash}"
                                                                avatarService="{AvatarService}"
                                                                options="{{
                                                                    class: 'w-full h-full rounded-full',
                                                                }}"
                                                            />
                                                        </div>
                                                        {#if checkin.goalsMet}
                                                            <div
                                                                class="absolute bottom-0 w-1/4 h-1/4 rounded-full shadow-md"
                                                            >
                                                                <svg
                                                                    xmlns="http://www.w3.org/2000/svg"
                                                                    viewBox="0 0 1024 1024"
                                                                >
                                                                    <circle
                                                                        class="text-white fill-current"
                                                                        cx="512"
                                                                        cy="512"
                                                                        r="512"
                                                                    ></circle>
                                                                    <circle
                                                                        fill="green"
                                                                        class="fill-current text-green-500"
                                                                        cx="512"
                                                                        cy="512"
                                                                        r="384"
                                                                    ></circle>
                                                                    <path
                                                                        class="text-white fill-current"
                                                                        d="M456.4 576.1L334.8 454.4l-81.1 81.1 121.6 121.7 81.1 81.1 81.1-81.1 243.3-243.3-81.1-81.1z"
                                                                    ></path>
                                                                </svg>
                                                            </div>
                                                        {/if}
                                                        {#if checkin.blockers != ''}
                                                            <BlockedPing />
                                                        {/if}
                                                        <div
                                                            class="hidden absolute top-0 right-0 w-1/4 h-1/4 bg-white rounded-full shadow-md"
                                                        >
                                                            <!-- emoji -->
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="ml-4">
                                                <div
                                                    class="text-sm font-medium"
                                                >
                                                    {checkin.user.name}
                                                </div>
                                                {#if checkin.user.id === $warrior.id}
                                                    <div>
                                                        <button
                                                            on:click="{() => {
                                                                toggleCheckin(
                                                                    checkin,
                                                                )
                                                            }}"
                                                            class="text-blue-500"
                                                            title="Edit"
                                                        >
                                                            <span
                                                                class="sr-only"
                                                                >Edit</span
                                                            >
                                                            <PencilIcon />
                                                        </button>
                                                        <button
                                                            on:click="{() => {
                                                                handleCheckinDelete(
                                                                    checkin.id,
                                                                )
                                                            }}"
                                                            class="text-red-500"
                                                            title="Delete"
                                                        >
                                                            <span
                                                                class="sr-only"
                                                                >Delete</span
                                                            >
                                                            <TrashIcon />
                                                        </button>
                                                    </div>
                                                {/if}
                                            </div>
                                        </div>
                                    </td>
                                    <td class="px-4 py-2 whitespace-nowrap">
                                        <div class="unreset">
                                            {@html checkin.yesterday}
                                        </div>
                                    </td>
                                    <td class="px-4 py-2 whitespace-nowrap">
                                        <div class="unreset">
                                            {@html checkin.today}
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="unreset">
                                            {@html checkin.blockers}
                                        </div>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap">
                                        <div class="unreset">
                                            {@html checkin.discuss}
                                        </div>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>

    {#if showCheckin}
        {#if selectedCheckin}
            <Checkin
                teamId="{team.id}"
                userId="{$warrior.id}"
                checkinId="{selectedCheckin.id}"
                yesterday="{selectedCheckin.yesterday}"
                today="{selectedCheckin.today}"
                blockers="{selectedCheckin.blockers}"
                discuss="{selectedCheckin.discuss}"
                goalsMet="{selectedCheckin.goalsMet}"
                toggleCheckin="{toggleCheckin}"
                handleCheckin="{handleCheckin}"
                handleCheckinEdit="{handleCheckinEdit}"
            />
        {:else}
            <Checkin
                teamId="{team.id}"
                userId="{$warrior.id}"
                toggleCheckin="{toggleCheckin}"
                handleCheckin="{handleCheckin}"
                handleCheckinEdit="{handleCheckinEdit}"
            />
        {/if}
    {/if}
</PageLayout>

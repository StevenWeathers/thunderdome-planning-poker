<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import Checkin from '../components/user/Checkin.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
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

    const battlesPageLimit = 1000
    const usersPageLimit = 1000
    const { AvatarService } = AppConfig

    let showCheckin = false
    let checkins = []
    let now = new Date()
    let maxNegativeDate
    let selectedDate

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
    let usersPage = 1
    let battlesPage = 1

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
            })
    }

    const score = 32
    // 100%=180° so: ° = % * 1.8
    // 45 is to add the needed rotation to have the green borders at the bottom
    $: barPercent = 45 + score * 1.8

    function toggleCheckin() {
        showCheckin = !showCheckin
    }

    function handleCheckin(checkin) {
        xfetch(`${teamPrefix}/checkins`, { body: checkin })
            .then(res => res.json())
            .then(function () {
                getCheckins()
                toggleCheckin()
            })
            .catch(function () {
                notifications.danger(`Error checking in`)
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
    })

    $: isAdmin =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole === 'ADMIN'
    $: isTeamMember =
        organizationRole === 'ADMIN' ||
        departmentRole === 'ADMIN' ||
        teamRole !== ''
</script>

<svelte:head>
    <title>{$_('team')} {team.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex sm:flex-wrap">
        <div class="md:grow">
            <h1
                class="text-4xl font-semibold font-rajdhani leading-none uppercase"
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
                <div class="text-2xl font-semibold font-rajdhani uppercase">
                    {$_('organization')}
                    <ChevronRight class="inline-block" />
                    <a
                        class="text-blue-500 hover:text-blue-800"
                        href="{appRoutes.organization}/{organization.id}"
                    >
                        {organization.name}
                    </a>
                    {#if departmentId}
                        &nbsp;
                        <ChevronRight class="inline-block" />
                        {$_('department')}
                        <ChevronRight class="inline-block" />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                        >
                            {department.name}
                        </a>
                        <ChevronRight class="inline-block" />
                        {$_('team')}
                        <ChevronRight class="inline-block" />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/department/{department.id}"
                        >
                            {team.name}
                        </a>
                    {:else}
                        <ChevronRight class="inline-block" />
                        {$_('team')}
                        <ChevronRight class="inline-block" />
                        <a
                            class="text-blue-500 hover:text-blue-800"
                            href="{appRoutes.organization}/{organization.id}/team/{team.id}"
                        >
                            {team.name}
                        </a>
                    {/if}
                </div>
            {:else}
                <div class="text-2xl font-semibold font-rajdhani uppercase">
                    {$_('team')}
                    <ChevronRight class="inline-block" />
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
                                            <div
                                                class="flex-shrink-0 h-10 w-10"
                                            >
                                                <UserAvatar
                                                    warriorId="{checkin.user
                                                        .id}"
                                                    avatar="{checkin.user
                                                        .avatar}"
                                                    avatarService="{AvatarService}"
                                                    options="{{
                                                        class: 'h-10 w-10 rounded-full',
                                                    }}"
                                                />
                                            </div>
                                            <div class="ml-4">
                                                <div
                                                    class="text-sm font-medium"
                                                >
                                                    {checkin.user.name}
                                                </div>
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
        <Checkin
            teamId="{team.id}"
            userId="{warrior.id}"
            toggleCheckin="{toggleCheckin}"
            handleCheckin="{handleCheckin}"
        />
    {/if}
</PageLayout>

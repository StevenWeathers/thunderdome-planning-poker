<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import Checkin from '../components/user/Checkin.svelte'
    import UserAvatar from '../components/user/UserAvatar.svelte'
    import { _ } from '../i18n'
    import { warrior } from '../stores.js'
    import { AppConfig, appRoutes } from '../config'
    import { validateUserIsRegistered } from '../validationUtils'
    import { onMount } from 'svelte'
    import { getTimezoneName, getTodaysDate } from '../dateUtils'

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
            `${teamPrefix}/checkins?date=${getTodaysDate()}&tz=${getTimezoneName()}`,
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

<style>
    .gauge {
        position: relative;
        text-align: center;
        width: 270px;
    }

    .bar-overflow {
        position: relative;
        width: 270px;
        height: 135px;
        margin-bottom: -40px;
        overflow: hidden;
    }

    .bar {
        position: absolute;
        top: 0;
        left: 0;
        width: 270px;
        height: 270px;
        border-radius: 50%;
        box-sizing: border-box;
        border: 20px solid #ccc;
    }

    .bar-blue {
        border-bottom-color: #0bf;
        border-right-color: #0bf;
    }

    .bar-green {
        border-bottom-color: greenyellow;
        border-right-color: greenyellow;
    }

    .bar-red {
        border-bottom-color: orangered;
        border-right-color: orangered;
    }
</style>

<svelte:head>
    <title>{$_('team')} {team.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex">
        <div class="flex-1">
            <h1 class="text-4xl font-semibold font-rajdhani leading-none">
                {new Date().toLocaleString([], {
                    day: 'numeric',
                    month: 'numeric',
                    year: 'numeric',
                })}
            </h1>
            <h2 class="text-2xl font-rajdhani uppercase leading-none">
                {$_('team')}: {team.name}
            </h2>
        </div>
        <div class="flex-1 text-right">
            <SolidButton
                additionalClasses="font-rajdhani uppercase text-2xl"
                onClick="{toggleCheckin}"
                >Check In
            </SolidButton>
        </div>
    </div>

    <div class="w-full mt-4">
        <div
            class="p-4 md:p-6 bg-white shadow-lg rounded grid grid-cols-3 gap-4 justify-items-center font-rajdhani text-3xl"
        >
            <div class="text-center">
                <h3 class="uppercase">Checked In</h3>
                <div class="gauge">
                    <div class="bar-overflow">
                        <div
                            class="bar bar-blue"
                            style="transform: rotate({barPercent}deg)"
                        ></div>
                    </div>
                    <span>7</span>
                </div>
            </div>
            <div class="text-center">
                <h3 class="uppercase">Met Goals</h3>
                <div class="gauge">
                    <div class="bar-overflow">
                        <div
                            class="bar bar-green"
                            style="transform: rotate({barPercent}deg)"
                        ></div>
                    </div>
                    <span>5</span>
                </div>
            </div>
            <div class="text-center">
                <h3 class="uppercase">Blocked</h3>
                <div class="gauge">
                    <div class="bar-overflow">
                        <div
                            class="bar bar-red"
                            style="transform: rotate({barPercent}deg)"
                        ></div>
                    </div>
                    <span>2</span>
                </div>
            </div>
        </div>
    </div>

    <div class="w-full mt-8">
        <div class="shadow border-b border-gray-200 sm:rounded-lg">
            <table class="min-w-full divide-y divide-gray-200">
                <thead
                    class="bg-gray-200 text-sm font-medium uppercase tracking-wider"
                >
                    <tr>
                        <th class="px-6 py-3">Name</th>
                        <th class="px-6 py-3">Yesterday</th>
                        <th class="px-6 py-3">Today</th>
                        <th class="px-6 py-3">Blockers</th>
                        <th class="px-6 py-3">Discuss</th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                    {#each checkins as checkin}
                        <tr>
                            <td class="px-4 py-2">
                                <div class="flex items-center">
                                    <div class="flex-shrink-0 h-10 w-10">
                                        <UserAvatar
                                            warriorId="{checkin.userId}"
                                            avatarService="{AvatarService}"
                                            options="{{
                                                class: 'h-10 w-10 rounded-full',
                                            }}"
                                        />
                                    </div>
                                    <div class="ml-4">
                                        <div class="text-sm font-medium">
                                            {checkin.userName}
                                        </div>
                                    </div>
                                </div>
                            </td>
                            <td class="px-4 py-2">
                                <div class="unreset">
                                    {@html checkin.yesterday}
                                </div>
                            </td>
                            <td class="px-4 py-2">
                                <div class="unreset">
                                    {@html checkin.today}
                                </div>
                            </td>
                            <td class="px-6 py-4">
                                <div class="unreset">
                                    {@html checkin.blockers}
                                </div>
                            </td>
                            <td class="px-6 py-4">
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

    {#if showCheckin}
        <Checkin
            teamId="{team.id}"
            userId="{warrior.id}"
            toggleCheckin="{toggleCheckin}"
            handleCheckin="{handleCheckin}"
        />
    {/if}
</PageLayout>

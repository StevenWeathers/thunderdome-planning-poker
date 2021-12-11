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
                                            warriorId="{checkin.user.id}"
                                            avatar="{checkin.user.avatar}"
                                            avatarService="{AvatarService}"
                                            options="{{
                                                class: 'h-10 w-10 rounded-full',
                                            }}"
                                        />
                                    </div>
                                    <div class="ml-4">
                                        <div class="text-sm font-medium">
                                            {checkin.user.name}
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

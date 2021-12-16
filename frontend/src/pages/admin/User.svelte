<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
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
    export let userId

    let user = {
        id: '',
        name: '',
        email: '',
        rank: '',
        avatar: '',
        verified: false,
        notificationsEnabled: true,
        country: '',
        locale: '',
        company: '',
        jobTitle: '',
        createdDate: '',
        updatedDate: '',
        lastActive: '',
    }

    function getUser() {
        xfetch(`/api/users/${userId}`)
            .then(res => res.json())
            .then(function (result) {
                user = result.data
            })
            .catch(function () {
                notifications.danger($_('getUserError'))
            })
    }

    const battlesPageLimit = 100
    let battleCount = 0
    let battles = []
    let battlesPage = 1
    let activeBattles = false

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit
        xfetch(
            `/api/users/${userId}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}&active=${activeBattles}`,
        )
            .then(res => res.json())
            .then(function (result) {
                battles = result.data
                battleCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getBattlesError'))
            })
    }

    const changeBattlesPage = evt => {
        battlesPage = evt.detail
        getBattles()
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

        getUser()
        getBattles()
    })
</script>

<svelte:head>
    <title>{$_('users')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="users">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani">
            {user.name}
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
                                        <th scope="col" class="px-6 py-3"> </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('email')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('type')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('createdDate')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('updatedDate')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('lastActive')}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    <tr>
                                        <td class="px-6 py-4 whitespace-nowrap">
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
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{user.email}
                                            {#if user.verified}
                                                <span
                                                    class="text-green-600"
                                                    title="{$_(
                                                        'pages.admin.registeredWarriors.verified',
                                                    )}"
                                                >
                                                    <CheckIcon />
                                                </span>
                                            {/if}
                                        </td>
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{user.rank}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{new Date(
                                                user.createdDate,
                                            ).toLocaleString()}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{new Date(
                                                user.updatedDate,
                                            ).toLocaleString()}</td
                                        >
                                        <td class="px-6 py-4 whitespace-nowrap"
                                            >{new Date(
                                                user.lastActive,
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
            <h4
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('battles')}
            </h4>

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
                                            {$_('dateCreated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="px-6 py-3 text-left text-sm font-medium text-gray-500 uppercase tracking-wider"
                                        >
                                            {$_('dateUpdated')}
                                        </th>
                                        <th
                                            scope="col"
                                            class="relative px-6 py-3"
                                        >
                                            <span class="sr-only">Actions</span>
                                        </th>
                                    </tr>
                                </thead>
                                <tbody
                                    class="bg-white divide-y divide-gray-200"
                                >
                                    {#each battles as battle, i}
                                        <tr class:bg-slate-100="{i % 2 !== 0}">
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                <a
                                                    href="{appRoutes.admin}/battles/{battle.id}"
                                                    class="no-underline text-blue-500 hover:text-blue-800"
                                                    >{battle.name}</a
                                                >
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    battle.createdDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap"
                                            >
                                                {new Date(
                                                    battle.updatedDate,
                                                ).toLocaleString()}
                                            </td>
                                            <td
                                                class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                            >
                                                <HollowButton
                                                    href="{appRoutes.battle}/{battle.id}"
                                                >
                                                    {$_('battleJoin')}
                                                </HollowButton>
                                            </td>
                                        </tr>
                                    {/each}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>

            {#if battleCount > battlesPageLimit}
                <div class="pt-6 flex justify-center">
                    <Pagination
                        bind:current="{battlesPage}"
                        num_items="{battleCount}"
                        per_page="{battlesPageLimit}"
                        on:navigate="{changeBattlesPage}"
                    />
                </div>
            {/if}
        </div>
    </div>
</AdminPageLayout>

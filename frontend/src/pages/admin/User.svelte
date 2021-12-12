<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'

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
        <h1 class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase">
            {$_('users')}
        </h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 mb-4 md:mb-6 bg-white shadow-lg rounded">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani mb-4 text-center"
            >
                {user.name}
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('email')}</th>
                        <th class="flex-1 p-2">{$_('type')}</th>
                        <th class="flex-1 p-2">{$_('createdDate')}</th>
                        <th class="flex-1 p-2">{$_('updatedDate')}</th>
                        <th class="flex-1 p-2">{$_('lastActive')}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td class="border p-2"
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
                        <td class="border p-2">{user.rank}</td>
                        <td class="border p-2"
                            >{new Date(user.createdDate).toLocaleString()}</td
                        >
                        <td class="border p-2"
                            >{new Date(user.updatedDate).toLocaleString()}</td
                        >
                        <td class="border p-2"
                            >{new Date(user.lastActive).toLocaleString()}</td
                        >
                    </tr>
                </tbody>
            </table>
        </div>
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <h4
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center"
            >
                {$_('battles')}
            </h4>
            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">{$_('name')}</th>
                        <th class="flex-1 p-2">{$_('dateCreated')}</th>
                        <th class="flex-1 p-2">{$_('dateUpdated')}</th>
                        <th class="flex-1 p-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each battles as battle}
                        <tr>
                            <td class="border p-2">
                                <a
                                    href="{appRoutes.admin}/battles/{battle.id}"
                                    class="no-underline text-blue-500 hover:text-blue-800"
                                    >{battle.name}</a
                                >
                            </td>
                            <td class="border p-2"
                                >{new Date(
                                    battle.createdDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border p-2"
                                >{new Date(
                                    battle.updatedDate,
                                ).toLocaleString()}</td
                            >
                            <td class="border p-2 text-right">
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

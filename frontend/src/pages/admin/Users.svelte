<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import CreateWarrior from '../../components/CreateWarrior.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CountryFlag from '../../components/CountryFlag.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import DeleteWarrior from '../../components/DeleteWarrior.svelte'
    import SolidButton from '../../components/SolidButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const usersPageLimit = 100

    let totalUsers = 0
    let users = []
    let showCreateUser = false
    let usersPage = 1
    let userDeleteId = null
    let showUserDeletion = false
    let searchEmail = ''

    const toggleDeleteUser = id => () => {
        showUserDeletion = !showUserDeletion
        userDeleteId = id
    }

    function toggleCreateUser() {
        showCreateUser = !showCreateUser
    }

    function createUser(
        warriorName,
        warriorEmail,
        warriorPassword1,
        warriorPassword2,
    ) {
        const body = {
            name: warriorName,
            email: warriorEmail,
            password1: warriorPassword1,
            password2: warriorPassword2,
        }

        xfetch('/api/admin/users', { body })
            .then(function () {
                eventTag('admin_create_warrior', 'engagement', 'success')

                getUsers()
                toggleCreateUser()
            })
            .catch(function () {
                notifications.danger('Error encountered creating warrior')
                eventTag('admin_create_warrior', 'engagement', 'failure')
            })
    }

    function getUsers() {
        const offset = (usersPage - 1) * usersPageLimit
        const isSearch = searchEmail !== ''
        const apiPrefix = isSearch
            ? `/api/admin/search/users/email?search=${searchEmail}&`
            : '/api/admin/users?'

        if (isSearch && searchEmail.length < 3) {
            notifications.danger('Search value must be at least 3 characters')
            return
        }

        xfetch(`${apiPrefix}limit=${usersPageLimit}&offset=${offset}`)
            .then(res => res.json())
            .then(function (result) {
                users = result.data
                totalUsers = result.meta.count
            })
            .catch(function () {
                notifications.danger('Error getting warriors')
            })
    }

    function promoteUser(userId) {
        return function () {
            xfetch(`/api/admin/users/${userId}/promote`, { method: 'PATCH' })
                .then(function () {
                    eventTag('admin_promote_warrior', 'engagement', 'success')

                    getUsers()
                })
                .catch(function () {
                    notifications.danger('Error encountered promoting warrior')
                    eventTag('admin_promote_warrior', 'engagement', 'failure')
                })
        }
    }

    function demoteUser(userId) {
        return function () {
            xfetch(`/api/admin/users/${userId}/demote`, { method: 'PATCH' })
                .then(function () {
                    eventTag('admin_demote_warrior', 'engagement', 'success')

                    getUsers()
                })
                .catch(function () {
                    notifications.danger('Error encountered demoting warrior')
                    eventTag('admin_demote_warrior', 'engagement', 'failure')
                })
        }
    }

    function handleDeleteUser() {
        xfetch(`/api/users/${userDeleteId}`, { method: 'DELETE' })
            .then(function () {
                eventTag('admin_demote_warrior', 'engagement', 'success')

                getUsers()
                toggleDeleteUser(null)()
            })
            .catch(function () {
                notifications.danger('Error encountered demoting warrior')
                eventTag('admin_demote_warrior', 'engagement', 'failure')
            })
    }

    function onSearchSubmit(evt) {
        evt.preventDefault()

        usersPage = 1
        getUsers()
    }

    const changePage = evt => {
        usersPage = evt.detail
        getUsers()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if ($warrior.rank !== 'GENERAL') {
            router.route(appRoutes.landing)
        }

        getUsers()
    })
</script>

<svelte:head>
    <title>{$_('users')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="users">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('users')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <div class="flex w-full">
                <div class="w-2/5">
                    <h2 class="text-2xl md:text-3xl font-bold mb-4">
                        {$_('pages.admin.registeredWarriors.title')}
                    </h2>
                </div>
                <div class="w-3/5">
                    <div class="text-right flex w-full">
                        <div class="w-3/4">
                            <form
                                on:submit="{onSearchSubmit}"
                                name="searchUsers"
                            >
                                <div class="mb-4">
                                    <label class="mb-2" for="searchEmail">
                                        <input
                                            bind:value="{searchEmail}"
                                            placeholder="Email"
                                            class="bg-gray-200 border-gray-200 border-2 appearance-none
                    rounded py-2 px-3 text-gray-700 leading-tight
                    focus:outline-none focus:bg-white focus:border-purple-500"
                                            id="searchEmail"
                                            name="searchEmail"
                                        />
                                    </label>
                                    <SolidButton type="submit">
                                        Search
                                    </SolidButton>
                                </div>
                            </form>
                        </div>
                        <div class="w-1/4">
                            <HollowButton onClick="{toggleCreateUser}">
                                {$_('warriorCreate')}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            </div>

            <table class="table-fixed w-full">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">
                            {$_('pages.admin.registeredWarriors.name')}
                        </th>
                        <th class="flex-1 p-2">
                            {$_('pages.admin.registeredWarriors.email')}
                        </th>
                        <th class="flex-1 p-2">
                            {$_('pages.admin.registeredWarriors.company')}
                        </th>
                        <th class="flex-1 p-2">
                            {$_('pages.admin.registeredWarriors.rank')}
                        </th>
                        <th class="flex-1 p-2"></th>
                    </tr>
                </thead>
                <tbody>
                    {#each users as user}
                        <tr>
                            <td class="border p-2">
                                {user.name}
                                {#if user.country}
                                    &nbsp;
                                    <CountryFlag
                                        country="{user.country}"
                                        size="{16}"
                                        additionalClass="inline-block"
                                    />
                                {/if}
                            </td>
                            <td class="border p-2">
                                {user.email}
                                {#if user.verified}
                                    &nbsp;
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
                            <td class="border p-2">
                                <div>{user.company}</div>
                                {#if user.jobTitle}
                                    <div class="text-gray-700 text-sm">
                                        {$_(
                                            'pages.admin.registeredWarriors.jobTitle',
                                        )}:
                                        {user.jobTitle}
                                    </div>
                                {/if}
                            </td>
                            <td class="border p-2">{user.rank}</td>
                            <td class="border p-2">
                                {#if user.rank !== 'GENERAL'}
                                    <HollowButton
                                        onClick="{promoteUser(user.id)}"
                                        color="blue"
                                    >
                                        {$_('promote')}
                                    </HollowButton>
                                {:else}
                                    <HollowButton
                                        onClick="{demoteUser(user.id)}"
                                        color="blue"
                                    >
                                        {$_('demote')}
                                    </HollowButton>
                                {/if}
                                <HollowButton
                                    color="red"
                                    onClick="{toggleDeleteUser(user.id)}"
                                >
                                    {$_('delete')}
                                </HollowButton>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>

            {#if totalUsers > usersPageLimit}
                <div class="pt-6 flex justify-center">
                    <Pagination
                        bind:current="{usersPage}"
                        num_items="{totalUsers}"
                        per_page="{usersPageLimit}"
                        on:navigate="{changePage}"
                    />
                </div>
            {/if}
        </div>
    </div>

    {#if showCreateUser}
        <CreateWarrior
            toggleCreate="{toggleCreateUser}"
            handleCreate="{createUser}"
            notifications
        />
    {/if}

    {#if showUserDeletion}
        <DeleteWarrior
            toggleDeleteAccount="{toggleDeleteUser(null)}"
            handleDeleteAccount="{handleDeleteUser}"
        />
    {/if}
</AdminPageLayout>

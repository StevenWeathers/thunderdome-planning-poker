<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import CreateWarrior from '../../components/user/CreateUser.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import VerifiedIcon from '../../components/icons/Verified.svelte'
    import DeleteConfirmation from '../../components/DeleteConfirmation.svelte'
    import SolidButton from '../../components/SolidButton.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import ProfileForm from '../../components/user/ProfileForm.svelte'
    import Modal from '../../components/Modal.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import RowCol from '../../components/table/RowCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'

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
    let showUserEdit = false
    let selectedUserProfile = {}

    const toggleDeleteUser = id => () => {
        showUserDeletion = !showUserDeletion
        userDeleteId = id
    }

    function toggleCreateUser() {
        showCreateUser = !showCreateUser
    }

    const toggleUserEdit = profile => () => {
        showUserEdit = !showUserEdit
        selectedUserProfile = profile
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
                notifications.danger($_('createUserError'))
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
            notifications.danger($_('searchLengthError'))
            return
        }

        xfetch(`${apiPrefix}limit=${usersPageLimit}&offset=${offset}`)
            .then(res => res.json())
            .then(function (result) {
                users = result.data
                totalUsers = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getUsersError'))
            })
    }

    function handleUserEdit(p) {
        xfetch(`/api/users/${selectedUserProfile.id}`, {
            body: p,
            method: 'PUT',
        })
            .then(res => res.json())
            .then(function () {
                notifications.success($_('pages.warriorProfile.updateSuccess'))
                eventTag('update_profile', 'engagement', 'success')
                getUsers()
                toggleUserEdit({})()
            })
            .catch(function () {
                notifications.danger($_('pages.warriorProfile.errorUpdating'))
                eventTag('update_profile', 'engagement', 'failure')
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
                    notifications.danger($_('promoteUserError'))
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
                    notifications.danger($_('demoteUserError'))
                    eventTag('admin_demote_warrior', 'engagement', 'failure')
                })
        }
    }

    function disableUser(userId) {
        return function () {
            xfetch(`/api/admin/users/${userId}/disable`, { method: 'PATCH' })
                .then(function () {
                    eventTag('admin_disable_user', 'engagement', 'success')

                    getUsers()
                })
                .catch(function () {
                    notifications.danger('Error disabling user')
                    eventTag('admin_disable_user', 'engagement', 'failure')
                })
        }
    }

    function enableUser(userId) {
        return function () {
            xfetch(`/api/admin/users/${userId}/enable`, { method: 'PATCH' })
                .then(function () {
                    eventTag('admin_enable_user', 'engagement', 'success')

                    getUsers()
                })
                .catch(function () {
                    notifications.danger('Error enabling user')
                    eventTag('admin_enable_user', 'engagement', 'failure')
                })
        }
    }

    function handleDeleteUser() {
        xfetch(`/api/users/${userDeleteId}`, { method: 'DELETE' })
            .then(function () {
                eventTag('admin_delete_warrior', 'engagement', 'success')

                getUsers()
                toggleDeleteUser(null)()
            })
            .catch(function () {
                notifications.danger('deleteUserError')
                eventTag('admin_delete_warrior', 'engagement', 'failure')
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
            return
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
            return
        }

        getUsers()
    })
</script>

<svelte:head>
    <title>{$_('users')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="users">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('users')}
        </h1>
    </div>

    <div class="w-full">
        <div class="flex w-full">
            <div class="w-2/5">
                <h2 class="text-2xl md:text-3xl font-bold mb-4 dark:text-white">
                    {$_('pages.admin.registeredWarriors.title')}
                </h2>
            </div>
            <div class="w-3/5">
                <div class="text-right flex w-full">
                    <div class="w-3/4">
                        <form on:submit="{onSearchSubmit}" name="searchUsers">
                            <div class="mb-4">
                                <label class="mb-2" for="searchEmail">
                                    <input
                                        bind:value="{searchEmail}"
                                        placeholder="{$_('email')}"
                                        class="border-2 border-slate-100 dark:border-gray-800 bg-gray-300 dark:bg-gray-800 appearance-none
                    rounded py-2 px-3 text-gray-600 dark:text-gray-400 leading-tight
                    focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                                        id="searchEmail"
                                        name="searchEmail"
                                    />
                                </label>
                                <SolidButton type="submit">
                                    {$_('search')}
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

        <Table>
            <tr slot="header">
                <HeadCol>
                    {$_('name')}
                </HeadCol>
                <HeadCol>
                    {$_('email')}
                </HeadCol>
                <HeadCol>
                    {$_('pages.admin.registeredWarriors.company')}
                </HeadCol>
                <HeadCol>
                    {$_('pages.admin.registeredWarriors.rank')}
                </HeadCol>
                <HeadCol type="action">
                    <span class="sr-only">{$_('actions')}</span>
                </HeadCol>
            </tr>
            <tbody slot="body" let:class="{className}" class="{className}">
                {#each users as user, i}
                    <TableRow itemIndex="{i}">
                        <RowCol>
                            <div class="flex items-center">
                                <div class="flex-shrink-0 h-10 w-10">
                                    <UserAvatar
                                        warriorId="{user.id}"
                                        avatar="{user.avatar}"
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
                                            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
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
                        </RowCol>
                        <RowCol>
                            {user.email}
                            {#if user.verified}
                                <span
                                    class="text-green-600"
                                    title="{$_(
                                        'pages.admin.registeredWarriors.verified',
                                    )}"
                                >
                                    <VerifiedIcon />
                                </span>
                            {/if}
                        </RowCol>
                        <RowCol>
                            <div
                                class="text-sm text-gray-900 dark:text-gray-400"
                            >
                                {user.company}
                            </div>
                            <div
                                class="text-sm text-gray-500 dark:text-gray-300"
                            >
                                {user.jobTitle}
                            </div>
                        </RowCol>
                        <RowCol>
                            <span class="text-gray-500 dark:text-gray-300"
                                >{user.rank}</span
                            >
                        </RowCol>
                        <RowCol type="action">
                            {#if user.rank !== 'ADMIN'}
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
                            {#if !user.disabled}
                                <HollowButton
                                    onClick="{disableUser(user.id)}"
                                    color="orange"
                                >
                                    Disable
                                </HollowButton>
                            {:else}
                                <HollowButton
                                    onClick="{enableUser(user.id)}"
                                    color="teal"
                                >
                                    Enable
                                </HollowButton>
                            {/if}
                            <HollowButton
                                color="green"
                                onClick="{toggleUserEdit(user)}"
                            >
                                {$_('edit')}
                            </HollowButton>
                            <HollowButton
                                color="red"
                                onClick="{toggleDeleteUser(user.id)}"
                            >
                                {$_('delete')}
                            </HollowButton>
                        </RowCol>
                    </TableRow>
                {/each}
            </tbody>
        </Table>

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

    {#if showCreateUser}
        <CreateWarrior
            toggleCreate="{toggleCreateUser}"
            handleCreate="{createUser}"
            notifications
        />
    {/if}

    {#if showUserEdit}
        <Modal closeModal="{toggleUserEdit({})}">
            <ProfileForm
                profile="{selectedUserProfile}"
                handleUpdate="{handleUserEdit}"
                xfetch="{xfetch}"
                notifications="{notifications}"
                eventTag="{eventTag}"
            />
        </Modal>
    {/if}

    {#if showUserDeletion}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteUser(null)}"
            handleDelete="{handleDeleteUser}"
            confirmText="{$_('pages.warriorProfile.delete.warningStatement')}"
            confirmBtnText="{$_('pages.warriorProfile.delete.confirmButton')}"
        />
    {/if}
</AdminPageLayout>

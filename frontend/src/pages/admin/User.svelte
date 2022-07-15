<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import VerifiedIcon from '../../components/icons/Verified.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import SolidButton from '../../components/SolidButton.svelte'
    import UpdatePasswordForm from '../../components/user/UpdatePasswordForm.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import Modal from '../../components/Modal.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { AppConfig, appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import RowCol from '../../components/table/RowCol.svelte'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let userId

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig

    let showUpdatePassword = false

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

    function toggleUpdatePassword() {
        showUpdatePassword = !showUpdatePassword
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

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit
        xfetch(
            `/api/users/${userId}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
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

    const retrosPageLimit = 100
    let retroCount = 0
    let retros = []
    let retrosPage = 1

    function getRetros() {
        const offset = (retrosPage - 1) * retrosPageLimit
        xfetch(
            `/api/users/${userId}/retros?limit=${retrosPageLimit}&offset=${offset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                retros = result.data
                retroCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getRetrosError'))
            })
    }

    const changeRetrosPage = evt => {
        retrosPage = evt.detail
        getRetros()
    }

    const storyboardsPageLimit = 100
    let storyboardCount = 0
    let storyboards = []
    let storyboardsPage = 1

    function getStoryboards() {
        const offset = (storyboardsPage - 1) * storyboardsPageLimit
        xfetch(
            `/api/users/${userId}/storyboards?limit=${storyboardsPageLimit}&offset=${offset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                storyboards = result.data
                storyboardCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getStoryboardsError'))
            })
    }

    const changeStoryboardsPage = evt => {
        storyboardsPage = evt.detail
        getStoryboards()
    }

    function updatePassword(password1, password2) {
        const body = {
            password1,
            password2,
        }

        xfetch(`/api/admin/users/${userId}/password`, { body, method: 'PATCH' })
            .then(function () {
                notifications.success(
                    $_('pages.warriorProfile.passwordUpdated'),
                    1500,
                )
                toggleUpdatePassword()
                eventTag('update_password', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger(
                    $_('pages.warriorProfile.passwordUpdateError'),
                )
                eventTag('update_password', 'engagement', 'failure')
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

        getUser()
        FeaturePoker && getBattles()
        FeatureRetro && getRetros()
        FeatureStoryboard && getStoryboards()
    })
</script>

<svelte:head>
    <title>{$_('users')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="users">
    <div class="w-full">
        <div class="flex px-4 md:px-6">
            <div class="flex-1">
                <h1
                    class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
                >
                    {user.name}
                </h1>
            </div>
            <div class="flex-1 text-right">
                <SolidButton onClick="{toggleUpdatePassword}"
                    >{$_('pages.warriorProfile.updatePasswordButton')}
                </SolidButton>
            </div>
        </div>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6">
            <Table>
                <tr slot="header">
                    <HeadCol />
                    <HeadCol>
                        {$_('pages.admin.registeredWarriors.country')}
                    </HeadCol>
                    <HeadCol>
                        {$_('email')}
                    </HeadCol>
                    <HeadCol>
                        {$_('type')}
                    </HeadCol>
                    <HeadCol>
                        {$_('createdDate')}
                    </HeadCol>
                    <HeadCol>
                        {$_('updatedDate')}
                    </HeadCol>
                    <HeadCol>
                        {$_('lastActive')}
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    <TableRow itemIndex="{0}">
                        <RowCol>
                            <div class="flex items-center flex-nowrap">
                                <div class="flex-shrink-0 h-10 w-10">
                                    <UserAvatar
                                        warriorId="{user.id}"
                                        avatar="{user.avatar}"
                                        gravatarHash="{user.gravatarHash}"
                                        width="48"
                                        class="h-10 w-10 rounded-full"
                                    />
                                </div>
                            </div>
                        </RowCol>
                        <RowCol>
                            {#if user.country}
                                &nbsp;
                                <CountryFlag
                                    country="{user.country}"
                                    additionalClass="inline-block"
                                    width="32"
                                    height="24"
                                />
                            {/if}
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
                            {user.rank}
                        </RowCol>
                        <RowCol>
                            {new Date(user.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(user.updatedDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(user.lastActive).toLocaleString()}
                        </RowCol>
                    </TableRow>
                </tbody>
            </Table>
        </div>

        {#if FeaturePoker}
            <div class="p-4 md:p-6">
                <h4
                    class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
                >
                    {$_('battles')}
                </h4>

                <Table>
                    <tr slot="header">
                        <HeadCol>
                            {$_('name')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateCreated')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateUpdated')}
                        </HeadCol>
                        <HeadCol type="action">
                            <span class="sr-only">{$_('actions')}</span>
                        </HeadCol>
                    </tr>
                    <tbody
                        slot="body"
                        let:class="{className}"
                        class="{className}"
                    >
                        {#each battles as battle, i}
                            <TableRow itemIndex="{i}">
                                <RowCol>
                                    <a
                                        href="{appRoutes.admin}/battles/{battle.id}"
                                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                        >{battle.name}</a
                                    >
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        battle.createdDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        battle.updatedDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol type="action">
                                    <HollowButton
                                        href="{appRoutes.battle}/{battle.id}"
                                    >
                                        {$_('battleJoin')}
                                    </HollowButton>
                                </RowCol>
                            </TableRow>
                        {/each}
                    </tbody>
                </Table>

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
        {/if}

        {#if FeatureRetro}
            <div class="p-4 md:p-6">
                <h4
                    class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
                >
                    {$_('retros')}
                </h4>

                <Table>
                    <tr slot="header">
                        <HeadCol>
                            {$_('name')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateCreated')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateUpdated')}
                        </HeadCol>
                        <HeadCol type="action">
                            <span class="sr-only">{$_('actions')}</span>
                        </HeadCol>
                    </tr>
                    <tbody
                        slot="body"
                        let:class="{className}"
                        class="{className}"
                    >
                        {#each retros as retro, i}
                            <TableRow itemIndex="{i}">
                                <RowCol>
                                    <a
                                        href="{appRoutes.admin}/retros/{retro.id}"
                                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                        >{retro.name}</a
                                    >
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        retro.createdDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        retro.updatedDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol type="action">
                                    <HollowButton
                                        href="{appRoutes.retro}/{retro.id}"
                                    >
                                        {$_('joinRetro')}
                                    </HollowButton>
                                </RowCol>
                            </TableRow>
                        {/each}
                    </tbody>
                </Table>

                {#if retroCount > retrosPageLimit}
                    <div class="pt-6 flex justify-center">
                        <Pagination
                            bind:current="{retrosPage}"
                            num_items="{retroCount}"
                            per_page="{retrosPageLimit}"
                            on:navigate="{changeRetrosPage}"
                        />
                    </div>
                {/if}
            </div>
        {/if}

        {#if FeatureStoryboard}
            <div class="p-4 md:p-6">
                <h4
                    class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
                >
                    {$_('storyboards')}
                </h4>

                <Table>
                    <tr slot="header">
                        <HeadCol>
                            {$_('name')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateCreated')}
                        </HeadCol>
                        <HeadCol>
                            {$_('dateUpdated')}
                        </HeadCol>
                        <HeadCol type="action">
                            <span class="sr-only">{$_('actions')}</span>
                        </HeadCol>
                    </tr>
                    <tbody
                        slot="body"
                        let:class="{className}"
                        class="{className}"
                    >
                        {#each storyboards as storyboard, i}
                            <TableRow itemIndex="{i}">
                                <RowCol>
                                    <a
                                        href="{appRoutes.admin}/storyboards/{storyboard.id}"
                                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                                        >{storyboard.name}</a
                                    >
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        storyboard.createdDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol>
                                    {new Date(
                                        storyboard.updatedDate,
                                    ).toLocaleString()}
                                </RowCol>
                                <RowCol type="action">
                                    <HollowButton
                                        href="{appRoutes.storyboard}/{storyboard.id}"
                                    >
                                        {$_('joinStoryboard')}
                                    </HollowButton>
                                </RowCol>
                            </TableRow>
                        {/each}
                    </tbody>
                </Table>

                {#if storyboardCount > storyboardsPageLimit}
                    <div class="pt-6 flex justify-center">
                        <Pagination
                            bind:current="{storyboardsPage}"
                            num_items="{storyboardCount}"
                            per_page="{storyboardsPageLimit}"
                            on:navigate="{changeStoryboardsPage}"
                        />
                    </div>
                {/if}
            </div>
        {/if}
    </div>

    {#if showUpdatePassword}
        <Modal closeModal="{toggleUpdatePassword}">
            <UpdatePasswordForm
                handleUpdate="{updatePassword}"
                toggleForm="{toggleUpdatePassword}"
                notifications="{notifications}"
            />
        </Modal>
    {/if}
</AdminPageLayout>

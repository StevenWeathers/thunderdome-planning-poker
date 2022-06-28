<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
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
    export let storyboardId

    let storyboard = {
        name: '',
        users: [],
        owner_id: '',
        createdDate: '',
        updatedDate: '',
    }

    function getStoryboard() {
        xfetch(`/api/storyboards/${storyboardId}`)
            .then(res => res.json())
            .then(function (result) {
                storyboard = result.data
            })
            .catch(function () {
                notifications.danger($_('getStoryboardErrorMessage'))
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

        getStoryboard()
    })
</script>

<svelte:head>
    <title>{$_('storyboard')} {$_('pages.admin.title')} | {$_('appName')}</title
    >
</svelte:head>

<AdminPageLayout activePage="storyboards">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
        >
            {storyboard.name}
        </h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6">
            <Table>
                <tr slot="header">
                    <HeadCol>
                        {$_('dateCreated')}
                    </HeadCol>
                    <HeadCol>
                        {$_('dateUpdated')}
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    <TableRow itemIndex="{0}">
                        <RowCol>
                            {new Date(storyboard.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(storyboard.updatedDate).toLocaleString()}
                        </RowCol>
                    </TableRow>
                </tbody>
            </Table>
        </div>
        <div class="p-4 md:p-6">
            <h3
                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 text-center dark:text-white"
            >
                {$_('users')}
            </h3>

            <Table>
                <tr slot="header">
                    <HeadCol>
                        {$_('name')}
                    </HeadCol>
                    <HeadCol>
                        {$_('active')}
                    </HeadCol>
                    <HeadCol>
                        {$_('abandoned')}
                    </HeadCol>
                    <HeadCol>
                        {$_('leader')}
                    </HeadCol>
                </tr>
                <tbody slot="body" let:class="{className}" class="{className}">
                    {#each storyboard.users as user, i}
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
                                {#if user.active}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </RowCol>
                            <RowCol>
                                {#if user.abandoned}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </RowCol>
                            <RowCol>
                                <RowCol>
                                    {#if storyboard.owner_id === user.id}
                                        <span class="text-green-600"
                                            ><CheckIcon /></span
                                        >
                                    {/if}
                                </RowCol>
                            </RowCol>
                        </TableRow>
                    {/each}
                </tbody>
            </Table>
        </div>
    </div>
</AdminPageLayout>

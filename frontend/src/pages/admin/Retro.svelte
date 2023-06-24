<script lang="ts">
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import UserAvatar from '../../components/user/UserAvatar.svelte'
    import CountryFlag from '../../components/user/CountryFlag.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.ts'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import RowCol from '../../components/table/RowCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import DeleteConfirmation from '../../components/DeleteConfirmation.svelte'
    import HollowButton from '../../components/HollowButton.svelte'

    export let xfetch
    export let router
    export let notifications
    export let retroId

    let showDeleteRetro = false

    let retro = {
        name: '',
        users: [],
        owner_id: '',
        createdDate: '',
        updatedDate: '',
    }

    function getRetro() {
        xfetch(`/api/retros/${retroId}`)
            .then(res => res.json())
            .then(function (result) {
                retro = result.data
            })
            .catch(function () {
                notifications.danger($_('getRetroErrorMessage'))
            })
    }

    function deleteRetro() {
        xfetch(`/api/retros/${retroId}`, { method: 'DELETE' })
            .then(res => res.json())
            .then(function () {
                router.route(appRoutes.adminRetros)
            })
            .catch(function () {
                notifications.danger($_('deleteRetroErrorMessage'))
            })
    }

    function toggleDeleteRetro() {
        showDeleteRetro = !showDeleteRetro
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

        getRetro()
    })
</script>

<svelte:head>
    <title>{$_('retro')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="retros">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani dark:text-white"
        >
            {retro.name}
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
                            {new Date(retro.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(retro.updatedDate).toLocaleString()}
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
                    {#each retro.users as user, i}
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
                                    <div class="ms-4">
                                        <div
                                            class="text-sm font-medium text-gray-900"
                                        >
                                            <a
                                                href="{appRoutes.adminUsers}/{user.id}"
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
                                {#if retro.owner_id === user.id}
                                    <span class="text-green-600"
                                        ><CheckIcon /></span
                                    >
                                {/if}
                            </RowCol>
                        </TableRow>
                    {/each}
                </tbody>
            </Table>

            <div class="text-center mt-4">
                <HollowButton
                    color="red"
                    onClick="{toggleDeleteRetro}"
                    testid="retro-delete"
                >
                    {$_('deleteRetro')}
                </HollowButton>
            </div>

            {#if showDeleteRetro}
                <DeleteConfirmation
                    toggleDelete="{toggleDeleteRetro}"
                    handleDelete="{deleteRetro}"
                    confirmText="{$_('confirmDeleteRetro')}"
                    confirmBtnText="{$_('deleteRetro')}"
                />
            {/if}
        </div>
    </div>
</AdminPageLayout>

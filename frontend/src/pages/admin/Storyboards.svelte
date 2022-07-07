<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import Pagination from '../../components/Pagination.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n.js'
    import { appRoutes } from '../../config.js'
    import { validateUserIsAdmin } from '../../validationUtils.js'
    import Table from '../../components/table/Table.svelte'
    import HeadCol from '../../components/table/HeadCol.svelte'
    import TableRow from '../../components/table/TableRow.svelte'
    import RowCol from '../../components/table/RowCol.svelte'

    export let xfetch
    export let router
    export let notifications
    // export let eventTag

    const storyboardsPageLimit = 100
    let storyboardCount = 0
    let storyboards = []
    let storyboardsPage = 1
    let activeStoryboards = false

    function getStoryboards() {
        const storyboardsOffset = (storyboardsPage - 1) * storyboardsPageLimit
        xfetch(
            `/api/storyboards?limit=${storyboardsPageLimit}&offset=${storyboardsOffset}&active=${activeStoryboards}`,
        )
            .then(res => res.json())
            .then(function (result) {
                storyboards = result.data
                storyboardCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getStoryboardsErrorMessage'))
            })
    }

    const changePage = evt => {
        storyboardsPage = evt.detail
        getStoryboards()
    }

    const changeActiveStoryboardsToggle = () => {
        storyboardsPage = 1
        getStoryboards()
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

        getStoryboards()
    })
</script>

<style>
    .toggle-checkbox:checked {
        @apply right-0;
        @apply border-green-400;
        border-color: #68d391;
    }

    .toggle-checkbox:checked + .toggle-label {
        @apply bg-green-400;
        background-color: #68d391;
    }
</style>

<svelte:head>
    <title
        >{$_('storyboards')} {$_('pages.admin.title')} | {$_('appName')}</title
    >
</svelte:head>

<AdminPageLayout activePage="storyboards">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('storyboards')}
        </h1>
    </div>

    <div class="w-full">
        <div class="text-right mb-4">
            <div
                class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in"
            >
                <input
                    type="checkbox"
                    name="activeStoryboards"
                    id="activeStoryboards"
                    bind:checked="{activeStoryboards}"
                    on:change="{changeActiveStoryboardsToggle}"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                    for="activeStoryboards"
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
            </div>
            <label for="activeStoryboards" class="dark:text-gray-300"
                >{$_('showActiveStoryboards')}</label
            >
        </div>

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
            <tbody slot="body" let:class="{className}" class="{className}">
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
                            {new Date(storyboard.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(storyboard.updatedDate).toLocaleString()}
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
                    on:navigate="{changePage}"
                />
            </div>
        {/if}
    </div>
</AdminPageLayout>

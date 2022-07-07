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

    const retrosPageLimit = 100
    let retroCount = 0
    let retros = []
    let retrosPage = 1
    let activeRetros = false

    function getRetros() {
        const retrosOffset = (retrosPage - 1) * retrosPageLimit
        xfetch(
            `/api/retros?limit=${retrosPageLimit}&offset=${retrosOffset}&active=${activeRetros}`,
        )
            .then(res => res.json())
            .then(function (result) {
                retros = result.data
                retroCount = result.meta.count
            })
            .catch(function () {
                notifications.danger($_('getRetrosErrorMessage'))
            })
    }

    const changePage = evt => {
        retrosPage = evt.detail
        getRetros()
    }

    const changeActiveRetrosToggle = () => {
        retrosPage = 1
        getRetros()
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

        getRetros()
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
    <title>{$_('retros')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="retros">
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('retros')}
        </h1>
    </div>

    <div class="w-full">
        <div class="text-right mb-4">
            <div
                class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in"
            >
                <input
                    type="checkbox"
                    name="activeRetros"
                    id="activeRetros"
                    bind:checked="{activeRetros}"
                    on:change="{changeActiveRetrosToggle}"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label
                    for="activeRetros"
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                >
                </label>
            </div>
            <label for="activeRetros" class="dark:text-gray-300"
                >{$_('showActiveRetros')}</label
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
                            {new Date(retro.createdDate).toLocaleString()}
                        </RowCol>
                        <RowCol>
                            {new Date(retro.updatedDate).toLocaleString()}
                        </RowCol>
                        <RowCol type="action">
                            <HollowButton href="{appRoutes.retro}/{retro.id}">
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
                    on:navigate="{changePage}"
                />
            </div>
        {/if}
    </div>
</AdminPageLayout>

<script lang="ts">
    import { onMount } from 'svelte'

    import PageLayout from '../../components/PageLayout.svelte'
    import CreateRetro from '../../components/retro/CreateRetro.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior as user } from '../../stores'
    import { appRoutes } from '../../config'
    import LL from '../../i18n/i18n-svelte'

    export let xfetch
    export let notifications
    export let router
    export let eventTag

    let retros = []

    xfetch(`/api/users/${$user.id}/retros`)
        .then(res => res.json())
        .then(function (bs) {
            retros = bs.data
        })
        .catch(function () {
            notifications.danger($LL.getRetrosErrorMessage())
            eventTag('fetch_retros', 'engagement', 'failure')
        })

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.login)
        }
    })
</script>

<svelte:head>
    <title>{$LL.yourRetros()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
    <h1
        class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
        {$LL.myRetros()}
    </h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
            {#each retros as retro}
                <div
                    class="bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                >
                    <div class="flex flex-wrap items-center p-4">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                        >
                            <span data-testid="retro-name">{retro.name}</span>
                            <div
                                class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
                            >
                                {#if $user.id === retro.ownerId}
                                    {$LL.owner()}
                                {/if}
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton href="{appRoutes.retro}/{retro.id}">
                                {$LL.joinRetro()}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 md:ps-2 xl:ps-4">
            <div
                class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
            >
                <h2
                    class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
                >
                    {$LL.createARetro()}
                </h2>
                <CreateRetro
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </div>
        </div>
    </div>
</PageLayout>

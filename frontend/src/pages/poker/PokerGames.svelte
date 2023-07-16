<script lang="ts">
    import { onMount } from 'svelte'

    import PageLayout from '../../components/PageLayout.svelte'
    import CreateBattle from '../../components/poker/CreatePokerGame.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { warrior } from '../../stores'
    import LL from '../../i18n/i18n-svelte'
    import { AppConfig, appRoutes } from '../../config'
    import LeaderIcon from '../../components/icons/LeaderIcon.svelte'
    import Pagination from '../../components/Pagination.svelte'

    export let xfetch
    export let notifications
    export let eventTag
    export let router

    const battlesPageLimit = 10
    let battleCount = 0
    let battlesPage = 1
    let battles = []

    function getBattles() {
        const battlesOffset = (battlesPage - 1) * battlesPageLimit

        xfetch(
            `/api/users/${$warrior.id}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
        )
            .then(res => res.json())
            .then(function (result) {
                battles = result.data
                battleCount = result.meta.count
            })
            .catch(function () {
                notifications.danger(
                    $LL.myBattlesError({
                        friendly: AppConfig.FriendlyUIVerbs,
                    }),
                )
                eventTag('fetch_battles', 'engagement', 'failure')
            })
    }

    const changePage = evt => {
        battlesPage = evt.detail
        getBattles()
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
            return
        }
        getBattles()
    })
</script>

<svelte:head>
    <title
        >{$LL.myBattles({ friendly: AppConfig.FriendlyUIVerbs })} | {$LL.appName()}</title
    >
</svelte:head>

<PageLayout>
    <h1
        class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
        {$LL.myBattles({ friendly: AppConfig.FriendlyUIVerbs })}
    </h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
            {#each battles as battle}
                <div
                    class="bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                >
                    <div class="flex flex-wrap items-center p-4">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                        >
                            {#if battle.leaders.includes($warrior.id)}
                                <LeaderIcon />
                                &nbsp;
                            {/if}
                            <span data-testid="battle-name">{battle.name}</span>
                            <div
                                class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
                            >
                                {$LL.countPlansPointed[
                                    AppConfig.FriendlyUIVerbs
                                ]({
                                    totalPointed: battle.plans.filter(
                                        p => p.points !== '',
                                    ).length,
                                    totalPlans: battle.plans.length,
                                })}
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton href="{appRoutes.game}/{battle.id}">
                                {$LL.battleJoin({
                                    friendly: AppConfig.FriendlyUIVerbs,
                                })}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
            {#if battleCount > battlesPageLimit}
                <div class="mt-6 pt-1 flex justify-center">
                    <Pagination
                        bind:current="{battlesPage}"
                        num_items="{battleCount}"
                        per_page="{battlesPageLimit}"
                        on:navigate="{changePage}"
                    />
                </div>
            {/if}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 md:ps-2 xl:ps-4">
            <div
                class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
            >
                <h2
                    class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
                >
                    {$LL.createBattle({ friendly: AppConfig.FriendlyUIVerbs })}
                </h2>
                <CreateBattle
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </div>
        </div>
    </div>
</PageLayout>

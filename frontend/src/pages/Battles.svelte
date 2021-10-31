<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import CreateBattle from '../components/CreateBattle.svelte'
    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import LeaderIcon from '../components/icons/LeaderIcon.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let xfetch
    export let notifications
    export let eventTag
    export let router

    let battles = []

    xfetch(`/api/users/${$warrior.id}/battles`)
        .then(res => res.json())
        .then(function(result) {
            battles = result.data
        })
        .catch(function() {
            notifications.danger($_('pages.myBattles.battlesError'))
            eventTag('fetch_battles', 'engagement', 'failure')
        })

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
    })
</script>

<svelte:head>
    <title>{$_('pages.myBattles.title')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <h1 class="mb-4 text-3xl font-bold">{$_('pages.myBattles.title')}</h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pr-4">
            {#each battles as battle}
                <div class="bg-white shadow-lg rounded mb-2">
                    <div
                        class="flex flex-wrap items-center p-4 border-gray-400
                        border-b">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight">
                            {#if battle.leaders.includes($warrior.id)}
                                <LeaderIcon />
                                &nbsp;
                            {/if}
                            {battle.name}
                            <div class="font-semibold md:text-sm text-gray-600">
                                {$_('pages.myBattles.countPlansPointed', {
                                    values: {
                                        totalPointed: battle.plans.filter(
                                            p => p.points !== '',
                                        ).length,
                                        totalPlans: battle.plans.length,
                                    },
                                })}
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton href="{appRoutes.battle}/{battle.id}">
                                {$_('battleJoin')}
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 md:pl-2 xl:pl-4">
            <div class="p-6 bg-white shadow-lg rounded">
                <h2 class="mb-4 text-2xl font-bold leading-tight">
                    {$_('pages.myBattles.createBattle.title')}
                </h2>
                <CreateBattle {notifications} {router} {eventTag} {xfetch} />
            </div>
        </div>
    </div>
</PageLayout>

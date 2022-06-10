<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import CreateStoryboard from '../components/storyboard/CreateStoryboard.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { warrior as user } from '../stores.js'
    import { appRoutes } from '../config'

    export let xfetch
    export let notifications
    export let router
    export let eventTag

    let storyboards = []

    xfetch(`/api/users/${$user.id}/storyboards`)
        .then(res => res.json())
        .then(function (bs) {
            storyboards = bs.data
        })
        .catch(function (error) {
            notifications.danger('Error finding your storyboards')
            eventTag('fetch_storyboards', 'engagement', 'failure')
        })

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.login)
        }
    })
</script>

<svelte:head>
    <title>Your Storyboards | Exothermic</title>
</svelte:head>

<PageLayout>
    <h1
        class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
        My Storyboards
    </h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pr-4">
            {#each storyboards as storyboard}
                <div
                    class="bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
                >
                    <div class="flex flex-wrap items-center p-4">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
                        >
                            {storyboard.name}
                            <div
                                class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
                            >
                                {#if $user.id === storyboard.owner_id}Owner{/if}
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton
                                href="{appRoutes.storyboard}/{storyboard.id}"
                            >
                                Join Storyboard
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 md:pl-2 xl:pl-4">
            <div
                class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
            >
                <h2
                    class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
                >
                    Create a Storyboard
                </h2>
                <CreateStoryboard
                    notifications="{notifications}"
                    router="{router}"
                    eventTag="{eventTag}"
                    xfetch="{xfetch}"
                />
            </div>
        </div>
    </div>
</PageLayout>

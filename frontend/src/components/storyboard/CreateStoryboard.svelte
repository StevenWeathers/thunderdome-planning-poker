<script>
    import { onMount } from 'svelte'

    import SolidButton from '../SolidButton.svelte'
    import { warrior as user } from '../../stores.js'
    import { appRoutes } from '../../config.js'
    import { _ } from '../../i18n.js'

    export let xfetch
    export let notifications
    export let eventTag
    export let router
    export let apiPrefix = '/api'

    let storyboardName = ''
    let joinCode = ''
    let facilitatorCode = ''

    function createStoryboard(e) {
        e.preventDefault()
        const body = {
            storyboardName,
            joinCode,
            facilitatorCode,
        }

        xfetch(`${apiPrefix}/users/${$user.id}/storyboards`, { body })
            .then(res => res.json())
            .then(function ({ data }) {
                eventTag('create_storyboard', 'engagement', 'success', () => {
                    router.route(`${appRoutes.storyboard}/${data.id}`)
                })
            })
            .catch(function (error) {
                if (Array.isArray(error)) {
                    error[1].json().then(function (result) {
                        notifications.danger(
                            `${$_('Error encountered creating storyboard')} : ${
                                result.error
                            }`,
                        )
                    })
                } else {
                    notifications.danger(
                        $_('Error encountered creating storyboard'),
                    )
                }
                eventTag('create_storyboard', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.register)
        }
    })
</script>

<form on:submit="{createStoryboard}" name="createStoryboard">
    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="storyboardName"
        >
            Storyboard Name
        </label>
        <div class="control">
            <input
                name="storyboardName"
                bind:value="{storyboardName}"
                placeholder="Enter a storyboard name"
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                    rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                    focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="storyboardName"
                required
            />
        </div>
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="joinCode"
        >
            {$_('passCode')}
        </label>
        <div class="control">
            <input
                name="joinCode"
                bind:value="{joinCode}"
                placeholder="{$_('optionalPasscodePlaceholder')}"
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="joinCode"
            />
        </div>
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="facilitatorCode"
        >
            {$_('facilitatorCodeOptional')}
        </label>
        <div class="control">
            <input
                name="facilitatorCode"
                bind:value="{facilitatorCode}"
                placeholder="{$_('facilitatorCodePlaceholder')}"
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="facilitatorCode"
            />
        </div>
    </div>

    <div class="text-right">
        <SolidButton type="submit">Create Storyboard</SolidButton>
    </div>
</form>

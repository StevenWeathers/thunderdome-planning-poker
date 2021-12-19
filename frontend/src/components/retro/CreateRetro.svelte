<script>
    import { onMount } from 'svelte'

    import SolidButton from '../SolidButton.svelte'
    import { warrior as user } from '../../stores.js'
    import { appRoutes } from '../../config'

    export let xfetch
    export let notifications
    export let eventTag
    export let router
    export let apiPrefix = '/api'

    let retroName = ''
    let joinCode = ''

    function createRetro(e) {
        e.preventDefault()
        const body = {
            retroName,
            format: 'worked_improve_question',
            joinCode,
        }

        xfetch(`${apiPrefix}/users/${$user.id}/retros`, { body })
            .then(res => res.json())
            .then(function ({ data }) {
                eventTag('create_retro', 'engagement', 'success', () => {
                    router.route(`${appRoutes.retro}/${data.id}`)
                })
            })
            .catch(function (error) {
                notifications.danger('Error encountered creating retro')
                eventTag('create_retro', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.register)
        }
    })
</script>

<form on:submit="{createRetro}" name="createRetro">
    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="retroName"
        >
            Retro Name
        </label>
        <div class="control">
            <input
                name="retroName"
                bind:value="{retroName}"
                placeholder="Enter a retro name"
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="retroName"
                required
            />
        </div>
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="joinCode"
        >
            Join Code (Optional)
        </label>
        <div class="control">
            <input
                name="joinCode"
                bind:value="{joinCode}"
                placeholder="Enter a join code"
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="joinCode"
            />
        </div>
    </div>

    <div class="text-right">
        <SolidButton type="submit">Create Retro</SolidButton>
    </div>
</form>

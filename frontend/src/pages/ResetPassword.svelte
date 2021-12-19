<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { validatePasswords } from '../validationUtils.js'
    import { _ } from '../i18n.js'
    import { appRoutes } from '../config.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let resetId

    let warriorPassword1 = ''
    let warriorPassword2 = ''

    function resetWarriorPassword(e) {
        e.preventDefault()
        const body = {
            resetId,
            password1: warriorPassword1,
            password2: warriorPassword2,
        }
        const validPasswords = validatePasswords(
            warriorPassword1,
            warriorPassword2,
        )

        let noFormErrors = true

        if (!validPasswords.valid) {
            noFormErrors = false
            notifications.danger(validPasswords.error, 1500)
        }

        if (noFormErrors) {
            xfetch('/api/auth/reset-password', { body, method: 'PATCH' })
                .then(function () {
                    eventTag('reset_password', 'engagement', 'success', () => {
                        router.route(appRoutes.login, true)
                    })
                })
                .catch(function () {
                    notifications.danger(
                        $_('pages.login.passwordReset.resetError'),
                    )
                    eventTag('reset_password', 'engagement', 'failure')
                })
        }
    }

    $: resetDisabled = warriorPassword1 === '' || warriorPassword2 === ''
</script>

<svelte:head>
    <title>{$_('pages.login.passwordReset.title')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            <form
                on:submit="{resetWarriorPassword}"
                class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                name="resetWarriorPassword"
            >
                <div
                    class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                    md:leading-tight text-center dark:text-white"
                >
                    {$_('pages.login.passwordReset.title')}
                </div>

                <div class="mb-4">
                    <label
                        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="yourPassword1"
                    >
                        {$_('pages.login.passwordReset.fields.password.label')}
                    </label>
                    <input
                        bind:value="{warriorPassword1}"
                        placeholder="{$_(
                            'pages.login.passwordReset.fields.password.placeholder',
                        )}"
                        class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        id="yourPassword1"
                        name="yourPassword1"
                        type="password"
                        required
                    />
                </div>

                <div class="mb-4">
                    <label
                        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                        for="yourPassword2"
                    >
                        {$_(
                            'pages.login.passwordReset.fields.confirmPassword.label',
                        )}
                    </label>
                    <input
                        bind:value="{warriorPassword2}"
                        placeholder="{$_(
                            'pages.login.passwordReset.fields.confirmPassword.placeholder',
                        )}"
                        class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                        id="yourPassword2"
                        name="yourPassword2"
                        type="password"
                        required
                    />
                </div>

                <div class="text-right">
                    <SolidButton type="submit" disabled="{resetDisabled}">
                        {$_('pages.login.passwordReset.saveButton')}
                    </SolidButton>
                </div>
            </form>
        </div>
    </div>
</PageLayout>

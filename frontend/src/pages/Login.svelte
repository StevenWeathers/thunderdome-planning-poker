<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'
    import { _, setupI18n } from '../i18n'
    import { AppConfig, appRoutes } from '../config'

    export let router
    export let xfetch
    export let notifications
    export let eventTag
    export let battleId

    const { AllowRegistration, LdapEnabled } = AppConfig
    const authEndpoint = LdapEnabled ? '/api/auth/ldap' : '/api/auth'

    let warriorEmail = ''
    let warriorPassword = ''

    let warriorResetEmail = ''
    let forgotPassword = false

    $: targetPage = battleId
        ? `${appRoutes.battle}/${battleId}`
        : appRoutes.battles

    function authWarrior(e) {
        e.preventDefault()
        const body = {
            email: warriorEmail,
            password: warriorPassword,
        }

        xfetch(authEndpoint, { body })
            .then(res => res.json())
            .then(function (result) {
                const newWarrior = result.data
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name,
                    email: newWarrior.email,
                    rank: newWarrior.rank,
                    locale: newWarrior.locale,
                    notificationsEnabled: newWarrior.notificationsEnabled,
                })

                eventTag('login', 'engagement', 'success', () => {
                    setupI18n({
                        withLocale: newWarrior.locale,
                    })
                    router.route(targetPage, true)
                })
            })
            .catch(function () {
                notifications.danger($_('pages.login.authError'))
                eventTag('login', 'engagement', 'failure')
            })
    }

    function toggleForgotPassword() {
        forgotPassword = !forgotPassword
        eventTag(
            'forgot_password_toggle',
            'engagement',
            `forgot: ${forgotPassword}`,
        )
    }

    function sendPasswordReset(e) {
        e.preventDefault()
        const body = {
            email: warriorResetEmail,
        }

        xfetch('/api/auth/forgot-password', { body })
            .then(function () {
                notifications.success(
                    $_('pages.login.sendResetSuccess', {
                        values: { email: warriorResetEmail },
                    }),
                    2000,
                )
                forgotPassword = !forgotPassword
                eventTag('forgot_password', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('pages.login.sendResetError'))
                eventTag('forgot_password', 'engagement', 'failure')
            })
    }

    $: loginDisabled = warriorEmail === '' || warriorPassword === ''
    $: resetDisabled = warriorResetEmail === ''
</script>

<svelte:head>
    <title>{$_('pages.login.title')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            {#if !forgotPassword}
                <form
                    on:submit="{authWarrior}"
                    class="bg-white shadow-lg rounded p-6 mb-4"
                    name="authWarrior"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center"
                    >
                        {$_('pages.login.title')}
                    </div>
                    {#if battleId && AllowRegistration}
                        <div
                            class="font-semibold font-rajdhani uppercase text-lg md:text-xl mb-2 md:mb-6
                            md:leading-tight text-center"
                        >
                            {@html $_('pages.login.registerForBattle', {
                                values: {
                                    registerOpen: `<a href="${appRoutes.register}/${battleId}" class="font-bold text-blue-500 hover:text-blue-800">`,
                                    registerClose: `</a>`,
                                },
                            })}
                        </div>
                    {/if}
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourEmail"
                        >
                            {$_('pages.login.fields.email.label')}
                        </label>
                        <input
                            bind:value="{warriorEmail}"
                            placeholder="{$_(
                                'pages.login.fields.email.placeholder',
                            )}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourEmail"
                            name="yourEmail"
                            type="email"
                            required
                        />
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourPassword"
                        >
                            {$_('pages.login.fields.password.label')}
                        </label>
                        <input
                            bind:value="{warriorPassword}"
                            placeholder="{$_(
                                'pages.login.fields.password.placeholder',
                            )}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourPassword"
                            name="yourPassword"
                            type="password"
                            required
                        />
                    </div>

                    <div class="text-right">
                        {#if !LdapEnabled}
                            <button
                                type="button"
                                class="inline-block align-baseline font-bold
                                text-sm text-blue-500 hover:text-blue-800 mr-4"
                                on:click="{toggleForgotPassword}"
                            >
                                {$_('pages.login.fields.password.forgotLabel')}
                            </button>
                        {/if}
                        <SolidButton type="submit" disabled="{loginDisabled}">
                            {$_('pages.login.button')}
                        </SolidButton>
                    </div>
                </form>
            {/if}

            {#if forgotPassword}
                <form
                    on:submit="{sendPasswordReset}"
                    class="bg-white shadow-lg rounded p-6 mb-4"
                    name="resetPassword"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center"
                    >
                        {$_('forgotPassword')}
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourResetEmail"
                        >
                            {$_('email')}
                        </label>
                        <input
                            bind:value="{warriorResetEmail}"
                            placeholder="{$_('enterYourEmail')}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourResetEmail"
                            name="yourResetEmail"
                            type="email"
                            required
                        />
                    </div>

                    <div class="text-right">
                        <button
                            type="button"
                            class="inline-block align-baseline font-bold text-sm
                            text-blue-500 hover:text-blue-800 mr-4"
                            on:click="{toggleForgotPassword}"
                        >
                            {$_('cancel')}
                        </button>
                        <SolidButton type="submit" disabled="{resetDisabled}">
                            {$_('sendResetEmail')}
                        </SolidButton>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</PageLayout>

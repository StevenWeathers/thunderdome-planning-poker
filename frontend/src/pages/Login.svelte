<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'
    import { _, setupI18n } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'

    export let router
    export let xfetch
    export let notifications
    export let eventTag
    export let battleId
    export let retroId
    export let storyboardId

    const { AllowRegistration, LdapEnabled } = AppConfig
    const authEndpoint = LdapEnabled ? '/api/auth/ldap' : '/api/auth'

    let warriorEmail = ''
    let warriorPassword = ''
    let mfaToken = ''

    let warriorResetEmail = ''
    let forgotPassword = false
    let mfaRequired = false
    let mfaUser = null
    let mfaSessionId = null

    function targetPage() {
        let tp = appRoutes.battles

        if (battleId) {
            tp = `${appRoutes.battle}/${battleId}`
        }

        if (retroId) {
            tp = `${appRoutes.retro}/${retroId}`
        }

        if (storyboardId) {
            tp = `${appRoutes.storyboard}/${storyboardId}`
        }

        return tp
    }

    function authMfa(e) {
        e.preventDefault()
        const body = {
            passcode: mfaToken,
            sessionId: mfaSessionId,
        }

        xfetch('/api/auth/mfa', { body, skip401Redirect: true })
            .then(res => res.json())
            .then(function () {
                warrior.create(mfaUser)
                eventTag('login_mfa', 'engagement', 'success', () => {
                    setupI18n({
                        withLocale: mfaUser.locale,
                    })
                    router.route(targetPage(), true)
                })
            })
            .catch(function () {
                notifications.danger($_('mfaAuthError'))
                eventTag('login_mfa', 'engagement', 'failure')
            })
    }

    function authUser(e) {
        e.preventDefault()
        const body = {
            email: warriorEmail,
            password: warriorPassword,
        }

        xfetch(authEndpoint, { body, skip401Redirect: true })
            .then(res => res.json())
            .then(function (result) {
                const u = result.data.user
                const newUser = {
                    id: u.id,
                    name: u.name,
                    email: u.email,
                    rank: u.rank,
                    locale: u.locale,
                    notificationsEnabled: u.notificationsEnabled,
                }
                if (result.data.mfaRequired) {
                    mfaRequired = true
                    mfaUser = newUser
                    mfaSessionId = result.data.sessionId
                } else {
                    warrior.create(newUser)
                    eventTag('login', 'engagement', 'success', () => {
                        setupI18n({
                            withLocale: newUser.locale,
                        })
                        router.route(targetPage(), true)
                    })
                }
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
    $: mfaLoginDisabled = mfaToken = ''
</script>

<svelte:head>
    <title>{$_('pages.login.title')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            {#if !forgotPassword && !mfaRequired}
                <form
                    on:submit="{authUser}"
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                    name="authWarrior"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                        data-formtitle="login"
                    >
                        {$_('pages.login.title')}
                    </div>
                    {#if battleId && AllowRegistration}
                        <div
                            class="font-semibold font-rajdhani uppercase text-lg md:text-xl mb-2 md:mb-6
                            md:leading-tight text-center dark:text-white"
                        >
                            {@html $_('pages.login.registerForBattle', {
                                values: {
                                    registerOpen: `<a href="${appRoutes.register}/battle/${battleId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                                    registerClose: `</a>`,
                                },
                            })}
                        </div>
                    {/if}
                    {#if retroId && AllowRegistration}
                        <div
                            class="font-semibold font-rajdhani uppercase text-lg md:text-xl mb-2 md:mb-6
                            md:leading-tight text-center dark:text-white"
                        >
                            {@html $_('registerForRetro', {
                                values: {
                                    registerOpen: `<a href="${appRoutes.register}/retro/${retroId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                                    registerClose: `</a>`,
                                },
                            })}
                        </div>
                    {/if}
                    {#if storyboardId && AllowRegistration}
                        <div
                            class="font-semibold font-rajdhani uppercase text-lg md:text-xl mb-2 md:mb-6
                            md:leading-tight text-center dark:text-white"
                        >
                            {@html $_('registerForStoryboard', {
                                values: {
                                    registerOpen: `<a href="${appRoutes.register}/storyboard/${storyboardId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                                    registerClose: `</a>`,
                                },
                            })}
                        </div>
                    {/if}
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourEmail"
                        >
                            {$_('pages.login.fields.email.label')}
                        </label>
                        <input
                            bind:value="{warriorEmail}"
                            placeholder="{$_(
                                'pages.login.fields.email.placeholder',
                            )}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                            id="yourEmail"
                            name="yourEmail"
                            type="email"
                            required
                        />
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourPassword"
                        >
                            {$_('pages.login.fields.password.label')}
                        </label>
                        <input
                            bind:value="{warriorPassword}"
                            placeholder="{$_(
                                'pages.login.fields.password.placeholder',
                            )}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                    name="resetPassword"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                        data-formtitle="forgotpassword"
                    >
                        {$_('forgotPassword')}
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourResetEmail"
                        >
                            {$_('email')}
                        </label>
                        <input
                            bind:value="{warriorResetEmail}"
                            placeholder="{$_('enterYourEmail')}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
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

            {#if mfaRequired}
                <form
                    on:submit="{authMfa}"
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                    name="authMfa"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                    >
                        {$_('pages.login.title')}
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourEmail"
                        >
                            {$_('mfaTokenLabel')}
                        </label>
                        <input
                            bind:value="{mfaToken}"
                            placeholder="{$_('mfaTokenPlaceholder')}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                            id="mfaToken"
                            name="mfaToken"
                            type="text"
                            required
                        />
                    </div>

                    <div class="text-right">
                        <SolidButton
                            type="submit"
                            disabled="{mfaLoginDisabled}"
                        >
                            {$_('pages.login.button')}
                        </SolidButton>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</PageLayout>

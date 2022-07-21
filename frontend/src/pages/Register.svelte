<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import WarriorRegisterForm from '../components/user/UserRegisterForm.svelte'
    import { warrior } from '../stores.js'
    import { validateName } from '../validationUtils.js'
    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'

    export let router
    export let xfetch
    export let notifications
    export let eventTag
    export let battleId
    export let retroId
    export let storyboardId

    const guestsAllowed = AppConfig.AllowGuests
    const registrationAllowed = AppConfig.AllowRegistration

    let warriorName = $warrior.name || ''

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

    function createUserGuest(e) {
        e.preventDefault()
        const body = {
            name: warriorName,
        }
        const validName = validateName(warriorName)

        let noFormErrors = true

        if (!validName.valid) {
            noFormErrors = false
            notifications.danger(validName.error, 1500)
        }

        if (noFormErrors) {
            xfetch('/api/auth/guest', { body })
                .then(res => res.json())
                .then(function (result) {
                    const newWarrior = result.data
                    warrior.create({
                        id: newWarrior.id,
                        name: newWarrior.name,
                        rank: newWarrior.rank,
                        notificationsEnabled: newWarrior.notificationsEnabled,
                    })

                    eventTag('register_guest', 'engagement', 'success', () => {
                        router.route(targetPage(), true)
                    })
                })
                .catch(function () {
                    notifications.danger(
                        $_('pages.createAccount.guestForm.createError'),
                    )
                    eventTag('register_guest', 'engagement', 'failure')
                })
        }
    }

    function createUserRegistered(
        warriorName,
        warriorEmail,
        warriorPassword1,
        warriorPassword2,
    ) {
        const body = {
            name: warriorName,
            email: warriorEmail,
            password1: warriorPassword1,
            password2: warriorPassword2,
        }

        xfetch('/api/auth/register', { body })
            .then(res => res.json())
            .then(function (result) {
                const newWarrior = result.data
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name,
                    email: newWarrior.email,
                    rank: newWarrior.rank,
                    notificationsEnabled: newWarrior.notificationsEnabled,
                })

                eventTag('register_account', 'engagement', 'success', () => {
                    router.route(targetPage(), true)
                })
            })
            .catch(function () {
                notifications.danger(
                    $_('pages.createAccount.createAccountForm.createError'),
                )
                eventTag('register_account', 'engagement', 'failure')
            })
    }

    $: registerDisabled = warriorName === ''
</script>

<svelte:head>
    <title>{$_('register')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="text-center px-2 mb-4">
        <h1
            class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
            {$_('register')}
        </h1>
        {#if battleId}
            <div
                class="font-semibold font-rajdhani uppercase text-md md:text-lg mb-2 md:mb-6 md:leading-tight
                text-center dark:text-white"
            >
                {@html $_('pages.createAccount.loginForBattle', {
                    values: {
                        loginOpen: `<a href="${appRoutes.login}/battle/${battleId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                        loginClose: `</a>`,
                    },
                })}
            </div>
        {/if}
        {#if retroId}
            <div
                class="font-semibold font-rajdhani uppercase text-md md:text-lg mb-2 md:mb-6 md:leading-tight
                text-center dark:text-white"
            >
                {@html $_('loginForRetro', {
                    values: {
                        loginOpen: `<a href="${appRoutes.login}/retro/${retroId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                        loginClose: `</a>`,
                    },
                })}
            </div>
        {/if}
        {#if storyboardId}
            <div
                class="font-semibold font-rajdhani uppercase text-md md:text-lg mb-2 md:mb-6 md:leading-tight
                text-center dark:text-white"
            >
                {@html $_('loginForStoryboard', {
                    values: {
                        loginOpen: `<a href="${appRoutes.login}/storyboard/${storyboardId}" class="font-bold text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600">`,
                        loginClose: `</a>`,
                    },
                })}
            </div>
        {/if}
    </div>
    <div class="flex flex-wrap justify-center">
        {#if !$warrior.id && (guestsAllowed || registrationAllowed)}
            <div class="w-full md:w-1/2 px-4">
                <form
                    on:submit="{createUserGuest}"
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 md:p-6 mb-4"
                    name="registerGuest"
                >
                    <h2
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl b-4 mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                    >
                        {$_('pages.createAccount.guestForm.title')}
                    </h2>

                    <div class="mb-6">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourName1"
                        >
                            {$_(
                                'pages.createAccount.guestForm.fields.name.label',
                            )}
                        </label>
                        <input
                            bind:value="{warriorName}"
                            placeholder="{$_(
                                'pages.createAccount.guestForm.fields.name.placeholder',
                            )}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                            id="yourName1"
                            name="yourName1"
                            required
                        />
                    </div>
                    <div>
                        <div class="text-right">
                            <SolidButton
                                type="submit"
                                disabled="{registerDisabled}"
                            >
                                {$_('pages.createAccount.guestForm.saveButton')}
                            </SolidButton>
                        </div>
                    </div>
                </form>
            </div>
        {/if}

        {#if registrationAllowed}
            <div class="w-full md:w-1/2 px-4">
                <div
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 md:p-6 mb-4"
                >
                    <h2
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                    >
                        {@html $_(
                            'pages.createAccount.createAccountForm.title',
                            {
                                values: {
                                    optionalOpen: `<span class="text-gray-500">`,
                                    optionalClose: `</span>`,
                                },
                            },
                        )}
                    </h2>

                    <WarriorRegisterForm
                        guestWarriorsName="{warriorName}"
                        handleSubmit="{createUserRegistered}"
                        notifications="{notifications}"
                    />
                </div>
            </div>
        {:else}
            <div class="w-full md:w-1/2 px-4">
                <h2
                    class="font-bold text-2xl md:text-3xl md:leading-tight
                    text-center dark:text-white"
                >
                    {$_('pages.createAccount.registrationDisabled')}
                </h2>
            </div>
        {/if}
    </div>
</PageLayout>

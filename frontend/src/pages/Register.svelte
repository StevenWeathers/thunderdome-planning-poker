<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import WarriorRegisterForm from '../components/WarriorRegisterForm.svelte'
    import { warrior } from '../stores.js'
    import { validateName, validatePasswords } from '../validationUtils.js'
    import { _ } from '../i18n'
    import { appRoutes } from '../config'

    export let router
    export let xfetch
    export let notifications
    export let eventTag
    export let battleId

    const guestsAllowed = appConfig.AllowGuests
    const registrationAllowed = appConfig.AllowRegistration

    let warriorName = $warrior.name || ''
    let warriorEmail = ''
    let warriorPassword1 = ''
    let warriorPassword2 = ''

    $: targetPage = battleId
        ? `${appRoutes.battle}/${battleId}`
        : appRoutes.battles

    function createWarriorPrivate(e) {
        e.preventDefault()
        const body = {
            warriorName,
        }
        const validName = validateName(warriorName)

        let noFormErrors = true

        if (!validName.valid) {
            noFormErrors = false
            notifications.danger(validName.error, 1500)
        }

        if (noFormErrors) {
            xfetch('/api/warrior', { body })
                .then(res => res.json())
                .then(function(newWarrior) {
                    warrior.create({
                        id: newWarrior.id,
                        name: newWarrior.name,
                        rank: newWarrior.rank,
                        notificationsEnabled: newWarrior.notificationsEnabled,
                    })

                    eventTag('register_guest', 'engagement', 'success', () => {
                        router.route(targetPage, true)
                    })
                })
                .catch(function(error) {
                    notifications.danger(
                        $_('pages.createAccount.guestForm.createError'),
                    )
                    eventTag('register_guest', 'engagement', 'failure')
                })
        }
    }

    function createWarriorCorporal(
        warriorName,
        warriorEmail,
        warriorPassword1,
        warriorPassword2,
    ) {
        const body = {
            warriorName,
            warriorEmail,
            warriorPassword1,
            warriorPassword2,
        }

        xfetch('/api/enlist', { body })
            .then(res => res.json())
            .then(function(newWarrior) {
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name,
                    email: newWarrior.email,
                    rank: newWarrior.rank,
                    notificationsEnabled: newWarrior.notificationsEnabled,
                })

                eventTag('register_account', 'engagement', 'success', () => {
                    router.route(targetPage, true)
                })
            })
            .catch(function(error) {
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
        <h1 class="text-3xl md:text-4xl font-bold">
            {$_('pages.createAccount.title')}
        </h1>
        {#if battleId}
            <div
                class="font-bold text-m md:text-l mb-2 md:mb-6 md:leading-tight
                text-center">
                {@html $_('pages.createAccount.loginForBattle', {
                    values: {
                        loginOpen: `<a href="${appRoutes.login}/${battleId}" class="font-bold text-blue-500 hover:text-blue-800">`,
                        loginClose: `</a>`,
                    },
                })}
            </div>
        {/if}
    </div>
    <div class="flex flex-wrap justify-center">
        {#if !$warrior.id && guestsAllowed && registrationAllowed}
            <div class="w-full md:w-1/2 px-4">
                <form
                    on:submit="{createWarriorPrivate}"
                    class="bg-white shadow-lg rounded p-4 md:p-6 mb-4"
                    name="registerGuest">
                    <h2
                        class="font-bold text-xl md:text-2xl b-4 mb-2 md:mb-6
                        md:leading-tight text-center">
                        {$_('pages.createAccount.guestForm.title')}
                    </h2>

                    <div class="mb-6">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourName1">
                            {$_('pages.createAccount.guestForm.fields.name.label')}
                        </label>
                        <input
                            bind:value="{warriorName}"
                            placeholder="{$_('pages.createAccount.guestForm.fields.name.placeholder')}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourName1"
                            name="yourName1"
                            required />
                    </div>
                    <div>
                        <div class="text-right">
                            <SolidButton
                                type="submit"
                                disabled="{registerDisabled}">
                                {$_('pages.createAccount.guestForm.saveButton')}
                            </SolidButton>
                        </div>
                    </div>
                </form>
            </div>
        {/if}

        {#if registrationAllowed}
            <div class="w-full md:w-1/2 px-4">
                <div class="bg-white shadow-lg rounded p-4 md:p-6 mb-4">
                    <h2
                        class="font-bold text-xl md:text-2xl mb-2 md:mb-6
                        md:leading-tight text-center">
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
                        handleSubmit="{createWarriorCorporal}"
                        {notifications} />
                </div>
            </div>
        {:else}
            <div class="w-full md:w-1/2 px-4">
                <h2
                    class="font-bold text-2xl md:text-3xl md:leading-tight
                    text-center">
                    {$_('pages.createAccount.registrationDisabled')}
                </h2>
            </div>
        {/if}
    </div>
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import DeleteConfirmation from '../components/DeleteConfirmation.svelte'
    import ProfileForm from '../components/user/ProfileForm.svelte'
    import CreateApiKey from '../components/user/CreateApiKey.svelte'
    import CheckIcon from '../components/icons/CheckIcon.svelte'
    import { warrior } from '../stores.js'
    import { validatePasswords } from '../validationUtils.js'
    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    let warriorProfile = {}
    let apiKeys = []
    let showApiKeyCreate = false
    let showAccountDeletion = false

    let updatePassword = false
    let warriorPassword1 = ''
    let warriorPassword2 = ''

    const { ExternalAPIEnabled, LdapEnabled } = AppConfig

    function toggleUpdatePassword() {
        updatePassword = !updatePassword
        eventTag(
            'update_password_toggle',
            'engagement',
            `update: ${updatePassword}`,
        )
    }

    function getProfile() {
        xfetch(`/api/users/${$warrior.id}`)
            .then(res => res.json())
            .then(function (result) {
                warriorProfile = result.data
            })
            .catch(function () {
                notifications.danger($_('pages.warriorProfile.errorRetreiving'))
                eventTag('fetch_profile', 'engagement', 'failure')
            })
    }

    function updateWarriorProfile(p) {
        const body = {
            ...p,
        }

        xfetch(`/api/users/${$warrior.id}`, { body, method: 'PUT' })
            .then(res => res.json())
            .then(function () {
                warrior.update({
                    id: warriorProfile.id,
                    name: p.name,
                    email: warriorProfile.email,
                    rank: warriorProfile.rank,
                    avatar: p.avatar,
                    notificationsEnabled: p.notificationsEnabled,
                    locale: p.locale,
                })

                notifications.success($_('pages.warriorProfile.updateSuccess'))
                eventTag('update_profile', 'engagement', 'success')
            })
            .catch(function () {
                notifications.danger($_('pages.warriorProfile.errorUpdating'))
                eventTag('update_profile', 'engagement', 'failure')
            })
    }

    function updateWarriorPassword(e) {
        e.preventDefault()
        const body = {
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
            xfetch('/api/auth/update-password', { body, method: 'PATCH' })
                .then(function () {
                    notifications.success(
                        $_('pages.warriorProfile.passwordUpdated'),
                        1500,
                    )
                    updatePassword = false
                    eventTag('update_password', 'engagement', 'success')
                })
                .catch(function () {
                    notifications.danger(
                        $_('pages.warriorProfile.passwordUpdateError'),
                    )
                    eventTag('update_password', 'engagement', 'failure')
                })
        }
    }

    function getApiKeys() {
        xfetch(`/api/users/${$warrior.id}/apikeys`)
            .then(res => res.json())
            .then(function (result) {
                apiKeys = result.data
            })
            .catch(function () {
                notifications.danger(
                    $_('pages.warriorProfile.apiKeys.errorRetreiving'),
                )
                eventTag('fetch_profile_apikeys', 'engagement', 'failure')
            })
    }

    function deleteApiKey(apk) {
        return function () {
            xfetch(`/api/users/${$warrior.id}/apikeys/${apk}`, {
                method: 'DELETE',
            })
                .then(res => res.json())
                .then(function (result) {
                    notifications.success(
                        $_('pages.warriorProfile.apiKeys.deleteSuccess'),
                    )
                    apiKeys = result.data
                })
                .catch(function () {
                    notifications.danger(
                        $_('pages.warriorProfile.apiKeys.deleteFailed'),
                    )
                })
        }
    }

    function toggleApiKeyActiveStatus(apk, active) {
        return function () {
            const body = {
                active: !active,
            }

            xfetch(`/api/users/${$warrior.id}/apikeys/${apk}`, {
                body,
                method: 'PUT',
            })
                .then(res => res.json())
                .then(function (result) {
                    notifications.success(
                        $_('pages.warriorProfile.apiKeys.updateSuccess'),
                    )
                    apiKeys = result.data
                })
                .catch(function () {
                    notifications.danger(
                        $_('pages.warriorProfile.apiKeys.updateFailed'),
                    )
                })
        }
    }

    function toggleCreateApiKey() {
        showApiKeyCreate = !showApiKeyCreate
    }

    function handleDeleteAccount() {
        xfetch(`/api/users/${$warrior.id}`, { method: 'DELETE' })
            .then(function () {
                warrior.delete()

                eventTag('delete_warrior', 'engagement', 'success')

                router.route(appRoutes.landing)
            })
            .catch(function () {
                notifications.danger($_('pages.warriorProfile.delete.error'))
                eventTag('delete_warrior', 'engagement', 'failure')
            })
    }

    function toggleDeleteAccount() {
        showAccountDeletion = !showAccountDeletion
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.register)
            return
        }

        getProfile()
        if (ExternalAPIEnabled) {
            getApiKeys()
        }
    })

    $: updatePasswordDisabled =
        warriorPassword1 === '' || warriorPassword2 === '' || LdapEnabled
</script>

<svelte:head>
    <title>{$_('pages.warriorProfile.title')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <h2
        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight dark:text-white"
    >
        {$_('pages.warriorProfile.title')}
    </h2>

    <div class="flex justify-center flex-wrap">
        <div class="w-full md:w-1/2 lg:w-1/3">
            {#if !updatePassword}
                <div
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 md:p-6 mb-4"
                >
                    <ProfileForm
                        profile="{warriorProfile}"
                        handleUpdate="{updateWarriorProfile}"
                        toggleUpdatePassword="{toggleUpdatePassword}"
                        xfetch="{xfetch}"
                        notifications="{notifications}"
                        eventTag="{eventTag}"
                    />
                </div>
            {/if}

            {#if updatePassword}
                <form
                    on:submit="{updateWarriorPassword}"
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                    name="updateWarriorPassword"
                >
                    <div
                        class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
                    >
                        {$_('pages.warriorProfile.updatePasswordForm.title')}
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="yourPassword1"
                        >
                            {$_(
                                'pages.warriorProfile.updatePasswordForm.fields.password.label',
                            )}
                        </label>
                        <input
                            bind:value="{warriorPassword1}"
                            placeholder="{$_(
                                'pages.warriorProfile.updatePasswordForm.fields.password.placeholder',
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
                                'pages.warriorProfile.updatePasswordForm.fields.confirmPassword.label',
                            )}
                        </label>
                        <input
                            bind:value="{warriorPassword2}"
                            placeholder="{$_(
                                'pages.warriorProfile.updatePasswordForm.fields.confirmPassword.placeholder',
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
                        <button
                            type="button"
                            class="inline-block align-baseline font-bold text-sm
                            text-blue-500 hover:text-blue-800 mr-4"
                            on:click="{toggleUpdatePassword}"
                        >
                            {$_(
                                'pages.warriorProfile.updatePasswordForm.cancelButton',
                            )}
                        </button>
                        <SolidButton
                            type="submit"
                            disabled="{updatePasswordDisabled}"
                        >
                            {$_(
                                'pages.warriorProfile.updatePasswordForm.saveButton',
                            )}
                        </SolidButton>
                    </div>
                </form>
            {/if}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/3">
            {#if ExternalAPIEnabled}
                <div class="ml-8">
                    <div class="flex w-full">
                        <div class="flex-1">
                            <h2
                                class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
                            >
                                {$_('pages.warriorProfile.apiKeys.title')}
                            </h2>
                        </div>
                        <div class="flex-1">
                            <div class="text-right">
                                <HollowButton
                                    href="/swagger/index.html"
                                    options="{{ target: '_blank' }}"
                                    color="blue"
                                >
                                    {$_('apiDocumentation')}
                                </HollowButton>
                                <HollowButton
                                    onClick="{toggleCreateApiKey}"
                                    testid="apikey-create"
                                >
                                    {$_(
                                        'pages.warriorProfile.apiKeys.createButton',
                                    )}
                                </HollowButton>
                            </div>
                        </div>
                    </div>

                    <div class="flex flex-col">
                        <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                            <div
                                class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8"
                            >
                                <div
                                    class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg"
                                >
                                    <table
                                        class="min-w-full divide-y divide-gray-200 dark:divide-gray-700"
                                    >
                                        <thead
                                            class="bg-gray-50 dark:bg-gray-800"
                                        >
                                            <tr>
                                                <th
                                                    scope="col"
                                                    class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                                >
                                                    {$_('name')}
                                                </th>
                                                <th
                                                    scope="col"
                                                    class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                                >
                                                    {$_(
                                                        'pages.warriorProfile.apiKeys.prefix',
                                                    )}
                                                </th>
                                                <th
                                                    scope="col"
                                                    class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                                >
                                                    {$_(
                                                        'pages.warriorProfile.apiKeys.active',
                                                    )}
                                                </th>
                                                <th
                                                    scope="col"
                                                    class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                                                >
                                                    {$_(
                                                        'pages.warriorProfile.apiKeys.updated',
                                                    )}
                                                </th>
                                                <th
                                                    scope="col"
                                                    class="relative px-6 py-3"
                                                >
                                                    <span class="sr-only"
                                                        >Actions</span
                                                    >
                                                </th>
                                            </tr>
                                        </thead>
                                        <tbody
                                            class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                                        >
                                            {#each apiKeys as apk, i}
                                                <tr
                                                    class:bg-slate-100="{i %
                                                        2 !==
                                                        0}"
                                                    class:dark:bg-gray-800="{i %
                                                        2 !==
                                                        0}"
                                                >
                                                    <td
                                                        class="px-6 py-4 whitespace-nowrap"
                                                        data-testid="apikey-name"
                                                        >{apk.name}</td
                                                    >
                                                    <td
                                                        class="px-6 py-4 whitespace-nowrap"
                                                        data-testid="apikey-prefix"
                                                    >
                                                        {apk.prefix}
                                                    </td>
                                                    <td
                                                        class="px-6 py-4 whitespace-nowrap"
                                                        data-testid="apikey-active"
                                                        data-active="{apk.active}"
                                                    >
                                                        {#if apk.active}
                                                            <span
                                                                class="text-green-600"
                                                                ><CheckIcon
                                                                /></span
                                                            >
                                                        {/if}
                                                    </td>
                                                    <td
                                                        class="px-6 py-4 whitespace-nowrap"
                                                    >
                                                        {new Date(
                                                            apk.updatedDate,
                                                        ).toLocaleString()}
                                                    </td>
                                                    <td
                                                        class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                                                    >
                                                        <HollowButton
                                                            onClick="{toggleApiKeyActiveStatus(
                                                                apk.id,
                                                                apk.active,
                                                            )}"
                                                            testid="apikey-activetoggle"
                                                        >
                                                            {#if !apk.active}
                                                                {$_(
                                                                    'pages.warriorProfile.apiKeys.activateButton',
                                                                )}
                                                            {:else}
                                                                {$_(
                                                                    'pages.warriorProfile.apiKeys.deactivateButton',
                                                                )}
                                                            {/if}
                                                        </HollowButton>
                                                        <HollowButton
                                                            color="red"
                                                            onClick="{deleteApiKey(
                                                                apk.id,
                                                            )}"
                                                            testid="apikey-delete"
                                                        >
                                                            {$_(
                                                                'pages.warriorProfile.apiKeys.deleteButton',
                                                            )}
                                                        </HollowButton>
                                                    </td>
                                                </tr>
                                            {/each}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        </div>

        <div class="w-full text-center mt-8">
            <HollowButton onClick="{toggleDeleteAccount}" color="red">
                {$_('pages.warriorProfile.delete.deleteButton')}
            </HollowButton>
        </div>
    </div>
    {#if showApiKeyCreate}
        <CreateApiKey
            toggleCreateApiKey="{toggleCreateApiKey}"
            handleApiKeyCreate="{getApiKeys}"
            notifications="{notifications}"
            xfetch="{xfetch}"
            eventTag="{eventTag}"
        />
    {/if}

    {#if showAccountDeletion}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteAccount}"
            handleDelete="{handleDeleteAccount}"
            confirmText="{$_('pages.warriorProfile.delete.warningStatement')}"
            confirmBtnText="{$_('pages.warriorProfile.delete.confirmButton')}"
        />
    {/if}
</PageLayout>

<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'
    import { validateName, validatePasswords } from '../validationUtils.js'

    export let router
    export let notifications

    let warriorProfile = {}

    let updatePassword = false
    let warriorPassword1 = ''
    let warriorPassword2 = ''

    function toggleUpdatePassword() {
        updatePassword = !updatePassword
    }

    fetch(`/api/warrior/${$warrior.id}`, {
        method: 'GET',
        credentials: 'same-origin',
    })
        .then(function(response) {
            if (!response.ok) {
                throw Error(response.statusText)
            }
            return response
        })
        .then(function(response) {
            return response.json()
        })
        .then(function(wp) {
            warriorProfile = wp
        })
        .catch(function(error) {
            notifications.danger('Error getting your profile')
        })

    function updateWarriorProfile(e) {
        e.preventDefault()
        const body = {
            warriorName: warriorProfile.name,
        }
        const validName = validateName(body.warriorName)

        let noFormErrors = true

        if (!validName.valid) {
            noFormErrors = false
            notifications.danger(validName.error, 1500)
        }

        if (noFormErrors) {
            fetch(`/api/warrior/${$warrior.id}`, {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(body),
            })
                .then(function(response) {
                    if (!response.ok) {
                        throw Error(response.statusText)
                    }
                    return response
                })
                .then(function(updatedWarrior) {
                    warrior.update({
                        id: warriorProfile.id,
                        name: warriorProfile.name,
                        email: warriorProfile.email,
                        rank: warriorProfile.rank,
                    })

                    notifications.success('Profile updated.', 1500)
                })
                .catch(function(error) {
                    notifications.danger(
                        'Error encountered updating your profile',
                    )
                })
        }
    }

    function updateWarriorPassword(e) {
        e.preventDefault()
        const body = {
            warriorPassword1,
            warriorPassword2,
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
            fetch('/api/auth/update-password', {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(body),
            })
                .then(function(response) {
                    if (!response.ok) {
                        throw Error(response.statusText)
                    }
                    return response
                })
                .then(function() {
                    notifications.success('Password updated.', 1500)
                    updatePassword = false
                })
                .catch(function(error) {
                    notifications.danger(
                        'Error encountered attempting to update password',
                    )
                })
        }
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route('/enlist')
        }
    })

    $: updateDisabled = warriorProfile.name === ''
    $: updatePasswordDisabled =
        warriorPassword1 === '' || warriorPassword2 === ''
</script>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            {#if !updatePassword}
                <form
                    on:submit="{updateWarriorProfile}"
                    class="bg-white shadow-lg rounded p-4 md:p-6 mb-4"
                    name="updateProfile">
                    <h2
                        class="font-bold text-xl md:text-2xl mb-2 md:mb-6
                        md:leading-tight">
                        Your Profile
                    </h2>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourName">
                            Name
                        </label>
                        <input
                            bind:value="{warriorProfile.name}"
                            placeholder="Enter your name"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourName"
                            name="yourName"
                            type="text"
                            required />
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourEmail">
                            Email
                        </label>
                        <input
                            bind:value="{warriorProfile.email}"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            cursor-not-allowed"
                            id="yourEmail"
                            name="yourEmail"
                            type="email"
                            disabled />
                    </div>

                    <div>
                        <div class="text-right">
                            <button
                                type="button"
                                class="inline-block align-baseline font-bold
                                text-sm text-blue-500 hover:text-blue-800 mr-4"
                                on:click="{toggleUpdatePassword}">
                                Update Password
                            </button>
                            <SolidButton
                                type="submit"
                                disabled="{updateDisabled}">
                                Update Profile
                            </SolidButton>
                        </div>
                    </div>
                </form>
            {/if}

            {#if updatePassword}
                <form
                    on:submit="{updateWarriorPassword}"
                    class="bg-white shadow-lg rounded p-6 mb-4"
                    name="updateWarriorPassword">
                    <div
                        class="font-bold text-xl md:text-2xl mb-2 md:mb-6
                        md:leading-tight text-center">
                        Update Password
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourPassword1">
                            Password
                        </label>
                        <input
                            bind:value="{warriorPassword1}"
                            placeholder="Enter a password"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourPassword1"
                            name="yourPassword1"
                            type="password"
                            required />
                    </div>

                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="yourPassword2">
                            Confirm Password
                        </label>
                        <input
                            bind:value="{warriorPassword2}"
                            placeholder="Confirm your password"
                            class="bg-gray-200 border-gray-200 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-purple-500"
                            id="yourPassword2"
                            name="yourPassword2"
                            type="password"
                            required />
                    </div>

                    <div class="text-right">
                        <button
                            type="button"
                            class="inline-block align-baseline font-bold text-sm
                            text-blue-500 hover:text-blue-800 mr-4"
                            on:click="{toggleUpdatePassword}">
                            Cancel
                        </button>
                        <SolidButton
                            type="submit"
                            disabled="{updatePasswordDisabled}">
                            Update
                        </SolidButton>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</PageLayout>

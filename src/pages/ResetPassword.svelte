<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'

    export let router
    export let notifications
    export let resetId

    const nameMin = 1
    const nameMax = 64
    const passMin = 6
    const passMax = 72
    const emailMax = 320

    let warriorPassword1 = ''
    let warriorPassword2 = ''


    function resetWarriorPassword(e) {
        e.preventDefault()
        const body = {
            resetId,
            warriorPassword1,
            warriorPassword2,
        }

        let noFormErrors = true

        if (
            warriorPassword1.length < passMin ||
            warriorPassword1.length > passMax
        ) {
            noFormErrors = false
            notifications.danger(
                `Password must be between ${passMin} and ${passMax} characters.`,
            )
        }

        if (warriorPassword1 !== warriorPassword2) {
            noFormErrors = false
            notifications.danger(`Password and Confirm Password do not match.`)
        }

        if (noFormErrors) {
            fetch('/api/auth/reset-password', {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(body),
            })
                .then(function(response) {
                    if (!response.ok) {
                        throw Error(response.statusText);
                    }
                    return response;
                })
                .then(function() {
                    router.route('/login', true)
                })
                .catch(function(error) {
                    notifications.danger(
                        'Error encountered attempting to reset password',
                    )
                })
        }
    }

    $: resetDisabled =
        warriorPassword1 === '' ||
        warriorPassword2 === ''
</script>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 lg:w-1/3">
            <form
                on:submit="{resetWarriorPassword}"
                class="bg-white shadow-lg rounded p-6 mb-4"
                name="resetWarriorPassword">
                <div
                    class="font-bold text-xl md:text-2xl mb-2 md:mb-6 md:leading-tight text-center"
                >
                    Reset Password
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
                        class="bg-gray-200 border-gray-200 border-2 appearance-none rounded w-full py-2
                        px-3 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500"
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
                        class="bg-gray-200 border-gray-200 border-2 appearance-none rounded w-full py-2
                        px-3 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-purple-500"
                        id="yourPassword2"
                        name="yourPassword2"
                        type="password"
                        required />
                </div>

                <div class="text-right">
                    <SolidButton type="submit" disabled="{resetDisabled}">
                        Reset
                    </SolidButton>
                </div>
            </form>
        </div>
    </div>
</PageLayout>

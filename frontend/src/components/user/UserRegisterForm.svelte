<script lang="ts">
    import { validateName, validatePasswords } from '../../validationUtils'
    import LL from '../../i18n/i18n-svelte'
    import SolidButton from '../SolidButton.svelte'

    export let notifications
    export let handleSubmit
    export let guestWarriorsName = ''

    let warriorName = guestWarriorsName
    let warriorEmail = ''
    let warriorPassword1 = ''
    let warriorPassword2 = ''

    function onSubmit(e) {
        e.preventDefault()

        const validName = validateName(warriorName)
        const validPasswords = validatePasswords(
            warriorPassword1,
            warriorPassword2,
        )

        let noFormErrors = true

        if (!validName.valid) {
            noFormErrors = false
            notifications.danger(validName.error, 1500)
        }

        if (!validPasswords.valid) {
            noFormErrors = false
            notifications.danger(validPasswords.error, 1500)
        }

        if (noFormErrors) {
            handleSubmit(
                warriorName,
                warriorEmail,
                warriorPassword1,
                warriorPassword2,
            )
        }
    }

    $: createDisabled =
        warriorName === '' ||
        warriorEmail === '' ||
        warriorPassword1 === '' ||
        warriorPassword2 === ''
</script>

<form on:submit="{onSubmit}" name="createAccount">
    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourName2"
        >
            {$LL.name()}
        </label>
        <input
            bind:value="{warriorName}"
            placeholder="{$LL.userNamePlaceholder()}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            id="yourName2"
            name="yourName2"
            required
        />
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourEmail"
        >
            {$LL.email()}
        </label>
        <input
            bind:value="{warriorEmail}"
            placeholder="{$LL.enterYourEmail()}"
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
            for="yourPassword1"
        >
            {$LL.password()}
        </label>
        <input
            bind:value="{warriorPassword1}"
            placeholder="{$LL.passwordPlaceholder()}"
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
            {$LL.confirmPassword()}
        </label>
        <input
            bind:value="{warriorPassword2}"
            placeholder="{$LL.confirmPasswordPlaceholder()}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            id="yourPassword2"
            name="yourPassword2"
            type="password"
            required
        />
    </div>

    <div>
        <div class="text-right">
            <SolidButton type="submit" disabled="{createDisabled}">
                {$LL.create()}
            </SolidButton>
        </div>
    </div>
</form>

<script>
    import SolidButton from '../SolidButton.svelte'
    import { validateName, validatePasswords } from '../../validationUtils.js'
    import { AppConfig } from '../../config.js'
    import { _ } from '../../i18n.js'

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
            {$_('pages.createAccount.createAccountForm.fields.name.label')}
        </label>
        <input
            bind:value="{warriorName}"
            placeholder="{$_(
                'pages.createAccount.createAccountForm.fields.name.placeholder',
            )}"
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
            {$_('pages.createAccount.createAccountForm.fields.email.label')}
        </label>
        <input
            bind:value="{warriorEmail}"
            placeholder="{$_(
                'pages.createAccount.createAccountForm.fields.email.placeholder',
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
            for="yourPassword1"
        >
            {$_('pages.createAccount.createAccountForm.fields.password.label')}
        </label>
        <input
            bind:value="{warriorPassword1}"
            placeholder="{$_(
                'pages.createAccount.createAccountForm.fields.password.placeholder',
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
                'pages.createAccount.createAccountForm.fields.confirmPassword.label',
            )}
        </label>
        <input
            bind:value="{warriorPassword2}"
            placeholder="{$_(
                'pages.createAccount.createAccountForm.fields.confirmPassword.placeholder',
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

    <div>
        <div class="text-right">
            <SolidButton type="submit" disabled="{createDisabled}">
                {$_('pages.createAccount.createAccountForm.saveButton')}
            </SolidButton>
        </div>
    </div>
</form>

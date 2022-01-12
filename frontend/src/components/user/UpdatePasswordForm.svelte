<script>
    import SolidButton from '../SolidButton.svelte'
    import { _ } from '../../i18n.js'
    import { AppConfig } from '../../config.js'
    import { validatePasswords } from '../../validationUtils.js'

    export let handleUpdate = () => {}
    export let toggleForm = () => {}
    export let notifications

    const { LdapEnabled } = AppConfig

    let warriorPassword1 = ''
    let warriorPassword2 = ''

    function updateWarriorPassword(e) {
        e.preventDefault()

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
            handleUpdate(warriorPassword1, warriorPassword2)
        }
    }

    $: updatePasswordDisabled =
        warriorPassword1 === '' || warriorPassword2 === '' || LdapEnabled
</script>

<form on:submit="{updateWarriorPassword}" name="updateWarriorPassword">
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
            on:click="{toggleForm}"
        >
            {$_('pages.warriorProfile.updatePasswordForm.cancelButton')}
        </button>
        <SolidButton type="submit" disabled="{updatePasswordDisabled}">
            {$_('pages.warriorProfile.updatePasswordForm.saveButton')}
        </SolidButton>
    </div>
</form>

<script>
    import DownCarrotIcon from '../icons/ChevronDown.svelte'
    import WarriorAvatar from './UserAvatar.svelte'
    import LocaleSwitcher from '../LocaleSwitcher.svelte'
    import SolidButton from '../SolidButton.svelte'
    import VerifiedIcon from '../icons/Verified.svelte'
    import HollowButton from '../HollowButton.svelte'
    import SetupMFA from '../user/SetupMFA.svelte'
    import DeleteConfirmation from '../DeleteConfirmation.svelte'
    import { countryList } from '../../country.js'
    import { AppConfig } from '../../config.js'
    import { _, locale, setupI18n } from '../../i18n.js'
    import { warrior } from '../../stores.js'
    import { validateName, validateUserIsAdmin } from '../../validationUtils.js'

    export let profile = {
        id: '',
        rank: '',
        name: '',
        email: '',
        company: '',
        country: '',
        jobTitle: '',
        notificationsEnabled: true,
        avatar: '',
        gravatarHash: '',
        verified: false,
        mfaEnabled: false,
    }
    export let handleUpdate = () => {}
    export let toggleUpdatePassword
    export let eventTag
    export let notifications
    export let xfetch
    export let ldapEnabled

    const { AvatarService } = AppConfig

    const configurableAvatarServices = [
        'dicebear',
        'gravatar',
        'robohash',
        'govatar',
    ]
    const isAvatarConfigurable =
        configurableAvatarServices.includes(AvatarService)
    const avatarOptions = {
        dicebear: [
            'male',
            'female',
            'human',
            'identicon',
            'bottts',
            'avataaars',
            'jdenticon',
            'gridy',
            'code',
        ],
        gravatar: [
            'mp',
            'identicon',
            'monsterid',
            'wavatar',
            'retro',
            'robohash',
        ],
        robohash: ['set1', 'set2', 'set3', 'set4'],
        govatar: ['male', 'female'],
    }

    let avatars = isAvatarConfigurable ? avatarOptions[AvatarService] : []

    function handleSubmit(e) {
        e.preventDefault()
        const validName = validateName(profile.name)
        let p = {
            name: profile.name,
            country: profile.country,
            company: profile.company,
            jobTitle: profile.jobTitle,
            notificationsEnabled: profile.notificationsEnabled,
            avatar: profile.avatar,
            locale: $locale,
        }

        if (userIsAdmin) {
            p.email = profile.email
        }

        if (!validName.valid) {
            notifications.danger(validName.error, 1500)
        } else {
            handleUpdate(p)
        }
    }

    let showMFASetup = false

    function toggleMfaSetup() {
        showMFASetup = !showMFASetup
    }

    function handleMfaSetupCompletion() {
        profile.mfaEnabled = true
        toggleMfaSetup()
    }

    let showMfaRemove = false

    function toggleMfaRemove() {
        showMfaRemove = !showMfaRemove
    }

    function handleMfaRemove() {
        xfetch('/api/auth/mfa', { method: 'DELETE' })
            .then(() => {
                profile.mfaEnabled = false
                toggleMfaRemove()
                notifications.success($_('mfa2faRemoveSuccess'))
            })
            .catch(() => {
                notifications.danger($_('mfa2faRemoveFailure'))
            })
    }

    function requestVerifyEmail(e) {
        e.preventDefault()
        xfetch(`/api/users/${profile.id}/request-verify`, { method: 'POST' })
            .then(function () {
                eventTag('user_verify_request', 'engagement', 'success')

                notifications.success($_('requestVerifyEmailSuccess'))
            })
            .catch(function () {
                notifications.danger($_('requestVerifyEmailFailure'))
                eventTag('user_verify_request', 'engagement', 'failure')
            })
    }

    $: updateDisabled = profile.name === ''
    $: userIsAdmin = validateUserIsAdmin($warrior)
</script>

<form on:submit="{handleSubmit}" name="updateProfile">
    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourName"
        >
            {$_('pages.warriorProfile.fields.name.label')}
        </label>
        <input
            bind:value="{profile.name}"
            placeholder="{$_('pages.warriorProfile.fields.name.placeholder')}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            id="yourName"
            name="yourName"
            type="text"
            disabled="{ldapEnabled}"
            required
        />
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourEmail"
        >
            {$_('pages.warriorProfile.fields.email.label')}
            {#if profile.verified}
                <span
                    class="font-bold text-green-600
                                    border-green-500 border py-1 px-2 rounded
                                    ml-1"
                    data-testid="user-verified"
                >
                    {$_('pages.warriorProfile.fields.email.verified')}
                    <VerifiedIcon class="inline fill-current h-4 w-4" />
                </span>
            {:else if profile.rank !== 'GUEST'}
                <button
                    class=" float-right inline-block align-baseline font-bold text-sm text-blue-500
                                        hover:text-blue-800"
                    on:click="{requestVerifyEmail}"
                    data-testid="request-verify"
                    type="button"
                    >{$_('requestVerifyEmail')}
                </button>
            {/if}
        </label>
        <input
            bind:value="{profile.email}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            class:cursor-not-allowed="{!userIsAdmin}"
            id="yourEmail"
            name="yourEmail"
            type="email"
            disabled="{!userIsAdmin}"
        />
    </div>

    {#if profile.rank !== 'GUEST'}
        <div class="mb-4">
            <p class="block text-gray-700 dark:text-gray-400 font-bold mb-2">
                {$_('mfa2faLabel')}
            </p>
            {#if !profile.mfaEnabled}
                <HollowButton color="teal" onClick="{toggleMfaSetup}"
                    >{$_('mfa2faSetup')}
                </HollowButton>
            {:else}
                <HollowButton color="red" onClick="{toggleMfaRemove}"
                    >{$_('mfa2faRemove')}
                </HollowButton>
            {/if}
        </div>
    {/if}

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourCountry"
        >
            {$_('pages.warriorProfile.fields.country.label')}
        </label>

        <div class="relative">
            <select
                bind:value="{profile.country}"
                class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
                id="yourCountry"
                name="yourCountry"
            >
                <option value="">
                    {$_('pages.warriorProfile.fields.country.placeholder')}
                </option>
                {#each countryList as item}
                    <option value="{item.abbrev}">
                        {item.name} [{item.abbrev}]
                    </option>
                {/each}
            </select>
            <div
                class="pointer-events-none absolute inset-y-0
                                right-0 flex items-center px-2 text-gray-700 dark:text-gray-300"
            >
                <DownCarrotIcon />
            </div>
        </div>
    </div>

    <div class="mb-4">
        <div class="text-gray-700 dark:text-gray-400 font-bold mb-2">
            {$_('pages.warriorProfile.fields.locale.label')}
        </div>
        <LocaleSwitcher
            selectedLocale="{$locale}"
            on:locale-changed="{e =>
                setupI18n({
                    withLocale: e.detail,
                })}"
        />
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourCompany"
        >
            {$_('pages.warriorProfile.fields.company.label')}
        </label>
        <input
            bind:value="{profile.company}"
            placeholder="{$_(
                'pages.warriorProfile.fields.company.placeholder',
            )}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            id="yourCompany"
            name="yourCompany"
            type="text"
        />
    </div>

    <div class="mb-4">
        <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="yourJobTitle"
        >
            {$_('pages.warriorProfile.fields.jobTitle.label')}
        </label>
        <input
            bind:value="{profile.jobTitle}"
            placeholder="{$_(
                'pages.warriorProfile.fields.jobTitle.placeholder',
            )}"
            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
            id="yourJobTitle"
            name="yourJobTitle"
            type="text"
        />
    </div>

    <div class="mb-4">
        <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2">
            <input
                bind:checked="{profile.notificationsEnabled}"
                type="checkbox"
                class="w-4 h-4 dark:accent-lime-400 mr-1"
            />
            <span>
                {$_('pages.warriorProfile.fields.enable_notifications.label')}
            </span>
        </label>
    </div>

    {#if isAvatarConfigurable}
        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 font-bold
                                mb-2"
                for="yourAvatar"
            >
                {$_('pages.warriorProfile.fields.avatar.label')}
            </label>
            <div class="flex">
                <div class="md:w-2/3 lg:w-3/4">
                    <div
                        class="relative"
                        class:hidden="{AvatarService === 'gravatar' &&
                            profile.email !== ''}"
                    >
                        <select
                            bind:value="{profile.avatar}"
                            class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
                            id="yourAvatar"
                            name="yourAvatar"
                        >
                            {#each avatars as item}
                                <option value="{item}">
                                    {item}
                                </option>
                            {/each}
                        </select>
                        <div
                            class="pointer-events-none absolute
                                            inset-y-0 right-0 flex items-center
                                            px-2 text-gray-700 dark:text-gray-300"
                        >
                            <DownCarrotIcon />
                        </div>
                    </div>
                </div>
                <div class="md:w-1/3 lg:w-1/4 ml-1">
                    <span class="float-right">
                        <WarriorAvatar
                            warriorId="{profile.id}"
                            avatar="{profile.avatar}"
                            gravatarHash="{profile.gravatarHash}"
                            width="48"
                            class="rounded-full"
                        />
                    </span>
                </div>
            </div>
        </div>
    {/if}

    <div>
        <div class="text-right">
            {#if !ldapEnabled && toggleUpdatePassword}
                <button
                    type="button"
                    class="inline-block align-baseline font-bold
                                    text-sm text-blue-500 hover:text-blue-800 mr-4"
                    on:click="{toggleUpdatePassword}"
                >
                    {$_('pages.warriorProfile.updatePasswordButton')}
                </button>
            {/if}
            <SolidButton type="submit" disabled="{updateDisabled}">
                {$_('pages.warriorProfile.saveProfileButton')}
            </SolidButton>
        </div>
    </div>
</form>

{#if showMFASetup}
    <SetupMFA
        notifications="{notifications}"
        xfetch="{xfetch}"
        eventTag="{eventTag}"
        toggleSetup="{toggleMfaSetup}"
        handleComplete="{handleMfaSetupCompletion}"
    />
{/if}

{#if showMfaRemove}
    <DeleteConfirmation
        toggleDelete="{toggleMfaRemove}"
        handleDelete="{handleMfaRemove}"
        confirmText="{$_('mfa2faRemoveText')}"
        confirmBtnText="{$_('remove')}"
    />
{/if}

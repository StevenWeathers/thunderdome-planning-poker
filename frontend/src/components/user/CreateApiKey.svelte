<script>
    import SolidButton from '../SolidButton.svelte'
    import ClipboardIcon from '../icons/ClipboardIcon.svelte'
    import Modal from '../Modal.svelte'
    import { _ } from '../../i18n.js'
    import { warrior } from '../../stores.js'

    export let handleApiKeyCreate = () => {}
    export let toggleCreateApiKey = () => {}
    export let eventTag = () => {}
    export let xfetch = () => {}
    export let notifications

    let keyName = ''
    let apiKey = ''

    function handleSubmit(event) {
        event.preventDefault()

        if (keyName === '') {
            notifications.danger(
                $_('pages.warriorProfile.apiKeys.fields.name.invalid'),
            )
            eventTag('create_api_key_name_invalid', 'engagement', 'failure')
            return false
        }

        const body = {
            name: keyName,
        }

        xfetch(`/api/users/${$warrior.id}/apikeys`, { body })
            .then(res => res.json())
            .then(function (result) {
                handleApiKeyCreate()
                apiKey = result.data.apiKey
                eventTag('create_api_key', 'engagement', 'success')
            })
            .catch(function (error, response) {
                if (Array.isArray(error)) {
                    error[1].json().then(function (result) {
                        let errMessage
                        switch (result.error) {
                            case 'USER_APIKEY_LIMIT_REACHED':
                                errMessage = $_(
                                    'pages.warriorProfile.apiKeys.limitReached',
                                )
                                break
                            case 'REQUIRES_VERIFIED_USER':
                                errMessage = $_(
                                    'pages.warriorProfile.apiKeys.unverifiedUser',
                                )
                                break
                            default:
                                errMessage = $_(
                                    'pages.warriorProfile.apiKeys.createFailed',
                                )
                        }

                        notifications.danger(errMessage)
                    })
                } else {
                    notifications.danger(
                        $_('pages.warriorProfile.apiKeys.createFailed'),
                    )
                }

                eventTag('create_api_key', 'engagement', 'failure')
            })
    }

    function copyKey() {
        const apk = document.getElementById('apiKey')

        if (!navigator.clipboard) {
            apk.select()
            document.execCommand('copy')
        } else {
            navigator.clipboard
                .writeText(apk.value)
                .then(function () {
                    notifications.success($_('apikeyCopySuccess'))
                })
                .catch(function () {
                    notifications.danger($_('apikeyCopyFailure'))
                })
        }
    }
</script>

<Modal closeModal="{toggleCreateApiKey}">
    {#if apiKey === ''}
        <form on:submit="{handleSubmit}" name="createApiKey">
            <div class="mb-4">
                <label
                    class="block dark:text-gray-400 font-bold mb-2"
                    for="keyName"
                >
                    {$_('pages.warriorProfile.apiKeys.fields.name.label')}
                </label>
                <input
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    type="text"
                    id="keyName"
                    name="keyName"
                    bind:value="{keyName}"
                    placeholder="{$_(
                        'pages.warriorProfile.apiKeys.fields.name.placeholder',
                    )}"
                    required
                />
            </div>
            <div class="text-right">
                <div>
                    <SolidButton type="submit">
                        {$_('pages.warriorProfile.apiKeys.fields.submitButton')}
                    </SolidButton>
                </div>
            </div>
        </form>
    {:else}
        <div class="mb-4">
            <p class="mb-3 mt-3 dark:text-white">
                {@html $_('pages.warriorProfile.apiKeys.createSuccess', {
                    values: {
                        keyName: `<span class="font-bold">${keyName}</span>`,
                        onlyNowOpen: '<span class="font-bold">',
                        onlyNowClose: '</span>',
                    },
                })}
            </p>
            <div class="flex flex-wrap items-stretch w-full mb-3">
                <input
                    class="flex-shrink flex-grow flex-auto leading-normal w-px
                    flex-1 border-2 h-10 bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-900 rounded
                    rounded-r-none px-4 appearance-none text-gray-800 dark:text-gray-400 font-bold
                    focus:outline-none focus:bg-white dark:focus:bg-gray-800 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    type="text"
                    value="{apiKey}"
                    id="apiKey"
                    readonly
                />
                <div class="invisible md:visible md:flex md:-mr-px">
                    <SolidButton
                        color="blue-copy"
                        onClick="{copyKey}"
                        additionalClasses="flex items-center leading-normal
                        whitespace-no-wrap text-sm"
                    >
                        <ClipboardIcon />
                    </SolidButton>
                </div>
            </div>
            <p class="dark:text-white">
                {$_('pages.warriorProfile.apiKeys.storeWarning')}
            </p>
        </div>
        <div class="text-right">
            <div>
                <SolidButton
                    onClick="{toggleCreateApiKey}"
                    testid="apikey-close"
                >
                    {$_('pages.warriorProfile.apiKeys.fields.closeButton')}
                </SolidButton>
            </div>
        </div>
    {/if}
</Modal>

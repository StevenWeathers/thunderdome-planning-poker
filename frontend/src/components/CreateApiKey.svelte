<script>
    import SolidButton from './SolidButton.svelte'
    import ClipboardIcon from './icons/ClipboardIcon.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import { _ } from '../i18n'
    import { warrior } from '../stores.js'

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
            notifications.danger('Please enter a key name')
            return false
        }

        const body = {
            name: keyName
        }

        xfetch(`/api/warrior/${$warrior.id}/apikey`, { body })
            .then(res => res.json())
            .then(function(apk) {
                handleApiKeyCreate()
                apiKey = apk.apiKey
            })
            .catch(function(error) {
                notifications.danger(
                    'Failed to create api key'
                )
            })
    }

    function copyKey() {
        const apk = document.getElementById('apiKey')
        apk.select()
        document.execCommand('copy')
    }
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/3 z-50
        max-h-full">
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{toggleCreateApiKey}"
                        class="text-gray-800">
                        <CloseIcon />
                    </button>
                </div>

                {#if apiKey === ''}
                    <form on:submit="{handleSubmit}" name="createApiKey">
                        <div class="mb-4">
                            <label
                                class="block text-sm font-bold mb-2"
                                for="keyName">
                                Key Name
                            </label>
                            <input
                                class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                type="text"
                                id="keyName"
                                name="keyName"
                                bind:value="{keyName}"
                                placeholder="Enter a key name"
                                required />
                        </div>
                        <div class="text-right">
                            <div>
                                <SolidButton type="submit">
                                    Create
                                </SolidButton>
                            </div>
                        </div>
                    </form>
                {:else}
                    <div class="mb-4">
                        <p class="mb-3 mt-3">New Api Key (<span class="font-bold">{keyName}</span>) created and <span class="font-bold">it will be displayed only now</span></p>
                        <div class="flex flex-wrap items-stretch w-full mb-3">
                            <input
                                class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
                                border-2 h-10 bg-gray-200 border-gray-200 rounded rounded-r-none px-4
                                appearance-none text-gray-800 font-bold focus:outline-none focus:bg-white
                                focus:border-purple-500 "
                                type="text"
                                value={apiKey}
                                id="apiKey"
                                readonly />
                            <div class="invisible md:visible md:flex md:-mr-px">
                                <SolidButton
                                    color="blue-copy"
                                    onClick="{copyKey}"
                                    additionalClasses="flex items-center leading-normal
                                    whitespace-no-wrap text-sm">
                                    <ClipboardIcon />
                                </SolidButton>
                            </div>
                        </div>
                        <p>
                            Please store it somewhere safe because as soon as you navigate away from this page, we will not be able to retrieve or restore this generated token.
                        </p>
                    </div>
                    <div class="text-right">
                        <div>
                            <SolidButton onClick={toggleCreateApiKey}>
                                Close
                            </SolidButton>
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    </div>
</div>

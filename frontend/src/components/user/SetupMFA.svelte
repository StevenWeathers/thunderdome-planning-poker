<script>
    import Modal from '../Modal.svelte'
    import SolidButton from '../SolidButton.svelte'

    export let toggleSetup = () => {}
    export let handleComplete = () => {}
    export let xfetch
    export let notifications
    export let eventTag

    let qrCode = ''
    let secret = ''
    let passcode = ''

    xfetch('/api/auth/mfa/setup/generate', { method: 'POST' })
        .then(res => res.json())
        .then(r => {
            qrCode = r.data.qrCode
            secret = r.data.secret
        })
        .catch(err => {
            console.log(err)
            notifications.danger('failed to start 2FA/MFA setup')
        })

    function onSubmit(e) {
        e.preventDefault()
        xfetch('/api/auth/mfa/setup/validate', { body: { secret, passcode } })
            .then(res => res.json())
            .then(r => {
                console.log(r.data.result)
                if (r.data.result === 'SUCCESS') {
                    notifications.success('successfully enabled 2FA/MFA')
                    handleComplete()
                } else {
                    notifications.danger(`${r.data.result}`)
                }
            })
            .catch(err => {
                notifications.danger('failed to validate 2FA/MFA setup token')
            })
    }

    $: submitDisabled = passcode === ''
</script>

<Modal closeModal="{toggleSetup}" widthClasses="md:w-2/3 lg:w-1/2">
    <div class="pt-12">
        <div class="dark:text-gray-300 text-center">
            <p class="font-rajdhani text-lg mb-2">
                Scan the following QR code with an Authenticator Application and
                enter the token generated
            </p>
            {#if qrCode !== ''}
                <img src="data:image/png;base64,{qrCode}" class="m-auto" />

                <p class="mt-2 font-rajdhani text-xl text-red-500">
                    Secret Key: {secret}
                </p>
            {/if}
        </div>
        <form on:submit="{onSubmit}" name="validateMFAPasscode" class="mt-8">
            <div class="mb-4">
                <label
                    class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                    for="mfaPasscode"
                >
                    Authenticator Generated Token
                </label>
                <input
                    bind:value="{passcode}"
                    placeholder="Enter generated token"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="mfaPasscode"
                    name="mfaPasscode"
                    type="password"
                    required
                />
            </div>

            <div>
                <div class="text-right">
                    <SolidButton type="submit" disabled="{submitDisabled}">
                        Confirm Token
                    </SolidButton>
                </div>
            </div>
        </form>
    </div>
</Modal>

<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'

    export let xfetch
    export let eventTag
    export let verifyId

    let accountVerified = false
    let verficationError = false

    xfetch('/api/auth/verify', { body: { verifyId } })
        .then(function() {
            accountVerified = true
            eventTag('account_verify', 'engagement', 'success')
        })
        .catch(function(error) {
            verficationError = true
            eventTag('account_verify', 'engagement', 'failure')
        })
</script>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 xl:w-1/3 py-4">
            {#if accountVerified}
                <div
                    class="bg-green-100 border border-green-400 text-green-700
                    px-4 py-3 rounded relative"
                    role="alert">
                    <strong class="font-bold">Account Verified</strong>
                    <p>Thanks for verifying your email.</p>
                </div>
            {:else if verficationError}
                <div
                    class="bg-red-100 border border-red-400 text-red-700 px-4
                    py-3 rounded relative"
                    role="alert">
                    <strong class="font-bold">Verification Failed</strong>
                    <p>
                        Something when wrong verifying your account, perhaps
                        this link expired or was already used.
                    </p>
                </div>
            {:else}
                <div class="text-center">
                    <h1 class="text-4xl text-teal-500 leading-tight font-bold">
                        Verifying Account...
                    </h1>
                </div>
            {/if}
        </div>
    </div>
</PageLayout>

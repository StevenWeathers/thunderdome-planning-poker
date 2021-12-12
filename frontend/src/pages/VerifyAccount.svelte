<script>
    import PageLayout from '../components/PageLayout.svelte'
    import { _ } from '../i18n.js'

    export let xfetch
    export let eventTag
    export let verifyId

    let accountVerified = false
    let verficationError = false

    xfetch('/api/auth/verify', { body: { verifyId }, method: 'PATCH' })
        .then(function () {
            accountVerified = true
            eventTag('account_verify', 'engagement', 'success')
        })
        .catch(function () {
            verficationError = true
            eventTag('account_verify', 'engagement', 'failure')
        })
</script>

<svelte:head>
    <title>{$_('verifyAccount')} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    <div class="flex justify-center">
        <div class="w-full md:w-1/2 xl:w-1/3 py-4">
            {#if accountVerified}
                <div
                    class="bg-green-100 border border-green-400 text-green-700
                    px-4 py-3 rounded relative"
                    role="alert"
                >
                    <strong class="font-bold">
                        {$_('pages.verifyAccount.verified.title')}
                    </strong>
                    <p>{$_('pages.verifyAccount.verified.thanks')}</p>
                </div>
            {:else if verficationError}
                <div
                    class="bg-red-100 border border-red-400 text-red-700 px-4
                    py-3 rounded relative"
                    role="alert"
                >
                    <strong class="font-bold">
                        {$_('pages.verifyAccount.failed.title')}
                    </strong>
                    <p>{$_('pages.verifyAccount.failed.error')}</p>
                </div>
            {:else}
                <div class="text-center">
                    <h1 class="text-4xl text-teal-500 leading-tight font-bold">
                        {$_('pages.verifyAccount.loading')}
                    </h1>
                </div>
            {/if}
        </div>
    </div>
</PageLayout>

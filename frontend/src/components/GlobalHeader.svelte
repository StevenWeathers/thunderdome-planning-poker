<script>
    import WarriorIcon from './icons/UserIcon.svelte'
    import HollowButton from './HollowButton.svelte'
    import LocaleSwitcher from './LocaleSwitcher.svelte'
    import { validateUserIsAdmin } from '../validationUtils'
    import { _, locale, setupI18n } from '../i18n'
    import { warrior } from '../stores'
    import { AppConfig, appRoutes } from '../config'

    export let xfetch
    export let router
    export let eventTag
    export let notifications

    const { AllowRegistration, PathPrefix } = AppConfig

    function logoutWarrior() {
        xfetch('/api/auth/logout', { method: 'DELETE' })
            .then(function () {
                eventTag('logout', 'engagement', 'success', () => {
                    warrior.delete()
                    router.route(appRoutes.landing, true)
                })
            })
            .catch(function () {
                notifications.danger($_('logoutError'))
                eventTag('logout', 'engagement', 'failure')
            })
    }
</script>

<style>
    :global(.nav-logo) {
        max-height: 3.75rem;
    }
</style>

<nav
    class="flex items-center justify-between flex-wrap bg-white px-6 py-2"
    role="navigation"
    aria-label="main navigation"
>
    <div class="flex items-center flex-shrink-0 mr-6">
        <a href="{appRoutes.landing}">
            <img
                src="{PathPrefix}/img/logo.svg"
                alt="Thunderdome"
                class="nav-logo"
            />
        </a>
    </div>
    <div class="text-right mt-4 md:mt-0">
        {#if $warrior.name}
            <span class="font-bold mr-2 text-xl">
                <WarriorIcon />
                <a href="{appRoutes.profile}" data-testid="userprofile-link"
                    >{$warrior.name}</a
                >
            </span>
            <HollowButton
                color="teal"
                href="{appRoutes.battles}"
                additionalClasses="mr-2"
            >
                {$_('pages.myBattles.nav')}
            </HollowButton>
            {#if $warrior.rank !== 'GUEST' && $warrior.rank !== 'PRIVATE'}
                <HollowButton
                    color="blue"
                    href="{appRoutes.organizations}"
                    additionalClasses="mr-2"
                >
                    {$_('organizations')} &amp; {$_('teams')}
                </HollowButton>
            {/if}
            {#if !$warrior.rank || $warrior.rank === 'GUEST' || $warrior.rank === 'PRIVATE'}
                {#if AllowRegistration}
                    <HollowButton
                        color="teal"
                        href="{appRoutes.register}"
                        additionalClasses="mr-2"
                    >
                        {$_('pages.createAccount.nav')}
                    </HollowButton>
                {/if}
                <HollowButton href="{appRoutes.login}">
                    {$_('pages.login.nav')}
                </HollowButton>
            {:else}
                {#if validateUserIsAdmin($warrior)}
                    <HollowButton
                        color="purple"
                        href="{appRoutes.admin}"
                        additionalClasses="mr-2"
                    >
                        {$_('pages.admin.nav')}
                    </HollowButton>
                {/if}
                <HollowButton color="red" onClick="{logoutWarrior}">
                    {$_('logout')}
                </HollowButton>
            {/if}
        {:else}
            {#if AllowRegistration}
                <HollowButton
                    color="teal"
                    href="{appRoutes.register}"
                    additionalClasses="mr-2"
                >
                    {$_('pages.createAccount.nav')}
                </HollowButton>
            {/if}
            <HollowButton href="{appRoutes.login}">
                {$_('pages.login.nav')}
            </HollowButton>
        {/if}
        <LocaleSwitcher
            class="ml-2"
            selectedLocale="{$locale}"
            on:locale-changed="{e =>
                setupI18n({
                    withLocale: e.detail,
                })}"
        />
    </div>
</nav>

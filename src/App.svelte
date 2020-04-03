<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'
    import Notifications from './components/Notifications.svelte'
    import WarriorIcon from './components/icons/WarriorIcon.svelte'
    import HollowButton from './components/HollowButton.svelte'

    import Landing from './pages/Landing.svelte'
    import Battles from './pages/Battles.svelte'
    import Battle from './pages/Battle.svelte'
    import Register from './pages/Register.svelte'
    import Login from './pages/Login.svelte'
    import ResetPassword from './pages/ResetPassword.svelte'
    import VerifyAccount from './pages/VerifyAccount.svelte'
    import WarriorProfile from './pages/WarriorProfile.svelte'
    import Admin from './pages/Admin.svelte'
    import { warrior } from './stores.js'

    const footerLinkClasses = 'no-underline text-teal-500 hover:text-teal-800'

    let notifications

    let activeWarrior

    const unsubscribe = warrior.subscribe(w => {
        activeWarrior = w
    })

    let currentPage = {
        route: Landing,
        params: {},
    }

    const router = Navaid('/')
        .on('/', () => {
            currentPage = {
                route: Landing,
                params: {},
            }
        })
        .on('/enlist/:battleId?', params => {
            currentPage = {
                route: Register,
                params,
            }
        })
        .on('/login/:battleId?', params => {
            currentPage = {
                route: Login,
                params,
            }
        })
        .on('/reset-password/:resetId', params => {
            currentPage = {
                route: ResetPassword,
                params,
            }
        })
        .on('/verify-account/:verifyId', params => {
            currentPage = {
                route: VerifyAccount,
                params,
            }
        })
        .on('/warrior-profile', params => {
            currentPage = {
                route: WarriorProfile,
                params,
            }
        })
        .on('/battles', () => {
            currentPage = {
                route: Battles,
                params: {},
            }
        })
        .on('/battle/:battleId', params => {
            currentPage = {
                route: Battle,
                params,
            }
        })
        .on('/admin', () => {
            currentPage = {
                route: Admin,
                params: {},
            }
        })
        .listen()

    function logoutWarrior() {
        fetch('/api/auth/logout', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json',
            },
        })
            .then(function() {
                warrior.delete()
                router.route('/', true)
            })
            .catch(function(error) {
                notifications.danger(
                    'Error encountered attempting to logout warrior',
                )
            })
    }

    onDestroy(router.unlisten)
</script>

<style>
    :global(.nav-logo) {
        max-height: 3.75rem;
    }
    :global(.text-yellow-thunder) {
        color: #ffdd57;
    }
    :global(.bg-yellow-thunder) {
        background-color: #ffdd57;
    }
</style>

<Notifications bind:this="{notifications}" />

<nav
    class="flex items-center justify-between flex-wrap bg-white p-6"
    role="navigation"
    aria-label="main navigation">
    <div class="flex items-center flex-shrink-0 mr-6">
        <a href="/">
            <img src="/img/logo.svg" alt="Thunderdome" class="nav-logo" />
        </a>
    </div>
    {#if activeWarrior.name}
        <div class="text-right mt-4 md:mt-0">
            <span class="font-bold mr-2 text-xl">
                <WarriorIcon />
                <a href="/warrior-profile">{activeWarrior.name}</a>
            </span>
            <HollowButton color="teal" href="/battles" additionalClasses="mr-2">
                My Battles
            </HollowButton>
            {#if !activeWarrior.rank || activeWarrior.rank === 'PRIVATE'}
                <HollowButton
                    color="teal"
                    href="/enlist"
                    additionalClasses="mr-2">
                    Create Account
                </HollowButton>
                <HollowButton href="/login">Login</HollowButton>
            {:else}
                {#if activeWarrior.rank === 'GENERAL'}
                    <HollowButton
                        color="purple"
                        href="/admin"
                        additionalClasses="mr-2">
                        Admin
                    </HollowButton>
                {/if}
                <HollowButton color="red" onClick="{logoutWarrior}">
                    Logout
                </HollowButton>
            {/if}
        </div>
    {:else}
        <div class="text-right mt-4 md:mt-0">
            <HollowButton color="teal" href="/enlist" additionalClasses="mr-2">
                Create Account
            </HollowButton>
            <HollowButton href="/login">Login</HollowButton>
        </div>
    {/if}
</nav>

<svelte:component
    this="{currentPage.route}"
    {...currentPage.params}
    {notifications}
    {router} />

<footer class="p-6 text-center">
    <a
        href="https://github.com/StevenWeathers/thunderdome-planning-poker"
        class="{footerLinkClasses}">
        Thunderdome
    </a>
    by
    <a href="http://stevenweathers.com" class="{footerLinkClasses}">
        Steven Weathers
    </a>
    . The source code is licensed
    <a href="http://www.apache.org/licenses/" class="{footerLinkClasses}">
        Apache 2.0
    </a>
    .
    <br />
    Powered by
    <a href="https://svelte.dev/" class="{footerLinkClasses}">Svelte</a>
    and
    <a href="https://golang.org/" class="{footerLinkClasses}">Go</a>
</footer>

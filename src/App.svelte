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
    import { warrior } from './stores.js'

    const footerLinkClasses = 'no-underline text-teal hover:text-teal-darker'

    let notifications

    let currentPage = {
        route: Landing,
        params: {}
    }

    const router = Navaid('/')
        .on('/', () => {
            currentPage = {
                route: Landing,
                params: {}
            }
        })
        .on('/enlist/:battleId?', params => {
            currentPage = {
                route: Register,
                params
            }
        })
        .on('/login/:battleId?', params => {
            currentPage = {
                route: Login,
                params
            }
        })
        .on('/battles', () => {
            currentPage = {
                route: Battles,
                params: {}
            }
        })
        .on('/battle/:battleId', params => {
            currentPage = {
                route: Battle,
                params
            }
        })
        .listen()

    function logoutWarrior() {
        fetch('/api/auth/logout', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(function() {
                warrior.delete()
                router.route('/', true)
            }).catch(function(error) {
                notifications.danger("Error encountered attempting to logout warrior")
            })
    }

    onDestroy(router.unlisten)
</script>

<style>
    :global(.nav-logo) {
        max-height: 3.75rem;
    }
    :global(.bg-yellow-thunder) {
        background-color: #ffdd57;
    }
</style>

<Notifications bind:this={notifications} />

<nav class="flex items-center justify-between flex-wrap bg-white p-6" role="navigation" aria-label="main navigation">
    <div class="flex items-center flex-no-shrink mr-6">
        <a href="/">
            <img src="/img/logo.svg" alt="Thunderdome" class="nav-logo"/>
        </a>
    </div>
    {#if $warrior.name}
        <div class="text-right mt-4 md:mt-0">
            <span class="font-bold mr-2 text-xl"><WarriorIcon />{$warrior.name}</span>
            <HollowButton color="teal" href="/battles">
                My Battles
            </HollowButton>
            {#if !$warrior.rank || $warrior.rank === 'PRIVATE'}
                <HollowButton color="teal" href="/enlist">
                    Create Account
                </HollowButton>
                <HollowButton href="/login">Login</HollowButton>
            {:else}
                <HollowButton color="red" onClick={logoutWarrior} additionalClasses="mr-2">
                    Logout
                </HollowButton>
            {/if}
        </div>
    {:else}
        <div class="text-right mt-4 md:mt-0">
            <HollowButton color="teal" href="/enlist">Create Account</HollowButton>
            <HollowButton href="/login">Login</HollowButton>
        </div>
    {/if}
</nav>

<svelte:component this={currentPage.route} {...currentPage.params} notifications={notifications} router={router} />

<footer class="p-6 text-center">
    <a href="https://github.com/StevenWeathers/thunderdome-planning-poker" class="{footerLinkClasses}">Thunderdome</a> by <a href="http://stevenweathers.com" class="{footerLinkClasses}">Steven Weathers</a>. The source code is licensed
    <a href="http://www.apache.org/licenses/" class="{footerLinkClasses}">Apache 2.0</a>.<br />
    Powered by <a href="https://svelte.dev/" class="{footerLinkClasses}">Svelte</a> and <a href="https://golang.org/" class="{footerLinkClasses}">Go</a>
</footer>
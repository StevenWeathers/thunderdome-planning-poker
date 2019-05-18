<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'
    import Notifications from './components/Notifications.svelte'
    

    import Landing from './pages/Landing.svelte'
    import Battles from './pages/Battles.svelte'
    import Battle from './pages/Battle.svelte'
    import RegisterPage from './pages/Register.svelte'
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

    onDestroy(router.unlisten)
</script>

<style>
    img {
        max-height: 3.75rem;
    }
</style>

<Notifications bind:this={notifications} />

<nav class="flex items-center justify-between flex-wrap bg-white p-6" role="navigation" aria-label="main navigation">
    <div class="flex items-center flex-no-shrink mr-6">
        <a href="/">
            <img src="/img/logo.svg" alt="Thunderdome"/>
        </a>
    </div>
    {#if $warrior.name}
        <div class="text-right mt-4 md:mt-0">
            <a
                href="/battles"
                class="inline-block mr-4 no-underline bg-transparent hover:bg-teal text-teal-dark font-semibold hover:text-white py-2 px-2 border border-teal hover:border-transparent rounded"
            >
                My Battles
            </a>
        </div>
    {/if}
</nav>

<section>
    <div class="container mx-auto px-4 py-6 md:py-10">
    {#if !$warrior.id}
        <RegisterPage notifications={notifications} />
    {:else}
        <svelte:component this={currentPage.route} {...currentPage.params} notifications={notifications} />
    {/if}
    </div>
</section>

<footer class="p-6 text-center">
    <a href="https://github.com/StevenWeathers/thunderdome-planning-poker" class="{footerLinkClasses}">Thunderdome</a> by <a href="http://stevenweathers.com" class="{footerLinkClasses}">Steven Weathers</a>. The source code is licensed
    <a href="http://www.apache.org/licenses/" class="{footerLinkClasses}">Apache 2.0</a>.<br />
    Powered by <a href="https://svelte.dev/" class="{footerLinkClasses}">Svelte</a> and <a href="https://golang.org/" class="{footerLinkClasses}">Go</a>
</footer>
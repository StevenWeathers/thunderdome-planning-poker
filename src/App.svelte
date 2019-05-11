<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import Landing from './pages/Landing.svelte'
    import Battle from './pages/Battle.svelte'
    import RegisterPage from './pages/Register.svelte'
    import { warrior } from './stores.js'

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

<nav class="flex items-center justify-between flex-wrap bg-white p-6" role="navigation" aria-label="main navigation">
    <div class="flex items-center flex-no-shrink mr-6">
        <a href="/">
            <img src="/img/logo.svg" alt="Thunderdome"/>
        </a>
    </div>
</nav>

<section>
    <div class="container mx-auto px-4 py-10">
    {#if !$warrior.id}
        <RegisterPage />
    {:else}
        <svelte:component this={currentPage.route} {...currentPage.params} />
    {/if}
    </div>
</section>

<footer class="p-6 text-center">
    <a href="https://github.com/StevenWeathers/thunderdome-planning-poker">Thunderdome</a> by <a href="http://stevenweathers.com">Steven Weathers</a>. The source code is licensed
    <a href="http://www.apache.org/licenses/">Apache 2.0</a>.<br />
    Powered by <a href="https://svelte.dev/">Svelte</a> and <a href="https://golang.org/">Go</a>
</footer>
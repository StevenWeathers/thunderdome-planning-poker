<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import Landing from './pages/Landing.svelte'
    import Battle from './pages/Battle.svelte'
    import RegisterPage from './pages/Register.svelte'
    import { user } from './stores.js'

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

<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            <img src="/img/logo.svg" alt="Thunderdome"/>
        </a>
    </div>
</nav>

<section class="section">
    <div class="container">
    {#if !$user.id}
        <RegisterPage />
    {:else}
        <svelte:component this={currentPage.route} {...currentPage.params} />
    {/if}
    </div>
</section>

<footer class="footer">
  <div class="content has-text-centered">
    <p>
      <a href="https://github.com/StevenWeathers/thunderdome-planning-poker">Thunderdome</a> by <a href="http://stevenweathers.com">Steven Weathers</a>. The source code is licensed
      <a href="http://www.apache.org/licenses/">Apache 2.0</a>.
    </p>
  </div>
</footer>
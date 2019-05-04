<script>
    import navaid from 'navaid'
    import LandingPage from './pages/Landing.svelte'
    import BattlePage from './pages/Battle.svelte'

    // Setup router
    let router = navaid()
    
    let currentPage = {
        name: 'landing',
        params: {},
    }

    router
        .on('/', () => {
            currentPage = {
                name: 'landing',
                params: {},
            }
        })
        .on('/battle/:battleId', params => {
            currentPage = {
                name: 'battle',
                params,
            }
        })

    router.listen()
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
        {#if currentPage.name === 'landing'}
            <LandingPage />
        {:else if currentPage.name === 'battle'}
            <BattlePage {...currentPage.params} />
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
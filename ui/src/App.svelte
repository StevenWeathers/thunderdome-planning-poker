<script>
    import navaid from 'navaid'
    import LandingPage from './pages/landing.svelte'
    import BattlePage from './pages/battle.svelte'

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

<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            Thunderdome
        </a>

        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false">
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
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
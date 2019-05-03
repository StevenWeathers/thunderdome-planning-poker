<script>
    import navaid from 'navaid'
    import LandingPage from './pages/landing.svelte'
    import BattlePage from './pages/battle.svelte'

    // Setup router
    let router = navaid()
    
    let currentPage = 'landing'
    let pageParams = {}

    router
        .on('/', () => {
            pageParams = {}
            currentPage = 'landing'
        })
        .on('/battle/:battleId', params => {
            pageParams = params
            currentPage = 'battle'
        })

    router.listen()
</script>

<style>
	:global(body) {
        font-size: 18px;
    }
</style>

<nav class="navbar" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/">
            Thunderdome
        </a>
    </div>
</nav>

<section class="section">
    <div class="container">
        {#if currentPage === 'landing'}
            <LandingPage />
        {:else if currentPage === 'battle'}
            <BattlePage {...pageParams} />
        {/if}
    </div>
</section>
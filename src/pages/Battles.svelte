<script>
    import { warrior } from '../stores.js'

    export let notifications
    let battles = []

    fetch('/api/battles', {
        method: 'GET',
    })
        .then(function(response) {
            return response.json()
        })
        .then(function(bs) {
            console.log(bs)
            battles = bs
        })
        .catch(function(error) {
            notifications.danger("Error finding your battles")
        })
</script>

<h1 class="mb-4">My Battles</h1>

<div class="bg-white shadow-md rounded">
    {#each battles as battle}
        <div class="flex flex-wrap p-4 border-grey-light border-b">
            <div class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold">
                {battle.name}
            </div>
            <div class="w-full md:w-1/2 mb-4 md:mb-0 md:text-right">
                <a
                    href="/battle/{battle.id}"
                    class="inline-block bg-transparent hover:bg-green text-green-dark font-semibold hover:text-white no-underline py-2 px-2 border border-green hover:border-transparent rounded"
                >
                    Go To Battle
                </a>
            </div>
        </div>
    {/each}
</div>

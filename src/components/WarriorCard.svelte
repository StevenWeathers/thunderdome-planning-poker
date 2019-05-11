<script>
    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let points = ''
</script>

<div class="card">
    <div class="card-content">
        <div class="media">
            <div class="media-left">
                <figure class="image is-48x48">
                <img src="https://api.adorable.io/avatars/48/{warrior.id}.png" alt="Placeholder image">
                </figure>
            </div>
            <div class="media-content">
                <p class="title is-4">
                    {warrior.name}
                    {#if voted && points === ''}
                        <span class="icon is-medium has-text-success">
                            <i class="fas fa-check-square"></i>
                        </span>
                    {:else if voted && points !== ''}
                        &nbsp;<span class="has-text-success	has-text-weight-bold">{points}</span>
                    {/if}
                </p>
                {#if isLeader}<p class="subtitle is-6">Leader</p>{/if}
            </div>
        </div>
    </div>
</div>
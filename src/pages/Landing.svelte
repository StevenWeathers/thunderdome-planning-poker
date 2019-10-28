<script>
    import CheckIcon from '../components/CheckIcon.svelte'

    import { warrior } from '../stores.js'
</script>

<section class="bg-yellow-thunder text-grey-darkest">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="flex flex-wrap items-center -mx-4">
            <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
                <h1 class="mb-2 lg:mb-4">Thunderdome is an Open Source Agile Planning Poker app with a fun theme</h1>

                <ul class="text-lg lg:text-2xl list-reset mb-4">
                    <li class="py-2 lg:py-4"><CheckIcon /> Built on cool tech like Svelte, Go, and WebSockets</li>
                    <li class="py-2 lg:py-4"><CheckIcon /> Simple <span class="line-through">User</span> Warrior Experience</li>
                    <li class="py-2 lg:py-4">
                        <CheckIcon /> Works on all <span class="font-semibold">modern browsers*</span>, including mobile<br />
                        <span class="text-sm">*sorry, not sorry IE</span>
                    </li>
                </ul>

                <div class="text-center mb-4 md:mb-0">
                    <a
                        class="w-full md:w-auto inline-block no-underline bg-transparent hover:bg-grey-darkest text-grey-darkest font-semibold hover:text-white py-4 px-4 border border-grey-darkest hover:border-transparent rounded"
                        href={$warrior.id ? "/battles" : "/enlist"}
                    >
                    Create a Battle</a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <img class="w-full" src="https://user-images.githubusercontent.com/846933/58061038-58d62d00-7b42-11e9-9679-ebd297a51c05.png" alt="Preview of Thunderdome, Agile Planning Poker UI Experience"/>
            </div>
        </div>
    </div>
</section>

<section>
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="text-center mb-6">
            <h2>Features</h2>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/3"><CheckIcon /> No limit on number of Warriors in a Battle</div>
            <div class="w-1/3"><CheckIcon /> Secure voting, keeps your Warriors from cheating</div>
            <div class="w-1/3"><CheckIcon /> Voting Results that are clear and helpful</div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/3"><CheckIcon /> Ability to promote any Warrior to Battle Leader</div>
            <div class="w-1/3"><CheckIcon /> Customizable pointing scales</div>
            <div class="w-1/3">
                <CheckIcon /> Open Source, feel free to <a href="https://github.com/StevenWeathers/thunderdome-planning-poker" class="no-underline text-teal hover:text-teal-darker">contribute and suggest features</a>
            </div>
        </div>
        <div class="flex text-center">
            <div class="w-full"><CheckIcon /> No fees &amp; Ad free</div>
        </div>
    </div>
</section>
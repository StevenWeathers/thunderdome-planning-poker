<script>
    import CheckIcon from '../components/CheckIcon.svelte'
</script>

<section class="bg-yellow-thunder text-grey-darkest">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="flex flex-wrap -mx-4">
            <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
                <h1 class="mb-2 lg:mb-4">Agile Planning Poker, with a fun theme</h1>
                <h2 class="mb-8 lg:mb-16">oh and its Open Source...</h2>

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
                        href="/battles"
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

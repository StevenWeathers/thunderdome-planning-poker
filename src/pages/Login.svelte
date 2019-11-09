<script>
    import PageLayout from '../components/PageLayout.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import { warrior } from '../stores.js'

    export let router
    export let notifications
    export let battleId

    let warriorEmail = ''
    let warriorPassword = ''

    $: targetPage = battleId ? `/battle/${battleId}` : '/battles'

    function authWarrior(e) {
        e.preventDefault()
        const body = {
            warriorEmail,
            warriorPassword
        }
        
        
        fetch('/api/auth', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(body)
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(newWarrior) {
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name,
                    email: newWarrior.email,
                    rank: newWarrior.rank
                })
                
                router.route(targetPage, true)
            }).catch(function(error) {
                notifications.danger("Error encountered attempting to authenticate warrior")
            })
    }

    $: loginDisabled = warriorEmail === '' || warriorPassword === ''
</script>

<PageLayout>
    <div class="text-center px-4 mb-4">
        <h1 class="text-4xl font-bold">Login</h1>
    </div>

    <div class="flex justify-center">
        <div class="w-1/3">
            <form on:submit={authWarrior} class="bg-white shadow-lg rounded p-6 mb-4" name="authWarrior">
                <div class="mb-4">
                    <label class="block text-gray-700 text-sm font-bold mb-2" for="yourEmail">Email</label>
                    <input
                        bind:value={warriorEmail}
                        placeholder="Enter your email"
                        class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                        id="yourEmail"
                        name="yourEmail"
                        type="email"
                        required
                    />
                </div>

                <div class="mb-4">
                    <label class="block text-gray-700 text-sm font-bold mb-2" for="yourPassword">Password</label>
                    <input
                        bind:value={warriorPassword}
                        placeholder="Enter your password"
                        class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                        id="yourPassword"
                        name="yourPassword"
                        type="password"
                        required
                    />
                </div>

                <div>
                    <div class="text-right">
                        <SolidButton type="submit" disabled={loginDisabled}>Login</SolidButton>
                    </div>
                </div>
            </form>
        </div>
    </div>
</PageLayout>
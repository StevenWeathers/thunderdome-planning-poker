<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import CheckIcon from '../../components/icons/CheckIcon.svelte'
    import { warrior } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'
    import { validateUserIsAdmin } from '../../validationUtils'

    export let xfetch
    export let router
    export let notifications
    export let eventTag
    export let userId

    let user = {
        id: '',
        name: '',
        email: '',
        rank: '',
        avatar: '',
        verified: false,
        notificationsEnabled: true,
        country: '',
        locale: '',
        company: '',
        jobTitle: '',
        createdDate: '',
        updatedDate: '',
        lastActive: '',
    }

    function getUser() {
        xfetch(`/api/users/${userId}`)
            .then(res => res.json())
            .then(function (result) {
                user = result.data
            })
            .catch(function () {
                notifications.danger('Error getting user')
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(appRoutes.login)
        }
        if (!validateUserIsAdmin($warrior)) {
            router.route(appRoutes.landing)
        }

        getUser()
    })
</script>

<svelte:head>
    <title>{$_('users')} {$_('pages.admin.title')} | {$_('appName')}</title>
</svelte:head>

<AdminPageLayout activePage="users">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">{$_('users')}</h1>
    </div>

    <div class="w-full">
        <div class="p-4 md:p-6 bg-white shadow-lg rounded">
            <h3 class="text-2xl md:text-3xl font-bold mb-4 text-center">
                {user.name}
            </h3>
            <table class="table-fixed w-full mb-4">
                <thead>
                    <tr>
                        <th class="flex-1 p-2">Email</th>
                        <th class="flex-1 p-2">Type</th>
                        <th class="flex-1 p-2">Created</th>
                        <th class="flex-1 p-2">Updated</th>
                        <th class="flex-1 p-2">Last Active</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td class="border p-2"
                            >{user.email}
                            {#if user.verified}
                                &nbsp;
                                <span
                                    class="text-green-600"
                                    title="{$_(
                                        'pages.admin.registeredWarriors.verified',
                                    )}"
                                >
                                    <CheckIcon />
                                </span>
                            {/if}
                        </td>
                        <td class="border p-2">{user.rank}</td>
                        <td class="border p-2"
                            >{new Date(user.createdDate).toLocaleString()}</td
                        >
                        <td class="border p-2"
                            >{new Date(user.updatedDate).toLocaleString()}</td
                        >
                        <td class="border p-2"
                            >{new Date(user.lastActive).toLocaleString()}</td
                        >
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
</AdminPageLayout>

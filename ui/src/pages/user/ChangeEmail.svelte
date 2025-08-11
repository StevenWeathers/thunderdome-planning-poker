<script lang="ts">
    import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
    import TextInput from '../../components/forms/TextInput.svelte';
    import SolidButton from '../../components/global/SolidButton.svelte';

    import type { NotificationService } from '../../types/notifications';

  interface Props {
    xfetch: any;
    verifyId: any;
    notifications: NotificationService;
    router: any;
    changeId: string;
  }

  let { xfetch, notifications, changeId, router }: Props = $props();

  let newEmail = $state('');
  let formDisabled = $state(true);
  let emailChanged = $state(false);

    $effect(() => {
        formDisabled = newEmail === '';
    });

    function changeUserEmail(event: Event) {
        event.preventDefault();

        if (formDisabled) return;

        xfetch(`/api/users/${$user.id}/email-change/${changeId}`, {
            method: 'POST',
            body: {
                email: newEmail,
            },
        }).then(res => res.json())
        .then(resp => {
                user.update({
                    ...$user,
                    email: resp.data.email,
                });
                emailChanged = true;
            })
            .catch(() => {
                notifications.danger($LL.errorChangingEmail());
            });
    }

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.login);
            return;
        }
    });
</script>

<svelte:head>
  <title>{$LL.changeEmail()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex justify-center">
    <div class="w-full md:w-1/2 xl:w-1/3 py-4">
        {#if emailChanged}
            <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4 text-center" role="alert">
                <div class="font-bold mb-2 lg:mb-4">{$LL.emailChanged()}</div>
                <div>{$LL.newEmailToLogin()}</div>
            </div>
        {:else}
            <form
                onsubmit={changeUserEmail}
                class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                name="changeUserEmail"
            >
                <div
                class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                            md:leading-tight text-center dark:text-white"
                >
                {$LL.changeEmail()}
                </div>
                <div class="mb-4 font-semibold text-gray-700 dark:text-gray-400">
                    Current email: <strong>{$user.email}</strong>
                </div>

                <div class="mb-4">
                <label
                    class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                    for="yourPassword1"
                >
                    {$LL.newEmail()}
                </label>
                <TextInput
                    bind:value="{newEmail}"
                    placeholder={$LL.enterYourNewEmail()}
                    id="newEmail"
                    name="email"
                    type="email"
                    required
                />
                </div>

                <div class="text-right">
                <SolidButton type="submit" disabled={formDisabled}>
                    {$LL.save()}
                </SolidButton>
                </div>
            </form>
        {/if}
    </div>
  </div>
</PageLayout>

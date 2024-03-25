<script lang="ts">
  import PageLayout from '../../components/global/PageLayout.svelte';
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';

  export let router;

  const { RepoURL } = AppConfig;

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }

    user.update({
      id: $user.id,
      name: $user.name,
      email: $user.email,
      rank: $user.rank,
      avatar: $user.avatar,
      verified: $user.verified,
      notificationsEnabled: $user.notificationsEnabled,
      locale: $user.locale,
      theme: $user.theme,
      subscribed: true,
    });
  });
</script>

<PageLayout>
  <div class="flex justify-center">
    <div class="w-full py-12">
      <div class="text-center text-lime-300">
        <h1 class="uppercase font-rajdhani text-5xl mb-4">
          Thank you for subscribing
        </h1>
        <div class="text-gray-700 dark:text-gray-300 text-2xl">
          <p class="mb-2">
            I greatly appreciate the support and hope to provide subscriber
            premium features soon.
          </p>
          <p>
            Got a suggested feature? <a
              class="underline text-lime-300"
              target="_blank"
              href="{RepoURL}/issues/new">Create a new github issue</a
            >.
          </p>
        </div>
      </div>
    </div>
  </div>
</PageLayout>

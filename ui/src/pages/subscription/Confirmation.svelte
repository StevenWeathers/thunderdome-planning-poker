<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import { AppConfig, appRoutes, PathPrefix } from '../../config';
  import type { SessionUser } from '../../types/user';

  interface Props {
    router: any;
  }

  let { router }: Props = $props();

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
    } as SessionUser);
  });
</script>

<section class="relative bg-indigo-600 py-12">
  <div class="relative mx-auto max-w-2xl px-6 text-center lg:max-w-7xl lg:px-8">
    <h1 class="text-5xl font-bold tracking-tight text-white sm:text-6xl">Thank you for subscribing!</h1>
    <p class="mt-4 text-xl text-indigo-100">
      You should receive an email shortly, once your subscription is activated.
    </p>
  </div>
</section>

<section class="bg-white dark:bg-transparent text-gray-700 dark:text-gray-300">
  <div class="container px-5 pb-24 pt-12 mx-auto max-w-2xl lg:max-w-4xl">
    <div class="mb-8">
      <h2 class="mb-4 font-bold text-2xl md:text-3xl lg:text-4xl text-green-500 dark:text-lime-400">
        Associate Team or Organization to Subscription
      </h2>
      <p class="mb-2 text-lg">
        If you purchased a
        <span class="font-bold">non Individual</span>
        subscription follow the below steps to associate it to the respective
        <span class="font-bold">Team</span>
        or <span class="font-bold">Organization</span>.
      </p>
    </div>
    <div class="mb-4">
      <h3 class="text-lg text-2xl font-bold mb-4 text-blue-500 dark:text-sky-400">
        1. Navigate to your profile by pressing the Ninja User in the top navigation bar
      </h3>
      <img
        src="{PathPrefix}/img/header_user_profile_icon_preview.png"
        alt="Preview of the ninja user icon in the top navigation bar"
        class="max-w-full md:max-w-sm border border-2px solid border-dashed border-blue-500 dark:border-sky-400 p-2"
      />
    </div>
    <div class="mb-4">
      <h3 class="text-lg text-2xl font-bold mb-4 text-blue-500 dark:text-sky-400">
        2. Look for the <span class="underline">Active Subscriptions</span>
        section and Press the <span class="underline">Associate to</span> button for the respective subscription type
      </h3>
      <img
        src="{PathPrefix}/img/active_subs_preview.png"
        alt="Preview of the active subscriptions section on user profile page"
        class="max-w-full lg:max-w-2xl border border-2px solid border-dashed border-blue-500 dark:border-sky-400 p-2"
      />
    </div>
    <div class="mb-4">
      <h3 class="text-lg text-2xl font-bold mb-4 text-blue-500 dark:text-sky-400">
        3. Select the Team or Organization to associate to the subscription and press save
      </h3>
      <img
        src="{PathPrefix}/img/sub_associate_team_preview.png"
        alt="Preview of the associate team to subscription modal"
        class="max-w-full lg:max-w-xl border border-2px solid border-dashed border-blue-500 dark:border-sky-400 p-2"
      />
    </div>
    <div>
      <h3 class="text-lg md:text-xl lg:text-2xl font-bold mb-4 text-green-500 dark:text-lime-400">
        Now the associated Team or Organization members can access premium features after they re-login.
      </h3>
    </div>
  </div>
</section>

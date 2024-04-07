<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../global/TextInput.svelte';
  import LockIcon from '../icons/LockIcon.svelte';

  export let registerLink = '';
  export let targetPage = appRoutes.landing;
  export let router;
  export let xfetch = async (url, ...options) => {};
  export let eventTag = () => {};
  export let notifications = () => {};

  declare global {
    interface Window {
      setTheme: Function;
    }
  }

  const { AllowRegistration, LdapEnabled } = AppConfig;
  const authEndpoint = LdapEnabled ? '/api/auth/ldap' : '/api/auth';

  let email = '';
  let password = '';
  let forgotPassword = false;
  let resetEmail = '';

  let mfaRequired = false;
  let mfaUser = null;
  let mfaSessionId = null;
  let mfaToken = '';

  function toggleForgotPassword() {
    forgotPassword = !forgotPassword;
    eventTag(
      'forgot_password_toggle',
      'engagement',
      `forgot: ${forgotPassword}`,
    );
  }

  function handleLoginSubmit(e) {
    e.preventDefault();

    const body = {
      email,
      password,
    };

    xfetch(authEndpoint, { body, skip401Redirect: true })
      .then((res: any) => res.json())
      .then(function (result) {
        const u = result.data.user;
        const newUser = {
          id: u.id,
          name: u.name,
          email: u.email,
          rank: u.rank,
          locale: u.locale,
          notificationsEnabled: u.notificationsEnabled,
          subscribed: result.data.subscribed,
        };
        if (result.data.mfaRequired) {
          mfaRequired = true;
          mfaUser = newUser;
          mfaSessionId = result.data.sessionId;
        } else {
          user.create(newUser);
          if (u.theme !== 'auto') {
            localStorage.setItem('theme', u.theme);
            window.setTheme();
          }
          eventTag('login', 'engagement', 'success', () => {
            // setupI18n({
            //     withLocale: newUser.locale,
            // })
            router.route(targetPage, true);
          });
        }
      })
      .catch(function () {
        notifications.danger(
          $LL.authError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
        eventTag('login', 'engagement', 'failure');
      });
  }

  function authMfa(e) {
    e.preventDefault();
    const body = {
      passcode: mfaToken,
      sessionId: mfaSessionId,
    };

    xfetch('/api/auth/mfa', { body, skip401Redirect: true })
      .then(res => res.json())
      .then(function () {
        user.create(mfaUser);
        eventTag('login_mfa', 'engagement', 'success', () => {
          // setupI18n({
          //     withLocale: mfaUser.locale,
          // })
          router.route(targetPage, true);
        });
      })
      .catch(function () {
        notifications.danger($LL.mfaAuthError());
        eventTag('login_mfa', 'engagement', 'failure');
      });
  }

  function sendPasswordReset(e) {
    e.preventDefault();
    const body = {
      email: resetEmail,
    };

    xfetch('/api/auth/forgot-password', { body })
      .then(function () {
        notifications.success(
          $LL.sendResetPasswordSuccess({
            email: resetEmail,
          }),
          2000,
        );
        forgotPassword = !forgotPassword;
        eventTag('forgot_password', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger($LL.sendResetPasswordError());
        eventTag('forgot_password', 'engagement', 'failure');
      });
  }
</script>

{#if !forgotPassword && !mfaRequired}
  <form
    on:submit="{handleLoginSubmit}"
    class="space-y-6"
    name="login"
    id="login"
  >
    <div>
      <label
        for="email"
        class="block text-sm font-medium text-gray-700 dark:text-white"
        >{$LL.email()}</label
      >
      <div class="mt-1">
        <input
          id="email"
          type="text"
          data-testid="username"
          placeholder="{$LL.enterYourEmail()}"
          required
          class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-300 dark:focus:border-indigo-400 dark:focus:ring-indigo-400 sm:text-sm"
          bind:value="{email}"
        />
      </div>
    </div>
    <div>
      <label
        for="password"
        class="block text-sm font-medium text-gray-700 dark:text-white"
        >{$LL.password()}</label
      >
      <div class="mt-1">
        <input
          id="password"
          name="password"
          type="password"
          data-testid="password"
          autocomplete="current-password"
          required
          class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-300 dark:focus:border-indigo-400 dark:focus:ring-indigo-400 sm:text-sm"
          placeholder="{$LL.yourPasswordPlaceholder()}"
          bind:value="{password}"
        />
      </div>
    </div>
    <div class="flex items-center justify-between">
      <div class="flex items-center">
        <!--                        <input id="remember_me" name="remember_me" type="checkbox"-->
        <!--                               class="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500 dark:text-white dark:border-gray-600 dark:focus:ring-indigo-400 disabled:cursor-wait disabled:opacity-50">-->
        <!--                        <label for="remember_me" class="ml-2 block text-sm text-gray-900 dark:text-white">Remember-->
        <!--                            me</label>-->
      </div>
      <div class="text-sm">
        {#if !LdapEnabled}
          <a
            class="font-medium text-indigo-400 hover:text-indigo-500 cursor-pointer"
            on:click="{toggleForgotPassword}"
          >
            {$LL.forgotPassword()}
          </a>
        {/if}
      </div>
    </div>
    <div>
      <button
        data-testid="login"
        type="submit"
        class="group relative flex w-full justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:bg-indigo-700 dark:border-transparent dark:hover:bg-indigo-600 dark:focus:ring-indigo-400 dark:focus:ring-offset-2 disabled:cursor-wait disabled:opacity-50"
      >
        <span class="absolute inset-y-0 left-0 flex items-center pl-3">
          <LockIcon
            class="h-5 w-5 text-indigo-500 group-hover:text-indigo-400"
          />
        </span>
        {$LL.login()}
      </button>
    </div>
  </form>
  <!--      <div class="mt-6">-->
  <!--        <div class="relative">-->
  <!--          <div class="absolute inset-0 flex items-center">-->
  <!--            <div class="w-full border-t border-gray-300"></div>-->
  <!--          </div>-->
  <!--          <div class="relative flex justify-center text-sm">-->
  <!--            <span-->
  <!--              class="bg-white dark:bg-gray-700 px-2 text-gray-500 dark:text-white"-->
  <!--              >Or continue with</span-->
  <!--            >-->
  <!--          </div>-->
  <!--        </div>-->
  <!--        <div class="mt-6 grid grid-cols-2 gap-3">-->
  <!--          <button-->
  <!--            class="inline-flex w-full items-center justify-center rounded-md border border-gray-300 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-500 dark:text-white shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 disabled:cursor-wait disabled:opacity-50"-->
  <!--          >-->
  <!--            <span class="sr-only">Sign in with Google</span>-->
  <!--            <img-->
  <!--              class="w-6 h-6"-->
  <!--              src="https://www.svgrepo.com/show/475656/google-color.svg"-->
  <!--              loading="lazy"-->
  <!--              alt="google logo"-->
  <!--            />-->
  <!--            &lt;!&ndash;                        <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">&ndash;&gt;-->
  <!--            &lt;!&ndash;                            <clipPath id="p.0">&ndash;&gt;-->
  <!--            &lt;!&ndash;                                <path d="m0 0l20.0 0l0 20.0l-20.0 0l0 -20.0z" clip-rule="nonzero"></path>&ndash;&gt;-->
  <!--            &lt;!&ndash;                            </clipPath>&ndash;&gt;-->
  <!--            &lt;!&ndash;                            <g clip-path="url(#p.0)">&ndash;&gt;-->
  <!--            &lt;!&ndash;                                <path fill="currentColor" fill-opacity="0.0" d="m0 0l20.0 0l0 20.0l-20.0 0z"&ndash;&gt;-->
  <!--            &lt;!&ndash;                                      fill-rule="evenodd"></path>&ndash;&gt;-->
  <!--            &lt;!&ndash;                                <path fill="currentColor"&ndash;&gt;-->
  <!--            &lt;!&ndash;                                      d="m19.850197 8.270351c0.8574047 4.880001 -1.987587 9.65214 -6.6881847 11.218641c-4.700598 1.5665016 -9.83958 -0.5449295 -12.08104 -4.963685c-2.2414603 -4.4187555 -0.909603 -9.81259 3.1310139 -12.6801605c4.040616 -2.867571 9.571754 -2.3443127 13.002944 1.2301085l-2.8127813 2.7000687l0 0c-2.0935059 -2.1808972 -5.468274 -2.500158 -7.933616 -0.75053835c-2.4653416 1.74962 -3.277961 5.040613 -1.9103565 7.7366734c1.3676047 2.6960592 4.5031037 3.9843292 7.3711267 3.0285425c2.868022 -0.95578575 4.6038647 -3.8674583 4.0807285 -6.844941z"&ndash;&gt;-->
  <!--            &lt;!&ndash;                                      fill-rule="evenodd"></path>&ndash;&gt;-->
  <!--            &lt;!&ndash;                                <path fill="currentColor" d="m10.000263 8.268785l9.847767 0l0 3.496233l-9.847767 0z"&ndash;&gt;-->
  <!--            &lt;!&ndash;                                      fill-rule="evenodd"></path>&ndash;&gt;-->
  <!--            &lt;!&ndash;                            </g>&ndash;&gt;-->
  <!--            &lt;!&ndash;                        </svg>&ndash;&gt;-->
  <!--          </button>-->
  <!--          <button-->
  <!--            class="inline-flex w-full justify-center rounded-md border border-gray-300 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-500 dark:text-white shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 disabled:cursor-wait disabled:opacity-50"-->
  <!--          >-->
  <!--            <span class="sr-only">Sign in with GitHub</span>-->
  <!--            <svg-->
  <!--              class="h-6 w-6"-->
  <!--              fill="currentColor"-->
  <!--              xmlns="http://www.w3.org/2000/svg"-->
  <!--              width="800px"-->
  <!--              height="800px"-->
  <!--              viewBox="0 0 32 32"-->
  <!--              version="1.1"-->
  <!--            >-->
  <!--              <path-->
  <!--                d="M16 1.375c-8.282 0-14.996 6.714-14.996 14.996 0 6.585 4.245 12.18 10.148 14.195l0.106 0.031c0.750 0.141 1.025-0.322 1.025-0.721 0-0.356-0.012-1.3-0.019-2.549-4.171 0.905-5.051-2.012-5.051-2.012-0.288-0.925-0.878-1.685-1.653-2.184l-0.016-0.009c-1.358-0.93 0.105-0.911 0.105-0.911 0.987 0.139 1.814 0.718 2.289 1.53l0.008 0.015c0.554 0.995 1.6 1.657 2.801 1.657 0.576 0 1.116-0.152 1.582-0.419l-0.016 0.008c0.072-0.791 0.421-1.489 0.949-2.005l0.001-0.001c-3.33-0.375-6.831-1.665-6.831-7.41-0-0.027-0.001-0.058-0.001-0.089 0-1.521 0.587-2.905 1.547-3.938l-0.003 0.004c-0.203-0.542-0.321-1.168-0.321-1.821 0-0.777 0.166-1.516 0.465-2.182l-0.014 0.034s1.256-0.402 4.124 1.537c1.124-0.321 2.415-0.506 3.749-0.506s2.625 0.185 3.849 0.53l-0.1-0.024c2.849-1.939 4.105-1.537 4.105-1.537 0.285 0.642 0.451 1.39 0.451 2.177 0 0.642-0.11 1.258-0.313 1.83l0.012-0.038c0.953 1.032 1.538 2.416 1.538 3.937 0 0.031-0 0.061-0.001 0.091l0-0.005c0 5.761-3.505 7.029-6.842 7.398 0.632 0.647 1.022 1.532 1.022 2.509 0 0.093-0.004 0.186-0.011 0.278l0.001-0.012c0 2.007-0.019 3.619-0.019 4.106 0 0.394 0.262 0.862 1.031 0.712 6.028-2.029 10.292-7.629 10.292-14.226 0-8.272-6.706-14.977-14.977-14.977-0.006 0-0.013 0-0.019 0h0.001z"-->
  <!--              ></path>-->
  <!--            </svg>-->
  <!--          </button>-->
  <!--        </div>-->
  <!--      </div>-->
  {#if registerLink !== ''}
    <div class="m-auto mt-6 w-fit md:mt-8">
      <span class="m-auto dark:text-gray-400"
        >{$LL.createAccountTagline()}
        <a
          class="font-semibold text-indigo-600 dark:text-indigo-100"
          href="{registerLink}">{$LL.createAccount()}</a
        >
      </span>
    </div>
  {/if}
{/if}

{#if forgotPassword}
  <form on:submit="{sendPasswordReset}" class="space-y-6" name="resetPassword">
    <div class="mb-4">
      <h2
        class="mb-1 text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white"
      >
        {$LL.forgotPassword()}
      </h2>
      <p class="font-light text-gray-600 dark:text-gray-300">
        {$LL.forgotPasswordSubtext()}
      </p>
    </div>
    <div>
      <label
        for="yourResetEmail"
        class="block text-sm font-medium text-gray-700 dark:text-white"
        >{$LL.email()}</label
      >
      <div class="mt-1">
        <input
          data-testid="resetemail"
          bind:value="{resetEmail}"
          placeholder="{$LL.enterYourEmail()}"
          id="yourResetEmail"
          name="yourResetEmail"
          type="email"
          required
          class="block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 placeholder-gray-400 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-300 dark:focus:border-indigo-400 dark:focus:ring-indigo-400 sm:text-sm"
        />
      </div>
    </div>

    <div class="text-right">
      <button
        type="button"
        class="inline-block align-baseline font-bold text-sm
                            text-indigo-400 hover:text-indigo-500 me-4"
        on:click="{toggleForgotPassword}"
      >
        {$LL.returnToLogin()}
      </button>
      <button
        type="submit"
        class="rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:bg-indigo-700 dark:border-transparent dark:hover:bg-indigo-600 dark:focus:ring-indigo-400 dark:focus:ring-offset-2 disabled:cursor-wait disabled:opacity-50"
      >
        {$LL.resetPassword()}
      </button>
    </div>
  </form>
{/if}

{#if mfaRequired}
  <form on:submit="{authMfa}" class="space-y-6" name="authMfa">
    <div
      class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight text-center dark:text-white"
    >
      {$LL.login()}
    </div>
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="mfaToken"
      >
        {$LL.mfaTokenLabel()}
      </label>
      <TextInput
        bind:value="{mfaToken}"
        placeholder="{$LL.mfaTokenPlaceholder()}"
        id="mfaToken"
        name="mfaToken"
        required
      />
    </div>

    <div class="text-right">
      <SolidButton type="submit">
        {$LL.login()}
      </SolidButton>
    </div>
  </form>
{/if}

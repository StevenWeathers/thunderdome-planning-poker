<script lang="ts">
  import { AppConfig, appRoutes, PathPrefix } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { Lock, Mail, Shield } from 'lucide-svelte';
  import PasswordInput from '../forms/PasswordInput.svelte';
  import Google from '../icons/Google.svelte';
  import TextInput from '../forms/TextInput.svelte';

  interface Props {
    registerLink?: string;
    targetPage?: any;
    router: any;
    xfetch?: any;
    notifications?: any;
  }

  let {
    registerLink = '',
    targetPage = appRoutes.landing,
    router,
    xfetch = async (url, ...options) => {},
    notifications = () => {}
  }: Props = $props();

  declare global {
    interface Window {
      setTheme: Function;
    }
  }

  const {
    LdapEnabled,
    GoogleAuthEnabled,
    HeaderAuthEnabled,
    OIDCAuthEnabled,
    OIDCProviderName,
  } = AppConfig;
  const authEndpoint = LdapEnabled ? '/api/auth/ldap' : '/api/auth';

  let email = $state('');
  let password = $state('');
  let forgotPassword = $state(false);
  let resetEmail = $state('');

  let mfaRequired = $state(false);
  let mfaUser = null;
  let mfaSessionId = null;
  let mfaToken = $state('');

  function googleLogin() {
    window.location = `${PathPrefix}/oauth/google/login`;
  }

  function oidcLogin() {
    window.location = `${PathPrefix}/oauth/${OIDCProviderName.toLowerCase()}/login`;
  }

  function toggleForgotPassword() {
    forgotPassword = !forgotPassword;
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
          } else {
            localStorage.removeItem('theme');
          }
          window.setTheme();
          router.route(targetPage, true);
        }
      })
      .catch(function () {
        notifications.danger($LL.authError());
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
        router.route(targetPage, true);
      })
      .catch(function () {
        notifications.danger($LL.mfaAuthError());
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
      })
      .catch(function () {
        notifications.danger($LL.sendResetPasswordError());
      });
  }
</script>

{#if OIDCAuthEnabled}
  <button
    onclick={oidcLogin}
    data-testid="login"
    class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
  >
    <span class="flex items-center pe-3">
      <Lock
        class="h-5 w-5 text-purple-300 group-hover:text-purple-200"
        aria-hidden="true"
      />
    </span>
    {$LL.loginWithSSO({ provider: OIDCProviderName })}
  </button>
{/if}

{#if !OIDCAuthEnabled && !forgotPassword && !mfaRequired}
  <form
    onsubmit={handleLoginSubmit}
    class="space-y-6"
    name="login"
    id="login"
  >
    <TextInput
      id="email"
      data-testid="username"
      placeholder={$LL.enterYourEmail()}
      required
      bind:value="{email}"
      icon={Mail}
      autocomplete="email"
    />

    <PasswordInput
      bind:value="{password}"
      placeholder={$LL.yourPasswordPlaceholder()}
      id="password"
      name="password"
      data-testid="password"
      required
    />

    <div class="flex items-center justify-between text-sm">
      <div class="flex items-center">
        <!--                      <input-->
        <!--                              id="remember_me"-->
        <!--                              name="remember_me"-->
        <!--                              type="checkbox"-->
        <!--                              class="h-4 w-4 rounded focus:ring-2 focus:ring-offset-2 transition-all duration-300 bg-purple-100 dark:bg-purple-900 border-purple-300 dark:border-purple-700 focus:ring-purple-500 dark:focus:ring-purple-400"-->
        <!--                      />-->
        <!--                      <label for="remember_me" class="ml-2 text-gray-700 dark:text-gray-300">-->
        <!--                          {$LL.rememberMe()}-->
        <!--                      </label>-->
      </div>
      {#if !LdapEnabled}
        <a
          class="font-medium text-purple-600 dark:text-purple-400 hover:text-purple-500 dark:hover:text-purple-300 transition-all duration-300 cursor-pointer"
          onclick={(e) => { 
            e.preventDefault();
            toggleForgotPassword();
          }}
          href="/login"
        >
          {$LL.forgotPassword()}
        </a>
      {/if}
    </div>

    <button
      data-testid="login"
      type="submit"
      class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span class="absolute left-0 inset-y-0 flex items-center pl-3">
        <Lock
          class="h-5 w-5 text-purple-300 group-hover:text-purple-200"
          aria-hidden="true"
        />
      </span>
      {$LL.login()}
    </button>
  </form>
  {#if GoogleAuthEnabled && !HeaderAuthEnabled && !LdapEnabled}
    <div class="w-full space-y-4 mt-4">
      <div class="flex items-center space-x-2">
        <div class="flex-grow h-px bg-gray-300"></div>
        <span
          class="uppercase text-sm text-gray-500 dark:text-gray-400 font-medium"
          >Or continue with</span
        >
        <div class="flex-grow h-px bg-gray-300"></div>
      </div>
      <button
        onclick={googleLogin}
        class="inline-flex w-full items-center justify-center rounded-md border border-gray-600 dark:border-gray-400 hover:border-indigo-600 dark:hover:border-purple-400 px-4 py-2 shadow-sm disabled:cursor-wait disabled:opacity-50"
      >
        <span class="sr-only">Sign in with Google</span>
        <Google class="w-8 h-8" />
      </button>
      <!--            <button-->
      <!--                    class="inline-flex w-full justify-center rounded-md border border-gray-300 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-500 dark:text-white shadow-sm hover:bg-gray-50 dark:hover:bg-gray-600 disabled:cursor-wait disabled:opacity-50"-->
      <!--            >-->
      <!--                <span class="sr-only">Sign in with GitHub</span>-->
      <!--                <svg-->
      <!--                        class="h-6 w-6"-->
      <!--                        fill="currentColor"-->
      <!--                        xmlns="http://www.w3.org/2000/svg"-->
      <!--                        width="800px"-->
      <!--                        height="800px"-->
      <!--                        viewBox="0 0 32 32"-->
      <!--                        version="1.1"-->
      <!--                >-->
      <!--                    <path-->
      <!--                            d="M16 1.375c-8.282 0-14.996 6.714-14.996 14.996 0 6.585 4.245 12.18 10.148 14.195l0.106 0.031c0.750 0.141 1.025-0.322 1.025-0.721 0-0.356-0.012-1.3-0.019-2.549-4.171 0.905-5.051-2.012-5.051-2.012-0.288-0.925-0.878-1.685-1.653-2.184l-0.016-0.009c-1.358-0.93 0.105-0.911 0.105-0.911 0.987 0.139 1.814 0.718 2.289 1.53l0.008 0.015c0.554 0.995 1.6 1.657 2.801 1.657 0.576 0 1.116-0.152 1.582-0.419l-0.016 0.008c0.072-0.791 0.421-1.489 0.949-2.005l0.001-0.001c-3.33-0.375-6.831-1.665-6.831-7.41-0-0.027-0.001-0.058-0.001-0.089 0-1.521 0.587-2.905 1.547-3.938l-0.003 0.004c-0.203-0.542-0.321-1.168-0.321-1.821 0-0.777 0.166-1.516 0.465-2.182l-0.014 0.034s1.256-0.402 4.124 1.537c1.124-0.321 2.415-0.506 3.749-0.506s2.625 0.185 3.849 0.53l-0.1-0.024c2.849-1.939 4.105-1.537 4.105-1.537 0.285 0.642 0.451 1.39 0.451 2.177 0 0.642-0.11 1.258-0.313 1.83l0.012-0.038c0.953 1.032 1.538 2.416 1.538 3.937 0 0.031-0 0.061-0.001 0.091l0-0.005c0 5.761-3.505 7.029-6.842 7.398 0.632 0.647 1.022 1.532 1.022 2.509 0 0.093-0.004 0.186-0.011 0.278l0.001-0.012c0 2.007-0.019 3.619-0.019 4.106 0 0.394 0.262 0.862 1.031 0.712 6.028-2.029 10.292-7.629 10.292-14.226 0-8.272-6.706-14.977-14.977-14.977-0.006 0-0.013 0-0.019 0h0.001z"-->
      <!--                    ></path>-->
      <!--                </svg>-->
      <!--            </button>-->
    </div>
  {/if}
  {#if registerLink !== ''}
    <div class="m-auto mt-6 w-fit md:mt-8">
      <span class="m-auto dark:text-gray-400"
        >{$LL.createAccountTagline()}
        <a
          class="font-semibold text-indigo-600 dark:text-indigo-100"
          href={registerLink}>{$LL.createAccount()}</a
        >
      </span>
    </div>
  {/if}
{/if}

{#if forgotPassword}
  <form
    onsubmit={sendPasswordReset}
    class="space-y-6 max-w-md"
    name="resetPassword"
  >
    <div class="mb-4">
      <h2
        class="font-rajdhani uppercase text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-pink-600 dark:from-purple-500 dark:to-pink-500 mb-2"
      >
        {$LL.forgotPassword()}
      </h2>
      <p class="font-light text-gray-600 dark:text-gray-300">
        {$LL.forgotPasswordSubtext()}
      </p>
    </div>

    <TextInput
      data-testid="resetemail"
      bind:value={resetEmail}
      placeholder={$LL.enterYourEmail()}
      id="yourResetEmail"
      name="yourResetEmail"
      type="email"
      required
      icon={Mail}
    />

    <div class="flex justify-between items-center">
      <button
        type="button"
        class="font-medium text-purple-600 dark:text-purple-400 hover:text-purple-500 dark:hover:text-purple-300 transition-all duration-300"
        onclick={toggleForgotPassword}
      >
        {$LL.returnToLogin()}
      </button>
      <button
        type="submit"
        class="px-6 py-2 text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {$LL.resetPassword()}
      </button>
    </div>
  </form>
{/if}

{#if mfaRequired}
  <form onsubmit={authMfa} class="space-y-6" name="authMfa">
    <div class="space-y-2">
      <label
        class="block text-sm font-medium text-gray-700 dark:text-gray-300"
        for="mfaToken"
      >
        {$LL.mfaTokenLabel()}
      </label>
      <TextInput
        bind:value={mfaToken}
        placeholder="Enter code"
        id="mfaToken"
        name="mfaToken"
        required
        icon={Shield}
        inputmode="numeric"
        pattern="[0-9]*"
        autocomplete="one-time-code"
      />
    </div>

    <div class="pt-4">
      <button
        type="submit"
        class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <span class="absolute left-0 inset-y-0 flex items-center pl-3">
          <Shield
            class="h-5 w-5 text-purple-300 group-hover:text-purple-200"
            aria-hidden="true"
          />
        </span>
        {$LL.login()}
      </button>
    </div>
  </form>
{/if}

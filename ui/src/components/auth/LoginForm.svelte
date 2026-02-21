<script lang="ts">
  import { AppConfig, appRoutes, PathPrefix } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { Lock, Mail, Shield } from '@lucide/svelte';
  import PasswordInput from '../forms/PasswordInput.svelte';
  import Google from '../icons/Google.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { SessionUser } from '../../types/user';
  import OtpForm from './OtpForm.svelte';
  import ForgotPasswordForm from './ForgotPasswordForm.svelte';

  interface Props {
    registerLink?: string;
    targetPage?: any;
    router: any;
    xfetch?: ApiClient;
    notifications?: NotificationService;
  }

  let {
    registerLink = '',
    targetPage = appRoutes.landing,
    router,
    xfetch = (async () => new Response()) as ApiClient,
    notifications = {
      success: () => {},
      danger: () => {},
      warning: () => {},
      info: () => {},
      show: () => {},
      removeToast: () => {},
    } as NotificationService,
  }: Props = $props();

  const { LdapEnabled, GoogleAuthEnabled, HeaderAuthEnabled, OIDCAuthEnabled, OIDCProviderName } = AppConfig;
  const authEndpoint = LdapEnabled ? '/api/auth/ldap' : '/api/auth';

  let email = $state('');
  let password = $state('');
  let forgotPassword = $state(false);

  let mfaRequired = $state(false);
  let mfaUser: any = null;
  let mfaSessionId: string | null = null;

  function googleLogin() {
    window.location.href = `${PathPrefix}/oauth/google/login`;
  }

  function oidcLogin() {
    window.location.href = `${PathPrefix}/oauth/${OIDCProviderName.toLowerCase()}/login`;
  }

  function toggleForgotPassword() {
    forgotPassword = !forgotPassword;
  }

  function handleLoginSubmit(e: Event) {
    e.preventDefault();

    const body = {
      email,
      password,
    };

    xfetch(authEndpoint, { body, skip401Redirect: true })
      .then((res: any) => res.json())
      .then(function (result: any) {
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
          user.create(newUser as SessionUser);
          if (u.theme !== 'auto') {
            localStorage.setItem('theme', u.theme);
          } else {
            localStorage.removeItem('theme');
          }
          (window as any).setTheme();
          router.route(targetPage, true);
        }
      })
      .catch(function () {
        notifications.danger($LL.authError());
      });
  }

  function authMfa(token: string) {
    const body = {
      passcode: token,
      sessionId: mfaSessionId,
    };

    xfetch('/api/auth/mfa', { body, skip401Redirect: true })
      .then((res: any) => res.json())
      .then(function () {
        user.create(mfaUser);
        router.route(targetPage, true);
      })
      .catch(function () {
        notifications.danger($LL.mfaAuthError());
      });
  }

  function sendPasswordReset(email: string) {
    const body = {
      email,
    };

    xfetch('/api/auth/forgot-password', { body })
      .then(function () {
        notifications.success(
          $LL.sendResetPasswordSuccess({
            email,
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
      <Lock class="h-5 w-5 text-purple-300 group-hover:text-purple-200" aria-hidden="true" />
    </span>
    {$LL.loginWithSSO({ provider: OIDCProviderName })}
  </button>
{/if}

{#if !OIDCAuthEnabled && !forgotPassword && !mfaRequired}
  <form onsubmit={handleLoginSubmit} class="space-y-6" name="login" id="login">
    <TextInput
      id="email"
      data-testid="username"
      placeholder={$LL.enterYourEmail()}
      required
      bind:value={email}
      icon={Mail}
      autocomplete="email"
    />

    <PasswordInput
      bind:value={password}
      placeholder={$LL.yourPasswordPlaceholder()}
      id="password"
      name="password"
      data-testid="password"
      required
    />

    {#if !LdapEnabled}
      <div class="flex items-center justify-end text-sm">
        <a
          class="font-medium text-purple-600 dark:text-purple-400 hover:text-purple-500 dark:hover:text-purple-300 transition-all duration-300 cursor-pointer"
          onclick={e => {
            e.preventDefault();
            toggleForgotPassword();
          }}
          href="/login"
        >
          {$LL.forgotPassword()}
        </a>
      </div>
    {/if}

    <button
      data-testid="login"
      type="submit"
      class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span class="absolute left-0 inset-y-0 flex items-center ps-3">
        <Lock class="h-5 w-5 text-purple-300 group-hover:text-purple-200" aria-hidden="true" />
      </span>
      {$LL.login()}
    </button>
  </form>
  {#if GoogleAuthEnabled && !HeaderAuthEnabled && !LdapEnabled}
    <div class="w-full space-y-4 mt-4">
      <div class="flex items-center space-x-2">
        <div class="flex-grow h-px bg-gray-300"></div>
        <span class="uppercase text-sm text-gray-500 dark:text-gray-400 font-medium">Or continue with</span>
        <div class="flex-grow h-px bg-gray-300"></div>
      </div>
      <button
        onclick={googleLogin}
        class="inline-flex w-full items-center justify-center rounded-md border border-gray-600 dark:border-gray-400 hover:border-indigo-600 dark:hover:border-purple-400 px-4 py-2 shadow-sm disabled:cursor-wait disabled:opacity-50"
      >
        <span class="sr-only">Sign in with Google</span>
        <Google class="w-8 h-8" />
      </button>
    </div>
  {/if}
  {#if registerLink !== ''}
    <div class="m-auto mt-6 w-fit md:mt-8">
      <span class="m-auto dark:text-gray-400"
        >{$LL.createAccountTagline()}
        <a class="font-semibold text-indigo-600 dark:text-indigo-100" href={registerLink}>{$LL.createAccount()}</a>
      </span>
    </div>
  {/if}
{/if}

{#if forgotPassword}
  <ForgotPasswordForm onSubmit={sendPasswordReset} {toggleForgotPassword} />
{/if}

{#if mfaRequired}
  <OtpForm {authMfa} />
{/if}

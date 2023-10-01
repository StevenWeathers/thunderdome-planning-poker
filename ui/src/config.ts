declare global {
  interface Window {
    appConfig: any;
  }
}

const locales = {
  de: 'Deutsch',
  en: 'English',
  es: 'Español',
  fr: 'Français',
  pt: 'Português',
  ru: 'Русский',
  fa: 'Persian',
  it: 'Italiano',
};
const rtlLanguages = ['fa'];

const AppConfig =
  typeof window.appConfig != 'undefined'
    ? window.appConfig
    : {
        PathPrefix: '',
        DefaultLocale: 'en',
      };
const {
  PathPrefix,
  DefaultLocale: fallbackLocale,
  FriendlyUIVerbs,
} = AppConfig;

const defaultAppRoutes = {
  landing: `${PathPrefix}/`,
  register: `${PathPrefix}/register`,
  login: `${PathPrefix}/login`,
  resetPwd: `${PathPrefix}/reset-password`,
  verifyAct: `${PathPrefix}/verify-account`,
  profile: `${PathPrefix}/profile`,
  games: `${PathPrefix}/battles`,
  game: `${PathPrefix}/battle`,
  retros: `${PathPrefix}/retros`,
  retro: `${PathPrefix}/retro`,
  storyboards: `${PathPrefix}/storyboards`,
  storyboard: `${PathPrefix}/storyboard`,
  teams: `${PathPrefix}/teams`,
  organization: `${PathPrefix}/organization`,
  team: `${PathPrefix}/team`,
  admin: `${PathPrefix}/admin`,
  adminPokerGames: `${PathPrefix}/admin/battles`,
  adminRetros: `${PathPrefix}/admin/retros`,
  adminStoryboards: `${PathPrefix}/admin/storyboards`,
  adminTeams: `${PathPrefix}/admin/teams`,
  adminOrganizations: `${PathPrefix}/admin/organizations`,
  adminApiKeys: `${PathPrefix}/admin/apikeys`,
  adminAlerts: `${PathPrefix}/admin/alerts`,
  adminUsers: `${PathPrefix}/admin/users`,
  subscriptionPricing: `${PathPrefix}/subscriptions/pricing`,
  subscriptionConfirmation: `${PathPrefix}/subscriptions/confirmation`,
};
const friendlyAppRoutes = {
  ...defaultAppRoutes,
  games: `${PathPrefix}/games`,
  game: `${PathPrefix}/game`,
};
const appRoutes = FriendlyUIVerbs ? friendlyAppRoutes : defaultAppRoutes;

export {
  locales,
  fallbackLocale,
  appRoutes,
  PathPrefix,
  AppConfig,
  rtlLanguages,
};

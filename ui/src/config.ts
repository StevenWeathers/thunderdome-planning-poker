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
        Subscription: {},
        RetroDefaultTemplateID: '',
      };

const { PathPrefix, DefaultLocale } = AppConfig;

const adminPrefix = `${PathPrefix}/admin`;
const appRoutes = {
  landing: `${PathPrefix}/`,
  dashboard: `${PathPrefix}/dashboard`,
  register: `${PathPrefix}/register`,
  login: `${PathPrefix}/login`,
  resetPwd: `${PathPrefix}/reset-password`,
  verifyAct: `${PathPrefix}/verify-account`,
  profile: `${PathPrefix}/profile`,
  battles: `${PathPrefix}/battles`,
  battle: `${PathPrefix}/battle`,
  games: `${PathPrefix}/games`,
  game: `${PathPrefix}/game`,
  retros: `${PathPrefix}/retros`,
  retro: `${PathPrefix}/retro`,
  storyboards: `${PathPrefix}/storyboards`,
  storyboard: `${PathPrefix}/storyboard`,
  teams: `${PathPrefix}/teams`,
  organization: `${PathPrefix}/organization`,
  team: `${PathPrefix}/team`,
  admin: adminPrefix,
  adminSupportTickets: `${adminPrefix}/support-tickets`,
  adminPokerGames: `${adminPrefix}/games`,
  adminRetros: `${adminPrefix}/retros`,
  adminStoryboards: `${adminPrefix}/storyboards`,
  adminTeams: `${adminPrefix}/teams`,
  adminOrganizations: `${adminPrefix}/organizations`,
  adminApiKeys: `${adminPrefix}/apikeys`,
  adminAlerts: `${adminPrefix}/alerts`,
  adminUsers: `${adminPrefix}/users`,
  adminSubscriptions: `${adminPrefix}/subscriptions`,
  adminEstimationScales: `${adminPrefix}/estimation-scales`,
  adminRetroTemplates: `${adminPrefix}/retro-templates`,
  adminProjects: `${adminPrefix}/projects`,
  subscriptionPricing: `${PathPrefix}/subscriptions/pricing`,
  subscriptionConfirmation: `${PathPrefix}/subscriptions/confirmation`,
  privacyPolicy: `${PathPrefix}/privacy-policy`,
  termsConditions: `${PathPrefix}/terms-conditions`,
  support: `${PathPrefix}/support`,
  openSource: `${PathPrefix}/open-source`,
  invite: `${PathPrefix}/invite`,
  projects: `${PathPrefix}/projects`,
};

export {
  locales,
  DefaultLocale,
  appRoutes,
  PathPrefix,
  AppConfig,
  rtlLanguages,
};

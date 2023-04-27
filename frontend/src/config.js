const locales = {
    de: 'Deutsch',
    en: 'English',
    es: 'Español',
    fr: 'Français',
    pt: 'Português',
    ru: 'Русский',
    fa: 'Persian',
    it: 'Italiano',
}

const AppConfig =
    typeof appConfig != 'undefined'
        ? appConfig
        : {
              PathPrefix: '',
              DefaultLocale: 'en',
          }
const { PathPrefix, DefaultLocale: fallbackLocale, FriendlyUIVerbs } = AppConfig

const defaultAppRoutes = {
    landing: `${PathPrefix}/`,
    register: `${PathPrefix}/register`,
    login: `${PathPrefix}/login`,
    resetPwd: `${PathPrefix}/reset-password`,
    verifyAct: `${PathPrefix}/verify-account`,
    profile: `${PathPrefix}/profile`,
    battles: `${PathPrefix}/battles`,
    battle: `${PathPrefix}/battle`,
    retros: `${PathPrefix}/retros`,
    retro: `${PathPrefix}/retro`,
    storyboards: `${PathPrefix}/storyboards`,
    storyboard: `${PathPrefix}/storyboard`,
    teams: `${PathPrefix}/teams`,
    organization: `${PathPrefix}/organization`,
    team: `${PathPrefix}/team`,
    admin: `${PathPrefix}/admin`,
    adminBattles: `${PathPrefix}/admin/battles`,
    adminRetros: `${PathPrefix}/admin/retros`,
    adminStoryboards: `${PathPrefix}/admin/storyboards`,
    adminTeams: `${PathPrefix}/admin/teams`,
    adminOrganizations: `${PathPrefix}/admin/organizations`,
    adminApiKeys: `${PathPrefix}/admin/apikeys`,
    adminAlerts: `${PathPrefix}/admin/alerts`,
    adminUsers: `${PathPrefix}/admin/users`,
}
const friendlyAppRoutes = {
    ...defaultAppRoutes,
    battles: `${PathPrefix}/games`,
    battle: `${PathPrefix}/game`,
}
const appRoutes = FriendlyUIVerbs ? friendlyAppRoutes : defaultAppRoutes

export { locales, fallbackLocale, appRoutes, PathPrefix, AppConfig }

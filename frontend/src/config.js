const locales = {
    de: 'Deutsch',
    en: 'English',
    es: 'Español',
    fr: 'Français',
    pt: 'Português',
    ru: 'Русский',
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
    admin: `${PathPrefix}/admin`,
    battles: `${PathPrefix}/battles`,
    battle: `${PathPrefix}/battle`,
    retros: `${PathPrefix}/retros`,
    retro: `${PathPrefix}/retro`,
    storyboards: `${PathPrefix}/storyboards`,
    storyboard: `${PathPrefix}/storyboard`,
    teams: `${PathPrefix}/teams`,
    organization: `${PathPrefix}/organization`,
    team: `${PathPrefix}/team`,
}
const friendlyAppRoutes = {
    ...defaultAppRoutes,
    battles: `${PathPrefix}/games`,
    battle: `${PathPrefix}/game`,
}
const appRoutes = FriendlyUIVerbs ? friendlyAppRoutes : defaultAppRoutes

export { locales, fallbackLocale, appRoutes, PathPrefix, AppConfig }

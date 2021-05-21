const locales = {
    de: 'Deutsch',
    en: 'English',
    ru: 'Русский',
}

const { PathPrefix, DefaultLocale: fallbackLocale } = appConfig

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
    organizations: `${PathPrefix}/organizations`,
    organization: `${PathPrefix}/organization`,
    team: `${PathPrefix}/team`,
    alerts: `${PathPrefix}/alerts`,
}
const friendlyAppRoutes = {
    ...defaultAppRoutes,
    battles: `${PathPrefix}/games`,
    battle: `${PathPrefix}/game`,
}
const appRoutes = appConfig.FriendlyUIVerbs
    ? friendlyAppRoutes
    : defaultAppRoutes

export { locales, fallbackLocale, appRoutes, PathPrefix }

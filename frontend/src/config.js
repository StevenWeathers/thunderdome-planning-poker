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

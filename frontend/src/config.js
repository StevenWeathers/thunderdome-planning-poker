const locales = {
    de: 'Deutsch',
    en: 'English',
    ru: 'Русский',
}

const fallbackLocale = appConfig.DefaultLocale

const defaultAppRoutes = {
    landing: '/',
    register: '/register',
    login: '/login',
    resetPwd: '/reset-password',
    verifyAct: '/verify-account',
    profile: '/profile',
    admin: '/admin',
    battles: '/battles',
    battle: '/battle',
}
const friendlyAppRoutes = {
    ...defaultAppRoutes,
    battles: '/games',
    battle: '/game',
}
const appRoutes = appConfig.FriendlyUIVerbs
    ? friendlyAppRoutes
    : defaultAppRoutes

export { locales, fallbackLocale, appRoutes }

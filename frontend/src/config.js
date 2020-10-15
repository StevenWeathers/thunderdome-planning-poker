const locales = {
    de: 'Deutsch',
    en: 'English',
    ru: 'Русский',
}

const fallbackLocale = appConfig.DefaultLocale

const defaultAppRoutes = {
    landing: '/',
    enlist: '/enlist',
    battles: '/battles',
    battle: '/battle',
    login: '/login',
    resetPwd: '/reset-password',
    verifyAct: '/verify-account',
    profile: '/warrior-profile',
    admin: '/admin',
}
const friendlyAppRoutes = {
    ...defaultAppRoutes,
    enlist: '/register',
    battles: '/sessions',
    battle: '/session',
    profile: '/user-profile',
}
const appRoutes = appConfig.FriendlyUIVerbs
    ? friendlyAppRoutes
    : defaultAppRoutes

export { locales, fallbackLocale, appRoutes }

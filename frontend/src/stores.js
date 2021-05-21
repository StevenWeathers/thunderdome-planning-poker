import { writable } from 'svelte/store'
import Cookies from 'js-cookie'

const { PathPrefix, CookieName } = appConfig
const cookiePath = `${PathPrefix}/`

function initWarrior() {
    const { subscribe, set, update } = writable(
        Cookies.getJSON(CookieName) || {},
    )

    return {
        subscribe,
        create: warrior => {
            Cookies.set(CookieName, warrior, {
                expires: 365,
                SameSite: 'strict',
                path: cookiePath,
            })
            set(warrior)
        },
        update: warrior => {
            Cookies.set(CookieName, warrior, {
                expires: 365,
                SameSite: 'strict',
                path: cookiePath,
            })
            update(w => (w = warrior))
        },
        delete: () => {
            Cookies.remove(CookieName, { path: cookiePath })
            set({})
        },
    }
}

export const warrior = initWarrior()

function initActiveAlerts() {
    const { subscribe, update } = writable(ActiveAlerts)

    return {
        subscribe,
        update: alerts => {
            update(a => (a = alerts))
        },
    }
}

export const activeAlerts = initActiveAlerts()

function initDismissedAlerts() {
    const dismissKey = 'dismissed_alerts'
    const dismissedAlerts = JSON.parse(localStorage.getItem(dismissKey)) || []
    const { subscribe, update } = writable(dismissedAlerts)

    return {
        subscribe,
        dismiss: (actives, dismisses) => {
            const validAlerts = actives.map((prev, alert) => alert.id)
            let alertsToDismiss = [
                ...dismisses.filter(alert => validAlerts.includes(alert.id)),
            ]
            localStorage.setItem(dismissKey, JSON.stringify(alertsToDismiss))
            update(a => (a = alertsToDismiss))
        },
    }
}

export const dismissedAlerts = initDismissedAlerts()

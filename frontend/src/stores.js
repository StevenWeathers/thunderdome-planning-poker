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

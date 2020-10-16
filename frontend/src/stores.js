import { writable } from 'svelte/store'
import Cookies from 'js-cookie'

const cookieName = appConfig.CookieName

function initWarrior() {
    const { subscribe, set, update } = writable(
        Cookies.getJSON(cookieName) || {},
    )

    return {
        subscribe,
        create: warrior => {
            Cookies.set(cookieName, warrior, {
                expires: 365,
                SameSite: 'strict',
            })
            set(warrior)
        },
        update: warrior => {
            Cookies.set(cookieName, warrior, {
                expires: 365,
                SameSite: 'strict',
            })
            update(w => (w = warrior))
        },
        delete: () => {
            Cookies.remove(cookieName)
            set({})
        },
    }
}

export const warrior = initWarrior()

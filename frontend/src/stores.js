import { writable } from 'svelte/store'
import Cookies from 'js-cookie'

function initWarrior() {
    const { subscribe, set, update } = writable(
        Cookies.getJSON('warrior') || {},
    )

    return {
        subscribe,
        create: warrior => {
            Cookies.set('warrior', warrior, {
                expires: 365,
                SameSite: 'strict',
            })
            set(warrior)
        },
        update: warrior => {
            Cookies.set('warrior', warrior, {
                expires: 365,
                SameSite: 'strict',
            })
            update(w => (w = warrior))
        },
        delete: () => {
            Cookies.remove('warrior')
            set({})
        },
    }
}

export const warrior = initWarrior()

import { writable } from 'svelte/store'
import Cookies from 'js-cookie'

function initUser () {
  const { subscribe, set } = writable(Cookies.getJSON('user') || {})

  return {
    subscribe,
    create: (user) => {
      Cookies.set('user', user)
      set(user)
    }
  }
}

export const user = initUser()

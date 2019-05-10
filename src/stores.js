import { writable } from 'svelte/store'
import Cookies from 'js-cookie'

function initWarrior () {
  const { subscribe, set } = writable(Cookies.getJSON('warrior') || {})

  return {
    subscribe,
    create: (warrior) => {
      Cookies.set('warrior', warrior)
      set(warrior)
    }
  }
}

export const warrior = initWarrior()

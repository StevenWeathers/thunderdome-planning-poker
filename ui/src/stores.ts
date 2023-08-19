import { derived, writable } from 'svelte/store';
import Cookies from 'js-cookie';
import { AppConfig, rtlLanguages } from './config';
import { locale } from './i18n/i18n-svelte';

const { PathPrefix, CookieName } = AppConfig;
const cookiePath = `${PathPrefix}/`;

declare global {
  let ActiveAlerts: any;
}

function initWarrior() {
  const { subscribe, set, update } = writable(
    JSON.parse(Cookies.get(CookieName) || '{}'),
  );

  return {
    subscribe,
    create: warrior => {
      Cookies.set(CookieName, JSON.stringify(warrior), {
        expires: 365,
        SameSite: 'strict',
        path: cookiePath,
      });
      set(warrior);
    },
    update: warrior => {
      Cookies.set(CookieName, JSON.stringify(warrior), {
        expires: 365,
        SameSite: 'strict',
        path: cookiePath,
      });
      update(w => (w = warrior));
    },
    delete: () => {
      Cookies.remove(CookieName, { path: cookiePath });
      set({});
    },
  };
}

export const user = initWarrior();

function initActiveAlerts() {
  const activeAlerts = typeof ActiveAlerts != 'undefined' ? ActiveAlerts : [];
  const { subscribe, update } = writable(activeAlerts);

  return {
    subscribe,
    update: alerts => {
      update(a => (a = alerts));
    },
  };
}

export const activeAlerts = initActiveAlerts();

function initDismissedAlerts() {
  const dismissKey = 'dismissed_alerts';
  const dismissedAlerts = JSON.parse(localStorage.getItem(dismissKey)) || [];
  const { subscribe, update } = writable(dismissedAlerts);

  return {
    subscribe,
    dismiss: (actives, dismisses) => {
      const validAlerts = actives.map((prev, alert) => alert.id);
      let alertsToDismiss = [
        ...dismisses.filter(alert => validAlerts.includes(alert.id)),
      ];
      localStorage.setItem(dismissKey, JSON.stringify(alertsToDismiss));
      update((a: any) => (a = alertsToDismiss));
    },
  };
}

export const dir = derived(locale, $locale =>
  rtlLanguages.includes($locale) ? 'rtl' : 'ltr',
);

export const dismissedAlerts = initDismissedAlerts();

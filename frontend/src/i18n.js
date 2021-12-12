import { derived, get, writable } from 'svelte/store'
import {
    _,
    addMessages,
    date,
    dictionary,
    getLocaleFromNavigator,
    init,
    locale,
    number,
} from 'svelte-i18n'

import { AppConfig, fallbackLocale, locales, PathPrefix } from './config.js'

const { AppVersion, FriendlyUIVerbs } = AppConfig
const verbsType = FriendlyUIVerbs ? 'friendly' : 'default'
const MESSAGE_FILE_URL_BASE = `${PathPrefix}/lang/{locale}.json?v=${AppVersion}`
const MESSAGE_FILE_URL_ADD = `${PathPrefix}/lang/${verbsType}/{locale}.json?v=${AppVersion}`

let _activeLocale

// Internal store for tracking network
// loading state
const isDownloading = writable(false)

function setupI18n(options = {}) {
    const locale_ = supported(options.withLocale || getLocaleFromNavigator())

    // Initialize svelte-i18n
    init({ initialLocale: locale_ })

    // Don't re-download translation files
    if (!hasLoadedLocale(locale_)) {
        isDownloading.set(true)

        const messagesFileUrlBase = MESSAGE_FILE_URL_BASE.replace(
            '{locale}',
            locale_,
        )
        const messageFileUrlAdd = MESSAGE_FILE_URL_ADD.replace(
            '{locale}',
            locale_,
        )

        return Promise.all([
            loadJson(messagesFileUrlBase),
            loadJson(messageFileUrlAdd),
        ]).then(([base, additional]) => {
            addMessages(locale_, base)
            addMessages(locale_, additional)

            // Configure svelte-i18n to use the locale
            _activeLocale = locale_
            locale.set(locale_)

            isDownloading.set(false)
        })
    }
}

const isLocaleLoaded = derived(
    [isDownloading, dictionary],
    ([$isDownloading, $dictionary]) =>
        !$isDownloading &&
        $dictionary[_activeLocale] &&
        Object.keys($dictionary[_activeLocale]).length > 0,
)

const dir = derived(locale, $locale => ($locale === 'ar' ? 'rtl' : 'ltr'))

function loadJson(url) {
    return fetch(url).then(response => response.json())
}

function hasLoadedLocale(locale) {
    // If the svelte-i18n dictionary has an entry for the
    // locale, then the locale has already been added
    return get(dictionary)[locale]
}

function language(locale) {
    return locale.replace('_', '-').split('-')[0]
}

function supported(locale) {
    const localeKeys = Object.keys(locales)

    // check if supported locales include exact variant or primary variant
    // e.g. checks for en-US then en
    if (localeKeys.includes(locale)) {
        return locale
    } else if (localeKeys.includes(language(locale))) {
        return language(locale)
    } else {
        return fallbackLocale
    }
}

// We expose the svelte-i18n _ store so that our app has
// a single API for i18n
export { _, setupI18n, isLocaleLoaded, locale, locales, dir, date, number }

// Most of this setup came from
// https://medium.com/i18n-and-l10n-resources-for-developers/a-step-by-step-guide-to-svelte-localization-with-svelte-i18n-v3-2c3ff0d645b8
// and
// https://lokalise.com/blog/svelte-i18n/

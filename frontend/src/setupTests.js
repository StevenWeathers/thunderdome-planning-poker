import en from '../public/lang/en.json'
import defaultEn from '../public/lang/default/en.json'
import { addMessages, init } from 'svelte-i18n'

addMessages('en', en)
addMessages('en', defaultEn)
init({
    fallbackLocale: 'en',
    initialLocale: 'en',
})

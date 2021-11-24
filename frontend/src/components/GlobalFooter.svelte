<script>
    import GithubIcon from './icons/Github.svelte'
    import { AppConfig } from '../config'
    import { _ } from '../i18n'

    const { AppVersion } = AppConfig
    const footerLinkClasses = 'no-underline text-teal-500 hover:text-teal-800'
</script>

<footer class="p-6 text-center">
    <GithubIcon />
    &nbsp;
    <a
        href="https://github.com/StevenWeathers/thunderdome-planning-poker"
        class="{footerLinkClasses}"
    >
        {$_('appName')}
    </a>
    {@html $_('footer.authoredBy', {
        values: {
            authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}">`,
            authorClose: `</a>`,
        },
    })}
    {@html $_('footer.license', {
        values: {
            licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}">`,
            licenseClose: `</a>`,
        },
    })}
    <br />
    {@html $_('footer.poweredBy', {
        values: {
            svelteOpen: `<a href="https://svelte.dev/" class="${footerLinkClasses}">`,
            svelteClose: `</a>`,
            goOpen: `<a href="https://golang.org/" class="${footerLinkClasses}">`,
            goClose: `</a>`,
        },
    })}
    <div class="text-sm text-gray-500">
        {$_('appVersion', { values: { version: AppVersion } })}
    </div>
</footer>

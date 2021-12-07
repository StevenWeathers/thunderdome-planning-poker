<script>
    import GithubIcon from './icons/Github.svelte'
    import { AppConfig } from '../config'
    import { _ } from '../i18n'

    const { AppVersion } = AppConfig
    const footerLinkClasses = 'no-underline text-blue-500 hover:text-blue-800'
</script>

<footer
    class="border-t px-16 py-10 text-center font-rajdhani text-xl"
    style="box-shadow: inset 0 7px 30px -7px rgba(0,0,0,0.1);"
>
    <GithubIcon />
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
    <br />
    {@html $_('footer.license', {
        values: {
            licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}">`,
            licenseClose: `</a>`,
        },
    })}
    <br />
    <div class="text-sm text-gray-500">
        {$_('appVersion', { values: { version: AppVersion } })}
    </div>
</footer>

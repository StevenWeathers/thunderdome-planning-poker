import svelte from 'rollup-plugin-svelte'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import livereload from 'rollup-plugin-livereload'
import { terser } from 'rollup-plugin-terser'
import copy from 'rollup-plugin-copy'
import del from 'rollup-plugin-delete'
import postcss from 'rollup-plugin-postcss'
import autoPreprocess from 'svelte-preprocess'
import html from 'rollup-plugin-bundle-html'

const production = !process.env.ROLLUP_WATCH

export default {
    input: 'frontend/src/main.js',
    output: {
        sourcemap: false,
        format: 'iife',
        name: 'app',
        file: `dist/js/bundle.[hash].js`,
    },
    plugins: [
        del({ targets: 'dist/*' }),
        svelte({
            preprocess: autoPreprocess({
                postcss: {
                    configFilePath: 'build/postcss.config.js'
                }
            }),
            // enable run-time checks when not in production
            dev: !production,
            // we'll extract any component CSS out into
            // a separate file — better for performance
            css: css => {
                css.write(`dist/css/bundle.[hash].css`, false)
            },
        }),
        postcss({
            extract: `dist/css/tailwind.[hash].css`,
            config: {
                path: 'build/postcss.config.js'
            }
        }),
        // If you have external dependencies installed from
        // npm, you'll most likely need these plugins. In
        // some cases you'll need additional configuration —
        // consult the documentation for details:
        // https://github.com/rollup/rollup-plugin-commonjs
        resolve(),
        commonjs(),

        html({
            template: 'frontend/public/index.html',
            dest: 'dist',
            filename: 'index.html',
            absolute: true,
        }),

        copy({
            targets: {
                'frontend/public/img': 'dist/img',
            },
        }),

        // Watch the `dist` directory and refresh the
        // browser on changes when not in production
        !production && livereload('dist'),

        // If we're building for production (npm run build
        // instead of npm run dev), minify
        production && terser(),
    ],
}

import svelte from 'rollup-plugin-svelte'
import resolve from '@rollup/plugin-node-resolve'
import commonjs from '@rollup/plugin-commonjs'
import livereload from 'rollup-plugin-livereload'
import { terser } from 'rollup-plugin-terser'
import copy from 'rollup-plugin-copy'
import del from 'rollup-plugin-delete'
import sveltePreprocess from 'svelte-preprocess'
import html from '@rollup/plugin-html'
import { template } from './buildHtmlTemplate.js'
import postcss from 'rollup-plugin-postcss'
import postcssNesting from 'postcss-nesting'

const production = !process.env.ROLLUP_WATCH

export default {
  input: 'frontend/src/main.js',
  output: {
    sourcemap: false,
    name: 'thunderdome',
    format: 'esm',
    dir: `dist/static/`,
    entryFileNames: '[name]-[hash].js',
    chunkFileNames: '[name].[hash].js',
    assetFileNames: '[name].[hash].[extension]'
  },
  plugins: [
    del({ targets: 'dist/*' }),
    svelte({
      preprocess: sveltePreprocess({
        sourceMap: !production,
        postcss: {
          plugins: [
            require('postcss-import'),
            require('tailwindcss/nesting'),
            require('tailwindcss'),
            require('autoprefixer')
          ],
        },
      }),
    }),
    postcss({
      plugins: [postcssNesting(), (production && require('cssnano'))],
      extract: true
    }),
    // If you have external dependencies installed from
    // npm, you'll most likely need these plugins. In
    // some cases you'll need additional configuration —
    // consult the documentation for details:
    // https://github.com/rollup/rollup-plugin-commonjs
    resolve({ browser: true, dedupe: ['svelte'] }),
    commonjs(),

    html({
      title: 'Thunderdome - Open Source Agile Planning Poker app',
      publicPath: '{{.AppConfig.PathPrefix}}/static/',
      template
    }),

    copy({
      targets: [
        {
          src: 'frontend/public/img',
          dest: 'dist'
        },
        {
          src: 'frontend/public/lang',
          dest: 'dist',
        }
      ]
    }),

    // Watch the `dist` directory and refresh the
    // browser on changes when not in production
    !production && livereload('dist'),

    // If we're building for production (npm run build
    // instead of npm run dev), minify
    production && terser(),
  ],
}

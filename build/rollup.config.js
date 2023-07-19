const svelte = require('rollup-plugin-svelte')
const resolve = require('@rollup/plugin-node-resolve')
const commonjs = require('@rollup/plugin-commonjs')
const livereload = require('rollup-plugin-livereload')
const terser = require('@rollup/plugin-terser')
const copy = require('rollup-plugin-copy')
const del = require('rollup-plugin-delete')
const sveltePreprocess = require('svelte-preprocess')
const html = require('@rollup/plugin-html')
const template = require('./buildHtmlTemplate.js')
const postcss = require('rollup-plugin-postcss')
const postcssNesting = require('postcss-nesting')
const typescript = require('@rollup/plugin-typescript')

const production = !process.env.ROLLUP_WATCH

module.exports = {
  input: 'ui/src/main.ts',
  output: {
    sourcemap: !production,
    name: 'thunderdome',
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
    typescript({ sourceMap: !production }),
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
          src: 'ui/public/img',
          dest: 'dist'
        },
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

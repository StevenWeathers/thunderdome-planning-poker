import svelte from 'rollup-plugin-svelte'
import resolve from 'rollup-plugin-node-resolve'
import commonjs from 'rollup-plugin-commonjs'
import { terser } from 'rollup-plugin-terser'
import copy from 'rollup-plugin-copy'
import del from 'rollup-plugin-delete'
import html from 'rollup-plugin-bundle-html'

export default {
  input: 'src/main.js',
  output: {
    sourcemap: false,
    format: 'iife',
    name: 'app',
    file: 'dist/js/bundle.[hash].js'
  },
  plugins: [
    del({ targets: 'dist/*' }),
    svelte({
      // enable run-time checks when not in production
      dev: false,
      // we'll extract any component CSS out into
      // a separate file — better for performance
      css: css => {
        css.write('dist/css/bundle.[hash].css', false)
      }
    }),

    // If you have external dependencies installed from
    // npm, you'll most likely need these plugins. In
    // some cases you'll need additional configuration —
    // consult the documentation for details:
    // https://github.com/rollup/rollup-plugin-commonjs
    resolve(),
    commonjs(),

    terser(),

    html({
      template: 'public/index.html',
      dest: 'dist',
      filename: 'index.html',
      absolute: true
    }),

    copy({
      targets: {
        'public/img': 'dist/img',
        'public/css': 'dist/css'
      }
    })
  ]
}

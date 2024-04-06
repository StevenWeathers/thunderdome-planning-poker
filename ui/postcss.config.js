import tailwind from 'tailwindcss'
import tailwindConfig from './tailwind.config.js'
import autoprefixer from 'autoprefixer'
import nesting from 'tailwindcss/nesting'
import cssimports from 'postcss-import'

export default {
  plugins: [
    cssimports,
    nesting,
    tailwind(tailwindConfig),
    autoprefixer
  ]
}
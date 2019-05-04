const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CopyPlugin = require('copy-webpack-plugin')
const path = require('path')

const mode = process.env.NODE_ENV || 'development'
const prod = mode === 'production'

module.exports = {
  entry: {
    bundle: ['./src/main.js']
  },
  resolve: {
    extensions: ['.mjs', '.js', '.svelte']
  },
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'js/[name].js',
    chunkFilename: 'js/[name].[id].js'
  },
  module: {
    rules: [
      {
        test: /\.svelte$/,
        exclude: /node_modules/,
        use: {
          loader: 'svelte-loader',
          options: {
            emitCss: true,
            hotReload: true
          }
        }
      },
      {
        test: /\.css$/,
        use: [
          prod ? MiniCssExtractPlugin.loader : 'style-loader',
          'css-loader'
        ]
      }
    ]
  },
  mode,
  plugins: [
    new MiniCssExtractPlugin({
      filename: 'css/[name].css'
    }),
    new CopyPlugin([
      { from: 'public', to: './' }
    ])
  ],
  devtool: prod ? false : 'source-map'
}

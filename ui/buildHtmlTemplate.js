const fs = require('fs')
const path = require('path')

const makeHtmlAttributes = (attributes) => {
  if (!attributes) {
    return ''
  }

  const keys = Object.keys(attributes)
  // eslint-disable-next-line no-param-reassign
  return keys.reduce((result, key) => (result += ` ${key}="${attributes[key]}"`), '')
}

module.exports = async ({
  attributes,
  files,
  meta,
  publicPath,
  title
}) => {
  const scripts = (files.js || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.script)
      return `<script src="${publicPath}${fileName}"${attrs} nonce="{{.Nonce}}"></script>`
    })
    .join('\n')

  const links = (files.css || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.link)
      return `<link href="${publicPath}${fileName}" rel="stylesheet"${attrs} nonce="{{.Nonce}}">`
    })
    .join('\n')

  const metas = meta
    .map((input) => {
      const attrs = makeHtmlAttributes(input)
      return `<meta${attrs}>`
    })
    .join('\n')

  const htmlFile = fs.readFileSync(path.resolve(__dirname, 'public/index.html'), 'utf8')
  return htmlFile.replace('${title}', title).replace('${metas}', metas).replace('${links}', links).replace('${scripts}', scripts)
}
import fs from 'fs'
import path from 'path'

const makeHtmlAttributes = (attributes) => {
  if (!attributes) {
    return ''
  }

  const keys = Object.keys(attributes)
  // eslint-disable-next-line no-param-reassign
  return keys.reduce((result, key) => (result += ` ${key}="${attributes[key]}"`), '')
}

export const template = async ({
  attributes,
  files,
  meta,
  publicPath,
  title
}) => {
  const scripts = (files.js || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.script)
      return `<script src="${publicPath}${fileName}"${attrs}></script>`
    })
    .join('\n')

  const links = (files.css || [])
    .map(({ fileName }) => {
      const attrs = makeHtmlAttributes(attributes.link)
      return `<link href="${publicPath}${fileName}" rel="stylesheet"${attrs}>`
    })
    .join('\n')

  const metas = meta
    .map((input) => {
      const attrs = makeHtmlAttributes(input)
      return `<meta${attrs}>`
    })
    .join('\n')

  const htmlFile = fs.readFileSync(path.resolve(__dirname, '../frontend/public/index.html'), 'utf8')
  return htmlFile.replace('${title}', title).replace('${metas}', metas).replace('${links}', links).replace('${scripts}', scripts)
}
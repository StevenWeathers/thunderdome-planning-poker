import { compute as computeScrollIntoView } from 'compute-scroll-into-view'

/* eslint no-bitwise: "off" */
/* eslint no-plusplus: "off" */

// https://github.com/lukeed/uid/blob/master/src/index.js
let IDX = 36
let HEX = ''

while (IDX--) {
  HEX += IDX.toString(36)
}

// Get a unique ID
/**
 * @param {number} [len]
 * @returns {string}
 */
const uid = (len) => {
  let str = ''
  let num = len || 11

  while (num--) {
    str += HEX[(Math.random() * 36) | 0]
  }

  return str
}

// Scroll an element into view if needed
/**
 * @param {HTMLElement | null} node
 * @param {HTMLElement} rootNode
 * @returns {void}
 */
const scrollIntoView = (node, rootNode) => {
  if (node === null) {
    return
  }

  const actions = computeScrollIntoView(node, {
    boundary: rootNode,
    block: 'center',
    scrollMode: 'if-needed'
  })

  // eslint-disable-next-line no-shadow
  actions.forEach(({ el, top }) => {
    el.scrollTop = top // eslint-disable-line no-param-reassign
  })
}

// Transform a string into a slug
/**
 * @param {any} str
 * @returns {string}
 */
const slugify = (str) =>
  str
    .toString()
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '') // Remove non-word [a-z0-9_], non-whitespace, non-hyphen characters
    .replace(/[\s_-]+/g, '-') // Swap any length of whitespace, underscore, hyphen characters with a single -
    .replace(/^-+|-+$/g, '') // Remove leading, trailing -

const keyCodes = {
  Enter: 13,
  Escape: 27,
  Space: 32,
  ArrowDown: 40,
  ArrowUp: 38,
  Backspace: 8,
  Characters: [
    48, // 0
    49, // 1
    50, // 2
    51, // 3
    52, // 4
    53, // 5
    54, // 6
    55, // 7
    56, // 8
    57, // 9
    65, // A
    66, // B
    67, // C
    68, // D
    69, // E
    70, // F
    71, // G
    72, // H
    73, // I
    74, // J
    75, // K
    76, // L
    77, // M
    78, // N
    79, // O
    80, // P
    81, // Q
    82, // R
    83, // S
    84, // T
    85, // U
    86, // V
    87, // W
    88, // X
    89, // Y
    90 // Z
  ]
}

/**
 * @param {Record<string, any>} timezones
 * @param {string[]} selection
 * @returns {Record<string, any>}
 */
const pick = (timezones, selection) => {
  const unique = Array.from(new Set([...selection]))

  return Object.keys(timezones).reduce((zones, zoneName) => {
    const picked = unique.includes(zoneName) ? timezones[zoneName] : {}
    return {
      ...zones,
      ...(Object.keys(picked).length > 0 && { [zoneName]: picked })
    }
  }, {})
}

// We take the grouped timezones and flatten them so that they can be easily searched
// e.g. { Europe: { 'London': 'Europe/London', 'Berlin': 'Europe/Berlin' } } => {'London': 'Europe/London', 'Berlin': 'Europe/Berlin' }
/**
 * @param {Record<string, any>} timezones
 * @returns {Record<string, any>}
 */
const ungroup = (timezones) =>
  Object.values(timezones).reduce(
    (values, zone) => ({ ...values, ...zone }),
    {}
  )

// Filter the list of zone labels to only those that match a search string
/**
 * @param {string} search
 * @param {Record<string, any>} zoneGroups
 * @returns {string[]}
 */
const filter = (search, zoneGroups) =>
  Object.entries(zoneGroups).reduce((zones, [zone, details]) => {
    if (details[0].toLowerCase().includes(search.toLowerCase())) {
      zones.push(zone)
    }
    return zones
  }, /** @type {string[]} */ ([]))

export { scrollIntoView, uid, slugify, keyCodes, ungroup, filter, pick }

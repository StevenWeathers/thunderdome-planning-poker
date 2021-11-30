const { Pool } = require('pg')

module.exports = function (db) {
  return new Pool({
    user: db.user,
    host: db.host,
    database: db.name,
    password: db.pass,
    port: db.port,
  })
}
import { Pool } from 'pg'

export const setupDB = function (db) {
    return new Pool({
        user: db.user,
        host: db.host,
        database: db.name,
        password: db.pass,
        port: db.port,
    })
}

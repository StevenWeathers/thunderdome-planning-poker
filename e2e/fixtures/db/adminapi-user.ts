export const adminAPIUser = {
    name: 'E2E Admin API User',
    email: 'e2eadminapi@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'ADMIN',
    apikey: 'Gssy-ffy.okeTA-3AJhCnY1sqeUvRPRHiNYIVUxs4',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
        [
            adminAPIUser.name,
            adminAPIUser.email,
            adminAPIUser.hashedPass,
            adminAPIUser.rank,
        ],
    )
    const id = newUser.rows[0].userid

    await pool.query('call verify_user_account($1);', [
        newUser.rows[0].verifyid,
    ])

    await pool.query(
        `INSERT INTO api_keys (id, user_id, name, active) VALUES ($1, $2, $3, TRUE);`,
        [
            'Gssy-ffy.e170ffced2ae5806aebc103f30255dc5cc1b9e203d6035aa817f2b7e6638f223',
            id,
            'test api key 2',
        ],
    )

    return {
        ...adminAPIUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(`SELECT id FROM users WHERE email = $1;`, [
        adminAPIUser.email,
    ])

    if (oldUser.rows.length) {
        await pool.query('call delete_user($1);', [oldUser.rows[0].id])
    }

    return {}
}

export const setupAdminAPIUser = {
    seed,
    teardown,
}

export const apiUser = {
    name: 'E2E API User',
    email: 'e2eapi@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'REGISTERED',
    apikey: '8MenPkY8.Vqvkh030vv7$rSyYs1gt++L0v7wKuVgR',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
        [apiUser.name, apiUser.email, apiUser.hashedPass, apiUser.rank],
    )
    const id = newUser.rows[0].userid

    await pool.query('call verify_user_account($1);', [
        newUser.rows[0].verifyid,
    ])

    await pool.query(
        `INSERT INTO api_keys (id, user_id, name, active) VALUES ($1, $2, $3, TRUE);`,
        [
            '8MenPkY8.cd737cbc4bdca1838bdcf1685b00a9a778261255c10193714d9ba1630b55b63c',
            id,
            'test apikey',
        ],
    )

    return {
        ...apiUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(`SELECT id FROM users WHERE email = $1;`, [
        apiUser.email,
    ])

    if (oldUser.rows.length) {
        await pool.query('call delete_user($1);', [oldUser.rows[0].id])
    }

    return {}
}

export const setupAPIUser = {
    seed,
    teardown,
}

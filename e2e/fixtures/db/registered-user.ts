export const registeredUser = {
    name: 'E2E Registered User',
    email: 'e2eregistered@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'REGISTERED',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
        [
            registeredUser.name,
            registeredUser.email,
            registeredUser.hashedPass,
            registeredUser.rank,
        ],
    )
    const id = newUser.rows[0].userid

    return {
        ...registeredUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(`SELECT id FROM users WHERE email = $1;`, [
        registeredUser.email,
    ])

    if (oldUser.rows.length) {
        await pool.query('call delete_user($1);', [oldUser.rows[0].id])
    }

    return {}
}

export const setupRegisteredUser = {
    seed,
    teardown,
}

export const registeredDeleteUser = {
    name: 'E2E Delete User',
    email: 'e2edelete@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'REGISTERED',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
        [
            registeredDeleteUser.name,
            registeredDeleteUser.email,
            registeredDeleteUser.hashedPass,
            registeredDeleteUser.rank,
        ],
    )
    const id = newUser.rows[0].userid

    return {
        ...registeredDeleteUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(`SELECT id FROM users WHERE email = $1;`, [
        registeredDeleteUser.email,
    ])

    if (oldUser.rows.length) {
        await pool.query('call delete_user($1);', [oldUser.rows[0].id])
    }

    return {}
}

export const setupDeleteRegisteredUser = {
    seed,
    teardown,
}

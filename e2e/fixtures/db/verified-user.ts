const verifiedUser = {
    name: 'E2E Verified User',
    email: 'e2everified@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'REGISTERED',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
        [
            verifiedUser.name,
            verifiedUser.email,
            verifiedUser.hashedPass,
            verifiedUser.rank,
        ],
    )

    await pool.query('call verify_user_account($1);', [
        newUser.rows[0].verifyid,
    ])
    const id = newUser.rows[0].userid

    return {
        ...verifiedUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(`SELECT id FROM users WHERE email = $1;`, [
        verifiedUser.email,
    ])

    if (oldUser.rows.length) {
        await pool.query('call delete_user($1);', [oldUser.rows[0].id])
    }

    return {}
}

export const setupVerifiedUser = {
    seed,
    teardown,
}

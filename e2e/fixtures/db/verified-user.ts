const verifiedUser = {
    name: 'E2E Verified User',
    email: 'e2everified@thunderdome.dev',
    password: 'kentRules!',
    hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
    rank: 'REGISTERED',
}

const seed = async pool => {
    const newUser = await pool.query(
        `SELECT userid, verifyid FROM thunderdome.user_register($1, $2, $3, $4);`,
        [
            verifiedUser.name,
            verifiedUser.email,
            verifiedUser.hashedPass,
            verifiedUser.rank,
        ],
    )

    await pool.query('call thunderdome.user_account_verify($1);', [
        newUser.rows[0].verifyid,
    ])
    const id = newUser.rows[0].userid

    return {
        ...verifiedUser,
        id,
    }
}

const teardown = async pool => {
    const oldUser = await pool.query(
        `SELECT id FROM thunderdome.users WHERE email = $1;`,
        [verifiedUser.email],
    )

    if (oldUser.rows.length) {
        await pool.query('DELETE FROM thunderdome.users WHERE id = $1;', [
            oldUser.rows[0].id,
        ])
    }

    return {}
}

export const setupVerifiedUser = {
    seed,
    teardown,
}

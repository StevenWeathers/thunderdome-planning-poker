const user = {
  name: 'Admin Test User',
  email: 'admin@thunderdome.dev',
  password: 'kentRules!',
  hashedPass: '$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m',
  rank: 'ADMIN'
}

const seed = async (pool) => {
  const newUser = await pool.query(
    `SELECT userid, verifyid FROM register_user($1, $2, $3, $4);`,
    [user.name, user.email, user.hashedPass, user.rank]
  )

  await pool.query('call verify_user_account($1);', [newUser.rows[0].verifyid])
  const id = newUser.rows[0].userid

  return {
    ...user,
    id
  }
}

const teardown = async (pool) => {
  const oldUser = await pool.query(
    `SELECT id FROM users WHERE email = $1;`,
    [user.email]
  )

  if (oldUser.rows.length) {
    await pool.query('call delete_user($1);', [oldUser.rows[0].id])
  }

  return {}
}

module.exports =
  {
    seed,
    teardown
  }
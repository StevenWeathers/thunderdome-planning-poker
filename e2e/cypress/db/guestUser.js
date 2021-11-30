const user = {
  name: 'Guest Test User'
}

const teardown = async (pool) => {
  const oldUser = await pool.query(
    `SELECT id FROM users WHERE name = $1;`,
    [user.name]
  )

  if (oldUser.rows.length) {
    await pool.query('call delete_user($1);', [oldUser.rows[0].id])
  }

  return {}
}

module.exports =
  {
    teardown
  }
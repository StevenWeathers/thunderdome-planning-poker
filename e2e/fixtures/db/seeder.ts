import crypto from "crypto";

export class ThunderdomeSeeder {
  constructor(pool) {
    this.pool = pool;
  }

  generateEmail(name) {
    const sanitizedName = name.toLowerCase().replace(/\s+/g, "");
    const randomString = crypto.randomBytes(4).toString("hex");
    return `${sanitizedName}.${randomString}@thunderometest.com`;
  }

  async createUser(name, email, type = "REGISTERED", verified = false) {
    const hashedPass =
      "$2a$10$3CvuzyoGIme3dJ4v9BnvyOIKFxEaYyjV2Lfunykv0VokGf/twxi9m"; // kentRules!
    const { rows } = await this.pool.query(
      `
      INSERT INTO thunderdome.users (name, email, type, verified)
      VALUES ($1, $2, $3, $4)
      RETURNING id
    `,
      [name, email, type, verified],
    );
    const id = rows[0].id;

    await this.pool.query(
      `INSERT INTO thunderdome.auth_credential (user_id, email, password, verified) VALUES ($1, $2, $3, $4);`,
      [id, email, hashedPass, verified],
    );

    return {
      id,
      name,
      email,
      type,
      verified,
    };
  }

  async addUserAPIKey(userId, apkId, name, active = true) {
    await this.pool.query(
      `
      INSERT INTO thunderdome.api_key (id, user_id, name, active) VALUES ($1, $2, $3, $4);
    `,
      [apkId, userId, name, active],
    );

    return {
      id: apkId,
      name,
      active,
    };
  }

  async deleteUserByName(name) {
    const { rows } = await this.pool.query(
      `
      DELETE FROM thunderdome.users
      WHERE name = $1
      RETURNING id
    `,
      [name],
    );
    return rows.length > 0 ? rows[0].id : null;
  }

  async deleteUserByEmail(email) {
    const { rows } = await this.pool.query(
      `
      DELETE FROM thunderdome.users
      WHERE email = $1
      RETURNING id
    `,
      [email],
    );
    return rows.length > 0 ? rows[0].id : null;
  }

  async createOrganization(name, ownerId) {
    const result = await this.pool.query(
      `
      INSERT INTO thunderdome.organization (name)
      VALUES ($1)
      RETURNING id
    `,
      [name],
    );
    const orgId = result.rows[0].id;
    await this.pool.query(
      `
      INSERT INTO thunderdome.organization_user (organization_id, user_id, role)
      VALUES ($1, $2, 'ADMIN')
    `,
      [orgId, ownerId],
    );
    return {
      id: orgId,
      name,
    };
  }

  async getOrganizationById(id) {
    const { rows } = await this.pool.query(
      `
      SELECT id, name, created_date, updated_date FROM thunderdome.organization WHERE id = $1;
    `,
      [id],
    );

    return rows.length > 0 ? rows[0] : null;
  }

  async createDepartment(name, orgId) {
    const result = await this.pool.query(
      `
      INSERT INTO thunderdome.organization_department (organization_id, name)
      VALUES ($1, $2)
      RETURNING id
    `,
      [orgId, name],
    );
    return {
      id: result.rows[0].id,
      organization_id: orgId,
      name,
    };
  }

  async getDepartmentById(id) {
    const { rows } = await this.pool.query(
      `
      SELECT id, organization_id, name, created_date, updated_date FROM thunderdome.organization_department WHERE id = $1;
    `,
      [id],
    );

    return rows.length > 0 ? rows[0] : null;
  }

  async createOrgTeam(name, orgId = null, deptId = null) {
    const result = await this.pool.query(
      `
      INSERT INTO thunderdome.team (name, organization_id, department_id)
      VALUES ($1, $2, $3)
      RETURNING id
    `,
      [name, orgId, deptId],
    );
    return {
      id: result.rows[0].id,
      name,
    };
  }

  async createUserTeam(name, userId) {
    const { rows } = await this.pool.query(
      `
      INSERT INTO thunderdome.team (name)
      VALUES ($1)
      RETURNING id
    `,
      [name],
    );
    const id = rows[0].id;

    await this.addUserToTeam(userId, id, "ADMIN");
    return {
      id,
      name,
    };
  }

  async addUserToOrg(userId, orgId, role = "MEMBER") {
    await this.pool.query(
      `
      INSERT INTO thunderdome.organization_user (organization_id, user_id, role)
      VALUES ($1, $2, $3)
    `,
      [orgId, userId, role],
    );
  }

  async addUserToDept(userId, deptId, role = "MEMBER") {
    await this.pool.query(
      `
      INSERT INTO thunderdome.department_user (department_id, user_id, role)
      VALUES ($1, $2, $3)
    `,
      [deptId, userId, role],
    );
  }

  async addUserToTeam(userId, teamId, role = "MEMBER") {
    await this.pool.query(
      `
      INSERT INTO thunderdome.team_user (team_id, user_id, role)
      VALUES ($1, $2, $3)
    `,
      [teamId, userId, role],
    );
  }

  async close() {
    await this.pool.end();
  }
}

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

  // Helper function to generate unique project keys with better uniqueness
  generateUniqueProjectKey(baseKey) {
    const timestamp = Date.now().toString();
    const random = Math.floor(Math.random() * 10000)
      .toString()
      .padStart(4, "0");
    const unique = `${timestamp.slice(-4)}${random}`.slice(0, 8); // 8 digit unique suffix
    const prefix = baseKey.substring(0, 2).toUpperCase(); // Take first 2 chars of base
    return `${prefix}${unique}`;
  }

  // Alternative: Check for existing keys and retry if needed
  async generateUniqueProjectKeyForScope(baseKey, scopeColumn, scopeId) {
    let attempts = 0;
    let uniqueKey;

    while (attempts < 10) {
      // Max 10 attempts
      uniqueKey = this.generateUniqueProjectKey(baseKey);

      // Check if key exists in this scope
      const { rows } = await this.pool.query(
        `SELECT id FROM thunderdome.project WHERE project_key = $1 AND ${scopeColumn} = $2`,
        [uniqueKey, scopeId],
      );

      if (rows.length === 0) {
        return uniqueKey; // Key is unique in this scope
      }

      attempts++;
      // Add a small delay to ensure different timestamps
      await new Promise((resolve) => setTimeout(resolve, 1));
    }

    // Fallback: use UUID suffix if all attempts failed
    const { v4: uuidv4 } = require("uuid");
    const uuid = uuidv4().replace(/-/g, "").substring(0, 8).toUpperCase();
    return `${baseKey.substring(0, 2)}${uuid}`;
  }

  // Create a project with organization association
  async createProject(projectKey, name, organizationId, description = null) {
    const uniqueKey = await this.generateUniqueProjectKeyForScope(
      projectKey,
      "organization_id",
      organizationId,
    );

    const { rows } = await this.pool.query(
      `
    INSERT INTO thunderdome.project (project_key, name, description, organization_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id, project_key, name, description, organization_id, created_at, updated_at
  `,
      [uniqueKey, name, description, organizationId],
    );

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      organizationId: rows[0].organization_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Create a project with organization association (explicit method)
  async createOrganizationProject(
    projectKey,
    name,
    organizationId,
    description = null,
  ) {
    const uniqueKey = await this.generateUniqueProjectKeyForScope(
      projectKey,
      "organization_id",
      organizationId,
    );

    const { rows } = await this.pool.query(
      `
    INSERT INTO thunderdome.project (project_key, name, description, organization_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id, project_key, name, description, organization_id, created_at, updated_at
  `,
      [uniqueKey, name, description, organizationId],
    );

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      organizationId: rows[0].organization_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Create a project with department association
  async createDepartmentProject(
    projectKey,
    name,
    departmentId,
    description = null,
  ) {
    const uniqueKey = await this.generateUniqueProjectKeyForScope(
      projectKey,
      "department_id",
      departmentId,
    );

    const { rows } = await this.pool.query(
      `
    INSERT INTO thunderdome.project (project_key, name, description, department_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id, project_key, name, description, department_id, created_at, updated_at
  `,
      [uniqueKey, name, description, departmentId],
    );

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      departmentId: rows[0].department_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Create a project with team association
  async createTeamProject(projectKey, name, teamId, description = null) {
    const uniqueKey = await this.generateUniqueProjectKeyForScope(
      projectKey,
      "team_id",
      teamId,
    );

    const { rows } = await this.pool.query(
      `
    INSERT INTO thunderdome.project (project_key, name, description, team_id)
    VALUES ($1, $2, $3, $4)
    RETURNING id, project_key, name, description, team_id, created_at, updated_at
  `,
      [uniqueKey, name, description, teamId],
    );

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      teamId: rows[0].team_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Helper function to generate unique department names within an organization
  async generateUniqueDepartmentName(baseName, organizationId) {
    let attempts = 0;
    let uniqueName;

    while (attempts < 10) {
      if (attempts === 0) {
        uniqueName = baseName; // Try original name first
      } else {
        const timestamp = Date.now().toString().slice(-4);
        const random = Math.floor(Math.random() * 100)
          .toString()
          .padStart(2, "0");
        uniqueName = `${baseName} ${timestamp}${random}`;
      }

      // Check if name exists in this organization
      const { rows } = await this.pool.query(
        `SELECT id FROM thunderdome.organization_department WHERE name = $1 AND organization_id = $2`,
        [uniqueName, organizationId],
      );

      if (rows.length === 0) {
        return uniqueName; // Name is unique in this organization
      }

      attempts++;
      // Add a small delay to ensure different timestamps
      await new Promise((resolve) => setTimeout(resolve, 1));
    }

    // Fallback: use UUID suffix if all attempts failed
    const { v4: uuidv4 } = require("uuid");
    const uuid = uuidv4().replace(/-/g, "").substring(0, 6).toUpperCase();
    return `${baseName} ${uuid}`;
  }

  // Create a project with multiple associations (for admin tests)
  async createProjectWithAssociations(
    projectKey,
    name,
    associations = {},
    description = null,
  ) {
    const { organizationId, departmentId, teamId } = associations;

    // Determine which scope to use for uniqueness check
    let scopeColumn, scopeId;
    if (teamId) {
      scopeColumn = "team_id";
      scopeId = teamId;
    } else if (departmentId) {
      scopeColumn = "department_id";
      scopeId = departmentId;
    } else if (organizationId) {
      scopeColumn = "organization_id";
      scopeId = organizationId;
    } else {
      // For admin tests without specific scope, use a simple unique key
      const uniqueKey = this.generateUniqueProjectKey(projectKey);

      const { rows } = await this.pool.query(
        `
      INSERT INTO thunderdome.project (project_key, name, description, organization_id, department_id, team_id)
      VALUES ($1, $2, $3, $4, $5, $6)
      RETURNING id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
    `,
        [
          uniqueKey,
          name,
          description,
          organizationId || null,
          departmentId || null,
          teamId || null,
        ],
      );

      return {
        id: rows[0].id,
        projectKey: rows[0].project_key,
        name: rows[0].name,
        description: rows[0].description,
        organizationId: rows[0].organization_id,
        departmentId: rows[0].department_id,
        teamId: rows[0].team_id,
        createdAt: rows[0].created_at,
        updatedAt: rows[0].updated_at,
      };
    }

    const uniqueKey = await this.generateUniqueProjectKeyForScope(
      projectKey,
      scopeColumn,
      scopeId,
    );

    const { rows } = await this.pool.query(
      `
    INSERT INTO thunderdome.project (project_key, name, description, organization_id, department_id, team_id)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
  `,
      [
        uniqueKey,
        name,
        description,
        organizationId || null,
        departmentId || null,
        teamId || null,
      ],
    );

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      organizationId: rows[0].organization_id,
      departmentId: rows[0].department_id,
      teamId: rows[0].team_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Delete a project (for cleanup in tests)
  async deleteProject(projectId) {
    await this.pool.query(
      `
    DELETE FROM thunderdome.project WHERE id = $1
  `,
      [projectId],
    );
  }

  // Get project by ID (for verification in tests)
  async getProjectById(projectId) {
    const { rows } = await this.pool.query(
      `
    SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
    FROM thunderdome.project
    WHERE id = $1
  `,
      [projectId],
    );

    if (rows.length === 0) {
      return null;
    }

    return {
      id: rows[0].id,
      projectKey: rows[0].project_key,
      name: rows[0].name,
      description: rows[0].description,
      organizationId: rows[0].organization_id,
      departmentId: rows[0].department_id,
      teamId: rows[0].team_id,
      createdAt: rows[0].created_at,
      updatedAt: rows[0].updated_at,
    };
  }

  // Get all projects for an organization (for verification)
  async getProjectsByOrganization(organizationId) {
    const { rows } = await this.pool.query(
      `
    SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
    FROM thunderdome.project
    WHERE organization_id = $1
    ORDER BY created_at DESC
  `,
      [organizationId],
    );

    return rows.map((row) => ({
      id: row.id,
      projectKey: row.project_key,
      name: row.name,
      description: row.description,
      organizationId: row.organization_id,
      departmentId: row.department_id,
      teamId: row.team_id,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }));
  }

  // Get all projects for a department (for verification)
  async getProjectsByDepartment(departmentId) {
    const { rows } = await this.pool.query(
      `
    SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
    FROM thunderdome.project
    WHERE department_id = $1
    ORDER BY created_at DESC
  `,
      [departmentId],
    );

    return rows.map((row) => ({
      id: row.id,
      projectKey: row.project_key,
      name: row.name,
      description: row.description,
      organizationId: row.organization_id,
      departmentId: row.department_id,
      teamId: row.team_id,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }));
  }

  // Get all projects for a team (for verification)
  async getProjectsByTeam(teamId) {
    const { rows } = await this.pool.query(
      `
    SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
    FROM thunderdome.project
    WHERE team_id = $1
    ORDER BY created_at DESC
  `,
      [teamId],
    );

    return rows.map((row) => ({
      id: row.id,
      projectKey: row.project_key,
      name: row.name,
      description: row.description,
      organizationId: row.organization_id,
      departmentId: row.department_id,
      teamId: row.team_id,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }));
  }

  async close() {
    await this.pool.end();
  }
}

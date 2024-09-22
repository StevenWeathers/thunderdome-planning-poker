import { ThunderdomeSeeder } from "./seeder";

const adminUser = {
  name: "E2E ADMIN",
  email: "e2eadmin@thunderdome.dev",
  type: "ADMIN",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    adminUser.name,
    adminUser.email,
    adminUser.type,
    true,
  );

  return {
    ...adminUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(adminUser.email);

  return {};
};

export const setupAdminUser = {
  seed,
  teardown,
};

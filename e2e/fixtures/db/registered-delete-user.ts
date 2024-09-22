import { ThunderdomeSeeder } from "./seeder";

export const registeredDeleteUser = {
  name: "E2E Delete User",
  email: "e2edelete@thunderdome.dev",
  type: "REGISTERED",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    registeredDeleteUser.name,
    registeredDeleteUser.email,
    registeredDeleteUser.type,
    false,
  );

  return {
    ...registeredDeleteUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(registeredDeleteUser.email);

  return {};
};

export const setupDeleteRegisteredUser = {
  seed,
  teardown,
};

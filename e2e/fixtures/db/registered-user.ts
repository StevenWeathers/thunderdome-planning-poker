import { ThunderdomeSeeder } from "./seeder";

export const registeredUser = {
  name: "E2E Registered User",
  email: "e2eregistered@thunderdome.dev",
  type: "REGISTERED",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    registeredUser.name,
    registeredUser.email,
    registeredUser.type,
    false,
  );

  return {
    ...registeredUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(registeredUser.email);

  return {};
};

export const setupRegisteredUser = {
  seed,
  teardown,
};

import { ThunderdomeSeeder } from "./seeder";

export const adminAPIUser = {
  name: "E2EAdminAPIUser",
  email: "e2eadminapi@thunderdome.dev",
  type: "ADMIN",
  apikey: "Gssy-ffy.okeTA-3AJhCnY1sqeUvRPRHiNYIVUxs4",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    adminAPIUser.name,
    adminAPIUser.email,
    adminAPIUser.type,
    true,
  );

  await seeder.addUserAPIKey(
    id,
    "Gssy-ffy.e170ffced2ae5806aebc103f30255dc5cc1b9e203d6035aa817f2b7e6638f223",
    "test api key 2",
  );

  return {
    ...adminAPIUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(adminAPIUser.email);

  return {};
};

export const setupAdminAPIUser = {
  seed,
  teardown,
};

import { ThunderdomeSeeder } from "./seeder";

export const apiUser = {
  name: "E2EAPIUser",
  email: "e2eapi@thunderdome.dev",
  type: "REGISTERED",
  apikey: "8MenPkY8.Vqvkh030vv7$rSyYs1gt++L0v7wKuVgR",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    apiUser.name,
    apiUser.email,
    apiUser.type,
    true,
  );

  await seeder.addUserAPIKey(
    id,
    "8MenPkY8.cd737cbc4bdca1838bdcf1685b00a9a778261255c10193714d9ba1630b55b63c",
    "test apikey",
  );

  return {
    ...apiUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(apiUser.email);

  return {};
};

export const setupAPIUser = {
  seed,
  teardown,
};

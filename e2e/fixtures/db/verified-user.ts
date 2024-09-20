import { ThunderdomeSeeder } from "./seeder";

const verifiedUser = {
  name: "E2E Verified User",
  email: "e2everified@thunderdome.dev",
  type: "REGISTERED",
};

const seed = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  const { id } = await seeder.createUser(
    verifiedUser.name,
    verifiedUser.email,
    verifiedUser.type,
    true,
  );

  return {
    ...verifiedUser,
    id,
  };
};

const teardown = async (pool) => {
  const seeder = new ThunderdomeSeeder(pool);
  await seeder.deleteUserByEmail(verifiedUser.email);

  return {};
};

export const setupVerifiedUser = {
  seed,
  teardown,
};

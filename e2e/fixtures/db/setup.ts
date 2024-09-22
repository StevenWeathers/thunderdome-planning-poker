import { Pool } from "pg";
import process from "process";

export const setupDB = function () {
  const dbConfig = {
    database: process.env.DB_NAME || "thunderdome",
    user: process.env.DB_USER || "thor",
    password: process.env.DB_PASS || "odinson",
    port: process.env.DB_PORT || "5432",
    host: process.env.DB_HOST || "localhost",
  };

  return new Pool(dbConfig);
};

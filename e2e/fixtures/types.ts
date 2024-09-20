export type TestUser = {
  id?: number;
  name?: string;
  email?: string;
  type?: string | "REGISTERED" | "GUEST" | "ADMIN";
  apikey?: string;
  verified?: boolean;
  orgs?: Array<any>;
  orgTeams?: Array<any>;
  depts?: Array<any>;
  deptTeams?: Array<any>;
  teams?: Array<any>;
};

export type TestUsers = {
  [key: string]: TestUser;
};

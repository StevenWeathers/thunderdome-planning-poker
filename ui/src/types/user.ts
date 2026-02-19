export type User = {
  avatar: string;
  company: string;
  country: string;
  createdDate: string;
  disabled: boolean;
  email: string;
  gravatarHash: string;
  picture?: string;
  id: string;
  jobTitle: string;
  lastActive: string;
  locale: string;
  mfaEnabled: boolean;
  name: string;
  notificationsEnabled: boolean;
  rank: string;
  updatedDate: string;
  verified: boolean;
  theme: string;
};

export type UserAPIKey = {
  active: boolean;
  apiKey: string;
  createdDate: string;
  id: string;
  name: string;
  prefix: string;
  updatedDate: string;
  userEmail: string;
  userId: string;
  userName: string;
};

export type SupportTicket = {
  id: string;
  userId: string | null;
  fullName: string;
  email: string;
  inquiry: string;
  assignedTo: string | null;
  resolvedAt: string | null;
  resolvedBy: string | null;
  notes: string | null;
  createdAt: string;
  updatedAt: string;
};

export type SessionUser = {
  avatar?: string;
  company?: string;
  country?: string;
  createdDate: string;
  disabled?: boolean;
  email?: string;
  gravatarHash?: string;
  picture?: string;
  id: string;
  jobTitle?: string;
  lastActive: string;
  locale: string;
  mfaEnabled?: boolean;
  name: string;
  notificationsEnabled?: boolean;
  rank: string;
  updatedDate: string;
  verified?: boolean;
  theme?: string;
  subscribed: boolean;
};

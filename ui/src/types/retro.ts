export type Retro = {
  actionItems: Array<RetroAction>;
  brainstormVisibility: string;
  createdDate: string;
  facilitatorCode: string;
  facilitators: Array<string>;
  format: string;
  groups: Array<RetroGroup>;
  id: string;
  items: Array<RetroItem>;
  joinCode: string;
  maxVotes: number;
  name: string;
  ownerId: string;
  phase: string;
  updatedDate: string;
  users: Array<RetroUser>;
  votes: Array<RetroVote>;
};

export type RetroAction = {
  comments: Array<RetroActionComment>;
  completed: boolean;
  content: string;
  assignees: Array<RetroUser>;
  id: string;
  retroId: string;
};

export type RetroActionComment = {
  comment: string;
  created_date: string;
  id: string;
  retro_id: string;
  updated_date: string;
  user_id: string;
};

export type RetroGroup = {
  id: string;
  name: string;
  votes?: Array<RetroVote>;
  items?: Array<RetroItem>;
};

export type RetroItem = {
  content: string;
  groupId: string;
  id: string;
  type: string;
  userId: string;
};

export type RetroUser = {
  active: boolean;
  avatar: string;
  gravatarHash: string;
  id: string;
  name: string;
};

export type RetroVote = {
  groupId: string;
  userId: string;
  count: number;
};

export type RetroTemplateColumn = {
  name: string;
  label: string;
  color: 'red' | 'blue' | 'green' | 'yellow' | 'purple' | 'teal' | 'orange' | '';
  icon: 'smiley' | 'frown' | 'question' | 'angry' | '';
};
export type RetroTemplateFormat = {
  columns: RetroTemplateColumn[];
};

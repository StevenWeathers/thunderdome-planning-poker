export type StoryboardColor = {
  color: string;
  legend: string;
};
export type Storyboard = {
  color_legend: Array<StoryboardColor>;
  createdDate?: Date;
  facilitatorCode: string;
  facilitators: Array<string>;
  goals: Array<StoryboardGoal>;
  id: string;
  joinCode?: string;
  name: string;
  owner_id: string;
  personas: Array<StoryboardPersona>;
  updatedDate?: Date;
  users: Array<StoryboardUser>;
  teamId?: string;
  teamName?: string;
};

export type StoryboardColumn = {
  id: string;
  name: string;
  personas: Array<StoryboardPersona>;
  sort_order: number;
  stories: Array<StoryboardStory>;
};

export type StoryboardGoal = {
  columns: Array<StoryboardColumn>;
  id: string;
  name: string;
  personas: Array<StoryboardPersona>;
  sort_order: number;
};

export type StoryboardPersona = {
  id: string;
  name: string;
  role?: string;
  description?: string;
};

export type StoryComment = {
  comment: string;
  created_date: string;
  id: string;
  story_id: string;
  updated_date: Date;
  user_id: string;
};

export type StoryboardStory = {
  annotations: Array<string>;
  closed: boolean;
  color: string;
  comments: Array<StoryComment>;
  content: string;
  id: string;
  link: string;
  name: string;
  points: number;
  sort_order: number;
};

export type StoryboardUser = {
  abandoned: boolean;
  active: boolean;
  avatar: string;
  gravatarHash: string;
  id: string;
  name: string;
};

export type ColorLegend = {
  color: string;
  legend: string;
};

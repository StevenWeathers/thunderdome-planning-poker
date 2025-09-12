export type PokerGame = {
  activePlanId?: string;
  autoFinishVoting: boolean;
  createdDate: Date;
  endedDate?: Date;
  hideVoterIdentity: boolean;
  id: string;
  joinCode?: string;
  leaderCode?: string;
  leaders: Array<String>;
  name: string;
  plans: Array<PokerStory>;
  pointAverageRounding: string;
  pointValuesAllowed: Array<string>;
  updatedDate: Date;
  users: Array<PokerUser>;
  votingLocked: boolean;
  teamId?: string;
};

export type PokerStory = {
  id: string;
  name: string;
  type: string;
  referenceId?: string;
  link?: string;
  description?: string;
  acceptanceCriteria?: string;
  active: boolean;
  points: string;
  priority: number;
  skipped: boolean;
  voteEndTime: Date;
  voteStartTime: Date;
  votes: Array<PokerStoryVote>;
  position: number;
};

export type PokerStoryVote = {
  vote: string;
  warriorId: string;
};

export type PokerUser = {
  abandoned: boolean;
  active: boolean;
  avatar: string;
  gravatarHash: string;
  id: string;
  name: string;
  rank: string;
  spectator: boolean;
};

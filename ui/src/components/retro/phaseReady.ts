export interface PhaseUser {
  id: string;
  active: boolean;
}

interface PhaseVote {
  userId: string;
  count: number;
}

interface PhaseState {
  phase: string;
  users: PhaseUser[];
  items: { groupId: string }[];
  votes: PhaseVote[];
  readyUsers: string[];
  maxVotes: number;
  allowCumulativeVoting: boolean;
}

// Never true without active users, otherwise an empty room would count as everyone being done
function everyActiveUser(users: PhaseUser[], isDone: (user: PhaseUser) => boolean): boolean {
  const activeUsers = users.filter(u => u.active);

  return activeUsers.length > 0 && activeUsers.every(isDone);
}

// Without cumulative voting each group takes one vote, so there may be fewer to cast than the limit
function votesAvailablePerUser({ items, maxVotes, allowCumulativeVoting }: PhaseState): number {
  const groupCount = new Set(items.map(item => item.groupId)).size;

  return allowCumulativeVoting && groupCount > 0 ? maxVotes : Math.min(maxVotes, groupCount);
}

function haveAllActiveUsersVoted(state: PhaseState): boolean {
  const required = votesAvailablePerUser(state);

  // votes are checked per user so those from users who have left don't count for anyone else
  return (
    required > 0 &&
    everyActiveUser(
      state.users,
      user => state.votes.reduce((total, v) => (v.userId === user.id ? total + v.count : total), 0) >= required,
    )
  );
}

/**
 * Whether every active user has finished the current phase, for phases that can auto advance.
 */
export function isPhaseComplete(state: PhaseState): boolean {
  switch (state.phase) {
    case 'brainstorm':
      return everyActiveUser(state.users, user => state.readyUsers.includes(user.id));
    case 'vote':
      return haveAllActiveUsersVoted(state);
    default:
      return false;
  }
}

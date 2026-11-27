import { describe, it, expect } from 'vitest';

import { isPhaseComplete } from './phaseReady';

const active = (id: string) => ({ id, active: true });
const inactive = (id: string) => ({ id, active: false });
const vote = (userId: string, count: number) => ({ userId, count });
const groups = (count: number) => Array.from({ length: count }, (_, i) => ({ groupId: `group-${i}` }));

const state = (overrides: Partial<Parameters<typeof isPhaseComplete>[0]>) => ({
  phase: 'vote',
  users: [active('a'), active('b')],
  items: groups(5),
  votes: [],
  readyUsers: [],
  maxVotes: 3,
  allowCumulativeVoting: false,
  ...overrides,
});

describe('isPhaseComplete vote phase', () => {
  it('is false when there are no active users', () => {
    expect(isPhaseComplete(state({ users: [inactive('a')] }))).toBe(false);
  });

  it('is false while an active user has votes left', () => {
    expect(isPhaseComplete(state({ votes: [vote('a', 3), vote('b', 2)] }))).toBe(false);
  });

  it('does not count votes from users who have left toward active users', () => {
    const users = [active('a'), active('b'), inactive('gone')];

    expect(isPhaseComplete(state({ users, votes: [vote('a', 3), vote('gone', 3)] }))).toBe(false);
  });

  it('is true once every active user has used all votes', () => {
    expect(isPhaseComplete(state({ votes: [vote('a', 2), vote('a', 1), vote('b', 3)] }))).toBe(true);
  });

  it('only requires one vote per group when there are fewer groups than votes', () => {
    expect(isPhaseComplete(state({ items: groups(2), votes: [vote('a', 2), vote('b', 2)] }))).toBe(true);
  });

  it('requires every vote with cumulative voting even when there are few groups', () => {
    const votes = [vote('a', 2), vote('b', 2)];

    expect(isPhaseComplete(state({ items: groups(2), votes, allowCumulativeVoting: true }))).toBe(false);
  });

  it('is false when there is nothing to vote on', () => {
    expect(isPhaseComplete(state({ items: [] }))).toBe(false);
    expect(isPhaseComplete(state({ items: [], allowCumulativeVoting: true }))).toBe(false);
  });
});

describe('isPhaseComplete brainstorm phase', () => {
  it('does not count a ready user who has left', () => {
    const users = [active('a'), active('b'), inactive('gone')];

    expect(isPhaseComplete(state({ phase: 'brainstorm', users, readyUsers: ['a', 'gone'] }))).toBe(false);
  });

  it('is true once every active user is ready', () => {
    expect(isPhaseComplete(state({ phase: 'brainstorm', readyUsers: ['a', 'b'] }))).toBe(true);
  });
});

describe('isPhaseComplete other phases', () => {
  it('never completes', () => {
    expect(isPhaseComplete(state({ phase: 'group', readyUsers: ['a', 'b'] }))).toBe(false);
  });
});

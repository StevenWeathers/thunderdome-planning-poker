import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { flushSync } from 'svelte';
import { render } from 'vitest-browser-svelte';

import { user } from '../../stores';
import Retro from './Retro.svelte';

const socket = vi.hoisted(() => ({
  onmessage: (_: { data: string }) => {},
  send: vi.fn(),
}));

vi.mock('../../lib/sockette', () => ({
  default: class {
    constructor(_: string, options: { onmessage: (e: { data: string }) => void }) {
      socket.onmessage = options.onmessage;
    }
    send = socket.send;
    close() {}
  },
}));

const users = (...inactive: string[]) =>
  ['facilitator', 'voter'].map(id => ({ id, name: id, avatar: '', active: !inactive.includes(id) }));

const votes = (...userIds: string[]) => userIds.map(userId => ({ userId, groupId: 'group-1', count: 1 }));

const receive = (type: string, value: unknown, userId = '') => {
  socket.onmessage({ data: JSON.stringify({ type, value: JSON.stringify(value), userId }) });
  flushSync();
};

const retro = (overrides = {}) => ({
  id: 'retro-1',
  name: 'Retro',
  phase: 'vote',
  phase_auto_advance: true,
  phase_time_limit_min: 0,
  phase_time_start: new Date().toISOString(),
  users: users(),
  facilitators: ['facilitator'],
  groups: [{ id: 'group-1', name: 'Deploys' }],
  items: [{ id: 'item-1', groupId: 'group-1', content: 'Ship it', type: 'worked', comments: [], reactions: [] }],
  votes: [],
  actionItems: [],
  readyUsers: [],
  maxVotes: 1,
  allowCumulativeVoting: false,
  template: { format: { columns: [{ name: 'worked', label: 'Worked', color: 'green' }] } },
  ...overrides,
});

const autoAdvances = () => socket.send.mock.calls.filter(([e]) => JSON.parse(e).type === 'phase_all_ready').length;

const setup = (overrides = {}) => {
  render(Retro, {
    retroId: 'retro-1',
    notifications: { success: vi.fn(), danger: vi.fn(), warning: vi.fn() } as any,
    router: { route: vi.fn() },
    xfetch: vi.fn() as any,
  });
  receive('init', retro(overrides));
};

describe('Retro phase auto advance', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    socket.send.mockClear();
    user.create({ id: 'facilitator', name: 'facilitator', rank: 'REGISTERED', locale: 'en' } as any);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('advances when the last vote is cast', () => {
    setup({ votes: votes('facilitator') });

    receive('votes_updated', votes('facilitator', 'voter'));

    expect(autoAdvances()).toBe(1);
  });

  it('does not advance on page load, even when everyone has voted', () => {
    setup({ votes: votes('facilitator', 'voter') });

    expect(autoAdvances()).toBe(0);
  });

  it('does not advance when going back to a phase everyone already finished', () => {
    setup({ phase: 'action', votes: votes('facilitator', 'voter') });

    receive('phase_updated', retro({ votes: votes('facilitator', 'voter') }));
    receive('user_left', users('voter'), 'voter');
    vi.advanceTimersByTime(10000);

    expect(autoAdvances()).toBe(0);
  });

  it('advances once the user being waited on has been gone for the grace period', () => {
    setup({ votes: votes('facilitator') });

    receive('user_left', users('voter'), 'voter');
    vi.advanceTimersByTime(9999);
    expect(autoAdvances()).toBe(0);

    vi.advanceTimersByTime(1);
    expect(autoAdvances()).toBe(1);
  });

  it('does not advance when the user being waited on rejoins in time', () => {
    setup({ votes: votes('facilitator') });

    receive('user_left', users('voter'), 'voter');
    receive('user_joined', users(), 'voter');
    vi.advanceTimersByTime(10000);

    expect(autoAdvances()).toBe(0);
  });

  it('keeps waiting for a refreshing user when someone else votes', () => {
    setup();

    receive('user_left', users('voter'), 'voter');
    receive('votes_updated', votes('facilitator'));

    expect(autoAdvances()).toBe(0);
  });
});

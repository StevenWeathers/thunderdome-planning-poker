import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { createPageLoader, reloadOnce } from '../pageLoader';

const First = {};
const Second = {};

function deferred() {
  let resolve!: (value: { default: any }) => void;
  let reject!: (reason: Error) => void;
  const promise = new Promise<{ default: any }>((res, rej) => {
    resolve = res;
    reject = rej;
  });

  return { load: () => promise, resolve, reject };
}

function setup() {
  const options = {
    onPage: vi.fn(),
    onPending: vi.fn(),
    onError: vi.fn(),
    reload: vi.fn(),
  };

  return { ...options, loadPage: createPageLoader(options) };
}

describe('pageLoader', () => {
  beforeEach(() => {
    sessionStorage.clear();
    vi.spyOn(console, 'error').mockImplementation(() => {});
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('provides the loaded page with its params and name', async () => {
    const { loadPage, onPage } = setup();

    await loadPage(async () => ({ default: First }), { teamId: '1' }, 'team');

    expect(onPage).toHaveBeenCalledWith({ route: First, params: { teamId: '1' }, name: 'team' });
  });

  it('ignores a slow page that finishes after a newer navigation', async () => {
    const { loadPage, onPage } = setup();
    const slow = deferred();
    const fast = deferred();

    const slowLoad = loadPage(slow.load, {}, 'first');
    const fastLoad = loadPage(fast.load, {}, 'second');
    fast.resolve({ default: Second });
    await fastLoad;
    slow.resolve({ default: First });
    await slowLoad;

    expect(onPage).toHaveBeenCalledTimes(1);
    expect(onPage).toHaveBeenCalledWith({ route: Second, params: {}, name: 'second' });
  });

  it('ignores the failure of a superseded page', async () => {
    const { loadPage, onError, reload } = setup();
    const slow = deferred();

    const slowLoad = loadPage(slow.load, {}, 'first');
    await loadPage(async () => ({ default: Second }), {}, 'second');
    slow.reject(new Error('gone'));
    await slowLoad;

    expect(reload).not.toHaveBeenCalled();
    expect(onError).not.toHaveBeenCalledWith(true);
  });

  it('reports pending only when loading is slow', async () => {
    vi.useFakeTimers();
    const { loadPage, onPending } = setup();
    const slow = deferred();

    const load = loadPage(slow.load, {}, 'first');
    vi.advanceTimersByTime(199);
    expect(onPending).not.toHaveBeenCalled();

    vi.advanceTimersByTime(1);
    expect(onPending).toHaveBeenLastCalledWith(true);

    slow.resolve({ default: First });
    await load;
    expect(onPending).toHaveBeenLastCalledWith(false);
  });

  it('reloads once for a failed page, then reports an error', async () => {
    const { loadPage, onError, reload } = setup();
    const fail = () => Promise.reject(new Error('missing chunk'));

    await loadPage(fail, {}, 'first');
    expect(reload).toHaveBeenCalledTimes(1);
    expect(onError).not.toHaveBeenCalledWith(true);

    await loadPage(fail, {}, 'first');
    expect(reload).toHaveBeenCalledTimes(1);
    expect(onError).toHaveBeenLastCalledWith(true);
  });

  it('allows another reload after a page loads successfully', async () => {
    const { loadPage, reload } = setup();
    const fail = () => Promise.reject(new Error('missing chunk'));

    await loadPage(fail, {}, 'first');
    await loadPage(async () => ({ default: Second }), {}, 'second');
    await loadPage(fail, {}, 'first');

    expect(reload).toHaveBeenCalledTimes(2);
  });

  it('does not reload when the guard cannot be stored', () => {
    vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
      throw new Error('blocked');
    });
    const reload = vi.fn();

    expect(reloadOnce('key', reload)).toBe(false);
    expect(reload).not.toHaveBeenCalled();
  });
});

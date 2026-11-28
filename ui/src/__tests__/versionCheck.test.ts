import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { reloadForNewVersion, startVersionCheck } from '../versionCheck';

function setup(version = 'v1') {
  const xfetch = vi.fn(async () => new Response(JSON.stringify({ data: { version } })));
  const onNewVersion = vi.fn();
  const stop = startVersionCheck({ xfetch, currentVersion: 'v1', onNewVersion, intervalMs: 1000 });

  return { xfetch, onNewVersion, stop };
}

describe('versionCheck', () => {
  beforeEach(() => {
    sessionStorage.clear();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('reports a different server version', async () => {
    const { xfetch, onNewVersion, stop } = setup('v2');

    await vi.advanceTimersByTimeAsync(1000);

    expect(xfetch).toHaveBeenCalledWith('/api/version');
    await vi.waitFor(() => expect(onNewVersion).toHaveBeenCalledWith('v2'));
    stop();
  });

  it('stays quiet when the version matches', async () => {
    const { onNewVersion, stop } = setup('v1');

    await vi.advanceTimersByTimeAsync(1000);

    expect(onNewVersion).not.toHaveBeenCalled();
    stop();
  });

  it('checks when the tab becomes visible, but not while hidden', async () => {
    const visibility = vi.spyOn(document, 'visibilityState', 'get').mockReturnValue('hidden');
    const { xfetch, stop } = setup();

    document.dispatchEvent(new Event('visibilitychange'));
    await vi.advanceTimersByTimeAsync(1000);
    expect(xfetch).not.toHaveBeenCalled();

    visibility.mockReturnValue('visible');
    document.dispatchEvent(new Event('visibilitychange'));
    expect(xfetch).toHaveBeenCalledTimes(1);
    stop();
  });

  it('keeps checking after a failed request', async () => {
    const { xfetch, onNewVersion, stop } = setup('v2');
    xfetch.mockRejectedValueOnce([new Error('Bad Gateway'), { status: 502 }]);

    await vi.advanceTimersByTimeAsync(2000);

    await vi.waitFor(() => expect(onNewVersion).toHaveBeenCalledTimes(1));
    stop();
  });

  it('stops checking', async () => {
    const { xfetch, stop } = setup();

    stop();
    document.dispatchEvent(new Event('visibilitychange'));
    await vi.advanceTimersByTimeAsync(5000);

    expect(xfetch).not.toHaveBeenCalled();
  });

  it('reloads at most once for each new version', () => {
    const reload = vi.fn();

    expect(reloadForNewVersion('', reload)).toBe(false);
    expect(reloadForNewVersion('v2', reload)).toBe(true);
    expect(reloadForNewVersion('v2', reload)).toBe(false);
    expect(reloadForNewVersion('v3', reload)).toBe(true);
    expect(reload).toHaveBeenCalledTimes(2);
  });
});

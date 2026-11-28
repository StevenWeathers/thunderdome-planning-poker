import { reloadOnce } from './pageLoader';
import type { ApiClient } from './types/apiclient';

export interface VersionCheckOptions {
  xfetch: ApiClient;
  currentVersion: string;
  onNewVersion: (version: string) => void;
  intervalMs?: number;
}

/**
 * Asks the server for its version while the tab is visible, and again whenever the tab becomes visible.
 * Returns a function that stops checking.
 */
export function startVersionCheck({
  xfetch,
  currentVersion,
  onNewVersion,
  intervalMs = 5 * 60 * 1000,
}: VersionCheckOptions): () => void {
  async function check() {
    if (document.visibilityState !== 'visible') return;

    try {
      const response = await xfetch('/api/version');
      const { data } = await response.json();
      if (data?.version && data.version !== currentVersion) {
        onNewVersion(data.version);
      }
    } catch {
      // offline or mid deploy, the next check will try again
    }
  }

  const timer = setInterval(check, intervalMs);
  document.addEventListener('visibilitychange', check);

  return () => {
    clearInterval(timer);
    document.removeEventListener('visibilitychange', check);
  };
}

/**
 * Reloads to pick up the given version, at most once per version so that servers running
 * mixed versions during a rolling deploy can't cause a reload loop.
 */
export function reloadForNewVersion(version: string, reload?: () => void): boolean {
  return version !== '' && reloadOnce(`versionReload:${version}`, reload);
}

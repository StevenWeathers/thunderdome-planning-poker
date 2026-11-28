export interface Page {
  route: any;
  params: any;
  name: string;
}

export type ModuleLoader = () => Promise<{ default: any }>;

export interface PageLoaderOptions {
  onPage: (page: Page) => void;
  // true while a page has been loading long enough to be worth showing a loader
  onPending: (pending: boolean) => void;
  onError: (failed: boolean) => void;
  reload?: () => void;
  pendingDelayMs?: number;
}

const reloadGuardKey = 'pageLoadReloaded';

/**
 * Reloads the app unless that was already tried for the given key, which prevents reload loops.
 * Returns whether a reload was started.
 */
export function reloadOnce(key: string, reload: () => void = () => window.location.reload()): boolean {
  try {
    if (sessionStorage.getItem(key)) {
      return false;
    }
    sessionStorage.setItem(key, 'true');
  } catch {
    // without storage there is no way to stop a loop, so don't reload
    return false;
  }

  reload();
  return true;
}

function clearReloadGuard(key: string) {
  try {
    sessionStorage.removeItem(key);
  } catch {
    // nothing was stored
  }
}

export function createPageLoader({ onPage, onPending, onError, reload, pendingDelayMs = 200 }: PageLoaderOptions) {
  let latestRequest = 0;
  let pendingTimer: ReturnType<typeof setTimeout> | undefined;

  function settle() {
    clearTimeout(pendingTimer);
    onPending(false);
  }

  return async function loadPage(load: ModuleLoader, params: any, name: string): Promise<void> {
    const request = ++latestRequest;

    clearTimeout(pendingTimer);
    onError(false);
    pendingTimer = setTimeout(() => onPending(true), pendingDelayMs);

    try {
      const module = await load();
      // a newer navigation happened while this page was loading
      if (request !== latestRequest) return;

      clearReloadGuard(reloadGuardKey);
      settle();
      onPage({ route: module.default, params, name });
    } catch (error) {
      if (request !== latestRequest) return;

      console.error(error);
      // most often a deploy replaced the hashed files this tab knows about
      if (reloadOnce(reloadGuardKey, reload)) return;

      settle();
      onError(true);
    }
  };
}

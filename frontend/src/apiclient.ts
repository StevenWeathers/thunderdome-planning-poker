import { PathPrefix } from './config';

/**
 * Extends fetch with common inputs e.g. credentials, content-type
 * and checks response status/ok for common errors
 * @param {function} handle401 401 handler, e.g. redirect to login
 */
export default function (handle401) {
  /**
   * Wrapper around fetch
   * @param {string} endpoint the endpoint to fetch
   * @param {object} config the optional fetch config e.g. body for post
   */
  return function (endpoint, customConfig: any = {}) {
    const headers = { 'content-type': 'application/json' };

    const config: RequestInit = {
      method: customConfig.body ? 'POST' : 'GET',
      credentials: 'same-origin',
      ...customConfig,
      headers: {
        ...headers,
        ...customConfig.headers,
      },
    };

    if (customConfig.body) {
      config.body = JSON.stringify(config.body);
    }

    return fetch(`${PathPrefix}${endpoint}`, config).then(response => {
      if (response.status === 401) {
        handle401(customConfig.skip401Redirect);
      }

      if (!response.ok) {
        throw [Error(response.statusText), response];
      }

      return response;
    });
  };
}

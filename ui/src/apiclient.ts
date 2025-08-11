import { PathPrefix } from './config';
import type {
  ApiClient,
  ApiClientConfig,
  ApiClientFactory,
} from './types/apiclient';

/**
 * Extends fetch with common inputs e.g. credentials, content-type
 * and checks response status/ok for common errors
 * @param {function} handle401 401 handler, e.g. redirect to login
 */
const apiclient: ApiClientFactory = handle401 => {
  /**
   * Wrapper around fetch
   * @param {string} endpoint the endpoint to fetch
   * @param {object} config the optional fetch config e.g. body for post
   */
  return function (
    endpoint: string,
    customConfig: ApiClientConfig = {},
  ): Promise<Response> {
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
};

export default apiclient;

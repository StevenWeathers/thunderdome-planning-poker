export interface ApiClientConfig {
  body?: any;
  headers?: Record<string, string>;
  method?: string;
  credentials?: RequestCredentials;
  skip401Redirect?: boolean;
  [key: string]: any;
}

export interface ApiClient {
  (endpoint: string, customConfig?: ApiClientConfig): Promise<Response>;
}

export type ApiClientFactory = (handle401: (skipRedirect?: boolean) => void) => ApiClient;

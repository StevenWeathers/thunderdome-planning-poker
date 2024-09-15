import { AppConfig, PathPrefix } from './config';

export const getWebsocketAddress = () => {
  const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws';
  const hostname = window.location.host;
  const socketDomain =
    AppConfig.WebsocketSubdomain !== ''
      ? `${AppConfig.WebsocketSubdomain}.${hostname}`
      : hostname;

  return `${socketExtension}://${socketDomain}${PathPrefix}`;
};

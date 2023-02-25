const URL = '/api/v1/';

const getWSURL = () => {
  const loc = window.location; let
    new_uri;
  if (loc.protocol === 'https:') {
    new_uri = 'wss:';
  } else {
    new_uri = 'ws:';
  }
  new_uri += `//${loc.host}${URL}`;
  return new_uri;
};
let apiUrl = URL;
let wsUrl = getWSURL();
let base = '';
if (process.env.NODE_ENV === 'development') {
  base = 'http://localhost:8080';
  apiUrl = base + URL;
  wsUrl = `ws://localhost:8080${URL}`;
}

export const BASE_URL = base;
export const API_URL = apiUrl;
export const WS_URL = wsUrl;

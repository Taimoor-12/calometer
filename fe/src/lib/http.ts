export const http_common = async (method: string, url: string, body?: Record<string, any>) => {
  const opts: RequestInit = {
      method,
      mode: 'cors',
      credentials: 'include',
      redirect: 'follow',
      cache: 'no-cache',
  };
  if (body) {
      opts.body = JSON.stringify(body);
  }
  const resp = await fetch(url, opts);
  const retResp = {
      httpCode: resp.status,
      httpHeaders: resp.headers,
      httpRedirected: resp.redirected,
  };
  try {
      const data = await resp.json();

      return {
          ...data,
          ...retResp,
      };
  } catch (err) {
      console.error('http_common: not json:', err);
      return retResp;
  }
};

export const http_get = async (url: string) => {
  return http_common('GET', url);
};

export const http_post = async (url: string, body: Record<string, any>) => {
  return http_common('POST', url, body);
};

export const http_put = async (url: string, body: Record<string, any>) => {
  return http_common('PUT', url, body);
};

export const http_delete = async (url: string, body: Record<string, any>) => {
  return http_common('DELETE', url, body);
};

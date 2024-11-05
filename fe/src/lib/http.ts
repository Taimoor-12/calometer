type HttpResponse = {
  httpCode: number;
  httpHeaders: Headers;
  httpRedirected: boolean;
}

type HttpData<T = any> = {
  code: Record<number, string>;
  data?: T
}

type HttpResponseWithData = HttpResponse & HttpData

export const http_common = async (
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  url: string,
  body?: Record<string, any>
): Promise<HttpResponseWithData> => {
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
  try {
      const data = await resp.json();
      const retResp: HttpResponse = {
        httpCode: resp.status,
        httpHeaders: resp.headers,
        httpRedirected: resp.redirected,
      };
      const respData: HttpData = data

      const httpResponseWithData: HttpResponseWithData = {
        ...retResp,
        ...respData,
      };
      
      return httpResponseWithData;
  } catch (err) {
      console.error('http_common: not json:', err);
      const retResp: HttpResponseWithData = {
        httpCode: resp.status,
        httpHeaders: resp.headers,
        httpRedirected: resp.redirected,
        code: {},
      };
      return retResp;
  }
};

export const http_get = async (url: string): Promise<HttpResponseWithData> => {
  return http_common('GET', url);
};

export const http_post = async (url: string, body: Record<string, any>): Promise<HttpResponseWithData> => {
  return http_common('POST', url, body);
};

export const http_put = async (url: string, body: Record<string, any>): Promise<HttpResponseWithData> => {
  return http_common('PUT', url, body);
};

export const http_delete = async (url: string, body: Record<string, any>): Promise<HttpResponseWithData> => {
  return http_common('DELETE', url, body);
};

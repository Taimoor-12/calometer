type RespData = {
  code: Record<number, string>,
  data?: unknown
}

type HttpInfo = {
  httpCode: number,
  httpHeaders: Headers,
  httpRedirected: boolean,
}

type RespDataWithHttpInfo = RespData & HttpInfo

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
  const retResp: HttpInfo = {
      httpCode: resp.status,
      httpHeaders: resp.headers,
      httpRedirected: resp.redirected,
  };
  try {
      const data: RespData = await resp.json();
      console.log(data)

      return {
          ...data,
          ...retResp,
      } as RespDataWithHttpInfo;
  } catch (err) {
      console.error('http_common: not json:', err);
      return retResp;
  }
};

export const http_get = async (url: string): Promise<RespDataWithHttpInfo | HttpInfo> => {
  return http_common('GET', url);
};

export const http_post = async (url: string, body: Record<string, any>): Promise<RespDataWithHttpInfo | HttpInfo>  => {
  return http_common('POST', url, body);
};

export const http_put = async (url: string, body: Record<string, any>): Promise<RespDataWithHttpInfo | HttpInfo> => {
  return http_common('PUT', url, body);
};

export const http_delete = async (url: string, body: Record<string, any>): Promise<RespDataWithHttpInfo | HttpInfo> => {
  return http_common('DELETE', url, body);
};

// Type guard function to check if resp is of type RespDataWithHttpInfo
export const isRespDataWithHttpInfo = (resp: RespDataWithHttpInfo | HttpInfo): resp is RespDataWithHttpInfo => {
  return (resp as RespDataWithHttpInfo).code !== undefined;
}
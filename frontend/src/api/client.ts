const API_BASE = '/api/v1';

export class ApiError extends Error {
  code: string;
  constructor(code: string, message: string) {
    super(message);
    this.code = code;
    this.name = 'ApiError';
  }
}

function getAccessToken(): string | null {
  return localStorage.getItem('access_token');
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
  auth = false
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (auth) {
    const token = getAccessToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    let errorData: { code?: string; message?: string } = {};
    try {
      errorData = await res.json();
    } catch {
      // ignore
    }
    throw new ApiError(
      errorData.code || 'UNKNOWN_ERROR',
      errorData.message || `HTTP ${res.status}`
    );
  }

  // some endpoints may return empty body
  const text = await res.text();
  if (!text) return {} as T;
  return JSON.parse(text);
}

export const api = {
  getCategories: () => request<import('../types').Category[]>('/categories'),

  getProducts: (params: { page?: number; category?: string } = {}) => {
    const search = new URLSearchParams();
    if (params.page) search.set('page', String(params.page));
    if (params.category) search.set('category', params.category);
    const query = search.toString();
    return request<import('../types').ProductsResponse>(
      `/products${query ? `?${query}` : ''}`
    );
  },

  // Auth
  login: (dto: import('../types').LoginDto) =>
    request<import('../types').AuthTokens>('/auth', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),

  register: (dto: import('../types').RegisterDto) =>
    request<{ message?: string } | Record<string, unknown>>('/users', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),

  refresh: () =>
    request<import('../types').AuthTokens>('/refresh', {
      method: 'POST',
    }),

  logout: () =>
    request<{ message: string }>('/logout', {
      method: 'POST',
    }, true),
};

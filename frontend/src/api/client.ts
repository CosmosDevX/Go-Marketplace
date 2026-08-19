const API_BASE = '/api/v1';

export class ApiError extends Error {
  constructor(message: string) {
    super(message);
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
    let message = `HTTP ${res.status}`;
    try {
      const errorData = await res.json();
      if (errorData?.message) message = errorData.message;
    } catch {
      // ignore
    }
    throw new ApiError(message);
  }

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
    request<Record<string, unknown>>('/users', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),

  refresh: () =>
    request<import('../types').AuthTokens>('/refresh', {
      method: 'POST',
    }),

  logout: () =>
    request<{ message: string }>(
      '/logout',
      { method: 'POST' },
      true
    ),

  // Cart (all require auth)
  getCart: () =>
    request<import('../types').CartItem[]>('/cart', {}, true),

  addToCart: (productId: number) =>
    request<{ cart_item_id: number }>(
      '/cart/items',
      {
        method: 'POST',
        body: JSON.stringify({ product_id: productId }),
      },
      true
    ),

  updateCartItem: (cartItemId: number, delta: 1 | -1) =>
    request<{ quantity: number }>(
      `/cart/items/${cartItemId}?delta=${delta}`,
      { method: 'PATCH' },
      true
    ),

  removeCartItem: (cartItemId: number) =>
    request<{ message: string }>(
      `/cart/items/${cartItemId}`,
      { method: 'DELETE' },
      true
    ),
};

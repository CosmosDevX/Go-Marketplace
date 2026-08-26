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
    ...(options.headers as Record<string, string>),
  };

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = headers['Content-Type'] || 'application/json';
  }

  if (auth) {
    const token = getAccessToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include',
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

  createCategory: (dto: { category_name: string; category_slug: string }) =>
    request<Record<string, unknown>>(
      '/categories',
      {
        method: 'POST',
        body: JSON.stringify(dto),
      },
      true
    ),

  getProducts: (
    params: {
      page?: number;
      category?: string;
      sortBy?: 'price' | 'name';
      asc?: boolean;
      search?: string;
    } = {}
  ) => {
    const search = new URLSearchParams();
    if (params.page) search.set('page', String(params.page));
    if (params.category) search.set('category', params.category);
    if (params.sortBy) {
      search.set('sortBy', params.sortBy);
      search.set('asc', String(params.asc ?? true));
    }
    if (params.search?.trim()) search.set('search', params.search.trim());
    const query = search.toString();
    return request<import('../types').ProductsResponse>(
      `/products${query ? `?${query}` : ''}`
    );
  },

  getSellerProducts: (page = 1) =>
    request<import('../types').ProductsResponse>(
      `/seller/products?page=${page}`,
      {},
      true
    ),

  createProduct: (formData: FormData) =>
    request<{ product_id: number }>(
      '/seller/products',
      { method: 'POST', body: formData },
      true
    ),

  updateProduct: (productId: number, formData: FormData) =>
    request<{ product_id?: number; message?: string }>(
      `/seller/products/${productId}`,
      { method: 'PUT', body: formData },
      true
    ),

  deleteProduct: (productId: number) =>
    request<{ message: string }>(
      `/seller/products/${productId}`,
      { method: 'DELETE' },
      true
    ),

  login: (dto: import('../types').LoginDto) =>
    request<import('../types').AuthTokens>('/login', {
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
    request<{ message: string }>('/logout', { method: 'POST' }, true),


  getOrders: () =>
    request<import('../types').Order[]>('/orders', {}, true),

  createOrder: () =>
    request<{ order_id: number }>('/orders', { method: 'POST' }, true),

  // Roles (admin)
  grantRole: (dto: { username: string; role: string }) =>
    request<{ message: string }>(
      '/roles',
      { method: 'POST', body: JSON.stringify(dto) },
      true
    ),

  revokeRole: (username: string, role: string) =>
    request<{ message: string }>(
      `/roles/user/${encodeURIComponent(username)}/role/${encodeURIComponent(role)}`,
      { method: 'DELETE' },
      true
    ),

  getUserRoles: (username: string) =>
    request<string[]>(`/roles/${encodeURIComponent(username)}`, {}, true),

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

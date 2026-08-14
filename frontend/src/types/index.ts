export interface Category {
  category_id: number;
  category_name: string;
  category_slug: string;
}

export interface Product {
  product_id: number;
  product_name: string;
  product_description: string;
  product_price: string;
  product_image: string | null;
  category: Category;
}

export interface ProductsResponse {
  limit: number;
  page: number;
  products: Product[];
}

export interface ApiError {
  code: string;
  message: string;
}

export interface AuthTokens {
  access_token: string;
}

export interface LoginDto {
  username: string;
  password: string;
}

export interface RegisterDto {
  username: string;
  password: string;
  email: string;
}

export interface UserInfo {
  username: string;
}

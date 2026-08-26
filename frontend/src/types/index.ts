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
  seller_id?: number;
}

export interface ProductsResponse {
  limit: number;
  page: number;
  products: Product[];
}

export interface CartItem {
  cart_item_id: number;
  cart_id: number;
  product: Product;
  quantity: number;
}

export interface OrderItem {
  order_item_id: number;
  product: Product;
  order_item_quantity: number;
  order_item_total: string;
}

export interface Order {
  order_id: number;
  order_status: string;
  user_id: number;
  order_total: string;
  order_items: OrderItem[];
}

export interface AuthTokens {
  access_token: string;
  roles: string[];
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
  roles: string[];
}

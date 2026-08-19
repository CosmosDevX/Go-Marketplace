import { useCallback, useEffect, useState } from 'react';
import { api, ApiError } from '../api/client';
import type { CartItem } from '../types';

export function useCart(isAuthenticated: boolean) {
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadCart = useCallback(async () => {
    if (!isAuthenticated) {
      setItems([]);
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const data = await api.getCart();
      setItems(Array.isArray(data) ? data : []);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Не удалось загрузить корзину');
      setItems([]);
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    loadCart();
  }, [loadCart]);

  const addItem = useCallback(
    async (productId: number) => {
      await api.addToCart(productId);
      await loadCart();
    },
    [loadCart]
  );

  const changeQuantity = useCallback(
    async (cartItemId: number, delta: 1 | -1) => {
      const res = await api.updateCartItem(cartItemId, delta);
      if (res.quantity === 0) {
        setItems((prev) => prev.filter((i) => i.cart_item_id !== cartItemId));
      } else {
        setItems((prev) =>
          prev.map((i) =>
            i.cart_item_id === cartItemId ? { ...i, quantity: res.quantity } : i
          )
        );
      }
    },
    []
  );

  const removeItem = useCallback(async (cartItemId: number) => {
    await api.removeCartItem(cartItemId);
    setItems((prev) => prev.filter((i) => i.cart_item_id !== cartItemId));
  }, []);

  const totalCount = items.reduce((sum, i) => sum + i.quantity, 0);

  const totalPrice = items.reduce((sum, i) => {
    const price = parseFloat(i.product.product_price) || 0;
    return sum + price * i.quantity;
  }, 0);

  return {
    items,
    loading,
    error,
    totalCount,
    totalPrice,
    loadCart,
    addItem,
    changeQuantity,
    removeItem,
  };
}

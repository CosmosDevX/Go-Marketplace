import { useState, useEffect } from 'react';
import { X, Plus, Minus, Trash2, ShoppingBag, Loader2 } from 'lucide-react';
import type { CartItem } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';
import { IMAGE_CONFIG } from '../config/images';
import { ApiError } from '../api/client';

interface Props {
  open: boolean;
  onClose: () => void;
  items: CartItem[];
  loading: boolean;
  error: string | null;
  totalPrice: number;
  onChangeQuantity: (cartItemId: number, delta: 1 | -1) => Promise<void>;
  onRemove: (cartItemId: number) => Promise<void>;
  onCheckout: () => Promise<void>;
}

export function CartDrawer({
  open,
  onClose,
  items,
  loading,
  error,
  totalPrice,
  onChangeQuantity,
  onRemove,
  onCheckout,
}: Props) {
  const [busyId, setBusyId] = useState<number | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [checkoutLoading, setCheckoutLoading] = useState(false);
  const [checkoutSuccess, setCheckoutSuccess] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      setCheckoutSuccess(null);
      setActionError(null);
    }
  }, [open]);

  useEffect(() => {
    if (items.length === 0) {
      // после успешного заказа корзина пустая — сообщение можно сбросить при следующем открытии
    }
  }, [items.length]);

  if (!open) return null;

  const handleDelta = async (id: number, delta: 1 | -1) => {
    setBusyId(id);
    setActionError(null);
    try {
      await onChangeQuantity(id, delta);
    } catch (e) {
      setActionError(e instanceof ApiError ? e.message : 'Ошибка обновления');
    } finally {
      setBusyId(null);
    }
  };

  const handleRemove = async (id: number) => {
    setBusyId(id);
    setActionError(null);
    try {
      await onRemove(id);
    } catch (e) {
      setActionError(e instanceof ApiError ? e.message : 'Ошибка удаления');
    } finally {
      setBusyId(null);
    }
  };

  return (
    <div className="fixed inset-0 z-[100] flex justify-end">
      <div className="absolute inset-0 bg-black/60" onClick={onClose} />

      <div className="relative flex h-full w-full max-w-md flex-col border-l border-[var(--color-border)] bg-[var(--color-background)] shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-[var(--color-border)] px-4 py-4">
          <h2 className="text-lg font-semibold">Корзина</h2>
          <button
            type="button"
            onClick={onClose}
            className="rounded-lg p-1.5 text-[var(--color-text-muted)] hover:bg-[var(--color-surface)] hover:text-[var(--color-text)]"
          >
            <X size={20} />
          </button>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto px-4 py-4">
          {loading && items.length === 0 ? (
            <div className="flex items-center justify-center py-16">
              <Loader2 size={24} className="animate-spin text-[var(--color-accent)]" />
            </div>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <ShoppingBag size={48} className="mb-3 text-[var(--color-text-muted)]" />
              <p className="font-medium">Корзина пуста</p>
              <p className="mt-1 text-sm text-[var(--color-text-muted)]">
                Добавьте товары из каталога
              </p>
            </div>
          ) : (
            <ul className="space-y-4">
              {items.map((item) => {
                const src = resolveImageUrl(item.product.product_image);
                const busy = busyId === item.cart_item_id;
                return (
                  <li
                    key={item.cart_item_id}
                    className="flex gap-3 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-3"
                  >
                    <div className="shrink-0 overflow-hidden rounded-lg bg-[#1a1a1a]"
                      style={{ width: IMAGE_CONFIG.CART_THUMB_SIZE, height: IMAGE_CONFIG.CART_THUMB_SIZE }}>
                      {src ? (
                        <img
                          src={src}
                          alt={item.product.product_name}
                          className="h-full w-full object-cover"
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center text-xs text-[var(--color-text-muted)]">
                          —
                        </div>
                      )}
                    </div>

                    <div className="flex min-w-0 flex-1 flex-col">
                      <p className="line-clamp-2 text-sm font-medium">
                        {item.product.product_name}
                      </p>
                      <p className="mt-0.5 text-sm font-semibold text-[var(--color-accent)]">
                        {formatPrice(item.product.product_price)}
                      </p>

                      <div className="mt-auto flex items-center justify-between pt-2">
                        <div className="flex items-center gap-1 rounded-lg border border-[var(--color-border)] bg-[var(--color-background)]">
                          <button
                            type="button"
                            disabled={busy}
                            onClick={() => handleDelta(item.cart_item_id, -1)}
                            className="p-1.5 text-[var(--color-text-muted)] hover:text-[var(--color-text)] disabled:opacity-40"
                          >
                            <Minus size={14} />
                          </button>
                          <span className="min-w-[1.5rem] text-center text-sm font-medium">
                            {item.quantity}
                          </span>
                          <button
                            type="button"
                            disabled={busy}
                            onClick={() => handleDelta(item.cart_item_id, 1)}
                            className="p-1.5 text-[var(--color-text-muted)] hover:text-[var(--color-text)] disabled:opacity-40"
                          >
                            <Plus size={14} />
                          </button>
                        </div>

                        <button
                          type="button"
                          disabled={busy}
                          onClick={() => handleRemove(item.cart_item_id)}
                          className="rounded-lg p-1.5 text-[var(--color-text-muted)] hover:text-red-400 disabled:opacity-40"
                          title="Удалить"
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </div>
                  </li>
                );
              })}
            </ul>
          )}

          {(error || actionError) && (
            <p className="mt-4 text-sm text-red-300">{actionError || error}</p>
          )}
        </div>

        {/* Footer */}
        {items.length > 0 && (
          <div className="border-t border-[var(--color-border)] px-4 py-4">
            <div className="mb-3 flex items-center justify-between">
              <span className="text-[var(--color-text-muted)]">Итого</span>
              <span className="text-lg font-semibold text-[var(--color-accent)]">
                {formatPrice(String(totalPrice))}
              </span>
            </div>
            {checkoutSuccess && (
              <p className="mb-2 text-center text-sm text-emerald-300">{checkoutSuccess}</p>
            )}
            <button
              type="button"
              disabled={checkoutLoading}
              className="btn-gradient flex w-full items-center justify-center gap-2 rounded-xl py-3 text-sm font-semibold text-black disabled:opacity-60"
              onClick={async () => {
                setActionError(null);
                setCheckoutSuccess(null);
                setCheckoutLoading(true);
                try {
                  await onCheckout();
                  setCheckoutSuccess('Заказ оформлен');
                  window.setTimeout(() => setCheckoutSuccess(null), 2500);
                } catch (e) {
                  setActionError(e instanceof ApiError ? e.message : 'Не удалось оформить заказ');
                } finally {
                  setCheckoutLoading(false);
                }
              }}
            >
              {checkoutLoading && <Loader2 size={16} className="animate-spin" />}
              Оформить заказ
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

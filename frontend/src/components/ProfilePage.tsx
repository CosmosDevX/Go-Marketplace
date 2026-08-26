import { useCallback, useEffect, useState } from 'react';
import {
  ArrowLeft,
  Loader2,
  Package,
  User,
  RefreshCw,
} from 'lucide-react';
import { api, ApiError } from '../api/client';
import type { Order, UserInfo } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';

interface Props {
  user: UserInfo;
  onBack: () => void;
}

const STATUS_LABELS: Record<string, string> = {
  Pending: 'В обработке',
  pending: 'В обработке',
  Completed: 'Завершён',
  completed: 'Завершён',
  Cancelled: 'Отменён',
  cancelled: 'Отменён',
  Shipped: 'Отправлен',
  shipped: 'Отправлен',
};

function statusLabel(status: string) {
  return STATUS_LABELS[status] ?? status;
}

function statusClass(status: string) {
  const s = status.toLowerCase();
  if (s === 'completed') return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-300';
  if (s === 'cancelled') return 'border-red-500/30 bg-red-500/10 text-red-300';
  if (s === 'shipped') return 'border-sky-500/30 bg-sky-500/10 text-sky-300';
  return 'border-amber-500/30 bg-amber-500/10 text-amber-300';
}

export function ProfilePage({ user, onBack }: Props) {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadOrders = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getOrders();
      // сохраняем порядок с бэкенда (order_id ASC)
      setOrders(Array.isArray(data) ? data : []);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Не удалось загрузить заказы');
      setOrders([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadOrders();
  }, [loadOrders]);

  return (
    <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 py-6 sm:py-8">
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={onBack}
          className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-2 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
        >
          <ArrowLeft size={16} />
          Каталог
        </button>
        <h1 className="text-2xl font-semibold tracking-tight">Профиль</h1>
      </div>

      <div className="mb-8 flex items-center gap-4 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-5">
        <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-amber-500/25 to-orange-600/15 text-[var(--color-accent)]">
          <User size={28} />
        </div>
        <div>
          <p className="text-lg font-semibold">{user.username}</p>
          <p className="text-sm text-[var(--color-text-muted)]">
            Роли: {user.roles?.length ? user.roles.join(', ') : 'user'}
          </p>
        </div>
      </div>

      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Мои заказы</h2>
        <button
          type="button"
          onClick={loadOrders}
          disabled={loading}
          className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] px-3 py-1.5 text-sm text-[var(--color-text-muted)] hover:text-[var(--color-text)] disabled:opacity-40"
        >
          <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
          Обновить
        </button>
      </div>

      {error && (
        <div className="mb-4 rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
          {error}
        </div>
      )}

      {loading && orders.length === 0 ? (
        <div className="flex justify-center py-16">
          <Loader2 size={28} className="animate-spin text-[var(--color-accent)]" />
        </div>
      ) : orders.length === 0 && !error ? (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] py-16 text-center">
          <Package size={40} className="mb-3 text-[var(--color-text-muted)]" />
          <p className="font-medium">Заказов пока нет</p>
          <p className="mt-1 text-sm text-[var(--color-text-muted)]">
            Оформите заказ из корзины
          </p>
        </div>
      ) : (
        <ul className="space-y-4">
          {orders.map((order) => (
            <li
              key={order.order_id}
              className="overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)]"
            >
              <div className="flex flex-wrap items-center justify-between gap-2 border-b border-[var(--color-border)] px-4 py-3">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-semibold">Заказ #{order.order_id}</span>
                  <span
                    className={`rounded-full border px-2.5 py-0.5 text-xs font-medium ${statusClass(order.order_status)}`}
                  >
                    {statusLabel(order.order_status)}
                  </span>
                </div>
                <span className="text-base font-semibold text-gradient">
                  {formatPrice(order.order_total)}
                </span>
              </div>

              <ul className="divide-y divide-[var(--color-border)]">
                {(order.order_items ?? []).map((item) => {
                  const p = item.product;
                  const src = resolveImageUrl(p.product_image);
                  return (
                    <li
                      key={item.order_item_id}
                      className="flex gap-3 px-4 py-3"
                    >
                      <div className="h-14 w-14 shrink-0 overflow-hidden rounded-lg bg-[#1a1a1a]">
                        {src ? (
                          <img
                            src={src}
                            alt=""
                            className="h-full w-full object-cover"
                          />
                        ) : (
                          <div className="flex h-full w-full items-center justify-center text-xs text-[var(--color-text-muted)]">
                            —
                          </div>
                        )}
                      </div>
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-medium">{p.product_name}</p>
                        <p className="mt-0.5 text-xs text-[var(--color-text-muted)]">
                          {formatPrice(p.product_price)} × {item.order_item_quantity}
                        </p>
                      </div>
                      <div className="shrink-0 text-right">
                        <p className="text-sm font-semibold text-[var(--color-accent)]">
                          {formatPrice(item.order_item_total)}
                        </p>
                        <p className="text-[10px] text-[var(--color-text-muted)]">
                          {item.order_item_quantity} шт.
                        </p>
                      </div>
                    </li>
                  );
                })}
              </ul>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

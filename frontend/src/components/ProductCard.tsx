import { useState, type CSSProperties } from 'react';
import { motion } from 'framer-motion';
import { ShoppingCart, ImageOff, Loader2, Plus, Minus } from 'lucide-react';
import type { Product } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';
import { ApiError } from '../api/client';
import { IMAGE_CONFIG } from '../config/images';

interface Props {
  product: Product;
  /** Если товар уже в корзине */
  cartItemId?: number | null;
  quantity?: number;
  onAddToCart: (productId: number) => Promise<void>;
  onChangeQuantity: (cartItemId: number, delta: 1 | -1) => Promise<void>;
  requireAuth: () => void;
  isAuthenticated: boolean;
  index?: number;
}

export function ProductCard({
  product,
  cartItemId,
  quantity = 0,
  onAddToCart,
  onChangeQuantity,
  requireAuth,
  isAuthenticated,
  index = 0,
}: Props) {
  const src = resolveImageUrl(product.product_image);
  const [failed, setFailed] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const inCart = cartItemId != null && quantity > 0;

  const handleAdd = async () => {
    if (!isAuthenticated) {
      requireAuth();
      return;
    }
    setBusy(true);
    setError(null);
    try {
      await onAddToCart(product.product_id);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Не удалось добавить');
    } finally {
      setBusy(false);
    }
  };

  const handleDelta = async (delta: 1 | -1) => {
    if (cartItemId == null) return;
    setBusy(true);
    setError(null);
    try {
      await onChangeQuantity(cartItemId, delta);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Ошибка обновления');
    } finally {
      setBusy(false);
    }
  };

  const imageStyle: CSSProperties = {
    aspectRatio: IMAGE_CONFIG.PRODUCT_CARD_ASPECT,
    ...(IMAGE_CONFIG.PRODUCT_CARD_MAX_HEIGHT
      ? { maxHeight: IMAGE_CONFIG.PRODUCT_CARD_MAX_HEIGHT }
      : {}),
  };

  return (
    <motion.article
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, delay: Math.min(index * 0.04, 0.3), ease: 'easeOut' }}
      whileHover={{ y: -6 }}
      className="card-glow group flex flex-col rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] overflow-hidden"
    >
      <div className="relative w-full overflow-hidden bg-[#141414]" style={imageStyle}>
        <div className="pointer-events-none absolute inset-0 z-[1] bg-gradient-to-t from-[var(--color-surface)]/80 via-transparent to-transparent opacity-60" />
        {src && !failed ? (
          <img
            src={src}
            alt={product.product_name}
            loading="lazy"
            decoding="async"
            onError={() => setFailed(true)}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full flex-col items-center justify-center gap-2 bg-gradient-to-br from-[#1a1a22] to-[#121218] text-[var(--color-text-muted)]">
            <ImageOff size={32} strokeWidth={1.5} />
            <span className="text-xs">Нет фото</span>
          </div>
        )}
      </div>

      <div className="flex flex-1 flex-col gap-1.5 p-3 sm:p-4">
        <h3 className="line-clamp-2 text-sm font-medium leading-snug text-[var(--color-text)]">
          {product.product_name}
        </h3>

        <p className="line-clamp-2 text-xs text-[var(--color-text-muted)]">
          {product.product_description}
        </p>

        {error && <p className="text-xs text-red-300">{error}</p>}

        <div className="mt-auto flex items-center justify-between gap-2 pt-2">
          <span className="min-w-0 truncate text-base sm:text-lg font-semibold text-gradient">
            {formatPrice(product.product_price)}
          </span>

          {inCart ? (
            <div className="flex shrink-0 items-center gap-0.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-background)]">
              <button
                type="button"
                disabled={busy}
                onClick={() => handleDelta(-1)}
                className="rounded-l-xl p-2 text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text)] disabled:opacity-40"
                aria-label="Уменьшить"
              >
                <Minus size={14} />
              </button>
              <span className="min-w-[1.5rem] text-center text-sm font-semibold tabular-nums">
                {quantity}
              </span>
              <button
                type="button"
                disabled={busy}
                onClick={() => handleDelta(1)}
                className="rounded-r-xl p-2 text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text)] disabled:opacity-40"
                aria-label="Увеличить"
              >
                <Plus size={14} />
              </button>
            </div>
          ) : (
            <motion.button
              type="button"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              disabled={busy}
              onClick={handleAdd}
              title="В корзину"
              className="btn-gradient flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-black disabled:opacity-60"
            >
              {busy ? (
                <Loader2 size={16} className="animate-spin" />
              ) : (
                <ShoppingCart size={16} strokeWidth={2.5} />
              )}
            </motion.button>
          )}
        </div>
      </div>
    </motion.article>
  );
}

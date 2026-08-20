import { useState, type CSSProperties } from 'react';
import { motion } from 'framer-motion';
import { ShoppingCart, ImageOff, Loader2, Check } from 'lucide-react';
import type { Product } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';
import { ApiError } from '../api/client';
import { IMAGE_CONFIG } from '../config/images';

interface Props {
  product: Product;
  onAddToCart: (productId: number) => Promise<void>;
  requireAuth: () => void;
  isAuthenticated: boolean;
  index?: number;
}

export function ProductCard({
  product,
  onAddToCart,
  requireAuth,
  isAuthenticated,
  index = 0,
}: Props) {
  const src = resolveImageUrl(product.product_image);
  const [failed, setFailed] = useState(false);
  const [adding, setAdding] = useState(false);
  const [added, setAdded] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAdd = async () => {
    if (!isAuthenticated) {
      requireAuth();
      return;
    }
    setAdding(true);
    setError(null);
    try {
      await onAddToCart(product.product_id);
      setAdded(true);
      setTimeout(() => setAdded(false), 1500);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Не удалось добавить');
    } finally {
      setAdding(false);
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

      <div className="flex flex-1 flex-col gap-2 p-4">
        <h3 className="line-clamp-2 text-sm sm:text-base font-medium leading-snug text-[var(--color-text)]">
          {product.product_name}
        </h3>

        <p className="line-clamp-2 text-xs text-[var(--color-text-muted)]">
          {product.product_description}
        </p>

        {error && <p className="text-xs text-red-300">{error}</p>}

        <div className="mt-auto flex items-center justify-between gap-2 pt-3">
          <span className="text-lg font-semibold text-gradient">
            {formatPrice(product.product_price)}
          </span>

          <motion.button
            type="button"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
            disabled={adding}
            onClick={handleAdd}
            className="btn-gradient flex items-center gap-1.5 rounded-xl px-3 py-2 text-sm font-medium text-black disabled:opacity-60"
          >
            {adding ? (
              <Loader2 size={16} className="animate-spin" />
            ) : added ? (
              <Check size={16} strokeWidth={2.5} />
            ) : (
              <ShoppingCart size={16} strokeWidth={2.5} />
            )}
            <span className="hidden sm:inline">
              {added ? 'Добавлено' : 'В корзину'}
            </span>
          </motion.button>
        </div>
      </div>
    </motion.article>
  );
}

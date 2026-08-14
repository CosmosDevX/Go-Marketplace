import { motion } from 'framer-motion';
import { ImageOff, ShoppingCart } from 'lucide-react';
import { useState } from 'react';
import type { Product } from '../types';
import { formatPrice } from '../utils/emoji';
import { resolveImageUrl } from '../utils/image';

interface Props {
  product: Product;
}

export function ProductCard({ product }: Props) {
  const src = resolveImageUrl(product.product_image);
  const [failed, setFailed] = useState(false);

  return (
    <motion.article
      whileHover={{ y: -4 }}
      transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      className="group flex flex-col rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] overflow-hidden transition-colors hover:border-[var(--color-accent)]/40 hover:bg-[var(--color-surface-hover)]"
    >
      {/* Image */}
      <div className="relative aspect-square w-full bg-[#1a1a1a] overflow-hidden">
        {src && !failed ? (
          <img
            src={src}
            alt={product.product_name}
            loading="lazy"
            onError={() => setFailed(true)}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        ) : (
          <div className="flex h-full w-full flex-col items-center justify-center gap-2 text-[var(--color-text-muted)]">
            <ImageOff size={32} strokeWidth={1.5} />
            <span className="text-xs">Нет фото</span>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="flex flex-1 flex-col gap-2 p-4">
        <h3 className="line-clamp-2 text-sm sm:text-base font-medium leading-snug text-[var(--color-text)]">
          {product.product_name}
        </h3>

        <p className="line-clamp-2 text-xs text-[var(--color-text-muted)]">
          {product.product_description}
        </p>

        <div className="mt-auto flex items-center justify-between gap-2 pt-3">
          <span className="text-lg font-semibold text-[var(--color-accent)]">
            {formatPrice(product.product_price)}
          </span>

          <button
            type="button"
            className="flex items-center gap-1.5 rounded-xl bg-[var(--color-accent)] px-3 py-2 text-sm font-medium text-black transition-colors hover:bg-[var(--color-accent-hover)] active:scale-95"
            onClick={() => {
              // пока ничего не делаем
            }}
          >
            <ShoppingCart size={16} strokeWidth={2.5} />
            <span className="hidden sm:inline">В корзину</span>
          </button>
        </div>
      </div>
    </motion.article>
  );
}

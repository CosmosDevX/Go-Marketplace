import { motion } from 'framer-motion';
import { ShoppingCart } from 'lucide-react';
import type { Product } from '../types';
import { getProductEmoji, formatPrice } from '../utils/emoji';

interface Props {
  product: Product;
}

export function ProductCard({ product }: Props) {
  const emoji = getProductEmoji(product.product_name, product.category.category_slug);

  return (
    <motion.article
      whileHover={{ y: -4 }}
      transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      className="group flex flex-col rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] overflow-hidden transition-colors hover:border-[var(--color-accent)]/40 hover:bg-[var(--color-surface-hover)]"
    >
      {/* Emoji "image" */}
      <div className="relative flex h-40 sm:h-44 items-center justify-center bg-gradient-to-b from-[#222] to-[#1a1a1a]">
        <span className="text-6xl sm:text-7xl select-none drop-shadow-lg transition-transform duration-300 group-hover:scale-110">
          {emoji}
        </span>
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

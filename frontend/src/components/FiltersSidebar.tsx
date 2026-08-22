import type { ReactNode, FormEvent } from 'react';
import { motion } from 'framer-motion';
import {
  ArrowDownAZ,
  ArrowUpAZ,
  ArrowDownWideNarrow,
  ArrowUpNarrowWide,
  RotateCcw,
  Search,
} from 'lucide-react';
import type { Category } from '../types';

export type SortBy = 'price' | 'name' | null;

interface Props {
  categories: Category[];
  categoriesLoading?: boolean;
  activeCategory: string | null;
  sortBy: SortBy;
  asc: boolean;
  searchInput: string;
  onSearchInputChange: (value: string) => void;
  onSearchSubmit: (value: string) => void;
  onCategoryChange: (slug: string | null) => void;
  onSortChange: (sortBy: SortBy, asc: boolean) => void;
  onReset: () => void;
}

export function FiltersSidebar({
  categories,
  categoriesLoading,
  activeCategory,
  sortBy,
  asc,
  searchInput,
  onSearchInputChange,
  onSearchSubmit,
  onCategoryChange,
  onSortChange,
  onReset,
}: Props) {
  const isPriceAsc = sortBy === 'price' && asc;
  const isPriceDesc = sortBy === 'price' && !asc;
  const isNameAsc = sortBy === 'name' && asc;
  const isNameDesc = sortBy === 'name' && !asc;

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    onSearchSubmit(searchInput);
  };

  return (
    <aside className="w-full shrink-0 lg:w-64">
      <div className="rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-4 lg:sticky lg:top-20">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-[var(--color-text-muted)]">
            Фильтры
          </h2>
          <button
            type="button"
            onClick={onReset}
            className="flex items-center gap-1 rounded-lg px-2 py-1 text-xs text-[var(--color-text-muted)] hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)]"
            title="Сбросить"
          >
            <RotateCcw size={12} />
            Сброс
          </button>
        </div>

        {/* Search */}
        <form onSubmit={handleSearch} className="mb-6">
          <p className="mb-2 text-sm font-medium">Поиск</p>
          <div className="flex gap-1.5">
            <input
              type="search"
              value={searchInput}
              onChange={(e) => onSearchInputChange(e.target.value)}
              placeholder="Название..."
              className="min-w-0 flex-1 rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3 py-2 text-sm outline-none focus:border-[var(--color-accent)]"
            />
            <button
              type="submit"
              className="btn-gradient flex items-center justify-center rounded-xl px-3 text-black"
              title="Искать"
            >
              <Search size={16} />
            </button>
          </div>
        </form>

        <div className="mb-6">
          <p className="mb-2 text-sm font-medium">Категория</p>
          <div className="flex flex-col gap-1">
            <FilterRow
              label="Все"
              active={activeCategory === null}
              onClick={() => onCategoryChange(null)}
            />
            {categoriesLoading
              ? Array.from({ length: 3 }).map((_, i) => (
                  <div
                    key={i}
                    className="h-9 rounded-xl border border-[var(--color-border)] bg-[var(--color-background)]"
                  />
                ))
              : categories.map((c) => (
                  <FilterRow
                    key={c.category_id}
                    label={c.category_name}
                    active={activeCategory === c.category_slug}
                    onClick={() => onCategoryChange(c.category_slug)}
                  />
                ))}
          </div>
        </div>

        <div className="mb-6">
          <p className="mb-2 text-sm font-medium">Сортировка по цене</p>
          <div className="flex flex-col gap-1">
            <FilterRow
              label="По возрастанию"
              icon={<ArrowUpNarrowWide size={14} />}
              active={isPriceAsc}
              onClick={() => onSortChange(isPriceAsc ? null : 'price', true)}
            />
            <FilterRow
              label="По убыванию"
              icon={<ArrowDownWideNarrow size={14} />}
              active={isPriceDesc}
              onClick={() => onSortChange(isPriceDesc ? null : 'price', false)}
            />
          </div>
        </div>

        <div>
          <p className="mb-2 text-sm font-medium">Сортировка по названию</p>
          <div className="flex flex-col gap-1">
            <FilterRow
              label="А → Я"
              icon={<ArrowDownAZ size={14} />}
              active={isNameAsc}
              onClick={() => onSortChange(isNameAsc ? null : 'name', true)}
            />
            <FilterRow
              label="Я → А"
              icon={<ArrowUpAZ size={14} />}
              active={isNameDesc}
              onClick={() => onSortChange(isNameDesc ? null : 'name', false)}
            />
          </div>
        </div>
      </div>
    </aside>
  );
}

function FilterRow({
  label,
  active,
  onClick,
  icon,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
  icon?: ReactNode;
}) {
  return (
    <motion.button
      type="button"
      whileTap={{ scale: 0.98 }}
      onClick={onClick}
      className={`
        flex w-full items-center gap-2 rounded-xl px-3 py-2 text-left text-sm transition-colors
        ${
          active
            ? 'border border-[var(--color-accent)]/30 bg-gradient-to-r from-amber-500/20 to-orange-600/10 text-[var(--color-accent)]'
            : 'border border-transparent text-[var(--color-text-muted)] hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)]'
        }
      `}
    >
      {icon}
      <span className="truncate">{label}</span>
    </motion.button>
  );
}

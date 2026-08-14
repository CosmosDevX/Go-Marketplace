import type { Category } from '../types';

interface Props {
  categories: Category[];
  activeSlug: string | null;
  onSelect: (slug: string | null) => void;
  loading?: boolean;
}

export function CategoryChips({ categories, activeSlug, onSelect, loading }: Props) {
  if (loading) {
    return (
      <div className="flex gap-2 overflow-x-auto pb-2">
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className="h-9 w-20 shrink-0 rounded-full bg-[var(--color-surface)] border border-[var(--color-border)]"
          />
        ))}
      </div>
    );
  }

  return (
    <div className="flex gap-2 overflow-x-auto pb-2 -mx-1 px-1">
      <Chip
        label="Все"
        active={activeSlug === null}
        onClick={() => onSelect(null)}
      />
      {categories.map((cat) => (
        <Chip
          key={cat.category_id}
          label={cat.category_name}
          active={activeSlug === cat.category_slug}
          onClick={() => onSelect(cat.category_slug)}
        />
      ))}
    </div>
  );
}

function Chip({
  label,
  active,
  onClick,
}: {
  label: string;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`
        shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-colors
        ${
          active
            ? 'bg-[var(--color-accent)] text-black'
            : 'bg-[var(--color-surface)] text-[var(--color-text-muted)] hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)] border border-[var(--color-border)]'
        }
      `}
    >
      {label}
    </button>
  );
}

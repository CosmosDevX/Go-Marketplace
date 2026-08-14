import { ChevronLeft, ChevronRight } from 'lucide-react';

interface Props {
  page: number;
  hasNext: boolean;
  onPrev: () => void;
  onNext: () => void;
  loading?: boolean;
}

export function Pagination({ page, hasNext, onPrev, onNext, loading }: Props) {
  return (
    <div className="flex items-center justify-center gap-4 pt-8 pb-4">
      <button
        type="button"
        disabled={page <= 1 || loading}
        onClick={onPrev}
        className="flex items-center gap-1 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2.5 text-sm font-medium text-[var(--color-text)] transition-colors enabled:hover:border-[var(--color-accent)]/50 enabled:hover:bg-[var(--color-surface-hover)] disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <ChevronLeft size={18} />
        <span className="hidden sm:inline">Предыдущая</span>
      </button>

      <span className="min-w-[5rem] text-center text-sm text-[var(--color-text-muted)]">
        Страница <span className="font-semibold text-[var(--color-text)]">{page}</span>
      </span>

      <button
        type="button"
        disabled={!hasNext || loading}
        onClick={onNext}
        className="flex items-center gap-1 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2.5 text-sm font-medium text-[var(--color-text)] transition-colors enabled:hover:border-[var(--color-accent)]/50 enabled:hover:bg-[var(--color-surface-hover)] disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <span className="hidden sm:inline">Следующая</span>
        <ChevronRight size={18} />
      </button>
    </div>
  );
}

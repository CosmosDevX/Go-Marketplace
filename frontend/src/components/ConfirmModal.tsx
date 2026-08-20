import { motion, AnimatePresence } from 'framer-motion';
import { AlertTriangle, Loader2, X } from 'lucide-react';

interface Props {
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  loading?: boolean;
  danger?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmModal({
  open,
  title,
  message,
  confirmLabel = 'Удалить',
  cancelLabel = 'Отмена',
  loading = false,
  danger = true,
  onConfirm,
  onCancel,
}: Props) {
  return (
    <AnimatePresence>
      {open && (
        <div className="fixed inset-0 z-[110] flex items-center justify-center p-4">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 bg-black/70 backdrop-blur-sm"
            onClick={loading ? undefined : onCancel}
          />

          <motion.div
            role="dialog"
            aria-modal="true"
            initial={{ opacity: 0, scale: 0.9, y: 10 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 8 }}
            transition={{ type: 'spring', duration: 0.35 }}
            className="relative w-full max-w-sm overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-6 shadow-2xl"
          >
            <div className="pointer-events-none absolute -right-10 -top-10 h-28 w-28 rounded-full bg-red-500/15 blur-2xl" />

            <button
              type="button"
              disabled={loading}
              onClick={onCancel}
              className="absolute right-3 top-3 rounded-lg p-1.5 text-[var(--color-text-muted)] hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)] disabled:opacity-40"
            >
              <X size={18} />
            </button>

            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gradient-to-br from-red-500/25 to-red-600/10 text-red-400">
              <AlertTriangle size={24} />
            </div>

            <h3 className="mb-2 text-lg font-semibold text-[var(--color-text)]">{title}</h3>
            <p className="mb-6 text-sm leading-relaxed text-[var(--color-text-muted)]">{message}</p>

            <div className="flex gap-3">
              <button
                type="button"
                disabled={loading}
                onClick={onCancel}
                className="flex-1 rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] py-2.5 text-sm font-medium text-[var(--color-text)] transition-colors hover:bg-[var(--color-surface-hover)] disabled:opacity-40"
              >
                {cancelLabel}
              </button>
              <motion.button
                type="button"
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
                disabled={loading}
                onClick={onConfirm}
                className={`flex flex-1 items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold transition-colors disabled:opacity-60 ${
                  danger
                    ? 'bg-gradient-to-r from-red-500 to-red-600 text-white shadow-lg shadow-red-500/20 hover:from-red-400 hover:to-red-500'
                    : 'btn-gradient text-black'
                }`}
              >
                {loading && <Loader2 size={16} className="animate-spin" />}
                {confirmLabel}
              </motion.button>
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}

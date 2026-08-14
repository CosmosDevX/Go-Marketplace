import { Store, LogIn, UserPlus, LogOut, User } from 'lucide-react';
import type { UserInfo } from '../types';

interface Props {
  user: UserInfo | null;
  onLoginClick: () => void;
  onRegisterClick: () => void;
  onLogout: () => void;
}

export function Header({ user, onLoginClick, onRegisterClick, onLogout }: Props) {
  return (
    <header className="sticky top-0 z-50 border-b border-[var(--color-border)] bg-[var(--color-background)]/90 backdrop-blur-md">
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6">
        <div className="flex items-center gap-2.5">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-[var(--color-accent)] text-black">
            <Store size={20} strokeWidth={2.5} />
          </div>
          <span className="text-lg font-semibold tracking-tight">
            Market<span className="text-[var(--color-accent)]">Place</span>
          </span>
        </div>

        <div className="flex items-center gap-2">
          {user ? (
            <>
              <div className="flex items-center gap-2 rounded-xl bg-[var(--color-surface)] px-3 py-1.5 text-sm">
                <User size={16} className="text-[var(--color-accent)]" />
                <span className="max-w-[120px] truncate font-medium">{user.username}</span>
              </div>
              <button
                type="button"
                onClick={onLogout}
                className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text-muted)] transition-colors hover:border-red-500/40 hover:text-red-400"
                title="Выйти"
              >
                <LogOut size={16} />
                <span className="hidden sm:inline">Выйти</span>
              </button>
            </>
          ) : (
            <>
              <button
                type="button"
                onClick={onLoginClick}
                className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text)] transition-colors hover:border-[var(--color-accent)]/50 hover:bg-[var(--color-surface-hover)]"
              >
                <LogIn size={16} />
                <span className="hidden sm:inline">Войти</span>
              </button>
              <button
                type="button"
                onClick={onRegisterClick}
                className="flex items-center gap-1.5 rounded-xl bg-[var(--color-accent)] px-3 py-1.5 text-sm font-medium text-black transition-colors hover:bg-[var(--color-accent-hover)]"
              >
                <UserPlus size={16} />
                <span className="hidden sm:inline">Регистрация</span>
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

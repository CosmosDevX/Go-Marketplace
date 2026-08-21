import { useState } from 'react';
import { motion } from 'framer-motion';
import {
  Store,
  LogIn,
  UserPlus,
  LogOut,
  User,
  ShoppingCart,
  LayoutDashboard,
  Shield,
} from 'lucide-react';
import type { UserInfo } from '../types';

interface Props {
  user: UserInfo | null;
  cartCount: number;
  canAccessSellerPanel: boolean;
  canAccessAdminPanel: boolean;
  onLogoClick: () => void;
  onLoginClick: () => void;
  onRegisterClick: () => void;
  onLogout: () => void;
  onCartClick: () => void;
  onSellerClick: () => void;
  onAdminClick: () => void;
}

export function Header({
  user,
  cartCount,
  canAccessSellerPanel,
  canAccessAdminPanel,
  onLogoClick,
  onLoginClick,
  onRegisterClick,
  onLogout,
  onCartClick,
  onSellerClick,
  onAdminClick,
}: Props) {
  const [logoFailed, setLogoFailed] = useState(false);

  return (
    <header className="sticky top-0 z-50 border-b border-[var(--color-border)] bg-[var(--color-background)]/80 backdrop-blur-xl">
      <div className="absolute inset-x-0 bottom-0 h-px bg-gradient-to-r from-transparent via-[var(--color-accent)]/40 to-transparent" />
      <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6">
        <motion.button
          type="button"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          onClick={onLogoClick}
          className="flex items-center gap-2.5 rounded-lg outline-none"
        >
          {!logoFailed ? (
            <img
              src="/marketplaceLogo.png"
              alt="MarketPlace"
              className="h-9 w-9 rounded-xl object-contain"
              onError={() => setLogoFailed(true)}
            />
          ) : (
            <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-amber-400 to-orange-600 text-black shadow-lg shadow-amber-500/20">
              <Store size={20} strokeWidth={2.5} />
            </div>
          )}
          <span className="text-lg font-semibold tracking-tight">
            Market<span className="text-gradient">Place</span>
          </span>
        </motion.button>

        <div className="flex items-center gap-2">
          {canAccessAdminPanel && (
            <motion.button
              type="button"
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
              onClick={onAdminClick}
              className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text)] transition-colors hover:border-[var(--color-accent)]/50 hover:bg-[var(--color-surface-hover)]"
            >
              <Shield size={16} />
              <span className="hidden sm:inline">Админ-панель</span>
            </motion.button>
          )}

          {canAccessSellerPanel && (
            <motion.button
              type="button"
              whileHover={{ scale: 1.03 }}
              whileTap={{ scale: 0.97 }}
              onClick={onSellerClick}
              className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text)] transition-colors hover:border-[var(--color-accent)]/50 hover:bg-[var(--color-surface-hover)]"
            >
              <LayoutDashboard size={16} />
              <span className="hidden sm:inline">Панель продавца</span>
            </motion.button>
          )}

          <motion.button
            type="button"
            whileHover={{ scale: 1.03 }}
            whileTap={{ scale: 0.97 }}
            onClick={onCartClick}
            className="relative flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text)] transition-colors hover:border-[var(--color-accent)]/50 hover:bg-[var(--color-surface-hover)]"
            title="Корзина"
          >
            <ShoppingCart size={16} />
            <span className="hidden sm:inline">Корзина</span>
            {cartCount > 0 && (
              <motion.span
                initial={{ scale: 0 }}
                animate={{ scale: 1 }}
                className="absolute -right-1.5 -top-1.5 flex h-5 min-w-5 items-center justify-center rounded-full bg-gradient-to-br from-amber-400 to-orange-500 px-1 text-[10px] font-bold text-black shadow-md shadow-amber-500/30"
              >
                {cartCount > 99 ? '99+' : cartCount}
              </motion.span>
            )}
          </motion.button>

          {user ? (
            <>
              <div className="flex items-center gap-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm">
                <User size={16} className="text-[var(--color-accent)]" />
                <span className="max-w-[120px] truncate font-medium">{user.username}</span>
              </div>
              <motion.button
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.97 }}
                onClick={onLogout}
                className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text-muted)] transition-colors hover:border-red-500/40 hover:text-red-400"
                title="Выйти"
              >
                <LogOut size={16} />
                <span className="hidden sm:inline">Выйти</span>
              </motion.button>
            </>
          ) : (
            <>
              <motion.button
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.97 }}
                onClick={onLoginClick}
                className="flex items-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] px-3 py-1.5 text-sm text-[var(--color-text)] transition-colors hover:border-[var(--color-accent)]/50 hover:bg-[var(--color-surface-hover)]"
              >
                <LogIn size={16} />
                <span className="hidden sm:inline">Войти</span>
              </motion.button>
              <motion.button
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.97 }}
                onClick={onRegisterClick}
                className="btn-gradient flex items-center gap-1.5 rounded-xl px-3 py-1.5 text-sm font-medium text-black"
              >
                <UserPlus size={16} />
                <span className="hidden sm:inline">Регистрация</span>
              </motion.button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}

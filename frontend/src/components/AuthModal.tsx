import { useState, FormEvent } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Loader2, Eye, EyeOff } from 'lucide-react';
import { api, ApiError } from '../api/client';

interface Props {
  mode: 'login' | 'register' | null;
  onClose: () => void;
  onSwitchMode: (mode: 'login' | 'register') => void;
  onSuccess: (token: string, username?: string, roles?: string[]) => void;
}

export function AuthModal({ mode, onClose, onSwitchMode, onSuccess }: Props) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isLogin = mode === 'login';
  const isOpen = mode !== null;

  const reset = () => {
    setUsername('');
    setPassword('');
    setEmail('');
    setError(null);
    setShowPassword(false);
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      if (isLogin) {
        const res = await api.login({ username, password });
        onSuccess(res.access_token, username, res.roles ?? []);
      } else {
        await api.register({ username, password, email });
        const res = await api.login({ username, password });
        onSuccess(res.access_token, username, res.roles ?? []);
      }
      reset();
    } catch (err) {
      const msg =
        err instanceof ApiError
          ? err.message
          : isLogin
            ? 'Ошибка входа'
            : 'Ошибка регистрации';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 bg-black/70 backdrop-blur-sm"
            onClick={handleClose}
          />

          <motion.div
            initial={{ opacity: 0, scale: 0.92, y: 12 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 8 }}
            transition={{ type: 'spring', duration: 0.4 }}
            className="relative w-full max-w-md overflow-hidden rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-6 shadow-2xl"
          >
            <div className="pointer-events-none absolute -right-16 -top-16 h-40 w-40 rounded-full bg-gradient-to-br from-amber-500/20 to-orange-600/10 blur-2xl" />
            <div className="pointer-events-none absolute -bottom-12 -left-12 h-32 w-32 rounded-full bg-violet-500/10 blur-2xl" />

            <button
              type="button"
              onClick={handleClose}
              className="absolute right-4 top-4 z-10 rounded-lg p-1.5 text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-hover)] hover:text-[var(--color-text)]"
            >
              <X size={18} />
            </button>

            <h2 className="relative mb-1 text-xl font-semibold">
              {isLogin ? 'Вход' : 'Регистрация'}
            </h2>
            <p className="relative mb-6 text-sm text-[var(--color-text-muted)]">
              {isLogin ? 'Войдите в свой аккаунт' : 'Создайте новый аккаунт'}
            </p>

            <form onSubmit={handleSubmit} className="relative space-y-4">
              {!isLogin && (
                <div>
                  <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                    Email
                  </label>
                  <input
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm text-[var(--color-text)] outline-none transition-colors placeholder:text-[var(--color-text-muted)] focus:border-[var(--color-accent)] focus:ring-1 focus:ring-[var(--color-accent)]/30"
                    placeholder="you@example.com"
                    autoComplete="email"
                  />
                </div>
              )}

              <div>
                <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                  Имя пользователя
                </label>
                <input
                  type="text"
                  required
                  minLength={3}
                  maxLength={40}
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm text-[var(--color-text)] outline-none transition-colors placeholder:text-[var(--color-text-muted)] focus:border-[var(--color-accent)] focus:ring-1 focus:ring-[var(--color-accent)]/30"
                  placeholder="username"
                  autoComplete="username"
                />
              </div>

              <div>
                <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
                  Пароль
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? 'text' : 'password'}
                    required
                    minLength={8}
                    maxLength={60}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 pr-10 text-sm text-[var(--color-text)] outline-none transition-colors placeholder:text-[var(--color-text-muted)] focus:border-[var(--color-accent)] focus:ring-1 focus:ring-[var(--color-accent)]/30"
                    placeholder="••••••••"
                    autoComplete={isLogin ? 'current-password' : 'new-password'}
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword((v) => !v)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>

              {error && (
                <motion.div
                  initial={{ opacity: 0, y: -4 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-300"
                >
                  {error}
                </motion.div>
              )}

              <motion.button
                type="submit"
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.98 }}
                disabled={loading}
                className="btn-gradient flex w-full items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold text-black disabled:opacity-60"
              >
                {loading && <Loader2 size={16} className="animate-spin" />}
                {isLogin ? 'Войти' : 'Зарегистрироваться'}
              </motion.button>
            </form>

            <p className="relative mt-5 text-center text-sm text-[var(--color-text-muted)]">
              {isLogin ? (
                <>
                  Нет аккаунта?{' '}
                  <button
                    type="button"
                    onClick={() => {
                      reset();
                      onSwitchMode('register');
                    }}
                    className="font-medium text-[var(--color-accent)] hover:underline"
                  >
                    Зарегистрироваться
                  </button>
                </>
              ) : (
                <>
                  Уже есть аккаунт?{' '}
                  <button
                    type="button"
                    onClick={() => {
                      reset();
                      onSwitchMode('login');
                    }}
                    className="font-medium text-[var(--color-accent)] hover:underline"
                  >
                    Войти
                  </button>
                </>
              )}
            </p>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}

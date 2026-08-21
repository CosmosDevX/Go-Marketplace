import { useState, FormEvent } from 'react';
import { motion } from 'framer-motion';
import { Loader2, Search, UserCog, ShieldPlus, ShieldMinus } from 'lucide-react';
import { api, ApiError } from '../api/client';

const AVAILABLE_ROLES = ['user', 'seller', 'admin'] as const;

export function AdminRoles() {
  const [username, setUsername] = useState('');
  const [role, setRole] = useState<string>('seller');
  const [roles, setRoles] = useState<string[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [actionLoading, setActionLoading] = useState<'grant' | 'revoke' | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const clearFeedback = () => {
    setMessage(null);
    setError(null);
  };

  const handleLookup = async (e?: FormEvent) => {
    e?.preventDefault();
    if (!username.trim()) return;
    clearFeedback();
    setLoading(true);
    setRoles(null);
    try {
      const data = await api.getUserRoles(username.trim());
      setRoles(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Не удалось получить роли');
      setRoles(null);
    } finally {
      setLoading(false);
    }
  };

  const handleGrant = async () => {
    if (!username.trim() || !role) return;
    clearFeedback();
    setActionLoading('grant');
    try {
      const res = await api.grantRole({
        username: username.trim(),
        role,
      });
      setMessage(res.message || 'role granted');
      const data = await api.getUserRoles(username.trim());
      setRoles(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Ошибка выдачи роли');
    } finally {
      setActionLoading(null);
    }
  };

  const handleRevoke = async () => {
    if (!username.trim() || !role) return;
    clearFeedback();
    setActionLoading('revoke');
    try {
      const res = await api.revokeRole(username.trim(), role);
      setMessage(res.message || 'role deleted');
      const data = await api.getUserRoles(username.trim());
      setRoles(Array.isArray(data) ? data : []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Ошибка отзыва роли');
    } finally {
      setActionLoading(null);
    }
  };

  return (
    <div className="mb-6 rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] p-5">
      <div className="mb-4 flex items-center gap-2">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-600/10 text-[var(--color-accent)]">
          <UserCog size={18} />
        </div>
        <div>
          <h2 className="text-lg font-medium">Управление ролями</h2>
          <p className="text-xs text-[var(--color-text-muted)]">
            Выдача и отзыв ролей по username
          </p>
        </div>
      </div>

      <form onSubmit={handleLookup} className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-end">
        <div className="flex-1">
          <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
            Username
          </label>
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)] focus:ring-1 focus:ring-[var(--color-accent)]/30"
            placeholder="username"
            required
          />
        </div>
        <div className="sm:w-40">
          <label className="mb-1.5 block text-sm text-[var(--color-text-muted)]">
            Роль
          </label>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="w-full rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-3.5 py-2.5 text-sm outline-none focus:border-[var(--color-accent)]"
          >
            {AVAILABLE_ROLES.map((r) => (
              <option key={r} value={r}>
                {r}
              </option>
            ))}
          </select>
        </div>
        <motion.button
          type="submit"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          disabled={loading}
          className="flex items-center justify-center gap-1.5 rounded-xl border border-[var(--color-border)] bg-[var(--color-background)] px-4 py-2.5 text-sm font-medium hover:border-[var(--color-accent)]/50 disabled:opacity-60"
        >
          {loading ? <Loader2 size={16} className="animate-spin" /> : <Search size={16} />}
          Найти
        </motion.button>
      </form>

      <div className="mb-4 flex flex-wrap gap-2">
        <motion.button
          type="button"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          disabled={!!actionLoading || !username.trim()}
          onClick={handleGrant}
          className="btn-gradient flex items-center gap-1.5 rounded-xl px-4 py-2.5 text-sm font-semibold text-black disabled:opacity-60"
        >
          {actionLoading === 'grant' ? (
            <Loader2 size={16} className="animate-spin" />
          ) : (
            <ShieldPlus size={16} />
          )}
          Выдать роль
        </motion.button>
        <motion.button
          type="button"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
          disabled={!!actionLoading || !username.trim()}
          onClick={handleRevoke}
          className="flex items-center gap-1.5 rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2.5 text-sm font-medium text-red-300 hover:bg-red-500/20 disabled:opacity-60"
        >
          {actionLoading === 'revoke' ? (
            <Loader2 size={16} className="animate-spin" />
          ) : (
            <ShieldMinus size={16} />
          )}
          Отозвать роль
        </motion.button>
      </div>

      {roles !== null && (
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <span className="text-sm text-[var(--color-text-muted)]">Текущие роли:</span>
          {roles.length === 0 ? (
            <span className="text-sm text-[var(--color-text-muted)]">нет</span>
          ) : (
            roles.map((r) => (
              <span
                key={r}
                className="rounded-full border border-[var(--color-accent)]/30 bg-[var(--color-accent)]/10 px-2.5 py-0.5 text-xs font-medium text-[var(--color-accent)]"
              >
                {r}
              </span>
            ))
          )}
        </div>
      )}

      {message && (
        <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-300">
          {message}
        </div>
      )}
      {error && (
        <div className="rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-300">
          {error}
        </div>
      )}
    </div>
  );
}

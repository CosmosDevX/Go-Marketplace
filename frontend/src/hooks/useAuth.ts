import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import type { UserInfo } from '../types';

const TOKEN_KEY = 'access_token';
const USER_KEY = 'auth_user';

function parseUsernameFromToken(token: string): string | null {
  try {
    const payload = token.split('.')[1];
    if (!payload) return null;
    const decoded = JSON.parse(
      atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    );
    return decoded.sub || decoded.username || decoded.preferred_username || null;
  } catch {
    return null;
  }
}

function readStoredUser(): UserInfo | null {
  const saved = localStorage.getItem(USER_KEY);
  if (!saved) return null;
  try {
    return JSON.parse(saved);
  } catch {
    return null;
  }
}

export function useAuth() {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem(TOKEN_KEY)
  );
  const [user, setUser] = useState<UserInfo | null>(() => readStoredUser());
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      const storedToken = localStorage.getItem(TOKEN_KEY);
      const storedUser = readStoredUser();

      // Сразу восстанавливаем сессию из localStorage — без ожидания refresh
      if (storedToken) {
        if (!cancelled) {
          setToken(storedToken);
          if (storedUser) setUser(storedUser);
        }
      } else {
        if (!cancelled) {
          setToken(null);
          setUser(null);
          setReady(true);
        }
        return;
      }

      // Пытаемся обновить токен (cookies через credentials: 'include')
      try {
        const res = await api.refresh();
        if (cancelled) return;

        if (res.access_token) {
          localStorage.setItem(TOKEN_KEY, res.access_token);
          setToken(res.access_token);

          const username =
            storedUser?.username ||
            parseUsernameFromToken(res.access_token) ||
            'user';
          const roles = Array.isArray(res.roles)
            ? res.roles
            : storedUser?.roles ?? [];
          const u: UserInfo = { username, roles };
          localStorage.setItem(USER_KEY, JSON.stringify(u));
          setUser(u);
        }
      } catch {
        // Refresh не удался — НЕ разлогиниваем.
        // Оставляем access_token из localStorage, пока бэкенд сам не ответит 401 на защищённых запросах.
      } finally {
        if (!cancelled) setReady(true);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(
    (accessToken: string, username?: string, roles: string[] = []) => {
      localStorage.setItem(TOKEN_KEY, accessToken);
      setToken(accessToken);

      const name = username || parseUsernameFromToken(accessToken) || 'user';
      const u: UserInfo = { username: name, roles };
      localStorage.setItem(USER_KEY, JSON.stringify(u));
      setUser(u);
    },
    []
  );

  const logout = useCallback(async () => {
    try {
      await api.logout();
    } catch {
      // ignore
    } finally {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      setToken(null);
      setUser(null);
    }
  }, []);

  const roles = user?.roles ?? [];
  const canAccessSellerPanel =
    roles.includes('seller') || roles.includes('admin');

  return {
    token,
    user,
    isAuthenticated: !!token,
    canAccessSellerPanel,
    ready,
    login,
    logout,
  };
}

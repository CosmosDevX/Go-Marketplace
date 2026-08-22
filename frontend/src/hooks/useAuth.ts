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
    const u = JSON.parse(saved);
    // drop legacy user_id if present
    return { username: u.username, roles: u.roles ?? [] };
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
        // keep local session
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
  const isAdmin = roles.includes('admin');
  const isSeller = roles.includes('seller');
  const canAccessSellerPanel = isSeller && !isAdmin;
  const canAccessAdminPanel = isAdmin;

  return {
    token,
    user,
    isAuthenticated: !!token,
    isAdmin,
    isSeller,
    canAccessSellerPanel,
    canAccessAdminPanel,
    ready,
    login,
    logout,
  };
}

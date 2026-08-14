import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import type { UserInfo } from '../types';

const TOKEN_KEY = 'access_token';
const USER_KEY = 'auth_user';

function parseUsernameFromToken(token: string): string | null {
  try {
    // JWT payload is the middle part
    const payload = token.split('.')[1];
    if (!payload) return null;
    const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
    return decoded.sub || decoded.username || decoded.preferred_username || null;
  } catch {
    return null;
  }
}

export function useAuth() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY));
  const [user, setUser] = useState<UserInfo | null>(() => {
    const saved = localStorage.getItem(USER_KEY);
    if (saved) {
      try {
        return JSON.parse(saved);
      } catch {
        return null;
      }
    }
    return null;
  });
  const [ready, setReady] = useState(false);

  // Try refresh on mount if we have a token
  useEffect(() => {
    let cancelled = false;

    (async () => {
      const stored = localStorage.getItem(TOKEN_KEY);
      if (!stored) {
        setReady(true);
        return;
      }

      try {
        const { access_token } = await api.refresh();
        if (cancelled) return;
        localStorage.setItem(TOKEN_KEY, access_token);
        setToken(access_token);

        const username = parseUsernameFromToken(access_token);
        if (username) {
          const u = { username };
          localStorage.setItem(USER_KEY, JSON.stringify(u));
          setUser(u);
        }
      } catch {
        // token invalid — clear
        if (!cancelled) {
          localStorage.removeItem(TOKEN_KEY);
          localStorage.removeItem(USER_KEY);
          setToken(null);
          setUser(null);
        }
      } finally {
        if (!cancelled) setReady(true);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback((accessToken: string, username?: string) => {
    localStorage.setItem(TOKEN_KEY, accessToken);
    setToken(accessToken);

    const name = username || parseUsernameFromToken(accessToken) || 'user';
    const u = { username: name };
    localStorage.setItem(USER_KEY, JSON.stringify(u));
    setUser(u);
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.logout();
    } catch {
      // ignore network errors on logout
    } finally {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      setToken(null);
      setUser(null);
    }
  }, []);

  return {
    token,
    user,
    isAuthenticated: !!token,
    ready,
    login,
    logout,
  };
}

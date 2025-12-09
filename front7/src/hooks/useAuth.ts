import { useEffect, useState } from 'react';

const AUTH_TOKEN_KEY = 'auth_token';
const AUTH_USER_KEY = 'auth_user';
const AUTH_USER_ID_KEY = 'auth_user_id';

export function useAuth() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [username, setUsername] = useState<string | null>(null);
  const [userID, setUserID] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem(AUTH_TOKEN_KEY);
    const user = localStorage.getItem(AUTH_USER_KEY);
    const userId = localStorage.getItem(AUTH_USER_ID_KEY);
    if (token && user) {
      setIsAuthenticated(true);
      setUsername(user);
      setUserID(userId);
    }
    setLoading(false);
  }, []);

  const login = (token: string, userId: string, username: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
    localStorage.setItem(AUTH_USER_KEY, username);
    localStorage.setItem(AUTH_USER_ID_KEY, userId);
    setIsAuthenticated(true);
    setUsername(username);
    setUserID(userId);
  };

  const logout = () => {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(AUTH_USER_KEY);
    localStorage.removeItem(AUTH_USER_ID_KEY);
    setIsAuthenticated(false);
    setUsername(null);
    setUserID(null);
  };

  return {
    isAuthenticated,
    username,
    userID,
    loading,
    login,
    logout,
  };
}


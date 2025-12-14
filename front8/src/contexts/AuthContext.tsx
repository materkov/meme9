'use client';

import React, { createContext, useContext, useLayoutEffect, useState, ReactNode } from 'react';
import { getAuthToken, LoginResponse } from '@/lib/api';
import { setAuthTokenCookie, removeAuthTokenCookie } from '@/lib/auth-client';

interface AuthContextType {
  isAuthenticated: boolean;
  username: string | null;
  userId: string | null;
  loading: boolean;
  login: (response: LoginResponse) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  // Always start with the same initial state on both server and client
  // This prevents hydration mismatches
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [username, setUsername] = useState<string | null>(null);
  const [userId, setUserId] = useState<string | null>(null);

  // Use useLayoutEffect to update state synchronously before browser paint
  // This ensures "Welcome {name}" appears immediately without visible delay
  useLayoutEffect(() => {
    const token = getAuthToken();
    const storedUsername = localStorage.getItem('auth_username');
    const storedUserId = localStorage.getItem('auth_user_id');

    if (token && storedUsername && storedUserId) {
      setIsAuthenticated(true);
      setUsername(storedUsername);
      setUserId(storedUserId);
    }
  }, []);

  const login = (response: LoginResponse) => {
    setAuthTokenCookie(response.token);
    localStorage.setItem('auth_username', response.username);
    localStorage.setItem('auth_user_id', response.user_id);
    setIsAuthenticated(true);
    setUsername(response.username);
    setUserId(response.user_id);
  };

  const logout = () => {
    removeAuthTokenCookie();
    setIsAuthenticated(false);
    setUsername(null);
    setUserId(null);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        username,
        userId,
        loading: false,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}

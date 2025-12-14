'use client';

import React, { createContext, useContext, useLayoutEffect, useState, ReactNode } from 'react';
import { getAuthToken, getAuthUsername, getAuthUserId } from '@/lib/authHelpers';
import type { LoginResponse } from '@/schema/auth';
import { setAuthTokenCookie, removeAuthTokenCookie } from '@/lib/authHelpers';

interface AuthContextType {
  isAuthenticated: boolean;
  username: string | null;
  userId: string | null;
  loading: boolean;
  login: (response: LoginResponse) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface InitialAuth {
  isAuthenticated: boolean;
  username: string | null;
  userId: string | null;
}

export function AuthProvider({ 
  children,
  initialAuth
}: { 
  children: ReactNode;
  initialAuth?: InitialAuth;
}) {
  // Use server-provided initial auth state to prevent hydration mismatch
  const [isAuthenticated, setIsAuthenticated] = useState(initialAuth?.isAuthenticated ?? false);
  const [username, setUsername] = useState<string | null>(initialAuth?.username ?? null);
  const [userId, setUserId] = useState<string | null>(initialAuth?.userId ?? null);


  // Sync with client-side storage on mount (in case it differs from server)
  // Only update if the values actually changed to avoid unnecessary re-renders
  useLayoutEffect(() => {
    const checkAuth = async () => {
      const token = await getAuthToken();
      const storedUsername = await getAuthUsername();
      const storedUserId = await getAuthUserId();

      // Only update if values changed (e.g., user logged in on another tab)
      if (token && storedUsername && storedUserId) {
        setIsAuthenticated(true);
        setUsername(storedUsername);
        setUserId(storedUserId);
      } else {
        // User not authenticated
        setIsAuthenticated(false);
        setUsername(null);
        setUserId(null);
      }
    };
    checkAuth();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const login = (response: LoginResponse) => {
    setAuthTokenCookie(response.token, response.username, response.userId);
    setIsAuthenticated(true);
    setUsername(response.username);
    setUserId(response.userId);
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
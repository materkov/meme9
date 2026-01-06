'use client';

import { useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { AuthClient, ApiError } from '@/lib/api-clients';
import type { LoginResponse } from '@/schema/auth';
import { useAuth } from '@/contexts/AuthContext';

interface AuthPopupProps {
  onClose: () => void;
}

export default function AuthPopup({ onClose }: AuthPopupProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const { login: authLogin } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      let response: LoginResponse;
      if (isLogin) {
        response = await AuthClient.Login({ username, password });
      } else {
        response = await AuthClient.Register({ username, password });
      }

      authLogin(response);
      onClose();
      // Refresh to trigger server-side re-render with new auth state
      // Cookies are set synchronously, so they'll be included in the refresh request
      router.refresh();
    } catch (err) {
      if (err instanceof ApiError && err.err === "username_exists") {
        setError('This username is already taken, please choose another one');
      } else if (err instanceof ApiError && err.err === "invalid_credentials") {
        setError('Invalid username or password');
      } else {
        setError('Something bad happened');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-xl max-w-md w-full p-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-black dark:text-zinc-50">
            Authorization
          </h2>
          <button
            onClick={onClose}
            className="text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="flex mb-6 border-b border-zinc-200 dark:border-zinc-700">
          <button
            onClick={() => {
              setIsLogin(true);
              setError(null);
            }}
            className={`flex-1 py-2 px-4 text-center font-medium transition-colors ${
              isLogin
                ? 'text-black dark:text-zinc-50 border-b-2 border-black dark:border-zinc-50'
                : 'text-zinc-500 dark:text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200'
            }`}
          >
            Login
          </button>

          <button
            onClick={() => {
              setIsLogin(false);
              setError(null);
            }}
            className={`flex-1 py-2 px-4 text-center font-medium transition-colors ${
              !isLogin
                ? 'text-black dark:text-zinc-50 border-b-2 border-black dark:border-zinc-50'
                : 'text-zinc-500 dark:text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200'
            }`}
          >
            Register
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
              <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
            </div>
          )}

          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1"
            >
              Username
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-black dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-zinc-50"
              placeholder="Enter your username"
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-1"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-black dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-zinc-50"
              placeholder="Enter your password"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-black dark:bg-zinc-50 text-white dark:text-black py-2 px-4 rounded-lg font-medium hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Loading...' : isLogin ? 'Login' : 'Register'}
          </button>
        </form>
      </div>
    </div>
  );
}


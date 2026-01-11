'use client';

import { useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { AuthClient, ApiError } from '@/lib/api-clients';
import type { LoginResponse } from '@/schema/auth';
import { useAuth } from '@/contexts/AuthContext';
import styles from './AuthPopup.module.css';

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
    <div className={styles.overlay}>
      <div className={styles.modal}>
        <div className={styles.header}>
          <h2 className={styles.title}>
            Authorization
          </h2>
          <button
            onClick={onClose}
            className={styles.closeButton}
          >
            <svg
              className={styles.closeIcon}
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

        <div className={styles.tabs}>
          <button
            onClick={() => {
              setIsLogin(true);
              setError(null);
            }}
            className={`${styles.tab} ${isLogin ? styles.tabActive : ''}`}
          >
            Login
          </button>

          <button
            onClick={() => {
              setIsLogin(false);
              setError(null);
            }}
            className={`${styles.tab} ${!isLogin ? styles.tabActive : ''}`}
          >
            Register
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          {error && (
            <div className={styles.error}>
              <p className={styles.errorText}>{error}</p>
            </div>
          )}

          <div className={styles.field}>
            <label
              htmlFor="username"
              className={styles.label}
            >
              Username
            </label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className={styles.input}
              placeholder="Enter your username"
            />
          </div>

          <div className={styles.field}>
            <label
              htmlFor="password"
              className={styles.label}
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className={styles.input}
              placeholder="Enter your password"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className={styles.submitButton}
          >
            {loading ? 'Loading...' : isLogin ? 'Login' : 'Register'}
          </button>
        </form>
      </div>
    </div>
  );
}


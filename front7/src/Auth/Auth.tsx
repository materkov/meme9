import { useState } from 'react';
import * as api from '../api/api';
import styles from './Auth.module.css';

interface AuthProps {
  onAuthSuccess: (token: string, userId: string, username: string) => void;
}

export function Auth({ onAuthSuccess }: AuthProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [usernameError, setUsernameError] = useState('');
  const [credentialsError, setCredentialsError] = useState('');
  const [loading, setLoading] = useState(false);

  const clearErrors = () => {
    setError('');
    setUsernameError('');
    setCredentialsError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    clearErrors();
    setLoading(true);

    try {
      const response = isLogin
        ? await api.login({ username, password })
        : await api.register({ username, password });

      onAuthSuccess(response.token, response.user_id, response.username);
    } catch (err) {
      if (err instanceof api.ApiError) {
        if (err.errorCode === 'username_exists') {
          setUsernameError('Username already exists');
        } else if (err.errorCode === 'invalid_credentials' && isLogin) {
          setCredentialsError('Invalid username or password');
        } else {
          setError(err.message);
        }
      } else {
        setError(err instanceof Error ? err.message : 'An error occurred');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.tabs}>
          <button
            className={`${styles.tab} ${isLogin ? styles.active : ''}`}
            onClick={() => {
              setIsLogin(true);
              clearErrors();
            }}
          >
            Login
          </button>
          <button
            className={`${styles.tab} ${!isLogin ? styles.active : ''}`}
            onClick={() => {
              setIsLogin(false);
              clearErrors();
            }}
          >
            Register
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="username">Username</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => {
                setUsername(e.target.value);
                setUsernameError('');
                setCredentialsError('');
              }}
              required
              disabled={loading}
              className={usernameError ? styles.inputError : ''}
            />
            {usernameError && <div className={styles.fieldError}>{usernameError}</div>}
          </div>

          <div className={styles.field}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                setCredentialsError('');
              }}
              required
              disabled={loading}
              className={credentialsError ? styles.inputError : ''}
            />
            {credentialsError && <div className={styles.fieldError}>{credentialsError}</div>}
          </div>

          {error && <div className={styles.error}>{error}</div>}

          <button type="submit" disabled={loading} className={styles.submit}>
            {loading ? 'Loading...' : isLogin ? 'Login' : 'Register'}
          </button>
        </form>
      </div>
    </div>
  );
}


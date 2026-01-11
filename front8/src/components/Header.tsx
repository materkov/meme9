'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import AuthPopup from './AuthPopup';
import styles from './Header.module.css';

export default function Header() {
  const { isAuthenticated } = useAuth();
  const [showAuth, setShowAuth] = useState(false);

  return (
    <>
      <header className={styles.header}>
        <div className={styles.container}>
          <Link href="/feed" className={styles.logo}>
            Meme9
          </Link>
          <div className={styles.authContainer} suppressHydrationWarning>
            {isAuthenticated ? 
              <Authentcated/> :
              <NotAuthentcated onAuthClick={() => setShowAuth(true)}/>
            }
          </div>
        </div>
      </header>

      {showAuth && <AuthPopup onClose={() => setShowAuth(false)} />}
    </>
  );
}

function Authentcated() {
  const router = useRouter();
  const { username, logout } = useAuth();

  const handleLogout = () => {
    logout();
    // Refresh the page to update server-rendered content
    router.refresh();
  };

  return (<>
    <span className={styles.welcomeText}>
      Welcome, <span className={styles.username}>{username}</span>
    </span>
    <button
      onClick={handleLogout}
      className={`${styles.button} ${styles.logoutButton}`}
    >
      Logout
    </button>
  </>)
}

function NotAuthentcated(props: {onAuthClick: () => void}) {
  return (
    <>
      <button
        onClick={props.onAuthClick}
        className={`${styles.button} ${styles.loginButton}`}
      >
        Login / Register
      </button>
    </>
  )
}
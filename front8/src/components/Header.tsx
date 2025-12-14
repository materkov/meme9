'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useIsAuthenticated, useAuthUsername, removeAuthTokenCookie } from '@/lib/authHelpers';
import { useRouter } from 'next/navigation';
import AuthPopup from './AuthPopup';

export default function Header() {
  const isAuthenticated = useIsAuthenticated();
  const [showAuth, setShowAuth] = useState(false);

  return (
    <>
      <header className="bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <Link href="/feed" className="text-2xl font-bold text-black dark:text-zinc-50 hover:opacity-80 transition-opacity no-underline">
              Meme9
            </Link>
            <div className="flex items-center gap-4 min-w-[140px] justify-end" suppressHydrationWarning>
              {isAuthenticated ? 
                <Authentcated/> :
                <NotAuthentcated onAuthClick={() => setShowAuth(true)}/>
              }
            </div>
          </div>
        </div>
      </header>

      {showAuth && <AuthPopup onClose={() => setShowAuth(false)} />}
    </>
  );
}

function Authentcated() {
  const username = useAuthUsername();
  const router = useRouter();

  const handleLogout = () => {
    removeAuthTokenCookie();
    router.refresh();
  };

  return (<>
    <span className="text-zinc-700 dark:text-zinc-300 whitespace-nowrap">
      Welcome, <span className="font-semibold">{username}</span>
    </span>
    <button
      onClick={handleLogout}
      className="px-4 py-2 bg-zinc-100 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 rounded-lg hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors whitespace-nowrap"
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
        className="px-4 py-2 bg-black dark:bg-zinc-50 text-white dark:text-black rounded-lg hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-colors whitespace-nowrap"
      >
        Login / Register
      </button>
    </>
  )
}
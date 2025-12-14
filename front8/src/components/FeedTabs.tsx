'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function FeedTabs() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { isAuthenticated } = useAuth();
  const currentFeedType = (searchParams.get('feed') || 'all') as 'all' | 'subscriptions';

  const setFeedType = (type: 'all' | 'subscriptions') => {
    const params = new URLSearchParams(searchParams.toString());
    if (type === 'all') {
      params.delete('feed');
    } else {
      params.set('feed', type);
    }
    const newUrl = params.toString() ? `/?${params.toString()}` : '/';
    router.replace(newUrl, { scroll: false });
  };

  return (
    <div className="flex border border-zinc-200 dark:border-zinc-800 rounded-lg overflow-hidden">
      <button
        onClick={() => setFeedType('all')}
        className={`px-4 py-2 text-sm font-medium transition-colors ${
          currentFeedType === 'all'
            ? 'bg-black dark:bg-zinc-50 text-white dark:text-black'
            : 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800'
        }`}
      >
        Global
      </button>
      <button
        onClick={() => setFeedType('subscriptions')}
        disabled={!isAuthenticated}
        className={`px-4 py-2 text-sm font-medium transition-colors ${
          currentFeedType === 'subscriptions'
            ? 'bg-black dark:bg-zinc-50 text-white dark:text-black'
            : 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800'
        } ${!isAuthenticated ? 'opacity-50 cursor-not-allowed' : ''}`}
        title={!isAuthenticated ? 'Login to view subscriptions' : ''}
      >
        Subscriptions
      </button>
    </div>
  );
}


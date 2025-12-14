'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useIsAuthenticated } from '@/lib/authHelpers';
import { FeedType } from '@/schema/posts';

export default function FeedTabs() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isAuthenticated = useIsAuthenticated();

  const typeParam = 'feed';
  
  // Convert URL param to enum value
  const feedParam = searchParams.get(typeParam);
  const currentFeedType = feedParam === 'subscriptions' 
    ? FeedType.SUBSCRIPTIONS 
    : FeedType.ALL;

  const setFeedType = (type: FeedType) => {
    const params = new URLSearchParams(searchParams.toString());
    if (type === FeedType.ALL) {
      params.delete(typeParam);
    } else {
      params.set(typeParam, 'subscriptions');
    }

    const newUrl = params.toString() ? `/feed?${params.toString()}` : '/feed';
    router.replace(newUrl, { scroll: false });
  };

  const classes = 'px-4 py-2 text-sm font-medium transition-colors';
  const activeClasses = 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800';
  const inactiveClasses = 'bg-black dark:bg-zinc-50 text-white dark:text-black';

  return (
    <div className="flex border border-zinc-200 dark:border-zinc-800 rounded-lg overflow-hidden">
      <button
        onClick={() => setFeedType(FeedType.ALL)}
        className={`${classes} ${currentFeedType === FeedType.ALL ? inactiveClasses : activeClasses}`}
      >
        Global
      </button>

      <button
        onClick={() => setFeedType(FeedType.SUBSCRIPTIONS)}
        disabled={!isAuthenticated}
        className={`${classes} ${currentFeedType === FeedType.SUBSCRIPTIONS ? inactiveClasses : activeClasses} ${!isAuthenticated ? 'opacity-50 cursor-not-allowed' : ''}`}
        title={!isAuthenticated ? 'Login to view subscriptions' : ''}
      >
        Subscriptions
      </button>
    </div>
  );
}


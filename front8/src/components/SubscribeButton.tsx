'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { SubscriptionsClient } from '@/lib/api-clients';
import { useAuth } from '@/contexts/AuthContext';

interface SubscribeButtonProps {
  userId: string;
  initialSubscribed: boolean;
}

export default function SubscribeButton({ userId, initialSubscribed }: SubscribeButtonProps) {
  const router = useRouter();
  const { userId: viewerId } = useAuth();
  const [isSubscribed, setIsSubscribed] = useState(initialSubscribed);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubscribe = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (!viewerId) {
      setError('Please login to subscribe');
      return;
    }
    
    setLoading(true);
    setError(null);

    try {
      if (isSubscribed) {
        await SubscriptionsClient.Unsubscribe({ userId });
      } else {
        await SubscriptionsClient.Subscribe({ userId });
      }

      const status = await SubscriptionsClient.GetStatus({ userId });
      setIsSubscribed(status.subscribed);
      
      // Refresh to update server-rendered content
      router.refresh();
    } catch (err) {
      setError('Failed to update subscription');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-end">
      <button
        type="button"
        onClick={handleSubscribe}
        disabled={loading}
        className="px-6 py-2 rounded-lg font-medium transition-all duration-300 bg-black dark:bg-zinc-50 text-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubscribed ? 'Unsubscribe' : 'Subscribe'}
      </button>
      {error && (
        <div className="mt-2 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-2">
          <p className="text-xs text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}
    </div>
  );
}


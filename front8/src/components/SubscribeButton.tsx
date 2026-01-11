'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { SubscriptionsClient } from '@/lib/api-clients';
import { useAuth } from '@/contexts/AuthContext';
import styles from './SubscribeButton.module.css';

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
    <div className={styles.container}>
      <button
        type="button"
        onClick={handleSubscribe}
        disabled={loading}
        className={styles.button}
      >
        {isSubscribed ? 'Unsubscribe' : 'Subscribe'}
      </button>
      {error && (
        <div className={styles.error}>
          <p className={styles.errorText}>{error}</p>
        </div>
      )}
    </div>
  );
}


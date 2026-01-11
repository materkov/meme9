import { FeedType } from '@/schema/posts';
import { getAuthToken } from '@/lib/authHelpers';
import FeedTabButton from './FeedTabButton';
import styles from './FeedTabs.module.css';

interface FeedTabsProps {
  searchParams: Promise<{ feed?: string }>;
}

export default async function FeedTabs({ searchParams }: FeedTabsProps) {
  const resolvedSearchParams = await searchParams;
  const feedParam = resolvedSearchParams?.feed;
  const currentFeedType = feedParam === 'subscriptions' 
    ? FeedType.SUBSCRIPTIONS 
    : FeedType.ALL;

  // Check authentication on server
  const token = await getAuthToken();
  const isAuthenticated = !!token;

  return (
    <div className={styles.tabs}>
      <FeedTabButton
        type={FeedType.ALL}
        label="Global"
        isActive={currentFeedType === FeedType.ALL}
      />
      <FeedTabButton
        type={FeedType.SUBSCRIPTIONS}
        label="Subscriptions"
        isActive={currentFeedType === FeedType.SUBSCRIPTIONS}
        disabled={!isAuthenticated}
        {...(!isAuthenticated && { disabledTitle: 'Login to view subscriptions' })}
      />
    </div>
  );
}


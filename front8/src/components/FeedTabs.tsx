import { FeedType } from '@/schema/posts';
import { getAuthToken } from '@/lib/authHelpers';
import FeedTabButton from './FeedTabButton';

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
    <div className="flex border border-zinc-200 dark:border-zinc-800 rounded-lg overflow-hidden">
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
        disabledTitle={!isAuthenticated ? 'Login to view subscriptions' : undefined}
      />
    </div>
  );
}


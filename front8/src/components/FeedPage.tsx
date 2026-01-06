import { Suspense } from 'react';
import { PostsClient } from '@/lib/api-clients';
import type { Post } from '@/schema/posts';
import { FeedType } from '@/schema/posts';
import { getAuthToken } from '@/lib/authHelpers';
import FeedTabs from '@/components/FeedTabs';
import Composer from '@/components/Composer';
import PostCard from './PostCard';

interface FeedPageProps {
  searchParams: Promise<{ feed?: string }>;
}

export default async function FeedPage({ searchParams }: FeedPageProps) {
  let posts: Post[] = [];
  let error: string | null = null;
  
  const resolvedSearchParams = await searchParams;
  const feedParam = resolvedSearchParams?.feed;
  const feedType = feedParam === 'subscriptions' ? FeedType.SUBSCRIPTIONS : FeedType.ALL;

  if (feedType === FeedType.SUBSCRIPTIONS) {
    const token = await getAuthToken();
    if (!token) {
      error = 'Please login, to view subscriptions feed';
    }
  }

  if (!error) {
    try {
      const response = await PostsClient.GetFeed({ type: feedType });
      posts = response.posts || [];
    } catch (err) {
      error = 'Failed to load feed';
    }
  }

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-black dark:text-zinc-50">Feed</h1>
        <Suspense fallback={<div className="w-[140px]" />}>
          <FeedTabs searchParams={searchParams} />
        </Suspense>
      </div>
      
      <Composer />

      {error ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-red-600 dark:text-red-400">Error: {error}</p>
        </div>
      ) : feedType === FeedType.SUBSCRIPTIONS && posts.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-zinc-600 dark:text-zinc-400">
            No posts from your subscriptions. Try following some users!
          </p>
        </div>
      ) : posts.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-zinc-600 dark:text-zinc-400">No posts found</p>
        </div>
      ) : (
        <>
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </>
      )}
    </div>
  );
}

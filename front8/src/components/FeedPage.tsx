import { Suspense } from 'react';
import { PostsClient } from '@/lib/api-clients';
import type { Post } from '@/schema/posts';
import { FeedType } from '@/schema/posts';
import FeedTabs from '@/components/FeedTabs';
import PostForm from '@/components/PostForm';
import PostCard from './PostCard';

interface FeedPageProps {
  searchParams: Promise<{ feed?: string }> | { feed?: string };
}

export default async function FeedPage({ searchParams }: FeedPageProps) {
  let posts: Post[] = [];
  let error: string | null = null;
  
  // Handle searchParams as either Promise or object (Next.js 15+ compatibility)
  const resolvedSearchParams = searchParams instanceof Promise 
    ? await searchParams 
    : searchParams;
  
  // Determine feed type from searchParams, default to ALL
  // Let the backend handle authentication - if user requests subscriptions but isn't authenticated,
  // the backend will return an error which we'll handle
  const feedParam = resolvedSearchParams?.feed;
  const feedType = feedParam === 'subscriptions' ? FeedType.SUBSCRIPTIONS : FeedType.ALL;

  try {
    // Standard client automatically reads token from cookies on server, localStorage on client
    const response = await PostsClient.GetFeed({ type: feedType });
    posts = response.posts || [];
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load feed';
    // If it's an authentication error for subscriptions, show a helpful message
    if (feedType === FeedType.SUBSCRIPTIONS && error.includes('authentication')) {
      error = 'Please login to view subscriptions feed';
    }
  }

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-black dark:text-zinc-50">Feed</h1>
        <Suspense fallback={<div className="w-[140px]" />}>
          <FeedTabs />
        </Suspense>
      </div>
      
      <PostForm />

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


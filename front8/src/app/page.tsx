import { Suspense } from 'react';
import { getFeed, FeedPost } from '@/lib/api';
import { getServerAuthToken, isServerAuthenticated } from '@/lib/auth-server';
import FeedPosts from '@/components/FeedPosts';
import FeedTabs from '@/components/FeedTabs';
import PostForm from '@/components/PostForm';

interface HomeProps {
  searchParams: { feed?: string };
}

export default async function Home({ searchParams }: HomeProps) {
  let posts: FeedPost[] = [];
  let error: string | null = null;
  
  // Get token from cookies
  const token = await getServerAuthToken();
  const isAuthenticated = await isServerAuthenticated();
  
  // Determine feed type from searchParams, default to 'all'
  const feedType = (searchParams.feed === 'subscriptions' && isAuthenticated) 
    ? 'subscriptions' 
    : 'all';

  try {
    // Fetch feed data on server
    posts = await getFeed(feedType, token);
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load feed';
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
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
          ) : feedType === 'subscriptions' && posts.length === 0 ? (
            <div className="flex items-center justify-center py-12">
              <p className="text-zinc-600 dark:text-zinc-400">
                No posts from your subscriptions. Try following some users!
              </p>
            </div>
          ) : (
            <FeedPosts posts={posts} />
          )}
        </div>
      </main>
    </div>
  );
}

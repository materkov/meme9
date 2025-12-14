import Feed from '@/components/Feed';
import { getFeed, FeedPost } from '@/lib/api';
import { getServerAuthToken } from '@/lib/auth-server';

export default async function Home() {
  let posts: FeedPost[] = [];
  let error: string | null = null;
  
  // Get token from cookies (no verification needed here, backend will verify on API call)
  const token = await getServerAuthToken();

  try {
    // Pre-fetch initial feed data for better initial load
    // Feed component will handle client-side switching
    posts = await getFeed('all', token);
  } catch (err) {
    // Don't show error on server, let client handle it
    error = err instanceof Error ? err.message : 'Failed to load feed';
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
        <Feed initialPosts={posts} initialFeedType="all" />
      </main>
    </div>
  );
}

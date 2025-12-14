import { notFound } from 'next/navigation';
import { UsersClient, PostsClient, SubscriptionsClient, ApiError } from '@/lib/api-clients';
import { getAuthToken } from '@/lib/authHelpers';
import type { Post } from '@/schema/posts';
import UserProfile from '@/components/UserProfile';

interface PageProps {
  params: Promise<{id: string}>;
}

export default async function UserPage({ params }: PageProps) {
  const { id } = await params;

  let user;
  let posts: Post[] = [];
  let subscriptionStatus = { subscribed: false };
  let error: string | null = null;

  try {
    // TODO: make this parallel
    user = await UsersClient.Get({ userId: id });
    
    const postsResponse = await PostsClient.GetByUsers({ userId: id });
    posts = postsResponse.posts || [];
    
    const token = await getAuthToken();
    if (token) {
      subscriptionStatus = await SubscriptionsClient.GetStatus({ userId: id });
    }
  } catch (err) {
    if (err instanceof ApiError && err.err === "user_not_found") {
      notFound();
    } else {
      error = 'Failed to load user info';
    }
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-black">
        <main className="container mx-auto px-4 py-8">
          <div className="flex items-center justify-center py-12">
            <p className="text-red-600 dark:text-red-400">Error: {error}</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
        <UserProfile 
          user={user} 
          initialPosts={posts} 
          initialSubscribed={subscriptionStatus.subscribed}
        />
      </main>
    </div>
  );
}

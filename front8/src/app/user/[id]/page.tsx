import { notFound } from 'next/navigation';
import { UsersClient, PostsClient, SubscriptionsClient } from '@/lib/api-clients';
import { UsersClientJSON as UsersClientJSONClass } from '@/schema/users.twirp';
import { PostsClientJSON as PostsClientJSONClass } from '@/schema/posts.twirp';
import { SubscriptionsClientJSON as SubscriptionsClientJSONClass } from '@/schema/subscriptions.twirp';
import { TwirpRpcImpl } from '@/lib/twirp-rpc';
import type { UserPostResponse as UserPost } from '@/schema/posts';
import { getServerAuthToken } from '@/lib/auth-server';
import UserProfile from '@/components/UserProfile';

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function UserPage({ params }: PageProps) {
  const { id } = await params;
  let user;
  let posts: UserPost[] = [];
  let subscriptionStatus = { subscribed: false };
  let error: string | null = null;

  const token = await getServerAuthToken();

  try {
    // Fetch user info
    // Create clients with server-side token
    const rpc = new TwirpRpcImpl(token);
    const usersClient = new UsersClientJSONClass(rpc);
    user = await usersClient.Get({ userId: id });
    
    // Fetch user posts
    const postsClient = new PostsClientJSONClass(rpc);
    const postsResponse = await postsClient.GetByUsers({ userId: id });
    posts = postsResponse.posts || [];
    
    // Fetch subscription status if authenticated
    if (token) {
      try {
        const subscriptionsClient = new SubscriptionsClientJSONClass(rpc);
        subscriptionStatus = await subscriptionsClient.GetStatus({ userId: id });
      } catch (err) {
        // Subscription status is optional, continue without it
        // Silently fail - subscription status will be checked client-side
      }
    }
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load user';
    if (error.includes('not found') || error.includes('NotFound')) {
      notFound();
    }
  }

  if (error || !user) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-black">
        <main className="container mx-auto px-4 py-8">
          <div className="flex items-center justify-center py-12">
            <p className="text-red-600 dark:text-red-400">Error: {error || 'User not found'}</p>
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

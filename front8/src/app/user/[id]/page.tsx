import { notFound } from 'next/navigation';
import { UsersClient, PostsClient, SubscriptionsClient } from '@/lib/api-clients';
import type { Post } from '@/schema/posts';
import UserProfile from '@/components/UserProfile';
import styles from './page.module.css';

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function UserPage({ params }: PageProps) {
  const { id } = await params;
  let user;
  let posts: Post[] = [];
  let subscriptionStatus = { subscribed: false };
  let error: string | null = null;

  try {
    // Standard clients automatically read token from cookies on server, localStorage on client
    // Fetch user info
    user = await UsersClient.Get({ userId: id });
    
    // Fetch user posts
    const postsResponse = await PostsClient.GetByUsers({ userId: id });
    posts = postsResponse.posts || [];
    
    // Fetch subscription status if authenticated
    try {
      subscriptionStatus = await SubscriptionsClient.GetStatus({ userId: id });
    } catch (err) {
      // Subscription status is optional, continue without it
      // Silently fail - subscription status will be checked client-side
    }
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load user';
    if (error.includes('not found') || error.includes('NotFound')) {
      notFound();
    }
  }

  if (error || !user) {
    return (
      <div className={styles.page}>
        <main className={styles.main}>
          <div className={styles.error}>
            <p className={styles.errorText}>Error: {error || 'User not found'}</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <UserProfile 
          user={user} 
          initialPosts={posts} 
          initialSubscribed={subscriptionStatus.subscribed}
        />
      </main>
    </div>
  );
}

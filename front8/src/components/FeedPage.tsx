import { Suspense } from 'react';
import { PostsClient } from '@/lib/api-clients';
import type { Post } from '@/schema/posts';
import { FeedType } from '@/schema/posts';
import { getAuthToken } from '@/lib/authHelpers';
import FeedTabs from '@/components/FeedTabs';
import Composer from '@/components/Composer';
import PostCard from './PostCard';
import styles from './FeedPage.module.css';

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
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>Feed</h1>
        <Suspense fallback={<div style={{ width: '140px' }} />}>
          <FeedTabs searchParams={searchParams} />
        </Suspense>
      </div>
      
      <Composer />

      {error ? (
        <div className={styles.error}>
          <p className={styles.errorText}>Error: {error}</p>
        </div>
      ) : feedType === FeedType.SUBSCRIPTIONS && posts.length === 0 ? (
        <div className={styles.empty}>
          <p className={styles.emptyText}>
            No posts from your subscriptions. Try following some users!
          </p>
        </div>
      ) : posts.length === 0 ? (
        <div className={styles.empty}>
          <p className={styles.emptyText}>No posts found</p>
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

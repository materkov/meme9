'use client';

import { useState } from 'react';
import { SubscriptionsClient } from '@/lib/api-clients';
import type { GetUserResponse as User } from '@/schema/users';
import type { Post } from '@/schema/posts';
import { useAuthUserId } from '@/lib/authHelpers';
import PostCard from './PostCard';

interface UserProfileHeaderProps {
  user: User;
  viewerId: string | null;
  isSubscribed: boolean;
  loading: boolean;
  error: string | null;
  onSubscribe: (e: React.MouseEvent) => void;
}

function UserProfileHeader({ 
  user, 
  viewerId, 
  isSubscribed, 
  loading, 
  error, 
  onSubscribe 
}: UserProfileHeaderProps) {
  const isOwnProfile = viewerId && viewerId === user.id;

  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-black dark:text-zinc-50 mb-2">
            {user.username}
          </h1>
        </div>
        
        {viewerId && !isOwnProfile ? (
          <button
            type="button"
            onClick={onSubscribe}
            disabled={loading}
            className="px-6 py-2 rounded-lg font-medium transition-all duration-300 bg-black dark:bg-zinc-50 text-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isSubscribed ? 'Unsubscribe' : 'Subscribe'}
          </button>
        ) : null}
      </div>

      {error && (
        <div className="mt-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
          <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
        </div>
      )}
    </div>
  );
}

interface UserPostsProps {
  posts: Post[];
}

function UserPosts({ posts }: UserPostsProps) {
  return (
    <div>
      <h2 className="text-2xl font-bold text-black dark:text-zinc-50 mb-4">
        Posts
      </h2>

      {posts.length === 0 ? (
        <div className="flex items-center justify-center py-12 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg">
          <p className="text-zinc-600 dark:text-zinc-400">No posts yet</p>
        </div>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
        </div>
      )}
    </div>
  );
}

interface UserProfileProps {
  user: User;
  initialPosts: Post[];
  initialSubscribed: boolean;
}

export default function UserProfile({ user, initialPosts, initialSubscribed }: UserProfileProps) {
  const viewerId = useAuthUserId();
  const [isSubscribed, setIsSubscribed] = useState(initialSubscribed);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubscribe = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (!viewerId) {
      setError('Please login to subscribe');
      return;
    }
    
    setLoading(true);
    setError(null);

    try {
      if (isSubscribed) {
        await SubscriptionsClient.Unsubscribe({ userId: user.id });
      } else {
        await SubscriptionsClient.Subscribe({ userId: user.id });
      }

      const status = await SubscriptionsClient.GetStatus({ userId: user.id });
      setIsSubscribed(status.subscribed);
    } catch (err) {
      setError('Failed to update subscription');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <UserProfileHeader
        user={user}
        viewerId={viewerId}
        isSubscribed={isSubscribed}
        loading={loading}
        error={error}
        onSubscribe={handleSubscribe}
      />

      <UserPosts posts={initialPosts} />
    </div>
  );
}

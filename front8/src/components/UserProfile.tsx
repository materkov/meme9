'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { SubscriptionsClient } from '@/lib/api-clients';
import type { GetUserResponse as User } from '@/schema/users';
import type { UserPostResponse as UserPost } from '@/schema/posts';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

interface UserProfileProps {
  user: User;
  initialPosts: UserPost[];
  initialSubscribed: boolean;
}

export default function UserProfile({ user, initialPosts, initialSubscribed }: UserProfileProps) {
  const { isAuthenticated, userId } = useAuth();
  const router = useRouter();
  const [posts, setPosts] = useState<UserPost[]>(initialPosts);
  const [isSubscribed, setIsSubscribed] = useState(initialSubscribed);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const isOwnProfile = isAuthenticated && userId === user.id;

  const handleSubscribe = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (!isAuthenticated) {
      setError('Please login to subscribe');
      return;
    }
    
    setLoading(true);
    setError(null);

    try {
      // SubscriptionsClient uses getAuthToken() automatically from localStorage
      if (isSubscribed) {
        const result = await SubscriptionsClient.Unsubscribe({ userId: user.id });
        setIsSubscribed(result.subscribed);
      } else {
        const result = await SubscriptionsClient.Subscribe({ userId: user.id });
        setIsSubscribed(result.subscribed);
      }
      // Refresh subscription status to ensure it's up to date
      try {
        const status = await SubscriptionsClient.GetStatus({ userId: user.id });
        setIsSubscribed(status.subscribed);
      } catch (statusErr) {
        // If status check fails, trust the subscribe/unsubscribe result
        console.warn('Failed to refresh subscription status:', statusErr);
      }
    } catch (err) {
      console.error('Subscribe error:', err);
      setError(err instanceof Error ? err.message : 'Failed to update subscription');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      {/* User Header */}
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-black dark:text-zinc-50 mb-2">
              {user.username}
            </h1>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              User ID: {user.id}
            </p>
          </div>
          
          {!isOwnProfile && isAuthenticated ? (
            <button
              type="button"
              onClick={handleSubscribe}
              disabled={loading}
              className={`px-6 py-2 rounded-lg font-medium transition-colors ${
                isSubscribed
                  ? 'bg-zinc-100 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-200 dark:hover:bg-zinc-700'
                  : 'bg-black dark:bg-zinc-50 text-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-200'
              } disabled:opacity-50 disabled:cursor-not-allowed`}
            >
              {loading ? '...' : isSubscribed ? 'Unsubscribe' : 'Subscribe'}
            </button>
          ) : !isAuthenticated ? (
            <span className="text-sm text-zinc-500 dark:text-zinc-400">
              Login to subscribe
            </span>
          ) : null}
        </div>

        {error && (
          <div className="mt-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
            <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
          </div>
        )}
      </div>

      {/* Posts Section */}
      <div>
        <h2 className="text-2xl font-bold text-black dark:text-zinc-50 mb-4">
          Posts ({posts.length})
        </h2>

        {posts.length === 0 ? (
          <div className="flex items-center justify-center py-12 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg">
            <p className="text-zinc-600 dark:text-zinc-400">No posts yet</p>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <Link
                key={post.id}
                href={`/post/${post.id}`}
                className="block bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <p className="font-semibold text-black dark:text-zinc-50">
                      {post.username || 'Unknown User'}
                    </p>
                    <p className="text-sm text-zinc-500 dark:text-zinc-400">
                      {post.userId}
                    </p>
                  </div>
                  <time className="text-sm text-zinc-500 dark:text-zinc-400">
                    {new Date(post.createdAt).toLocaleDateString('en-US', {
                      year: 'numeric',
                      month: 'short',
                      day: 'numeric',
                      hour: '2-digit',
                      minute: '2-digit',
                    })}
                  </time>
                </div>
                <p className="text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap">
                  {post.text}
                </p>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

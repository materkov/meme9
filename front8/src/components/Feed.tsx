'use client';

import { useState, useEffect } from 'react';
import { PostsClient } from '@/lib/api-clients';
import type { Post } from '@/schema/posts';
import PostForm from './PostForm';
import { useAuth } from '@/contexts/AuthContext';
import PostCard from './PostCard';

interface FeedProps {
  initialPosts?: Post[];
  initialFeedType?: 'all' | 'subscriptions';
}

export default function Feed({ initialPosts = [], initialFeedType = 'all' }: FeedProps) {
  const { isAuthenticated } = useAuth();
  const [feedType, setFeedType] = useState<'all' | 'subscriptions'>(initialFeedType);
  const [posts, setPosts] = useState<Post[]>(initialPosts);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadFeed = async () => {
      setLoading(true);
      setError(null);
      
      try {
        // PostsClient uses getAuthToken() automatically from localStorage
        const response = await PostsClient.GetFeed({ type: feedType });
        const feedPosts = response.posts || [];
        setPosts(feedPosts);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load feed');
        setPosts([]);
      } finally {
        setLoading(false);
      }
    };

    loadFeed();
  }, [feedType]);

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-black dark:text-zinc-50">Feed</h1>
        
        {/* Tabs */}
        <div className="flex border border-zinc-200 dark:border-zinc-800 rounded-lg overflow-hidden">
          <button
            onClick={() => setFeedType('all')}
            className={`px-4 py-2 text-sm font-medium transition-colors ${
              feedType === 'all'
                ? 'bg-black dark:bg-zinc-50 text-white dark:text-black'
                : 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800'
            }`}
          >
            Global
          </button>
          <button
            onClick={() => setFeedType('subscriptions')}
            disabled={!isAuthenticated}
            className={`px-4 py-2 text-sm font-medium transition-colors ${
              feedType === 'subscriptions'
                ? 'bg-black dark:bg-zinc-50 text-white dark:text-black'
                : 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800'
            } ${!isAuthenticated ? 'opacity-50 cursor-not-allowed' : ''}`}
            title={!isAuthenticated ? 'Login to view subscriptions' : ''}
          >
            Subscriptions
          </button>
        </div>
      </div>
      
      <PostForm />

      {loading ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-zinc-600 dark:text-zinc-400">Loading...</p>
        </div>
      ) : error ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-red-600 dark:text-red-400">Error: {error}</p>
        </div>
      ) : posts.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-zinc-600 dark:text-zinc-400">
            {feedType === 'subscriptions' 
              ? 'No posts from your subscriptions. Try following some users!'
              : 'No posts found'}
          </p>
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

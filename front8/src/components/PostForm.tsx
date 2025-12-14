'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { PostsClient } from '@/lib/api-clients';
import { useRouter } from 'next/navigation';

export default function PostForm() {
  const { isAuthenticated } = useAuth();
  const [text, setText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  if (!isAuthenticated) {
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    
    if (!text.trim()) {
      setError('Post text cannot be empty');
      return;
    }

    setLoading(true);

    try {
      // PostsClient uses getAuthToken() automatically from localStorage
      await PostsClient.Publish({ text: text.trim() });
      setText('');
      setError(null);
      
      // Refresh the page to show the new post
        router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to publish post');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm mb-6">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="post-text" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
            What's on your mind?
          </label>
          <textarea
            id="post-text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Share your thoughts..."
            rows={4}
            className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-black dark:text-zinc-50 focus:outline-none focus:ring-2 focus:ring-black dark:focus:ring-zinc-50 resize-none"
            disabled={loading}
          />
          <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
            {text.length} characters
          </p>
        </div>

        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
            <p className="text-sm text-red-600 dark:text-red-400">{error}</p>
          </div>
        )}

        <div className="flex justify-end">
          <button
            type="submit"
            disabled={loading || !text.trim()}
            className="px-6 py-2 bg-black dark:bg-zinc-50 text-white dark:text-black rounded-lg font-medium hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Publishing...' : 'Publish'}
          </button>
        </div>
      </form>
    </div>
  );
}

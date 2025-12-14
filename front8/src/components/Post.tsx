import { Post, User } from '@/lib/api';
import Link from 'next/link';

interface PostProps {
  post: Post;
  user: User | null;
}

export default function PostComponent({ post, user }: PostProps) {
  return (
    <div className="w-full max-w-2xl mx-auto">
      <div className="mb-6">
        <Link
          href="/"
          className="text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors"
        >
          ← Back to Feed
        </Link>
      </div>
      
      <article className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-8 shadow-sm">
        <div className="flex items-start justify-between mb-6">
          <div>
            <Link
              href={`/user/${post.user_id}`}
              className="text-2xl font-semibold text-black dark:text-zinc-50 mb-2 hover:underline block"
            >
              {user ? user.username : 'Unknown User'}
            </Link>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {post.user_id}
            </p>
          </div>
          <time className="text-sm text-zinc-500 dark:text-zinc-400">
            {new Date(post.created_at).toLocaleDateString('en-US', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            })}
          </time>
        </div>
        
        <div className="prose dark:prose-invert max-w-none">
          <p className="text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap text-lg leading-relaxed">
            {post.text}
          </p>
        </div>
      </article>
    </div>
  );
}

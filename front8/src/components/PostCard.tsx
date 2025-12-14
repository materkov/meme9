'use client';

import { useRouter } from 'next/navigation';
import Link from 'next/link';
import type { Post } from '@/schema/posts';
import FormattedDate from './FormattedDate';

interface PostCardProps {
  post: Post;
  clickable?: boolean;
}

export default function PostCard({ post, clickable = true }: PostCardProps) {
  const router = useRouter();
  
  // Get username from Post
  const username = post.userName || 'Unknown User';

  const handleClick = () => {
    if (clickable) {
      router.push(`/post/${post.id}`);
    }
  };

  return (
    <div
      onClick={handleClick}
      className={`bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm ${clickable ? 'hover:shadow-md transition-shadow cursor-pointer' : ''}`}
    >
      <div className="flex items-start justify-between mb-3">
        <div>
          <Link
            href={`/user/${post.userId}`}
            onClick={(e) => e.stopPropagation()}
            className="font-semibold text-black dark:text-zinc-50 hover:underline"
          >
            {username}
          </Link>
          <p className="text-sm text-zinc-500 dark:text-zinc-400">
            {post.userId}
          </p>
        </div>
        <FormattedDate date={post.createdAt}/>
      </div>
      <p className="text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap">
        {post.text}
      </p>
    </div>
  );
}


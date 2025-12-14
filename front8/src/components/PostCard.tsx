'use client';

import { useRouter } from 'next/navigation';
import type { FeedPostResponse as FeedPost } from '@/schema/feed';
import type { GetPostResponse as Post } from '@/schema/posts';
import type { GetUserResponse as User } from '@/schema/users';
import UserLink from './UserLink';
import FormattedDate from './FormattedDate';

interface PostCardProps {
  post: FeedPost | Post;
  user?: User | null;
  clickable?: boolean;
  showBackLink?: boolean;
}

export default function PostCard({ post, user, clickable = true, showBackLink = false }: PostCardProps) {
  const router = useRouter();
  
  // Get username from FeedPost or from user prop
  const username = 'username' in post ? post.username : (user?.username || 'Unknown User');

  const handleClick = () => {
    if (clickable) {
      router.push(`/post/${post.id}`);
    }
  };

  const cardContent = (
    <div
      onClick={handleClick}
      className={`bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg ${showBackLink ? 'p-8' : 'p-6'} shadow-sm ${clickable ? 'hover:shadow-md transition-shadow cursor-pointer' : ''}`}
    >
        <div className="flex items-start justify-between mb-3">
          <div>
            <UserLink
              href={`/user/${post.userId}`}
              className={`${showBackLink ? 'text-2xl font-semibold mb-2 block' : 'font-semibold'} text-black dark:text-zinc-50 hover:underline`}
            >
              {username}
            </UserLink>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              {post.userId}
            </p>
          </div>
          <FormattedDate 
            date={post.createdAt} 
            month={showBackLink ? 'long' : 'short'}
            className="text-sm text-zinc-500 dark:text-zinc-400"
          />
        </div>
        <p className={`text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap ${showBackLink ? 'text-lg leading-relaxed' : ''}`}>
          {post.text}
        </p>
      </div>
  );

  if (showBackLink) {
    return (
      <div className="w-full max-w-2xl mx-auto">
        <div className="mb-6">
          <button
            onClick={() => router.push('/feed')}
            className="text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors"
          >
            ← Back to Feed
          </button>
        </div>
        {cardContent}
      </div>
    );
  }

  return cardContent;
}


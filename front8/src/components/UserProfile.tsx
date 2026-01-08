import type { GetUserResponse as User } from '@/schema/users';
import type { Post } from '@/schema/posts';
import { getAuthUserId } from '@/lib/authHelpers';
import PostCard from './PostCard';
import SubscribeButton from './SubscribeButton';
import AvatarUploadButton from './AvatarUploadButton';

interface UserProfileProps {
  user: User;
  initialPosts: Post[];
  initialSubscribed: boolean;
}

export default async function UserProfile({ user, initialPosts, initialSubscribed }: UserProfileProps) {
  // Get viewer ID on server
  const viewerId = await getAuthUserId();
  const isOwnProfile = viewerId && viewerId === user.id;

  return (
    <div className="w-full max-w-2xl mx-auto space-y-6">
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm">
        <div className="flex items-center justify-between mb-2">
          <div className="flex items-center gap-4">
            {user.avatarUrl ? (
              <img
                src={user.avatarUrl}
                alt={`${user.username}'s avatar`}
                className="w-16 h-16 rounded-full object-cover border-2 border-zinc-200 dark:border-zinc-700"
              />
            ) : (
              <div className="w-16 h-16 rounded-full bg-zinc-200 dark:bg-zinc-800 flex items-center justify-center border-2 border-zinc-200 dark:border-zinc-700">
                <span className="text-2xl font-bold text-zinc-500 dark:text-zinc-400">
                  {user.username.charAt(0).toUpperCase()}
                </span>
              </div>
            )}
            <div>
              <h1 className="text-3xl font-bold text-black dark:text-zinc-50">
                {user.username}
              </h1>
            </div>
          </div>
          
          {viewerId && !isOwnProfile ? (
            <SubscribeButton userId={user.id} initialSubscribed={initialSubscribed} />
          ) : isOwnProfile ? (
            <AvatarUploadButton userId={user.id} />
          ) : null}
        </div>
      </div>

      <div>
        <h2 className="text-2xl font-bold text-black dark:text-zinc-50 mb-4">
          Posts
        </h2>

        {initialPosts.length === 0 ? (
          <div className="flex items-center justify-center py-12 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg">
            <p className="text-zinc-600 dark:text-zinc-400">No posts yet</p>
          </div>
        ) : (
          <div className="space-y-4">
            {initialPosts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

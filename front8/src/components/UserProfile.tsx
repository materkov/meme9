import type { GetUserResponse as User } from '@/schema/users';
import type { Post } from '@/schema/posts';
import { getAuthUserId } from '@/lib/authHelpers';
import PostCard from './PostCard';
import SubscribeButton from './SubscribeButton';
import AvatarUploadButton from './AvatarUploadButton';
import styles from './UserProfile.module.css';

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
    <div className={styles.container}>
      <div className={styles.profileCard}>
        <div className={styles.profileHeader}>
          <div className={styles.profileInfo}>
            {user.avatarUrl ? (
              <img
                src={user.avatarUrl}
                alt={`${user.username}'s avatar`}
                className={styles.avatar}
              />
            ) : (
              <div className={styles.avatarPlaceholder}>
                <span className={styles.avatarInitial}>
                  {user.username.charAt(0).toUpperCase()}
                </span>
              </div>
            )}
            <div>
              <h1 className={styles.username}>
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

      <div className={styles.postsSection}>
        <h2 className={styles.postsTitle}>
          Posts
        </h2>

        {initialPosts.length === 0 ? (
          <div className={styles.emptyPosts}>
            <p className={styles.emptyPostsText}>No posts yet</p>
          </div>
        ) : (
          <div className={styles.postsList}>
            {initialPosts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

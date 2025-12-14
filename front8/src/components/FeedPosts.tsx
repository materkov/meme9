import { FeedPost } from '@/lib/api';
import PostCard from './PostCard';

interface FeedPostsProps {
  posts: FeedPost[];
}

export default function FeedPosts({ posts }: FeedPostsProps) {
  if (posts.length === 0) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-zinc-600 dark:text-zinc-400">No posts found</p>
      </div>
    );
  }

  return (
    <>
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
    </>
  );
}


import Link from "next/link";
import type { Post } from "@/schema/posts";
import FormattedDate from "./FormattedDate";
import PostCardInteractive from "./PostCardInteractive";
import PostCardClickable from "./PostCardClickable";

interface PostCardProps {
  post: Post;
  clickable?: boolean;
  onLikeChange?: (postId: string, liked: boolean, count: number) => void;
  onDelete?: () => void;
}

export default function PostCard({
  post,
  clickable = true,
  onLikeChange,
  onDelete,
}: PostCardProps) {
  const username = post.userName || "Unknown User";

  return (
    <div className="relative">
      <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm">
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-full bg-zinc-200 dark:bg-zinc-700 overflow-hidden flex items-center justify-center">
              {post.userAvatar ? (
                <img
                  src={post.userAvatar}
                  alt={`${username}'s avatar`}
                  className="w-full h-full object-cover"
                />
              ) : (
                <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
                  {username?.[0]?.toUpperCase() || "?"}
                </span>
              )}
            </div>
            <Link
              href={`/user/${post.userId}`}
              className="font-semibold text-black dark:text-zinc-50 hover:underline"
            >
              {username}
            </Link>
          </div>
          <div className="flex items-center gap-3">
            <FormattedDate date={post.createdAt} />
            <PostCardInteractive post={post} onDelete={onDelete} />
          </div>
        </div>
        <PostCardClickable post={post} clickable={clickable}>
          <p className="text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap">
            {post.text}
          </p>
        </PostCardClickable>
        <PostCardInteractive
          post={post}
          showLikeButton
          onLikeChange={onLikeChange}
        />
      </div>
    </div>
  );
}

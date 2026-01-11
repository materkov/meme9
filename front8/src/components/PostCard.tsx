import Link from "next/link";
import type { Post } from "@/schema/posts";
import FormattedDate from "./FormattedDate";
import PostCardInteractive from "./PostCardInteractive";
import PostCardClickable from "./PostCardClickable";
import styles from "./PostCard.module.css";

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
    <div className={styles.card}>
      <div className={styles.cardInner}>
        <div className={styles.header}>
          <div className={styles.userInfo}>
            <div className={styles.avatar}>
              {post.userAvatar ? (
                <img
                  src={post.userAvatar}
                  alt={`${username}'s avatar`}
                  className={styles.avatarImage}
                />
              ) : (
                <span className={styles.avatarInitial}>
                  {username?.[0]?.toUpperCase() || "?"}
                </span>
              )}
            </div>
            <Link
              href={`/user/${post.userId}`}
              className={styles.username}
            >
              {username}
            </Link>
          </div>
          <div className={styles.meta}>
            <FormattedDate date={post.createdAt} />
            <PostCardInteractive 
              post={post} 
              {...(onDelete && { onDelete })}
            />
          </div>
        </div>
        <PostCardClickable post={post} clickable={clickable}>
          <p className={styles.text}>
            {post.text}
          </p>
        </PostCardClickable>
        <PostCardInteractive
          post={post}
          showLikeButton
          {...(onLikeChange && { onLikeChange })}
        />
      </div>
    </div>
  );
}

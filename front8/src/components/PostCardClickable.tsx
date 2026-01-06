'use client';

import { useRouter } from 'next/navigation';
import type { Post } from '@/schema/posts';

interface PostCardClickableProps {
  post: Post;
  clickable: boolean;
  children: React.ReactNode;
}

export default function PostCardClickable({ post, clickable, children }: PostCardClickableProps) {
  const router = useRouter();
  
  const handleClick = () => {
    if (clickable) {
      router.push(`/post/${post.id}`);
    }
  };

  if (!clickable) {
    return <>{children}</>;
  }

  return (
    <div 
      onClick={handleClick}
      className="hover:shadow-md transition-shadow cursor-pointer"
    >
      {children}
    </div>
  );
}


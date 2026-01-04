'use client';

import { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import type { Post } from '@/schema/posts';
import FormattedDate from './FormattedDate';
import { useAuth } from '@/contexts/AuthContext';
import { LikesClient } from '@/lib/api-clients';
import { ApiError } from '@/lib/api-clients';
import LikersPopup from './LikersPopup';

interface PostCardProps {
  post: Post;
  clickable?: boolean;
  onLikeChange?: (postId: string, liked: boolean, count: number) => void;
}

export default function PostCard({ post, clickable = true, onLikeChange }: PostCardProps) {
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const [likesCount, setLikesCount] = useState(post.likesCount ?? 0);
  const [isLiked, setIsLiked] = useState(post.isLiked ?? false);
  const [isLiking, setIsLiking] = useState(false);
  const [showLikersPopup, setShowLikersPopup] = useState(false);
  const likeButtonRef = useRef<HTMLButtonElement>(null);
  const popupTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
  
  const username = post.userName || 'Unknown User';

  const handleClick = () => {
    if (clickable) {
      router.push(`/post/${post.id}`);
    }
  };

  const handleLikeClick = async (e: React.MouseEvent) => {
    e.stopPropagation();
    
    if (!isAuthenticated || isLiking) {
      return;
    }

    setIsLiking(true);
    try {
      if (isLiked) {
        await LikesClient.Unlike({ postId: post.id });
        const newCount = Math.max(0, likesCount - 1);
        setLikesCount(newCount);
        setIsLiked(false);
        onLikeChange?.(post.id, false, newCount);
      } else {
        await LikesClient.Like({ postId: post.id });
        setLikesCount(likesCount + 1);
        setIsLiked(true);
        onLikeChange?.(post.id, true, likesCount + 1);
      }
    } catch (error) {
      if (error instanceof ApiError && error.err === 'auth_required') {
        // User needs to log in
        return;
      }
      console.error('Failed to toggle like:', error);
    } finally {
      setIsLiking(false);
    }
  };

  const handleLikeButtonHover = () => {
    if (popupTimeoutRef.current) {
      clearTimeout(popupTimeoutRef.current);
      popupTimeoutRef.current = null;
    }
    
    if (likesCount > 0 && likeButtonRef.current) {
      const rect = likeButtonRef.current.getBoundingClientRect();
      setPopupPosition({
        top: rect.bottom + 8,
        left: rect.left,
      });
      setShowLikersPopup(true);
    }
  };

  const handleLikeButtonLeave = () => {
    // Delay closing to allow moving to popup
    popupTimeoutRef.current = setTimeout(() => {
      setShowLikersPopup(false);
    }, 150);
  };

  const handlePopupEnter = () => {
    if (popupTimeoutRef.current) {
      clearTimeout(popupTimeoutRef.current);
      popupTimeoutRef.current = null;
    }
  };

  const handlePopupLeave = () => {
    setShowLikersPopup(false);
  };

  return (
    <div className="relative">
      <div
        className={`bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-6 shadow-sm`}
      >
        <div className="flex items-start justify-between mb-3">
          <div>
            <Link
              href={`/user/${post.userId}`}
              className="font-semibold text-black dark:text-zinc-50 hover:underline"
            >
              {username}
            </Link>
          </div>
          <FormattedDate date={post.createdAt}/>
        </div>
        <p 
          className={`text-zinc-800 dark:text-zinc-200 whitespace-pre-wrap  ${clickable ? 'hover:shadow-md transition-shadow cursor-pointer' : ''}`} 
          onClick={handleClick}
        >
          {post.text}
        </p>
        <div className="mt-4 flex items-center gap-2">
          <button
            ref={likeButtonRef}
            onClick={handleLikeClick}
            onMouseEnter={handleLikeButtonHover}
            onMouseLeave={handleLikeButtonLeave}
            disabled={!isAuthenticated || isLiking}
            className={`flex items-center gap-1 px-3 py-1 rounded-full transition-colors ${
              isLiked
                ? 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
                : 'bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-200 dark:hover:bg-zinc-700'
            } ${!isAuthenticated ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
          >
            <svg
              className="w-5 h-5"
              fill={isLiked ? 'currentColor' : 'none'}
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              />
            </svg>
            <span className="text-sm font-medium">{likesCount}</span>
          </button>
        </div>
      </div>
      
      <LikersPopup
        postId={post.id}
        likesCount={likesCount}
        isVisible={showLikersPopup}
        onClose={() => setShowLikersPopup(false)}
        onMouseEnter={handlePopupEnter}
        onMouseLeave={handlePopupLeave}
        position={popupPosition}
      />
    </div>
  );
}


'use client';

import { useState, useRef, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import type { Post } from '@/schema/posts';
import { useAuth } from '@/contexts/AuthContext';
import { useSnackbar } from '@/contexts/SnackbarContext';
import { LikesClient, PostsClient } from '@/lib/api-clients';
import { ApiError } from '@/lib/api-clients';
import LikersPopup from './LikersPopup';
import PostMenuPopup from './PostMenuPopup';
import styles from './PostCardInteractive.module.css';

interface PostCardInteractiveProps {
  post: Post;
  showLikeButton?: boolean;
  onLikeChange?: (postId: string, liked: boolean, count: number) => void;
  onDelete?: () => void;
}

export default function PostCardInteractive({ 
  post, 
  showLikeButton = false,
  onLikeChange,
  onDelete 
}: PostCardInteractiveProps) {
  const router = useRouter();
  const { isAuthenticated, userId } = useAuth();
  const { showSnackbar } = useSnackbar();
  
  // Use props as source of truth (SSR), state only for optimistic updates
  const [optimisticLikesCount, setOptimisticLikesCount] = useState<number | null>(null);
  const [optimisticIsLiked, setOptimisticIsLiked] = useState<boolean | null>(null);
  const [isLiking, setIsLiking] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [showLikersPopup, setShowLikersPopup] = useState(false);
  const likeButtonRef = useRef<HTMLButtonElement>(null);
  const popupTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
  
  // Use optimistic state if available, otherwise use SSR props
  const isLiked = optimisticIsLiked !== null ? optimisticIsLiked : (post.isLiked ?? false);
  const likesCount = optimisticLikesCount !== null ? optimisticLikesCount : (post.likesCount ?? 0);
  
  const isOwner = userId === post.userId;

  // Reset optimistic state when props change (e.g., after refresh)
  useEffect(() => {
    setOptimisticIsLiked(null);
    setOptimisticLikesCount(null);
  }, [post.isLiked, post.likesCount, post.id]);

  const handleLikeClick = async (e: React.MouseEvent) => {
    e.stopPropagation();
    
    if (!isAuthenticated || isLiking) {
      return;
    }

    setIsLiking(true);
    
    // Optimistic update
    const currentIsLiked = optimisticIsLiked !== null ? optimisticIsLiked : (post.isLiked ?? false);
    const currentLikesCount = optimisticLikesCount !== null ? optimisticLikesCount : (post.likesCount ?? 0);
    
    if (currentIsLiked) {
      setOptimisticIsLiked(false);
      setOptimisticLikesCount(Math.max(0, currentLikesCount - 1));
    } else {
      setOptimisticIsLiked(true);
      setOptimisticLikesCount(currentLikesCount + 1);
    }
    
    try {
      if (currentIsLiked) {
        await LikesClient.Unlike({ postId: post.id });
        const newCount = Math.max(0, currentLikesCount - 1);
        setOptimisticLikesCount(newCount);
        setOptimisticIsLiked(false);
        onLikeChange?.(post.id, false, newCount);
      } else {
        await LikesClient.Like({ postId: post.id });
        setOptimisticLikesCount(currentLikesCount + 1);
        setOptimisticIsLiked(true);
        onLikeChange?.(post.id, true, currentLikesCount + 1);
      }
    } catch (error) {
      if (error instanceof ApiError && error.err === 'auth_required') {
        // User needs to log in
        return;
      }
      console.error('Failed to toggle like:', error);
      // Revert optimistic update on error
      setOptimisticIsLiked(null);
      setOptimisticLikesCount(null);
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

  const handleDelete = async () => {
    if (!isAuthenticated || isDeleting) {
      return;
    }

    setIsDeleting(true);
    try {
      await PostsClient.Delete({ postId: post.id });
      
      // Show success snackbar
      showSnackbar('Post deleted successfully');
      
      // If we're on a post detail page, redirect to feed after deletion
      if (typeof window !== 'undefined' && window.location.pathname.startsWith('/post/')) {
        router.push('/feed');
      } else {
        // Refresh the page to get updated posts from server (maintains SSR)
        router.refresh();
      }
      
      onDelete?.();
    } catch (error) {
      if (error instanceof ApiError && error.err === 'auth_required') {
        return;
      }
    } finally {
      setIsDeleting(false);
    }
  };

  if (showLikeButton) {
    return (
      <>
        <div className={styles.likeButtonContainer}>
          <button
            ref={likeButtonRef}
            onClick={handleLikeClick}
            onMouseEnter={handleLikeButtonHover}
            onMouseLeave={handleLikeButtonLeave}
            disabled={!isAuthenticated || isLiking}
            className={`${styles.likeButton} ${
              isLiked ? styles.likeButtonLiked : styles.likeButtonNotLiked
            }`}
          >
            <svg
              className={styles.likeIcon}
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
            <span className={styles.likeCount}>{likesCount}</span>
          </button>
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
      </>
    );
  }

  // Menu button (delete)
  if (isOwner) {
    return <PostMenuPopup onDelete={handleDelete} isDeleting={isDeleting} />;
  }

  return null;
}


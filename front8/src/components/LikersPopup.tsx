"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { LikesClient } from "@/lib/api-clients";
import type { GetLikersResponse_Liker } from "@/schema/likes";
import styles from "./LikersPopup.module.css";

interface LikersPopupProps {
  postId: string;
  likesCount: number;
  isVisible: boolean;
  onClose: () => void;
  onMouseEnter?: () => void;
  onMouseLeave?: () => void;
  position: { top: number; left: number };
}

export default function LikersPopup({
  postId,
  likesCount,
  isVisible,
  onClose,
  onMouseEnter,
  onMouseLeave,
  position,
}: LikersPopupProps) {
  const [likers, setLikers] = useState<GetLikersResponse_Liker[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(false);
  const [pageToken, setPageToken] = useState("");
  const popupRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isVisible && likesCount > 0) {
      fetchLikers();
    } else {
      // Reset state when popup is hidden
      setLikers([]);
      setPageToken("");
      setHasMore(false);
      setError(null);
    }
  }, [isVisible, postId]);

  // Close popup when clicking outside or pressing Escape
  useEffect(() => {
    if (!isVisible) return;

    const handleClickOutside = (event: MouseEvent) => {
      if (
        popupRef.current &&
        !popupRef.current.contains(event.target as Node)
      ) {
        onClose();
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };

    // Use a small delay to allow mouse to move to popup
    const timeoutId = setTimeout(() => {
      document.addEventListener("mousedown", handleClickOutside);
    }, 100);

    document.addEventListener("keydown", handleEscape);

    return () => {
      clearTimeout(timeoutId);
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [isVisible, onClose]);

  const fetchLikers = async (loadMore = false) => {
    if (loading) return;

    setLoading(true);
    setError(null);

    try {
      const data = await LikesClient.GetLikers({
        postId,
        pageToken: loadMore ? pageToken : "",
        count: 20,
      });

      if (loadMore) {
        setLikers((prev) => [...prev, ...data.likers]);
      } else {
        setLikers(data.likers);
      }

      setPageToken(data.pageToken || "");
      setHasMore(!!data.pageToken);
    } catch (err) {
      console.error("Failed to fetch likers:", err);
      setError("Failed to load likers");
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    if (hasMore && !loading) {
      fetchLikers(true);
    }
  };

  if (!isVisible || likesCount === 0) {
    return null;
  }

  return (
    <div
      ref={popupRef}
      data-likers-popup
      className={styles.popup}
      style={{
        top: `${position.top}px`,
        left: `${position.left}px`,
      }}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <div className={styles.header}>
        <h3 className={styles.title}>
          {likesCount === 1 ? "1 like" : `${likesCount} likes`}
        </h3>
      </div>

      <div className={styles.content}>
        {loading && likers.length === 0 ? (
          <div className={styles.loading}>
            Loading...
          </div>
        ) : error ? (
          <div className={styles.error}>
            {error}
          </div>
        ) : likers.length === 0 ? (
          <div className={styles.empty}>
            No likers found
          </div>
        ) : (
          <div className={styles.likersList}>
            {likers.map((liker) => (
              <Link
                key={liker.userId}
                href={`/user/${liker.userId}`}
                className={styles.likerLink}
                onClick={onClose}
              >
                <div className={styles.likerInfo}>
                  <div className={styles.likerAvatar}>
                    {liker.userAvatar ? (
                      <img
                        src={liker.userAvatar}
                        alt={`${liker.username || "User"} avatar`}
                        className={styles.likerAvatarImage}
                      />
                    ) : (
                      <span className={styles.likerAvatarInitial}>
                        {liker.username?.[0]?.toUpperCase() || "?"}
                      </span>
                    )}
                  </div>
                  <span className={styles.likerName}>
                    {liker.username || "Unknown User"}
                  </span>
                </div>
              </Link>
            ))}
          </div>
        )}

        {hasMore && (
          <div className={styles.loadMoreContainer}>
            <button
              onClick={loadMore}
              disabled={loading}
              className={styles.loadMoreButton}
            >
              {loading ? "Loading..." : "Load more"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

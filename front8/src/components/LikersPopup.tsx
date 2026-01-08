"use client";

import { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { LikesClient } from "@/lib/api-clients";
import type { GetLikersResponse_Liker } from "@/schema/likes";

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
      className="fixed z-50 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg shadow-lg min-w-[250px] max-w-[350px] max-h-[400px] overflow-hidden flex flex-col"
      style={{
        top: `${position.top}px`,
        left: `${position.left}px`,
      }}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <div className="p-3 border-b border-zinc-200 dark:border-zinc-800">
        <h3 className="text-sm font-semibold text-zinc-900 dark:text-zinc-50">
          {likesCount === 1 ? "1 like" : `${likesCount} likes`}
        </h3>
      </div>

      <div className="overflow-y-auto flex-1">
        {loading && likers.length === 0 ? (
          <div className="p-4 text-center text-sm text-zinc-600 dark:text-zinc-400">
            Loading...
          </div>
        ) : error ? (
          <div className="p-4 text-center text-sm text-red-600 dark:text-red-400">
            {error}
          </div>
        ) : likers.length === 0 ? (
          <div className="p-4 text-center text-sm text-zinc-600 dark:text-zinc-400">
            No likers found
          </div>
        ) : (
          <div className="py-2">
            {likers.map((liker) => (
              <Link
                key={liker.userId}
                href={`/user/${liker.userId}`}
                className="block px-4 py-2 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors"
                onClick={onClose}
              >
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-zinc-200 dark:bg-zinc-700 overflow-hidden flex items-center justify-center">
                    {liker.userAvatar ? (
                      <img
                        src={liker.userAvatar}
                        alt={`${liker.username || "User"} avatar`}
                        className="w-full h-full object-cover"
                      />
                    ) : (
                      <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
                        {liker.username?.[0]?.toUpperCase() || "?"}
                      </span>
                    )}
                  </div>
                  <span className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                    {liker.username || "Unknown User"}
                  </span>
                </div>
              </Link>
            ))}
          </div>
        )}

        {hasMore && (
          <div className="p-2 border-t border-zinc-200 dark:border-zinc-800">
            <button
              onClick={loadMore}
              disabled={loading}
              className="w-full px-3 py-2 text-sm text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? "Loading..." : "Load more"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

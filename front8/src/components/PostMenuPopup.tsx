'use client';

import { useState, useRef } from 'react';

interface PostMenuPopupProps {
  onDelete: () => void;
  isDeleting?: boolean;
}

export default function PostMenuPopup({ onDelete, isDeleting = false }: PostMenuPopupProps) {
  const [showMenuPopup, setShowMenuPopup] = useState(false);
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const menuPopupTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const [menuPopupPosition, setMenuPopupPosition] = useState<{ top: number; right?: number; left?: number }>({ top: 0 });

  const handleMenuButtonHover = () => {
    if (menuPopupTimeoutRef.current) {
      clearTimeout(menuPopupTimeoutRef.current);
      menuPopupTimeoutRef.current = null;
    }
    
    if (menuButtonRef.current) {
      const rect = menuButtonRef.current.getBoundingClientRect();
      setMenuPopupPosition({
        top: rect.bottom - rect.top + 8,
        left: 0,
      });
      setShowMenuPopup(true);
    }
  };

  const handleMenuButtonLeave = () => {
    menuPopupTimeoutRef.current = setTimeout(() => {
      setShowMenuPopup(false);
    }, 150);
  };

  const handleMenuPopupEnter = () => {
    if (menuPopupTimeoutRef.current) {
      clearTimeout(menuPopupTimeoutRef.current);
      menuPopupTimeoutRef.current = null;
    }
  };

  const handleMenuPopupLeave = () => {
    setShowMenuPopup(false);
  };

  const handleDeleteClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    setShowMenuPopup(false);
    onDelete();
  };

  return (
    <div className="relative">
      <button
        ref={menuButtonRef}
        onMouseEnter={handleMenuButtonHover}
        onMouseLeave={handleMenuButtonLeave}
        className="p-1 text-zinc-500 dark:text-zinc-400 hover:text-zinc-700 dark:hover:text-zinc-200 transition-colors cursor-pointer"
        title="More options"
      >
        <svg
          className="w-5 h-5"
          fill="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z" />
        </svg>
      </button>
      
      {showMenuPopup && (
        <div
          className="absolute z-50 bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg shadow-lg py-1 min-w-[120px]"
          style={{
            top: menuPopupPosition.top,
            ...(menuPopupPosition.right !== undefined ? { right: menuPopupPosition.right } : {}),
            ...(menuPopupPosition.left !== undefined ? { left: menuPopupPosition.left } : {}),
          }}
          onMouseEnter={handleMenuPopupEnter}
          onMouseLeave={handleMenuPopupLeave}
        >
          <button
            onClick={handleDeleteClick}
            disabled={isDeleting}
            className="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-zinc-100 dark:hover:bg-zinc-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer flex items-center gap-2"
          >
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
            Delete
          </button>
        </div>
      )}
    </div>
  );
}


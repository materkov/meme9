'use client';

import { useState, useRef } from 'react';
import styles from './PostMenuPopup.module.css';

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
    <div className={styles.container}>
      <button
        ref={menuButtonRef}
        onMouseEnter={handleMenuButtonHover}
        onMouseLeave={handleMenuButtonLeave}
        className={styles.menuButton}
        title="More options"
      >
        <svg
          className={styles.menuIcon}
          fill="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z" />
        </svg>
      </button>
      
      {showMenuPopup && (
        <div
          className={styles.menu}
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
            className={styles.deleteButton}
          >
            <svg
              className={styles.deleteIcon}
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


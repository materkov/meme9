'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';
import { FeedType } from '@/schema/posts';
import styles from './FeedTabButton.module.css';

interface FeedTabButtonProps {
  type: FeedType;
  label: string;
  isActive: boolean;
  disabled?: boolean;
  disabledTitle?: string;
}

export default function FeedTabButton({ 
  type, 
  label, 
  isActive, 
  disabled = false,
  disabledTitle 
}: FeedTabButtonProps) {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  
  const typeParam = 'feed';
  const params = new URLSearchParams(searchParams.toString());
  
  let href = '/feed';
  if (type === FeedType.SUBSCRIPTIONS) {
    params.set(typeParam, 'subscriptions');
    href = `/feed?${params.toString()}`;
  } else {
    params.delete(typeParam);
    const newParams = params.toString();
    href = newParams ? `/feed?${newParams}` : '/feed';
  }

  if (disabled) {
    return (
      <button
        disabled
        className={`${styles.button} ${isActive ? styles.buttonInactive : styles.buttonActive}`}
        title={disabledTitle}
      >
        {label}
      </button>
    );
  }

  return (
    <Link
      href={href}
      className={`${styles.button} ${styles.link} ${isActive ? styles.buttonInactive : styles.buttonActive}`}
    >
      {label}
    </Link>
  );
}


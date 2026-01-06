'use client';

import Link from 'next/link';
import { usePathname, useSearchParams } from 'next/navigation';
import { FeedType } from '@/schema/posts';

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

  const classes = 'px-4 py-2 text-sm font-medium transition-colors';
  const activeClasses = 'bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800';
  const inactiveClasses = 'bg-black dark:bg-zinc-50 text-white dark:text-black';

  if (disabled) {
    return (
      <button
        disabled
        className={`${classes} ${isActive ? inactiveClasses : activeClasses} opacity-50 cursor-not-allowed`}
        title={disabledTitle}
      >
        {label}
      </button>
    );
  }

  return (
    <Link
      href={href}
      className={`${classes} ${isActive ? inactiveClasses : activeClasses} block text-center`}
    >
      {label}
    </Link>
  );
}


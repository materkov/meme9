'use client';

import Link from 'next/link';

interface UserLinkProps {
  href: string;
  children: React.ReactNode;
  className?: string;
}

export default function UserLink({ href, children, className }: UserLinkProps) {
  return (
    <Link
      href={href}
      onClick={(e) => e.stopPropagation()}
      className={className}
    >
      {children}
    </Link>
  );
}


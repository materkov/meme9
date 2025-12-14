'use client';

import { useState, useEffect } from 'react';

interface FormattedDateProps {
  date: string | Date;
  month?: 'short' | 'long';
  className?: string;
}

export default function FormattedDate({ date, month = 'short', className }: FormattedDateProps) {
  const [formattedDate, setFormattedDate] = useState<string>('');

  useEffect(() => {
    // Format date only on client to prevent hydration mismatch
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    const formatted = dateObj.toLocaleDateString('en-US', {
      year: 'numeric',
      month: month,
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
    setFormattedDate(formatted);
  }, [date, month]);

  return (
    <time className={className} suppressHydrationWarning>
      {formattedDate || ''}
    </time>
  );
}


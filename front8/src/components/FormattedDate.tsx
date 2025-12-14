interface FormattedDateProps {
  date: string;
}

export default function FormattedDate({ date }: FormattedDateProps) {
  const className = "text-sm text-zinc-500 dark:text-zinc-400";

  const dateObj = new Date(date);
  const formattedDate = dateObj.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <time className={className}>
      {formattedDate}
    </time>
  );
}

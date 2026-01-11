import styles from './FormattedDate.module.css';

interface FormattedDateProps {
  date: string;
}

export default function FormattedDate({ date }: FormattedDateProps) {
  const dateObj = new Date(date);
  const formattedDate = dateObj.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <time className={styles.date}>
      {formattedDate}
    </time>
  );
}

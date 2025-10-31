import styles from './Post.module.css';

interface PostProps {
  text: string;
  createdAt: string;
}

export function Post({ text, createdAt }: PostProps) {
  const formattedDate = new Date(createdAt).toLocaleString();

  return (
    <article className={styles.post}>
      <p className={styles.text}>{text}</p>
      <time className={styles.date}>{formattedDate}</time>
    </article>
  );
}


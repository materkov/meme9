import styles from './Post.module.css';

interface PostProps {
  text: string;
  username: string;
  createdAt: string;
}

export function Post({ text, username, createdAt }: PostProps) {
  const formattedDate = new Date(createdAt).toLocaleString();

  return (
    <article className={styles.post}>
      <div className={styles.header}>
        <span className={styles.username}>{username || 'Unknown'}</span>
        <time className={styles.date}>{formattedDate}</time>
      </div>
      <p className={styles.text}>{text}</p>
    </article>
  );
}


import styles from './Post.module.css';

interface PostProps {
  text: string;
  username: string;
  createdAt: string;
  userID?: string;
  onUsernameClick?: (userID: string) => void;
}

export function Post({ text, username, createdAt, userID, onUsernameClick }: PostProps) {
  const formattedDate = new Date(createdAt).toLocaleString();

  const handleUsernameClick = () => {
    if (userID && onUsernameClick) {
      onUsernameClick(userID);
    }
  };

  return (
    <article className={styles.post}>
      <div className={styles.header}>
        <span 
          className={`${styles.username} ${userID && onUsernameClick ? styles.usernameClickable : ''}`}
          onClick={handleUsernameClick}
        >
          {username || 'Unknown'}
        </span>
        <time className={styles.date}>{formattedDate}</time>
      </div>
      <p className={styles.text}>{text}</p>
    </article>
  );
}


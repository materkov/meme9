import { useEffect, useState } from 'react';
import styles from './UserPostsPage.module.css';
import { Post } from '../Post/Post';
import * as api from '../api/api';

interface UserPostsPageProps {
  userID: string;
  onBack: () => void;
}

export function UserPostsPage({ userID, onBack }: UserPostsPageProps) {
  const [posts, setPosts] = useState<api.Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    loadPosts();
  }, [userID]);

  const loadPosts = () => {
    setLoading(true);
    setError(null);
    api.fetchUserPosts(userID)
      .then(data => {
        setPosts(data);
        if (data.length > 0) {
          setUsername(data[0].username);
        }
        setLoading(false);
      })
      .catch(err => {
        console.error('Error fetching user posts:', err);
        setError(err instanceof api.ApiError ? err.errorDetails : 'Failed to load posts');
        setLoading(false);
      });
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <button onClick={onBack} className={styles.backButton}>
          ← Back
        </button>
        <h1 className={styles.title}>
          {username ? `${username}'s Posts` : 'User Posts'}
        </h1>
      </header>
      <main className={styles.main}>
        {loading ? (
          <div className={styles.loading}>Loading posts...</div>
        ) : error ? (
          <div className={styles.error}>{error}</div>
        ) : posts.length === 0 ? (
          <div className={styles.empty}>No posts yet</div>
        ) : (
          <div className={styles.feed}>
            {posts.map(post => (
              <Post 
                key={post.id} 
                text={post.text} 
                username={post.username} 
                createdAt={post.createdAt}
                userID={post.user_id}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}


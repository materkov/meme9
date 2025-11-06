import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './FeedPage.module.css';
import { Post } from '../Post/Post';
import { PostForm } from '../PostForm/PostForm';
import * as api from '../api/api';

interface FeedPageProps {
  username: string;
  onLogout: () => void;
}

export function FeedPage({ username, onLogout }: FeedPageProps) {
  const navigate = useNavigate();
  const [posts, setPosts] = useState<api.Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPosts();
  }, []);

  const loadPosts = () => {
    setLoading(true);
    api.fetchPosts()
      .then(data => {
        setPosts(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Error fetching posts:', err);
        setLoading(false);
      });
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Posts Feed</h1>
        <div className={styles.userInfo}>
          <span className={styles.username}>{username}</span>
          <button onClick={onLogout} className={styles.logout}>
            Logout
          </button>
        </div>
      </header>
      <main className={styles.main}>
        <PostForm onPostCreated={loadPosts} />
        {loading ? (
          <div className={styles.loading}>Loading posts...</div>
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
                onUsernameClick={(id) => navigate(`/users/${id}`)}
              />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}


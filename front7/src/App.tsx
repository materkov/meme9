import { useEffect, useState } from 'react';
import styles from './App.module.css';
import { Post } from './Post/Post';
import { PostForm } from './PostForm/PostForm';
import { Auth } from './Auth/Auth';
import * as api from './api/api';

const AUTH_TOKEN_KEY = 'auth_token';
const AUTH_USER_KEY = 'auth_user';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [username, setUsername] = useState<string | null>(null);
  const [posts, setPosts] = useState<api.Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is already authenticated
    const token = localStorage.getItem(AUTH_TOKEN_KEY);
    const user = localStorage.getItem(AUTH_USER_KEY);
    if (token && user) {
      setIsAuthenticated(true);
      setUsername(user);
      loadPosts();
    } else {
      setLoading(false);
    }
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

  const handleAuthSuccess = (token: string, userId: string, username: string) => {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
    localStorage.setItem(AUTH_USER_KEY, username);
    setIsAuthenticated(true);
    setUsername(username);
    loadPosts();
  };

  const handleLogout = () => {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(AUTH_USER_KEY);
    setIsAuthenticated(false);
    setUsername(null);
    setPosts([]);
  };

  if (!isAuthenticated) {
    return <Auth onAuthSuccess={handleAuthSuccess} />;
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Posts Feed</h1>
        <div className={styles.userInfo}>
          <span className={styles.username}>{username}</span>
          <button onClick={handleLogout} className={styles.logout}>
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
              <Post key={post.id} text={post.text} username={post.username} createdAt={post.createdAt} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}

export default App;


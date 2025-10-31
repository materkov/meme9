import { useEffect, useState } from 'react';
import styles from './App.module.css';
import { Post } from './Post/Post';
import { PostForm } from './PostForm/PostForm';
import * as api from './api/api';

function App() {
  const [posts, setPosts] = useState<api.Post[]>([]);
  const [loading, setLoading] = useState(true);

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

  useEffect(() => {
    loadPosts();
  }, []);

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1>Posts Feed</h1>
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
              <Post key={post.id} text={post.text} createdAt={post.createdAd} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
}

export default App;


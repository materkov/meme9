import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import styles from './UserPostsPage.module.css';
import { Post } from '../Post/Post';
import * as api from '../api/api';
import { useAuth } from '../hooks/useAuth';

export function UserPostsPage() {
  const { id: userID } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { userID: currentUserID } = useAuth();
  const [posts, setPosts] = useState<api.Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [username, setUsername] = useState<string | null>(null);
  const [isSubscribed, setIsSubscribed] = useState<boolean | null>(null);
  const [subscriptionLoading, setSubscriptionLoading] = useState(false);

  if (!userID) {
    return <div>Invalid user ID</div>;
  }

  useEffect(() => {
    loadPosts();
    if (currentUserID && userID && currentUserID !== userID) {
      loadSubscriptionStatus();
    } else {
      setIsSubscribed(null);
    }
  }, [userID, currentUserID]);

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

  const loadSubscriptionStatus = () => {
    if (!userID) return;
    api.getSubscriptionStatus(userID)
      .then(response => {
        setIsSubscribed(response.subscribed);
      })
      .catch(err => {
        console.error('Error fetching subscription status:', err);
        setIsSubscribed(null);
      });
  };

  const handleSubscribe = async () => {
    if (!userID) return;
    setSubscriptionLoading(true);
    try {
      await api.subscribe(userID);
      setIsSubscribed(true);
    } catch (err) {
      console.error('Error subscribing:', err);
      alert(err instanceof api.ApiError ? err.errorDetails : 'Failed to subscribe');
    } finally {
      setSubscriptionLoading(false);
    }
  };

  const handleUnsubscribe = async () => {
    if (!userID) return;
    setSubscriptionLoading(true);
    try {
      await api.unsubscribe(userID);
      setIsSubscribed(false);
    } catch (err) {
      console.error('Error unsubscribing:', err);
      alert(err instanceof api.ApiError ? err.errorDetails : 'Failed to unsubscribe');
    } finally {
      setSubscriptionLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <button onClick={() => navigate('/')} className={styles.backButton}>
          ← Back
        </button>
        <h1 className={styles.title}>
          {username ? `${username}'s Posts` : 'User Posts'}
        </h1>
        {currentUserID && userID && currentUserID !== userID && (
          <div className={styles.subscribeSection}>
            {isSubscribed === null ? (
              <button 
                onClick={loadSubscriptionStatus} 
                className={styles.subscribeButton}
                disabled={subscriptionLoading}
              >
                Check Subscription
              </button>
            ) : isSubscribed ? (
              <button 
                onClick={handleUnsubscribe} 
                className={styles.unsubscribeButton}
                disabled={subscriptionLoading}
              >
                {subscriptionLoading ? 'Unsubscribing...' : 'Unsubscribe'}
              </button>
            ) : (
              <button 
                onClick={handleSubscribe} 
                className={styles.subscribeButton}
                disabled={subscriptionLoading}
              >
                {subscriptionLoading ? 'Subscribing...' : 'Subscribe'}
              </button>
            )}
          </div>
        )}
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


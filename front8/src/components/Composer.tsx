'use client';

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { PostsClient } from '@/lib/api-clients';
import { useRouter } from 'next/navigation';
import styles from './Composer.module.css';

export default function Composer() {
  const { userId: viewerId } = useAuth();
  const [text, setText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  if (!viewerId) {
    return null;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    setLoading(true);

    try {
      await PostsClient.Publish({ text: text.trim() });
      setText('');
      setError(null);
      
      router.refresh();
    } catch (err) {
      setError('Failed to publish post');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <form onSubmit={handleSubmit} className={styles.form}>
        <div>
          <label htmlFor="post-text" className={styles.label}>
            What's on your mind?
          </label>
          <textarea
            id="post-text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Share your thoughts..."
            rows={4}
            className={styles.textarea}
            disabled={loading}
          />
          <p className={styles.charCount}>
            {text.length} characters
          </p>
        </div>

        {error && (
          <div className={styles.error}>
            <p className={styles.errorText}>{error}</p>
          </div>
        )}

        <div className={styles.buttonContainer}>
          <button
            type="submit"
            disabled={loading || !text.trim()}
            className={styles.submitButton}
          >
            {loading ? 'Publishing...' : 'Publish'}
          </button>
        </div>
      </form>
    </div>
  );
}


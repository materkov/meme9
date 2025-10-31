import { useState } from 'react';
import styles from './PostForm.module.css';
import * as api from '../api/api';

interface PostFormProps {
  onPostCreated: () => void;
}

export function PostForm({ onPostCreated }: PostFormProps) {
  const [text, setText] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!text.trim() || submitting) {
      return;
    }

    setSubmitting(true);
    
    try {
      await api.publishPost({ text });
      setText('');
      onPostCreated();
    } catch (error) {
      console.error('Error creating post:', error);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <textarea
        className={styles.textarea}
        value={text}
        onChange={(e) => setText(e.target.value)}
        placeholder="What's on your mind?"
        rows={4}
        disabled={submitting}
      />
      <button 
        type="submit" 
        className={styles.button}
        disabled={submitting || !text.trim()}
      >
        {submitting ? 'Posting...' : 'Post'}
      </button>
    </form>
  );
}


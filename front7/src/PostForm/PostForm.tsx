import { useState } from 'react';
import styles from './PostForm.module.css';
import * as api from '../api/api';

const MAX_TEXT_LENGTH = 1000;

interface PostFormProps {
  onPostCreated: () => void;
}

export function PostForm({ onPostCreated }: PostFormProps) {
  const [text, setText] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const textLength = text.length;
  const isValid = text.trim().length > 0 && textLength <= MAX_TEXT_LENGTH;

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newText = e.target.value;
    setText(newText);
    setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!isValid || submitting) {
      return;
    }

    setSubmitting(true);
    setError(null);
    
    try {
      await api.publishPost({ text });
      setText('');
      setError(null);
      onPostCreated();
    } catch (error: any) {
      const errorMessage = error?.message || 'Failed to create post';
      setError(errorMessage);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <textarea
        className={`${styles.textarea} ${textLength > MAX_TEXT_LENGTH ? styles.textareaError : ''}`}
        value={text}
        onChange={handleTextChange}
        placeholder="What's on your mind?"
        rows={4}
        disabled={submitting}
        maxLength={MAX_TEXT_LENGTH}
      />
      <div className={styles.footer}>
        <div className={styles.meta}>
          {error && <div className={styles.error}>{error}</div>}
          <div className={`${styles.counter} ${textLength > MAX_TEXT_LENGTH ? styles.counterError : ''}`}>
            {textLength} / {MAX_TEXT_LENGTH}
          </div>
        </div>
        <button 
          type="submit" 
          className={styles.button}
          disabled={submitting || !isValid}
        >
          {submitting ? 'Posting...' : 'Post'}
        </button>
      </div>
    </form>
  );
}


import FeedPage from '@/components/FeedPage';
import styles from './page.module.css';

interface FeedProps {
  searchParams: Promise<{ feed?: string }>;
}

export default async function Feed({ searchParams }: FeedProps) {
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <FeedPage searchParams={searchParams} />
      </main>
    </div>
  );
}


import FeedPage from '@/components/FeedPage';

interface FeedProps {
  searchParams: Promise<{ feed?: string }> | { feed?: string };
}

export default async function Feed({ searchParams }: FeedProps) {
  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
        <FeedPage searchParams={searchParams} />
      </main>
    </div>
  );
}


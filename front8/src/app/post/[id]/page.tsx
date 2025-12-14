import { notFound } from 'next/navigation';
import PostCard from '@/components/PostCard';
import { ApiError, PostsClient } from '@/lib/api-clients';

interface PageProps {
  params: Promise<{id: string}>;
}

export default async function PostPage({ params }: PageProps) {
  const { id } = await params;
  let post;
  let error: string | null = null;

  try {
    post = await PostsClient.Get({ postId: id });
  } catch (err) {
    if (err instanceof ApiError && err.err == "post_not_found") {
      notFound();
    } else {
      error = 'Failed to load post';
    }
  }

  if (error || !post) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-black">
        <main className="container mx-auto px-4 py-8">
          <div className="flex items-center justify-center py-12">
            <p className="text-red-600 dark:text-red-400">Error: {error}</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
        <PostCard post={post} clickable={false} />
      </main>
    </div>
  );
}

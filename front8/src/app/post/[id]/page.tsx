import { notFound } from 'next/navigation';
import PostCard from '@/components/PostCard';
import { PostsClient, UsersClient } from '@/lib/api-clients';

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export default async function PostPage({ params }: PageProps) {
  const { id } = await params;
  let post;
  let user: { id: string; username: string } | null = null;
  let error: string | null = null;

  try {
    // Standard clients automatically read token from cookies on server, localStorage on client
    post = await PostsClient.Get({ postId: id });
    
    // Fetch user info to get username
    try {
      user = await UsersClient.Get({ userId: post.userId });
    } catch (err) {
      // If user fetch fails, continue without username
      // Silently fail - post will display without username
    }
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load post';
    if (error.includes('not found') || error.includes('NotFound')) {
      notFound();
    }
  }

  if (error || !post) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-black">
        <main className="container mx-auto px-4 py-8">
          <div className="flex items-center justify-center py-12">
            <p className="text-red-600 dark:text-red-400">Error: {error || 'Post not found'}</p>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black">
      <main className="container mx-auto px-4 py-8">
        <PostCard post={post} user={user} clickable={false} showBackLink={true} />
      </main>
    </div>
  );
}

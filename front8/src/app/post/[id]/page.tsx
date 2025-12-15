import { notFound } from 'next/navigation';
import type { Metadata } from 'next';
import PostCard from '@/components/PostCard';
import { PostsClient } from '@/lib/api-clients';

interface PageProps {
  params: Promise<{
    id: string;
  }>;
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { id } = await params;
  
  try {
    const post = await PostsClient.Get({ postId: id });
    const username = post.userName || 'Unknown User';
    const postText = post.text || '';
    
    // Truncate text for description (OpenGraph recommends 200 chars)
    const description = postText.length > 200 
      ? postText.substring(0, 197) + '...' 
      : postText;
    
    return {
      title: `${username} on Meme9`,
      description: description || `Post by ${username}`,
      openGraph: {
        title: `${username} on Meme9`,
        description: description || `Post by ${username}`,
        type: 'article',
      },
    };
  } catch (err) {
    return {
      title: 'Post on Meme9',
      description: 'View this post on Meme9',
    };
  }
}

export default async function PostPage({ params }: PageProps) {
  const { id } = await params;
  let post;
  let error: string | null = null;

  try {
    // Standard clients automatically read token from cookies on server, localStorage on client
    post = await PostsClient.Get({ postId: id });
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
        <PostCard post={post} clickable={false} />
      </main>
    </div>
  );
}

//const API_BASE_URL = 'http://localhost:8080';
//const API_BASE_URL = 'https://meme.mmaks.me';
const API_BASE_URL = window.API_BASE_URL;

export interface Post {
  id: string;
  text: string;
  createdAd: string;
}

export async function fetchPosts(): Promise<Post[]> {
  const response = await fetch(`${API_BASE_URL}/feed`);
  if (!response.ok) {
    throw new Error('Failed to fetch posts');
  }
  return response.json();
}

export interface PublishPostRequest {
  text: string;
}

export interface PublishPostResponse {
  id: string;
}

export async function publishPost(data: PublishPostRequest): Promise<PublishPostResponse> {
  const response = await fetch(`${API_BASE_URL}/publish`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error('Failed to create post');
  }

  return response.json();
}


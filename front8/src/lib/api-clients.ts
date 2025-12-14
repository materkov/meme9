import { AuthClientJSON } from '@/schema/auth.twirp';
import { PostsClientJSON } from '@/schema/posts.twirp';
import { UsersClientJSON } from '@/schema/users.twirp';
import { SubscriptionsClientJSON } from '@/schema/subscriptions.twirp';
import { getAuthToken } from './authHelpers';

function getApiBaseUrl(): string {
  // Server-side
  if (typeof window === 'undefined') {
    return 'http://localhost:8080';
  }

  // Client-side
  if (window.location.hostname == "localhost") {
    return 'http://localhost:8080';
  } else {
    return 'https://meme2.mmaks.me';
  }
}

class TwirpRpcImpl {
  async request(
    service: string,
    method: string,
    contentType: string,
    data: any
  ): Promise<any> {
    const baseURL = getApiBaseUrl();
    const url = `${baseURL}/twirp/${service}/${method}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    const token = await getAuthToken();
    if (token) {
      headers['Authorization'] = token;
    }
    
    try {
    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ 
        code: 'unknown',
        msg: `HTTP ${response.status}: ${response.statusText}` 
      }));
      throw new Error(errorData.msg || errorData.error || `Request failed: ${response.statusText}`);
    }

    return response.json();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Network error';
      throw new Error(`Failed to connect to server: ${errorMessage}. Please check if the server is running at ${baseURL}`);
    }
  }
}

const rpcImpl = new TwirpRpcImpl();

export const AuthClient = new AuthClientJSON(rpcImpl);
export const PostsClient = new PostsClientJSON(rpcImpl);
export const UsersClient = new UsersClientJSON(rpcImpl);
export const SubscriptionsClient = new SubscriptionsClientJSON(rpcImpl);
